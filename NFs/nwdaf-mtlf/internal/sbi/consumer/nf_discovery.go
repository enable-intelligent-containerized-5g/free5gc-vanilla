package consumer

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/enable-intelligent-containerized-5g/openapi/Nnrf_NFDiscovery"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/logger"
)

func sendSearchNfInstances(nrfUri string, targetNfType, requestNfType models.NfType,
	param Nnrf_NFDiscovery.SearchNFInstancesParamOpts,
) (models.SearchResult, error) {
	// Set client and set url
	configuration := Nnrf_NFDiscovery.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

	var res *http.Response
	result, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requestNfType, &param)
	if res != nil && res.StatusCode == http.StatusTemporaryRedirect {
		err = fmt.Errorf("temporary redirect for Non NRF consumer")
	}
	defer func() {
		if rspCloseErr := res.Body.Close(); rspCloseErr != nil {
			logger.ConsumerLog.Errorf("SearchNFInstances response body cannot close: %+v", rspCloseErr)
		}
	}()

	return result, err
}

// func SendSearchNFInstances(nrfUri string, targetNfType, requestNfType models.NfType,
// 	param Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (models.SearchResult, error) {
// 	// Set client and set url
// 	configuration := Nnrf_NFDiscovery.NewConfiguration()
// 	configuration.SetBasePath(nrfUri)
// 	client := Nnrf_NFDiscovery.NewAPIClient(configuration)

// 	var res *http.Response
// 	result, res, err := client.NFInstancesStoreApi.SearchNFInstances(context.TODO(), targetNfType, requestNfType, &param)
// 	if res != nil && res.StatusCode == http.StatusTemporaryRedirect {
// 		err = fmt.Errorf("temporary redirect for Non NRF consumer")
// 	}
// 	defer func() {
// 		if rspCloseErr := res.Body.Close(); rspCloseErr != nil {
// 			logger.ConsumerLog.Errorf("SearchNFInstances response body cannot close: %+v", rspCloseErr)
// 		}
// 	}()

// 	return result, err
// }

func SearchAllNfInstance(nfInstances *[]models.NfProfile, nrfUri string, targetNfType, requestNfType models.NfType,
	param Nnrf_NFDiscovery.SearchNFInstancesParamOpts,
) error {
	resp, localErr := sendSearchNfInstances(nrfUri, targetNfType, requestNfType, param)
	if localErr != nil {
		return localErr
	}

	*nfInstances = append(*nfInstances, resp.NfInstances...)

	if len(*nfInstances) == 0 {
		return errors.New("no NF instances found")
	}

	return nil
}
