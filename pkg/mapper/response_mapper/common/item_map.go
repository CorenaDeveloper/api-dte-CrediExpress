package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

func MapCommonItems(item interfaces.Item) structs.DTEItem {
	result := structs.DTEItem{
		NumItem:         item.GetNumber(),
		TipoItem:        item.GetType(),
		Descripcion:     item.GetDescription(),
		Cantidad:        item.GetQuantity(),
		UniMedida:       item.GetUnitMeasure(),
		PrecioUni:       item.GetUnitPrice(),
		MontoDescu:      item.GetDiscount(),
		NumeroDocumento: item.GetRelatedDoc(),
	}

	// Mapear tributos si existen
	if item.GetTaxes() != nil {
		MapTaxCodes(item.GetTaxes())
		result.Tributos = item.GetTaxes()
	} else {
		result.Tributos = nil
	}

	return result
}
