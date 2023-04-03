package pkg

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/hktalent/go-utils/bigdb/blevExp"
)

func CreateHttp3Server(static1 embed.FS) {
	var router *gin.Engine
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	router.Use(gin.Recovery())
	router.UseH2C = true

	blevExp.DoIndexDb(router, static1, blevExp.DataDir, func(r001 *gin.Engine) {})

	RunHttp3(":8080", router)
}
