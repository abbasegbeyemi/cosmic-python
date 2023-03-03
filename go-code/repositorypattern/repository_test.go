package repositorypattern

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getTestDB(t testing.TB) *gorm.DB {
	t.Helper()

	// Create a temporary file for the database
	dir, _ := os.Getwd()
	file, err := os.CreateTemp(dir, "test-*.db")

	if err != nil {
		t.Fatal(err)
	}

	// Delete the file when the test is complete
	t.Cleanup(func() {
		if err := os.Remove(file.Name()); err != nil {
			t.Fatal(err)
		}
	})

	// Open a connection to the database
	db, err := gorm.Open(sqlite.Open(file.Name()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	// AutoMigrate the schema for the OrderLine model
	db.AutoMigrate(&Batch{})
	return db
}

func TestRepository(t *testing.T) {
	t.Run("Can save a batch", func(t *testing.T) {
		db := getTestDB(t)

		// Test database queries should be in a transaction
		tx := db.Begin()
		// Create a new repository
		repo := NewGormRepository(tx)

		// Create a batch of order lines
		batch := Batch{
			Reference: "batch1",
			SKU:       "RUSTY_SNAPPISH",
			Quantity:  100,
			ETA:       "",
		}

		err := repo.Add(&batch)
		assert.Nil(t, err)

		err = tx.Commit().Error
		assert.Nil(t, err)

		// Check that the batch was saved
		var created Batch

		// Load the batch from the database
		tx = db.Raw("SELECT reference, sku, quantity FROM batches WHERE reference = ?", "batch1").Scan(&created)

		assert.Nil(t, tx.Error)

		assert.Equal(t, batch, created)
	})
}
