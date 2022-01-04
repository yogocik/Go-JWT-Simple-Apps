package process

import (
	u "integration/utility/reusable"

	m "integration/process/model"

	"gorm.io/gorm"
)

type IUserRepository interface {
	GetAll() ([]m.User, error)
	GetByUserName(name string) (m.User, error)
	CreateOne(*m.User) (m.User, error)
	Migrate(*m.User)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(resource *gorm.DB) IUserRepository {
	memberRepository := &UserRepository{db: resource}
	return memberRepository
}

func (tr *UserRepository) Migrate(table *m.User) {
	tr.db.AutoMigrate(table)
}

func (t *UserRepository) GetAll() ([]m.User, error) {
	userList := []m.User{}
	result := t.db.Find(&userList)
	return userList, u.HandleGormError(result)
}

func (t *UserRepository) GetByUserName(name string) (m.User, error) {
	var userList m.User
	result := t.db.Where("username = ?", name).Find(&userList)
	return userList, u.HandleGormError(result)
}

func (tr *UserRepository) CreateOne(user *m.User) (m.User, error) {
	result := tr.db.Create(user)
	return *user, u.HandleGormError(result)
}
