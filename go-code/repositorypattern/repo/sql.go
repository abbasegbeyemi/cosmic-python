package repositorypattern

import (
	"database/sql"
	"log"
)

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db}
}

func (r *SQLRepository) Add(batch *Batch) error {
	// Add the batch to the database
	result, err := r.db.Exec(
		"insert into batches (reference, sku, quantity, eta) values (?, ?, ?, ?)",
		batch.Reference, batch.SKU, batch.Quantity, batch.ETA,
	)
	log.Println(result)
	return err
}
