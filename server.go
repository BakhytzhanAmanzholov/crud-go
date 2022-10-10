package main

import (
	"crud-golang/internal/dto"
	"crud-golang/internal/repositories"
	services "crud-golang/internal/services"
	"crud-golang/pkg/client/mongodb"
	"crud-golang/pkg/logging"
	"crud-golang/pkg/responses"
	swagger "github.com/arsmn/fiber-swagger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber"
	jwtware "github.com/gofiber/jwt/v3"
	"net/http"
	"os"
	"strconv"
)

var (
	repository = repositories.NewRepository()
	service    = services.NewService(repository)
	validate   = validator.New()
)

const jwtSecret = "secret"

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
	mongodb.ConnectDB()
	app := fiber.New()

	app.Post("/login", signIn)
	accountApi := app.Group("/accounts")
	accountApi.Get("/", findAll)
	accountApi.Post("/", createAccount)
	accountApi.Get("/:accountId", findOne)
	accountApi.Put("/:accountId", updateAccount)
	accountApi.Delete("/:accountId", deleteAccount)

	app.Get("/hello", func(ctx *fiber.Ctx) {
		jwtware.New(jwtware.Config{
			SigningKey: []byte("secret"),
		})
		ctx.Send("Hello, world!")
	})

	registerSwagger(app)

	app.Listen(":8181")
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
// @Param input body todo.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func signIn(ctx *fiber.Ctx) {
	var accountDto dto.LoginDto

	if err := ctx.BodyParser(&accountDto); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(responses.AccountResponse{Status: http.StatusBadRequest,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}

	if validationErr := validate.Struct(&accountDto); validationErr != nil {
		ctx.Status(http.StatusBadRequest).JSON(responses.AccountResponse{Status: http.StatusBadRequest,
			Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
		return
	}

	account, err := service.Login(accountDto.Email, accountDto.Password)
	if err != nil {
		ctx.Status(http.StatusUnauthorized).JSON(responses.AccountResponse{Status: http.StatusInternalServerError,
			Message: "error", Data: &fiber.Map{"data": err}})
		return
	}
	ctx.Status(http.StatusOK).JSON(responses.AccountResponse{Status: http.StatusCreated,
		Message: "success", Data: &fiber.Map{"data": account}})

	token, exp, err := services.CreateJWTToken(account)
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(responses.AccountResponse{Status: http.StatusInternalServerError,
			Message: "error", Data: &fiber.Map{"data": err}})
		return
	}

	ctx.JSON(fiber.Map{"token": token, "exp": exp, "user": account})
}

func findAll(ctx *fiber.Ctx) {
	accounts, err := service.FindAll()
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(responses.AccountResponse{Status: http.StatusInternalServerError,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}
	ctx.Status(http.StatusOK).JSON(
		responses.AccountResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": accounts}},
	)
	return
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body dto.AccountDto true "account info"
// @Success 200 {object} models.Account
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /accounts/ [post]
func createAccount(ctx *fiber.Ctx) {
	var accountDto dto.AccountDto

	if err := ctx.BodyParser(&accountDto); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(responses.AccountResponse{Status: http.StatusBadRequest,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}

	if validationErr := validate.Struct(&accountDto); validationErr != nil {
		ctx.Status(http.StatusBadRequest).JSON(responses.AccountResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
		return
	}

	account, err := service.Create(accountDto)
	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(responses.AccountResponse{Status: http.StatusInternalServerError,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}
	ctx.Status(http.StatusCreated).JSON(responses.AccountResponse{Status: http.StatusCreated,
		Message: "success", Data: &fiber.Map{"data": account}})
}

// @Summary Find by ID
// @Tags auth
// @Description Find account by id
// @ID find-account-id
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Account
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /accounts/:accountId [get]
func findOne(ctx *fiber.Ctx) {
	accountId := ctx.Params("accountId")
	account, err := service.FindOne(accountId)

	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(responses.AccountResponse{Status: http.StatusInternalServerError,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}

	ctx.Status(http.StatusOK).JSON(responses.AccountResponse{Status: http.StatusOK,
		Message: "success", Data: &fiber.Map{"data": account}})
}

// @Summary Update account
// @Tags auth
// @Description Update account by id
// @ID update-account
// @Accept  json
// @Produce  json
// @Param input body dto.AccountDto true "account info"
// @Success 200 {object} models.Account
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /accounts/:accountId [put]
func updateAccount(ctx *fiber.Ctx) {
	accountId := ctx.Params("accountId")
	var accountDto dto.AccountDto

	if err := ctx.BodyParser(&accountDto); err != nil {
		ctx.Status(http.StatusBadRequest).JSON(responses.AccountResponse{Status: http.StatusBadRequest,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}

	if validationErr := validate.Struct(&accountDto); validationErr != nil {
		ctx.Status(http.StatusBadRequest).JSON(responses.AccountResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
		return
	}
	account, err := service.Update(accountDto, accountId)

	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(responses.AccountResponse{Status: http.StatusInternalServerError,
			Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}

	ctx.Status(http.StatusOK).JSON(responses.AccountResponse{Status: http.StatusOK,
		Message: "success", Data: &fiber.Map{"data": account}})
}

// @Summary Delete account
// @Tags auth
// @Description delete account by id
// @ID delete-account
// @Accept  json
// @Produce  json
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /accounts/:accountId [put]
func deleteAccount(ctx *fiber.Ctx) {
	userId := ctx.Params("userId")

	err := service.Delete(userId)

	if err != nil {
		ctx.Status(http.StatusInternalServerError).JSON(responses.AccountResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		return
	}

	ctx.Status(http.StatusOK).JSON(responses.AccountResponse{Status: http.StatusOK,
		Message: "success", Data: &fiber.Map{"data": "Successfully deleted"}})

}
