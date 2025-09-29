package services

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"main.go/repositories"
	"main.go/schemas"
	"mime/multipart"
	"time"
)

type Service struct {
	repository *repositories.Repository
}

func NewService(repository *repositories.Repository) *Service {
	return &Service{repository: repository}
}

func (r *Service) GetDocuments(page int, pageSize int, order string) ([]schemas.DocumentMetadata, error) {
	docs, err := r.repository.GetDocuments(page, pageSize, order)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get documents")
	}
	return docs, nil
}

func (r *Service) DeleteDocument(documentId uuid.UUID) error {
	err := r.repository.DeleteDocument(documentId)
	if err != nil {
		return errors.Wrap(err, "failed to delete document")
	}

	return nil
}

func (r *Service) UploadFile(doc *multipart.FileHeader, minim, maxim int, name string) error {
	uid := uuid.New()
	contents, err := doc.Open()
	if err != nil {
		return errors.Wrap(err, "failed to open uploaded file")
	}

	now := time.Now().UTC()
	document := &schemas.DocumentMetadata{
		ID:        uid,
		Name:      name,
		Status:    repositories.StatusProcessing,
		Progress:  0,
		Min:       minim,
		Max:       maxim,
		CreatedAt: now,
	}
	err = r.repository.SaveToMinio(doc, uid, contents)
	if err != nil {
		return errors.Wrap(err, "failed to save file to minio")
	}
	var pages []schemas.PageMetadata
	err = r.repository.SaveToPg(document, &pages)

	return nil
}
