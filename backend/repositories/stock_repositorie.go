package repositories

import (
	"backend/db"
	"backend/models"
	"fmt"

	"gorm.io/gorm/clause"
)

func StoreStock(stocks []models.Stock) (string, error) {
	if len(stocks) == 0 {
		return "", nil
	}
	DB, err := db.Conect()
	if err != nil {
		return "", fmt.Errorf("can't conect to database: %v", err)
	}

	result := DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "ticker"},
			{Name: "time"},
		},
		DoNothing: true,
	}).Create(&stocks)

	if result.Error != nil {
		return "", fmt.Errorf("can't insert data: %v", result.Error)
	}

	message := fmt.Sprintf("Inserted: %d, Ignored: %d\n", result.RowsAffected, len(stocks)-int(result.RowsAffected))
	return message, nil
}

func GetAll(page, pageSize int) ([]models.Stock, int, int, int, error) {
	DB, err := db.Conect()
	if err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't get conection: %v", err)
	}

	offset := (page - 1) * pageSize

	var stocks []models.Stock
	if err := DB.
		Offset(offset).
		Limit(pageSize).
		Find(&stocks).
		Error; err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't get all data:  %v", err)
	}

	var totalItems int64
	DB.Model(&models.Stock{}).Count(&totalItems)

	return stocks, page, offset, int(totalItems),  nil
}

func GetByTicker(ticker string, page, pageSize int) ([]models.Stock,int, int, int, error) {
	DB, err := db.Conect()
	if err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't get conection: %v", err)
	}

	offset := (page - 1) * pageSize

	var stocks []models.Stock
	if err := DB.
		Offset(offset).
		Limit(pageSize).
		Where("ticker ILIKE ?", "%"+ticker+"%").
		Find(&stocks).
		Error; err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't find %v", err)
	}

	var totalItems int64
	DB.Model(&models.Stock{}).
		Where("ticker ILIKE ?", "%"+ticker+"%").
		Count(&totalItems)


	return stocks, page, offset, int(totalItems), nil
}

func GetByCompany(company string, page, pageSize int) ([]models.Stock,int, int, int, error) {
	DB, err := db.Conect()
	if err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't get conection: %v", err)
	}

	offset := (page - 1) * pageSize

	var stocks []models.Stock
	if err := DB.
		Offset(offset).
		Limit(pageSize).
		Where("company ILIKE ?", "%"+company+"%").
		Find(&stocks).
		Error; err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't find %v", err)
	}

	var totalItems int64
	DB.Model(&models.Stock{}).
		Where("company ILIKE ?", "%"+company+"%").
		Count(&totalItems)

	return stocks, page, offset, int(totalItems), nil
}

func GetByBrokerage(brokerage string, page, pageSize int) ([]models.Stock,int, int, int, error) {
	DB, err := db.Conect()
	if err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't get conection: %v", err)
	}

	offset := (page - 1) * pageSize

	var stocks []models.Stock
	if err := DB.
		Offset(offset).
		Limit(pageSize).
		Where("brokerage ILIKE ?", "%"+brokerage+"%").
		Find(&stocks).
		Error; err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't find %v", err)
	}

	var totalItems int64
	DB.Model(&models.Stock{}).
		Where("brokerage ILIKE ?", "%"+brokerage+"%").
		Count(&totalItems)

	return stocks, page, offset, int(totalItems), nil
}

func GetByAction(action string, page, pageSize int) ([]models.Stock,int, int, int, error) {
	DB, err := db.Conect()
	if err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't get conection: %v", err)
	}

	offset := (page - 1) * pageSize

	var stocks []models.Stock
	if err := DB.
		Offset(offset).
		Limit(pageSize).
		Where("action ILIKE ?", "%"+action+"%").
		Find(&stocks).
		Error; err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't find %v", err)
	}

	var totalItems int64
	DB.Model(&models.Stock{}).
		Where("action ILIKE ?", "%"+action+"%").
		Count(&totalItems)

	return stocks, page, offset, int(totalItems), nil
}

func GetByRatingTo(ratingTo string, page, pageSize int) ([]models.Stock,int, int, int, error) {
	DB, err := db.Conect()
	if err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't get conection: %v", err)
	}

	offset := (page - 1) * pageSize

	var stocks []models.Stock
	if err := DB.
		Offset(offset).
		Limit(pageSize).
		Where("rating_to ILIKE ?", "%"+ratingTo+"%").
		Find(&stocks).
		Error; err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't find %v", err)
	}

	var totalItems int64
	DB.Model(&models.Stock{}).
		Where("rating_to ILIKE ?", "%"+ratingTo+"%").
		Count(&totalItems)

	return stocks, page, offset, int(totalItems), nil
}

func GetByRatingFrom(ratingFrom string, page, pageSize int) ([]models.Stock,int, int, int, error) {
	DB, err := db.Conect()
	if err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't get conection: %v", err)
	}

	offset := (page - 1) * pageSize

	var stocks []models.Stock
	if err := DB.
		Offset(offset).
		Limit(pageSize).
		Where("rating_from ILIKE ?", "%"+ratingFrom+"%").
		Find(&stocks).
		Error; err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't find %v", err)
	}

	var totalItems int64
	DB.Model(&models.Stock{}).
		Where("rating_from ILIKE ?", "%"+ratingFrom+"%").
		Count(&totalItems)

	return stocks, page, offset, int(totalItems), nil
}

func GetByPrice(min, max float64, page, pageSize int) ([]models.Stock,int, int, int, error) {
	DB, err := db.Conect()
	if err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't get conection: %v", err)
	}

	offset := (page - 1) * pageSize

	var stocks []models.Stock
	if err := DB.
		Offset(offset).
		Limit(pageSize).
		Where("target_to BETWEEN ? AND ?", min, max).
		Find(&stocks).
		Error; err != nil {
		return nil, 1, 20, 0, fmt.Errorf("can't find %v", err)
	}

	var totalItems int64
	DB.Model(&models.Stock{}).
		Where("target_to BETWEEN ? AND ?", min, max).
		Count(&totalItems)

	return stocks, page, offset, int(totalItems), nil
}

func GetByRecommendation() ([]models.Stock, error) {
	DB, err := db.Conect()
	if err != nil {
		return nil, fmt.Errorf("can't get conection: %v", err)
	}

	var stocks []models.Stock
	if err := DB.Model(&models.Stock{}).
		Order("time DESC").
		Find(&stocks).
		Error; err != nil {
		return nil, fmt.Errorf("can't find %v", err)
	}

	return stocks, nil
}

