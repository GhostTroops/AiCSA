package main

import (
	"github.com/hktalent/AiCSA/pkg"
	util "github.com/hktalent/go-utils"
	"github.com/hktalent/go-utils/bigdb/blevExp"
	"os"
)

func main() {
	util.DoInitAll()
	os.Args = []string{"", "src/a611b13e2af8684075349c35ed6a4147"}
	// control + c 退出时做些什么
	//util.NewExit().RegClose(func() error {
	//	return nil
	//})

	// init index db
	blevExp.CreateIndex4Name(pkg.IndexName)
	blevExp.InitIndexDb(func(s string) bool {
		return false
	})
	defer pkg.Limiter.Stop()

	pkg.WalkDir(os.Args[1])
	util.Wg.Wait()
	util.CloseAll()
}
