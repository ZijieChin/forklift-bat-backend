package router

import (
	v1 "forklift-bat-backend/controller/v1"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	GroupV1 := r.Group("/v1")
	{
		GroupV1.GET("/userauth", v1.UserAuth)
		GroupV1.GET("/jwt", v1.CheckJWT)
		GroupV1.GET("/forklifts", v1.GetForklift)
		GroupV1.POST("/forklift", v1.AddForklift)
		GroupV1.POST("/warehouse", v1.AddWarehouse)
		GroupV1.POST("/forkcat", v1.AddForkCat)
		GroupV1.POST("/battery", v1.AddBattery)
		GroupV1.POST("/switch", v1.SwitchBattery) // forklift SN, warehouse, battery SN, switch(on/off), userid
	}
	return r
}
