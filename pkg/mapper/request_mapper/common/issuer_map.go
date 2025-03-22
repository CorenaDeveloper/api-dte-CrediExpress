package common

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/core/dte"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/models"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/base"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/document"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/identification"
)

// MapCommonIssuer mapea un emisor común a un modelo de emisor -> Origen: Request
func MapCommonIssuer(client *dte.IssuerDTE) (*models.Issuer, error) {
	address, err := MapClientAddress(client.Address)
	if err != nil {
		return nil, err
	}

	return &models.Issuer{
		NIT:                 *identification.NewValidatedNIT(client.NIT),
		NRC:                 *identification.NewValidatedNRC(client.NRC),
		Name:                client.BusinessName,
		ActivityCode:        *identification.NewValidatedActivityCode(client.EconomicActivity),
		ActivityDescription: client.EconomicActivityDesc,
		EstablishmentType:   *document.NewValidatedEstablishmentType(client.EstablishmentType),
		Address:             address,
		Phone:               *base.NewValidatedPhone(*client.Phone),
		Email:               *base.NewValidatedEmail(*client.Email),
		CommercialName:      client.CommercialName,
		EstablishmentCode:   client.EstablishmentCode,
		EstablishmentMHCode: client.EstablishmentCodeMH,
		POSCode:             client.POSCode,
		POSMHCode:           client.POSCodeMH,
	}, nil
}
