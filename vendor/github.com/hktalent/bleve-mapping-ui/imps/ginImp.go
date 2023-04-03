package imps

import (
	"github.com/gin-gonic/gin"
	xx "github.com/hktalent/bleve-mapping-ui"
	"net/http"
)

type GinImp struct {
	router *gin.Engine
	path   string
	f      func(http.ResponseWriter, *http.Request)
}

func NewGinImp(r *gin.Engine) *GinImp {
	return &GinImp{router: r}
}
func (r *GinImp) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) xx.MethodsFace {
	r.path = path
	r.f = f
	var k xx.MethodsFace = r
	return k
}

func (r *GinImp) Methods(methods ...string) {
	fnCbk := func(cc *gin.Context) {
		r.f(cc.Writer, cc.Request)
	}
	if 0 < len(methods) {
		r.router.Handle(methods[0], r.path, fnCbk)
	}
}
