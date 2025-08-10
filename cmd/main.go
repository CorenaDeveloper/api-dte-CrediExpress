// cmd/main.go

// @title API Facturación El Salvador
// @version 1.0
// @swagger 2.0
// @description API para la gestión de Documentos Tributarios Electrónicos (DTE) que cumple con los requisitos establecidos por la autoridad fiscal de El Salvador.
// @termsOfService http://swagger.io/terms/

// @contact.name Soporte API DTE
// @contact.url https://github.com/CorenaDeveloper/api-dte-CrediExpress
// @contact.email soporte@empresa.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host 185.222.241.94:7319
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/bootstrap"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
	"os"
)

func main() {
	// Crear e inicializar la aplicación
	app := bootstrap.NewApplication()
	if err := app.Initialize(); err != nil {
		logs.Fatal("Failed to initialize application", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}

	// Iniciar la aplicación
	if err := app.Start(); err != nil {
		logs.Fatal("Application error", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
}