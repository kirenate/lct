package repositories

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"main.go/schemas"
	"main.go/utils/settings_utils"
	"mime/multipart"
)

var StatusProcessing = "processing"
var StatusComplete = "complete"
var StatusFailed = "failed"

type Repository struct {
	minio *minio.Client
	db    *gorm.DB
}

func NewRepository(minio *minio.Client, db *gorm.DB) *Repository {
	return &Repository{minio: minio, db: db}
}

func (r *Repository) GetDocuments(page, pageSize int, order string) ([]schemas.DocumentMetadata, error) {
	var docs *[]schemas.DocumentMetadata
	err := r.db.Table("document_metadata").
		Order("name " + order).
		Offset(page * pageSize).
		Limit(pageSize).
		Find(&docs).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find documents in db")
	}
	return *docs, nil
}

func (r *Repository) DeleteDocument(documentId uuid.UUID) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := r.minio.RemoveObject(settings_utils.Settings.MinioBucketName, documentId.String())
		if err != nil {
			return errors.Wrap(err, "failed to delete document from minio")
		}

		err = tx.Table("document_metadata").Delete("WHERE id = (?)", documentId).Error
		if err != nil {
			return errors.Wrap(err, "failed to delete document metadata")
		}

		err = tx.Table("page_metadata").Delete("WHERE document_id = (?)", documentId).Error
		if err != nil {
			return errors.Wrap(err, "failed to delete pages metadata")
		}

		err = tx.Table("attribute").Delete("WHERE document_id = (?)", documentId).Error
		if err != nil {
			return errors.Wrap(err, "failed to delete attributes")
		}

		err = tx.Table("text").Delete("WHERE document_id = (?)", documentId).Error
		if err != nil {
			return errors.Wrap(err, "failed to delete text")
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to delete document")
	}
	return nil
}

func (r *Repository) SaveToMinio(doc *multipart.FileHeader, uid uuid.UUID, contents multipart.File) error {
	_, err := r.minio.PutObject(settings_utils.Settings.MinioBucketName, uid.String()+".png", contents, doc.Size, minio.PutObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to put object in minio")
	}

	return nil
}

func (r *Repository) SaveDocToPg(document *schemas.DocumentMetadata) error {
	err := r.db.Table("document_metadata").Save(&document).Error
	if err != nil {
		return errors.Wrap(err, "failed to save document metadata")
	}

	return nil
}

func (r *Repository) SavePageToPg(page *schemas.PageMetadata) error {
	err := r.db.Table("page_metadata").Save(&page).Error
	if err != nil {
		return errors.Wrap(err, "failed to save page")
	}

	return nil
}

func (r *Repository) CreateFolder(folderName string) error {
	var b []byte
	reader := bytes.NewReader(b)
	_, err := r.minio.PutObject(settings_utils.Settings.MinioBucketName, folderName+"/", reader, 0, minio.PutObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to create folder in minio")
	}

	return nil
}
