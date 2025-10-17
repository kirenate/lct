package services

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/teadove/teasutils/service_utils/logger_utils"
	"golang.org/x/image/draw"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"main.go/repositories"
	"main.go/schemas"
	"main.go/utils/settings_utils"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

type Export struct {
	ID   uuid.UUID  `json:"documentId"`
	Text [][]string `json:"text"`
	Ext  string     `json:"ext"`
}
type Service struct {
	repository *repositories.Repository
	reader     *kafka.Reader
	writer     *kafka.Writer
}

func NewService(repository *repositories.Repository, reader *kafka.Reader, writer *kafka.Writer) (*Service, error) {
	service := &Service{repository: repository, reader: reader, writer: writer}
	err := service.repository.CreateBucket(settings_utils.Settings.MinioBucketName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create bucket")
	}
	go service.BackgroundConsumer(logger_utils.NewLoggedCtx())

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

	go r.PageLoaderChecker(logger_utils.NewLoggedCtx(), uid)

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

	err = r.SendToQueue(uid, documentId)
	if err != nil {
		return errors.Wrap(err, "failed to send msg to queue")
	}

	page := &schemas.PageMetadata{
		ID:         uid,
		DocumentId: documentId,
		Thumb:      uThumb,
		Original:   u,
		Number:     number,
		Progress:   repositories.StatusProcessing,
	}

	err = r.repository.SavePageToPg(page)
	if err != nil {
		return errors.Wrap(err, "failed to save page to postgres")
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

func (r *Service) GetPages(id uuid.UUID, page, pageSize int) ([]schemas.PageMetadata, int64, error) {
	pages, err := r.repository.GetPages(id, page, pageSize)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get pages")
	}
	count, err := r.repository.CountPages(id)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to count pages")
	}
	return pages, count, nil
}

func getOriginalLink(name string) string {
	u := "/" + settings_utils.Settings.MinioBucketName + "/" + name
	return u
}

func (r *Service) UpdateDocument(doc *schemas.DocumentMetadata, id uuid.UUID) error {
	err := r.repository.UpdateDocument(doc, id)
	if err != nil {
		return errors.Wrap(err, "failed to update document")
	}
	return nil
}

func (r *Service) UpdatePage(page string, pageId uuid.UUID) error {
	err := r.repository.UpdatePage(page, pageId)
	if err != nil {
		return errors.Wrap(err, "failed to update page")
	}

	return nil
}

func (r *Service) CreatePageAndCsv(export *Export) ([]string, error) {
	page, err := r.repository.GetSinglePage(export.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pages")
	}

	link := getOriginalLink(page.ID.String())
	out, err := os.Create("./downloads/" + page.ID.String() + ".jpg")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create jpg file")
	}
	defer out.Close()
	resp, err := http.Get(link)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get image from minio")
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to copy response")
	}
	defer resp.Body.Close()

	file, err := os.Create("./downloads/" + page.ID.String() + ".csv")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create csv file")
	}
	defer file.Close()

	w := csv.NewWriter(file)
	err = w.WriteAll(export.Text)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write csv file")
	}

	return []string{"./downloads/" + page.ID.String() + ".jpg", "./downloads/" + page.ID.String() + ".csv"}, nil
}

func (r *Service) ZipFiles(paths []string) error {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	defer w.Close()
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return errors.Wrap(err, "failed to open file")
		}
		f, err := w.Create(file.Name())
		if err != nil {
			return errors.Wrap(err, "failed to create file inside zip")
		}
		_, err = file.WriteTo(f)
		if err != nil {
			return errors.Wrap(err, "failed to write to zip file")
		}
	}
	//n := strings.Split(paths[0], "/")
	//name := strings.Replace(n[len(n)-1], ".jpg", "", 1)
	//out, err := os.Create(name)
	//if err != nil {
	//	return errors.Wrap(err,"failed to create zip archive")
	//}
	//
	//w.Copy(out)
	return nil
}
