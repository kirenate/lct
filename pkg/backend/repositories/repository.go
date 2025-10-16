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
var StatusComplete = "success"
var StatusFailed = "error"

type Repository struct {
	minio *minio.Client
	db    *gorm.DB
}

func NewRepository(minio *minio.Client, db *gorm.DB) *Repository {
	repository := &Repository{minio: minio, db: db}
	return repository
}

func (r *Repository) CreateBucket(bucketName string) error {
	if ok, err := r.minio.BucketExists(bucketName); ok {
		if err != nil {
			return errors.Wrap(err, "failed to check if bucket exists")
		}
		err = r.minio.SetBucketPolicy(bucketName, settings_utils.Settings.BucketPolicy)
		if err != nil {
			return errors.Wrap(err, "setting bucket policy failed")
		}
		return nil
	}
	err := r.minio.MakeBucket(bucketName, "")
	if err != nil {
		return errors.Wrap(err, "failed to create bucket")
	}
	err = r.minio.SetBucketPolicy(bucketName, settings_utils.Settings.BucketPolicy)
	if err != nil {
		return errors.Wrap(err, "setting bucket policy failed")
	}
	return nil
}

func (r *Repository) GetDocuments(page, pageSize int, order, sorting string) ([]schemas.DocumentMetadata, error) {
	var docs *[]schemas.DocumentMetadata
	err := r.db.Table("document_metadata").
		Order(sorting + " " + order).
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

		err = tx.Table("document_metadata").Where("id", documentId).Delete(&schemas.DocumentMetadata{}, documentId).Error
		if err != nil {
			return errors.Wrap(err, "failed to delete document metadata")
		}

		err = tx.Table("page_metadata").Where("document_id", documentId).Delete(&schemas.PageMetadata{}, documentId).Error
		if err != nil {
			return errors.Wrap(err, "failed to delete pages metadata")
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to delete document")
	}
	return nil
}

func (r *Repository) SaveToMinio(doc *multipart.FileHeader, name string, contents multipart.File) error {

	_, err := r.minio.PutObject(settings_utils.Settings.MinioBucketName, name,
		contents, doc.Size, minio.PutObjectOptions{})
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
	if folderName == "" {
		return nil
	}
	var b []byte
	reader := bytes.NewReader(b)
	_, err := r.minio.PutObject(settings_utils.Settings.MinioBucketName, folderName+"/",
		reader, 0, minio.PutObjectOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to create folder in minio")
	}

	return nil
}

func (r *Repository) SaveAttribute(attr *schemas.Attribute) error {
	err := r.db.Table("attribute").Save(&attr).Error
	if err != nil {
		return errors.Wrap(err, "failed to save attribute")
	}

	return nil
}

func (r *Repository) SaveText(text *[]schemas.Text) error {
	err := r.db.Table("text").Save(&text).Error
	if err != nil {
		return errors.Wrap(err, "failed to save text")
	}

	return nil
}

func (r *Repository) SearchDocuments(page, pageSize int, order, name, status, sorting string) (*[]schemas.DocumentMetadata, error) {
	var docs *[]schemas.DocumentMetadata

	name = "%" + name + "%"
	stmp := r.db.Table("document_metadata").
		Where("name ILIKE (?)", name).
		Order(sorting + " " + order).
		Offset(page * pageSize).
		Limit(pageSize)

	if status != "" {
		stmp = stmp.Where("status", status)
	}

	err := stmp.Find(&docs).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find documents in db")
	}
	return docs, nil
}

func (r *Repository) GetDocumentById(id uuid.UUID) (*schemas.DocumentMetadata, error) {
	var doc *schemas.DocumentMetadata
	err := r.db.Table("document_metadata").Where("id", id).Find(&doc).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to find document by id")
	}

	return doc, nil
}

func (r *Repository) SaveThumbToMinio(img *bytes.Buffer, name string) error {
	_, err := r.minio.PutObject(settings_utils.Settings.MinioBucketName, name,
		img, int64(img.Len()), minio.PutObjectOptions{})

	if err != nil {
		return errors.Wrap(err, "failed to put object in minio")
	}

	return nil
}

func (r *Repository) GetPages(documentId uuid.UUID, page, pageSize int) ([]schemas.PageMetadata, error) {
	var pages *[]schemas.PageMetadata

	err := r.db.Table("page_metadata").
		Order("number ASC").
		Offset(page*pageSize).
		Limit(pageSize).
		Where("document_id", documentId).
		Find(&pages).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pages")
	}

	return *pages, nil
}

func (r *Repository) ChangeStatus(documentId uuid.UUID, status string) error {
	err := r.db.Table("document_metadata").Where("id", documentId).Update("status", status).Error
	if err != nil {
		return errors.Wrap(err, "failed to update status")
	}

	return nil
}

func (r *Repository) CountDocs() (int64, error) {
	var count int64
	err := r.db.Table("document_metadata").Count(&count).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to count documents")
	}

	return count, nil
}

func (r *Repository) CountPages(documentId uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Table("page_metadata").Where("document_id", documentId).Count(&count).Error
	if err != nil {
		return 0, errors.Wrap(err, "failed to count pages")
	}

	return count, nil
}

func (r *Repository) UpdateDocument(doc *schemas.DocumentMetadata, documentId uuid.UUID) error {
	stmp := r.db.Table("document_metadata").Where("id", documentId)
	if doc.Name != "" {
		stmp = stmp.Update("name", doc.Name)
	}
	if doc.Code != "" {
		stmp = stmp.Update("code", doc.Code)
	}
	if doc.Status != "" {
		stmp = stmp.Update("status", doc.Status)
	}
	if doc.Max != 0 {
		stmp = stmp.Update("code", doc.Code)
	}
	if doc.Min != 0 {
		stmp = stmp.Update("min", doc.Min)
	}
	if doc.Progress != 0 {
		stmp = stmp.Update("progress", doc.Progress)
	}
	err := stmp.Error
	if err != nil {
		return errors.Wrap(err, "failed to update document")
	}

	return nil
}

func (r *Repository) UpdatePage(text *string, pageId uuid.UUID) error {
	err := r.db.Table("page_metadata").Where("id", pageId).Update("full_text", *text).Error
	if err != nil {
		return errors.Wrap(err, "failed to update document")
	}

	return nil
}
