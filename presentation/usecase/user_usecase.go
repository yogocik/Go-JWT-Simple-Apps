package presentation

import (
	m "integration/process/model"
	repository "integration/process/repository"
)

type IUserUseCase interface {
	GetAll() ([]m.User, error)
	GetByUserName(name string) (m.User, error)
	CreateOne(*m.User) (m.User, error)
	CreateRecord(*m.UserHistory) (m.UserHistory, error)
}

type UserUseCase struct {
	repo    repository.IUserRepository
	cache   repository.ICacheRepository
	history repository.IUserHistoryRepository
}

func NewUserUseCase(clientRepository repository.IUserRepository,
	cache repository.ICacheRepository,
	history repository.IUserHistoryRepository) IUserUseCase {
	return &UserUseCase{clientRepository, cache, history}
}

func (c *UserUseCase) CreateOne(mod *m.User) (m.User, error) {
	return c.repo.CreateOne(mod)
}

func (c *UserUseCase) GetAll() ([]m.User, error) {
	return c.repo.GetAll()
}

func (c *UserUseCase) CreateRecord(m *m.UserHistory) (m.UserHistory, error) {
	return c.history.CreateOne(m)
}

func (c *UserUseCase) GetByUserName(name string) (m.User, error) {
	return c.repo.GetByUserName(name)
}
