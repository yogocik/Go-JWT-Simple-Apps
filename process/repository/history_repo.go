package process

import (
	u "integration/utility/reusable"

	m "integration/process/model"

	"gorm.io/gorm"
)

type IUserHistoryRepository interface {
	GetAll() ([]m.UserHistory, error)
	GetByUserName(name string) (m.UserHistory, error)
	CreateOne(*m.UserHistory) (m.UserHistory, error)
	Migrate(*m.UserHistory)
}

type UserHistoryRepository struct {
	db *gorm.DB
}

func NewUserHistoryRepository(resource *gorm.DB) IUserHistoryRepository {
	historyRepository := &UserHistoryRepository{db: resource}
	return historyRepository
}

func (tr *UserHistoryRepository) Migrate(table *m.UserHistory) {
	tr.db.AutoMigrate(table)
}

func (t *UserHistoryRepository) GetAll() ([]m.UserHistory, error) {
	userList := []m.UserHistory{}
	result := t.db.Find(&userList)
	return userList, u.HandleGormError(result)
}

func (t *UserHistoryRepository) GetByUserName(name string) (m.UserHistory, error) {
	var userList m.UserHistory
	result := t.db.Find(&userList).Where("user_name = ?", name)
	return userList, u.HandleGormError(result)
}

func (tr *UserHistoryRepository) CreateOne(hist *m.UserHistory) (m.UserHistory, error) {
	result := tr.db.Create(hist)
	return *hist, u.HandleGormError(result)
}
