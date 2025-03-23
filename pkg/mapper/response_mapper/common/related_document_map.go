package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/response_mapper/structs"
)

// MapCommonResponseRelatedDocuments mapea los documentos relacionados a una invoice electrónica a un modelo de documento relacionado -> Origen: Response
func MapCommonResponseRelatedDocuments(docs []interfaces.RelatedDocument) []structs.DTERelatedDocument {
	result := make([]structs.DTERelatedDocument, len(docs))
	for i, doc := range docs {
		result[i] = structs.DTERelatedDocument{
			TipoDocumento:   doc.GetDocumentType(),
			TipoGeneracion:  doc.GetGenerationType(),
			NumeroDocumento: doc.GetDocumentNumber(),
			FechaEmision:    doc.GetEmissionDate().Format("2006-01-02"),
		}
	}
	return result
}
