package v1

import (
	"fmt"
	"go_photo_gallary/constant"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Add a new photo
func AddPhoto(context *gin.Context) {
	log.Printf("addPhoto clicked!!!")
	responseCode := constant.INVALID_PARAMS

	photoFile, header, err := context.Request.FormFile("photo")
	if err != nil {
		log.Println(err)
	}

	// header调用Filename方法，就可以得到文件名
	fileName := header.Filename
	fmt.Println(photoFile, err, fileName)

	if err != nil {
		fmt.Println("err", err)
	}

	// log, _ := zap.NewDevelopment()

	/*
		photo := models.Photo{}
		paramErr := context.ShouldBindWith(&photo, binding.Form)

		if err != nil || paramErr != nil {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": responseCode,
				"data": make(map[string]string),
				"msg":  constant.GetMessage(responseCode),
			})
			return
		}

		validCheck := validation.Validation{}
		validCheck.Required(photo.AuthID, "auth_id").Message("Must have auth id")
		validCheck.Required(photo.BucketID, "bucket_id").Message("Must have bucket id")
		validCheck.Required(photo.Name, "photo_name").Message("Must have photo name")
		validCheck.MaxSize(photo.Name, 255, "photo_name").Message("Photo name len must not exceed 255")

		data := make(map[string]interface{})
		photoToAdd := &models.Photo{BucketID: photo.BucketID, AuthID: photo.AuthID,
			Name: photo.Name, Description: photo.Description,
			Tag: strings.Join(photo.Tags, ";")}

		if !validCheck.HasErrors() {
			if photoToAdd, uploadID, err := models.AddPhoto(photoToAdd, photoFile); err != nil {
				if err == models.PhotoExistsError {
					responseCode = constant.PHOTO_ALREADY_EXIST
				} else {
					responseCode = constant.INTERNAL_SERVER_ERROR
				}
			} else {
				responseCode = constant.PHOTO_ADD_IN_PROCESS
				data["photo"] = *photoToAdd
				data["photo_upload_id"] = uploadID
			}
		} else {
			for _, err := range validCheck.Errors {
				log.Println(err.Message)
			}
		}
	*/
	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	out, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()

	// 将file的内容拷贝到out
	_, err = io.Copy(out, photoFile)
	if err != nil {
		log.Fatal(err)
	}

	context.String(http.StatusCreated, "upload successful \n")

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": "data",
		"msg":  constant.GetMessage(responseCode),
	})
}

// Delete an existed photo.
func DeletePhoto(context *gin.Context) {
	// ......
}

// Update an existed photo.
func UpdatePhoto(context *gin.Context) {
	// ......
}

// Get a photo by photo id.
func GetPhotoByID(context *gin.Context) {
	// ......
}

// Get photos by bucket id.
func GetPhotoByBucketID(context *gin.Context) {
	// ......
}

// Get the upload status of a photo by upload id.
func GetPhotoUploadStatus(context *gin.Context) {
	// ......
}
