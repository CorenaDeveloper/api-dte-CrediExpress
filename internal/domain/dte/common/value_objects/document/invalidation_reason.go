package document

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type InvalidationReason struct {
	Value string
}

func NewInvalidationReason(value string) (*InvalidationReason, error) {
	ir := &InvalidationReason{Value: value}
	if ir.IsValid() {
		return ir, nil
	}
	return nil, dte_errors.NewValidationError("InvalidInvalidationReason", value)
}

func (ir *InvalidationReason) IsValid() bool {
	return len(ir.Value) >= 5 && len(ir.Value) <= 250
}

func (ir *InvalidationReason) GetValue() string {
	return ir.Value
}

func (ir *InvalidationReason) Equals(other interfaces.ValueObject[string]) bool {
	return ir.Value == other.GetValue()
}

func (ir *InvalidationReason) ToString() string {
	return ir.Value
}
