package repository

import (
	"github.com/alir32a/jupiter/internal/model"
	"gorm.io/gorm"
	"math"
)

func Paginate(req *model.Pagination) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if req.CurrentPage <= 0 {
			req.CurrentPage = 1
		}

		if req.PageSize > 100 || req.PageSize <= 0 {
			req.PageSize = 100
		}

		var total int64
		if err := tx.Session(&gorm.Session{Initialized: true}).Count(&total).Error; err != nil {
			return nil
		}

		req.Total = int(total)
		req.TotalPages = int(math.Ceil(float64(total) / float64(req.PageSize)))

		offset := (req.CurrentPage - 1) * req.PageSize

		return tx.Offset(offset).Limit(req.PageSize)
	}
}
