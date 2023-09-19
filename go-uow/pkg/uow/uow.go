package uow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type RepositoryFactory func(tx *sql.Tx) interface{}

type UowInterface interface {
	Register(name string, fc RepositoryFactory)
	GetRepository(ctx context.Context, name string) (interface{}, error)
	Do(ctx context.Context, fn func(uow *Uow) error) error
	CommitOrRollback() error
	Rollback() error
	UnRegister(name string)
}

type Uow struct {
	Db           *sql.DB
	Tx           *sql.Tx
	Repositories map[string]RepositoryFactory
}

func NewUow(ctx context.Context, db *sql.DB) *Uow {
	return &Uow{
		Db:           db,
		Repositories: make(map[string]RepositoryFactory),
	}
}

func (u *Uow) Register(name string, fc RepositoryFactory) {
	u.Repositories[name] = fc
}

func (u *Uow) UnRegister(name string) {
	delete(u.Repositories, name)
}

func (u *Uow) GetRepository(ctx context.Context, name string) (interface{}, error) {
	if u.Tx == nil {
		tx, err := u.Db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
		u.Tx = tx
	}
	repo := u.Repositories[name](u.Tx)
	return repo, nil
}

func (u *Uow) Do(ctx context.Context, fn func(uow *Uow) error) error {
	if u.Tx != nil {
		return errors.New("transaction already in progress")
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	u.Tx = tx

	if err = fn(u); err != nil {
		if errRb := u.Rollback(); errRb != nil {
			return fmt.Errorf("original error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
	}
	return u.CommitOrRollback()
}

func (u *Uow) Rollback() error {
	if u.Tx != nil {
		return errors.New("no transaction to rollback")
	}
	if err := u.Tx.Rollback(); err != nil {
		return err
	}
	u.Tx = nil
	return nil
}

func (u *Uow) CommitOrRollback() error {
	if err := u.Tx.Commit(); err != nil {
		if errRb := u.Rollback(); errRb != nil {
			return fmt.Errorf("original error: %s, rollback error: %s", err.Error(), errRb.Error())
		}
		return err
	}
	u.Tx = nil
	return nil
}
