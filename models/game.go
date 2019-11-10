package models

import (
	"log"
	"scratch_maker_server/constant"
	"scratch_maker_server/utils"

	"github.com/pkg/errors"
)

//游戏实体定义
type Game struct {
	BaseModel
	Name        string `json:"name" gorm:"type:varchar(64)" form:"name"`
	GameID      string `json:"gameId" gorm:"type:varchar(64)" form:"gameId"`
	CategoryID  int    `json:"categoryId" gorm:"type:tinyint(1)" form:"categoryId"`
	Path        string `json:"path" gorm:"type:varchar(255)" form:"path"`
	Status      int    `json:"status" gorm:"type:int(1)" form:"status"`
	Stars       int    `json:"stars" gorm:"type:double(2)" form:"stars"`
	LikeCount   int    `json:"likeCount" gorm:"type:int" form:"likeCount"`
	Description string `json:"description" gorm:"type:text" form:"description"`
}

var GameExistsError = errors.New("bucket already exists")
var NoSuchBucketError = errors.New("no such bucket")

// Add a new bucket.
func InsertGame(newGame *Game) error {
	trx := Db.Begin()
	defer trx.Commit()

	// check if the bucket exists, select with a WRITE LOCK.
	game := Game{}
	trx.Set("gorm:query_option", "FOR UPDATE").
		Where("name = ?", newGame.Name).
		First(&game)
	if game.ID > 0 {
		return GameExistsError
	}

	game.GameID = utils.GetCurrentTimestamp() + utils.GetRandomNum(6)
	game.Name = newGame.Name
	game.CategoryID = newGame.CategoryID
	game.Path = newGame.Path
	game.Status = 0
	game.Stars = newGame.Stars
	game.LikeCount = newGame.LikeCount
	game.Description = newGame.Description
	if err := trx.Create(&game).Error; err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Delete an existed bucket.
func DeleteBucket(bucketID uint) error {
	trx := Db.Begin()
	defer trx.Commit()

	result := trx.Where("id = ? and state = ?", bucketID, 1).Delete(Game{})
	if err := result.Error; err != nil {
		return err
	}
	if affected := result.RowsAffected; affected == 0 {
		return NoSuchBucketError
	}
	return nil
}

// Update an existed bucket.
func UpdateBucket(bucketToUpdate *Game) error {
	trx := Db.Begin()
	defer trx.Commit()

	bucket := Game{}
	bucket.ID = bucketToUpdate.ID
	result := trx.Model(&bucket).Updates(*bucketToUpdate)
	if err := result.Error; err != nil {
		return err
	}
	if affected := result.RowsAffected; affected == 0 {
		return NoSuchBucketError
	}
	return nil
}

// Get a bucket by bucket id.
func GetBucketByID(bucketID uint) (Game, error) {
	trx := Db.Begin()
	defer trx.Commit()

	bucket := Game{}
	found := NoSuchBucketError
	trx.Where("id = ?", bucketID).First(&bucket)
	if bucket.ID > 0 {
		found = nil
	}
	return bucket, found
}

// Get all buckets of the given user.
func GetBucketByAuthID(authID uint, offset int) ([]Game, error) {
	trx := Db.Begin()
	defer trx.Commit()

	buckets := make([]Game, 0, constant.PAGE_SIZE)
	err := trx.Where("auth_id = ?", authID).
		Offset(offset).
		Limit(constant.PAGE_SIZE).
		Find(&buckets).Error

	if err != nil {
		log.Println(err)
		return buckets, err
	}
	return buckets, nil
}
