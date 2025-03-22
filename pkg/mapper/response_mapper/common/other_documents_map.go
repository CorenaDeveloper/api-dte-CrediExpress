package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// MapCommonResponseOtherDocuments mapea los documentos asociados a una invoice electrónica a un modelo de documento asociado -> Origen: Response
func MapCommonResponseOtherDocuments(docs []interfaces.OtherDocuments) []structs.DTEOtherDocument {
	result := make([]structs.DTEOtherDocument, len(docs))
	for i, doc := range docs {
		result[i] = structs.DTEOtherDocument{
			CodDocAsociado: doc.GetAssociatedDocument(),
		}

		// Mapear campos opcionales si existen
		if doc.GetDescription() != "" {
			result[i].Description = utils.ToStringPointer(doc.GetDescription())
		}
		if doc.GetDetail() != "" {
			result[i].Detail = utils.ToStringPointer(doc.GetDetail())
		}
		if doc.GetDoctor() != nil && doc.GetDoctor().GetName() != "" && doc.GetDoctor().GetServiceType() != 0 {
			result[i].Doctor = &structs.DTEDoctor{
				Nombre:       doc.GetDoctor().GetName(),
				NIT:          utils.ToStringPointer(doc.GetDoctor().GetNIT()),
				TipoServicio: doc.GetDoctor().GetServiceType(),
			}
			if doc.GetDoctor().GetIdentification() != "" {
				identi := doc.GetDoctor().GetIdentification()
				result[i].Doctor.DocIdentificacion = &identi
			}
		}
	}
	return result
}
