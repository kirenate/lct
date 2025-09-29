package services

import (
	"github.com/minio/minio-go"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Service struct {
	minio *minio.Client
	db    *mongo.Client
}

func NewService(minio *minio.Client, db *mongo.Client) *Service {
	return &Service{minio: minio, db: db}
}
