package repository

import (
	"database/sql"
	"errors"

	uuid "github.com/satori/go.uuid"

	storage "github.com/overridesh/checkout-challenger/pkg/storage/sql"
)

var (
	ErrMerchantNotFound = errors.New("invalid key")
)

const (
	getByApiKey string = `
	 	SELECT id
	 	FROM merchants
	 	WHERE apikey = $1
		AND active = true
	`
)

type MerchantRepository interface {
	GetIDByApiKey(apikey uuid.UUID) (uuid.UUID, error)
}

type merchantRepository struct {
	db storage.DB
}

func NewMerchantRepository(db storage.DB) MerchantRepository {
	return &merchantRepository{
		db: db,
	}
}

func (pg *merchantRepository) GetIDByApiKey(apikey uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID

	if err := pg.db.QueryRow(
		getByApiKey,
		apikey,
	).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, ErrMerchantNotFound
		}
		return uuid.Nil, err
	}

	return id, nil
}
