package presentation

import r "integration/process/resource"

type UseCaseManager interface {
	UserUseCase() IUserUseCase
	HistoryUseCase() IHistoryUseCase
}

type useCaseManager struct {
	repo r.RepoManager
}

func (uc *useCaseManager) UserUseCase() IUserUseCase {
	return NewUserUseCase(uc.repo.UserRepo(), uc.repo.CacheRepo(), uc.repo.HistoryRepo())
}

func (uc *useCaseManager) HistoryUseCase() IHistoryUseCase {
	return NewHistoryUseCase(uc.repo.HistoryRepo())
}

func NewUseCaseManger(manager r.RepoManager) UseCaseManager {
	return &useCaseManager{repo: manager}
}
