package routes

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterDTERoutes(r *mux.Router, h *handlers.DTEHandler) {
	// Rutas específicas para creación de DTEs (para Swagger)
	r.HandleFunc("/dte/invoices", h.GenericHandler.CreateInvoice).Methods(http.MethodPost)
	r.HandleFunc("/dte/ccf", h.GenericHandler.CreateCCF).Methods(http.MethodPost)
	r.HandleFunc("/dte/creditnote", h.GenericHandler.CreateCreditNote).Methods(http.MethodPost)
	r.HandleFunc("/dte/retention", h.GenericHandler.CreateRetention).Methods(http.MethodPost)
	
	// Rutas de consulta de DTE e Invalidación
	r.HandleFunc("/dte/invalidation", h.InvalidateDocument).Methods(http.MethodPost)
	r.HandleFunc("/dte/{id}", h.GetByGenerationCode).Methods(http.MethodGet)
	r.HandleFunc("/dte", h.GetAll).Methods(http.MethodGet)
}