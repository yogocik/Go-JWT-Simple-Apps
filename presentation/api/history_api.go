package presentation

import (
	response "integration/presentation/model/response"
	usecase "integration/presentation/usecase"

	"github.com/gin-gonic/gin"
)

type HistoryApi struct {
	usecase usecase.IHistoryUseCase
}

func NewHistoryApi(usecase usecase.IHistoryUseCase) IDelivery {
	HistoryApi := HistoryApi{
		usecase: usecase,
	}
	return &HistoryApi
}
func (api *HistoryApi) InitRouter(publicRoute *gin.RouterGroup) {
	userRoute := publicRoute.Group("/history")
	userRoute.GET("/name/:UserName", api.getByName)
}

func (api *HistoryApi) getByName(c *gin.Context) {
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
