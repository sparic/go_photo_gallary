package routers

import (
	v1 "go_photo_gallary/apis/v1"
	"go_photo_gallary/middleware"

	"github.com/gin-gonic/gin"
)

// a global router
var Router *gin.Engine

// Init router, adding paths to it.
func init() {
	Router = gin.Default()
	checkAuthMdw := middleware.GetAuthMiddleware()        // middleware for authentication
	refreshMdw := middleware.GetRefreshMiddleware()       // middleware for refresh auth token
	paginationMdw := middleware.GetPaginationMiddleware() // middleware for pagination

	// api group for v1
	v1AdminGroup := Router.Group("/admin/api/v1")
	v1FrontGroup := Router.Group("api/v1")
	{
		// 前台用户模块
		userFrontGroup := v1FrontGroup.Group("/user")
		{
			userFrontGroup.POST("/register", v1.Register)
			userFrontGroup.POST("/login", v1.FrontLogin)
			userFrontGroup.POST("/update", v1.FrontUpdateUser)
		}
		// 后台用户模块
		userAdminGroup := v1AdminGroup.Group("/user")
		{
			userAdminGroup.POST("/login", v1.AdminLogin)
			userAdminGroup.GET("/list", checkAuthMdw, refreshMdw, v1.ListUser)
			userAdminGroup.GET("/getOne", checkAuthMdw, refreshMdw, v1.UserDetail)
			userAdminGroup.PUT("/update", checkAuthMdw, refreshMdw, v1.AdminUpdateUser)
		}

		// 前台游戏模块
		gameFrontGroup := v1FrontGroup.Group("/game")
		{
			// must check auth & refresh auth token before any operation
			gameFrontGroup.POST("/add", checkAuthMdw, refreshMdw, v1.FrontAddGame)
			//添加/取消收藏游戏
			gameFrontGroup.DELETE("/favorate", checkAuthMdw, refreshMdw, v1.Favourate)
			gameFrontGroup.GET("/getOne", checkAuthMdw, refreshMdw, v1.GameDetail)
			gameFrontGroup.GET("/list", checkAuthMdw, refreshMdw, paginationMdw, v1.FrontGameList)
			//评论游戏
			gameFrontGroup.POST("/comment", checkAuthMdw, refreshMdw, paginationMdw, v1.FrontGameComment)
			//删除游戏评论
			gameFrontGroup.DELETE("/deleteComment", checkAuthMdw, refreshMdw, paginationMdw, v1.FrontDeleteComment)
		}
		// 后台游戏模块
		gameAdminGroup := v1AdminGroup.Group("/game")
		{
			//新增游戏
			gameAdminGroup.POST("/add", checkAuthMdw, refreshMdw, v1.AdminAddGame)
			gameAdminGroup.PUT("/update", checkAuthMdw, refreshMdw, v1.AdminUpdateGame)
			gameAdminGroup.GET("/getOne", checkAuthMdw, refreshMdw, v1.AdminGameDetail)
			gameAdminGroup.PUT("/list", checkAuthMdw, refreshMdw, v1.AdminGameList)
			gameAdminGroup.DELETE("/delete", checkAuthMdw, refreshMdw, v1.AdminGameDelete)
		}

		// 前台上传文件
		fileFrontGroup := v1FrontGroup.Group("/file")
		{
			// must check auth & refresh auth token before any operation
			fileFrontGroup.POST("/upload", checkAuthMdw, refreshMdw, v1.AddPhoto)
		}

		// 后台上传文件
		fileAdminGroup := v1AdminGroup.Group("/file")
		{
			// must check auth & refresh auth token before any operation
			fileAdminGroup.POST("/upload", checkAuthMdw, refreshMdw, v1.AddPhoto)
		}
	}
}
