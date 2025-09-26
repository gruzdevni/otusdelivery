package restapi

import (
	"strconv"

	"otusdelivery/internal/models"
	query "otusdelivery/internal/repo"
	deliveryOps "otusdelivery/internal/restapi/operations/delivery"
	"otusdelivery/internal/restapi/operations/other"
	"otusdelivery/internal/service/api/delivery"

	"github.com/go-openapi/runtime/middleware"
)

type Handler struct {
	delivertSrv delivery.Service
}

func NewHandler(delivertSrv delivery.Service) *Handler {
	return &Handler{
		delivertSrv: delivertSrv,
	}
}

func (h *Handler) GetHealth(_ other.GetHealthParams) middleware.Responder {
	return other.NewGetHealthOK().WithPayload(&models.DefaultStatusResponse{Code: "01", Message: "OK"})
}

func (h *Handler) CheckOrderDelivery(params deliveryOps.GetCheckDeliveryStatusOrderIDParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	res, err := h.delivertSrv.CheckOrderDelivery(ctx, params.OrderID)
	if err != nil {
		return deliveryOps.NewGetCheckDeliveryStatusOrderIDInternalServerError()
	}

	return deliveryOps.NewGetCheckDeliveryStatusOrderIDOK().WithPayload(res)
}

func (h *Handler) ReserveSlotForOrder(params deliveryOps.PostReserveSlotParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	slotID, err := strconv.Atoi(params.Request.SlotID)
	if err != nil {
		return deliveryOps.NewPostReserveSlotBadRequest().WithPayload(&models.DefaultStatusResponse{Code: "03", Message: err.Error()})
	}

	err = h.delivertSrv.ReserveSlotForOrder(ctx, query.ReserveSlotForOrderParams{
		OrderID: params.Request.OrderID,
		SlotID:  int32(slotID),
	})
	if err != nil {
		return deliveryOps.NewPostReserveSlotInternalServerError().WithPayload(&models.DefaultStatusResponse{Code: "03", Message: err.Error()})
	}

	return deliveryOps.NewPostReserveSlotOK()
}

func (h *Handler) UnreserveSlotForOrder(params deliveryOps.DeleteUnreserveSlotOrderIDParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	err := h.delivertSrv.UnreserveSlotForOrder(ctx, params.OrderID)
	if err != nil {
		return deliveryOps.NewDeleteUnreserveSlotOrderIDInternalServerError().WithPayload(&models.DefaultStatusResponse{Code: "03", Message: err.Error()})
	}

	return deliveryOps.NewDeleteUnreserveSlotOrderIDOK()
}
