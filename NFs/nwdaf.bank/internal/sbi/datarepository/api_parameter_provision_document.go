/*
 * Nnwdaf_DataRepository API OpenAPI file
 *
 * Unified Data Repository Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package datarepository

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi/producer"
	"github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/util/httpwrapper"
)

// HTTPGetppData - Read the profile of a given UE
func HTTPGetppData(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	req.Params["ueId"] = c.Params.ByName("ueId")

	rsp := producer.HandleGetppData(req)

	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.DataRepoLog.Errorln(err)
		problemDetails := models.ProblemDetails{
			Status: http.StatusInternalServerError,
			Cause:  "SYSTEM_FAILURE",
			Detail: err.Error(),
		}
		c.JSON(http.StatusInternalServerError, problemDetails)
	} else {
		c.Data(rsp.Status, "application/json", responseBody)
	}
}
