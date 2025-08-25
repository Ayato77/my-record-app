package database

import (
	"fmt"
	"my-record-app/models"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type SearchType string

const (
	SearchTitle   SearchType = "title"
	SearchContent SearchType = "content"
	SearchTag     SearchType = "tag"
	SearchGet     SearchType = "get" //Just get some records
)

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

func CreateRecord(record models.Record) (uint, error) {
	if err := DB.Create(&record).Error; err != nil {
		return 0, err
	}

	return record.ID, nil
}

// Returns the achived records, length of the total found records, error
func GetWithPaginationDB(userId, page, limit, offset int, searchType SearchType, keywords, tags []string, sort string) ([]models.Record, int64, error) {
	var records []models.Record
	query := DB.Model(&models.Record{}).Preload("Tags")

	var column string
	if len(keywords) > 0 {
		switch searchType {
		case SearchTitle:
			column = string(SearchTitle)
		case SearchContent:
			column = string(SearchContent)
		default:
			return nil, 0, fmt.Errorf("unknown search field: %s", searchType)
		}

		//search with keywords
		if DB.Dialector.Name() == "postgres" {
			for i, keyword := range keywords {
				condition := fmt.Sprintf("%s ILIKE ?", column)
				value := "%" + keyword + "%"
				if i == 0 {
					query = query.Where(condition, value)
				} else {
					query = query.Or(condition, value)
				}
			}
		} else {
			for i, keyword := range keywords {
				condition := fmt.Sprintf("LOWER(%s) LIKE LOWER(?)", column)
				value := "%" + keyword + "%"
				if i == 0 {
					query = query.Where(condition, strings.ToLower(value))
				} else {
					query = query.Or(condition, strings.ToLower(value))
				}
			}
		}
	}

	//search with tags
	if len(tags) > 0 && searchType == SearchTag {
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

func DeleteSingleRecord(userId, recordId int) (int, error) {
	res := DB.Where("records.user_id = ?", userId).Where("records.id = ?", recordId).Delete(&models.Record{})

	if res.Error != nil {
		return 0, res.Error
	}

	return int(res.RowsAffected), nil
}
