package services

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"main.go/repositories"
	"main.go/schemas"
	"main.go/utils/settings_utils"
	"mime/multipart"
	"strings"
	"time"
)

type Service struct {
	repository *repositories.Repository
}

func NewService(repository *repositories.Repository) (*Service, error) {
	service := &Service{repository: repository}
	err := service.repository.CreateBucket(settings_utils.Settings.MinioBucketName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bucket")
	}
	return service, nil
}

func (r *Service) GetDocuments(page int, pageSize int, sortBy string) ([]schemas.DocumentMetadata, int64, error) {
	order := "DESC"
	sorting, ok := strings.CutPrefix(sortBy, "-")
	if ok {
		order = "ASC"
	}
	if sorting == "createdAt" {
		sorting = "created_at"
	}
	docs, err := r.repository.GetDocuments(page, pageSize, order, sorting)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get documents")
	}
	count, err := r.repository.CountDocs()
	if err != nil {
		return nil, 0, errors.Wrap(err, "get documents")
	}
	return docs, count, nil
}

func (r *Service) DeleteDocument(documentId uuid.UUID) error {
	err := r.repository.DeleteDocument(documentId)
	if err != nil {
		return errors.Wrap(err, "failed to delete document")
	}

	return nil
}

func (r *Service) UploadDocument(minim, maxim int, name, code string) (*uuid.UUID, error) {
	uid := uuid.New()

	now := time.Now().UTC()
	document := &schemas.DocumentMetadata{
		ID:        uid,
		Name:      name,
		Code:      code,
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

	err = r.repository.CreateFolder(name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create folder")
	}

	return &uid, nil
}

func (r *Service) UploadPage(doc *multipart.FileHeader, documentId uuid.UUID, number int) error {
	uid := uuid.New()

	contents, err := doc.Open()
	if err != nil {
		return errors.Wrap(err, "failed to open uploaded page")
	}

	err = r.repository.SaveToMinio(doc, uid.String()+".jpg", contents)
	if err != nil {
		return errors.Wrap(err, "failed to save page to minio")
	}
	u := getOriginalLink(uid.String() + ".jpg")

	processedImg, err := compressImage(contents)
	if err != nil {
		return errors.Wrap(err, "failed to compress image")
	}
	err = r.repository.SaveThumbToMinio(processedImg, uid.String()+"_thumb.jpg")
	if err != nil {
		return errors.Wrap(err, "failed to save thumbnail in minio")
	}

	uThumb := getOriginalLink(uid.String() + "_thumb.jpg")

	text, err := r.ProcessWithML(doc)
	if err != nil {
		return errors.Wrap(err, "failed to send image")
	}
	page := &schemas.PageMetadata{
		ID:         uid,
		DocumentId: documentId,
		Thumb:      uThumb,
		Original:   u,
		Number:     number,
		FullText:   *text,
	}

	err = r.repository.SavePageToPg(page)
	if err != nil {
		return errors.Wrap(err, "failed to save page to postgres")
	}

	err = r.repository.StatusSuccess(documentId)
	if err != nil {
		return errors.Wrap(err, "failed to change status")
	}

	return nil
}

func (r *Service) SearchDocuments(page, pageSize int, sortBy, name, status string) (*[]schemas.DocumentMetadata, int64, error) {
	order := "DESC"
	sorting, ok := strings.CutPrefix(sortBy, "-")
	if ok {
		order = "ASC"
	}
	if sorting == "createdAt" {
		sorting = "created_at"
	}
	docs, err := r.repository.SearchDocuments(page, pageSize, order, name, status, sorting)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to search documents")
	}
	count, err := r.repository.CountDocs()
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to count docs")
	}
	return docs, count, nil
}

func (r *Service) GetSingleDocument(id uuid.UUID) (*schemas.DocumentMetadata, error) {
	doc, err := r.repository.GetDocumentById(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve document")
	}

	return doc, nil
}

func compressImage(file multipart.File) (*bytes.Buffer, error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to seek file")
	}

	src, err := jpeg.Decode(file)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, errors.Wrap(err, "failed to decode image")
	}

	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/5, src.Bounds().Max.Y/5))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)
	w := new(bytes.Buffer)
	err = png.Encode(w, dst)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode image")
	}

	return w, nil
}

func (r *Service) GetPages(id uuid.UUID, page, pageSize int) ([]schemas.PageMetadata, error) {
	pages, err := r.repository.GetPages(id, page, pageSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pages")
	}

	return pages, nil
}

func getOriginalLink(name string) string {
	u := "/" + settings_utils.Settings.MinioBucketName + "/" + name
	return u
}

func (r *Service) UpdateDocument(doc *schemas.DocumentMetadata) error {
	err := r.repository.SaveDocToPg(doc)
	if err != nil {
		return errors.Wrap(err, "failed to update document")
	}
	return nil
}

func (r *Service) UpdatePage(page *schemas.PageMetadata) error {
	err := r.repository.SavePageToPg(page)
	if err != nil {
		return errors.Wrap(err, "failed to update page")
	}

	return nil
}
