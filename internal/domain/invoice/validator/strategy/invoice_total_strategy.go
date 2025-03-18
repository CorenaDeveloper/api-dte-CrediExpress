package strategy

import (
	"github.com/shopspring/decimal"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/invoice/invoice_models"
)

type InvoiceTotalsStrategy struct {
	Document *invoice_models.ElectronicInvoice
}

func (s *InvoiceTotalsStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil {
		return nil
	}

	validations := []func() *dte_errors.DTEError{
		s.validateSubTotal,
		s.validateDiscounts,
		s.validateTotalOperation,
	}

	for _, validate := range validations {
		if err := validate(); err != nil {
			return err
		}
	}

	return nil
}

// validateSubTotal valida el subtotal de la invoice electrónica
func (s *InvoiceTotalsStrategy) validateSubTotal() *dte_errors.DTEError {
	expectedSubTotal := decimal.NewFromFloat(s.Document.InvoiceSummary.TotalNonSubject.GetValue()).
		Add(decimal.NewFromFloat(s.Document.InvoiceSummary.TotalExempt.GetValue())).
		Add(decimal.NewFromFloat(s.Document.InvoiceSummary.TotalTaxed.GetValue())).
		Sub(decimal.NewFromFloat(s.Document.InvoiceSummary.ExemptDiscount.GetValue())).
		Sub(decimal.NewFromFloat(s.Document.InvoiceSummary.NonSubjectDiscount.GetValue())).
		Sub(decimal.NewFromFloat(s.Document.InvoiceSummary.TaxedDiscount.GetValue()))

	actualSubTotal := decimal.NewFromFloat(s.Document.InvoiceSummary.SubTotal.GetValue())

	if !s.compareTotalsWithTolerance(expectedSubTotal, actualSubTotal, 0.0001) {
		return dte_errors.NewDTEErrorSimple("InvalidSubTotal",
			actualSubTotal.InexactFloat64(),
			expectedSubTotal.InexactFloat64())
	}
	return nil
}

// validateDiscounts valida los descuentos de la invoice electrónica
func (s *InvoiceTotalsStrategy) validateDiscounts() *dte_errors.DTEError {
	totalDiscount := decimal.NewFromFloat(s.Document.InvoiceSummary.TotalDiscount.GetValue())
	subTotal := decimal.NewFromFloat(s.Document.InvoiceSummary.SubTotal.GetValue())

	if totalDiscount.LessThan(decimal.Zero) {
		return dte_errors.NewDTEErrorSimple("NegativeDiscount", totalDiscount)
	}

	if totalDiscount.GreaterThan(subTotal) {
		return dte_errors.NewDTEErrorSimple("ExcessiveDiscount", totalDiscount, subTotal)
	}

	return nil
}

// validateTotalOperation valida el total de la operación de la invoice electrónica
func (s *InvoiceTotalsStrategy) validateTotalOperation() *dte_errors.DTEError {
	expectedTotal := decimal.NewFromFloat(s.Document.InvoiceSummary.SubTotalSales.GetValue())
	expectedTotal = expectedTotal.Sub(decimal.NewFromFloat(s.Document.InvoiceSummary.ExemptDiscount.GetValue()))
	expectedTotal = expectedTotal.Sub(decimal.NewFromFloat(s.Document.InvoiceSummary.NonSubjectDiscount.GetValue()))
	expectedTotal = expectedTotal.Sub(decimal.NewFromFloat(s.Document.InvoiceSummary.TaxedDiscount.GetValue()))

	actualTotal := decimal.NewFromFloat(s.Document.InvoiceSummary.TotalOperation.GetValue())

	if !s.compareTotalsWithTolerance(expectedTotal, actualTotal, 0.0001) {
		return dte_errors.NewDTEErrorSimple("InvalidTotalOperation",
			actualTotal.InexactFloat64(),
			expectedTotal.InexactFloat64())
	}
	return nil
}

// calculateExpectedTotal calcula el total esperado de la operación
func (s *InvoiceTotalsStrategy) calculateExpectedTotal() decimal.Decimal {
	subTotalSales := decimal.NewFromFloat(s.Document.InvoiceSummary.SubTotalSales.GetValue())

	nonTaxed := decimal.NewFromFloat(s.Document.InvoiceSummary.TotalNonTaxed.GetValue())

	var totalTributos decimal.Decimal
	for _, tax := range s.Document.InvoiceSummary.TotalTaxes {
		taxValue := decimal.NewFromFloat(tax.GetValue())
		totalTributos = totalTributos.Add(taxValue)
	}

	expectedTotal := subTotalSales.Add(nonTaxed).Add(totalTributos)

	return expectedTotal
}

// compareTotalsWithTolerance compara dos totales con una tolerancia
func (s *InvoiceTotalsStrategy) compareTotalsWithTolerance(expected, actual decimal.Decimal, tolerance float64) bool {
	diff := expected.Sub(actual).Abs()
	return diff.LessThanOrEqual(decimal.NewFromFloat(tolerance))
}
