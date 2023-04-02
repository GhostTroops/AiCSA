package blevExp

import (
	"embed"
	"flag"
	"github.com/blevesearch/bleve/v2"
	bleveHttp "github.com/blevesearch/bleve/v2/http"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	bleveMappingUI "github.com/hktalent/bleve-mapping-ui"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// var DataDir = flag.String("DataDir", "/Volumes/Home/data", "data directory")
var StaticEtag = flag.String("StaticEtag", "", "optional static etag value.")
var StaticPath = flag.String("static", "",
	"optional path to static directory for web resources")
var StaticBleveMappingPath = flag.String("staticBleveMapping", "./static/", "optional path to static-bleve-mapping directory for web resources")

// 初始化索引库
func InitIndexDb(isSkip func(s string) bool) {
	// walk the data dir and register index names
	dirEntries, err := ioutil.ReadDir(*DataDir)
	if err != nil {
		log.Fatalf("error reading data dir: %v", err)
	}
	var wg1 sync.WaitGroup
	for _, dirInfo := range dirEntries {
		indexPath := *DataDir + string(os.PathSeparator) + dirInfo.Name()

		a1 := strings.Split(indexPath, string(os.PathSeparator))
		if a1[len(a1)-1] == "cache" || nil != isSkip && isSkip(a1[len(a1)-1]) {
			continue
		}
		// skip single files in data dir since a valid index is a directory that
		// contains multiple files
		if !dirInfo.IsDir() {
			log.Printf("not registering %s, skipping", indexPath)
			continue
		}
		//wg1.Add(1)
		func(s1 string, di string) {
			//defer wg1.Done()
			i, err := bleve.Open(s1)
			if err != nil {
				log.Printf("error opening index %s: %v", s1, err)
			} else {
				log.Printf("registered index: %s", di)
				bleveHttp.RegisterIndexName(di, i)
				// set correct name in stats
				i.SetName(di)
			}
		}(indexPath, dirInfo.Name())
	}
	wg1.Wait()
}

func GetStaticBleveMapping() http.Handler {
	// default to bindata for static-bleve-mapping resources.
	staticBleveMapping := http.FileServer(bleveMappingUI.AssetFS())
	if *StaticBleveMappingPath != "" {
		fi, err := os.Stat(*StaticBleveMappingPath)
		if err == nil && fi.IsDir() {
			log.Printf("using static-bleve-mapping resources from %s", *StaticBleveMappingPath)
			staticBleveMapping = http.FileServer(http.Dir(*StaticBleveMappingPath))
		}
	}
	return staticBleveMapping
}

func GetStaticHandler(staticDir embed.FS) gin.HandlerFunc {
	static := static.Serve("/static", EmbedFolder(staticDir, "static"))
	//static := http.FileServer(AssetFS())
	//if *StaticPath != "" {
	//	fi, err := os.Stat(*StaticPath)
	//	if err == nil && fi.IsDir() {
	//		log.Printf("using static resources from %s", *StaticPath)
	//		static = http.FileServer(http.Dir(*StaticPath))
	//	}
	//}
	return static
}

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(prefix string, path string) bool {
	if strings.HasPrefix(path, prefix) {
		path = path[len(prefix):]
	}
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}
func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	fsys, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fsys),
	}
}
