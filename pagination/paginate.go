package pagination

import "gorm.io/gorm"

func Paginate[T any](tx *gorm.DB, page int, pageSize int, result *PageResult[T]) (*gorm.DB, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	var total int64

	if err := tx.Count(&total).Error; err != nil {
		return tx, err
	}

	offset := (page - 1) * pageSize

	if err := tx.Limit(pageSize).Offset(offset).Find(&result.Records).Error; err != nil {
		return tx, err
	}

	result.Total = total
	result.Page = page
	result.PageSize = pageSize

	return tx, nil
}
