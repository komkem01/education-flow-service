package paymenttransactions

import (
	"database/sql"
	"errors"

	"education-flow/app/modules/entities/ent"
	"education-flow/app/utils"
	"education-flow/app/utils/base"
	ci18n "education-flow/config/i18n"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

const dateTimeLayout = "2006-01-02T15:04:05Z07:00"

type Controller struct {
	tracer trace.Tracer
	svc    *Service
}

func newController(trace trace.Tracer, svc *Service) *Controller {
	return &Controller{tracer: trace, svc: svc}
}

type createRequest struct {
	InvoiceID          string   `json:"invoice_id" binding:"required,uuid"`
	AmountPaid         *float64 `json:"amount_paid"`
	PaymentMethod      string   `json:"payment_method" binding:"omitempty,oneof=cash transfer qr_code"`
	EvidenceURL        *string  `json:"evidence_url"`
	TransactionDate    *string  `json:"transaction_date" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	ProcessedByStaffID *string  `json:"processed_by_staff_id" binding:"omitempty,uuid"`
}

type updateRequest struct {
	InvoiceID          string   `json:"invoice_id" binding:"required,uuid"`
	AmountPaid         *float64 `json:"amount_paid"`
	PaymentMethod      string   `json:"payment_method" binding:"required,oneof=cash transfer qr_code"`
	EvidenceURL        *string  `json:"evidence_url"`
	TransactionDate    *string  `json:"transaction_date" binding:"omitempty,datetime=2006-01-02T15:04:05Z07:00"`
	ProcessedByStaffID *string  `json:"processed_by_staff_id" binding:"omitempty,uuid"`
}

type response struct {
	ID                 string   `json:"id"`
	InvoiceID          string   `json:"invoice_id"`
	AmountPaid         *float64 `json:"amount_paid"`
	PaymentMethod      string   `json:"payment_method"`
	EvidenceURL        *string  `json:"evidence_url"`
	TransactionDate    *string  `json:"transaction_date"`
	ProcessedByStaffID *string  `json:"processed_by_staff_id"`
}

func (c *Controller) Create(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, _, ok := parseIDs(ctx, false)
	if !ok {
		return
	}

	var req createRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	invoiceID, err := uuid.Parse(req.InvoiceID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	transactionDate, err := utils.ParseTimePtrWithLayout(req.TransactionDate, dateTimeLayout)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	processedByStaffID, err := utils.ParseUUIDPtr(req.ProcessedByStaffID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	paymentMethod := ent.PaymentMethodCash
	if req.PaymentMethod != "" {
		paymentMethod = ent.ToPaymentMethod(req.PaymentMethod)
	}

	item, err := c.svc.Create(ctx.Request.Context(), &CreateInput{StudentID: studentID, InvoiceID: invoiceID, AmountPaid: req.AmountPaid, PaymentMethod: paymentMethod, EvidenceURL: req.EvidenceURL, TransactionDate: transactionDate, ProcessedByStaffID: processedByStaffID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.StudentInvoiceNotFound, nil)
			return
		}
		log.Errf("payment-transactions.create.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func (c *Controller) List(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, _, ok := parseIDs(ctx, false)
	if !ok {
		return
	}
	items, err := c.svc.ListByStudentID(ctx.Request.Context(), studentID)
	if err != nil {
		log.Errf("payment-transactions.list.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	responseList := make([]response, 0, len(items))
	for _, item := range items {
		responseList = append(responseList, toResponse(item))
	}
	base.Success(ctx, responseList)
}

func (c *Controller) Update(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}
	var req updateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	invoiceID, err := uuid.Parse(req.InvoiceID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}
	transactionDate, err := utils.ParseTimePtrWithLayout(req.TransactionDate, dateTimeLayout)
	if err != nil {
		base.BadRequest(ctx, ci18n.BadRequest, nil)
		return
	}
	processedByStaffID, err := utils.ParseUUIDPtr(req.ProcessedByStaffID)
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return
	}

	item, err := c.svc.UpdateByID(ctx.Request.Context(), studentID, childID, &UpdateInput{StudentID: studentID, InvoiceID: invoiceID, AmountPaid: req.AmountPaid, PaymentMethod: ent.ToPaymentMethod(req.PaymentMethod), EvidenceURL: req.EvidenceURL, TransactionDate: transactionDate, ProcessedByStaffID: processedByStaffID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.PaymentTransactionNotFound, nil)
			return
		}
		log.Errf("payment-transactions.update.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, toResponse(item))
}

func (c *Controller) Delete(ctx *gin.Context) {
	_, log := utils.LogSpanFromGin(ctx)
	studentID, childID, ok := parseIDs(ctx, true)
	if !ok {
		return
	}
	if err := c.svc.DeleteByID(ctx.Request.Context(), studentID, childID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			base.ValidateFailed(ctx, ci18n.PaymentTransactionNotFound, nil)
			return
		}
		log.Errf("payment-transactions.delete.error: %v", err)
		base.InternalServerError(ctx, ci18n.InternalServerError, nil)
		return
	}
	base.Success(ctx, gin.H{"id": childID.String()})
}

func parseIDs(ctx *gin.Context, childRequired bool) (uuid.UUID, uuid.UUID, bool) {
	studentID, err := utils.ParsePathUUID(ctx, "id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, false
	}
	if !childRequired {
		return studentID, uuid.Nil, true
	}
	childID, err := utils.ParsePathUUID(ctx, "child_id")
	if err != nil {
		base.BadRequest(ctx, ci18n.InvalidID, nil)
		return uuid.Nil, uuid.Nil, false
	}
	return studentID, childID, true
}

func toResponse(item *ent.PaymentTransaction) response {
	return response{ID: item.ID.String(), InvoiceID: item.InvoiceID.String(), AmountPaid: item.AmountPaid, PaymentMethod: string(item.PaymentMethod), EvidenceURL: item.EvidenceURL, TransactionDate: utils.TimeToStringPtrWithLayout(item.TransactionDate, dateTimeLayout), ProcessedByStaffID: utils.UUIDToStringPtr(item.ProcessedByStaffID)}
}
