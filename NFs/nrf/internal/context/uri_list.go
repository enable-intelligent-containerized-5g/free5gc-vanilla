/*
 * NRF UriList
 */

package context

import (
	"github.com/enable-intelligent-containerized-5g/openapi/models"
)

type UriList struct {
	NfType models.NfType `json:"nfType" bson:"nfType"`
	Link   Links         `json:"_link" bson:"_link" mapstructure:"_link"`
}
