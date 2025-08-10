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
// @Summary Crear factura electrónica
// @Description Crea una nueva factura electrónica DTE tipo 01
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
// @Summary Crear Comprobante de Crédito Fiscal
// @Description Crea un nuevo CCF electrónico DTE tipo 03
// @Tags DTE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param ccf body map[string]interface{} true "Datos del CCF"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} response.APIError
// @Failure 401 {object} response.APIError
// @Failure 500 {object} response.APIError
// @Router /dte/ccf [post]
func (h *GenericCreatorDTEHandler) CreateCCF(w http.ResponseWriter, r *http.Request) {
	h.HandleCreate(w, r)
}

// CreateCreditNote godoc
// @Summary Crear Nota de Crédito
// @Description Crea una nueva Nota de Crédito electrónica DTE tipo 05
// @Tags DTE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param creditnote body map[string]interface{} true "Datos de la nota de crédito"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} response.APIError
// @Failure 401 {object} response.APIError
// @Failure 500 {object} response.APIError
// @Router /dte/creditnote [post]
func (h *GenericCreatorDTEHandler) CreateCreditNote(w http.ResponseWriter, r *http.Request) {
	h.HandleCreate(w, r)
}

// CreateRetention godoc
// @Summary Crear Comprobante de Retención
// @Description Crea un nuevo Comprobante de Retención electrónico DTE tipo 06
// @Tags DTE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param retention body map[string]interface{} true "Datos del comprobante de retención"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} response.APIError
// @Failure 401 {object} response.APIError
// @Failure 500 {object} response.APIError
// @Router /dte/retention [post]
func (h *GenericCreatorDTEHandler) CreateRetention(w http.ResponseWriter, r *http.Request) {
	h.HandleCreate(w, r)
}

// HandleCreate maneja la creación de cualquier tipo de documento
func (h *GenericCreatorDTEHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	// 1. Determinar qué tipo de documento se está creando basado en la ruta
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

	// 4. Invocar el caso de uso genérico
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

	// 6. Responder con éxito
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

	// 2. Actualizar la identificación de contingencia en el JSON del DTE
	updatedDTE, err := utils.UpdateContingencyIdentification(dte, contiType, reason)
	if err != nil {
		return err
	}

	// 3. Responder con la respuesta de la creación del DTE
	h.respWriter.Success(w, http.StatusCreated, updatedDTE, options)
	return nil
}

// GetDocumentConfigs devuelve la configuración de documentos registrada
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