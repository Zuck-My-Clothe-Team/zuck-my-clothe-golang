package controller

import (
	"fmt"
	"zuck-my-clothe/zuck-my-clothe-backend/model"

	"github.com/gofiber/fiber/v2"
	_ "github.com/golang-jwt/jwt/v5"
)

type PaymentController interface {
	CreatePayment(c *fiber.Ctx) error
	FindByPaymentID(c *fiber.Ctx) error
	UpdatePaymenstatus(c *fiber.Ctx) error
}

type paymentController struct {
	paymentUsecase model.PaymentUsecase
}

func CreateNewPaymentController(paymentUsecase model.PaymentUsecase) PaymentController {
	return &paymentController{paymentUsecase: paymentUsecase}
}

//	@Summary		Add new payment
//	@Description	Add a new payment record to db [mockup]
//	@Tags			Payment
//	@Accept			json
//	@Produce		json
//	@Param			PaymentModel	body		model.Payments	true	"New Payment Data"
//	@Success		201				{object}	model.Payments	"Created"
//	@Failure		202				{string}	string			"Accepted"
//	@Router			/payment/add [post]
func (u *paymentController) CreatePayment(c *fiber.Ctx) error {
	newPayment := new(model.Payments)
	if err := c.BodyParser(newPayment); err != nil {
		return c.SendStatus(fiber.StatusNotAcceptable)
	}
	createdPayment, err := u.paymentUsecase.CreatePayment(*newPayment)
	if err != nil {
		return c.Status(fiber.StatusAccepted).SendString(err.Error())
	}
	fmt.Println(createdPayment)
	return c.Status(fiber.StatusCreated).JSON(createdPayment)
}

//	@Summary		Find payment by id
//	@Description	Find payment by paymentID [mockup]
//	@Tags			Payment
//	@Produce		json
//	@Param			paymentID	path		string			true	"PaymentID"
//	@Success		200			{object}	model.Payments	"OK"
//	@Failure		204			{string}	string			"no content"
//	@Failure		500			{string}	string			"Internal Server Error"
//	@Router			/payment/detail/{paymentID} [get]
func (u *paymentController) FindByPaymentID(c *fiber.Ctx) error {
	paymentID := c.Params("paymentID")
	data, err := u.paymentUsecase.FindByPaymentID(paymentID)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNoContent).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(data)
}

//	@Summary		Update payment status
//	@Description	Set the status of payment
//	@Tags			Payment
//	@Param			paymentID	path		string			true	"Machine Serial ID"
//	@Param			status		path		string			true	"Set status (Pending/Paid/Expired/Cancel)"
//	@Success		200			{object}	model.Payments	"OK"
//	@Failure		202			{string}	string			"Accepted"
//	@Failure		406			{string}	string			"err: not valid status"
//	@Router			/payment/update/{paymentID}/setstatus/{status} [put]
func (u *paymentController) UpdatePaymenstatus(c *fiber.Ctx) error {
	paymentID := c.Params("paymentID")
	status := c.Params("status")
	if status != "Pending" && status != "Paid" && status != "Expired" && status != "Cancel" {
		return c.Status(fiber.StatusNotAcceptable).SendString("err: not valid status")
	}
	reponse, err := u.paymentUsecase.UpdatePaymentStatus(paymentID, model.PaymentStatus(status))
	if err != nil {
		return c.Status(fiber.StatusAccepted).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(reponse)
}
