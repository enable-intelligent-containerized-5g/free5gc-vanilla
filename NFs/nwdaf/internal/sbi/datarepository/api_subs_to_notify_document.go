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
// "net/http"
//
// "github.com/gin-gonic/gin"
//
// "github.com/free5gc/nwdaf/internal/logger"
// "github.com/free5gc/nwdaf/internal/sbi/producer"
// "github.com/free5gc/openapi"
// "github.com/free5gc/openapi/models"
// "github.com/free5gc/util/httpwrapper"
)

// HTTPRemovesubscriptionDataSubscriptions - Deletes a subscriptionDataSubscriptions
// func HTTPRemovesubscriptionDataSubscriptions(c *gin.Context) {
// 	req := httpwrapper.NewRequest(c.Request, nil)
// 	req.Params["ueId"] = c.Params.ByName("ueId")
//
// 	rsp := producer.HandleRemovesubscriptionDataSubscriptions(req)
//
// 	responseBody, err := openapi.Serialize(rsp.Body, "application/json")
// 	if err != nil {
// 		logger.DataRepoLog.Errorln(err)
// 		problemDetails := models.ProblemDetails{
// 			Status: http.StatusInternalServerError,
// 			Cause:  "SYSTEM_FAILURE",
// 			Detail: err.Error(),
// 		}
// 		c.JSON(http.StatusInternalServerError, problemDetails)
// 	} else {
// 		c.Data(rsp.Status, "application/json", responseBody)
// 	}
// }
