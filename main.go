package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/jaaaaason/hmblog/configer"
	"github.com/jaaaaason/hmblog/database"
	"github.com/jaaaaason/hmblog/handler"
	"github.com/jaaaaason/hmblog/logger"
)

func main() {
	// get config file's path with commandline arg
	confFilepath := flag.String("c", "", "the config file's path")

	var err error
	if *confFilepath == "" {
		// use default config file's path
		err = configer.Initialize("config.json")
	} else {
		err = configer.Initialize(*confFilepath)
	}

	if err != nil {
		fmt.Println("Can not initializes configuration")
		return
	}

	logFile := os.Stdout
	// desire log file is given
	if configer.Config.LogFile != "" {
		logFile, err = os.OpenFile(configer.Config.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println("Can not initializes log")
			return
		}
	}
	defer logFile.Close()

	logger.Initialize(logFile)

	// initialize the database's connection
	err = database.Initialize()
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
	defer database.CloseSession()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(handler.CORSMiddleware())
	registerRoute(r)

	adminRouter := r.Group("/admin")
	adminRouter.Use(handler.JWTMiddleware())
	registerAdminRoute(adminRouter)

	var addr string
	if configer.Config.Listen > 0 {
		// use given port in config file
		addr = fmt.Sprintf(":%d", configer.Config.Listen)
	} else {
		// use default port 8080
		addr = ":8080"
	}

	r.Run(addr)
}

// registerRoute registers api route
func registerRoute(r *gin.Engine) {
	r.POST("/admin/login", handler.PostLogin)

	// category
	r.GET("/categories", handler.GetCategories)
	r.GET("/categories/:id", handler.GetCategory)

	// post
	r.GET("/posts", handler.GetPosts)
	r.GET("/posts/:id", handler.GetPost)
	r.GET("/categories/:id/posts", handler.GetCategoryPosts)
}

// registerAdminRoute registers admin api route
func registerAdminRoute(r *gin.RouterGroup) {
	// admin category
	r.GET("/categories", handler.GetAdminCategories)
	r.GET("/categories/:id", handler.GetAdminCategory)
	r.POST("/categories", handler.PostCategory)
	r.PUT("/categories/:id", handler.UpdateCategory)
	r.PATCH("/categories/:id", handler.UpdateCategory)
	r.DELETE("/categories/:id", handler.DeleteCategory)

	// admin post
	r.GET("/posts", handler.GetAdminPosts)
	r.GET("/posts/:id", handler.GetAdminPost)
	r.GET("/categories/:id/posts", handler.GetAdminCategoryPosts)
	r.POST("/categories/:id/posts", handler.PostCategoryPost)
	r.POST("/posts", handler.PostPost)
	r.PUT("/posts/:id", handler.UpdatePost)
	r.PATCH("/posts/:id", handler.UpdatePost)
	r.DELETE("/posts/:id", handler.DeletePost)

	// admin user
	r.PUT("/users/:id", handler.UpdateUser)
	r.PATCH("/users/:id", handler.UpdateUser)
	r.PUT("/users/:id/password", handler.UpdateUserPassword)
}
