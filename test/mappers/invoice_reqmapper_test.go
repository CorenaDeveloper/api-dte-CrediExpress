package mappers

import (
	"testing"

	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/test"
	"github.com/MarlonG1/api-facturacion-sv/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestMapToInvoiceData(t *testing.T) {
	test.TestMain(t)

	// Emisor por defecto para todas las pruebas
	issuer := fixtures.CreateDefaultIssuer()

	// Definir casos de prueba
	tests := []struct {
		name      string
		req       func() *structs.CreateInvoiceRequest
		wantErr   bool
		errorCode string
	}{
		// ------ VALIDACIONES BÁSICAS ------
		{
			name: "Valid invoice request",
			req: func() *structs.CreateInvoiceRequest {
				return fixtures.CreateDefaultInvoiceRequest()
			},
			wantErr: false,
		},
		{
			name: "Null invoice request",
			req: func() *structs.CreateInvoiceRequest {
				return nil
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Invoice without items",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				req.Items = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Invoice without summary",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				req.Summary = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},

		// ------ VALIDACIONES DE RECEPTOR ------
		{
			name: "Invoice with DocumentType but no DocumentNumber",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				docType := "36" // NIT
				req.Receiver.DocumentType = &docType
				req.Receiver.DocumentNumber = nil
				return req
			},
			wantErr:   true,
			errorCode: "InvalidDocumentTypeAndNumber",
		},
		{
			name: "Invoice with invalid email format",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				invalidEmail := "not-an-email"
				req.Receiver.Email = &invalidEmail
				return req
			},
			wantErr:   true,
			errorCode: "InvalidEmail",
		},
		{
			name: "Invoice with invalid phone",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				invalidPhone := "123" // Demasiado corto
				req.Receiver.Phone = &invalidPhone
				return req
			},
			wantErr:   true,
			errorCode: "InvalidPhone",
		},

		// ------ VALIDACIONES DE DIRECCIÓN ------
		{
			name: "Invoice with invalid municipality",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				req.Receiver.Address = &structs.AddressRequest{
					Department:   "06",
					Municipality: "99", // Inválido
					Complement:   "Dirección de prueba",
				}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidMunicipality",
		},
		{
			name: "Invoice with empty address fields",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				req.Receiver.Address = &structs.AddressRequest{
					Department:   "",
					Municipality: "",
					Complement:   "",
				}
				return req
			},
			wantErr:   true,
			errorCode: "ErrorMapping",
		},

		// ------ VALIDACIONES DE ITEMS ------
		{
			name: "Invoice with invalid item type",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.Type = 99 // Tipo inválido
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidItemType",
		},
		{
			name: "Invoice with negative quantity",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.Quantity = -5 // Cantidad negativa
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidQuantity",
		},
		{
			name: "Invoice with invalid unit measure",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.UnitMeasure = 0 // Inválido (debe ser 1-99)
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidNumberRange",
		},

		// ------ VALIDACIONES DE CAMPOS FINANCIEROS ------
		{
			name: "Invoice with negative amount",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.UnitPrice = -10.0 // Precio negativo
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidAmount",
		},
		{
			name: "Invoice with negative discount",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.Discount = -5.0 // Descuento negativo
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidDiscount",
		},
		{
			name: "Invoice with invalid discount (>100%)",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.Discount = 150.0 // Más de 100%
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidDiscount",
		},

		// ------ VALIDACIONES DE CAMPOS OPCIONALES ------
		{
			name: "Invoice with all valid optional fields",
			req: func() *structs.CreateInvoiceRequest {
				return fixtures.CreateInvoiceRequestWithAllOptionalFields()
			},
			wantErr: false,
		},
		{
			name: "Invoice with invalid extension",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
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
			name: "Invoice with invalid third party sale",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				req.ThirdPartySale = &structs.ThirdPartySaleRequest{
					NIT:  "", // Requerido pero vacío
					Name: "Empresa Tercero",
				}
				return req
			},
			wantErr:   true,
			errorCode: "ErrorMapping",
		},
		{
			name: "Invoice with invalid appendix",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
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
			name: "Invoice with invalid appendix label length",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
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
		{
			name: "Invoice with invalid payment type",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				req.Payments = []structs.PaymentRequest{
					{
						Code:   "77", // Código inválido
						Amount: 100.0,
					},
				}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidLength",
		},
		{
			name: "Invoice with invalid associated document code",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				description := "Documento adicional"
				detail := "Detalle del documento adicional"
				req.OtherDocs = []structs.OtherDocRequest{
					{
						DocumentCode: 99, // Inválido (debe ser 1-4)
						Description:  &description,
						Detail:       &detail,
					},
				}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidAssociatedDocumentCode",
		},
		{
			name: "Invoice with missing doctor for medical document",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				description := "Documento médico"
				detail := "Detalle médico"
				req.OtherDocs = []structs.OtherDocRequest{
					{
						DocumentCode: 3, // Documento médico
						Description:  &description,
						Detail:       &detail,
						Doctor:       nil, // Requerido pero nulo
					},
				}
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "Invoice with invalid doctor info",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				nit := "12345678901234"
				req.OtherDocs = []structs.OtherDocRequest{
					{
						DocumentCode: 3, // Documento médico
						Doctor: &structs.DoctorRequest{
							Name:        "", // Requerido pero vacío
							NIT:         &nit,
							ServiceType: 1,
						},
					},
				}
				return req
			},
			wantErr:   true,
			errorCode: "ErrorMapping",
		},
		{
			name: "Invoice with missing doctor documents",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				req.OtherDocs = []structs.OtherDocRequest{
					{
						DocumentCode: 3, // Documento médico
						Doctor: &structs.DoctorRequest{
							Name:        "Dr. Juan Pérez",
							ServiceType: 1,
							// NIT e IdentificationDoc ambos nulos
						},
					},
				}
				return req
			},
			wantErr:   true,
			errorCode: "ErrorMapping",
		},
		// ------ CASOS VÁLIDOS PARA RECEPTOR ------
		{
			name: "Invoice with valid DocumentType and DocumentNumber",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				docType := "36" // NIT
				docNumber := "06141804941035"
				req.Receiver.DocumentType = &docType
				req.Receiver.DocumentNumber = &docNumber
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with valid email format",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				validEmail := "test.valid@google.com"
				req.Receiver.Email = &validEmail
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with valid phone",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				validPhone := "22123456"
				req.Receiver.Phone = &validPhone
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with NIT",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				validNIT := "06141804941035"
				req.Receiver.NIT = &validNIT
				return req
			},
			wantErr: true,
		},

		// ------ CASOS VÁLIDOS PARA DIRECCIÓN ------
		{
			name: "Invoice with valid municipality",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				req.Receiver.Address = &structs.AddressRequest{
					Department:   "06",
					Municipality: "20",
					Complement:   "Colonia Escalón, Calle La Reforma #123",
				}
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with valid complete address",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				req.Receiver.Address = &structs.AddressRequest{
					Department:   "05",
					Municipality: "23",
					Complement:   "Residencial Santa Elena, Calle Principal #45",
				}
				return req
			},
			wantErr: false,
		},

		// ------ CASOS VÁLIDOS PARA ITEMS ------
		{
			name: "Invoice with valid item type (product)",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.Type = 1 // Producto (válido)
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with valid item type (service)",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.Type = 2 // Servicio (válido)
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with positive quantity",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.Quantity = 10.5 // Cantidad positiva
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with valid unit measure",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.UnitMeasure = 59 // Unidades (válido)
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr: false,
		},

		// ------ CASOS VÁLIDOS PARA CAMPOS FINANCIEROS ------
		{
			name: "Invoice with positive amount",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.UnitPrice = 100.50 // Precio positivo
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with zero discount",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.Discount = 0.0 // Sin descuento
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with valid discount (10%)",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				item := fixtures.CreateDefaultInvoiceItem(0)
				item.Discount = 10.0 // 10% de descuento
				req.Items = []structs.InvoiceItemRequest{item}
				return req
			},
			wantErr: false,
		},

		// ------ CASOS VÁLIDOS PARA CAMPOS OPCIONALES ------
		{
			name: "Invoice with valid extension",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				observation := "Observación válida"
				req.Extension = &structs.ExtensionRequest{
					DeliveryName:     "Juan Martínez",
					DeliveryDocument: "12345678-9",
					ReceiverName:     "María González",
					ReceiverDocument: "98765432-1",
					Observation:      &observation,
				}
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with valid third party sale",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				req.ThirdPartySale = &structs.ThirdPartySaleRequest{
					NIT:  "06141804941035",
					Name: "Tercero Válido, S.A. de C.V.",
				}
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with valid appendix",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
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
		{
			name: "Invoice with valid payment type",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				reference := "REF-123"
				req.Payments = []structs.PaymentRequest{
					{
						Code:      "01", // Código válido
						Amount:    100.0,
						Reference: &reference,
					},
				}
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with valid associated document code",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				description := "Documento adicional válido"
				detail := "Detalle válido"
				req.OtherDocs = []structs.OtherDocRequest{
					{
						DocumentCode: 1, // Código válido
						Description:  &description,
						Detail:       &detail,
					},
				}
				return req
			},
			wantErr: false,
		},
		{
			name: "Invoice with valid doctor for medical document",
			req: func() *structs.CreateInvoiceRequest {
				req := fixtures.CreateDefaultInvoiceRequest()
				nit := "06141804941035"
				req.OtherDocs = []structs.OtherDocRequest{
					{
						DocumentCode: 3, // Documento médico
						Doctor: &structs.DoctorRequest{
							Name:        "Dr. Juan Pérez",
							NIT:         &nit,
							ServiceType: 1,
						},
					},
				}
				return req
			},
			wantErr: false,
		},
	}

	// Ejecutar casos de prueba
	mapper := request_mapper.NewInvoiceMapper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.req()

			got, err := mapper.MapToInvoiceData(req, issuer)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errorCode != "" {
					test.AssertErrorCode(t, err, tt.errorCode)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.NotNil(t, got.InputDataCommon)
			assert.NotNil(t, got.InputDataCommon.Issuer)
			assert.NotNil(t, got.InputDataCommon.Identification)
			assert.NotNil(t, got.InputDataCommon.Receiver)
			assert.Len(t, got.Items, len(req.Items))
			assert.NotNil(t, got.InvoiceSummary)

			// Verificar campos opcionales si están presentes
			if req.ThirdPartySale != nil {
				assert.NotNil(t, got.ThirdPartySale)
			}

			if req.Extension != nil {
				assert.NotNil(t, got.Extension)
			}

			if req.RelatedDocs != nil {
				assert.NotNil(t, got.RelatedDocs)
				assert.Len(t, got.RelatedDocs, len(req.RelatedDocs))
			}

			if req.OtherDocs != nil {
				assert.NotNil(t, got.OtherDocs)
				assert.Len(t, got.OtherDocs, len(req.OtherDocs))
			}

			if req.Appendixes != nil {
				assert.NotNil(t, got.Appendixes)
				assert.Len(t, got.Appendixes, len(req.Appendixes))
			}
		})
	}
}
