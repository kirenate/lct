package schemas

import (
	"github.com/google/uuid"
	"time"
)

type Attribute struct {
	ID         uuid.UUID `json:"id,omitempty"`
	PageId     uuid.UUID `json:"pageId,omitempty"`
	DocumentId uuid.UUID `json:"documentId,omitempty"`
	Name       string    `json:"name,omitempty"`
	Value      string    `json:"value,omitempty"`
	Confidence int       `json:"confidence,omitempty"`
}

type Text struct {
	RecScores []int       `json:"rec_scores" gorm:"serializer:json"`
	RecPolys  [][]float64 `json:"rec_polys" gorm:"serializer:json"`
	RecTexts  []string    `json:"rec_texts" gorm:"serializer:json"`
}

type TextJson struct {
	ID         uuid.UUID `json:"id"`
	Text       string    `json:"text"`
	Confidence int       `json:"confidence"`
	X1         float64   `json:"x1"`
	Y1         float64   `json:"y1"`
	X2         float64   `json:"x2"`
	Y2         float64   `json:"y2"`
	X3         float64   `json:"x3"`
	Y3         float64   `json:"y3"`
	X4         float64   `json:"x4"`
	Y4         float64   `json:"y4"`
}

type PageMetadata struct {
	ID         uuid.UUID   `json:"id"`
	DocumentId uuid.UUID   `json:"documentId"`
	Thumb      string      `json:"thumb"`
	Original   string      `json:"original"`
	Number     int         `json:"number"`
	FullText   []byte      `json:"fullText" gorm:"serializer:json"`
	K          int         `json:"k"`
	Attrs      []Attribute `json:"attrs,omitempty" gorm:"serializer:json"`
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
