from dataclasses import dataclass, field
from typing import Optional
from datetime import date


@dataclass(frozen=True)
class OrderLine:
    """
    Order line
    """

    order_id: str
    sku: str
    quantity: int


@dataclass(frozen=True)
class Batch:
    """
    Batch of stock
    """

    reference: str
    sku: str
    quantity: int
    allocations: set[OrderLine] = field(default_factory=set)
    eta: Optional[date] = None

    @property
    def available_quantity(self) -> int:
        return self.quantity - self._total_allocated_quantity()

    def _total_allocated_quantity(self) -> int:
        return sum(allocation.quantity for allocation in self.allocations)

    def allocate(self, line: OrderLine):
        if self.can_allocate(line):
            self.allocations.add(line)

    def deallocate(self, line: OrderLine):
        if line in self.allocations:
            self.allocations.remove(line)

    def can_allocate(self, line: OrderLine) -> bool:
        return (
            line not in self.allocations
            and self.sku == line.sku
            and self.available_quantity >= line.quantity
        )
