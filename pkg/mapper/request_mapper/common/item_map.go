package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	itemVO "github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/item"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
)

// MapCommonRequestItems mapea un arreglo de items comunes a un modelo de item -> Origen: Request
func MapCommonRequestItems(items []structs.ItemRequest) ([]models.Item, error) {
	result := make([]models.Item, len(items))

	for i, item := range items {
		itemResult, err := MapCommonRequestItem(item, i)
		if err != nil {
			return nil, err
		}
		result[i] = *itemResult
	}
	return result, nil
}

// MapCommonRequestItem mapea un item común a un modelo de item -> Origen: Request
func MapCommonRequestItem(item structs.ItemRequest, index int) (*models.Item, error) {
	var err error

	code := itemVO.NewValidatedItemCode("")
	if item.Code != nil {
		code, err = itemVO.NewItemCode(*item.Code)
		if err != nil {
			return nil, err
		}
	}

	taxCode := financial.NewValidatedTaxType("")
	if item.TaxCode != nil {
		taxCode, err = financial.NewTaxType(*item.TaxCode)
		if err != nil {
			return nil, err
		}
	}

	discount, err := financial.NewDiscount(item.Discount)
	if err != nil {
		return nil, err
	}

	quantity, err := itemVO.NewQuantity(item.Quantity)
	if err != nil {
		return nil, err
	}

	unitMeasure, err := itemVO.NewUnitMeasure(item.UnitMeasure)
	if err != nil {
		return nil, err
	}

	unitPrice, err := financial.NewAmount(item.UnitPrice)
	if err != nil {
		return nil, err
	}

	itemType, err := itemVO.NewItemType(item.Type)
	if err != nil {
		return nil, err
	}

	taxes := make([]string, len(item.Taxes))
	if item.Taxes != nil {
		for i, tax := range item.Taxes {
			taxVO, err := financial.NewTaxType(tax)
			if err != nil {
				return nil, err
			}
			taxes[i] = taxVO.GetValue()
		}
	}

	return &models.Item{
		Number:      *itemVO.NewValidatedItemNumber(index + 1),
		Type:        *itemType,
		Quantity:    *quantity,
		UnitMeasure: *unitMeasure,
		UnitPrice:   *unitPrice,
		Discount:    *discount,
		Code:        code,
		Taxes:       taxes,
		TaxCode:     taxCode,
		Description: item.Description,
		RelatedDoc:  item.RelatedDoc,
	}, nil
}
