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
// @Summary Crear factura electronica
// @Description // @Description Este endpoint permite crear y emitir una Factura Electrónica.
// @Description 
// @Description ## Ejemplo de Solicitud
// @Description ```json
// @Description {
// @Description   "items": [
// @Description     {
// @Description       "type": 1,
// @Description       "description": "CODO PVC 3/4",
// @Description       "quantity": 12,
// @Description       "unit_measure": 59,
// @Description       "unit_price": 0.65,
// @Description       "discount": 0,
// @Description       "code": "COD1",
// @Description       "non_subject_sale": 0,
// @Description       "exempt_sale": 0,
// @Description       "taxed_sale": 7.8,
// @Description       "suggested_price": 0,
// @Description       "non_taxed": 0,
// @Description       "iva_item": 0.9
// @Description     },
// @Description     {
// @Description       "type": 1,
// @Description       "description": "CODO PVC 1",
// @Description       "quantity": 123,
// @Description       "unit_measure": 59,
// @Description       "unit_price": 0.75,
// @Description       "discount": 0,
// @Description       "code": "COD2",
// @Description       "non_subject_sale": 0,
// @Description       "exempt_sale": 0,
// @Description       "taxed_sale": 92.25,
// @Description       "suggested_price": 0,
// @Description       "non_taxed": 0,
// @Description       "iva_item": 10.61
// @Description     }
// @Description   ],
// @Description   "receiver": {
// @Description     "document_type": "13",
// @Description     "document_number": "00000000-0",
// @Description     "name": "CLIENTE DE PRUEBA",
// @Description     "address": {
// @Description       "department": "08",
// @Description       "municipality": "23",
// @Description       "complement": "SOYAPANGO, SAN SALVADOR"
// @Description     },
// @Description     "phone": "21212121",
// @Description     "email": "cliente@gmail.com"
// @Description   },
// @Description   "summary": {
// @Description     "total_non_subject": 0,
// @Description     "total_exempt": 0,
// @Description     "total_taxed": 100.05,
// @Description     "sub_total": 100.05,
// @Description     "non_subject_discount": 0,
// @Description     "exempt_discount": 0,
// @Description     "taxed_discount": 0,
// @Description     "discount_percentage": 0,
// @Description     "total_discount": 0,
// @Description     "sub_total_sales": 100.05,
// @Description     "total_operation": 100.05,
// @Description     "total_non_taxed": 0,
// @Description     "total_to_pay": 99.05,
// @Description     "operation_condition": 1,
// @Description     "iva_retention": 1.00,
// @Description     "total_iva": 11.51,
// @Description     "payment_types": [
// @Description       {
// @Description         "code": "01",
// @Description         "amount": 99.05
// @Description       }
// @Description     ]
// @Description   },
// @Description 
// @Description 
// @Description   "third_party_sale": null,
// @Description   "related_docs": null,
// @Description   "other_docs": null,
// @Description   "appendixes": null
// @Description }
// @Description ```
// @Description 
// @Description ## Ejemplo de Respuesta
// @Description ```json
// @Description {
// @Description     "success": true,
// @Description   "reception_stamp": "202533D8A3CB39484D...",
// @Description   "qr_link": "https://admin.factura.gob.sv/consultaPublica?ambiente=00&codGen=UUID-GENERADO&fechaEmi=FECHA-EMISION",
// @Description     "data": {
// @Description         "identificacion": {
// @Description             "version": 1,
// @Description             "ambiente": "00",
// @Description             "tipoDte": "01",
// @Description             "numeroControl": "DTE-01-00000000-000000000000001",
// @Description             "codigoGeneracion": "29378712-E876-4...",
// @Description             "tipoModelo": 1,
// @Description             "tipoOperacion": 1,
// @Description             "tipoContingencia": null,
// @Description             "motivoContin": null,
// @Description             "fecEmi": "2025-04-16",
// @Description             "horEmi": "15:00:48",
// @Description             "tipoMoneda": "USD"
// @Description         },
// @Description         "resumen": {
// @Description             "totalNoSuj": 0,
// @Description             "totalExenta": 0,
// @Description             "totalGravada": 100.05,
// @Description             "subTotalVentas": 100.05,
// @Description             "descuNoSuj": 0,
// @Description             "descuExenta": 0,
// @Description             "descuGravada": 0,
// @Description             "porcentajeDescuento": 0,
// @Description             "totalDescu": 0,
// @Description             "tributos": [],
// @Description             "subTotal": 100.05,
// @Description             "reteRenta": 0,
// @Description             "ivaRete1": 1,
// @Description             "montoTotalOperacion": 100.05,
// @Description             "totalNoGravado": 0,
// @Description             "totalPagar": 99.05,
// @Description             "totalLetras": "NOVENTA Y NUEVE 05/100",
// @Description             "totalIva": 11.51,
// @Description             "saldoFavor": 0,
// @Description             "condicionOperacion": 1,
// @Description             "pagos": [
// @Description                 {
// @Description                     "codigo": "01",
// @Description                     "montoPago": 99.05,
// @Description                     "referencia": null,
// @Description                     "plazo": null,
// @Description                     "periodo": null
// @Description                 }
// @Description             ],
// @Description             "numPagoElectronico": null
// @Description         },
// @Description         "emisor": {
// @Description             "nit": "00000000000000",
// @Description             "nrc": "0000000",
// @Description             "nombre": "EMPRESA DE PRUEBAS SA DE CV 2",
// @Description             "codActividad": "00000",
// @Description             "descActividad": "Venta al por mayor de otros productos",
// @Description             "tipoEstablecimiento": "01",
// @Description             "direccion": {
// @Description                 "departamento": "06",
// @Description                 "municipio": "20",
// @Description                 "complemento": "BOULEVARD SANTA ELENA SUR, SANTA TECLA"
// @Description             },
// @Description             "telefono": "21212828",
// @Description             "correo": "facturacion@empresa.com.sv",
// @Description             "nombreComercial": "EJEMPLO",
// @Description             "codEstableMH": null,
// @Description             "codEstable": null,
// @Description             "codPuntoVentaMH": null,
// @Description             "codPuntoVenta": null
// @Description         },
// @Description         "receptor": {
// @Description             "nombre": "CLIENTE DE PRUEBA",
// @Description             "tipoDocumento": "13",
// @Description             "numDocumento": "00000000-0",
// @Description             "nrc": null,
// @Description             "codActividad": null,
// @Description             "descActividad": null,
// @Description             "direccion": {
// @Description                 "departamento": "08",
// @Description                 "municipio": "23",
// @Description                 "complemento": "SOYAPANGO, SAN SALVADOR"
// @Description             },
// @Description             "telefono": "21212121",
// @Description             "correo": "cliente@gmail.com"
// @Description         },
// @Description         "cuerpoDocumento": [
// @Description             {
// @Description                 "numItem": 1,
// @Description                 "tipoItem": 1,
// @Description                 "numeroDocumento": null,
// @Description                 "codigo": null,
// @Description                 "codTributo": null,
// @Description                 "descripcion": "CODO PVC 3/4",
// @Description                 "cantidad": 12,
// @Description                 "uniMedida": 59,
// @Description                 "precioUni": 0.65,
// @Description                 "montoDescu": 0,
// @Description                 "ventaNoSuj": 0,
// @Description                 "ventaExenta": 0,
// @Description                 "ventaGravada": 7.8,
// @Description                 "tributos": null,
// @Description                 "psv": 0,
// @Description                 "noGravado": 0,
// @Description                 "ivaItem": 0.9
// @Description             },
// @Description             {
// @Description                 "numItem": 2,
// @Description                 "tipoItem": 1,
// @Description                 "numeroDocumento": null,
// @Description                 "codigo": null,
// @Description                 "codTributo": null,
// @Description                 "descripcion": "CODO PVC 1",
// @Description                 "cantidad": 123,
// @Description                 "uniMedida": 59,
// @Description                 "precioUni": 0.75,
// @Description                 "montoDescu": 0,
// @Description                 "ventaNoSuj": 0,
// @Description                 "ventaExenta": 0,
// @Description                 "ventaGravada": 92.25,
// @Description                 "tributos": null,
// @Description                 "psv": 0,
// @Description                 "noGravado": 0,
// @Description                 "ivaItem": 10.61
// @Description             }
// @Description         ],
// @Description         "documentoRelacionado": null,
// @Description         "otrosDocumentos": null,
// @Description         "ventaTercero": null,
// @Description         "extension": null,
// @Description         "apendice": [
// @Description             {
// @Description                 "campo": "Datos del documento",
// @Description                 "etiqueta": "Sello de recepción",
// @Description                 "valor": "202533D8A3CB39484D..."
// @Description             }
// @Description         ]
// @Description     }
// @Description }
// @Description ```
// @Description 
// @Description Para ver ejemplos completos, consulta: /jsonExamples/
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
// @Summary Crear Comprobante de Credito Fiscal
// @Description // @Description Este endpoint permite crear y emitir un Comprobante de Crédito Fiscal (CCF) electrónico.
// @Description 
// @Description ## Ejemplo de Solicitud
// @Description ```json
// @Description {
// @Description   "items": [
// @Description     {
// @Description       "type": 2,
// @Description       "description": "Mantenimiento de computadora",
// @Description       "quantity": 1,
// @Description       "unit_measure": 59,
// @Description       "unit_price": 200.00,
// @Description       "taxed_sale": 200.00,
// @Description       "exempt_sale": 0,
// @Description       "non_subject_sale": 0,
// @Description       "taxes": ["20"]
// @Description     }
// @Description   ],
// @Description   "receiver": {
// @Description     "nrc": "3625871",
// @Description     "nit": "051283596",
// @Description     "name": "Mauricio Antonio Corena Gomez",
// @Description     "commercial_name": "Contabvs",
// @Description     "activity_code": "4259",
// @Description     "activity_description": "ACTIVIDAD ECONOMICA DE EJEMPLO",
// @Description     "address": {
// @Description       "department": "06",
// @Description       "municipality": "22",
// @Description       "complement": "Dirección de Prueba 1, N° 1234"
// @Description     },
// @Description     "phone": "61032136",
// @Description     "email": "corenadeveloper@gmail.com"
// @Description   },
// @Description   "summary": {
// @Description     "operation_condition": 1,
// @Description     "total_taxed": 200.00,
// @Description     "total_exempt": 0,
// @Description     "total_non_taxed": 0,
// @Description     "total_non_subject": 0,
// @Description     "sub_total_sales": 200.00,
// @Description     "sub_total": 200.00,
// @Description     "iva_perception": 0,
// @Description     "iva_retention": 0,
// @Description     "income_retention": 0,
// @Description     "total_operation": 260.00,
// @Description     "total_to_pay": 260.00,
// @Description     "taxes": [
// @Description       {
// @Description         "code": "20",
// @Description         "description": "IVA 13%",
// @Description         "value": 26.00
// @Description       }
// @Description     ],
// @Description     "payment_types": [
// @Description       {
// @Description         "code": "02",
// @Description         "amount": 260.00
// @Description       }
// @Description     ]
// @Description   }
// @Description }
// @Description ```
// @Description 
// @Description ## Ejemplo de Respuesta
// @Description ```json
// @Description {
// @Description   "success": true,
// @Description   "reception_stamp": "202533D8A3CB39484D...",
// @Description   "qr_link": "https://admin.factura.gob.sv/consultaPublica?ambiente=00&codGen=UUID-GENERADO&fechaEmi=FECHA-EMISION",
// @Description   "data": {
// @Description     "identificacion": {
// @Description       "version": 3,
// @Description       "ambiente": "00",
// @Description       "tipoDte": "03",
// @Description       "numeroControl": "DTE-03-C0020000-000000000000001",
// @Description       "codigoGeneracion": "96B3D5DD-EE92-4...",
// @Description       "tipoModelo": 1,
// @Description       "tipoOperacion": 1,
// @Description       "tipoContingencia": null,
// @Description       "motivoContin": null,
// @Description       "fecEmi": "2025-04-12",
// @Description       "horEmi": "19:38:32",
// @Description       "tipoMoneda": "USD"
// @Description     },
// @Description     "emisor": {
// @Description       "nit": "00000000000000",
// @Description       "nrc": "0000000",
// @Description       "nombre": "EMPRESA DE PRUEBAS SA DE CV 2",
// @Description       "codActividad": "00000",
// @Description       "descActividad": "Venta al por mayor de otros productos",
// @Description       "tipoEstablecimiento": "01",
// @Description       "direccion": {
// @Description         "departamento": "06",
// @Description         "municipio": "20",
// @Description         "complemento": "BOULEVARD SANTA ELENA SUR, SANTA TECLA"
// @Description       },
// @Description       "telefono": "21212121",
// @Description       "correo": "facturacion@empresa.com.sv",
// @Description       "nombreComercial": "EJEMPLO",
// @Description       "codEstableMH": null,
// @Description       "codEstable": "C002",
// @Description       "codPuntoVentaMH": null,
// @Description       "codPuntoVenta": null
// @Description     },
// @Description     "receptor": {
// @Description       "nombre": "CLIENTE DE PRUEBA",
// @Description       "nrc": "0000",
// @Description       "nit": "00000000000000",
// @Description       "codActividad": "00000",
// @Description       "descActividad": "ACTIVIDAD ECONOMICA DE EJEMPLO",
// @Description       "direccion": {
// @Description         "departamento": "06",
// @Description         "municipio": "22",
// @Description         "complemento": "Dirección de Prueba 1, N° 1234"
// @Description       },
// @Description       "telefono": "21212828",
// @Description       "correo": "cliente@gmail.com",
// @Description       "nombreComercial": "EJEMPLO S.A de S.V"
// @Description     },
// @Description     "cuerpoDocumento": [
// @Description       {
// @Description         "numItem": 1,
// @Description         "tipoItem": 1,
// @Description         "numeroDocumento": null,
// @Description         "codigo": null,
// @Description         "codTributo": null,
// @Description         "descripcion": "Venta gravada",
// @Description         "cantidad": 1,
// @Description         "uniMedida": 59,
// @Description         "precioUni": 2000,
// @Description         "montoDescu": 0,
// @Description         "ventaNoSuj": 0,
// @Description         "ventaExenta": 0,
// @Description         "ventaGravada": 2000,
// @Description         "tributos": [
// @Description           "20"
// @Description         ],
// @Description         "psv": 0,
// @Description         "noGravado": 0
// @Description       }
// @Description     ],
// @Description     "resumen": {
// @Description       "totalNoSuj": 0,
// @Description       "totalExenta": 0,
// @Description       "totalGravada": 2000,
// @Description       "subTotalVentas": 2000,
// @Description       "descuNoSuj": 0,
// @Description       "descuExenta": 0,
// @Description       "descuGravada": 0,
// @Description       "porcentajeDescuento": 0,
// @Description       "totalDescu": 0,
// @Description       "tributos": [
// @Description         {
// @Description           "codigo": "20",
// @Description           "descripcion": "IVA 13%",
// @Description           "valor": 260
// @Description         }
// @Description       ],
// @Description       "subTotal": 2000,
// @Description       "ivaRete1": 0,
// @Description       "ivaPerci1": 0,
// @Description       "reteRenta": 0,
// @Description       "montoTotalOperacion": 2260,
// @Description       "totalNoGravado": 0,
// @Description       "totalPagar": 2260,
// @Description       "totalLetras": "DOS MIL DOSCIENTOS SESENTA 00/100",
// @Description       "saldoFavor": 0,
// @Description       "condicionOperacion": 1,
// @Description       "pagos": [
// @Description         {
// @Description           "codigo": "01",
// @Description           "montoPago": 2260,
// @Description           "referencia": null,
// @Description           "plazo": null,
// @Description           "periodo": null
// @Description         }
// @Description       ],
// @Description       "numPagoElectronico": null
// @Description     },
// @Description     "documentoRelacionado": null,
// @Description     "otrosDocumentos": null,
// @Description     "ventaTercero": null,
// @Description     "extension": null,
// @Description     "apendice": [
// @Description       {
// @Description         "campo": "Datos del documento",
// @Description         "etiqueta": "Sello de recepción",
// @Description         "valor": "202533D8A3CB39484D..."
// @Description       }
// @Description     ]
// @Description   }
// @Description }
// @Description ```
// @Description 
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
// @Summary Crear Nota de Credito
// @Description // @Description Este endpoint permite crear y emitir una Nota de Crédito electrónica.
// @Description 
// @Description ## Ejemplo de Solicitud
// @Description ```json
// @Description {
// @Description     "items": [
// @Description         {
// @Description             "type": 1,
// @Description             "description": "Venta gravada",
// @Description             "quantity": 1,
// @Description             "unit_measure": 59,
// @Description             "unit_price": 1000.00,
// @Description             "taxed_sale": 1000.00,
// @Description             "exempt_sale": 0,
// @Description             "non_subject_sale": 0,
// @Description             "taxes": [
// @Description                 "20"
// @Description             ],
// @Description             "related_doc" : "DE4BD411-DEBF-4EB8-B..."
// @Description         }
// @Description     ],
// @Description     "receiver": {
// @Description         "nrc": "0000",
// @Description         "nit": "00000000000000",
// @Description         "name": "CLIENTE DE PRUEBA",
// @Description         "commercial_name": "EJEMPLO S.A de S.V",
// @Description         "activity_code": "47190",
// @Description         "activity_description": "ACTIVIDADES JURÍDICAS Y CONTABLES",
// @Description         "address": {
// @Description             "department": "06",
// @Description             "municipality": "22",
// @Description             "complement": "Dirección de Prueba 1, N° 1234"
// @Description         },
// @Description         "phone": "21212828",
// @Description         "email": "cliente@gmail.com"
// @Description     },
// @Description     "summary": {
// @Description         "operation_condition": 1,
// @Description         "total_taxed": 1000.00,
// @Description         "iva_retention": 10,
// @Description         "sub_total_sales": 1000.00,
// @Description         "sub_total": 1000.00,
// @Description         "total_operation": 1120.00,
// @Description         "taxes": [
// @Description             {
// @Description                 "code": "20",
// @Description                 "description": "IVA 13%",
// @Description                 "value": 130.00
// @Description             }
// @Description         ]
// @Description     },
// @Description     "extension": {
// @Description         "delivery_name": "Juan Pérez",
// @Description         "delivery_document": "06141809931020",
// @Description         "receiver_name": "María López",
// @Description         "receiver_document": "06142509882011",
// @Description         "observation": "Entrega en oficina central"
// @Description     },
// @Description     "related_docs": [
// @Description         {
// @Description             "document_type": "03",
// @Description             "generation_type": 2,
// @Description             "document_number": "DE4BD411-DEBF-4..."
// @Description         }
// @Description     ],
// @Description     "appendixes": null,
// @Description     "third_party_sale": null,
// @Description     "other_docs": null
// @Description 
// @Description }
// @Description 
// @Description ```
// @Description 
// @Description ## Ejemplo de Respuesta
// @Description ```json
// @Description {
// @Description     "success": true,
// @Description     "reception_stamp": "202534D1BECF3321453...",
// @Description     "qr_link": "https://admin.factura.gob.sv/consultaPublica?ambiente=00&codGen=5367521F-DD80-4B6B-9...&fechaEmi=FECHA-DE-EMISION",
// @Description     "data": {
// @Description         "identificacion": {
// @Description             "version": 3,
// @Description             "ambiente": "00",
// @Description             "tipoDte": "05",
// @Description             "numeroControl": "DTE-05-C0020000-000000000000001",
// @Description             "codigoGeneracion": "5367521F-DD80-4B6B-9...",
// @Description             "tipoModelo": 1,
// @Description             "tipoOperacion": 1,
// @Description             "tipoContingencia": null,
// @Description             "motivoContin": null,
// @Description             "fecEmi": "2025-04-16",
// @Description             "horEmi": "17:54:19",
// @Description             "tipoMoneda": "USD"
// @Description         },
// @Description         "emisor": {
// @Description             "nit": "00000000000000",
// @Description             "nrc": "0000000",
// @Description             "nombre": "EMPRESA DE PRUEBAS SA DE CV 2",
// @Description             "codActividad": "00000",
// @Description             "descActividad": "Venta al por mayor de otros productos",
// @Description             "tipoEstablecimiento": "01",
// @Description             "direccion": {
// @Description                 "departamento": "06",
// @Description                 "municipio": "20",
// @Description                 "complemento": "BOULEVARD SANTA ELENA SUR, SANTA TECLA"
// @Description             },
// @Description             "telefono": "21212828",
// @Description             "correo": "facturacion@empresa.com.sv",
// @Description             "nombreComercial": "EJEMPLO"
// @Description         },
// @Description         "receptor": {
// @Description             "nombre": "CLIENTE DE PRUEBA",
// @Description             "nrc": "0000",
// @Description             "nit": "00000000000000",
// @Description             "codActividad": "00000",
// @Description             "descActividad": "ACTIVIDADES JURÍDICAS Y CONTABLES",
// @Description             "direccion": {
// @Description                 "departamento": "06",
// @Description                 "municipio": "22",
// @Description                 "complemento": "Dirección de Prueba 1, N° 1234"
// @Description             },
// @Description             "telefono": "21212828",
// @Description             "correo": "cliente@gmail.com",
// @Description             "nombreComercial": "EJEMPLO S.A de S.V"
// @Description         },
// @Description         "cuerpoDocumento": [
// @Description             {
// @Description                 "numItem": 1,
// @Description                 "tipoItem": 1,
// @Description                 "numeroDocumento": "DE4BD411-DEBF-4EB8-B...",
// @Description                 "codigo": null,
// @Description                 "codTributo": null,
// @Description                 "descripcion": "Venta gravada",
// @Description                 "cantidad": 1,
// @Description                 "uniMedida": 59,
// @Description                 "precioUni": 1000,
// @Description                 "montoDescu": 0,
// @Description                 "ventaNoSuj": 0,
// @Description                 "ventaExenta": 0,
// @Description                 "ventaGravada": 1000,
// @Description                 "tributos": [
// @Description                     "20"
// @Description                 ]
// @Description             }
// @Description         ],
// @Description         "resumen": {
// @Description             "totalNoSuj": 0,
// @Description             "totalExenta": 0,
// @Description             "totalGravada": 1000,
// @Description             "subTotalVentas": 1000,
// @Description             "descuNoSuj": 0,
// @Description             "descuExenta": 0,
// @Description             "descuGravada": 0,
// @Description             "totalDescu": 0,
// @Description             "tributos": [
// @Description                 {
// @Description                     "codigo": "20",
// @Description                     "descripcion": "IVA 13%",
// @Description                     "valor": 130
// @Description                 }
// @Description             ],
// @Description             "subTotal": 1000,
// @Description             "ivaRete1": 10,
// @Description             "ivaPerci1": 0,
// @Description             "reteRenta": 0,
// @Description             "montoTotalOperacion": 1120,
// @Description             "totalLetras": "UN MIL CIENTO VEINTE 00/100",
// @Description             "condicionOperacion": 1
// @Description         },
// @Description         "documentoRelacionado": [
// @Description             {
// @Description                 "tipoDocumento": "03",
// @Description                 "tipoGeneracion": 2,
// @Description                 "numeroDocumento": "DE4BD411-DEBF-4EB8-B...",
// @Description                 "fechaEmision": "2025-04-16"
// @Description             }
// @Description         ],
// @Description         "ventaTercero": null,
// @Description         "extension": {
// @Description             "nombEntrega": "Juan Pérez",
// @Description             "docuEntrega": "06141809931020",
// @Description             "nombRecibe": "María López",
// @Description             "docuRecibe": "06142509882011",
// @Description             "observaciones": "Entrega en oficina central"
// @Description         },
// @Description         "apendice": [
// @Description             {
// @Description                 "campo": "Datos del documento",
// @Description                 "etiqueta": "Sello de recepción",
// @Description                 "valor": "202534D1BECF33214..."
// @Description             }
// @Description         ]
// @Description     }
// @Description }
// @Description ```
// @Description 
// @Description Para ver ejemplos completos, consulta: /jsonExamples/
// @Tags DTE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param creditnote body object true "Datos de la nota de credito"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} response.APIError
// @Failure 401 {object} response.APIError
// @Failure 500 {object} response.APIError
// @Router /dte/creditnote [post]
func (h *GenericCreatorDTEHandler) CreateCreditNote(w http.ResponseWriter, r *http.Request) {
	h.HandleCreate(w, r)
}

// CreateRetention godoc
// @Summary Crear Comprobante de Retencion
// @Description // @Description Este endpoint permite crear y emitir un Comprobante de Retención electrónico.
// @Description 
// @Description ## Ejemplo de Solicitud
// @Description ```json
// @Description {
// @Description     "items": [
// @Description         {
// @Description             "type": 2,
// @Description             "document_number": "1EEAB582-AA75-4D9C-A...",
// @Description             "description": "Compra de equipos informáticos",
// @Description             "retention_code": "22"
// @Description         }
// @Description     ],
// @Description     "receiver": {
// @Description     "document_type": "36",
// @Description     "document_number": "00000000000000",
// @Description     "nrc": "000000",
// @Description     "name": "EJEMPLO S.A de S.V",
// @Description     "commercial_name": "EJEMPLO",
// @Description     "activity_code": "00000",
// @Description     "activity_description": "ACTIVIDADES JURÍDICAS Y CONTABLES",
// @Description     "address": {
// @Description       "department": "06",
// @Description       "municipality": "20",
// @Description       "complement": "Dirección de Prueba 1, N° 1234"
// @Description     },
// @Description     "phone": "21212121",
// @Description     "email": "cliente@gmail.com"
// @Description   },
// @Description    "extension": {
// @Description         "delivery_name": "Juan Pérez",
// @Description         "delivery_document": "06141809931020",
// @Description         "receiver_name": "María López",
// @Description         "receiver_document": "06142509882011",
// @Description         "observation": "Entrega en oficina central"
// @Description     },
// @Description 
// @Description     "appendixes": null
// @Description }
// @Description ```
// @Description 
// @Description ## Ejemplo de Respuesta
// @Description ```json
// @Description {
// @Description     "success": true,
// @Description     "reception_stamp": "202534D1BECF3321453...",
// @Description     "qr_link": "https://admin.factura.gob.sv/consultaPublica?ambiente=00&codGen=5367521F-DD80-4B6B-9...&fechaEmi=FECHA-DE-EMISION",
// @Description     "data":{
// @Description   "identificacion": {
// @Description     "version": 1,
// @Description     "ambiente": "00",
// @Description     "tipoDte": "07",
// @Description     "numeroControl": "DTE-07-C0020000-000000000000001",
// @Description     "codigoGeneracion": "EB399033-E184-41D2-B01...",
// @Description     "tipoModelo": 1,
// @Description     "tipoOperacion": 1,
// @Description     "tipoContingencia": null,
// @Description     "motivoContin": null,
// @Description     "fecEmi": "2025-04-16",
// @Description     "horEmi": "20:40:56",
// @Description     "tipoMoneda": "USD"
// @Description   },
// @Description   "resumen": {
// @Description     "totalSujetoRetencion": 80097.35,
// @Description     "totalIVAretenido": 800.97,
// @Description     "totalIVAretenidoLetras": "OCHOCIENTOS 97/100"
// @Description   },
// @Description   "emisor": {
// @Description     "nit": "00000000000000",
// @Description     "nrc": "0000000",
// @Description     "nombre": "EMPRESA DE PRUEBAS SA DE CV 2",
// @Description     "codActividad": "00000",
// @Description     "descActividad": "Venta al por mayor de otros productos",
// @Description     "tipoEstablecimiento": "01",
// @Description     "direccion": {
// @Description       "departamento": "06",
// @Description       "municipio": "20",
// @Description       "complemento": "BOULEVARD SANTA ELENA SUR, SANTA TECLA"
// @Description     },
// @Description     "telefono": "21212828",
// @Description     "correo": "facturacion@empresa.com.sv",
// @Description     "nombreComercial": "JEMEPLO",
// @Description     "codigoMH": null,
// @Description     "codigo": "C002",
// @Description     "puntoVentaMH": null,
// @Description     "puntoVenta": null
// @Description   },
// @Description   "receptor": {
// @Description     "nombre": "EJEMPLO S.A DE C.V",
// @Description     "tipoDocumento": "36",
// @Description     "numDocumento": "00000000000000",
// @Description     "nrc": "000000",
// @Description     "codActividad": "00000",
// @Description     "descActividad": "ACTIVIDADES JURÍDICAS Y CONTABLES",
// @Description     "direccion": {
// @Description       "departamento": "06",
// @Description       "municipio": "20",
// @Description       "complemento": "Dirección de Prueba 1, N° 1234"
// @Description     },
// @Description     "telefono": "21212121",
// @Description     "correo": "cliente@gmail.com",
// @Description     "nombreComercial": "EJEMPLO"
// @Description   },
// @Description   "cuerpoDocumento": [
// @Description     {
// @Description       "numItem": 1,
// @Description       "tipoDte": "03",
// @Description       "tipoDoc": 2,
// @Description       "numDocumento": "1EEAB582-AA75-4D9C-AFF...",
// @Description       "fechaEmision": "2025-04-11",
// @Description       "montoSujetoGrav": 80097.35,
// @Description       "codigoRetencionMH": "22",
// @Description       "ivaRetenido": 800.97,
// @Description       "descripcion": "Compra de equipos informáticos"
// @Description     }
// @Description   ],
// @Description   "extension": {
// @Description     "nombEntrega": "Juan Pérez",
// @Description     "docuEntrega": "06141809931020",
// @Description     "nombRecibe": "María López",
// @Description     "docuRecibe": "06142509882011",
// @Description     "observaciones": "Entrega en oficina central"
// @Description   },
// @Description         "apendice": [
// @Description             {
// @Description                 "campo": "Datos del documento",
// @Description                 "etiqueta": "Sello de recepción",
// @Description                 "valor": "202534D1BECF33214..."
// @Description             }
// @Description         ]
// @Description   }
// @Description }
// @Description ```
// @Description 
// @Description Para ver ejemplos completos, consulta: /jsonExamples/
// @Tags DTE
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer token"
// @Param retention body object true "Datos del comprobante de Retencion"
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



