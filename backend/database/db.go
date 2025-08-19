package database

import (
	"fmt"
	"my-record-app/models"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	dsn := os.Getenv("DB_URL")
	fmt.Println("DB_URL: ", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}

// TODO: test it!
// Returns the achived records, length of the total found records, error
func GetWithPaginationDB(userId, page, limit, offset int, tags []string, sort string) ([]models.Record, int64, error) {
	var records []models.Record
	query := DB.Model(&models.Record{})
	if len(tags) > 0 {
		query = query.Joins("JOIN record_tag rt ON rt.record_id = records.id").
			Joins("JOIN tags t ON t.id = rt.tag_id").
			Where("records.user_id = ?", userId).
			Where("t.name IN ?", tags).
			Group("records.id").
			Having("COUNT(DISTINCT t.id) = ?", len(tags)). // AND condition for tags
			Preload("Tags").
			Find(&records)
	}

	var total int64
	query.Count(&total)

	if err := query.Order(sort).Limit(limit).Offset(offset).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}
