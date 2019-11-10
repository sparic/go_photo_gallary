package v1

import (
	"log"
	"net/http"
	"scratch_maker_server/constant"
	"scratch_maker_server/models"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 添加游戏
func FrontAddGame(context *gin.Context) {
	responseCode := constant.INVALID_PARAMS
	gameToAdd := models.Game{}
	if err := context.ShouldBindWith(&gameToAdd, binding.Form); err != nil {
		log.Println(err)
		context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": responseCode,
			"data": make(map[string]string),
			"msg":  constant.GetMessage(responseCode),
		})
		return
	}

	validCheck := validation.Validation{}

	// validCheck.Required(gameToAdd.Name, "bucket_name").Message("Must have bucket name")
	// validCheck.MaxSize(gameToAdd.Name, 64, "bucket_name").Message("Bucket name length can not exceed 64")

	if !validCheck.HasErrors() {
		if err := models.InsertGame(&gameToAdd); err != nil {
			if err == models.GameExistsError {
				responseCode = constant.BUCKET_ALREADY_EXIST
			} else {
				responseCode = constant.INTERNAL_SERVER_ERROR
			}
		} else {
			responseCode = constant.BUCKET_ADD_SUCCESS
		}
	} else {
		for _, err := range validCheck.Errors {
			log.Println(err.Message)
		}
	}

	data := make(map[string]string)
	data["bucket_name"] = gameToAdd.Name

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": data,
		"msg":  constant.GetMessage(responseCode),
	})
}

// 将游戏加入收藏或取消收藏
func Favourate(context *gin.Context) {
	// ......
}

// Update an existed bucket.
func UpdateBucket(context *gin.Context) {
	// ......
}

// 查看游戏详情
func GameDetail(context *gin.Context) {
	// ......
}

// 前台游戏列表
func FrontGameList(context *gin.Context) {
	// ......
}

// 前台游戏评论
func FrontGameComment(context *gin.Context) {

}

// 前台删除游戏评论
func FrontDeleteComment(context *gin.Context) {

}

//-------------- 后台API ----------//
func AdminAddGame(context *gin.Context) {

}

func AdminUpdateGame(context *gin.Context) {

}

func AdminGameDetail(context *gin.Context) {

}

func AdminGameList(context *gin.Context) {

}

func AdminGameDelete(context *gin.Context) {

}
