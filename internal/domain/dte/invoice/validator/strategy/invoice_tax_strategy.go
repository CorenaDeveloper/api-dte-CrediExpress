package strategy

import (
	"fmt"
	"github.com/shopspring/decimal"

	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/constants"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/dte_errors"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/interfaces"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/validator/strategy"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/invoice_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/logs"
)

// InvoiceTaxStrategy implementa la validación de impuestos de una invoice electrónica
type InvoiceTaxStrategy struct {
	*strategy.TaxCalculationStrategy
	Document *invoice_models.ElectronicInvoice
}

// Validate valida los impuestos de la invoice, sobreescribiendo el método de la interfaz
func (s *InvoiceTaxStrategy) Validate() *dte_errors.DTEError {
	if s.Document == nil || len(s.Document.InvoiceItems) == 0 {
		return nil
	}

	// 1. Validar totales base
	if err := s.validateBaseTotals(); err != nil {
		return err
	}

	// 2. Validar IVA
	if err := s.validateIVA(); err != nil {
		return err
	}

	if err := s.validatePerception(); err != nil {
		return err
	}

	// 3. Validar montos monetarios
	if err := s.validateMonetaryAmounts(); err != nil {
		return err
	}

	// 4. Validar montos totales
	if err := s.validateTotalAmounts(); err != nil {
		return err
	}

	// 5. Validar estrategia de items
	for _, item := range s.Document.InvoiceItems {
		if len(s.Document.RelatedDocuments) > 0 {

			if item.GetRelatedDoc() == nil {
				logs.Error("Missing related document in item when related_docs is present", map[string]interface{}{
					"itemNumber": item.GetNumber(),
				})
				return dte_errors.NewDTEErrorSimple("MissingItemRelatedDoc", item.GetNumber())
			}

			found := false
			itemRelatedDoc := *item.GetRelatedDoc()

			for _, relatedDoc := range s.Document.RelatedDocuments {
				if relatedDoc.GetDocumentNumber() == itemRelatedDoc {
					found = true
					break
				}
			}

			if !found {
				logs.Error("Item related document not found in document related docs", map[string]interface{}{
					"itemNumber": item.GetNumber(),
					"relatedDoc": itemRelatedDoc,
				})
				return dte_errors.NewDTEErrorSimple("InvalidItemRelatedDoc",
					item.GetNumber(),
					itemRelatedDoc)
			}
		}
	}

	return nil
}

func (s *InvoiceTaxStrategy) validateTotalAmounts() *dte_errors.DTEError {
	// Obtener total operación
	totalOperation := decimal.NewFromFloat(s.Document.InvoiceSummary.TotalOperation.GetValue())

	// Obtener montos que afectan el total a pagar
	taxedAmount := decimal.NewFromFloat(s.Document.InvoiceSummary.TotalTaxed.GetValue())

	expectedSubTotal := decimal.NewFromFloat(s.Document.InvoiceSummary.SubTotalSales.GetValue()).
		Sub(decimal.NewFromFloat(s.Document.InvoiceSummary.TaxedDiscount.GetValue())).
		Sub(decimal.NewFromFloat(s.Document.InvoiceSummary.ExemptDiscount.GetValue())).
		Sub(decimal.NewFromFloat(s.Document.InvoiceSummary.NonSubjectDiscount.GetValue()))

	actualSubTotal := decimal.NewFromFloat(s.Document.InvoiceSummary.SubTotal.GetValue())
	if !expectedSubTotal.Equal(actualSubTotal) {
		logs.Error("Invalid subtotal calculation with discounts", map[string]interface{}{
			"expected":           expectedSubTotal,
			"actual":             actualSubTotal,
			"taxedDiscount":      s.Document.InvoiceSummary.TaxedDiscount.GetValue(),
			"exemptDiscount":     s.Document.InvoiceSummary.ExemptDiscount.GetValue(),
			"nonSubjectDiscount": s.Document.InvoiceSummary.NonSubjectDiscount.GetValue(),
		})
		return dte_errors.NewDTEErrorSimple("InvalidSubTotalCalculation",
			expectedSubTotal.InexactFloat64(),
			actualSubTotal.InexactFloat64())
	}

	if taxedAmount.GreaterThan(decimal.Zero) {
		taxedWithDiscount := taxedAmount.
			Sub(decimal.NewFromFloat(s.Document.InvoiceSummary.TaxedDiscount.GetValue()))
		expectedIVA := taxedWithDiscount.Mul(decimal.NewFromFloat(0.13))

		for _, tax := range s.Document.InvoiceSummary.TotalTaxes {
			if tax.GetCode() == constants.TaxIVA {
				actualIVA := decimal.NewFromFloat(tax.GetValue())
				if !expectedIVA.Equal(actualIVA) {
					logs.Error("Invalid IVA calculation with discount", map[string]interface{}{
						"expected":      expectedIVA,
						"actual":        actualIVA,
						"taxedAmount":   taxedAmount,
						"taxedDiscount": s.Document.InvoiceSummary.TaxedDiscount.GetValue(),
					})
					return dte_errors.NewDTEErrorSimple("InvalidIVACalculation",
						expectedIVA.InexactFloat64(),
						actualIVA.InexactFloat64())
				}
				break
			}
		}
	}

	// Inicializar el total a pagar con el total operación
	totalToPay := totalOperation

	if taxedAmount.GreaterThan(decimal.Zero) {
		// Agregar percepción
		perception := decimal.NewFromFloat(s.Document.InvoiceSummary.IVAPerception.GetValue())
		totalToPay = totalToPay.Add(perception)

		// Restar retención IVA
		ivaRetention := decimal.NewFromFloat(s.Document.InvoiceSummary.IVARetention.GetValue())
		totalToPay = totalToPay.Sub(ivaRetention)

		// Restar retención de renta
		incomeRetention := decimal.NewFromFloat(s.Document.InvoiceSummary.IncomeRetention.GetValue())
		totalToPay = totalToPay.Sub(incomeRetention)
	}

	// Agregar monto no gravado si existe
	totalNonTaxed := decimal.NewFromFloat(s.Document.InvoiceSummary.TotalNonTaxed.GetValue())
	if totalNonTaxed.GreaterThan(decimal.Zero) {
		totalToPay = totalToPay.Add(totalNonTaxed)
	}

	actualTotalToPay := decimal.NewFromFloat(s.Document.InvoiceSummary.TotalToPay.GetValue())

	// Usar una pequeña tolerancia para comparaciones con decimales
	diff := totalToPay.Sub(actualTotalToPay).Abs()
	if diff.GreaterThan(decimal.NewFromFloat(0.01)) {
		logs.Error("Invalid total to pay", map[string]interface{}{
			"calculated":      totalToPay,
			"declared":        actualTotalToPay,
			"difference":      diff,
			"operation":       totalOperation,
			"perception":      s.Document.InvoiceSummary.IVAPerception.GetValue(),
			"ivaRetention":    s.Document.InvoiceSummary.IVARetention.GetValue(),
			"incomeRetention": s.Document.InvoiceSummary.IncomeRetention.GetValue(),
		})
		return dte_errors.NewDTEErrorSimple("InvalidTotalToPayCalculation",
			totalToPay.InexactFloat64(),
			actualTotalToPay.InexactFloat64())
	}

	return nil
}

func (s *InvoiceTaxStrategy) validatePerception() *dte_errors.DTEError {
	baseTaxed := decimal.NewFromFloat(s.Document.InvoiceSummary.TotalTaxed.GetValue())
	if baseTaxed.GreaterThan(decimal.Zero) && s.Document.InvoiceSummary.IVAPerception.GetValue() > 0 {
		expectedPerception := baseTaxed.Mul(decimal.NewFromFloat(0.01))
		actualPerception := decimal.NewFromFloat(s.Document.InvoiceSummary.IVAPerception.GetValue())

		diff := expectedPerception.Sub(actualPerception).Abs()
		if diff.GreaterThan(decimal.NewFromFloat(0.01)) {
			logs.Error("Invalid perception amount", map[string]interface{}{
				"expected":   expectedPerception,
				"actual":     actualPerception,
				"taxedBase":  baseTaxed,
				"difference": diff,
			})
			return dte_errors.NewDTEErrorSimple("InvalidPerceptionAmount",
				actualPerception.StringFixed(2),
				expectedPerception.StringFixed(2))
		}
	}

	return nil
}

func ValidateMonetaryAmount(amount float64, fieldName string) *dte_errors.DTEError {
	decValue := decimal.NewFromFloat(amount)
	multiplier := decimal.NewFromInt(100)
	scaled := decValue.Mul(multiplier)

	if !scaled.Equal(decimal.NewFromInt(scaled.IntPart())) {
		return dte_errors.NewDTEErrorSimple("InvalidMonetaryAmount",
			fieldName,
			fmt.Sprintf("Must be multiple of 0.01, got: %v", amount))
	}

	return nil
}

func (s *InvoiceTaxStrategy) validateIVA() *dte_errors.DTEError {
	baseTaxed := decimal.NewFromFloat(s.Document.InvoiceSummary.TotalTaxed.GetValue())

	// Si no hay monto gravado, no debe haber impuestos
	if !baseTaxed.GreaterThan(decimal.Zero) {
		if len(s.Document.InvoiceSummary.TotalTaxes) > 0 {
			return dte_errors.NewDTEErrorSimple("InvalidTaxes")
		}
		return nil
	}

	// Validar cada impuesto
	for _, tax := range s.Document.InvoiceSummary.TotalTaxes {
		if err := s.validateTaxCalculation(tax, baseTaxed); err != nil {
			return err
		}
	}

	// Validar los totales de impuestos
	if err := s.validateSummaryTaxes(); err != nil {
		return err
	}

	return nil
}

func (s *InvoiceTaxStrategy) validateMonetaryAmounts() *dte_errors.DTEError {
	// Validar Total Operation
	if err := ValidateMonetaryAmount(s.Document.InvoiceSummary.TotalOperation.GetValue(), "total_operation"); err != nil {
		return err
	}

	// Validar IVA Retention
	if err := ValidateMonetaryAmount(s.Document.InvoiceSummary.IVARetention.GetValue(), "iva_retention"); err != nil {
		return err
	}

	// Validar Income Retention
	if err := ValidateMonetaryAmount(s.Document.InvoiceSummary.IncomeRetention.GetValue(), "income_retention"); err != nil {
		return err
	}

	// Validar Total To Pay
	if err := ValidateMonetaryAmount(s.Document.InvoiceSummary.TotalToPay.GetValue(), "total_to_pay"); err != nil {
		return err
	}

	// Validar Payment Amounts
	for _, payment := range s.Document.InvoiceSummary.GetPaymentTypes() {
		if err := ValidateMonetaryAmount(payment.GetAmount(), "payment_amount"); err != nil {
			return err
		}
	}

	return nil
}

func (s *InvoiceTaxStrategy) validateTaxCalculation(tax interfaces.Tax, baseTaxed decimal.Decimal) *dte_errors.DTEError {
	var expectedTax decimal.Decimal

	switch tax.GetCode() {
	case constants.TaxIVA:
		expectedTax = baseTaxed.Mul(decimal.NewFromFloat(constants.TaxIvaAmount))
	case constants.TaxIVAExport:
		expectedTax = baseTaxed.Mul(decimal.NewFromFloat(constants.TaxIVAExportAmount))
	case constants.TaxTourism:
		expectedTax = baseTaxed.Mul(decimal.NewFromFloat(constants.TaxTourismAmount))
	case constants.TaxTourismAirport:
		expectedTax = decimal.NewFromFloat(constants.TaxTourismAirportAmount)
	case constants.TaxFOVIAL:
		expectedTax = baseTaxed.Mul(decimal.NewFromFloat(constants.TaxFOVIALAmount))
	case constants.TaxCOTRANS:
		expectedTax = decimal.NewFromFloat(constants.TaxCOTRANSAmount)
	case constants.TaxSpecialOther:
		return nil
	}

	actualTax := decimal.NewFromFloat(tax.GetValue())
	if !actualTax.Equal(expectedTax) {
		return dte_errors.NewDTEErrorSimple("InvalidTaxCalculation",
			tax.GetCode(),
			expectedTax.InexactFloat64(),
			actualTax.InexactFloat64())
	}

	return nil
}

func (s *InvoiceTaxStrategy) validateBaseTotals() *dte_errors.DTEError {

	//Verificar que los descuentos no sobrepasen el subtotal
	if decimal.NewFromFloat(s.Document.InvoiceSummary.SubTotal.GetValue()).LessThan(decimal.NewFromFloat(s.Document.InvoiceSummary.TaxedDiscount.GetValue())) {
		logs.Error("Invalid taxed discount", map[string]interface{}{
			"taxedDiscount": s.Document.InvoiceSummary.TaxedDiscount.GetValue(),
			"subTotal":      s.Document.InvoiceSummary.SubTotal.GetValue(),
		})
		return dte_errors.NewDTEErrorSimple("DiscountExceedsSubtotal",
			"TaxedDiscount",
			s.Document.InvoiceSummary.TaxedDiscount.GetValue(),
			s.Document.InvoiceSummary.SubTotal.GetValue())
	}

	if decimal.NewFromFloat(s.Document.InvoiceSummary.SubTotal.GetValue()).LessThan(decimal.NewFromFloat(s.Document.InvoiceSummary.ExemptDiscount.GetValue())) {
		logs.Error("Invalid exempt discount", map[string]interface{}{
			"exemptDiscount": s.Document.InvoiceSummary.ExemptDiscount.GetValue(),
			"subTotal":       s.Document.InvoiceSummary.SubTotal.GetValue(),
		})
		return dte_errors.NewDTEErrorSimple("DiscountExceedsSubtotal",
			"ExemptDiscount",
			s.Document.InvoiceSummary.ExemptDiscount.GetValue(),
			s.Document.InvoiceSummary.SubTotal.GetValue())
	}

	if decimal.NewFromFloat(s.Document.InvoiceSummary.SubTotal.GetValue()).LessThan(decimal.NewFromFloat(s.Document.InvoiceSummary.NonSubjectDiscount.GetValue())) {
		logs.Error("Invalid non subject discount", map[string]interface{}{
			"nonSubjectDiscount": s.Document.InvoiceSummary.NonSubjectDiscount.GetValue(),
			"subTotal":           s.Document.InvoiceSummary.SubTotal.GetValue(),
		})
		return dte_errors.NewDTEErrorSimple("DiscountExceedsSubtotal",
			"NonSubjectDiscount",
			s.Document.InvoiceSummary.NonSubjectDiscount.GetValue(),
			s.Document.InvoiceSummary.SubTotal.GetValue())
	}

	return nil
}

// validateSummaryTaxes valida los totales de impuestos del resumen
func (s *InvoiceTaxStrategy) validateSummaryTaxes() *dte_errors.DTEError {
	// 1. Calcular IVA desde los items
	var totalIVAFromItems decimal.Decimal
	for _, item := range s.Document.InvoiceItems {
		itemIVA := decimal.NewFromFloat(item.IVAItem.GetValue())
		totalIVAFromItems = totalIVAFromItems.Add(itemIVA)
	}

	// 2. Obtener el IVA total declarado
	summaryIVA := decimal.NewFromFloat(s.Document.InvoiceSummary.TotalIva.GetValue())

	// 3. Validar que el IVA calculado coincida con el total declarado
	if !s.CompareTaxWithTolerance(totalIVAFromItems, summaryIVA, 0.01) {
		logs.Error("Invalid IVA total", map[string]interface{}{
			"expectedFromItems": totalIVAFromItems.InexactFloat64(),
			"actual":            summaryIVA.InexactFloat64(),
		})
		return dte_errors.NewDTEErrorSimple("InvalidTotalIVA",
			summaryIVA.InexactFloat64(),
			totalIVAFromItems.InexactFloat64())
	}

	return nil
}

// CompareTaxWithTolerance Compara dos impuestos con una tolerancia dada, normalmente 0.0001
func (s *InvoiceTaxStrategy) CompareTaxWithTolerance(expected, actual decimal.Decimal, tolerance float64) bool {
	diff := expected.Sub(actual).Abs()
	return diff.LessThanOrEqual(decimal.NewFromFloat(tolerance))
}
