package mlmodeltraining

import (
	"github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/internal/sbi/producer"
	"github.com/free5gc/util/httpwrapper"
	"github.com/gin-gonic/gin"
	"net/http"

	// "database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type MlModelInfo struct {
	URI          string `json:"uri,omitempty" yaml:"uri" bson:"uri" mapstructure:"uri"`
	Accuracy     string `json:"accuracy,omitempty" yaml:"accuracy" bson:"accuracy" mapstructure:"accuracy"`
	NF           string `json:"nf,omitempty" yaml:"nf" bson:"nf" mapstructure:"nf"`
	TargetPeriod string `json:"targetPeriod,omitempty" yaml:"targetPeriod" bson:"targetPeriod" mapstructure:"targetPeriod"`
	EventId      string `json:"eventId,omitempty" yaml:"eventId" bson:"eventId" mapstructure:"eventId"`
}

func HTTPSaveMlModel(c *gin.Context) {
	var modelInfo MlModelInfo

    // Step 1: Retrieve HTTP request body
    requestBody, err := c.GetRawData()
    if err != nil {
        problemDetail := models.ProblemDetails{
            Title:  "System failure",
            Status: http.StatusInternalServerError,
            Detail: err.Error(),
            Cause:  "SYSTEM_FAILURE",
        }
        logger.ManagementLog.Errorf("Get Request Body error: %+v", err)
        c.JSON(http.StatusInternalServerError, problemDetail)
        return
    }

    // Step 2: Convert requestBody to MlModelInfo model
    err = openapi.Deserialize(&modelInfo, requestBody, "application/json")
    if err != nil {
        problemDetail := "[Request Body] " + err.Error()
        rsp := models.ProblemDetails{
            Title:  "Malformed request syntax",
            Status: http.StatusBadRequest,
            Detail: problemDetail,
        }
        logger.ManagementLog.Errorln(problemDetail)
        c.JSON(http.StatusBadRequest, rsp)
        return
    }

    logger.UtilLog.Warn("Deserialized ModelInfo: ", modelInfo)

    req := httpwrapper.NewRequest(c.Request, modelInfo)
    httpResponse := producer.HandleSaveMlModel(req)
    responseBody, err := openapi.Serialize(httpResponse.Body, "application/json")
    if err != nil {
        logger.ManagementLog.Errorln(err)
        problemDetails := models.ProblemDetails{
            Status: http.StatusInternalServerError,
            Cause:  "SYSTEM_FAILURE",
            Detail: err.Error(),
        }
        c.JSON(http.StatusInternalServerError, problemDetails)
    } else {
        c.Data(httpResponse.Status, "application/json", responseBody)
    }
}