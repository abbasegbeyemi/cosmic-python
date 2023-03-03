package domainmodelling

type Batch struct {
	Reference   string
	SKU         string
	Quantity    int
	ETA         string
	Allocations []*OrderLine
}

func NewBatch(ref, sku string, qty int) *Batch {
	return &Batch{ref, sku, qty, "", nil}
}

func (b *Batch) Allocate(line *OrderLine) {
	if b.CanAllocate(line) {
		b.Allocations = append(b.Allocations, line)
	}
}

func (b *Batch) CanAllocate(line *OrderLine) bool {
	// Check that an allocation is not already made for this order line
	for _, a := range b.Allocations {
		if a.OrderID == line.OrderID {
			return false
		}
	}
	return b.SKU == line.SKU && b.Quantity >= line.Quantity
}

func (b *Batch) AvailableQuantity() int {
	return b.Quantity - b.totalAllocatedQuantity()
}

func (b *Batch) totalAllocatedQuantity() int {
	sum := 0
	for _, a := range b.Allocations {
		sum += a.Quantity
	}
	return sum
}

type OrderLine struct {
	OrderID  string
	SKU      string
	Quantity int
}

func NewOrderLine(orderID, sku string, qty int) *OrderLine {
	return &OrderLine{orderID, sku, qty}
}

func (b *Batch) Deallocate(line *OrderLine) {
	// Find the allocation for this order line
	for i, a := range b.Allocations {
		if a.OrderID == line.OrderID {
			// Remove the allocation from the slice
			b.Allocations = append(b.Allocations[:i], b.Allocations[i+1:]...)
			return
		}
	}
}
