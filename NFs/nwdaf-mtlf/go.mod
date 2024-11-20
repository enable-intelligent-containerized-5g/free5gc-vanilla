module github.com/free5gc/nwdaf

go 1.14

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/enable-intelligent-containerized-5g/openapi v1.0.41
	github.com/free5gc/util v1.0.3
	github.com/gin-contrib/cors v1.7.2
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/validator/v10 v10.20.0
	github.com/google/uuid v1.3.0
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/mitchellh/mapstructure v1.4.1
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli v1.22.5
	go.mongodb.org/mongo-driver v1.8.4
	golang.org/x/text v0.19.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/sqlite v1.5.6
	gorm.io/gorm v1.25.12
)

// replace github.com/enable-intelligent-containerized-5g/openapi => ../../../openapi
