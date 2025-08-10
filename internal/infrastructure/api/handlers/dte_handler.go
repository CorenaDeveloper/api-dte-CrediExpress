package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MarlonG1/api-facturacion-sv/internal/application/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/helpers"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/response"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

type DTEHandler struct {
	GenericHandler      *GenericCreatorDTEHandler
	dteConsultUseCase   *dte.DTEConsultUseCase
	invalidationUseCase *dte.InvalidationUseCase
	respWriter          *response.ResponseWriter
}

func NewDTEHandler(
	dteConsultUseCase *dte.DTEConsultUseCase,
	invalidationUseCase *dte.InvalidationUseCase,
	genericHandler *GenericCreatorDTEHandler,
) *DTEHandler {
	return &DTEHandler{
		GenericHandler:      genericHandler,
		dteConsultUseCase:   dteConsultUseCase,
		invalidationUseCase: invalidationUseCase,
		respWriter:          response.NewResponseWriter(),
	}
}

// NOTA IMPORTANTE:
// Ahora la creacion de un DTE se maneja en el GenericCreatorDTEHandler y sus rutas en el router

// GetByGenerationCode maneja la solicitud HTTP para obtener un DTE por su código de generación
// GetByGenerationCode godoc
// @Summary Obtener DTE por código de generación
// @Description Obtiene un DTE específico por su código de generación
// @Tags DTE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param id path string true "Código de generación del DTE" format(uuid)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.APIError
// @Failure 401 {object} response.APIError
// @Failure 404 {object} response.APIError
// @Failure 500 {object} response.APIError
// @Router /dte/{id} [get]
func (h *DTEHandler) GetByGenerationCode(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener el código de generación
	generationCode := helpers.GetRequestVar(r, "id")

	// 2. Obtener DTE ejecutando el caso de uso
	dte, err := h.dteConsultUseCase.GetByGenerationCode(r.Context(), generationCode)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, dte, nil)
}

// GetAll maneja la solicitud HTTP para obtener todos los DTEs
// GetAll godoc
// @Summary Listar DTEs
// @Description Obtiene lista paginada de DTEs del usuario autenticado
// @Tags DTE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param page query int false "Número de página" default(1)
// @Param limit query int false "Elementos por página" default(10)
// @Param status query string false "Filtrar por estado" Enums(PROCESSED,PENDING,ERROR)
// @Param dte_type query string false "Filtrar por tipo DTE" Enums(01,03,05,06,11,14)
// @Param date_from query string false "Fecha inicio (YYYY-MM-DD)"
// @Param date_to query string false "Fecha fin (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.APIError
// @Failure 401 {object} response.APIError
// @Failure 500 {object} response.APIError
// @Router /dte [get]
func (h *DTEHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// 1. Obtener todos los DTEs ejecutando el caso de uso
	dtes, err := h.dteConsultUseCase.GetAllDTEs(r.Context(), r)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, dtes, nil)
}

// InvalidateDocument maneja la solicitud HTTP para invalidar un DTE
// InvalidateDocument godoc
// @Summary Invalidar DTE
// @Description Invalida un DTE existente
// @Tags DTE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param invalidation body map[string]interface{} true "Datos para invalidación"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.APIError
// @Failure 401 {object} response.APIError
// @Failure 500 {object} response.APIError
// @Router /dte/invalidation [post]
func (h *DTEHandler) InvalidateDocument(w http.ResponseWriter, r *http.Request) {
	// 1. Decodificar la solicitud de invalidación de documento a un DTO de solicitud
	var req structs.CreateInvalidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logs.Error("Failed to decode request body", map[string]interface{}{"error": err.Error()})
		h.respWriter.Error(w, http.StatusBadRequest, "Invalid request format", nil)
		return
	}

	// 2. Ejecutar el caso de uso de invalidación de documento
	invalidation, err := h.invalidationUseCase.InvalidateDocument(r.Context(), req)
	if err != nil {
		h.respWriter.HandleError(w, err)
		return
	}

	h.respWriter.Success(w, http.StatusOK, invalidation, nil)
}
