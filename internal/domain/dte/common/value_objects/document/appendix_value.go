package document

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type AppendixValue struct {
	Value string `json:"value"`
}

func NewAppendixValue(value string) (*AppendixValue, error) {
	v := &AppendixValue{Value: value}
	if v.IsValid() {
		return v, nil
	}
	return &AppendixValue{}, dte_errors.NewValidationError("InvalidLength", "AppendixValue", "1 a 150", value)
}

func NewValidatedAppendixValue(value string) *AppendixValue {
	return &AppendixValue{Value: value}
}

func (v *AppendixValue) IsValid() bool {
	return len(v.Value) >= 1 && len(v.Value) <= 150
}

func (v *AppendixValue) Equals(other interfaces.ValueObject[string]) bool {
	return v.Value == other.GetValue()
}

func (v *AppendixValue) GetValue() string {
	return v.Value
}

func (v *AppendixValue) ToString() string {
	return v.Value
}
