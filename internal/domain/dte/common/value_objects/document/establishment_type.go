package document

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type EstablishmentType struct {
	Value string `json:"value"`
}

func NewEstablishmentType(value string) (*EstablishmentType, error) {
	est := &EstablishmentType{Value: value}
	if est.IsValid() {
		return est, nil
	}
	return &EstablishmentType{}, dte_errors.NewValidationError("InvalidEstablishmentType", value)
}

func NewValidatedEstablishmentType(value string) *EstablishmentType {
	return &EstablishmentType{Value: value}
}

// IsValid valida que el valor de EstablishmentType sea 01, 02, 04, 07 o 20
func (et *EstablishmentType) IsValid() bool {
	for _, v := range constants.AllowedEstablishmentTypes {
		if et.Value == v {
			return true
		}
	}
	return false
}

func (et *EstablishmentType) Equals(other interfaces.ValueObject[string]) bool {
	return et.GetValue() == other.GetValue()
}

func (et *EstablishmentType) GetValue() string {
	return et.Value
}

func (et *EstablishmentType) ToString() string {
	return et.Value
}
