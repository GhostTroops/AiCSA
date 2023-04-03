package pkg

import (
	"embed"
	"github.com/gin-gonic/gin"
	util "github.com/hktalent/go-utils"
	"github.com/hktalent/go-utils/bigdb/blevExp"
	"net/http"
)

// curl -v -sk https://127.0.0.1:8080/api/AiCSA/b1f4bdaab215468b0ca7821236eda2e3d4e7f48d -o- |jq
func CreateHttp3Server(static1 embed.FS) {
	var router *gin.Engine
	gin.SetMode(gin.ReleaseMode)
	router = gin.New()
	router.Use(gin.Recovery())
	router.UseH2C = true
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/indexes")
	})

	blevExp.DoIndexDb(router, static1, blevExp.DataDir, func(r001 *gin.Engine) {})

	RunHttp3(":"+util.GetVal("HttpPort"), router)
}
