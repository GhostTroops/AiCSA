package main

import (
	"embed"
	"github.com/hktalent/AiCSA/pkg"
	util "github.com/hktalent/go-utils"
	"github.com/hktalent/go-utils/bigdb/blevExp"
	"os"
)

// 搜索引擎
//
//go:embed static/*
var static1 embed.FS

func main() {
	util.DoInitAll()
	os.Args = []string{"", "src/a611b13e2af8684075349c35ed6a4147", "src/cd0a7d64c4b34ae5d9ffab8b94c9db7f"}
	// control + c 退出时做些什么
	//util.NewExit().RegClose(func() error {
	//	return nil
	//})

	// init index db
	blevExp.CreateIndex4Name(pkg.IndexName)

	// web server
	pkg.CreateHttp3Server(static1)

	//blevExp.InitIndexDb(func(s string) bool {
	//	return false
	//})
	defer pkg.Limiter.Stop()

	for _, x := range os.Args[1:] {
		pkg.WalkDir(x)
	}

	util.Wg.Wait()
	util.CloseAll()
}
