package presentation

import (
	m "integration/process/model"
	repository "integration/process/repository"
)

type IHistoryUseCase interface {
	GetAll() ([]m.UserHistory, error)
	GetByUserName(name string) (m.UserHistory, error)
	CreateOne(*m.UserHistory) (m.UserHistory, error)
}

type HistoryUseCase struct {
	repo repository.IUserHistoryRepository
}

func NewHistoryUseCase(clientRepository repository.IUserHistoryRepository) IHistoryUseCase {
	return &HistoryUseCase{clientRepository}
}

func (c *HistoryUseCase) CreateOne(mod *m.UserHistory) (m.UserHistory, error) {
	return c.repo.CreateOne(mod)
}

func (c *HistoryUseCase) GetAll() ([]m.UserHistory, error) {
	return c.repo.GetAll()
}

func (c *HistoryUseCase) GetByUserName(name string) (m.UserHistory, error) {
	return c.repo.GetByUserName(name)
}
