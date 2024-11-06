/*
 * Nudr_DataRepository API OpenAPI file
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
	"github.com/free5gc/nwdaf/internal/util"
	"github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/util/httpwrapper"
)

func sendResponse(c *gin.Context, rsp *httpwrapper.Response) {
	for k, v := range rsp.Header {
		// TODO: concatenate all values
		c.Header(k, v[0])
	}
	serializedBody, err := openapi.Serialize(rsp.Body, "application/json")
	if err != nil {
		logger.DataRepoLog.Errorf("Serialize Response Body error: %+v", err)
		pd := util.ProblemDetailsSystemFailure(err.Error())
		c.JSON(http.StatusInternalServerError, pd)
	} else {
		c.Data(rsp.Status, "application/json", serializedBody)
	}
}

func getDataFromRequestBody(c *gin.Context, data interface{}) error {
	reqBody, err := c.GetRawData()
	if err != nil {
		logger.DataRepoLog.Errorf("Get Request Body error: %+v", err)
		pd := util.ProblemDetailsSystemFailure(err.Error())
		c.JSON(http.StatusInternalServerError, pd)
		return err
	}

	err = openapi.Deserialize(data, reqBody, "application/json")
	if err != nil {
		logger.DataRepoLog.Errorf("Deserialize Request Body error: %+v", err)
		pd := util.ProblemDetailsMalformedReqSyntax(err.Error())
		c.JSON(http.StatusBadRequest, pd)
		return err
	}
	return err
}

// HTTPApplicationDataInfluenceDataGet -
func HTTPApplicationDataInfluenceDataGet(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	rsp := producer.HandleApplicationDataInfluenceDataGet(queryParams)
	sendResponse(c, rsp)
}

// HTTPApplicationDataInfluenceDataInfluenceIdDelete -
func HTTPApplicationDataInfluenceDataInfluenceIdDelete(c *gin.Context) {
	rsp := producer.HandleApplicationDataInfluenceDataInfluenceIdDelete(c.Params.ByName("influenceId"))
	sendResponse(c, rsp)
}

// HTTPApplicationDataInfluenceDataInfluenceIdPatch -
func HTTPApplicationDataInfluenceDataInfluenceIdPatch(c *gin.Context) {
	var trInfluDataPatch models.TrafficInfluDataPatch

	if err := getDataFromRequestBody(c, &trInfluDataPatch); err != nil {
		return
	}

	rsp := producer.HandleApplicationDataInfluenceDataInfluenceIdPatch(c.Params.ByName("influenceId"),
		&trInfluDataPatch)

	sendResponse(c, rsp)
}

// HTTPApplicationDataInfluenceDataInfluenceIdPut -
func HTTPApplicationDataInfluenceDataInfluenceIdPut(c *gin.Context) {
	var trInfluData models.TrafficInfluData

	if err := getDataFromRequestBody(c, &trInfluData); err != nil {
		return
	}

	rsp := producer.HandleApplicationDataInfluenceDataInfluenceIdPut(c.Params.ByName("influenceId"), &trInfluData)

	sendResponse(c, rsp)
}

// HTTPApplicationDataInfluenceDataSubsToNotifyGet -
func HTTPApplicationDataInfluenceDataSubsToNotifyGet(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	rsp := producer.HandleApplicationDataInfluenceDataSubsToNotifyGet(queryParams)
	sendResponse(c, rsp)
}

// HTTPApplicationDataInfluenceDataSubsToNotifyPost -
// func HTTPApplicationDataInfluenceDataSubsToNotifyPost(c *gin.Context) {
// 	var trInfluSub models.TrafficInfluSub

// 	if err := getDataFromRequestBody(c, &trInfluSub); err != nil {
// 		return
// 	}

// 	rsp := producer.HandleApplicationDataInfluenceDataSubsToNotifyPost(&trInfluSub)

// 	sendResponse(c, rsp)
// }

// HTTPApplicationDataInfluenceDataSubsToNotifySubscriptionIdDelete -
func HTTPApplicationDataInfluenceDataSubsToNotifySubscriptionIdDelete(c *gin.Context) {
	rsp := producer.HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdDelete(
		c.Params.ByName("subscriptionId"))
	sendResponse(c, rsp)
}

// HTTPApplicationDataInfluenceDataSubsToNotifySubscriptionIdGet -
func HTTPApplicationDataInfluenceDataSubsToNotifySubscriptionIdGet(c *gin.Context) {
	rsp := producer.HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdGet(
		c.Params.ByName("subscriptionId"))
	sendResponse(c, rsp)
}

// HTTPApplicationDataInfluenceDataSubsToNotifySubscriptionIdPut -
func HTTPApplicationDataInfluenceDataSubsToNotifySubscriptionIdPut(c *gin.Context) {
	var trInfluSub models.TrafficInfluSub

	if err := getDataFromRequestBody(c, &trInfluSub); err != nil {
		return
	}

	rsp := producer.HandleApplicationDataInfluenceDataSubsToNotifySubscriptionIdPut(
		c.Params.ByName("subscriptionId"), &trInfluSub)

	sendResponse(c, rsp)
}

// HTTPApplicationDataPfdsAppIdDelete -
func HTTPApplicationDataPfdsAppIdDelete(c *gin.Context) {
	rsp := producer.HandleApplicationDataPfdsAppIdDelete(c.Params.ByName("appId"))
	sendResponse(c, rsp)
}

// HTTPApplicationDataPfdsAppIdGet -
func HTTPApplicationDataPfdsAppIdGet(c *gin.Context) {
	rsp := producer.HandleApplicationDataPfdsAppIdGet(c.Params.ByName("appId"))
	sendResponse(c, rsp)
}

// HTTPApplicationDataPfdsAppIdPut -
func HTTPApplicationDataPfdsAppIdPut(c *gin.Context) {
	var pfdDataforApp models.PfdDataForApp

	if err := getDataFromRequestBody(c, &pfdDataforApp); err != nil {
		return
	}

	rsp := producer.HandleApplicationDataPfdsAppIdPut(c.Params.ByName("appId"), &pfdDataforApp)

	sendResponse(c, rsp)
}

// HTTPApplicationDataPfdsGet -
func HTTPApplicationDataPfdsGet(c *gin.Context) {
	query := c.Request.URL.Query()
	rsp := producer.HandleApplicationDataPfdsGet(query["appId"])
	sendResponse(c, rsp)
}

// HTTPExposureDataSubsToNotifyPost -
func HTTPExposureDataSubsToNotifyPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// HTTPExposureDataSubsToNotifySubIdDelete - Deletes a subcription for notifications
func HTTPExposureDataSubsToNotifySubIdDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// HTTPExposureDataSubsToNotifySubIdPut - updates a subcription for notifications
func HTTPExposureDataSubsToNotifySubIdPut(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// HTTPPolicyDataBdtDataBdtReferenceIdDelete -
func HTTPPolicyDataBdtDataBdtReferenceIdDelete(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	req.Params["bdtReferenceId"] = c.Params.ByName("bdtReferenceId")

	rsp := producer.HandlePolicyDataBdtDataBdtReferenceIdDelete(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataBdtDataBdtReferenceIdGet -
func HTTPPolicyDataBdtDataBdtReferenceIdGet(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	req.Params["bdtReferenceId"] = c.Params.ByName("bdtReferenceId")

	rsp := producer.HandlePolicyDataBdtDataBdtReferenceIdGet(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataBdtDataBdtReferenceIdPut -
func HTTPPolicyDataBdtDataBdtReferenceIdPut(c *gin.Context) {
	var bdtData models.BdtData

	if err := getDataFromRequestBody(c, &bdtData); err != nil {
		return
	}

	req := httpwrapper.NewRequest(c.Request, bdtData)
	req.Params["bdtReferenceId"] = c.Params.ByName("bdtReferenceId")

	rsp := producer.HandlePolicyDataBdtDataBdtReferenceIdPut(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataBdtDataGet -
func HTTPPolicyDataBdtDataGet(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)

	rsp := producer.HandlePolicyDataBdtDataGet(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataPlmnsPlmnIdUePolicySetGet -
func HTTPPolicyDataPlmnsPlmnIdUePolicySetGet(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	req.Params["plmnId"] = c.Params.ByName("plmnId")

	rsp := producer.HandlePolicyDataPlmnsPlmnIdUePolicySetGet(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataSponsorConnectivityDataSponsorIdGet -
func HTTPPolicyDataSponsorConnectivityDataSponsorIdGet(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	req.Params["sponsorId"] = c.Params.ByName("sponsorId")

	rsp := producer.HandlePolicyDataSponsorConnectivityDataSponsorIdGet(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataSubsToNotifyPost -
// func HTTPPolicyDataSubsToNotifyPost(c *gin.Context) {
// 	var policyDataSubscription models.PolicyDataSubscription

// 	if err := getDataFromRequestBody(c, &policyDataSubscription); err != nil {
// 		return
// 	}

// 	req := httpwrapper.NewRequest(c.Request, policyDataSubscription)
// 	req.Params["ueId"] = c.Params.ByName("ueId")

// 	rsp := producer.HandlePolicyDataSubsToNotifyPost(req)

// 	for key, val := range rsp.Header {
// 		c.Header(key, val[0])
// 	}

// 	sendResponse(c, rsp)
// }

// HTTPPolicyDataSubsToNotifySubsIdDelete -
// func HTTPPolicyDataSubsToNotifySubsIdDelete(c *gin.Context) {
// 	req := httpwrapper.NewRequest(c.Request, nil)
// 	req.Params["subsId"] = c.Params.ByName("subsId")

// 	rsp := producer.HandlePolicyDataSubsToNotifySubsIdDelete(req)

// 	sendResponse(c, rsp)
// }

// HTTPPolicyDataSubsToNotifySubsIdPut -
// func HTTPPolicyDataSubsToNotifySubsIdPut(c *gin.Context) {
// 	var policyDataSubscription models.PolicyDataSubscription

// 	if err := getDataFromRequestBody(c, &policyDataSubscription); err != nil {
// 		return
// 	}

// 	req := httpwrapper.NewRequest(c.Request, policyDataSubscription)
// 	req.Params["subsId"] = c.Params.ByName("subsId")

// 	rsp := producer.HandlePolicyDataSubsToNotifySubsIdPut(req)

// 	sendResponse(c, rsp)
// }

// HTTPPolicyDataUesUeIdAmDataGet -
func HTTPPolicyDataUesUeIdAmDataGet(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	req.Params["ueId"] = c.Params.ByName("ueId")

	rsp := producer.HandlePolicyDataUesUeIdAmDataGet(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataUesUeIdOperatorSpecificDataGet -
func HTTPPolicyDataUesUeIdOperatorSpecificDataGet(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	req.Params["ueId"] = c.Params.ByName("ueId")

	rsp := producer.HandlePolicyDataUesUeIdOperatorSpecificDataGet(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataUesUeIdOperatorSpecificDataPatch - Need to be fixed
func HTTPPolicyDataUesUeIdOperatorSpecificDataPatch(c *gin.Context) {
	var patchItemArray []models.PatchItem

	if err := getDataFromRequestBody(c, &patchItemArray); err != nil {
		return
	}

	req := httpwrapper.NewRequest(c.Request, patchItemArray)
	req.Params["ueId"] = c.Params.ByName("ueId")

	rsp := producer.HandlePolicyDataUesUeIdOperatorSpecificDataPatch(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataUesUeIdOperatorSpecificDataPut -
func HTTPPolicyDataUesUeIdOperatorSpecificDataPut(c *gin.Context) {
	var operatorSpecificDataContainerMap map[string]models.OperatorSpecificDataContainer

	if err := getDataFromRequestBody(c, &operatorSpecificDataContainerMap); err != nil {
		return
	}

	req := httpwrapper.NewRequest(c.Request, operatorSpecificDataContainerMap)
	req.Params["ueId"] = c.Params.ByName("ueId")

	rsp := producer.HandlePolicyDataUesUeIdOperatorSpecificDataPut(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataUesUeIdSmDataGet -
func HTTPPolicyDataUesUeIdSmDataGet(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	req.Params["ueId"] = c.Params.ByName("ueId")

	rsp := producer.HandlePolicyDataUesUeIdSmDataGet(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataUesUeIdSmDataPatch - Need to be fixed
func HTTPPolicyDataUesUeIdSmDataPatch(c *gin.Context) {
	var usageMonDataMap map[string]models.UsageMonData

	if err := getDataFromRequestBody(c, &usageMonDataMap); err != nil {
		return
	}

	req := httpwrapper.NewRequest(c.Request, usageMonDataMap)
	req.Params["ueId"] = c.Params.ByName("ueId")

	rsp := producer.HandlePolicyDataUesUeIdSmDataPatch(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataUesUeIdSmDataUsageMonIdDelete -
func HTTPPolicyDataUesUeIdSmDataUsageMonIdDelete(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	req.Params["ueId"] = c.Params.ByName("ueId")
	req.Params["usageMonId"] = c.Params.ByName("usageMonId")

	rsp := producer.HandlePolicyDataUesUeIdSmDataUsageMonIdDelete(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataUesUeIdSmDataUsageMonIdGet -
func HTTPPolicyDataUesUeIdSmDataUsageMonIdGet(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	req.Params["ueId"] = c.Params.ByName("ueId")
	req.Params["usageMonId"] = c.Params.ByName("usageMonId")

	rsp := producer.HandlePolicyDataUesUeIdSmDataUsageMonIdGet(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataUesUeIdSmDataUsageMonIdPut -
func HTTPPolicyDataUesUeIdSmDataUsageMonIdPut(c *gin.Context) {
	var usageMonData models.UsageMonData

	if err := getDataFromRequestBody(c, &usageMonData); err != nil {
		return
	}

	req := httpwrapper.NewRequest(c.Request, usageMonData)
	req.Params["ueId"] = c.Params.ByName("ueId")
	req.Params["usageMonId"] = c.Params.ByName("usageMonId")

	rsp := producer.HandlePolicyDataUesUeIdSmDataUsageMonIdPut(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataUesUeIdUePolicySetGet -
func HTTPPolicyDataUesUeIdUePolicySetGet(c *gin.Context) {
	req := httpwrapper.NewRequest(c.Request, nil)
	req.Params["ueId"] = c.Params.ByName("ueId")

	rsp := producer.HandlePolicyDataUesUeIdUePolicySetGet(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataUesUeIdUePolicySetPatch -
func HTTPPolicyDataUesUeIdUePolicySetPatch(c *gin.Context) {
	var uePolicySet models.UePolicySet

	if err := getDataFromRequestBody(c, &uePolicySet); err != nil {
		return
	}

	req := httpwrapper.NewRequest(c.Request, uePolicySet)
	req.Params["ueId"] = c.Params.ByName("ueId")

	rsp := producer.HandlePolicyDataUesUeIdUePolicySetPatch(req)

	sendResponse(c, rsp)
}

// HTTPPolicyDataUesUeIdUePolicySetPut -
func HTTPPolicyDataUesUeIdUePolicySetPut(c *gin.Context) {
	var uePolicySet models.UePolicySet

	if err := getDataFromRequestBody(c, &uePolicySet); err != nil {
		return
	}

	req := httpwrapper.NewRequest(c.Request, uePolicySet)
	req.Params["ueId"] = c.Params.ByName("ueId")

	rsp := producer.HandlePolicyDataUesUeIdUePolicySetPut(req)

	sendResponse(c, rsp)
}