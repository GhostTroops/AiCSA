package main

import (
	"github.com/hktalent/AiCSA/pkg"
	util "github.com/hktalent/go-utils"
	"os"
)

func main() {
	util.DoInitAll()
	// control + c 退出
	util.NewExit().RegClose(func() error {
		//pkg.GptApi.
		return nil
	})
	pkg.WalkDir(os.Args[1])
	util.Wg.Wait()
	util.CloseAll()
}
