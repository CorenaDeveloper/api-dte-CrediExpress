package strategy

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/retention/retention_models"
	"github.com/shopspring/decimal"
)

type RetentionTotalStrategy struct {
	Document *retention_models.RetentionModel
}

func (s *RetentionTotalStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil {
		return nil
	}

	validations := []func() *dte_errors.DTEError{
		s.validateTotalAmounts,
	}

	for _, validate := range validations {
		if err := validate(); err != nil {
			return err
		}
	}

	return nil
}

// validateTotalAmounts valida los totales de retención del documento
func (s *RetentionTotalStrategy) validateTotalAmounts() *dte_errors.DTEError {
	// 1. Obtener los totales de retención esperados del documento
	expectedTotalSubjectRetention := s.Document.RetentionSummary.TotalSubjectRetention.GetValueAsDecimal()
	actualIvaRetention := s.Document.RetentionSummary.TotalIVARetention.GetValueAsDecimal()

	// 2. Validar que el total de retenciones concuerde con el total calculado por los items
	actualTotalSubjectRetention, expectedIvaRetention := s.Document.GetTotalByItems()
	if !s.compareTotalsWithTolerance(expectedTotalSubjectRetention, actualTotalSubjectRetention, 0.01) {
		return dte_errors.NewDTEErrorSimple("InvalidTotalSubjectRetention",
			actualTotalSubjectRetention.InexactFloat64(),
			expectedTotalSubjectRetention.InexactFloat64())
	}

	// 3. Validar que el total de IVA retenciones concuerde con el total calculado por los items
	if !s.compareTotalsWithTolerance(expectedIvaRetention, actualIvaRetention, 0.01) {
		return dte_errors.NewDTEErrorSimple("InvalidTotalIVARetention",
			actualIvaRetention.InexactFloat64(),
			expectedIvaRetention.InexactFloat64())
	}

	return nil
}

// compareTotalsWithTolerance compara dos totales con una tolerancia especificada
func (s *RetentionTotalStrategy) compareTotalsWithTolerance(expected, actual decimal.Decimal, tolerance float64) bool {
	diff := expected.Sub(actual).Abs()
	return diff.LessThanOrEqual(decimal.NewFromFloat(tolerance))
}
