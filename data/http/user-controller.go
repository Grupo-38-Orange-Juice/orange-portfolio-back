package http

import (
	"fmt"
	"net/http"

	"github.com/Grupo-38-Orange-Juice/orange-portfolio-back/domain/entities"
	usecases "github.com/Grupo-38-Orange-Juice/orange-portfolio-back/domain/use-cases"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	userUseCase usecases.UserUseCase
	crypto      usecases.Crypto
	token       usecases.Token
}

func NewUserController(userUseCase usecases.UserUseCase, crypto usecases.Crypto, token usecases.Token) UserController {
	return UserController{
		userUseCase: userUseCase,
		crypto:      crypto,
		token:       token,
	}
}

// @Summary Create a new user
// @Description Create a new user with the provided information
// @Accept json
// @Produce json
// @Param user body CreateUserDTO true "User object to be created"
// @Success 201 {object} Response "User created successfully" {"message": "User created successfully"}
// @Failure 400 {object} Response "Bad Request" {"message": "Bad Request"}
// @Failure 500 {object} Response "Internal Server Error" {"message": "Internal Server Error"}
// @Router /users [post]
func (u UserController) CreateUser(c *gin.Context) {
	var user entities.User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error"})
		return
	}

	newUser, err := entities.NewUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	newUser.Password, err = u.crypto.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error"})
		return
	}

	newUser.ID = uuid.New().String()
	err = u.userUseCase.CreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuário criado com sucesso"})
}

// @Summary Login
// @Description Login with the provided information and get a api token
// @Param user body LoginDTO true "User object to be logged in"
// @Router /login [post]
// @Success 201 {object} Response "Token" {"token": "user_token"}
// @Failure 401 {object} Response "Unauthorized" {"message": "Unauthorized"}
// @Failure 500 {object} Response "Internal Server Error" {"message": "Internal Server Error"}
func (u UserController) Login(c *gin.Context) {
	loginRequest := LoginDTO{}
	err := c.ShouldBind(&loginRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error"})
		return
	}
	fmt.Println(loginRequest)
	user, err := u.userUseCase.FindUserByEmail(loginRequest.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Email não encontrado"})
		return
	}

	if err = u.crypto.CompareHashAndPassword(user.Password, loginRequest.Password); err != nil {
		println(err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Senha inválida"})
		return
	}
	token, err := u.token.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Get user information
// @Description Get user information by api token
// @Security ApiKeyAuth
// @Router /me [get]
// @Success 200 {object} UserDTO "User information"
// @Failure 401 {object} Response "Unauthorized" {"message": "Unauthorized"}
// @Failure 404 {object} Response "Not Found" {"message": "Usuário não encontrado"}
// @Failure 500 {object} Response "Internal Server Error" {"message": "Internal Server Error"}
func (u UserController) Me(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	userId, err := u.token.GetPayload(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	user, err := u.userUseCase.FindUserById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Usuário não encontrado"})
		return
	}
	c.JSON(http.StatusOK, userEntityToDTO(*user))
}
