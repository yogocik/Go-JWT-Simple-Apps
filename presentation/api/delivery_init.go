package presentation

import (
	usecase "integration/presentation/usecase"
	authenticator "integration/process/authentication"

	"github.com/gin-gonic/gin"
)

type IDelivery interface {
	InitRouter(publicRoute *gin.RouterGroup)
}

type Routes struct {
	routers       []IDelivery
	RouterEngine  *gin.Engine
	publicRoute   *gin.RouterGroup
	tokenServices authenticator.Token
}

func NewServer(useCaseManager usecase.UseCaseManager, tokenServices authenticator.Token) *Routes {
	newServer := new(Routes)

	r := gin.Default()
	publicRoute := r.Group("/api")
	routers := []IDelivery{
		NewUserApi(useCaseManager.UserUseCase(), tokenServices),
		NewHistoryApi(useCaseManager.HistoryUseCase()),
	}
	newServer.routers = routers
	newServer.RouterEngine = r
	newServer.publicRoute = publicRoute
	newServer.initAppRoutes()
	return newServer
}
func (app *Routes) initAppRoutes() {
	for _, rt := range app.routers {
		rt.InitRouter(app.publicRoute)
	}
}

func (app *Routes) AttachToken(token authenticator.Token) {
	app.tokenServices = token
}
