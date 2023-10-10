package router

import (
	"goweb/app/account/api"

	"github.com/gin-gonic/gin"
)

func InitAccountRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	{
		userRouter.GET("/:userId", api.UserInfo)
		userRouter.POST("/create", api.UserCreate)
		userRouter.POST("/:userId/update", api.UserUpdate)
		userRouter.GET("/:userId/delete", api.UserDelete)
		userRouter.GET("/list", api.UserList)
		userRouter.POST("/search", api.UserSearch)

		//TODO
		userRouter.POST("/register", api.Todo)
		userRouter.POST("/login", api.Todo)
		userRouter.POST("/logout", api.Todo)
		userRouter.POST("/pwd/modify", api.Todo)
		userRouter.POST("/pwd/reset", api.Todo)

		//M2M User2Group
		userRouter.GET("/groupsList", api.GroupsListByUser)
		userRouter.POST("/groupsAdd", api.GroupsAddToUser)
		userRouter.POST("/groupsDelete", api.GroupsDeleteFromUser)
	}

	groupRouter := Router.Group("group")
	{
		groupRouter.GET("/:groupId", api.GroupInfo)
		groupRouter.POST("/create", api.GroupCreate)
		groupRouter.POST("/:groupId/update", api.GroupUpdate)
		groupRouter.GET("/:groupId/delete", api.GroupDelete)
		groupRouter.GET("/list", api.GroupList)
		groupRouter.POST("/search", api.GroupSearch)

		//M2M User2Group
		groupRouter.GET("/usersList", api.UsersListByGroup)
		groupRouter.POST("/usersAdd", api.UsersAddToGroup)
		groupRouter.POST("/usersDelete", api.UsersDeleteFromGroup)
	}

}
