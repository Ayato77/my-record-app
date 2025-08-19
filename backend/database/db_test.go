package database

import (
	"my-record-app/models"

	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitTestDB() *gorm.DB {
	//create a db in local memory with sqlite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.AutoMigrate(&models.Record{}, &models.User{}); err != nil {
		panic("failed to automigrate database")
	}

	DB = db
	return db
}

func TestGetWithPaginationDB(t *testing.T) {
	InitTestDB()

	//Create New Record Obj
	var tags1 = []models.Tag{
		{
			Name: "tamago",
			ID:   1,
		},
		{
			Name: "niku",
			ID:   2,
		},
	}

	var tags2 = []models.Tag{
		{
			Name: "tori",
			ID:   3,
		},
		{
			Name: "nasubi",
			ID:   4,
		},
	}

	var tags3 = []models.Tag{
		{
			Name: "tori",
			ID:   3,
		},
		{
			Name: "negi",
			ID:   5,
		},
	}

	newRecord1 := models.Record{
		UserID:  1,
		Title:   "A",
		Content: "Content A",
		Rating:  1,
		Tags:    tags1,
	}

	newRecord2 := models.Record{
		UserID:  1,
		Title:   "B",
		Content: "Content B",
		Rating:  3,
		Tags:    tags2,
	}

	newRecord3 := models.Record{
		UserID:  1,
		Title:   "C",
		Content: "Content C",
		Rating:  4,
		Tags:    tags3,
	}

	newRecord4 := models.Record{
		UserID:  2,
		Title:   "D",
		Content: "Content D",
		Rating:  4,
		Tags:    tags3,
	}

	//Insert records by DB.Create
	if err := DB.Create(&newRecord1).Error; err != nil {
		t.Errorf("Error: DB Create: %s", err)
		return
	}

	if err := DB.Create(&newRecord2).Error; err != nil {
		t.Errorf("Error: DB Create: %s", err)
		return
	}

	if err := DB.Create(&newRecord3).Error; err != nil {
		t.Errorf("Error: DB Create: %s", err)
		return
	}

	if err := DB.Create(&newRecord4).Error; err != nil {
		t.Errorf("Error: DB Create: %s", err)
		return
	}

	//Call GetWithPaginationDB
	var tagSingle = []string{"tori"}
	var tagDouble = []string{"tori", "nasubi"}
	var tagInvalid = []string{"onion"}
	records, total, err := GetWithPaginationDB(1, 1, 10, 0, tagSingle, "")
	if err != nil {
		t.Errorf("Error: GetWithPaginationDB: %s", err)
		return
	}

	if total != 2 {
		t.Error("Too many or few records are found")
		return
	}

	var tagCounter int = 0
	for _, item := range records[0].Tags {
		if tagSingle[0] == item.Name {
			tagCounter = tagCounter + 1
		}
	}

	if tagCounter != 1 {
		t.Error("A record without desired tag is included or no record found (Double)")
	} else {
		t.Log(records)
		tagCounter = 0
	}

	records, total, err = GetWithPaginationDB(1, 1, 10, 0, tagDouble, "")
	if err != nil {
		t.Errorf("Error: GetWithPaginationDB Double: %s", err)
		return
	}

	if total > 1 {
		t.Error("Too many records are found (Double)")
		return
	}

	for _, item := range records[0].Tags {
		for i := range tagDouble {
			if tagDouble[i] == item.Name {
				tagCounter = tagCounter + 1
			}
		}
	}

	if tagCounter != 2 {
		t.Error("A record without desired tags is included")
	} else {
		t.Log(records)
		tagCounter = 0
	}

	records, total, err = GetWithPaginationDB(1, 1, 10, 0, tagInvalid, "")
	if err != nil {
		t.Errorf("Error: GetWithPaginationDB Invalid: %s", err)
		return
	}

	if total > 0 {
		t.Error("Too many records are found (Invalid)")
		return
	}

	t.Log(records)

}
