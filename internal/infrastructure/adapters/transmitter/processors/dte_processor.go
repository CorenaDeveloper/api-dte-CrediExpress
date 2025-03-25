package processors

import (
	"github.com/MarlonG1/api-facturacion-sv/config"
	models2 "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/transmitter/models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/shared_error"
)

type DTEProcessor struct{}

func (p *DTEProcessor) ProcessRequest(signedDoc string, document interface{}) (*models2.HaciendaRequest, error) {
	version, dteType, generationCode, sequenceNumber, err := GetDocumentRequestData(document)
	if err != nil {
		logs.Error("Failed to get document request data", map[string]interface{}{"error": err.Error()})
		return nil, err
	}

	return &models2.HaciendaRequest{
		Ambient:        config.Server.AmbientCode,
		SendID:         sequenceNumber,
		Version:        version,
		Document:       signedDoc,
		DTEType:        dteType,
		GenerationCode: generationCode,
		URL:            config.MHPaths.ReceptionURL,
	}, nil
}

func (p *DTEProcessor) ProcessResponse(resp *models2.HaciendaResponse) (*models2.TransmitResult, error) {
	if resp == nil {
		return nil, shared_error.NewGeneralServiceError("InvalidationProcessor", "ProcessResponse", "nil response", nil)
	}

	return &models2.TransmitResult{
		Status:         resp.Status,
		ReceptionStamp: &resp.ReceptionStamp,
		ProcessingDate: resp.ProcessingDate,
		MessageCode:    resp.MessageCode,
		MessageDesc:    resp.DescriptionMessage,
		Observations:   resp.Observations,
	}, nil
}
