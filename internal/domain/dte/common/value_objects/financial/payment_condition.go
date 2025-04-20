package financial

import (
	"fmt"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
)

type PaymentCondition struct {
	Value int
}

func NewPaymentCondition(value int) (*PaymentCondition, error) {
	paymentCondition := &PaymentCondition{Value: value}
	if paymentCondition.IsValid() {
		return paymentCondition, nil
	}
	return &PaymentCondition{}, dte_errors.NewValidationError("InvalidNumberRange", "payment.condition", "1-3", fmt.Sprintf("%d", value))
}

func NewValidatedPaymentCondition(value int) *PaymentCondition {
	return &PaymentCondition{Value: value}
}

func (p *PaymentCondition) GetValue() int {
	return p.Value
}

func (p *PaymentCondition) IsValid() bool {
	return p.Value >= 1 && p.Value <= 3
}

func (p *PaymentCondition) Equals(other interfaces.ValueObject[int]) bool {
	return p.Value == other.GetValue()
}

func (p *PaymentCondition) ToString() string {
	return fmt.Sprintf("%d", p.Value)
}
