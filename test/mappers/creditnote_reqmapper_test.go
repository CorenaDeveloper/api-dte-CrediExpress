package mappers

import (
	"testing"

	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/test"
	"github.com/MarlonG1/api-facturacion-sv/test/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestMapToCreditNoteData(t *testing.T) {
	test.TestMain(t)

	// Emisor por defecto para todas las pruebas
	issuer := fixtures.CreateDefaultIssuer()

	// Definir casos de prueba
	tests := []struct {
		name      string
		req       func() *structs.CreateCreditNoteRequest
		wantErr   bool
		errorCode string
	}{
		// ------ VALIDACIONES BÁSICAS ------
		{
			name: "Valid CreditNote request",
			req: func() *structs.CreateCreditNoteRequest {
				return fixtures.CreateDefaultCreditNoteRequest()
			},
			wantErr: false,
		},
		{
			name: "Null CreditNote request",
			req: func() *structs.CreateCreditNoteRequest {
				return nil
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote without items",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Items = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote without summary",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Summary = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote without receiver",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Receiver = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote without related documents",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.RelatedDocs = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote with empty related documents",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.RelatedDocs = []structs.RelatedDocRequest{}
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},

		// ------ VALIDACIONES DE RECEPTOR ------
		{
			name: "CreditNote without receiver name",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Receiver.Name = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote without receiver email",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Receiver.Email = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote without receiver address",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Receiver.Address = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote without receiver NIT",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Receiver.NIT = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote without receiver NRC",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Receiver.NRC = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote without activity code",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Receiver.ActivityCode = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote without activity description",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Receiver.ActivityDesc = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote without commercial name",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Receiver.CommercialName = nil
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote with invalid NIT format",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				invalidNIT := "123456" // Formato inválido
				req.Receiver.NIT = &invalidNIT
				return req
			},
			wantErr:   true,
			errorCode: "InvalidPattern",
		},
		{
			name: "CreditNote with invalid NRC format",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				invalidNRC := "ABC123" // Formato inválido
				req.Receiver.NRC = &invalidNRC
				return req
			},
			wantErr:   true,
			errorCode: "InvalidFormat",
		},
		{
			name: "CreditNote with invalid email format",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				invalidEmail := "not-an-email"
				req.Receiver.Email = &invalidEmail
				return req
			},
			wantErr:   true,
			errorCode: "InvalidEmail",
		},
		{
			name: "CreditNote with invalid phone",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				invalidPhone := "123" // Demasiado corto
				req.Receiver.Phone = &invalidPhone
				return req
			},
			wantErr:   true,
			errorCode: "InvalidPhone",
		},

		// ------ VALIDACIONES DE DIRECCIÓN ------
		{
			name: "CreditNote with invalid municipality",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
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
			name: "CreditNote with empty address fields",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
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
			name: "CreditNote with invalid item type",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				item := fixtures.CreateDefaultCreditNoteItem(0)
				item.Type = 99 // Tipo inválido
				req.Items = []structs.CreditNoteItemRequest{item}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidItemType",
		},
		{
			name: "CreditNote with negative quantity",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				item := fixtures.CreateDefaultCreditNoteItem(0)
				item.Quantity = -5 // Cantidad negativa
				req.Items = []structs.CreditNoteItemRequest{item}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidQuantity",
		},
		{
			name: "CreditNote with invalid unit measure",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				item := fixtures.CreateDefaultCreditNoteItem(0)
				item.UnitMeasure = 0 // Inválido (debe ser 1-99)
				req.Items = []structs.CreditNoteItemRequest{item}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidNumberRange",
		},

		// ------ VALIDACIONES DE DOCUMENTOS RELACIONADOS ------
		{
			name: "CreditNote with invalid related document type",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.RelatedDocs = []structs.RelatedDocRequest{
					{
						DocumentType:   "99", // Tipo inválido
						GenerationType: 1,
						DocumentNumber: "12345678",
						EmissionDate:   "2023-01-15",
					},
				}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidDTEType",
		},
		{
			name: "CreditNote with invalid related document generation type",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.RelatedDocs = []structs.RelatedDocRequest{
					{
						DocumentType:   "03", // CCF (válido)
						GenerationType: 99,   // Tipo inválido
						DocumentNumber: "12345678",
						EmissionDate:   "2023-01-15",
					},
				}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidLength",
		},
		{
			name: "CreditNote with empty document number in related doc",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.RelatedDocs = []structs.RelatedDocRequest{
					{
						DocumentType:   "03", // CCF (válido)
						GenerationType: 1,
						DocumentNumber: "", // Vacío (inválido)
						EmissionDate:   "2023-01-15",
					},
				}
				return req
			},
			wantErr:   true,
			errorCode: "RequiredField",
		},
		{
			name: "CreditNote with invalid emission date format in related doc",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.RelatedDocs = []structs.RelatedDocRequest{
					{
						DocumentType:   "03", // CCF (válido)
						GenerationType: 2,
						DocumentNumber: "0408DCE7-8E96-47AA-92B2-B0F0C8FBDAF3",
						EmissionDate:   "15/01/2023", // Formato incorrecto
					},
				}
				return req
			},
			wantErr:   true,
			errorCode: "ErrorMapping",
		},
		{
			name: "CreditNote with future emission date in related doc",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.RelatedDocs = []structs.RelatedDocRequest{
					{
						DocumentType:   "03", // CCF (válido)
						GenerationType: 1,
						DocumentNumber: "12345678",
						EmissionDate:   "2099-01-15", // Fecha futura
					},
				}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidDateTime",
		},

		// ------ VALIDACIONES DE CAMPOS FINANCIEROS ------
		{
			name: "CreditNote with negative amount",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				item := fixtures.CreateDefaultCreditNoteItem(0)
				item.UnitPrice = -10.0 // Precio negativo
				req.Items = []structs.CreditNoteItemRequest{item}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidAmount",
		},
		{
			name: "CreditNote with negative discount",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				item := fixtures.CreateDefaultCreditNoteItem(0)
				item.Discount = -5.0 // Descuento negativo
				req.Items = []structs.CreditNoteItemRequest{item}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidDiscount",
		},
		{
			name: "CreditNote with invalid discount (>100%)",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				item := fixtures.CreateDefaultCreditNoteItem(0)
				item.Discount = 150.0 // Más de 100%
				req.Items = []structs.CreditNoteItemRequest{item}
				return req
			},
			wantErr:   true,
			errorCode: "InvalidDiscount",
		},
		{
			name: "CreditNote with invalid payment condition",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Summary.OperationCondition = 99 // Condición inválida
				return req
			},
			wantErr:   true,
			errorCode: "InvalidNumberRange",
		},

		// ------ VALIDACIONES DE CAMPOS OPCIONALES ------
		{
			name: "CreditNote with all valid optional fields",
			req: func() *structs.CreateCreditNoteRequest {
				return fixtures.CreateCreditNoteRequestWithAllOptionalFields()
			},
			wantErr: false,
		},
		{
			name: "CreditNote with invalid extension",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
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
			name: "CreditNote with invalid third party sale",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
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
			name: "CreditNote with invalid appendix",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
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
			name: "CreditNote with invalid appendix label length",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
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
			name: "CreditNote with invalid payment type",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
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
			name: "CreditNote with invalid associated document code",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
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
			name: "CreditNote with missing doctor for medical document",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
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

		// ------ CASOS VÁLIDOS PARA RECEPTOR ------
		{
			name: "CreditNote with valid NIT",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				validNIT := "06141804941035"
				req.Receiver.NIT = &validNIT
				return req
			},
			wantErr: false,
		},
		{
			name: "CreditNote with valid NRC",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				validNRC := "1234567"
				req.Receiver.NRC = &validNRC
				return req
			},
			wantErr: false,
		},
		{
			name: "CreditNote with valid email",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				validEmail := "cliente@gmail.com"
				req.Receiver.Email = &validEmail
				return req
			},
			wantErr: false,
		},
		{
			name: "CreditNote with valid phone",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				validPhone := "22345678"
				req.Receiver.Phone = &validPhone
				return req
			},
			wantErr: false,
		},

		// ------ CASOS VÁLIDOS PARA DIRECCIÓN ------
		{
			name: "CreditNote with valid address",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Receiver.Address = &structs.AddressRequest{
					Department:   "06",
					Municipality: "20",
					Complement:   "Colonia Escalón, Calle La Reforma #123",
				}
				return req
			},
			wantErr: false,
		},

		// ------ CASOS VÁLIDOS PARA ITEMS ------
		{
			name: "CreditNote with valid item type (product)",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				item := fixtures.CreateDefaultCreditNoteItem(0)
				item.Type = 1 // Producto (válido)
				req.Items = []structs.CreditNoteItemRequest{item}
				return req
			},
			wantErr: false,
		},
		{
			name: "CreditNote with valid item type (service)",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				item := fixtures.CreateDefaultCreditNoteItem(0)
				item.Type = 2 // Servicio (válido)
				req.Items = []structs.CreditNoteItemRequest{item}
				return req
			},
			wantErr: false,
		},
		{
			name: "CreditNote with valid unit measure",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				item := fixtures.CreateDefaultCreditNoteItem(0)
				item.UnitMeasure = 59 // Unidades (válido)
				req.Items = []structs.CreditNoteItemRequest{item}
				return req
			},
			wantErr: false,
		},

		// ------ CASOS VÁLIDOS PARA DOCUMENTOS RELACIONADOS ------
		{
			name: "CreditNote with valid related document (CCF)",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.RelatedDocs = []structs.RelatedDocRequest{
					{
						DocumentType:   "03", // CCF (válido)
						GenerationType: 1,
						DocumentNumber: "12345678",
						EmissionDate:   "2023-01-15",
					},
				}
				return req
			},
			wantErr: false,
		},
		{
			name: "CreditNote with valid related document (Withholding Receipt)",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.RelatedDocs = []structs.RelatedDocRequest{
					{
						DocumentType:   "07", // Comprobante de Retención (válido)
						GenerationType: 1,
						DocumentNumber: "12345678",
						EmissionDate:   "2023-01-15",
					},
				}
				return req
			},
			wantErr: false,
		},
		{
			name: "CreditNote with multiple valid related documents",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.RelatedDocs = []structs.RelatedDocRequest{
					{
						DocumentType:   "03", // CCF (válido)
						GenerationType: 1,
						DocumentNumber: "12345678",
						EmissionDate:   "2023-01-15",
					},
					{
						DocumentType:   "07", // Comprobante de Retención (válido)
						GenerationType: 1,
						DocumentNumber: "87654321",
						EmissionDate:   "2023-01-16",
					},
				}
				return req
			},
			wantErr: false,
		},

		// ------ CASOS VÁLIDOS PARA CAMPOS FINANCIEROS ------
		{
			name: "CreditNote with valid amount",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				item := fixtures.CreateDefaultCreditNoteItem(0)
				item.UnitPrice = 100.50 // Precio válido
				req.Items = []structs.CreditNoteItemRequest{item}
				return req
			},
			wantErr: false,
		},
		{
			name: "CreditNote with valid discount",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				item := fixtures.CreateDefaultCreditNoteItem(0)
				item.Discount = 10.5 // Descuento válido
				req.Items = []structs.CreditNoteItemRequest{item}
				return req
			},
			wantErr: false,
		},
		{
			name: "CreditNote with valid payment condition",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.Summary.OperationCondition = 1 // Contado (válido)
				return req
			},
			wantErr: false,
		},

		// ------ CASOS VÁLIDOS PARA CAMPOS OPCIONALES ------
		{
			name: "CreditNote with valid extension",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				observation := "Observación válida"
				vehiculePlate := "P123456"
				req.Extension = &structs.ExtensionRequest{
					DeliveryName:     "Juan Martínez",
					DeliveryDocument: "12345678-9",
					ReceiverName:     "María González",
					ReceiverDocument: "98765432-1",
					Observation:      &observation,
					VehiculePlate:    &vehiculePlate,
				}
				return req
			},
			wantErr: false,
		},
		{
			name: "CreditNote with valid third party sale",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				req.ThirdPartySale = &structs.ThirdPartySaleRequest{
					NIT:  "06141804941035",
					Name: "Tercero Válido, S.A. de C.V.",
				}
				return req
			},
			wantErr: false,
		},
		{
			name: "CreditNote with valid appendix",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
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
			name: "CreditNote with valid payment type",
			req: func() *structs.CreateCreditNoteRequest {
				req := fixtures.CreateDefaultCreditNoteRequest()
				reference := "REF-123"
				req.Payments = []structs.PaymentRequest{
					{
						Code:      "01", // Efectivo (válido)
						Amount:    100.0,
						Reference: &reference,
					},
				}
				return req
			},
			wantErr: false,
		},
	}

	// Ejecutar casos de prueba
	mapper := request_mapper.NewCreditNoteMapper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.req()

			got, err := mapper.MapToCreditNoteData(req, issuer)

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
			assert.Len(t, got.Items, len(req.Items))
			assert.NotNil(t, got.CreditSummary)
			assert.NotNil(t, got.RelatedDocs)
			assert.Len(t, got.RelatedDocs, len(req.RelatedDocs))

			// Verificar campos opcionales si están presentes
			if req.ThirdPartySale != nil {
				assert.NotNil(t, got.ThirdPartySale)
			}

			if req.Extension != nil {
				assert.NotNil(t, got.Extension)
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
