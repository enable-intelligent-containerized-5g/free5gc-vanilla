module github.com/free5gc/udm

go 1.14

require (
	github.com/antihax/optional v1.0.0
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/enable-intelligent-containerized-5g/openapi v1.0.41
	github.com/free5gc/util v1.0.3
	github.com/gin-gonic/gin v1.7.3
	github.com/google/uuid v1.3.0
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli v1.22.5
	golang.org/x/crypto v0.14.0
	gopkg.in/yaml.v2 v2.4.0
)

// replace github.com/enable-intelligent-containerized-5g/openapi => ../../../openapi
