
# Función para procesar JSON y generar descripción
function Process-Json {
    param(
        [string]$JsonFile,
        [string]$Title
    )
    
    if (Test-Path $JsonFile) {
        $content = Get-Content $JsonFile -Raw
        return @"
## $Title
``````json
$content
``````

"@
    } else {
        return "// Error: No se encontró el archivo $JsonFile`n"
    }
}

# Función para generar descripción completa
function Generate-Description {
    param(
        [string]$RequestFile,
        [string]$ResponseFile,
        [string]$Description
    )
    
    $result = "$Description`n`n"
    $result += Process-Json $RequestFile "Ejemplo de Solicitud"
    $result += Process-Json $ResponseFile "Ejemplo de Respuesta"
    $result += "Para ver ejemplos completos, consulta: /jsonExamples/"
    
    return $result
}

# Función para reemplazar en archivo
function Replace-InFile {
    param(
        [string]$FilePath,
        [string]$Placeholder,
        [string]$NewContent
    )
    
    if (Test-Path $FilePath) {
        # Leer contenido del archivo
        $content = Get-Content $FilePath -Raw
        
        # Preparar el contenido para Swagger (agregar // @Description a cada línea)
        $lines = $NewContent -split "`n"
        $swaggerContent = ""
        foreach ($line in $lines) {
            if ($line.Trim() -ne "") {
                $swaggerContent += "// @Description $line`n"
            } else {
                $swaggerContent += "// @Description `n"
            }
        }
        
        # Reemplazar placeholder
        $content = $content -replace "\{\{$Placeholder\}\}", $swaggerContent.TrimEnd()
        
        # Escribir archivo actualizado
        Set-Content -Path $FilePath -Value $content -Encoding UTF8
        
        Write-Host "✅ Reemplazado $Placeholder en $FilePath" -ForegroundColor Green
    } else {
        Write-Host "❌ No se encontró el archivo $FilePath" -ForegroundColor Red
    }
}

# Archivo de destino
$HandlerFile = "internal\infrastructure\api\handlers\generic_creator_handler.go"

Write-Host "🔄 Iniciando actualización de descripciones Swagger..." -ForegroundColor Yellow

# Verificar que existe el archivo de destino
if (-not (Test-Path $HandlerFile)) {
    Write-Host "❌ Error: No se encontró $HandlerFile" -ForegroundColor Red
    exit 1
}

# Generar descripciones para cada endpoint
Write-Host "📝 Generando descripción para CCF..." -ForegroundColor Cyan
$ccfDesc = Generate-Description "jsonExamples\ccf_request.json" "jsonExamples\ccf_response.json" "Este endpoint permite crear y emitir un Comprobante de Crédito Fiscal (CCF) electrónico."

Write-Host "📝 Generando descripción para Invoice..." -ForegroundColor Cyan
$invoiceDesc = Generate-Description "jsonExamples\invoice_request.json" "jsonExamples\invoice_response.json" "Este endpoint permite crear y emitir una Factura Electrónica."

Write-Host "📝 Generando descripción para Credit Note..." -ForegroundColor Cyan
$creditnoteDesc = Generate-Description "jsonExamples\creditnote_request.json" "jsonExamples\creditnote_response.json" "Este endpoint permite crear y emitir una Nota de Crédito electrónica."

Write-Host "📝 Generando descripción para Retention..." -ForegroundColor Cyan
$retentionDesc = Generate-Description "jsonExamples\retention_request.json" "jsonExamples\retention_response.json" "Este endpoint permite crear y emitir un Comprobante de Retención electrónico."

# Reemplazar todas las descripciones
Replace-InFile $HandlerFile "CCF_DESCRIPTION" $ccfDesc
Replace-InFile $HandlerFile "INVOICE_DESCRIPTION" $invoiceDesc
Replace-InFile $HandlerFile "CREDITNOTE_DESCRIPTION" $creditnoteDesc
Replace-InFile $HandlerFile "RETENTION_DESCRIPTION" $retentionDesc

Write-Host "`n✅ ¡Todas las descripciones han sido actualizadas automáticamente!" -ForegroundColor Green
Write-Host "📋 Próximos pasos:" -ForegroundColor Yellow
Write-Host "   1. swag init -g cmd/main.go -o ./docs --parseDependency --parseInternal" -ForegroundColor White
Write-Host "   2. go run cmd/main.go" -ForegroundColor White