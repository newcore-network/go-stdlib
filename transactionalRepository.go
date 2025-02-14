package stdlib

import "gorm.io/gorm"

// TransactionalRepository defines the methods for managing transactions in a database.

type TransactionalRepository interface {

	// BeginTransaction initializes a transaction by returning the transaction context or a possible error
	BeginTransaction() (*gorm.DB, error)

	// CommitTransaction commits a transaction, returning a possible error
	CommitTransaction(tx *gorm.DB) error
	// RollbackTransaction rolls back a transaction leaving changes uncommitted
	RollbackTransaction(tx *gorm.DB) error

	// ExecuteInTransaction executes a function within a transaction context, logging whether there is a possible error
	ExecuteInTransaction(fn func(tx *gorm.DB) error) error
}

type transactionalRepositoryImpl struct {
	gorm *gorm.DB
}

func NewTransactionalRepository(gorm *gorm.DB) TransactionalRepository {
	return &transactionalRepositoryImpl{gorm: gorm}
}

func (repo *transactionalRepositoryImpl) BeginTransaction() (*gorm.DB, error) {
	tx := repo.gorm.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (repo *transactionalRepositoryImpl) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (repo *transactionalRepositoryImpl) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (repo *transactionalRepositoryImpl) ExecuteInTransaction(fn func(tx *gorm.DB) error) error {
	tx, err := repo.BeginTransaction()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil || err != nil {
			repo.RollbackTransaction(tx)
		}
	}()

	err = fn(tx)
	if err != nil {
		repo.RollbackTransaction(tx)
		return err
	}

	return repo.CommitTransaction(tx)
}
