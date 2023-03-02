package repositorypattern

import (
	"testing"

	"domainmodelling"
)

func TestOrderLineMapper(t *testing.T) {
	expected := []domainmodelling.OrderLine{
		{OrderID: "order1", SKU: "RED-CHAIR", Quantity: 12},
		{OrderID: "order1", SKU: "RED-TABLE", Quantity: 13},
		{OrderID: "order2", SKU: "BLUE-LIPSTICK", Quantity: 14},
	}
}
