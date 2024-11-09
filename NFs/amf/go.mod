module github.com/free5gc/amf

go 1.21

toolchain go1.22.4

require (
	git.cs.nctu.edu.tw/calee/sctp v1.1.0
	github.com/antihax/optional v1.0.0
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/enable-intelligent-containerized-5g/nas v0.0.0-00010101000000-000000000000
	github.com/enable-intelligent-containerized-5g/ngap v1.0.8
	github.com/enable-intelligent-containerized-5g/openapi v1.0.8
	github.com/free5gc/aper v1.0.6-0.20240503143507-2c4c4780b98f
	github.com/free5gc/util v1.0.3
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.3
	github.com/google/uuid v1.3.0
	github.com/mitchellh/mapstructure v1.4.2
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli v1.22.5
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/aead/cmac v0.0.0-20160719120800-7af84192f0b1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0-20190314233015-f79a8a8ca69d // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/golang-jwt/jwt v3.2.1+incompatible // indirect
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/h2non/parth v0.0.0-20190131123155-b4df798d6542 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/russross/blackfriday/v2 v2.0.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/tim-ywliu/nested-logrus-formatter v1.3.2 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.0.0-20210810183815-faf39c7919d5 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/h2non/gock.v1 v1.1.2 // indirect
)

replace github.com/enable-intelligent-containerized-5g/openapi => ../../../openapi

replace github.com/enable-intelligent-containerized-5g/ngap => ../../../ngap

replace github.com/enable-intelligent-containerized-5g/nas => ../../../nas
