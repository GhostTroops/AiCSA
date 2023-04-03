package blevExp

import (
	"embed"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/gin-gonic/gin"
	bleveMappingUI "github.com/hktalent/bleve-mapping-ui"
	util "github.com/hktalent/go-utils"
	"log"
	"net/http"
	"strings"

	bleveHttp "github.com/blevesearch/bleve/v2/http"
	"github.com/hktalent/bleve-mapping-ui/imps"
	// import general purpose configuration
	_ "github.com/blevesearch/bleve/v2/config"
)

// 初始化搜索引擎
func DoIndexDb(router *gin.Engine, staticDir embed.FS, DataDir *string, args ...func(*gin.Engine)) {
	isDev := util.GetValAsBool("devDebug")
	var fnInitDb = func() {
		InitIndexDb(func(s string) bool {
			// return isDev && (s == "fileIndex" || s == "sgk")
			return isDev && (s == "fileIndex") || ".DS_Store" == s
		})
	}
	//if isDev {
	fnInitDb()
	//}
	// https://127.0.0.1:8081/initDb
	router.GET("/initDb", func(c *gin.Context) {
		go fnInitDb()
	})
	if nil == router {
		return
	}
	// 必须在上面初始化索引后
	defer func() {
		for i, fnCbk := range args {
			log.Println("必须在上面初始化索引后 func: ", i)
			util.DoSyncFunc(func() {
				fnCbk(router)
			})
		}
	}()

	//router := gin.Default()
	router.UseH2C = true

	staticBleveMapping := GetStaticBleveMapping()
	oStatic := http.Dir("./static")
	router.StaticFS("/static", oStatic)
	router.StaticFS("/template", oStatic)
	router.StaticFS("/overview", oStatic)
	router.StaticFS("/search", oStatic)
	router.StaticFS("/indexes", oStatic)
	router.StaticFS("/analysis", oStatic)
	router.StaticFS("/monitor", oStatic)

	// 像个静态目录的处理
	router.Use(func(c *gin.Context) {
		k1 := "/static-bleve-mapping/"
		switch {
		case strings.HasPrefix(c.Request.URL.Path, k1):
			s1 := c.Request.URL.Path
			c.Request.URL.Path = c.Request.URL.Path[len(k1):]
			if nil != staticBleveMapping {
				staticBleveMapping.ServeHTTP(c.Writer, c.Request)
			}
			c.Request.URL.Path = s1

			//staticDir1 := static.Serve(k1, EmbedFolder(staticDir, "static-bleve-mapping"))
			//staticDir1(c)
			c.Abort()
			//http.StripPrefix(k1, MyFileHandler{staticBleveMapping}).ServeHTTP(c.Writer, c.Request)
			return
			//case strings.HasPrefix(c.Request.URL.Path, "/static/"):
			//	static.ServeHTTP(c.Writer, c.Request)
			//	return
			//c.Next()
		}
	})

	// 各种 api
	x1 := imps.NewGinImp(router)
	bleveMappingUI.RegisterHandlers(x1, "/api")

	// 创建索引
	createIndexHandler := bleveHttp.NewCreateIndexHandler(*DataDir)
	createIndexHandler.IndexNameLookup = IndexNameLookup
	router.Handle("PUT", "/api/:indexName", CreateHandle(createIndexHandler))

	getIndexHandler := bleveHttp.NewGetIndexHandler()
	getIndexHandler.IndexNameLookup = IndexNameLookup
	router.Handle("GET", "/api/:indexName", CreateHandle(getIndexHandler))

	//  先注释，避免恶意的删除
	//deleteIndexHandler := bleveHttp.NewDeleteIndexHandler(*DataDir)
	//deleteIndexHandler.IndexNameLookup = IndexNameLookup
	//router.Handle("DELETE", "/api/:indexName", CreateHandle(deleteIndexHandler))

	docIndexHandler := bleveHttp.NewDocIndexHandler("")
	docIndexHandler.IndexNameLookup = IndexNameLookup
	docIndexHandler.DocIDLookup = DocIDLookup
	router.Handle("PUT", "/api/:indexName/:docID", CreateHandle(docIndexHandler))

	docCountHandler := bleveHttp.NewDocCountHandler("")
	docCountHandler.IndexNameLookup = IndexNameLookup
	router.Handle("GET", "/api/:indexName/_count", CreateHandle(docCountHandler))

	docGetHandler := bleveHttp.NewDocGetHandler("")
	docGetHandler.IndexNameLookup = IndexNameLookup
	docGetHandler.DocIDLookup = DocIDLookup
	router.Handle("GET", "/api/:indexName/:docID", CreateHandle(docGetHandler))

	docDeleteHandler := bleveHttp.NewDocDeleteHandler("")
	docDeleteHandler.IndexNameLookup = IndexNameLookup
	docDeleteHandler.DocIDLookup = DocIDLookup
	var delDoc = CreateHandle(docDeleteHandler)
	// 支持查询条件进行删除
	router.Handle("DELETE", "/api/:indexName/:docID", func(g *gin.Context) {
		id := g.Param("docID")
		if GetDoc(DefaulIndexName, id) != nil {
			delDoc(g)
		} else {
			go Delete4Query(DefaulIndexName, id, func(match *search.DocumentMatch) bool {
				return true
			})
		}
	})

	searchHandler := bleveHttp.NewSearchHandler("")
	searchHandler.IndexNameLookup = IndexNameLookup
	router.Handle("POST", "/api/:indexName/_search", CreateHandle(searchHandler))

	listFieldsHandler := bleveHttp.NewListFieldsHandler("")
	listFieldsHandler.IndexNameLookup = IndexNameLookup
	router.Handle("GET", "/api/:indexName/_fields", CreateHandle(listFieldsHandler))

	debugHandler := bleveHttp.NewDebugDocumentHandler("")
	debugHandler.IndexNameLookup = IndexNameLookup
	debugHandler.DocIDLookup = DocIDLookup
	router.Handle("GET", "/api/:indexName/:docID/_debug", CreateHandle(debugHandler))

	aliasHandler := bleveHttp.NewAliasHandler()
	router.Handle("POST", "/api/_aliases", CreateHandle(aliasHandler))

	listIndexesHandler := bleveHttp.NewListIndexesHandler()
	router.Handle("GET", "/api", CreateHandle(listIndexesHandler))
	router.Use(func(c *gin.Context) {
		if strings.HasSuffix(c.Request.URL.Path, "/") {
			c.Request.URL.Path = c.Request.URL.Path[0 : len(c.Request.URL.Path)-1]
		}
	})
	//router.NoRoute(func(c *gin.Context) {
	//	http.RedirectHandler("/static/index.html", 302).ServeHTTP(c.Writer, c.Request)
	//	log.Println("no route ", c.Request.URL.Path)
	//})
}
