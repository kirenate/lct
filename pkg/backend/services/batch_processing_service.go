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

func (r *Service) UploadDocument(minim, maxim int, name string) (*uuid.UUID, error) {
	uid := uuid.New()

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
	err := r.repository.SaveDocToPg(document)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save document to postgres")
	}
	return &uid, nil
}

func (r *Service) UploadPage(doc *multipart.FileHeader, documentId uuid.UUID, number int) error {
	uid := uuid.New()

	contents, err := doc.Open()
	if err != nil {
		return errors.Wrap(err, "failed to open uploaded page")
	}
	page := &schemas.PageMetadata{
		ID:         uid,
		DocumentId: documentId,
		Thumb:      "",
		Original:   "",
		Number:     number,
		FullText:   nil,
		Attributes: nil,
	}

	err = r.repository.SaveToMinio(doc, uid, contents)
	if err != nil {
		return errors.Wrap(err, "failed to save page to minio")
	}
	err = r.repository.SavePageToPg(page)
	if err != nil {
		return errors.Wrap(err, "failed to save page to postgres")
	}

	return nil
}
