package repositorypattern

import (
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

type Batch struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	Reference   string
	SKU         string
	Quantity    int
	ETA         string
	Allocations []*OrderLine
}

type OrderLine struct {
	ID       uint `gorm:"primaryKey;autoIncrement"`
	OrderID  string
	SKU      string
	Quantity int
	BatchID  uint
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db}
}

func (r *GormRepository) Add(batch *Batch) error {
	return r.db.Create(batch).Error
}

func (r *GormRepository) Get(reference string) *Batch {
	var batch Batch
	r.db.Preload("Allocations").First(&batch, "reference = ?", reference)
	return &batch
}

func (r *GormRepository) List() []*Batch {
	var batches []*Batch
	r.db.Preload("Allocations").Find(&batches)
	return batches
}
