package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/config"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/usecases"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type AuthenController interface {
	SignIn(c *fiber.Ctx) error
	Me(c *fiber.Ctx) error
	GoogleCallback(c *fiber.Ctx) error
}

type authenUsecase struct {
	usecase     model.AuthenUsecase
	userUsecase usecases.UserUsecases
	cfg         *config.Config
}

func CreateNewAuthenController(usecase model.AuthenUsecase, userUsecase usecases.UserUsecases, cfg *config.Config) AuthenController {
	return &authenUsecase{usecase: usecase, userUsecase: userUsecase, cfg: cfg}
}

// @Summary		Sign in to the application
// @Description	Sign in user with credentials
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			authenPayload	body		model.AuthenPayload	true	"Authentication Payload"
// @Success		200				{object}	model.AuthenResponse
// @Failure		400				{string}	string	"Missing body"
// @Failure		401				{string}	string	"Unauthorized"
// @Router			/auth/signin [post]
func (u *authenUsecase) SignIn(c *fiber.Ctx) error {
	payLoad := new(model.AuthenPayload)
	if err := c.BodyParser(payLoad); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Missing body")
	}
	authResponse, err := u.usecase.SignIn(payLoad)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	token, err := jwtSigner(authResponse.UserId, authResponse.Role, u.cfg.JWT_ACCESS_TOKEN)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 336),
		HTTPOnly: true})

	return c.JSON(model.AuthenResponse{
		Token: token,
	})
}

// @Summary		Google OAuth2 Callback
// @Description	Handle Google OAuth2 callback and log in or create a user
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Param			requestBody	body		model.RequestBody	true	"Google OAuth2 Request Body"
// @Success		200			{object}	model.AuthenResponse
// @Failure		204			{string}	string	"User Data Fetch Failed"
// @Failure		500			{string}	string	"Internal Server Error"
// @Router			/auth/google/callback [post]
func (u *authenUsecase) GoogleCallback(c *fiber.Ctx) error {

	requestBody := new(model.RequestBody)

	if err := c.BodyParser(requestBody); err != nil {
		return err
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + requestBody.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusNoContent).SendString("User Data Fetch Failed")
	}
	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("JSON Parsing Failed")
	}

	newGoogleUser := new(model.GoogleUser)
	err = json.Unmarshal(userData, newGoogleUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("JSON Unmarshal failed")
	}

	_, Geterr := u.userUsecase.FindUserByGoogleID(newGoogleUser.GoogleID)
	if Geterr != nil {
		if errors.Is(Geterr, gorm.ErrRecordNotFound) {
			// fmt.Println("Yedtood")
			newUser := new(model.Users)
			newUser.GoogleID = newGoogleUser.GoogleID
			newUser.Email = newGoogleUser.Email
			newUser.FirstName = newGoogleUser.Name
			newUser.LastName = newGoogleUser.Surname
			newUser.Role = "Client"
			newUser.CreateAt = time.Now()
			newUser.UpdateAt = time.Now()
			newUser.ProfileImageURL = newGoogleUser.ImageUrl
			newUser.Password = ""
			// fmt.Println(newUser)

			if err := u.userUsecase.CreateUser(*newUser); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
		}
	}

	test, testerr := u.userUsecase.FindUserByGoogleID(newGoogleUser.GoogleID)

	if testerr != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	token, err := jwtSigner(test.UserID, test.Role, u.cfg.JWT_ACCESS_TOKEN)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 336),
		HTTPOnly: true})

	return c.JSON(model.AuthenResponse{
		Token: token,
	})

}

// @Summary		Extract User Data from JWT token
// @Description	handle user data extraction and token expiration check
// @Tags			Authentication
// @Accept			json
// @Produce		json
// @Security		BearerAuth
// @Param			Authorization	header		string	true	"Bearer {access_token}"
// @Success		200				{object}	model.AuthenResponse
// @Failure		204				{string}	string	"User Data Fetch Failed"
// @Failure		401				{string}	string	"token expired"
// @Failure		500				{string}	string	"Internal Server Error"
// @Router			/auth/me [get]
func (u *authenUsecase) Me(c *fiber.Ctx) error {
	reqToken := c.Locals("user").(*jwt.Token)
	claims := reqToken.Claims.(jwt.MapClaims)
	response, err := u.usecase.Me(claims["userID"].(string))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func jwtSigner(userID string, role model.Roles, access_token string) (string, error) {
	day := time.Hour * 24
	claims := jwt.MapClaims{
		"userID":     userID,
		"positionID": role,
		"exp":        time.Now().Add(day * 14).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(access_token))
	if err != nil {
		return "", err
	}
	return t, nil
}
