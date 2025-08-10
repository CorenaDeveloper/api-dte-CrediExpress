package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/helpers"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// GenericCreatorDTEHandler maneja las solicitudes para crear cualquier tipo de documento DTE
type GenericCreatorDTEHandler struct {
	documentConfigs    map[string]helpers.DocumentConfig
	respWriter         *response.ResponseWriter
	contingencyHandler *helpers.ContingencyHandler
}

// NewGenericDTEHandler crea una nueva instancia de GenericCreatorDTEHandler
func NewGenericDTEHandler(contingencyHandler *helpers.ContingencyHandler) *GenericCreatorDTEHandler {
	return &GenericCreatorDTEHandler{
		documentConfigs:    make(map[string]helpers.DocumentConfig),
		respWriter:         response.NewResponseWriter(),
		contingencyHandler: contingencyHandler,
	}
}

// RegisterDocument registra un nuevo tipo de documento para ser manejado
func (h *GenericCreatorDTEHandler) RegisterDocument(path string, config helpers.DocumentConfig) {
	h.documentConfigs[path] = config
}

// CreateInvoice godoc
// @Summary Crear factura electrÃ³nica
// @Description Crea una nueva factura electrÃ³nica DTE tipo 01
// @Tags DTE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param invoice body map[string]interface{} true "Datos de la factura"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} response.APIError
// @Failure 401 {object} response.APIError
// @Failure 500 {object} response.APIError
// @Router /dte/invoices [post]
func (h *GenericCreatorDTEHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	h.HandleCreate(w, r)
}

// CreateCCF godoc
// @Summary Crear Comprobante de CrÃ©dito Fiscal
// @Description // @Description Este endpoint permite crear y emitir una Factura ElectrÃ³nica.
// @Description 
// @Description ## Ejemplo de Solicitud
// @Description ```json
// @Description 
// @Description ```
// @Description ## Ejemplo de Respuesta
// @Description ```json
// @Description 
// @Description ```
// @Description Para ver ejemplos completos, consulta: /jsonExamples/
// @Tags DTE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param ccf body object true "Datos de CCF"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} response.APIError
// @Failure 401 {object} response.APIError
// @Failure 500 {object} response.APIError
// @Router /dte/ccf [post]
func (h *GenericCreatorDTEHandler) CreateCCF(w http.ResponseWriter, r *http.Request) {
	h.HandleCreate(w, r)
}

// CreateCreditNote godoc
// @Summary Crear Nota de CrÃ©dito
// @Description // @Description Este endpoint permite crear y emitir una Factura ElectrÃ³nica.
// @Description 
// @Description ## Ejemplo de Solicitud
// @Description ```json
// @Description 
// @Description ```
// @Description ## Ejemplo de Respuesta
// @Description ```json
// @Description 
// @Description ```
// @Description Para ver ejemplos completos, consulta: /jsonExamples/
// @Tags DTE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param creditnote body object true "Datos de la nota de crÃ©dito"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} response.APIError
// @Failure 401 {object} response.APIError
// @Failure 500 {object} response.APIError
// @Router /dte/creditnote [post]
func (h *GenericCreatorDTEHandler) CreateCreditNote(w http.ResponseWriter, r *http.Request) {
	h.HandleCreate(w, r)
}

// CreateRetention godoc
// @Summary Crear Comprobante de RetenciÃ³n
// @Description // @Description Este endpoint permite crear y emitir una Factura ElectrÃ³nica.
// @Description 
// @Description ## Ejemplo de Solicitud
// @Description ```json
// @Description 
// @Description ```
// @Description ## Ejemplo de Respuesta
// @Description ```json
// @Description 
// @Description ```
// @Description Para ver ejemplos completos, consulta: /jsonExamples/
// @Tags DTE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param retention body object true "Datos del comprobante de retenciÃ³n"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} response.APIError
// @Failure 401 {object} response.APIError
// @Failure 500 {object} response.APIError
// @Router /dte/retention [post]
func (h *GenericCreatorDTEHandler) CreateRetention(w http.ResponseWriter, r *http.Request) {
	h.HandleCreate(w, r)
}

// HandleCreate maneja la creaciÃ³n de cualquier tipo de documento
func (h *GenericCreatorDTEHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// 1. Determinar quÃ© tipo de documento se estÃ¡ creando basado en la ruta
	path := r.URL.Path
	config, err := h.getDocumentTypeFromPath(path)
	if err != nil {
		h.respWriter.Error(w, http.StatusNotFound, "Document type not supported", nil)
		return
	}

	// 2. Crear una nueva instancia del tipo de solicitud
	requestType := reflect.TypeOf(config.RequestType)
	request := reflect.New(requestType.Elem()).Interface()

	// 3. Decodificar el JSON en la estructura de solicitud
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		logs.Error("Failed to decode request body", map[string]interface{}{"error": err.Error()})
		h.respWriter.Error(w, http.StatusBadRequest, "Invalid request format", nil)
		return
	}

	// 4. Invocar el caso de uso genÃ©rico
	resp, options, err := config.UseCase.Create(r.Context(), request)
	if err != nil {
		logs.Warn("Error processing document because", map[string]interface{}{"error": err.Error()})

		// 5. Si aplica contingencia, manejarla
		if config.UsesContingency {
			err = h.handleErrorForContingency(r.Context(), resp, config.DocumentType, options, err, w)
			if err != nil {
				h.respWriter.HandleError(w, err)
				return
			}
			return
		} else {
			h.respWriter.HandleError(w, err)
			return
		}
	}

	// 6. Responder con Ã©xito
	h.respWriter.Success(w, http.StatusCreated, resp, options)
}

// handleErrorForContingency maneja el error en caso de que se aplique una contingencia
func (h *GenericCreatorDTEHandler) handleErrorForContingency(ctx context.Context, dte interface{}, dteType string, options *response.SuccessOptions, err error, w http.ResponseWriter) error {
	// 1. Verificar si aplica a contingencia
	logs.Warn("Error transmitting DTE because", map[string]interface{}{
		"error": err.Error(),
	})

	contiType, reason := h.contingencyHandler.HandleContingency(ctx, dte, dteType, err)
	if contiType == nil || reason == nil {
		logs.Error("Error creating DTE contingency", map[string]interface{}{"error": err.Error()})
		return err
	}

	// 2. Actualizar la identificaciÃ³n de contingencia en el JSON del DTE
	updatedDTE, err := utils.UpdateContingencyIdentification(dte, contiType, reason)
	if err != nil {
		return err
	}

	// 3. Responder con la respuesta de la creaciÃ³n del DTE
	h.respWriter.Success(w, http.StatusCreated, updatedDTE, options)
	return nil
}

// GetDocumentConfigs devuelve la configuraciÃ³n de documentos registrada
func (h *GenericCreatorDTEHandler) GetDocumentConfigs() map[string]helpers.DocumentConfig {
	return h.documentConfigs
}

// getDocumentTypeFromPath obtiene el tipo de documento basado en la ruta
func (h *GenericCreatorDTEHandler) getDocumentTypeFromPath(path string) (helpers.DocumentConfig, error) {
	for key := range h.documentConfigs {
		if strings.Contains(path, key) {
			return h.documentConfigs[key], nil
		}
	}

	return helpers.DocumentConfig{}, errors.New("error was found in the path")
}



