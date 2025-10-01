package schemas

import (
	"github.com/google/uuid"
	"time"
)

type Attribute struct {
	ID         uuid.UUID `json:"id"`
	PageId     uuid.UUID `json:"pageId"`
	DocumentId uuid.UUID `json:"documentId"`
	Name       string    `json:"name"`
	Value      string    `json:"value"`
	Confidence int       `json:"confidence"`
	X1         int       `json:"x1"`
	Y1         int       `json:"y1"`
	X2         int       `json:"x2"`
	Y2         int       `json:"y2"`
}

type Text struct {
	ID         uuid.UUID `json:"id"`
	PageId     uuid.UUID `json:"pageId"`
	DocumentId uuid.UUID `json:"documentId"`
	Text       string    `json:"text"`
	X1         int       `json:"x1"`
	Y1         int       `json:"y1"`
	X2         int       `json:"x2"`
	Y2         int       `json:"y2"`
}
type PageMetadata struct {
	ID         uuid.UUID `json:"id"`
	DocumentId uuid.UUID `json:"documentId"`
	Thumb      string    `json:"thumb"`
	Original   string    `json:"original"`
	Number     int       `json:"number"`
}

type Page struct {
	ID   uuid.UUID     `json:"id"`
	Meta *PageMetadata `json:"meta"`
}

type DocumentMetadata struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Progress  int       `json:"progress"`
	Min       int       `json:"min"`
	Max       int       `json:"max"`
	CreatedAt time.Time `json:"createdAt"`
}
