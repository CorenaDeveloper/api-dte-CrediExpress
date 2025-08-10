package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func processJSON(jsonFile, title string) string {
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		return fmt.Sprintf("Error: No se encontró el archivo %s\n", jsonFile)
	}

	content, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return fmt.Sprintf("Error leyendo %s: %v\n", jsonFile, err)
	}

	return fmt.Sprintf(`## %s
`+"```json\n%s\n```\n\n", title, string(content))
}

func generateDescription(requestFile, responseFile, description string) string {
	result := description + "\n\n"
	result += processJSON(requestFile, "Ejemplo de Solicitud")
	result += processJSON(responseFile, "Ejemplo de Respuesta")
	result += "Para ver ejemplos completos, consulta: /jsonExamples/"
	return result
}

func replaceInFile(filePath, placeholder, newContent string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error leyendo archivo: %v", err)
	}

	// Preparar contenido para Swagger
	lines := strings.Split(newContent, "\n")
	var swaggerLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			swaggerLines = append(swaggerLines, "// @Description "+line)
		} else {
			swaggerLines = append(swaggerLines, "// @Description ")
		}
	}
	swaggerContent := strings.Join(swaggerLines, "\n")

	// Reemplazar placeholder
	newFileContent := strings.ReplaceAll(string(content), "{{"+placeholder+"}}", swaggerContent)

	// Escribir archivo
	err = ioutil.WriteFile(filePath, []byte(newFileContent), 0644)
	if err != nil {
		return fmt.Errorf("error escribiendo archivo: %v", err)
	}

	fmt.Printf("✅ Reemplazado %s\n", placeholder)
	return nil
}

func main() {
	handlerFile := "internal/infrastructure/api/handlers/generic_creator_handler.go"

	fmt.Println("🔄 Iniciando actualización de descripciones Swagger...")

	// Verificar que existe el archivo
	if _, err := os.Stat(handlerFile); os.IsNotExist(err) {
		fmt.Printf("❌ Error: No se encontró %s\n", handlerFile)
		os.Exit(1)
	}

	// Generar descripciones
	endpoints := map[string]struct {
		requestFile, responseFile, description string
	}{
		"CCF_DESCRIPTION": {
			"jsonExamples/ccf_request.json",
			"jsonExamples/ccf_response.json",
			"Este endpoint permite crear y emitir un Comprobante de Crédito Fiscal (CCF) electrónico.",
		},
		"INVOICE_DESCRIPTION": {
			"jsonExamples/invoice_request.json",
			"jsonExamples/invoice_response.json",
			"Este endpoint permite crear y emitir una Factura Electrónica.",
		},
		"CREDITNOTE_DESCRIPTION": {
			"jsonExamples/creditnote_request.json",
			"jsonExamples/creditnote_response.json",
			"Este endpoint permite crear y emitir una Nota de Crédito electrónica.",
		},
		"RETENTION_DESCRIPTION": {
			"jsonExamples/retention_request.json",
			"jsonExamples/retention_response.json",
			"Este endpoint permite crear y emitir un Comprobante de Retención electrónico.",
		},
	}

	// Procesar cada endpoint
	for placeholder, config := range endpoints {
		fmt.Printf("📝 Generando descripción para %s...\n", placeholder)
		desc := generateDescription(config.requestFile, config.responseFile, config.description)
		
		if err := replaceInFile(handlerFile, placeholder, desc); err != nil {
			fmt.Printf("❌ Error procesando %s: %v\n", placeholder, err)
		}
	}

	fmt.Println("\n✅ ¡Todas las descripciones han sido actualizadas!")
	fmt.Println("📋 Próximos pasos:")
	fmt.Println("   1. swag init -g cmd/main.go -o ./docs --parseDependency --parseInternal")
	fmt.Println("   2. go run cmd/main.go")
}