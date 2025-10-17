package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"main.go/schemas"
	"strings"
)

func (r *Presentation) getDocumentPages(c *fiber.Ctx) error {
	path := c.Path()
	l := strings.Split(path, "/")
	idStr := l[len(l)-2]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}

	page := c.QueryInt("page")
	pageSize := c.QueryInt("pageSize")

	pages, total, err := r.service.GetPages(id, page, pageSize)
	if err != nil {
		return errors.Wrap(err, "pages not found")
	}
	return c.JSON(fiber.Map{"data": pages, "total": total})
}

func (r *Presentation) getSingleDocument(c *fiber.Ctx) error {
	path := c.Path()
	l := strings.Split(path, "/")
	idStr := l[len(l)-1]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}

	doc, err := r.service.GetSingleDocument(id)
	if err != nil {
		return errors.Wrap(err, "failed to find document")
	}

	return c.JSON(fiber.Map{"data": doc})
}

func (r *Presentation) getDocuments(c *fiber.Ctx) error {
	page := c.QueryInt("page")
	pageSize := c.QueryInt("pageSize")
	sortBy := c.Query("sortBy", "name")
	name := c.Query("query")
	status := c.Query("status")
	if name != "" || status != "" {
		docs, count, err := r.service.SearchDocuments(page, pageSize, sortBy, name, status)
		if err != nil {
			return errors.Wrap(err, "failed to search documents")
		}
		return c.JSON(fiber.Map{"data": *docs, "total": count})
	}
	docs, count, err := r.service.GetDocuments(page, pageSize, sortBy)
	if err != nil {
		return errors.Wrap(err, "failed to get documents")
	}

	return c.JSON(fiber.Map{"data": docs, "total": count})
}

func (r *Presentation) deleteDocument(c *fiber.Ctx) error {
	documentId := c.Query("documentId")
	id, err := uuid.Parse(documentId)
	if err != nil {
		return errors.Wrap(err, "failed to parse uuid")
	}
	err = r.service.DeleteDocument(id)
	if err != nil {
		return errors.Wrap(err, "failed to delete document")
	}
	return nil
}

func (r *Presentation) uploadDocument(c *fiber.Ctx) error {
	doc, err := c.MultipartForm()
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}
	minim := c.QueryInt("min", 0)
	maxim := c.QueryInt("max", 100)
	name := c.Query("name")
	code := c.Query("code")
	if name == "" {
		return errors.New("document must have a name")
	}
	for _, v := range doc.File {
		documentId, err := r.service.UploadDocument(minim, maxim, name, code)
		if err != nil {
			return errors.Wrap(err, "failed to upload document")
		}
		for i, vv := range v {
			err = r.service.UploadPage(vv, *documentId, i)
			if err != nil {
				return errors.Wrap(err, "failed to upload file")
			}
		}
	}

	return nil
}

func (r *Presentation) editDocument(c *fiber.Ctx) error {
	var body *schemas.DocumentMetadata
	err := c.BodyParser(&body)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}
	path := c.Path()
	l := strings.Split(path, "/")
	idStr := l[len(l)-1]
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse uuid")
	}

	err = r.service.UpdateDocument(body, id)
	if err != nil {
		return errors.Wrap(err, "failed to update")
	}

	return nil
}

func (r *Presentation) editPage(c *fiber.Ctx) error {
	var body *string
	err := c.BodyParser(&body)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusUnprocessableEntity,
			Message: err.Error(),
		}
	}

	path := c.Path()
	l := strings.Split(path, "/")
	idStr := l[len(l)-1]
	pageId, err := uuid.Parse(idStr)
	if err != nil {
		return errors.Wrap(err, "failed to parse uuid")
	}

	err = r.service.UpdatePage(body, pageId)
	if err != nil {
		return errors.Wrap(err, "failed to update")
	}

	return nil
}

//
//func (r *Presentation) exportDocument(c *fiber.Ctx) error {
//	id := c.Query("documentId")
//	uid, err := uuid.Parse(id)
//	if err != nil {
//		return &fiber.Error{
//			Code:    fiber.StatusUnprocessableEntity,
//			Message: err.Error(),
//		}
//	}
//	ext := c.Query("ext")
//	textByte := c.Body()
//	var text [][]string
//	err = json.Unmarshal(textByte, &text)
//	if err != nil {
//		return &fiber.Error{
//			Code:    fiber.StatusUnprocessableEntity,
//			Message: err.Error(),
//		}
//	}
//
//	export := &services.Export{ID: uid, Text: text, Ext: ext}
//	csv, err := r.service.CreatePageAndCsv(export)
//	if err != nil {
//		return errors.Wrap(err, "failed to export document")
//	}
//
//	return nil
//}
