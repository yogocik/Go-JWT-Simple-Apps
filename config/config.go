package config

import (
	"fmt"
	delivery "integration/presentation/api"
	usecase "integration/presentation/usecase"
	authenticator "integration/process/authentication"
	model "integration/process/model"
	manager "integration/process/resource"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

type Config struct {
	InfraManager   *manager.Infra
	RepoManager    manager.RepoManager
	UseCaseManager usecase.UseCaseManager
	Routes         *delivery.Routes
	ApiBaseUrl     string
	TokenConfig    authenticator.TokenConfig
	TokenServices  authenticator.Token
}

func NewConfig() *Config {
	apiHost := os.Getenv("API_HOST")
	apiPort := os.Getenv("API_PORT")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")

	redisPORT := os.Getenv("REDIS_PORT")
	redisHOST := os.Getenv("REDIS_HOST")
	redisPASS := os.Getenv("REDIS_PASSWORD")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	redisTIMEOUT, _ := strconv.Atoi(os.Getenv("REDIS_TIMEOUT"))

	kvs_opt := &redis.Options{
		Addr:     redisHOST + ":" + redisPORT,
		Password: redisPASS,
		DB:       redisDB,
	}

	MyInfra := &manager.Infra{}
	MyInfra.NewClientSQL(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName))
	MyInfra.NewClientKVS(kvs_opt)
	fmt.Println("My REDIS ->", MyInfra.KVStorage())
	repoManager := manager.NewRepoManager(MyInfra, MyInfra)
	useCaseManager := usecase.NewUseCaseManger(repoManager)
	repoManager.UserRepo().Migrate(&model.User{})
	repoManager.HistoryRepo().Migrate(&model.UserHistory{})
	config := new(Config)
	config.InfraManager = MyInfra
	config.RepoManager = repoManager
	config.UseCaseManager = useCaseManager
	tokenConfig := authenticator.TokenConfig{
		ApplicationName:     "ENIGMA",
		JwtSignatureKey:     "P@ssw0rd",
		JwtSigningMethod:    jwt.SigningMethodHS256,
		AccessTokenLifeTime: time.Duration(redisTIMEOUT) * time.Second,
		Client:              config.RepoManager.CacheRepo(),
	}
	tokenService := authenticator.NewTokenService(tokenConfig)
	config.TokenConfig = tokenConfig
	config.TokenServices = tokenService
	fmt.Println("MY TOKEN CONFIG -> ", tokenConfig.Client)
	router := delivery.NewServer(useCaseManager, config.TokenServices)
	config.Routes = router

	config.ApiBaseUrl = fmt.Sprintf("%s:%s", apiHost, apiPort)
	return config
}
