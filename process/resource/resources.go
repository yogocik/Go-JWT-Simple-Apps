package process

import (
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type InfraSQL interface {
	NewClientSQL(connectionString string)
	SqlDb() *gorm.DB
}

type InfraKeyValue interface {
	KVStorage() *redis.Client
	NewClientKVS(*redis.Options)
}

type Infra struct {
	db    *gorm.DB
	redis *redis.Client
}

func (i *Infra) SqlDb() *gorm.DB {
	return i.db
}

func (i *Infra) NewClientSQL(connectionString string) {
	conn, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed connect to database", err)
	}
	i.db = conn
}

func (i *Infra) KVStorage() *redis.Client {
	return i.redis
}

func (i *Infra) NewClientKVS(opt *redis.Options) {
	client := redis.NewClient(opt)
	i.redis = client
}
