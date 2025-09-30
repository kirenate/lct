package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"strconv"
)

func (r *Presentation) getDocumentPages(c *fiber.Ctx) error {
	return nil
}

func (r *Presentation) getDocuments(c *fiber.Ctx) error {
	params := c.AllParams()
	page, err := strconv.Atoi(params["page"])
	if err != nil {
		return errors.Wrap(err, "failed to convert param page to int")
	}
	pageSize, err := strconv.Atoi(params["pageSize"])
	if err != nil {
		return errors.Wrap(err, "failed to convert param pageSize to int")
	}
	order := "DESC"
	if params["sortBy"] == "-name" {
		order = "ASC"
	}
	docs, err := r.service.GetDocuments(page, pageSize, order)
	if err != nil {
		return errors.Wrap(err, "failed to get documents")
	}
	return c.JSON(fiber.Map{"documents": docs})
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
	minim, err := strconv.Atoi(c.Params("min", "0"))
	if err != nil {
		return errors.Wrap(err, "failed to convert min value to int")
	}
	maxim, err := strconv.Atoi(c.Params("max", "100"))
	if err != nil {
		return errors.Wrap(err, "failed to convert max value to int")
	}

	name := c.Params("name")

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
