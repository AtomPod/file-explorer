package repository

//DataRepository 数据仓库集合
type DataRepository interface {
	File() (FileRepository, error)
	User() (UserRepository, error)
	VerificationCode() (VerificationCodeRepository, error)
}

//UnitOfWork 单元工作
type UnitOfWork interface {
	DataRepository
	Commit() error
	Rollback() error
}

//DataContext 数据上下文
type DataContext interface {
	DataRepository
	Unit() (UnitOfWork, error)
}
