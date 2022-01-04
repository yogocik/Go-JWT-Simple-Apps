package presentation

import (
	"errors"
	"fmt"
	request "integration/presentation/model/request"
	response "integration/presentation/model/response"
	usecase "integration/presentation/usecase"
	auth "integration/process/authentication"
	auth_model "integration/process/authentication/model"
	model "integration/process/model"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserApi struct {
	usecase  usecase.IUserUseCase
	tokenJWT auth.Token
}

func NewUserApi(usecase usecase.IUserUseCase, token auth.Token) IDelivery {
	UserApi := UserApi{
		usecase:  usecase,
		tokenJWT: token,
	}
	return &UserApi
}
func (api *UserApi) InitRouter(publicRoute *gin.RouterGroup) {
	userRoute := publicRoute.Group("/user")
	userRoute.GET("/all", api.getAllUser)
	userRoute.GET("/name/:UserName", api.getByName)
	userRoute.POST("/login", api.login)
	userRoute.POST("", api.createUser)
	userRoute.POST("/logout", api.logout)
}

func (api *UserApi) createUser(c *gin.Context) {
	var user request.UserForm
	err := c.BindJSON(&user)
	if err != nil {
		errModel := response.NewBadRequestError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Bad Request",
			nil, errModel))
		return
	}
	userModel := model.NewUserDB(user.Username, user.Password, user.Email)
	newUser, err := api.usecase.CreateOne(userModel)
	if err != nil {
		errModel := response.NewInternalServerError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Internal Server Error",
			nil, errModel))
		return
	}
	successModel := response.NewResponse("Success", newUser, response.NewError{})
	c.JSON(successModel["StatusCode"].(int), successModel)
}

func (api *UserApi) getAllUser(c *gin.Context) {
	allUser, err := api.usecase.GetAll()
	if err != nil {
		errModel := response.NewInternalServerError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Internal Server Error",
			nil, errModel))
		return
	}
	successModel := response.NewResponse("Success", allUser, response.NewError{})
	c.JSON(successModel["StatusCode"].(int), successModel)
}

func (api *UserApi) getByName(c *gin.Context) {
	name := c.Param("UserName")
	allUser, err := api.usecase.GetByUserName(name)
	if err != nil {
		errModel := response.NewInternalServerError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Internal Server Error",
			nil, errModel))
		return
	}
	successModel := response.NewResponse("Success", allUser, response.NewError{})
	c.JSON(successModel["StatusCode"].(int), successModel)
}

func (api *UserApi) login(c *gin.Context) {
	var user request.UserLogin
	err := c.BindJSON(&user)
	if err != nil {
		errModel := response.NewBadRequestError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Bad Request",
			nil, errModel))
		return
	}
	searchData := auth.AccessDetails{AccessUuid: user.Username}
	data, err := api.tokenJWT.FetchAccessToken(&searchData)
	fmt.Println("MY LOGIN REDIS ->", data, err)
	if data != "" {
		fmt.Println("ERROR KEISINI")
		errModel := response.NewUnauthorizedError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Concurrent Login",
			data, errModel))
		return
	}
	allUser, err := api.usecase.GetByUserName(user.Username)
	fmt.Println("GET BY NAME", allUser)
	if err != nil {
		errModel := response.NewInternalServerError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Internal Server Error",
			nil, errModel))
		return
	}
	if allUser.Username == "" {
		errModel := response.NewBadRequestError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Username is not found or wrong",
			nil, errModel))
		return
	}
	fmt.Println("User and Password from DB ->", allUser.Username, allUser.Password, "User Pasword Form ->", user.Username, user.Password)
	if allUser.Password != user.Password {
		errModel := response.NewUnauthorizedError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Invalid password",
			nil, errModel))
		return
	}
	credModel := auth_model.CredentialModel{
		Username: allUser.Username,
		Password: allUser.Password,
	}
	token, err := api.tokenJWT.CreateAccessToken(&credModel)
	if err != nil {
		errModel := response.NewUnauthorizedError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Something wrong in token creation",
			nil, errModel))
		return
	}
	err = api.tokenJWT.StoreAccessToken(user.Username, token)
	if err != nil {
		errModel := response.NewInternalServerError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Something wrong in token caching process",
			nil, errModel))
		return
	}
	record := model.NewUserHistoryDB(allUser.Username, true, true)
	rec, err := api.usecase.CreateRecord(record)
	if err != nil {
		errModel := response.NewInternalServerError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Failed login record",
			nil, errModel))
		return
	}
	successModel := response.NewResponse("Success", gin.H{"token": token, "record": rec}, response.NewError{})
	c.JSON(successModel["StatusCode"].(int), successModel)
}

func (api *UserApi) logout(c *gin.Context) {
	header := struct {
		AuthorizationHeader string `header:"Authorization"`
	}{}
	if err := c.ShouldBindHeader(&header); err != nil {
		errModel := response.NewUnauthorizedError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("No valid token provided",
			nil, errModel))
		return
	}
	fmt.Println("My Token Header", header.AuthorizationHeader)
	tokenString := strings.Replace(header.AuthorizationHeader, "Bearer ", "", -1)
	fmt.Println("My Token String", tokenString)
	if tokenString == "" {
		errModel := response.NewUnauthorizedError(errors.New("Empty Bearer Token"))
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("No valid token provided",
			nil, errModel))
		return
	}
	token, err := api.tokenJWT.VerifyAccessToken(tokenString)
	if err != nil {
		errModel := response.NewUnauthorizedError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Invalid token",
			nil, errModel))
		return
	}
	numCode, err := api.tokenJWT.RevokeAccessToken(token)
	if numCode == -1 || err != nil {
		errModel := response.NewInternalServerError(err)
		c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Something Wrong in token revocation",
			nil, errModel))
		return
	}
	successModel := response.NewResponse("Success revocation", gin.H{"token": token}, response.NewError{})
	c.JSON(successModel["StatusCode"].(int), successModel)
}
