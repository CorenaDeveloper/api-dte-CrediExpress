package mappers

import (
	"testing"

	"github.com/MarlonG1/api-facturacion-sv/config"
	"github.com/MarlonG1/api-facturacion-sv/internal/i18n"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
	"github.com/MarlonG1/api-facturacion-sv/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestMapToRetentionData(t *testing.T) {
	// Inicialización estándar
	rootPath := utils.FindProjectRoot()
	err := config.InitEnvConfig(rootPath)
	if err != nil {
		t.Fatalf("Error initializing environment config: %v", err)
	}

	err = utils.TimeInit()
	if err != nil {
		t.Fatalf("Error initializing time: %v", err)
	}

	err = i18n.InitTranslations(rootPath+"/internal/i18n", "en")
	if err != nil {
		t.Fatalf("Error initializing translations: %v", err)
	}

	// Emisor por defecto para todas las pruebas
	issuer := fixtures.CreateDefaultIssuer()

	// Definir casos de prueba
	tests := []struct {
		name      string
		req       func() *structs.CreateRetentionRequest
		wantErr   bool
		errorCode string
	}{
		// ------ VALIDACIONES BÁSICAS ------
		{
			name: "Valid retention request with physical documents",
			req: func() *structs.CreateRetentionRequest {
				return fixtures.CreatePhysicalDocumentsRetentionRequest()
			},
			wantErr: false,
		},
		{
			name: "Valid retention request with electronic documents",
			req: func() *structs.CreateRetentionRequest {
				return fixtures.CreateElectronicDocumentsRetentionRequest()
			},
			wantErr: false,
		},
		{
			name: "Valid retention request with mixed documents",
			req: func() *structs.CreateRetentionRequest {
				return fixtures.CreateMixedDocumentsRetentionRequest()
			},
			wantErr: false,
		},
		{
			name: "Null retention request",
			req: func() *structs.CreateRetentionRequest {
				return nil
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Retention without items",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Items = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Retention without receiver",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Receiver = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Physical documents without summary",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Summary = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},

		// ------ VALIDACIONES DE RECEPTOR ------
		{
			name: "Retention without receiver document type",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Receiver.DocumentType = nil
				return req
			},
			wantErr:   true,
			errorCode: "InvalidField",
		},
		{
			name: "Retention without receiver document number",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Receiver.DocumentNumber = nil
				return req
			},
			wantErr:   true,
			errorCode: "InvalidField",
		},
		{
			name: "Retention without receiver name",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Receiver.Name = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Retention without receiver email",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Receiver.Email = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Retention without receiver address",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Receiver.Address = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Retention without receiver activity code",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Receiver.ActivityCode = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Retention without receiver activity description",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Receiver.ActivityDesc = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Retention with invalid receiver email",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				invalidEmail := "not-an-email"
				req.Receiver.Email = &invalidEmail
				return req
			},
			wantErr:   true,
			errorCode: "InvalidEmail",
		},
		{
			name: "Retention with invalid receiver phone",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				invalidPhone := "123" // Demasiado corto
				req.Receiver.Phone = &invalidPhone
				return req
			},
			wantErr:   true,
			errorCode: "InvalidPhone",
		},

		// ------ VALIDACIONES DE ITEMS FÍSICOS ------
		{
			name: "Physical item without taxed amount",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Items[0].TaxedAmount = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Physical item without IVA amount",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Items[0].IvaAmount = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Physical item without emission date",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Items[0].EmissionDate = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Physical item without DTE type",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Items[0].DTEType = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Physical item with invalid document number",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Items[0].DocumentNumber = "" // Vacío (inválido)
				return req
			},
			wantErr:   true,
			errorCode: "InvalidDocumentNumberItem",
		},
		{
			name: "Physical item with invalid DTE type",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				invalidDTEType := "99" // Tipo inválido
				req.Items[0].DTEType = &invalidDTEType
				return req
			},
			wantErr:   true,
			errorCode: "InvalidDTEType",
		},
		{
			name: "Physical item with invalid emission date format",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				invalidDate := "20/04/2025" // Formato incorrecto
				req.Items[0].EmissionDate = &invalidDate
				return req
			},
			wantErr:   true,
			errorCode: "InvalidDateTime",
		},
		{
			name: "Physical item with future emission date",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				futureDate := "2099-01-01" // Fecha futura
				req.Items[0].EmissionDate = &futureDate
				return req
			},
			wantErr:   true,
			errorCode: "InvalidDateTime",
		},

		// ------ VALIDACIONES DE ITEMS ELECTRÓNICOS ------
		{
			name: "Electronic item without description",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreateElectronicDocumentsRetentionRequest()
				req.Items[0].Description = "" // Vacío (inválido)
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Electronic item with invalid document number",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreateElectronicDocumentsRetentionRequest()
				req.Items[0].DocumentNumber = "123-not-valid-uuid" // Formato inválido para documento electrónico
				return req
			},
			wantErr:   true,
			errorCode: "InvalidDocumentNumberItem",
		},
		{
			name: "Electronic item with invalid document type",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreateElectronicDocumentsRetentionRequest()
				req.Items[0].DocumentType = 99 // Tipo inválido
				return req
			},
			wantErr:   true,
			errorCode: "InvalidNumberRange",
		},

		// ------ VALIDACIONES DE CÓDIGOS DE RETENCIÓN ------
		{
			name: "Item with invalid retention code",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Items[0].RetentionCode = "99" // Código inválido
				return req
			},
			wantErr:   true,
			errorCode: "InvalidRetentionCode",
		},
		{
			name: "Item with valid retention code 22 (IVA 1%)",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Items[0].RetentionCode = "22" // Código válido
				return req
			},
			wantErr: false,
		},
		{
			name: "Item with valid retention code C4 (IVA 13%)",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Items[0].RetentionCode = "C4" // Código válido
				return req
			},
			wantErr: false,
		},
		{
			name: "Item with valid retention code C9 (Otros casos)",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Items[0].RetentionCode = "C9" // Código válido
				return req
			},
			wantErr: false,
		},

		// ------ VALIDACIONES DE CAMPOS FINANCIEROS ------
		{
			name: "Summary with negative total retention amount",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Summary.TotalRetentionAmount = -100.0 // Monto negativo
				return req
			},
			wantErr:   true,
			errorCode: "InvalidAmount",
		},
		{
			name: "Summary with negative total retention IVA",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Summary.TotalRetentionIVA = -13.0 // Monto negativo
				return req
			},
			wantErr:   true,
			errorCode: "InvalidAmount",
		},
		{
			name: "Physical item with negative taxed amount",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				negativeAmount := -100.0
				req.Items[0].TaxedAmount = &negativeAmount
				return req
			},
			wantErr:   true,
			errorCode: "InvalidAmount",
		},
		{
			name: "Physical item with negative IVA amount",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				negativeAmount := -13.0
				req.Items[0].IvaAmount = &negativeAmount
				return req
			},
			wantErr:   true,
			errorCode: "InvalidAmount",
		},

		// ------ VALIDACIONES DE CAMPOS OPCIONALES ------
		{
			name: "Retention with invalid extension",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Extension = &structs.ExtensionRequest{
					DeliveryName:     "", // Requerido pero vacío
					DeliveryDocument: "123456",
					ReceiverName:     "Ana López",
					ReceiverDocument: "98765432-1",
				}
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Retention with extension including vehicle plate (not allowed)",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				plate := "P123456"
				req.Extension = &structs.ExtensionRequest{
					DeliveryName:     "Juan Pérez",
					DeliveryDocument: "12345678-9",
					ReceiverName:     "Ana López",
					ReceiverDocument: "98765432-1",
					VehiculePlate:    &plate, // No permitido en retenciones
				}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidFieldValue",
		},
		{
			name: "Retention with invalid appendix",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Appendixes = []structs.AppendixRequest{
					{
						Field: "", // Requerido pero vacío
						Label: "Etiqueta",
						Value: "Valor",
					},
				}
				return req
			},
			wantErr:   true,
			errorCode: "ErrorMapping",
		},
		{
			name: "Retention with invalid appendix label length",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Appendixes = []structs.AppendixRequest{
					{
						Field: "campo",
						Label: "ab", // Demasiado corto (mínimo 3)
						Value: "Valor",
					},
				}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidAppendixLabel",
		},

		// ------ CASOS VÁLIDOS ADICIONALES ------
		{
			name: "Retention with valid extension",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				observation := "Observación válida"
				req.Extension = &structs.ExtensionRequest{
					DeliveryName:     "Juan Pérez",
					DeliveryDocument: "12345678-9",
					ReceiverName:     "Ana López",
					ReceiverDocument: "98765432-1",
					Observation:      &observation,
				}
				return req
			},
			wantErr: false,
		},
		{
			name: "Retention with valid appendix",
			req: func() *structs.CreateRetentionRequest {
				req := fixtures.CreatePhysicalDocumentsRetentionRequest()
				req.Appendixes = []structs.AppendixRequest{
					{
						Field: "campo_valido",
						Label: "Etiqueta válida",
						Value: "Valor válido para este apéndice",
					},
				}
				return req
			},
			wantErr: false,
		},
	}

	// Ejecutar casos de prueba
	mapper := request_mapper.NewRetentionMapper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.req()

			got, err := mapper.MapToRetentionData(req, issuer)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorCode != "" {
					assertErrorCode(t, err, tt.errorCode)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.NotNil(t, got.InputDataCommon)
			assert.NotNil(t, got.InputDataCommon.Issuer)
			assert.NotNil(t, got.InputDataCommon.Identification)
			assert.NotNil(t, got.InputDataCommon.Receiver)
			assert.NotNil(t, got.RetentionItems)
			assert.Len(t, got.RetentionItems, len(req.Items))

			// Verificar campos específicos según el tipo de documento
			isAllPhysical := true
			for _, item := range req.Items {
				if item.DocumentType == 2 {
					isAllPhysical = false
					break
				}
			}

			if isAllPhysical {
				assert.NotNil(t, got.RetentionSummary)
				assert.NotZero(t, got.RetentionSummary.TotalSubjectRetention.GetValue())
				assert.NotZero(t, got.RetentionSummary.TotalIVARetention.GetValue())
			}

			// Verificar campos opcionales si están presentes
			if req.Extension != nil {
				assert.NotNil(t, got.Extension)
			}

			if req.Appendixes != nil {
				assert.NotNil(t, got.Appendixes)
				assert.Len(t, got.Appendixes, len(req.Appendixes))
			}
		})
	}
}
