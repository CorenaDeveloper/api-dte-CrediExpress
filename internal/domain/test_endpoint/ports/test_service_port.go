package ports

import (
	"context"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/test_endpoint/models"
)

// TestManager es una interfaz que define los métodos que debe implementar un servicio de pruebas del sistema.
type TestManager interface {
	RunSystemTest(ctx context.Context) (*models.TestResult, error) // RunSystemTest ejecuta una serie de pruebas del sistema y retorna un resumen de los resultados.
}
