package database

import (
	"my-record-app/models"

	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitTestDB() *gorm.DB {
	if DB != nil {
		return DB
	}
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
		},
		{
			Name: "niku",
		},
	}

	var tags2 = []models.Tag{
		{
			Name: "tori",
		},
		{
			Name: "nasubi",
		},
	}

	var tags3 = []models.Tag{
		{
			Name: "tori",
		},
		{
			Name: "negi",
		},
	}

	newRecord1 := models.Record{
		UserID:  1,
		Title:   "オムライス",
		Content: "とろーり卵",
		Rating:  1,
		Tags:    tags1,
	}

	newRecord2 := models.Record{
		UserID:  1,
		Title:   "とりと野菜の煮込み",
		Content: "和風のチキン。ナス",
		Rating:  3,
		Tags:    tags2,
	}

	newRecord3 := models.Record{
		UserID:  1,
		Title:   "親子丼",
		Content: "とりもも肉、その他",
		Rating:  4,
		Tags:    tags3,
	}

	newRecord4 := models.Record{
		UserID:  2,
		Title:   "ねぎま",
		Content: "串、炭火",
		Rating:  4,
		Tags:    tags3,
	}

	//Insert records by DB.Create
	if err := DB.Create(&newRecord1).Error; err != nil {
		t.Errorf("Error: DB Create: %s", err)
	}

	if err := DB.Create(&newRecord2).Error; err != nil {
		t.Errorf("Error: DB Create: %s", err)
	}

	if err := DB.Create(&newRecord3).Error; err != nil {
		t.Errorf("Error: DB Create: %s", err)
	}

	if err := DB.Create(&newRecord4).Error; err != nil {
		t.Errorf("Error: DB Create: %s", err)
	}

	//Call GetWithPaginationDB
	var tagEmpty []string
	var tagSingle = []string{"tori"}
	var tagDouble = []string{"tori", "nasubi"}
	var tagInvalid = []string{"onion"}
	var keywordEmpty []string
	var keywordSingle = []string{"オム"}
	var keywordDouble = []string{"とり", "チキン"}
	records, total, err := GetWithPaginationDB(1, 1, 10, 0, SearchTag, keywordEmpty, tagSingle, "")
	if err != nil {
		t.Errorf("Error: GetWithPaginationDB: %s", err)
	}

	if total != 2 {
		t.Error("Too many or few records are found")
	}

	t.Log(records)

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

	records, total, err = GetWithPaginationDB(1, 1, 10, 0, SearchTag, keywordEmpty, tagDouble, "")
	if err != nil {
		t.Errorf("Error: GetWithPaginationDB Double: %s", err)
	}

	if total > 1 {
		t.Error("Too many records are found (Double)")
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

	records, total, err = GetWithPaginationDB(1, 1, 10, 0, SearchTag, keywordEmpty, tagInvalid, "")
	if err != nil {
		t.Errorf("Error: GetWithPaginationDB Invalid: %s", err)
	}

	if total > 0 {
		t.Error("Too many records are found (Invalid)")
	}

	t.Log(records)

	//tests search with keys
	records, total, err = GetWithPaginationDB(1, 1, 10, 0, SearchTitle, keywordSingle, tagEmpty, "")
	if err != nil {
		t.Errorf("Error: GetWithPaginationDB Invalid: %s", err)
	}

	if total != 1 {
		t.Error("Too many or few records are found")
	}

	if records[0].ID != newRecord1.ID {
		t.Error("Error: wrong record returned")
	}

	t.Log(records)

	records, total, err = GetWithPaginationDB(1, 1, 10, 0, SearchContent, keywordDouble, tagEmpty, "")
	if err != nil {
		t.Errorf("Error: GetWithPaginationDB Invalid: %s", err)
	}

	if total != 2 {
		t.Error("Too many or few records are found")
	}

	t.Log(records)

}

func TestDeleteSingleRecord(t *testing.T) {
	InitTestDB()

	userId := 1
	invalidUserId := 2
	//create record
	var tags1 = []models.Tag{
		{
			Name: "tamago",
		},
		{
			Name: "niku",
		},
	}

	newRecord1 := models.Record{
		UserID:  uint(userId),
		Title:   "E",
		Content: "Content E",
		Rating:  1,
		Tags:    tags1,
	}

	recordId, err := CreateRecord(newRecord1)
	if err != nil {
		t.Errorf("Error: DB Create: %s", err)
	}
	//try to delete a record with not existing user_id
	rowsAffected, err := DeleteSingleRecord(invalidUserId, int(recordId))
	if err == nil && rowsAffected != 0 {
		t.Errorf("Error: DB Delete with a invalid user ID was succeeded: %s", err)
	}

	//try to delete a record with not existing record_id
	rowsAffected, err = DeleteSingleRecord(userId, int(recordId+100))
	if err == nil && rowsAffected != 0 {
		t.Errorf("Error: DB Delete with a invalid record ID was succeeded: %s", err)
	}
	//delete a record with a valid id
	rowsAffected, err = DeleteSingleRecord(userId, int(recordId))
	if err != nil {
		t.Errorf("Error: DB Delete with valid IDs failed: %s", err)
	}

	if rowsAffected == 0 {
		t.Error("No record is deleted")
	}

	t.Log(rowsAffected)
}
