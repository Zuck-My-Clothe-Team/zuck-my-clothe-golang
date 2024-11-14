package controller

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	validatorboi "zuck-my-clothe/zuck-my-clothe-backend/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserAddressesController interface {
	AddUserAddress(c *fiber.Ctx) error
	FindByAddressID(c *fiber.Ctx) error
	FindUserAddresByOwnerID(c *fiber.Ctx) error
	UpdateUserAddressData(c *fiber.Ctx) error
	DeleteUserAddress(c *fiber.Ctx) error
}

type userAddressesController struct {
	userAddressesUsecase model.UserAddressesUsecase
}

func CreateNewUserAddressesController(userAddressesUsecase model.UserAddressesUsecase) userAddressesController {
	return userAddressesController{userAddressesUsecase: userAddressesUsecase}
}

//	@Summary		Add new user address
//	@Description	Add a new user address information to system
//	@Tags			UserAddress
//	@Accept			json
//	@Produce		json
//	@Param			UserAddresses	body		model.AddUserAddressDTO	true	"New Address Data"
//	@Success		201				{object}	model.UserAddressDetail	"Created"
//	@Failure		406				{string}	string					"Not Acceptable"
//	@Failure		500				{string}	string					"internal server error"
//	@Router			/address/add [post]
func (u *userAddressesController) AddUserAddress(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	var ownerID = claims["userID"].(string)

	newUserAddress := new(model.AddUserAddressDTO)
	if err := c.BodyParser(newUserAddress); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	if err := validatorboi.Validate(newUserAddress); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}
	createdRecord, err := u.userAddressesUsecase.AddUserAddress(ownerID, newUserAddress)
	if err != nil {
		return c.Status(fiber.StatusAccepted).SendString(err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(createdRecord)
}

//	@Summary		Find address by Address id
//	@Description	Find requested address by address id
//	@Tags			UserAddress
//	@Param			addressID	path	string	true	"Address ID"
//	@Produce		json
//	@Success		200	{object}	model.UserAddressDetail
//	@Failure		202	{string}	string	"Accepted"
//	@Router			/detail/aid/{addressID} [get]
func (u *userAddressesController) FindByAddressID(c *fiber.Ctx) error {
	addressID := c.Params("addressID")

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	var ownerID = claims["userID"].(string)

	addressDetail, err := u.userAddressesUsecase.FindUserAddressByID(addressID, ownerID)
	if err != nil {
		return c.Status(fiber.StatusAccepted).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(addressDetail)
}

//	@Summary		Find owned address
//	@Description	List all owned address of thant user
//	@Tags			UserAddress
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.UserAddressDetail[]
//	@Failure		202	{string}	string	"Accepted"
//	@Router			/address/detail/owner [get]
func (u *userAddressesController) FindUserAddresByOwnerID(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	var ownerID = claims["userID"].(string)

	addresses, err := u.userAddressesUsecase.FindUserAddresByOwnerID(ownerID)
	if err != nil {
		return c.Status(fiber.StatusAccepted).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(addresses)
}

//	@Summary		Update address information
//	@Description	Update requested address information
//	@Tags			UserAddress
//	@Accept			json
//	@Produce		json
//	@Param			UserAddresses	body		model.UpdateUserAddressDTO	true	"New Address Data"
//	@Success		200				{object}	model.UserAddressDetail
//	@Failure		202				{string}	string	"Accepted"
//	@Failure		204				{string}	string	"record not found"
//	@Failure		406				{string}	string	"Not Acceptable"
//	@Router			/address/update [put]
func (u *userAddressesController) UpdateUserAddressData(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	var ownerID = claims["userID"].(string)

	updateData := new(model.AddUserAddressDTO)
	if err := c.BodyParser(updateData); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(err.Error())
	}

	updatedData, err := u.userAddressesUsecase.UpdateUserAddressData(ownerID, updateData)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNoContent).SendString("record not found")
		}
		return c.Status(fiber.StatusAccepted).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(updatedData)
}

//	@Summary		Delete address
//	@Description	Delete requested address from system
//	@Tags			UserAddress
//	@Accept			json
//	@Produce		json
//	@Param			address_id	path		string	true	"Address ID"
//	@Success		200			{string}	string	"Ok"
//	@Failure		202			{string}	string	"Accepted"
//	@Failure		204			{string}	string	"record not found"
//	@Router			/address/delete/{addressID} [delete]
func (u *userAddressesController) DeleteUserAddress(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	var ownerID = claims["userID"].(string)

	var addressID string = c.Params("addressID")
	err := u.userAddressesUsecase.DeleteUserAddress(ownerID, addressID)
	if err != nil {
		if err.Error() == "record not found" {
			return c.Status(fiber.StatusNoContent).SendString("record not found")
		}
		return c.Status(fiber.StatusAccepted).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusOK)
}
