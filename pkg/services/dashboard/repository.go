package dashboard

import (
	"enterprise.sidooh/pkg/datastore"
)

type Repository interface {
	CountAccounts() int64
}

type repository struct {
}

func (r *repository) CountAccounts() (count int64) {
	datastore.DB.Table("accounts").Count(&count)
	return
}

func NewRepo() Repository {
	return &repository{}
}
