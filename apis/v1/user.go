package v1

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"scratch_maker_server/constant"
	"scratch_maker_server/models"
	"scratch_maker_server/service"
	utils "scratch_maker_server/utils"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type UserResponse struct {
	UserName string `json:"userName,omitempty"`
	Token    string `json:"token,omitempty"`
}

// Add a new auth.
func Register(context *gin.Context) {
	var requestBody models.UserReq
	if context.ShouldBind(&requestBody) == nil {
		validCheck := validation.Validation{}
		validCheck.Required(requestBody.UserName, "user_name").Message("Must have user name")
		validCheck.MaxSize(requestBody.UserName, 16, "user_name").Message("User name length can not exceed 16")
		validCheck.MinSize(requestBody.UserName, 6, "user_name").Message("User name length is at least 6")
		validCheck.Required(requestBody.Password, "password").Message("Must have password")
		validCheck.MaxSize(requestBody.Password, 16, "password").Message("Password length can not exceed 16")
		validCheck.MinSize(requestBody.Password, 6, "password").Message("Password length is at least 6")
		validCheck.Required(requestBody.Email, "email").Message("Must have email")
		validCheck.MaxSize(requestBody.Email, 128, "email").Message("Email can not exceed 128 chars")

		responseCode := constant.INVALID_PARAMS
		if !validCheck.HasErrors() {
			userEntity := models.UserReq2User(requestBody)
			if err := models.InsertUser(userEntity); err == nil {
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
			"data": "",
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
		if models.CheckUser(userName, password) {
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

	var users []models.User

	users = service.GetUserList(pageNum, pageSize, nameQuery)
	userCount := service.GetUserCount(nameQuery)

	var userPageResp = models.UserPageResp{utils.ConvertString2Int(pageNum),
		utils.ConvertString2Int(pageSize), userCount, users}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": userPageResp,
		"msg":  "query success",
	})
}

//查看用户详情
func UserDetail(context *gin.Context) {
	id := context.Query("id")
	var users []models.User

	models.Db.First(&users, id)

	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": users,
		"msg":  "query success",
	})
}

// 前台更新用户
func FrontUpdateUser(context *gin.Context) {
	var requestBody models.UserReq
	if context.ShouldBind(&requestBody) == nil {
		var user = models.UserReq2User(requestBody)
		//校验参数
		responseCode := constant.INVALID_PARAMS
		if CheckUpdateParams(requestBody) {
			trx := models.Db.Begin()
			//检查是否存在这样的用户
			userDb := models.SelectUserById(user.ID, trx)
			if userDb.ID > 0 {
				userByName := models.SelectUserByName(user.UserName, trx)
				if userByName.ID != 0 && userByName.ID != user.ID {
					context.JSON(http.StatusOK, gin.H{
						"code": "1008",
						"data": "",
						"msg":  "用户名已存在",
					})
				} else {
					// service.UpdateUser(user)
					models.Db.LogMode(true)
					//更新
					hash := md5.New()
					io.WriteString(hash, user.Password) // for safety, don't just save the plain text
					user.Password = fmt.Sprintf("%x", hash.Sum(nil))

					models.Db.Model(&user).Updates(user)
					//返回
					context.JSON(http.StatusOK, gin.H{
						"code": 200,
						"data": "",
						"msg":  "操作成功",
					})
				}
			} else {
				responseCode = constant.USER_NOT_EXIST
			}
			defer trx.Commit()
			context.JSON(http.StatusOK, gin.H{
				"code": responseCode,
				"data": "",
				"msg":  constant.GetMessage(responseCode),
			})
		} else {
			context.JSON(http.StatusOK, gin.H{
				"code": responseCode,
				"data": "",
				"msg":  constant.GetMessage(responseCode),
			})
		}
	}
}

// 后台更新用户
func AdminUpdateUser(context *gin.Context) {

}

func AdminDelUser(context *gin.Context) {
	var requestBody models.UserReq
	if context.ShouldBind(&requestBody) == nil {
		service.DeleteUser(requestBody.ID)
		context.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": "",
			"msg":  "删除成功",
		})
	}
}

func AdminLogin(context *gin.Context) {

}

func CheckUpdateParams(userReq models.UserReq) bool {
	// 校验参数
	validCheck := validation.Validation{}
	validCheck.Required(userReq.ID, "id").Message("ID缺失")

	validCheck.Required(userReq.UserName, "userName").Message("用户名缺失")
	// validCheck.Required(userReq.Email, "email").Message("密码缺失")
	// validCheck.Required(userReq.NickName, "nickName").Message("")
	// validCheck.Required(userReq.Sex, "sex").Message("Password length is at least 6")
	return !validCheck.HasErrors()
}
