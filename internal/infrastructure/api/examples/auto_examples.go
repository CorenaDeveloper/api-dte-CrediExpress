// 1. Crear archivo: internal/infrastructure/api/examples/auto_examples.go

package examples

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Estructura para mapear endpoints a sus archivos JSON
var EndpointExamples = map[string]ExampleConfig{
	"ccf": {
		RequestFile:  "jsonExamples/ccf_request.json",
		ResponseFile: "jsonExamples/ccf_response.json",
		Title:        "Comprobante de Crédito Fiscal",
		Description:  "Este endpoint permite crear y emitir un Comprobante de Crédito Fiscal (CCF) electrónico.",
	},
	"invoice": {
		RequestFile:  "jsonExamples/invoice_request.json", 
		ResponseFile: "jsonExamples/invoice_response.json",
		Title:        "Factura Electrónica",
		Description:  "Este endpoint permite crear y emitir una Factura Electrónica.",
	},
	"creditnote": {
		RequestFile:  "jsonExamples/creditnote_request.json",
		ResponseFile: "jsonExamples/creditnote_response.json", 
		Title:        "Nota de Crédito",
		Description:  "Este endpoint permite crear y emitir una Nota de Crédito electrónica.",
	},
	"retention": {
		RequestFile:  "jsonExamples/retention_request.json",
		ResponseFile: "jsonExamples/retention_response.json",
		Title:        "Comprobante de Retención", 
		Description:  "Este endpoint permite crear y emitir un Comprobante de Retención electrónico.",
	},
}

type ExampleConfig struct {
	RequestFile  string
	ResponseFile string
	Title        string
	Description  string
}

// Función que genera automáticamente la descripción Swagger desde JSON
func GenerateSwaggerDescription(endpointKey string) string {
	config, exists := EndpointExamples[endpointKey]
	if !exists {
		return "Endpoint no configurado"
	}

	var parts []string
	
	// Descripción inicial
	parts = append(parts, config.Description)
	parts = append(parts, "")
	
	// Agregar ejemplo de request
	if requestExample := loadAndFormatJSON(config.RequestFile, "Ejemplo de Solicitud"); requestExample != "" {
		parts = append(parts, requestExample)
	}
	
	// Agregar ejemplo de response
	if responseExample := loadAndFormatJSON(config.ResponseFile, "Ejemplo de Respuesta"); responseExample != "" {
		parts = append(parts, responseExample)
	}
	
	// Link a archivos completos
	parts = append(parts, fmt.Sprintf("Para ver ejemplos completos, consulta: /%s", filepath.Dir(config.RequestFile)))
	
	return strings.Join(parts, "\n// @Description ")
}

// Helper para cargar y formatear JSON
func loadAndFormatJSON(filePath string, title string) string {
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("// Error loading %s", filePath)
	}

	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return "// Error parsing JSON"
	}

	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "// Error formatting JSON"
	}

	lines := strings.Split(string(formatted), "\n")
	var result []string

	result = append(result, fmt.Sprintf("## %s", title))
	result = append(result, "```json")
	for _, line := range lines {
		result = append(result, line)
	}
	result = append(result, "```")
	result = append(result, "")

	return strings.Join(result, "\n// @Description ")
}