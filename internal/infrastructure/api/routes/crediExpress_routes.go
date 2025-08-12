																
package routes

import (
	"github.com/gorilla/mux"
	"github.com/MarlonG1/api-facturacion-sv/internal/infrastructure/api/controllers/credi_express/cliente_controller"
)

// SetupPagosRoutes configura las rutas para el módulo de pagos
func SetupPagosRoutes(router *mux.Router) {
	// Crear controlador de clientes
	clienteController := pagos.NewClienteController()

	// Subrutador para el módulo de pagos
	pagosRouter := router.PathPrefix("/api/pagos").Subrouter()

	// Rutas de clientes
	clientesRouter := pagosRouter.PathPrefix("/clientes").Subrouter()
	
	// GET /api/pagos/clientes - Obtener todos los clientes
	clientesRouter.HandleFunc("", clienteController.GetAllClientes).Methods("GET")
	
	// GET /api/pagos/clientes/activos - Obtener clientes activos
	clientesRouter.HandleFunc("/activos", clienteController.GetActiveClientes).Methods("GET")
	
	// GET /api/pagos/clientes/search - Buscar clientes
	clientesRouter.HandleFunc("/search", clienteController.SearchClientes).Methods("GET")
	
	// GET /api/pagos/clientes/dui/{dui} - Obtener cliente por DUI
	clientesRouter.HandleFunc("/dui/{dui}", clienteController.GetClienteByDUI).Methods("GET")
	
	// GET /api/pagos/clientes/{id} - Obtener cliente por ID (debe ir al final)
	clientesRouter.HandleFunc("/{id:[0-9]+}", clienteController.GetClienteByID).Methods("GET")
}

