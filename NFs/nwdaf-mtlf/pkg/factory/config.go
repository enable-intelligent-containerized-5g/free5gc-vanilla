/*
 * NWDAF Configuration Factory
 */

package factory

import (
	"fmt"

	"github.com/asaskevich/govalidator"

	logger_util "github.com/free5gc/util/logger"
)

const (
	NWDAF_EXPECTED_CONFIG_VERSION = "1.0.0"
)

type Config struct {
	Info          *Info               `yaml:"info" valid:"required"`
	Configuration *Configuration      `yaml:"configuration" valid:"required"`
	Logger        *logger_util.Logger `yaml:"logger" valid:"optional"`
}

func (c *Config) Validate() (bool, error) {
	if info := c.Info; info != nil {
		if result, err := info.validate(); err != nil {
			return result, err
		}
	}

	if configuration := c.Configuration; configuration != nil {
		if result, err := configuration.validate(); err != nil {
			return result, err
		}
	}

	if logger := c.Logger; logger != nil {
		if result, err := logger.Validate(); err != nil {
			return result, err
		}
	}

	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
}

type Info struct {
	Version     string `yaml:"version,omitempty" valid:"type(string),required"`
	Description string `yaml:"description,omitempty" valid:"type(string),optional"`
}

func (i *Info) validate() (bool, error) {
	result, err := govalidator.ValidateStruct(i)
	return result, appendInvalid(err)
}

const (
	NWDAF_DEFAULT_IPV4     = "127.0.0.23"
	NWDAF_DEFAULT_PORT     = "8000"
	NWDAF_DEFAULT_PORT_INT = 8000
)

type Configuration struct {
	ContainerName string `yaml:"containerName,omitempty"`
	// SqlLiteTableName string   `yaml:"SqlLiteTableName" valid:"type(string),required"`
	SqlLiteDB           string               `yaml:"sqlLiteDB" valid:"type(string),required"`
	Name                string               `yaml:"name,omitempty"`
	NwdafName           string               `yaml:"nwdafName,omitempty"`
	Sbi                 *Sbi                 `yaml:"sbi" valid:"required"`
	NrfUri              string               `yaml:"nrfUri" valid:"url,required"`
	OamUri              string               `yaml:"oamUri" valid:"url,required"`
	KsmInstance         string               `yaml:"ksmInstance" valid:"required"`
	Namespace           string               `yaml:"namespace" valid:"required"`
	ServiceNameList     []string             `yaml:"serviceNameList" valid:"required"`
	MlModelTrainingInfo *MlModelTrainingInfo `yaml:"mlModelTrainingInfo" valid:"required"`
}

func (c *Configuration) validate() (bool, error) {
	govalidator.TagMap["scheme"] = govalidator.Validator(func(str string) bool {
		return str == "https" || str == "http"
	})
	result, err := govalidator.ValidateStruct(c)
	return result, appendInvalid(err)
}

type Sbi struct {
	Scheme       string `yaml:"scheme" valid:"scheme,required"`
	RegisterIPv4 string `yaml:"registerIPv4,omitempty" valid:"host,optional"` // IP that is registered at NRF.
	// IPv6Addr string `yaml:"ipv6Addr,omitempty"`
	BindingIPv4 string `yaml:"bindingIPv4,omitempty" valid:"host,optional"` // IP used to run the server in the node.
	Port        int    `yaml:"port" valid:"port,required"`
	Tls         *Tls   `yaml:"tls,omitempty" valid:"optional"`
}

type Tls struct {
	Pem string `yaml:"pem,omitempty" valid:"type(string),minstringlength(1),required"`
	Key string `yaml:"key,omitempty" valid:"type(string),minstringlength(1),required"`
}

type MlModelTrainingInfo struct {
	TimeSteps int64 `yaml:"timeSteps,omitempty" valid:"type(int64),required"`
}

func appendInvalid(err error) error {
	var errs govalidator.Errors

	if err == nil {
		return nil
	}

	es := err.(govalidator.Errors).Errors()
	for _, e := range es {
		errs = append(errs, fmt.Errorf("Invalid %w", e))
	}

	return error(errs)
}

func (c *Config) GetVersion() string {
	if c.Info != nil && c.Info.Version != "" {
		return c.Info.Version
	}
	return ""
}
