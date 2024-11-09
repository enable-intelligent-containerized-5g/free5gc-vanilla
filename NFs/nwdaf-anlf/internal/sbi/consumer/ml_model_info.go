package consumer

import (
	"context"
	"fmt"
	"net/http"

	Nnwaf_MLModelInfo "github.com/enable-intelligent-containerized-5g/openapi/Nnwdaf_MLModelInfo"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
)

func sendGetMlModelInfoList(mtlfUri string) (models.MlModelDataResponse, error) {
	configuration := Nnwaf_MLModelInfo.NewConfiguration()
	configuration.SetBasePath(mtlfUri)
	client := Nnwaf_MLModelInfo.NewAPIClient(configuration)

	var res *http.Response
	result, res, err := client.MLModelInfoStoreApi.SearhMlModelInfoList(context.TODO())

	if res == nil {
		return result, fmt.Errorf(err.Error())
	}

	defer func() {
		if rspCloseErr := res.Body.Close(); rspCloseErr != nil {
			logger.ConsumerLog.Errorf("SearchMlModelInfoList response body cannot close: %+v", rspCloseErr)
		}
	}()
	
	return result, nil
}

func SendGetMlModelInfoList(mlModelInfoList *[]models.MlModelData, mtlfUri string) error {
	resp, localErr := sendGetMlModelInfoList(mtlfUri)

	if localErr != nil {
		return localErr
	}

	*mlModelInfoList = resp.MlModels

	return nil
}
