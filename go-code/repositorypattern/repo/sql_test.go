package repositorypattern

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getSqlDb(t testing.TB) *sql.DB {
	t.Helper()
	// Create a temporary file for the database
	dir, _ := os.Getwd()
	file, err := os.CreateTemp(dir, "test-*.db")
	assert.Nil(t, err)

	// Open a connection to the database
	db, err := sql.Open("sqlite3", file.Name())
	assert.Nil(t, err)

	t.Cleanup(func() {
		// Delete the file when the test is complete
		if err := os.Remove(file.Name()); err != nil {
			t.Fatal(err)
		}
	})

	// Create test tables
	sqlStmt := `
	create table if not exists batches (id integer not null primary key, reference text, sku text, quantity integer, eta text);
	create table if not exists order_lines (id integer not null primary key, order_id text, sku text, quantity integer, batch_id integer);
	`

	_, err = db.Exec(sqlStmt)
	assert.Nil(t, err)

	return db
}
func TestSQLRepository(t *testing.T) {

	t.Run("Can save a batch", func(t *testing.T) {
		// Test database queries should be in a transaction
		db := getSqlDb(t)

		// Create a new repository
		repo := NewSQLRepository(db)

		// Create a batch
		batch := Batch{
			Reference: "batch1",
			SKU:       "RUSTY_SNAPPISH",
			Quantity:  100,
			ETA:       "",
		}

		err := repo.Add(&batch)
		assert.Nil(t, err)

		var createdBatch Batch

		// Check that the batch was saved
		err = db.QueryRow("select reference, sku, quantity, eta from batches where reference = ?", "batch1").Scan(
			&createdBatch.Reference,
			&createdBatch.SKU,
			&createdBatch.Quantity,
			&createdBatch.ETA,
		)

		assert.Nil(t, err)
		assert.Equal(t, batch, createdBatch)
	})

	t.Run("Can get a batch by reference", func(t *testing.T) {
		db := getSqlDb(t)

		// Create a batch
		batch := Batch{
			Reference: "batch1",
			SKU:       "RUSTY_SNAPPISH",
			Quantity:  100,
			ETA:       "",
		}

		sqlInsertBatch(t, db, batch.Reference, batch.SKU, batch.Quantity, batch.ETA)

		// Create a new repository
		repo := NewSQLRepository(db)

		err := repo.Add(&batch)
		assert.Nil(t, err)

		// Check that the batch was saved
		var createdBatch Batch
		err = db.QueryRow("select reference, sku, quantity, eta from batches where reference = ?", "batch1").Scan(
			&createdBatch.Reference,
			&createdBatch.SKU,
			&createdBatch.Quantity,
			&createdBatch.ETA,
		)
		assert.Nil(t, err)
		assert.Equal(t, batch, createdBatch)
	})
}

func sqlInsertBatch(t testing.TB, db *sql.DB, reference, sku string, quantity int, eta string) {
	t.Helper()

	// Add the batch to the database
	_, err := db.Exec(
		"insert into batches (reference, sku, quantity, eta) values (?, ?, ?, ?)",
		reference, sku, quantity, eta,
	)
	assert.Nil(t, err)
}
