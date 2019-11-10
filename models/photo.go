package models

import (
	// "bufio"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"scratch_maker_server/constant"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var NoSuchPhotoError = errors.New("no such photo")
var PhotoExistsError = errors.New("photo already exists")
var PhotoFileBrokenError = errors.New("photo file is broken")

// The photo model.
type Photo struct {
	BaseModel
	AuthID      uint     `json:"auth_id" gorm:"type:int" form:"auth_id"`
	BucketID    uint     `json:"bucket_id" gorm:"type:int" form:"bucket_id"`
	Name        string   `json:"name" gorm:"type:varchar(255)" form:"name"`
	Tag         string   `json:"tag" gorm:"type:varchar(255)" form:"tag"`
	Tags        []string `json:"tags" gorm:"-" form:"tags"`
	Url         string   `json:"url" gorm:"type:varchar(255)" form:"url"`
	Description string   `json:"description" gorm:"type:text" form:"description"`
	State       int      `json:"state" gorm:"type:tinyint(1)" form:"state"`
}

// Add a new photo
func AddPhoto(photoToAdd *Photo, photoFileHeader *multipart.FileHeader, context *gin.Context) (*Photo, string, error) {
	trx := Db.Begin()
	defer trx.Commit()
	// check if the photo exists, select with a WRITE LOCK
	photo := Photo{}
	trx.Set("gorm:query_option", "FOR UPDATE").
		Where("bucket_id = ? AND name = ?", photoToAdd.BucketID, photoToAdd.Name).
		First(&photo)
	if photo.ID > 0 {
		return nil, "", PhotoExistsError
	}

	photo.AuthID = photoToAdd.AuthID
	photo.BucketID = photoToAdd.BucketID
	photo.Name = photoToAdd.Name
	photo.Tag = photoToAdd.Tag
	photo.Description = photoToAdd.Description
	photo.State = 1

	err := trx.Create(&photo).Error
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	err = trx.Model(&Game{}).Where("id = ?", photoToAdd.BucketID).
		Update("size", gorm.Expr("size + ?", 1)).
		Error
	if err != nil {
		trx.Rollback()
		log.Println(err)
		return nil, "", err
	}

	// header调用Filename方法，就可以得到文件名
	filename := photoFileHeader.Filename
	fmt.Println(filename, err, filename)

	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	out, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	// 将file的内容拷贝到out
	// _, err = io.Copy(out, filename)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	context.String(http.StatusCreated, "upload successful \n")

	// TODO: upload to the tencent cloud COS
	// ......
	return nil, "", err
}

// Delete a photo by photo id.
func DeletePhotoByID(photoID uint) error {
	trx := Db.Begin()
	defer trx.Commit()

	result := trx.Where("id = ? AND state = ?", photoID, 1).Delete(Photo{})
	if err := result.Error; err != nil {
		log.Println(err)
		return err
	}
	if affected := result.RowsAffected; affected == 0 {
		return NoSuchPhotoError
	}
	return nil
}

// Delete a photo by its bucket id & its name.
func DeletePhotoByBucketAndName(bucketID uint, name string) error {
	trx := Db.Begin()
	defer trx.Commit()

	result := trx.Where("bucket_id = ? AND name = ?", bucketID, name).Delete(Photo{})
	if err := result.Error; err != nil {
		return err
	}
	if affected := result.RowsAffected; affected == 0 {
		return NoSuchPhotoError
	}
	return nil
}

// Update a photo.
func UpdatePhoto(photoToUpdate *Photo) (*Photo, error) {
	trx := Db.Begin()
	defer trx.Commit()

	photo := Photo{}
	photo.ID = photoToUpdate.ID

	result := trx.Model(&photo).Updates(*photoToUpdate)
	if err := result.Error; err != nil {
		log.Println(err)
		return &photo, err
	}
	if affected := result.RowsAffected; affected == 0 {
		return &photo, NoSuchPhotoError
	}

	return &photo, nil
}

// Update the url for a photo.
func UpdatePhotoUrl(photoID uint, url string) error {
	trx := Db.Begin()
	defer trx.Commit()

	photo := Photo{}
	photo.ID = photoID
	err := trx.Model(&photo).Update("url", url).Error
	if err != nil {
		return err
	}
	return nil
}

// Get a photo by its photo id.
func GetPhotoByID(photoID uint) (*Photo, error) {
	trx := Db.Begin()
	defer trx.Commit()

	photo := Photo{}
	err := trx.Where("id = ?", photoID).First(&photo).Error
	found := NoSuchPhotoError
	if err != nil || photo.ID == 0 {
		log.Println(err)
		found = err
	}
	found = nil
	return &photo, found
}

// Get photos by bucket id.
func GetPhotoByBucketID(bucketID uint, offset int) ([]Photo, error) {
	trx := Db.Begin()
	defer trx.Commit()

	photos := make([]Photo, 0, constant.PAGE_SIZE)
	err := trx.Where("bucket_id = ?", bucketID).
		Offset(offset).
		Limit(constant.PAGE_SIZE).
		Find(&photos).
		Error
	if err != nil {
		return photos, err
	}
	return photos, nil
}
