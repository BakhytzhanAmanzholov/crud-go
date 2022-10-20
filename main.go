package main

import (
	"crud-golang/internal/dto"
	"crud-golang/internal/repositories/mongo"
	services "crud-golang/internal/services"
	"crud-golang/internal/services/database"
	"crud-golang/pkg/client/mongodb"
	"crud-golang/pkg/logging"
	"crud-golang/pkg/responses"
	swagger "github.com/arsmn/fiber-swagger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	repository = mongo.NewRepository()
	service    = database.NewService(repository)
	validate   = validator.New()
)

// @title CRUD account
// @version 1.0
// @description API Server for CRUD Application

// @host localhost:8181
// @BasePath /accounts

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logging.NewLogger()
	_, err := mongodb.ConnectDB()
	if err != nil {
		return
	}
	app := fiber.New()

	app.Post("/login", signIn)
	accountApi := app.Group("/accounts")
	accountApi.Get("/", findAll)
	accountApi.Post("/", createAccount)
	accountApi.Get("/:accountId", findOne)
	accountApi.Put("/:accountId", updateAccount)
	accountApi.Delete("/:accountId", deleteAccount)

	//app.Use("/hello", func(ctx *fiber.Ctx) {
	//	jwtware.New(jwtware.Config{
	//		SigningKey: []byte(os.Getenv("SECRET")),
	//	})
	//
	//	ctx.Send("Hello, world!")
	//})
	app.Get("/private", private)

	registerSwagger(app)

	app.Listen(":8181")
}

func private(ctx *fiber.Ctx) {
	authorization := ctx.Fasthttp.Request.Header.Peek("Authorization")
	_, str, _ := strings.Cut(string(authorization), "Bearer ")

	prime, err := services.VerifyJWT(str)
	if err != nil {
		ctx.Status(fiber.StatusUnauthorized)
	}
	if prime {
		ctx.JSON(fiber.Map{"Data": "Hello world!"})
	}
}

func registerSwagger(app *fiber.App) {
	enableSwagger := os.Getenv("ENABLE_SWAGGER")
	if enabled, _ := strconv.ParseBool(enableSwagger); enabled {
		route := app.Group("/swagger")
		route.Get("*", swagger.Handler)
		logging.Infof("Swagger Started")
	}
}

// @Summary SignIn
// @Tags auth
// @Description create account
// @ID signin
// @Accept  json
// @Produce  json
// @Param input body dto.Login true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /auth/sign-up [post]
func signIn(ctx *fiber.Ctx) {
	var accountDto dto.Login

	if err := ctx.BodyParser(&accountDto); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(responses.Account{Status: http.StatusBadRequest,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}

	if validationErr := validate.Struct(&accountDto); validationErr != nil {
		ctx.Status(http.StatusBadRequest).JSON(responses.Account{Status: http.StatusBadRequest,
			Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
		return
	}

	account, err := service.Login(accountDto.Email, accountDto.Password)
	if err != nil {
		ctx.Status(http.StatusUnauthorized).JSON(responses.Account{Status: http.StatusInternalServerError,
			Message: "error", Data: &fiber.Map{"data": err}})
		return
	}
	ctx.Status(http.StatusOK).JSON(responses.Account{Status: http.StatusCreated,
		Message: "success", Data: &fiber.Map{"data": account}})

	token, exp, err := services.CreateJWTToken(account.Id)
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(responses.Account{Status: http.StatusInternalServerError,
			Message: "error", Data: &fiber.Map{"data": err}})
		return
	}

	ctx.JSON(fiber.Map{"token": token, "exp": exp, "user": account})
}

func findAll(ctx *fiber.Ctx) {
	accounts, err := service.FindAll()
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(responses.Account{Status: http.StatusInternalServerError,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}
	ctx.Status(http.StatusOK).JSON(
		responses.Account{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": accounts}},
	)
	return
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body dto.Registration true "Registration info"
// @Success 200 {object} models.Account
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /accounts/ [post]
func createAccount(ctx *fiber.Ctx) {
	var accountDto dto.Registration

	if err := ctx.BodyParser(&accountDto); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(responses.Account{Status: http.StatusBadRequest,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}

	if validationErr := validate.Struct(&accountDto); validationErr != nil {
		ctx.Status(http.StatusBadRequest).JSON(responses.Account{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
		return
	}

	account, err := service.Create(accountDto)
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(responses.Account{Status: http.StatusInternalServerError,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}
	ctx.Status(http.StatusCreated).JSON(responses.Account{Status: http.StatusCreated,
		Message: "success", Data: &fiber.Map{"data": account}})
}

// @Summary Find by ID
// @Tags auth
// @Description Find account by id
// @ID find-account-id
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Account
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /accounts/:accountId [get]
func findOne(ctx *fiber.Ctx) {
	accountId := ctx.Params("accountId")
	account, err := service.FindOne(accountId)

	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(responses.Account{Status: http.StatusInternalServerError,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}

	ctx.Status(http.StatusOK).JSON(responses.Account{Status: http.StatusOK,
		Message: "success", Data: &fiber.Map{"data": account}})
}

// @Summary Update account
// @Tags auth
// @Description Update account by id
// @ID update-account
// @Accept  json
// @Produce  json
// @Param input body dto.Registration true "account info"
// @Success 200 {object} models.Account
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /accounts/:accountId [put]
func updateAccount(ctx *fiber.Ctx) {
	accountId := ctx.Params("accountId")
	var accountDto dto.Registration

	if err := ctx.BodyParser(&accountDto); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(responses.Account{Status: http.StatusBadRequest,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}

	if validationErr := validate.Struct(&accountDto); validationErr != nil {
		ctx.Status(http.StatusBadRequest).JSON(responses.Account{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
		return
	}
	account, err := service.Update(accountDto, accountId)

	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(responses.Account{Status: http.StatusInternalServerError,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}

	ctx.Status(http.StatusOK).JSON(responses.Account{Status: http.StatusOK,
		Message: "success", Data: &fiber.Map{"data": account}})
}

// @Summary Delete account
// @Tags auth
// @Description delete account by id
// @ID delete-account
// @Accept  json
// @Produce  json
// @Failure 400,404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Failure default {object} HTTPError
// @Router /accounts/:accountId [put]
func deleteAccount(ctx *fiber.Ctx) {
	userId := ctx.Params("userId")

	err := service.Delete(userId)

	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(responses.Account{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}

	ctx.Status(http.StatusOK).JSON(responses.Account{Status: http.StatusOK,
		Message: "success", Data: &fiber.Map{"data": "Successfully deleted"}})

}

type HTTPError struct {
	status  string
	message string
}
