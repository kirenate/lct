package main

import (
	"fmt"
	minio2 "github.com/minio/minio-go"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"main.go/presentations/web"
	"main.go/repositories"
	"main.go/schemas"
	"main.go/services"
	"main.go/utils/settings_utils"
)

func main() {
	minio, err := minio2.New(settings_utils.Settings.MinioEndpoint, settings_utils.Settings.MinioAccessKeyID,
		settings_utils.Settings.MinioSecretAccessKey, true)
	if err != nil {
		panic(errors.Wrap(err, "failed to initiate minio client"))
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		settings_utils.Settings.PgHost, settings_utils.Settings.PgPort, settings_utils.Settings.PgUsername,
		settings_utils.Settings.PgPassword, settings_utils.Settings.PgDb)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		panic(errors.Wrap(err, "failed to connect database"))
	}

	err = db.AutoMigrate(&schemas.DocumentMetadata{}, &schemas.PageMetadata{}, &schemas.Attribute{}, &schemas.Text{})
	if err != nil {
		panic(errors.Wrap(err, "failed to merge database"))
	}

	repository := repositories.NewRepository(minio, db)
	service, err := services.NewService(repository)
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize service"))
	}

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
