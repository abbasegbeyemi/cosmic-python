package domainmodelling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeBatchAndLine(t *testing.T, sku string, batchQty, lineQty int) (*Batch, *OrderLine) {
	t.Helper()
	return NewBatch("batch-001", sku, batchQty), NewOrderLine("order-123", sku, lineQty)
}

func TestBatch_Allocate(t *testing.T) {
	t.Run("allocate if available", func(t *testing.T) {
		batch, line := makeBatchAndLine(t, "SMALL-TABLE", 20, 2)
		batch.Allocate(line)
		assert.Equal(t, 18, batch.AvailableQuantity())
	})
}

func TestBatch_CanAllocate(t *testing.T) {
	t.Run("can not allocate if insufficient available quantity", func(t *testing.T) {
		batch, line := makeBatchAndLine(t, "SMALL-TABLE", 2, 20)
		assert.False(t, batch.CanAllocate(line))
	})
	t.Run("can not allocate if sku is different", func(t *testing.T) {
		batch := NewBatch("batch-001", "SMALL-TABLE", 20)
		line := NewOrderLine("order-123", "BIG-TABLE", 2)
		assert.False(t, batch.CanAllocate(line))
	})
	t.Run("can not allocate order line twice", func(t *testing.T) {
		batch, line := makeBatchAndLine(t, "SMALL-TABLE", 20, 2)
		batch.Allocate(line)
		assert.False(t, batch.CanAllocate(line))
	})
}

func TestBatch_Deallocate(t *testing.T) {
	t.Run("deallocate an allocated line", func(t *testing.T) {
		batch, line := makeBatchAndLine(t, "SMALL-TABLE", 20, 2)
		batch.Allocate(line)
		batch.Deallocate(line)
		assert.Equal(t, 20, batch.AvailableQuantity())
	})
}
