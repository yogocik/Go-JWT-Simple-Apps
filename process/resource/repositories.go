package process

import r "integration/process/repository"

type RepoManager interface {
	UserRepo() r.IUserRepository
	HistoryRepo() r.IUserHistoryRepository
	CacheRepo() r.ICacheRepository
}

type repoManager struct {
	infraSQL InfraSQL
	infrakv  InfraKeyValue
}

func (rm *repoManager) UserRepo() r.IUserRepository {
	return r.NewUserRepository(rm.infraSQL.SqlDb())
}

func (rm *repoManager) HistoryRepo() r.IUserHistoryRepository {
	return r.NewUserHistoryRepository(rm.infraSQL.SqlDb())
}

func (rm *repoManager) CacheRepo() r.ICacheRepository {
	return r.NewCacheRepository(rm.infrakv.KVStorage())
}

func NewRepoManager(infra InfraSQL, infrakv InfraKeyValue) RepoManager {
	return &repoManager{infra, infrakv}
}
