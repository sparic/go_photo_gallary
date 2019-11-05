package v1

import (
	"fmt"
	"go_photo_gallary/constant"
	"go_photo_gallary/models"
	"go_photo_gallary/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	UserName string `json:"userName,omitempty"`
	Token    string `json:"token,omitempty"`
}

// Add a new auth.
func Register(context *gin.Context) {

	// userName := context.PostForm("user_name")
	// password := context.PostForm("password")
	// email := context.PostForm("email")

	var requestBody models.Auth
	if context.ShouldBind(&requestBody) == nil {
		userName := requestBody.UserName
		password := requestBody.Password
		email := requestBody.Email
		// set up param validation

		validCheck := validation.Validation{}
		validCheck.Required(userName, "user_name").Message("Must have user name")
		validCheck.MaxSize(userName, 16, "user_name").Message("User name length can not exceed 16")
		validCheck.MinSize(userName, 6, "user_name").Message("User name length is at least 6")
		validCheck.Required(password, "password").Message("Must have password")
		validCheck.MaxSize(password, 16, "password").Message("Password length can not exceed 16")
		validCheck.MinSize(password, 6, "password").Message("Password length is at least 6")
		validCheck.Required(email, "email").Message("Must have email")
		validCheck.MaxSize(email, 128, "email").Message("Email can not exceed 128 chars")

		responseCode := constant.INVALID_PARAMS
		if !validCheck.HasErrors() {
			if err := models.AddAuth(userName, password, email); err == nil {
				responseCode = constant.USER_ADD_SUCCESS
			} else {
				responseCode = constant.USER_ALREADY_EXIST
			}
		} else {
			for _, err := range validCheck.Errors {
				log.Println(err)
			}
		}

		context.JSON(http.StatusOK, gin.H{
			"code": responseCode,
			"data": "aaa",
			"msg":  constant.GetMessage(responseCode),
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": "",
			"msg":  "添加失败，参数有误",
		})
	}
}

// Check if an auth is valid.
func FrontLogin(context *gin.Context) {

	userName := context.PostForm("user_name")
	password := context.PostForm("password")

	// 校验参数
	validCheck := validation.Validation{}
	validCheck.Required(userName, "user_name").Message("Must have user name")
	validCheck.MaxSize(userName, 16, "user_name").Message("User name length can not exceed 16")
	validCheck.MinSize(userName, 6, "user_name").Message("User name length is at least 6")
	validCheck.Required(password, "password").Message("Must have password")
	validCheck.MaxSize(password, 16, "password").Message("Password length can not exceed 16")
	validCheck.MinSize(password, 6, "password").Message("Password length is at least 6")

	responseCode := constant.INVALID_PARAMS
	jwtStringRs := ""
	if !validCheck.HasErrors() {
		if models.CheckAuth(userName, password) {
			if jwtString, err := utils.GenerateJWT(userName); err != nil {
				responseCode = constant.JWT_GENERATION_ERROR
			} else {
				// pass auth validation
				// 1. set JWT to user's cookie
				// 2. add user to the Redis
				jwtStringRs = jwtString
				log.Printf("jwtString->%v\n", jwtString)
				context.SetCookie(constant.JWT, jwtString,
					constant.COOKIE_MAX_AGE, constant.SERVER_PATH,
					constant.SERVER_DOMAIN, true, true)
				if err = utils.AddAuthToRedis(userName); err != nil {
					responseCode = constant.INTERNAL_SERVER_ERROR
				} else {
					responseCode = constant.USER_AUTH_SUCCESS
				}
			}
		} else {
			responseCode = constant.USER_AUTH_ERROR
		}
	} else {
		for _, err := range validCheck.Errors {
			log.Println(err)
		}
	}

	user := UserResponse{}
	user.UserName = userName
	user.Token = jwtStringRs

	context.JSON(http.StatusOK, gin.H{
		"code": responseCode,
		"data": user,
		"msg":  constant.GetMessage(responseCode),
	})
}

// type UserQuery struct {
// 	UserName string `gorm:"column:user_name"`
// }

func ListUser(context *gin.Context) {
	nameQuery := context.Query("userName")
	pageNum := context.Query("pageNum")
	pageSize := context.Query("pageSize")

	// var conditions = map[string]interface{}{
	// 	"user_name like ?": "%" + nameQuery + "%",
	// }

	var users []models.Auth
	// firstname := context.DefaultQuery("firstname", "Guest")
	// lastname := context.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
	models.Db.LogMode(true)

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		fmt.Println("err", err)
	}

	pageNumInt, err := strconv.Atoi(pageNum)
	if err != nil {
		fmt.Println("err", err)
	}

	models.Db.Where("user_name LIKE ?", "%"+nameQuery+"%").Order("id asc").Offset((pageNumInt - 1) * pageSizeInt).Limit(pageSize).Find(&users)
	// context.String(http.StatusOK, "Hello %s %s", firstname, lastname)

	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": users,
		"msg":  "query success",
	})
}

//查看用户详情
func UserDetail(context *gin.Context) {
	id := context.Query("id")
	var users []models.Auth

	models.Db.First(&users, id)

	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": users,
		"msg":  "query success",
	})
}

// 前台更新用户
func FrontUpdateUser(context *gin.Context) {

}

// 后台更新用户
func AdminUpdateUser(context *gin.Context) {

}

func AdminLogin(context *gin.Context) {

}
