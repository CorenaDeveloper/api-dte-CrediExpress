package ccf

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/ccf/ccf_models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

func MapCCFItems(item []structs.CreditItemRequest) ([]ccf_models.CreditItem, error) {
	result := make([]ccf_models.CreditItem, len(item))

	for i, invoiceItem := range item {
		itemMapped, err := MapCCFRequestItem(invoiceItem, i)
		if err != nil {
			return nil, err
		}
		result[i] = *itemMapped
	}

	return result, nil
}

// MapCCFRequestItem mapea un item de Comprobante de Crédito Fiscal -> Origen: Request
func MapCCFRequestItem(item structs.CreditItemRequest, index int) (*ccf_models.CreditItem, error) {

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
	}, index)

	if err != nil {
		return nil, err
	}

	baseItem.RelatedDoc = item.RelatedDoc

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

	return &ccf_models.CreditItem{
		Item:           baseItem,
		NonSubjectSale: *nonSubjectSale,
		ExemptSale:     *exemptSale,
		TaxedSale:      *taxedSale,
		SuggestedPrice: *suggestedPrice,
		NonTaxed:       *nonTaxed,
	}, nil
}
