package middleware

import (
	"forklift-bat-backend/model"

	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
)

var Enforcer *casbin.Enforcer

func CasbinInit() {
	a, errGorm := gormadapter.NewAdapterByDB(model.DB)
	if errGorm != nil {
		panic("Error initializing casbin adapter")
	}
	e := casbin.NewEnforcer("./middleware/rbca_model.conf", a)
	Enforcer = e
}

func Authorize() *gin.HandlerFunc {
	return func(c *gin.Context) {
		act := c.Request.Method

	}
}
