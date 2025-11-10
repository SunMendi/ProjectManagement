package main

import (
	"ProjectManagement/User"
	"ProjectManagement/chat"
	"ProjectManagement/database"
	"ProjectManagement/task"
	"ProjectManagement/team"
	"ProjectManagement/upload" 
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)


func main () {
	   if _, err := os.Stat("uploads"); os.IsNotExist(err) {
        os.Mkdir("uploads", 0755)
    }


	 db, err := database.ConnectDB()
	 if err!= nil {
		 log.Panic("Database Initialization Issue")
	 }
	 gin.SetMode(gin.ReleaseMode)
	 log.Println("successfully conntected to database")
	 router:= gin.Default() 
	 router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"}, // or use "*" to allow any origin
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,               // allow cookies/credentials if you need them
        MaxAge:           12 * time.Hour,     // cache preflight response for 12 hours
    }))

	 user.RegisterRoutes(router, db)
	 team.RegisterRoutes(router, db)
	 task.RegisterRoutes(router, db)
	 chat.RegisterRoutes(router, db)
	 upload.RegisterRoutes(router)

	 router.Run(":8081")
}