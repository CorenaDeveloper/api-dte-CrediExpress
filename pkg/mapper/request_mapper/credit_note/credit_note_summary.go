package credit_note

import (
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/common/value_objects/financial"
	"github.com/MarlonG1/api-facturacion-sv/internal/domain/dte/credit_note/credit_note_models"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/common"
	"github.com/MarlonG1/api-facturacion-sv/pkg/mapper/request_mapper/structs"
	"github.com/MarlonG1/api-facturacion-sv/pkg/shared/utils"
)

// MapCreditNoteRequestSummary mapea un resumen de Nota de Crédito a un modelo de resumen de Nota de Crédito -> Origen: Request
func MapCreditNoteRequestSummary(summary *structs.CreditNoteSummaryRequest) (*credit_note_models.CreditNoteSummary, error) {
	if summary.TotalInWords == nil {
		inLetters := utils.InLetters(summary.TotalOperation)
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
		TotalToPay:         1,
		OperationCondition: summary.OperationCondition,
		Taxes:              summary.Taxes,
		PaymentTypes:       []structs.PaymentRequest{},
		TotalInWords:       summary.TotalInWords,
	})

	if err != nil {
		return nil, err
	}

	taxedDiscount, err := financial.NewAmountForTotal(summary.TaxedDiscount)
	if err != nil {
		return nil, err
	}

	ivaPerception, err := financial.NewAmountForTotal(summary.IVAPerception)
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

	balanceInFavor, err := financial.NewAmountForTotal(summary.BalanceInFavor)
	if err != nil {
		return nil, err
	}

	return &credit_note_models.CreditNoteSummary{
		Summary:         baseSummary,
		TaxedDiscount:   *taxedDiscount,
		IVAPerception:   *ivaPerception,
		IVARetention:    *ivaRetention,
		IncomeRetention: *incomeRetention,
		BalanceInFavor:  *balanceInFavor,
	}, nil
}
