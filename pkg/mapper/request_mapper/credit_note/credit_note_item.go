package credit_note

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/credit_note_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

func MapCreditNoteItems(item []structs.CreditNoteItemRequest) ([]credit_note_models.CreditNoteItem, error) {
	result := make([]credit_note_models.CreditNoteItem, len(item))

	for i, noteItem := range item {
		itemMapped, err := MapCreditNoteRequestItem(noteItem, i)
		if err != nil {
			return nil, err
		}
		result[i] = *itemMapped
	}

	return result, nil
}

// MapCreditNoteRequestItem mapea un item de Nota de Crédito -> Origen: Request
func MapCreditNoteRequestItem(item structs.CreditNoteItemRequest, index int) (*credit_note_models.CreditNoteItem, error) {
	baseItem, err := common.MapCommonRequestItem(structs.ItemRequest{
		Type:        item.Type,
		Quantity:    item.Quantity,
		UnitMeasure: item.UnitMeasure,
		UnitPrice:   item.UnitPrice,
		Discount:    item.Discount,
		Code:        item.Code,
		Taxes:       item.Taxes,
		TaxCode:     item.TaxCode,
		Description: item.Description,
		RelatedDoc:  item.RelatedDoc,
	}, index)

	if err != nil {
		return nil, err
	}

	nonSubjectSale, err := financial.NewAmount(item.NonSubjectSale)
	if err != nil {
		return nil, err
	}

	exemptSale, err := financial.NewAmount(item.ExemptSale)
	if err != nil {
		return nil, err
	}

	taxedSale, err := financial.NewAmount(item.TaxedSale)
	if err != nil {
		return nil, err
	}

	suggestedPrice, err := financial.NewAmount(item.SuggestedPrice)
	if err != nil {
		return nil, err
	}

	nonTaxed, err := financial.NewAmount(item.NonTaxed)
	if err != nil {
		return nil, err
	}

	return &credit_note_models.CreditNoteItem{
		Item:           baseItem,
		NonSubjectSale: *nonSubjectSale,
		ExemptSale:     *exemptSale,
		TaxedSale:      *taxedSale,
		SuggestedPrice: *suggestedPrice,
		NonTaxed:       *nonTaxed,
	}, nil
}
