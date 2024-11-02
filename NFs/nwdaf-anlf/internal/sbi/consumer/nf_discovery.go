package consumer

import (
	"context"
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/enable-intelligent-containerized-5g/openapi/Nnrf_NFDiscovery"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
	"github.com/free5gc/nwdaf/internal/util"
)

func SendSearchNFInstances(nrfUri string, targetNfType, requestNfType models.NfType,
	param Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (models.SearchResult, error) {
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

func SearchAllNfInstance(nrfUri string, targetNfType, requestNfType models.NfType,
	param Nnrf_NFDiscovery.SearchNFInstancesParamOpts) error {
	resp, localErr := SendSearchNFInstances(nrfUri, targetNfType, requestNfType, param)

	// Convertir `allNfs` a string en formato JSON
    allNfsJSON, errJson := json.Marshal(resp)
	if errJson != nil {
        fmt.Println("Error al convertir a JSON:", localErr)
        return nil
    }

	if localErr != nil {
		return localErr
	}

	// Imprimir la representación en string de `allNfs`
	fmt.Println("resp", string(allNfsJSON))

	return nil
}

func SearchAnaliticsInfoInstance(nrfUri string, targetNfType, requestNfType models.NfType,
	param Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (err error) {
	resp, err := SendSearchNFInstances(nrfUri, targetNfType, requestNfType, param)

	// Convertir `allNfs` a string en formato JSON
    // allNfsJSON, errJson := json.Marshal(resp)
	// if errJson != nil {
    //     fmt.Println("Error al convertir a JSON:", localErr)
    //     return nil
    // }

	// if localErr != nil {
	// 	return localErr
	// }

	// Imprimir la representación en string de `allNfs`
	// fmt.Println("resp", string(allNfsJSON))

	// select the first AMF, TODO: select base on other info
	var amfUri string
	// var TargetMtlfProfile models.NfProfile // Reemplazo de ue.TargetAmfProfile
	var TargetMtlfUri string // Reemplazo de ue.TargetAmfUri
	for _, nfProfile := range resp.NfInstances {
		// ue.TargetAmfProfile = &nfProfile
		// TargetMtlfProfile = nfProfile
		amfUri = util.SearchNFServiceUri(nfProfile, models.ServiceName_NNWDAF_ANALYTICSINFO, models.NfServiceStatus_REGISTERED)
		if amfUri != "" {
			break
		}
	}

	// fmt.Println("Profile Found: ", TargetMtlfProfile)

	// ue.TargetAmfUri = amfUri
	TargetMtlfUri = amfUri
	if TargetMtlfUri == "" {
		err = fmt.Errorf("NF can not select an target AnLF by NRF")
	}

	fmt.Println("TargetMtlfUri found: ", TargetMtlfUri)
	return
}

func SearchMlModelInfo(nrfUri string, targetNfType, requestNfType models.NfType,
	param Nnrf_NFDiscovery.SearchNFInstancesParamOpts) (err error) {
	resp, err := SendSearchNFInstances(nrfUri, targetNfType, requestNfType, param)

	// Convertir `allNfs` a string en formato JSON
    // allNfsJSON, errJson := json.Marshal(resp)
	// if errJson != nil {
    //     fmt.Println("Error al convertir a JSON:", localErr)
    //     return nil
    // }

	// if localErr != nil {
	// 	return localErr
	// }

	// Imprimir la representación en string de `allNfs`
	// fmt.Println("resp", string(allNfsJSON))

	// select the first AMF, TODO: select base on other info
	var amfUri string
	// var TargetMtlfProfile models.NfProfile // Reemplazo de ue.TargetAmfProfile
	var TargetMtlfUri string // Reemplazo de ue.TargetAmfUri
	for _, nfProfile := range resp.NfInstances {
		// ue.TargetAmfProfile = &nfProfile
		// TargetMtlfProfile = nfProfile
		amfUri = util.SearchNFServiceUri(nfProfile, models.ServiceName_NNWDAF_MLMODELPROVISION, models.NfServiceStatus_REGISTERED)
		if amfUri != "" {
			break
		}
	}

	// fmt.Println("Profile Found: ", TargetMtlfProfile)

	// ue.TargetAmfUri = amfUri
	TargetMtlfUri = amfUri
	if TargetMtlfUri == "" {
		err = fmt.Errorf("NF can not select an target MTLF by NRF")
	}

	fmt.Println("TargetMtlfUri found: ", TargetMtlfUri)
	return
}
