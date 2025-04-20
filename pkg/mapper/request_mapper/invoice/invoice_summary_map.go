package invoice

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/invoice/invoice_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// MapInvoiceRequestSummary mapea un resumen de invoice a un modelo de resumen de invoice -> Origen: Request
func MapInvoiceRequestSummary(summary *structs.InvoiceSummaryRequest) (*invoice_models.InvoiceSummary, error) {
	if summary.TotalInWords == nil {
		inLetters := utils.InLetters(summary.TotalToPay)
		summary.TotalInWords = &inLetters
	}

	baseSummary, err := common.MapCommonRequestSummary(structs.SummaryRequest{
		TotalNonSubject:    summary.TotalNonSubject,
		TotalExempt:        summary.TotalExempt,
		TotalTaxed:         summary.TotalTaxed,
		SubTotal:           summary.SubTotal,
		NonSubjectDiscount: summary.NonSubjectDiscount,
		ExemptDiscount:     summary.ExemptDiscount,
		DiscountPercentage: summary.DiscountPercentage,
		TotalDiscount:      summary.TotalDiscount,
		TotalOperation:     summary.TotalOperation,
		TotalNonTaxed:      summary.TotalNonTaxed,
		SubTotalSales:      summary.SubTotalSales,
		TotalToPay:         summary.TotalToPay,
		OperationCondition: summary.OperationCondition,
		Taxes:              summary.Taxes,
		PaymentTypes:       summary.PaymentTypes,
		TotalInWords:       summary.TotalInWords,
	})
	if err != nil {
		return nil, err
	}

	taxedDiscount, err := financial.NewAmountForTotal(summary.TaxedDiscount)
	if err != nil {
		return nil, err
	}

	ivaRetention, err := financial.NewAmountForTotal(summary.IVARetention)
	if err != nil {
		return nil, err
	}

	incomeRetention, err := financial.NewAmountForTotal(summary.IncomeRetention)
	if err != nil {
		return nil, err
	}

	totalIva, err := financial.NewAmountForTotal(summary.TotalIVA)
	if err != nil {
		return nil, err
	}

	balanceInFavor, err := financial.NewAmountForTotal(summary.BalanceInFavor)
	if err != nil {
		return nil, err
	}

	return &invoice_models.InvoiceSummary{
		Summary:                 baseSummary,
		TaxedDiscount:           *taxedDiscount,
		IVARetention:            *ivaRetention,
		IncomeRetention:         *incomeRetention,
		TotalIva:                *totalIva,
		BalanceInFavor:          *balanceInFavor,
		ElectronicPaymentNumber: new(string),
	}, nil
}
