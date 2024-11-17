module github.com/free5gc/nwdaf

go 1.14

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/enable-intelligent-containerized-5g/openapi v1.0.4
	github.com/free5gc/util v1.0.3
	github.com/gin-contrib/cors v1.7.2
	github.com/gin-gonic/gin v1.9.1
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/prometheus/client_golang v1.18.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli v1.22.5
	go.mongodb.org/mongo-driver v1.8.4
	golang.org/x/text v0.19.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/sqlite v1.5.6 // indirect
	gorm.io/gorm v1.25.12 // indirect
)

replace github.com/enable-intelligent-containerized-5g/openapi => ../../../openapi
