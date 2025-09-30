package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (r *Presentation) getDocumentPages(c *fiber.Ctx) error {
	return nil
}

func (r *Presentation) getDocuments(c *fiber.Ctx) error {
	page := c.QueryInt("page")
	pageSize := c.QueryInt("pageSize")
	sortBy := c.Query("sortBy", "name")
	name := c.Query("query")
	status := c.Query("status")
	order := "DESC"
	if sortBy == "-name" {
		order = "ASC"
	}
	if name != "" || status != "" {
		docs, err := r.service.SearchDocuments(page, pageSize, order, name, status)
		if err != nil {
			return errors.Wrap(err, "failed to search documents")
		}
		return c.JSON(fiber.Map{"documents": docs})
	}
	docs, err := r.service.GetDocuments(page, pageSize, order)
	if err != nil {
		return errors.Wrap(err, "failed to get documents")
	}
	return c.JSON(fiber.Map{"documents": docs, "documentsLoaded": len(docs)})
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
	if name == "" {
		return errors.New("document must have a name")
	}
	for _, v := range doc.File {
		documentId, err := r.service.UploadDocument(minim, maxim, name)
		if err != nil {
			return errors.Wrap(err, "failed to upload document")
		}
		for i, vv := range v {
			err = r.service.UploadPage(vv, *documentId, i, name)
			if err != nil {
				return errors.Wrap(err, "failed to upload file")
			}
		}
	}

	return nil
}
