package main

import (
	minio2 "github.com/minio/minio-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"main.go/presentations/web"
	"main.go/services"
	"main.go/utils/settings_utils"
)

func main() {
	minio, err := minio2.New(settings_utils.Settings.MinioEndpoint, settings_utils.Settings.MinioAccessKeyID,
		settings_utils.Settings.MinioSecretAccessKey, true)
	if err != nil {
		panic(errors.Wrap(err, "failed to initiate minio client"))
	}

	mongoClient, err := mongo.Connect(options.Client().ApplyURI(settings_utils.Settings.MongoURI))
	if err != nil {
		panic(errors.Wrap(err, "failed to connect database"))
	}

	service := services.NewService(minio, mongoClient)
	presentation := web.NewPresentation(service)

	app := presentation.BuildApp()

	if settings_utils.Settings.TLS {
		err = app.ListenTLS(
			settings_utils.Settings.URL,
			settings_utils.Settings.CertFile,
			settings_utils.Settings.KeyFile,
		)
		if err != nil {
			panic(errors.Wrap(err, "failed to start tls server"))
		}

		return
	}

	err = app.Listen(settings_utils.Settings.URL)
	if err != nil {
		panic(errors.Wrap(err, "failed to start server"))
	}
}
