package schemas

import (
	"github.com/google/uuid"
	"time"
)

type Attribute struct {
	id         uuid.UUID
	name       string
	value      string
	confidence int
	x1         int
	y1         int
	x2         int
	y2         int
}

type Text struct {
	id   uuid.UUID
	text string
	x1   int
	y1   int
	x2   int
	y2   int
}
type PageMetadata struct {
	id         uuid.UUID
	documentId uuid.UUID
	thumb      string
	original   string
	number     int
	fullText   []Text
	attributes []Attribute
}

type Page struct {
	id   uuid.UUID
	meta *PageMetadata
}

type DocumentMetadata struct {
	id        uuid.UUID
	code      string
	name      string
	status    string
	progress  int
	min       int
	max       int
	createdAt time.Time
}
