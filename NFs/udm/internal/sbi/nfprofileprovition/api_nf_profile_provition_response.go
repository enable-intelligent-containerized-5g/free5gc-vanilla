package nfprofileprovition

import (
	"net/http"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/udm/internal/context"
	"github.com/free5gc/udm/internal/logger"
	"github.com/free5gc/util/httpwrapper"
	"github.com/gin-gonic/gin"
)

func HTTPNfProfileProvitionRequest(c *gin.Context) {
	self := context.UDM_Self()
	if self == nil {
		logger.NfPer.Error("NWDAF self is not initialized")
		problemDetails := models.ProblemDetails{
			Title:  "NWDAF self is not initialized",
			Status: http.StatusInternalServerError,
			Detail: "NWDAF self is not initialized",
			Cause:  "SYSTEM_FAILURE",
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	}

	rsp := httpwrapper.NewResponse(http.StatusOK, nil, self.NfId)

	c.JSON(rsp.Status, rsp)
}
