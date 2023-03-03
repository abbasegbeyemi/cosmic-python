package repositorypattern

import (
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

type Batch struct {
	Reference string `gorm:"primaryKey"`
	SKU       string
	Quantity  int
	ETA       string
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db}
}

func (r *GormRepository) Add(batch *Batch) error {
	return r.db.Create(batch).Error
}
