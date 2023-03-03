package repositorypattern

// import (
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// type OrderLine struct {
// 	OrderID  string `gorm:"column:order_id;primaryKey"`
// 	SKU      string `gorm:"column:sku"`
// 	Quantity int    `gorm:"column:quantity"`
// }

// func getTestDB(t testing.TB) *gorm.DB {
// 	t.Helper()

// 	// Create a temporary file for the database
// 	dir, _ := os.Getwd()
// 	file, err := ioutil.TempFile(dir, "test-*.db")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Delete the file when the test is complete
// 	t.Cleanup(func() {
// 		if err := os.Remove(file.Name()); err != nil {
// 			t.Fatal(err)
// 		}
// 	})

// 	// Open a connection to the database
// 	db, err := gorm.Open(sqlite.Open(file.Name()), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("failed to connect database: %v", err)
// 	}
// 	// AutoMigrate the schema for the OrderLine model
// 	db.AutoMigrate(&OrderLine{})
// 	return db
// }

// func TestOrderLineMapper(t *testing.T) {

// 	t.Run("Can load order lines from the database", func(t *testing.T) {

// 		db := getTestDB(t)
// 		tx := db.Exec(`
//     INSERT INTO order_lines (order_id, sku, quantity)
//     VALUES
//         (?, ?, ?),
//         (?, ?, ?),
//         (?, ?, ?)`,
// 			"order1", "RED-CHAIR", 12,
// 			"order2", "RED-TABLE", 13,
// 			"order3", "BLUE-LIPSTICK", 14)
// 		if tx.Error != nil {
// 			t.Fatal(tx.Error)
// 		}

// 		orders := []OrderLine{}
// 		expected := []OrderLine{
// 			{OrderID: "order1", SKU: "RED-CHAIR", Quantity: 12},
// 			{OrderID: "order2", SKU: "RED-TABLE", Quantity: 13},
// 			{OrderID: "order3", SKU: "BLUE-LIPSTICK", Quantity: 14},
// 		}

// 		// Query the database for all the order lines
// 		db.Find(&orders)

// 		assert.Equal(t, expected, orders)
// 	})

// 	t.Run("Can save order lines to the database", func(t *testing.T) {
// 		db := getTestDB(t)
// 		newLine := OrderLine{OrderID: "order4", SKU: "CAT_TREE", Quantity: 12}

// 		result := db.Create(&newLine)
// 		assert.Nil(t, result.Error)

// 		// Query the database for the order line
// 		var order OrderLine
// 		db.First(&order, "order_id = ?", "order4")
// 		assert.Equal(t, newLine, order)
// 	})
// }
