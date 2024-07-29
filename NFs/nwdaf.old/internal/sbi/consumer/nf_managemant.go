package consumer

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	nwdaf_context "github.com/free5gc/nwdaf/internal/context"
	"github.com/free5gc/nwdaf/internal/logger"
	"github.com/free5gc/nwdaf/pkg/factory"
	"github.com/enable-intelligent-containerized-5g/openapi"
	"github.com/enable-intelligent-containerized-5g/openapi/Nnrf_NFManagement"
	"github.com/enable-intelligent-containerized-5g/openapi/models"
)

func BuildNFInstance(context *nwdaf_context.NWDAFContext) models.NfProfile {
	config := factory.NwdafConfig

	profile := models.NfProfile{
		NfInstanceId:  context.NfId,
		NfType:        models.NfType_NWDAF,
		NfStatus:      models.NfStatus_REGISTERED,
		Ipv4Addresses: []string{context.RegisterIPv4},
		// NwdafInfo: &models.NwdafInfo{
		// 	SupportedDataSets: []models.DataSetId{
		// 		// models.DataSetId_APPLICATION,
		// 		// models.DataSetId_EXPOSURE,
		// 		// models.DataSetId_POLICY,
		// 		models.DataSetId_SUBSCRIPTION,
		// 	}
		// },
	}

	version := config.Info.Version
	tmpVersion := strings.Split(version, ".")
	versionUri := "v" + tmpVersion[0]
	apiPrefix := fmt.Sprintf("%s://%s:%d", context.UriScheme, context.RegisterIPv4, context.SBIPort)
	profile.NfServices = &[]models.NfService{
		{
			ServiceInstanceId: "datarepository",
			ServiceName:       models.ServiceName_NNWDAF_ANALYTICSINFO,
			Versions: &[]models.NfServiceVersion{
				{
					ApiFullVersion:  version,
					ApiVersionInUri: versionUri,
				},
			},
			Scheme:          context.UriScheme,
			NfServiceStatus: models.NfServiceStatus_REGISTERED,
			ApiPrefix:       apiPrefix,
			IpEndPoints: &[]models.IpEndPoint{
				{
					Ipv4Address: context.RegisterIPv4,
					Transport:   models.TransportProtocol_TCP,
					Port:        int32(context.SBIPort),
				},
			},
		},
	}

	// TODO: finish the Nwdaf Info
	return profile
}

func SendRegisterNFInstance(nrfUri, nfInstanceId string, profile models.NfProfile) (string, string, error) {
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)
	var resouceNrfUri string
	var retrieveNfInstanceId string 

	for {
		_, res, err := client.NFInstanceIDDocumentApi.RegisterNFInstance(context.TODO(), nfInstanceId, profile)
		if err != nil || res == nil {
			// TODO : add log
			fmt.Println(fmt.Errorf("NWDAF register to NRF Error[%s]", err.Error()))
			time.Sleep(2 * time.Second)
			continue
		}
		defer func() {
			if rspCloseErr := res.Body.Close(); rspCloseErr != nil {
				logger.ConsumerLog.Errorf("RegisterNFInstance response body cannot close: %+v", rspCloseErr)
			}
		}()

		status := res.StatusCode
		if status == http.StatusOK {
			// NFUpdate
			return resouceNrfUri, retrieveNfInstanceId, err
		} else if status == http.StatusCreated {
			// NFRegister
			resourceUri := res.Header.Get("Location")
			resouceNrfUri = resourceUri[:strings.Index(resourceUri, "/nnrf-nfm/")]
			retrieveNfInstanceId = resourceUri[strings.LastIndex(resourceUri, "/")+1:]
			return resouceNrfUri, retrieveNfInstanceId, err
		} else {
			fmt.Println("handler returned wrong status code", status)
			fmt.Println("NRF return wrong status code", status)
		}
	}
}

func SendDeregisterNFInstance() (problemDetails *models.ProblemDetails, err error) {
	logger.ConsumerLog.Infof("Send Deregister NFInstance")

	nwdafSelf := nwdaf_context.NWDAF_Self()
	// Set client and set url
	configuration := Nnrf_NFManagement.NewConfiguration()
	configuration.SetBasePath(nwdafSelf.NrfUri)
	client := Nnrf_NFManagement.NewAPIClient(configuration)

	var res *http.Response

	res, err = client.NFInstanceIDDocumentApi.DeregisterNFInstance(context.Background(), nwdafSelf.NfId)
	if err == nil {
		return
	} else if res != nil {
		defer func() {
			if rspCloseErr := res.Body.Close(); rspCloseErr != nil {
				logger.ConsumerLog.Errorf("DeregisterNFInstance response body cannot close: %+v", rspCloseErr)
			}
		}()

		if res.Status != err.Error() {
			return
		}
		problem := err.(openapi.GenericOpenAPIError).Model().(models.ProblemDetails)
		problemDetails = &problem
	} else {
		err = openapi.ReportError("server no response")
	}
	return
}
