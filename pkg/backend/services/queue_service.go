package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"github.com/teadove/teasutils/service_utils/logger_utils"
	"main.go/repositories"
	"time"
)

type Msg struct {
	PageId     uuid.UUID
	DocumentId uuid.UUID
}

func (r *Service) SendToQueue(pageId uuid.UUID, documentId uuid.UUID) error {
	msg := &Msg{PageId: pageId, DocumentId: documentId}
	txt, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal kafka message")
	}

	kafkaMsg := kafka.Message{Value: txt}

	err = r.writer.WriteMessages(context.Background(), kafkaMsg)
	if err != nil {
		return errors.Wrap(err, "failed to write message to kafka")
	}

	zerolog.Ctx(logger_utils.NewLoggedCtx()).Info().Msg("kafka.msg.sent")

	return nil
}

func (r *Service) BackgroundConsumer(ctx context.Context) {
	zerolog.Ctx(ctx).Info().Msg("kafka.background.consumer.started")
	for {
		err := r.fetch(ctx)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Send()
		}

		time.Sleep(300 * time.Millisecond)
	}
}

func (r *Service) fetch(ctx context.Context) error {
	msg, err := r.reader.FetchMessage(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to fetch msg")
	}

	innerCtx := logger_utils.NewLoggedCtx()
	innerCtx = logger_utils.WithValue(ctx, "kafka", fmt.Sprintf("%s:%d:%d", msg.Topic, msg.Partition, msg.Offset))

	zerolog.Ctx(innerCtx).
		Info().
		Msg("kafka.msg.recieved")

	err = r.consume(innerCtx, &msg)
	if err != nil {
		return errors.Wrap(err, "failed to consume msg")
	}
	err = r.reader.CommitMessages(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "failed to commit msg")
	}
	zerolog.Ctx(innerCtx).Info().Msg("kafka.msg.processed")

	return nil
}

func (r *Service) consume(ctx context.Context, kafkaMsg *kafka.Message) error {
	var msg *Msg
	err := json.Unmarshal(kafkaMsg.Value, &msg)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal kafka message")
	}

	obj, err := r.repository.GetObjFromMinio(msg.PageId.String() + ".jpg")
	if err != nil {
		return errors.Wrap(err, "failed to get object from minio")
	}

	text, err := r.ProcessWithML(ctx, obj)
	if err != nil {
		errs := r.repository.ChangePageStatus(msg.DocumentId, repositories.StatusFailed)
		if errs != nil {
			return errors.Wrap(err, errs.Error())
		}
		return errors.Wrap(err, "failed to process image with ML")
	}

	err = r.repository.UpdatePage(text, msg.PageId)
	if err != nil {
		return errors.Wrap(err, "failed to save page to postgres")
	}

	err = r.repository.ChangePageStatus(msg.PageId, repositories.StatusComplete)
	if err != nil {
		return errors.Wrap(err, "failed to change status")
	}

	zerolog.Ctx(ctx).
		Info().
		Msg("msg.consumed")

	return nil
}

func (r *Service) PageLoaderChecker(ctx context.Context, id uuid.UUID) {
	for {
		time.Sleep(10 * time.Second)
		count, err := r.repository.CheckPageLoading(id)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Send()
			break
		}

		if count == 0 {
			err = r.repository.ChangeStatus(id, repositories.StatusComplete)
			if err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Send()
			}
			break
		}
	}
}
