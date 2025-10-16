package services

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"io"
	"main.go/repositories"
	"mime/multipart"
)

type Msg struct {
	Doc        *multipart.FileHeader
	Contents   []byte
	Uid        uuid.UUID
	DocumentId uuid.UUID
}

func (r *Service) SendToQueue(doc *multipart.FileHeader, uid uuid.UUID, documentId uuid.UUID) error {
	file, err := doc.Open()
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}

	contents, err := io.ReadAll(file)
	if err != nil {
		return errors.Wrap(err, "failed to read file")
	}

	msg := &Msg{Doc: doc, Uid: uid, DocumentId: documentId, Contents: contents}
	txt, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal kafka message")
	}

	kafkaMsg := kafka.Message{Value: txt}

	err = r.writer.WriteMessages(context.Background(), kafkaMsg)
	if err != nil {
		return errors.Wrap(err, "failed to write message to kafka")
	}

	return nil
}

func (r *Service) BackgroundConsumer(ctx context.Context) {
	zerolog.Ctx(ctx).Info().Msg("kafka.background.consumer.started")
	for {
		err := r.fetch(ctx)
		if err != nil {
			zerolog.Ctx(ctx).Error().Err(err).Send()
		}
	}
}

func (r *Service) fetch(ctx context.Context) error {
	msg, err := r.reader.FetchMessage(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to fetch msg")
	}
	err = r.consume(&msg)
	if err != nil {
		return errors.Wrap(err, "failed to consume msg")
	}
	err = r.reader.CommitMessages(ctx, msg)
	if err != nil {
		return errors.Wrap(err, "failed to commit msg")
	}
	return nil
}

func (r *Service) consume(kafkaMsg *kafka.Message) error {
	var msg *Msg
	err := json.Unmarshal(kafkaMsg.Value, &msg)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal kafka message")
	}

	text, err := r.ProcessWithML(msg.Doc, msg.Contents)
	if err != nil {
		errs := r.repository.ChangeStatus(msg.DocumentId, repositories.StatusFailed)
		if errs != nil {
			return errors.Wrap(err, errs.Error())
		}
		return errors.Wrap(err, "failed to process image with ML")
	}

	err = r.repository.UpdatePage(text, msg.Uid)
	if err != nil {
		return errors.Wrap(err, "failed to save page to postgres")
	}

	err = r.repository.ChangeStatus(msg.DocumentId, repositories.StatusComplete)
	if err != nil {
		return errors.Wrap(err, "failed to change status")
	}

	return nil
}
