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
	db.AutoMigrate(&Batch{}, &OrderLine{})
	return db
}

func TestGormRepository(t *testing.T) {
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
		tx = db.Raw("SELECT * FROM batches WHERE reference = ?", "batch1").Scan(&created)

		assert.Nil(t, tx.Error)

		assert.Equal(t, batch, created)
	})

	t.Run("Can retrieve a batch", func(t *testing.T) {
		db := getTestDB(t)
		expectedOrderLine := OrderLine{
			OrderID:  "order1",
			SKU:      "CAT-TREE",
			Quantity: 12,
		}
		orderLineId := gormInsertOrderLine(t, db, expectedOrderLine.OrderID, expectedOrderLine.SKU, expectedOrderLine.Quantity)
		batchId := gormInsertBatch(t, db, "batch1", "CAT-TREE", 100, "2021-01-01")

		gormInsertBatch(t, db, "batch2", "CAT-HOUSE", 100, "2022-01-01")
		gormInsertAllocation(t, db, orderLineId, batchId)

		repo := NewGormRepository(db)
		retrieved := repo.Get("batch1")

		assert.Equal(t, "batch1", retrieved.Reference)
		assert.Equal(t, "CAT-TREE", retrieved.SKU)
		assert.Equal(t, 100, retrieved.Quantity)
		assert.Equal(t, "2021-01-01", retrieved.ETA)

		retrievedOrderLine := retrieved.Allocations[0]
		assert.Equal(t, expectedOrderLine.OrderID, retrievedOrderLine.OrderID)
		assert.Equal(t, expectedOrderLine.SKU, retrievedOrderLine.SKU)
		assert.Equal(t, expectedOrderLine.Quantity, retrievedOrderLine.Quantity)
	})

	t.Run("Can retrieve a list of batches", func(t *testing.T) {
		db := getTestDB(t)

		gormInsertBatch(t, db, "batch1", "CAT-TREE", 100, "2021-01-01")
		gormInsertBatch(t, db, "batch2", "CAT-HOUSE", 100, "2022-01-01")

		repo := NewGormRepository(db)
		retrieved := repo.List()
		assert.Equal(t, 2, len(retrieved))
	})
}

func gormInsertOrderLine(t testing.TB, db *gorm.DB, orderId string, sku string, quantity int) uint {
	t.Helper()
	orderLine := &OrderLine{
		OrderID:  orderId,
		SKU:      sku,
		Quantity: quantity,
	}
	tx := db.Create(orderLine)
	assert.Nil(t, tx.Error)
	return orderLine.ID
}

func gormInsertBatch(t testing.TB, db *gorm.DB, reference string, sku string, quantity int, eta string) uint {
	t.Helper()
	batch := &Batch{
		Reference: reference,
		SKU:       sku,
		Quantity:  quantity,
		ETA:       eta,
	}
	tx := db.Create(batch)
	assert.Nil(t, tx.Error)
	return batch.ID
}

func gormInsertAllocation(t testing.TB, db *gorm.DB, orderLineId uint, batchId uint) {
	t.Helper()
	// Get the order line
	var orderLine OrderLine
	tx := db.First(&orderLine, orderLineId)

	assert.Nil(t, tx.Error)
	db.Model(&Batch{
		ID: batchId,
	}).Association("Allocations").Append(&orderLine)
	assert.Nil(t, tx.Error)
}
