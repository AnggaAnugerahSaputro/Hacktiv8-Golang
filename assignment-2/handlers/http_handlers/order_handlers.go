package http_handlers

import (
	"assignment_2/dto"
	"assignment_2/pkg/errs"
	"assignment_2/pkg/helpers"

	// "assignment_2/pkg/helpers"
	"assignment_2/services"

	"github.com/gin-gonic/gin"
)

type orderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(orderService services.OrderService) *orderHandler {
	return &orderHandler{orderService: orderService}
}


// CreateOrders godoc
//
//	@Summary		Create a Orders
//	@Description	Create a Orders by json
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			orders  body        dto.NewOrderRequest	    true	   "Create orders request body"
//	@Success		201		{object}	dto.NewOrderResponse
//	@Failure		401		{object}	errs.MessageErrData
//	@Failure		422		{object}	errs.MessageErrData
//	@Failure		500		{object}	errs.MessageErrData
//	@Router			/orders [post]
func (o *orderHandler) CreateOrder(ctx *gin.Context) {
	var orderRequest dto.NewOrderRequest

		if err := ctx.ShouldBindJSON(&orderRequest); err != nil {
			errMessage := errs.NewUnprocessableEntity("Invalid request body")

			ctx.JSON(errMessage.StatusCode(), errMessage)
			return
		}

		result, err := o.orderService.CreateOrder(orderRequest)

		if err != nil {
			errReq := errs.NewBadRequest("Bad Request")

			ctx.JSON(errReq.StatusCode(), errReq)
			return
		}
		ctx.JSON(result.StatusCode, result)

}

// GetAllOrders godoc
//
//	@Summary		Get all orders
//	@Description	Get all orders by json
//	@Tags			orders
//	@Produce		json
//	@Success		200		{object}	dto.GetAllOrderResponse
//	@Failure		401		{object}	errs.MessageErrData
//	@Failure		500		{object}	errs.MessageErrData
//	@Router			/orders [get]
func (o *orderHandler) GetAllOrders(ctx *gin.Context) {
	result, err := o.orderService.GetAllOrder()
	if err != nil {
		ctx.JSON(err.StatusCode(), err)
		return
	}
	ctx.JSON(result.StatusCode, result)
}

// GetorderById godoc
//
//	@Summary		Get order by id
//	@Description	Get order by id json
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			orderId		path		uint						true	"orders ID request"
//	@Success		200			{object}	dto.NewOrderUpdateResponse
//	@Failure		401			{object}	errs.MessageErrData
//	@Failure		400			{object}	errs.MessageErrData
//	@Failure		422			{object}	errs.MessageErrData
//	@Failure		500			{object}	errs.MessageErrData
//	@Router			/orders/{orderId} [get]
func (o *orderHandler) GetOrderById(ctx *gin.Context) {
	orderId, err := helpers.GetParamId(ctx, "orderId")

	if err != nil {
		ctx.JSON(err.StatusCode(), err)
		return
	}

	response, err := o.orderService.GetOrderById(orderId)
	if err != nil {
		ctx.JSON(err.StatusCode(), err)
		return
	}

	ctx.JSON(response.StatusCode, response)
}

// UpdateOrders godoc
//
//	@Summary		Update a orders 
//	@Description	Update a orders 
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Param			orders		body		dto.NewOrderRequest			true	"Update orders request body"
//	@Param			orderId		path		uint						true	"orders ID request"
//	@Success		200			{object}	dto.NewOrderUpdateResponse
//	@Failure		401			{object}	errs.MessageErrData
//	@Failure		400			{object}	errs.MessageErrData
//	@Failure		422			{object}	errs.MessageErrData
//	@Failure		500			{object}	errs.MessageErrData
//	@Router			/orders/{orderId} [put]
func (o *orderHandler) UpdateOrder(ctx *gin.Context) {
	var orderRequest dto.NewOrderRequest

	if err := ctx.ShouldBindJSON(&orderRequest); err != nil {
		error := errs.NewUnprocessableEntity("invalid request body")
		ctx.JSON(error.StatusCode(), error)
		return
	}

	orderId, err := helpers.GetParamId(ctx, "orderId")
	if err != nil {
		ctx.JSON(err.StatusCode(), err)
		return
	}

	result, err := o.orderService.UpdateOrder(orderId, orderRequest)

	if err != nil {
		ctx.JSON(err.StatusCode(), err)
		return
	}

	ctx.JSON(result.StatusCode, result)
}

// DeleteOrders godoc
//
//	@Summary		Delete a orders
//	@Description	Delete a orders by param
//	@Tags			orders
//	@Produce		json
//	@Param			orderId	 	path		uint						true	"order ID request"
//	@Success		200			{object}	dto.DeleteOrderResponse
//	@Failure		401			{object}	errs.MessageErrData
//	@Failure		400			{object}	errs.MessageErrData
//	@Router			/orders/{orderId} [delete]
func (o *orderHandler) DeleteOrders(ctx *gin.Context) {
	orderID, err := helpers.GetParamId(ctx, "orderId")
	if err != nil {
		ctx.JSON(err.StatusCode(), err)
		return
	}

	result, err := o.orderService.DeleteOrders(orderID)
	if err != nil {
		ctx.JSON(err.StatusCode(), err)
		return
	}

	ctx.JSON(result.StatusCode, result)
}
