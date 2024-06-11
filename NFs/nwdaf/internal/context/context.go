package context

import (
	"fmt"

	"github.com/enable-intelligent-and-containerized-5g/openapi/models"
)

var nwdafContext = NWDAFContext{}

func init() {
	nwdafContext.Name = "nwdaf"
}

type NWDAFContext struct {
	Name                                    string
	UriScheme                               models.UriScheme
	BindingIPv4                             string
	SBIPort                                 int
	RegisterIPv4                            string // IP register to NRF
	HttpIPv6Address                         string
	NfId                                    string
	NrfUri                                  string
}

// Reset NWDAF Context
func (context *NWDAFContext) Reset() {
	context.UriScheme = models.UriScheme_HTTPS
	context.Name = "udr"
}

func (context *NWDAFContext) GetIPv4Uri() string {
	return fmt.Sprintf("%s://%s:%d", context.UriScheme, context.RegisterIPv4, context.SBIPort)
}

func NWDAF_Self() *NWDAFContext {
	return &nwdafContext
}
