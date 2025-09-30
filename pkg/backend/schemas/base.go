package schemas

import (
	"github.com/google/uuid"
	"time"
)

type Attribute struct {
	ID         uuid.UUID
	PageId     uuid.UUID
	DocumentId uuid.UUID
	Name       string
	Value      string
	Confidence int
	X1         int
	Y1         int
	X2         int
	Y2         int
}

type Text struct {
	ID         uuid.UUID
	PageId     uuid.UUID
	DocumentId uuid.UUID
	Text       string
	X1         int
	Y1         int
	X2         int
	Y2         int
}
type PageMetadata struct {
	ID         uuid.UUID
	DocumentId uuid.UUID
	Thumb      string
	Original   string
	Number     int
}

type Page struct {
	ID   uuid.UUID
	Meta *PageMetadata
}

type DocumentMetadata struct {
	ID        uuid.UUID
	Code      string
	Name      string
	Status    string
	Progress  int
	Min       int
	Max       int
	CreatedAt time.Time
}
