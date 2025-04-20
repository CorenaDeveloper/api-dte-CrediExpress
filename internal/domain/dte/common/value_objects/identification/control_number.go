package identification

import (
	"regexp"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type ControlNumber struct {
	Value string `json:"value"`
}

func NewControlNumber(value string) (*ControlNumber, error) {
	controlNumber := &ControlNumber{Value: value}
	if controlNumber.IsValid() {
		return controlNumber, nil
	}
	return &ControlNumber{}, dte_errors.NewValidationError("InvalidPattern", "control_number", "DTE-00-XXXXXXXX-000000000000000", value)
}

func NewValidatedControlNumber(value string) *ControlNumber {
	return &ControlNumber{Value: value}
}

// IsValid válido si cumple con el patrón DTE-00-XXXXXXXX-000000000000000 donde X es un caracter alfanumérico y 0 es un dígito
func (cn *ControlNumber) IsValid() bool {
	pattern := `^DTE-[0-9]{2}-[A-Z0-9]{8}-[0-9]{15}$`
	matched, _ := regexp.MatchString(pattern, cn.Value)
	return matched && len(cn.Value) == 31
}

func (cn *ControlNumber) Equals(other interfaces.ValueObject[string]) bool {
	return cn.GetValue() == other.GetValue()
}

func (cn *ControlNumber) GetValue() string {
	return cn.Value
}

func (cn *ControlNumber) ToString() string {
	return cn.Value
}
