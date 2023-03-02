"""
Test the allocation system
"""

from domain_modelling.models import Batch, OrderLine


class TestBatch:
    def test_allocating_to_a_batch_reduces_the_available_quantity(self):
        batch = Batch(reference="batch-001", sku="SMALL-TABLE", quantity=20, eta=None)
        line = OrderLine(order_id="order-ref", sku="SMALL-TABLE", quantity=2)
        batch.allocate(line)
        assert batch.available_quantity == 18

    def test_cant_allocate_if_skus_do_not_match(self):
        batch = Batch(reference="batch-001", sku="SMALL-TABLE", quantity=20, eta=None)
        line = OrderLine(order_id="order-ref", sku="LARGE-TABLE", quantity=2)
        assert batch.can_allocate(line) is False

    def test_cant_allocate_if_available_quantity_is_too_low(self):
        batch = Batch(reference="batch-001", sku="SMALL-TABLE", quantity=20, eta=None)
        line = OrderLine(order_id="order-ref", sku="SMALL-TABLE", quantity=21)
        assert batch.can_allocate(line) is False

    def test_cant_allocate_order_line_twice(self):
        batch = Batch(reference="batch-001", sku="SMALL-TABLE", quantity=20, eta=None)
        line = OrderLine(order_id="order-ref", sku="SMALL-TABLE", quantity=2)
        batch.allocate(line)
        assert batch.can_allocate(line) is False

    def test_can_deallocate(self):
        batch = Batch(reference="batch-001", sku="SMALL-TABLE", quantity=20, eta=None)
        line = OrderLine(order_id="order-ref", sku="SMALL-TABLE", quantity=2)
        batch.allocate(line)
        batch.deallocate(line)
        assert batch.available_quantity == 20
