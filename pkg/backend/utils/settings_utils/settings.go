package settings_utils

import (
	"github.com/teadove/teasutils/service_utils/settings_utils"
)

type settings struct {
	MinioAccessKeyID     string `env:"MINIO_ACCESS_KEY_ID"`
	MinioSecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY"`
	MinioEndpoint        string `env:"MINIO_ENDPOINT"      envDefault:"0.0.0.0:9000"`
	MinioBucketName      string `env:"MINIO_BUCKET_NAME"   envDefault:"scans"`
	MinioSecure          bool   `env:"MINIO_SECURE" envDefault:"false"`

	PgDb       string `env:"PG_DB"`
	PgHost     string `env:"PG_HOST" envDefault:"localhost"`
	PgPort     int    `env:"PG_PORT"`
	PgUsername string `env:"PG_USERNAME"`
	PgPassword string `env:"PG_PASSWORD"`

	TLS      bool   `env:"TLS"                             envDefault:"false"`
	URL      string `env:"URL"                             envDefault:"0.0.0.0:8080"`
	CertFile string `env:"SSL_CERTFILE"                       envDefault:"./.data/cert.pem"`
	KeyFile  string `env:"SSL_KEYFILE"                        envDefault:"./.data/cert.key"`

	MlUrl string `env:"ML_URL"`
}

var Settings = settings_utils.MustGetSetting[settings]("LCT_") //nolint:gochecknoglobals // required
