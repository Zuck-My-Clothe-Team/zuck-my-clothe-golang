package controller

import (
	"net/http"
	"strings"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"
	vboi "zuck-my-clothe/zuck-my-clothe-backend/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type OrderController interface {
	CreateNewOrder(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetByHeaderID(c *fiber.Ctx) error
	GetByBranchID(c *fiber.Ctx) error
	GetByUserID(c *fiber.Ctx) error
	UpdateStatus(c *fiber.Ctx) error
	UpdateReview(c *fiber.Ctx) error
	SoftDelete(c *fiber.Ctx) error
}

type orderController struct {
	orderUsecase usecases.OrderUsecase
}

func CreateOrderController(orderUsecase usecases.OrderUsecase) OrderController {
	return &orderController{orderUsecase: orderUsecase}
}

func getCookieData(c *fiber.Ctx, key string) string {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	return claims[key].(string)
}

//	@Summary		Add new order
//	@Description	Add a new order to the system
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			NewOrder	body		model.NewOrder	true	"New Order Data"
//	@Success		201			{object}	model.FullOrder	"Created"
//	@Failure		400			{string}	string			"Bad Request - Invalid input"
//	@Failure		418			{string}	string			"ERR: mai wang ja"
//	@Failure		406			{string}	string			"Not Acceptable - Validation failed"
//	@Failure		500			{string}	string			"Internal Server Error"
//	@Router			/order/new [post]
func (u *orderController) CreateNewOrder(c *fiber.Ctx) error {
	newOrder := new(model.NewOrder)

	if err := c.BodyParser(newOrder); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	if err := vboi.Validate(newOrder); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	newOrder.UserID = getCookieData(c, "userID")

	response, err := u.orderUsecase.CreateNewOrder(newOrder)

	if err != nil {
		if err.Error() == "ERR: mai wang ja" {
			return c.Status(http.StatusTeapot).SendString(err.Error())
		} else if err.Error() == "null detected on one or more essential field(s)" {
			return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

//	@Summary		Get all orders
//	@Description	Retrieve all orders in the system
//	@Tags			Order
//	@Produce		json
//	@Success		200	{array}		model.FullOrder	"OK"
//	@Failure		404	{string}	string			"Not Found - No orders available"
//	@Failure		500	{string}	string			"Internal Server Error"
//	@Router			/order/all [get]
func (u *orderController) GetAll(c *fiber.Ctx) error {
	result, err := u.orderUsecase.GetAll()

	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNoContent)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

//	@Summary		Get full order by order header id
//	@Description	Retrieve full order by order header id
//	@Tags			Order
//	@Produce		json
//	@Param			order_header_id	path		string			true	"Order Header ID"
//	@Param			option			path		string			true	"Option = full, header, detail"
//	@Success		200				{array}		model.FullOrder	"OK"
//	@Failure		404				{string}	string			"Not Found - No orders available"
//	@Failure		500				{string}	string			"Internal Server Error"
//	@Router			/order/{order_header_id}/{option} [get]
func (u *orderController) GetByHeaderID(c *fiber.Ctx) error {
	orderHeaderID := c.Params("order_header_id")
	option := c.Params("option")

	if option != "full" && option != "header" && option != "detail" {
		return c.Status(fiber.StatusNotAcceptable).SendString("err: not valid option")
	}

	role := getCookieData(c, "positionID")

	var isAdminView bool = false

	if role == "BranchManager" || role == "SuperAdmin" {
		isAdminView = true
	}

	result, err := u.orderUsecase.GetByHeaderID(orderHeaderID, isAdminView, option)

	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNoContent)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

//	@Summary		Get full order by branch id
//	@Description	Retrieve full order by branch id
//	@Tags			Order
//	@Produce		json
//	@Param			branch_id	path		string			true	"branch id"
//	@Success		200			{array}		model.FullOrder	"OK"
//	@Failure		404			{string}	string			"Not Found - No orders available"
//	@Failure		500			{string}	string			"Internal Server Error"
//	@Router			/order/branch/{branch_id} [get]
func (u *orderController) GetByBranchID(c *fiber.Ctx) error {
	branchID := c.Params("branch_id")

	managerID := getCookieData(c, "userID")

	result, err := u.orderUsecase.GetByBranchID(branchID, managerID)

	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNoContent)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

//	@Summary		Get full order by user id
//	@Description	Retrieve full order by user id
//	@Tags			Order
//	@Produce		json
//	@Param			status	query		string			true	"status: waiting, processing, completed, expired"
//	@Success		200		{array}		model.FullOrder	"OK"
//	@Failure		404		{string}	string			"Not Found - No orders available"
//	@Failure		500		{string}	string			"Internal Server Error"
//	@Router			/order/me [get]
func (u *orderController) GetByUserID(c *fiber.Ctx) error {
	userID := getCookieData(c, "userID")

	status := c.Query("status")

	status = strings.ToUpper(status[:1]) + status[1:]

	if status != string(model.Waiting) &&
		status != string(model.Processing) &&
		status != string(model.Completed) &&
		status != string(model.Expired) &&
		status != "" {
		return c.Status(fiber.StatusBadRequest).SendString("ERR: status option is not valid")
	}

	result, err := u.orderUsecase.GetByUserID(userID, status)

	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNoContent)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

//	@Summary		Update order status
//	@Description	Update order status
//	@Tags			Order
//	@Param			UpdateOrder	body		model.UpdateOrder	true	"Updated order field"
//	@Success		200			{object}	model.FullOrder		"OK"
//	@Failure		400			{string}	string				"Bad Request"
//	@Failure		404			{string}	string				"Not Found"
//	@Failure		500			{string}	string				"Internal Server Error"
//	@Router			/order/update [put]
func (u *orderController) UpdateStatus(c *fiber.Ctx) error {
	order := new(model.UpdateOrder)

	if err := c.BodyParser(order); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	if err := vboi.Validate(order); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	order.UpdatedBy = getCookieData(c, "userID")

	if order.UpdatedBy == "" {
		return c.SendStatus(fiber.StatusNotAcceptable)
	}

	result, err := u.orderUsecase.UpdateStatus(*order)

	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNoContent)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

//	@Summary		Update order review
//	@Description	Update order review
//	@Tags			Order
//	@Param			OrderReview	body		model.OrderReview	true	"Updated order field"
//	@Success		200			{object}	model.FullOrder		"OK"
//	@Failure		400			{string}	string				"Bad Request"
//	@Failure		404			{string}	string				"Not Found"
//	@Failure		500			{string}	string				"Internal Server Error"
//	@Router			/order/review [put]
func (u *orderController) UpdateReview(c *fiber.Ctx) error {
	order := new(model.OrderReview)
	userID := getCookieData(c, "userID")
	order.UserID = userID
	if err := c.BodyParser(order); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	if err := vboi.Validate(order); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	result, err := u.orderUsecase.UpdateReview(*order)

	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNoContent)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

//	@Summary		Delete an order
//	@Description	Delete an order
//	@Tags			Order
//	@Accept			json
//	@Produce		json
//	@Param			order_header_id	path		string	true	"Order Header ID"
//	@Success		200				{string}	string	"ok"
//	@Failure		404				{string}	string	"record not found"
//	@Failure		500				{string}	string	"internal server error"
//	@Router			/order/delete/{order_header_id} [delete]
func (u *orderController) SoftDelete(c *fiber.Ctx) error {
	orderHeaderID := c.Params("order_header_id")

	deleteBy := getCookieData(c, "userID")

	result, err := u.orderUsecase.SoftDelete(orderHeaderID, deleteBy)
	if err != nil {
		if err.Error() == "record not found" {
			return c.SendStatus(fiber.StatusNoContent)
		} else {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
