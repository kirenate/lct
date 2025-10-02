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
	RecScores []int    `json:"rec_scores" gorm:"serializer:json"`
	RecPolys  [][]int  `json:"rec_polys" gorm:"serializer:json"`
	RecTexts  []string `json:"rec_texts" gorm:"serializer:json"`
}
type PageMetadata struct {
	ID         uuid.UUID `json:"id"`
	DocumentId uuid.UUID `json:"documentId"`
	Thumb      string    `json:"thumb"`
	Original   string    `json:"original"`
	Number     int       `json:"number"`
	Text       Text      `json:"text" gorm:"serializer:json"`
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
