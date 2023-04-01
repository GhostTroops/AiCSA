package pkg

import (
	util "github.com/hktalent/go-utils"
	"log"
	"os"
	"path/filepath"
)

func DoOneJava(s, data string) {
	util.Wg.Add(1)
	util.DefaultPool.Submit(func() {
		defer util.Wg.Done()
		if s1, err := GptNew(data); nil != err {

		} else {

		}
	})
}

func WalkDir(s string) {
	err := filepath.Walk(s, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果该路径是一个文件且扩展名为.java，则读取该文件的内容
		if !info.IsDir() && filepath.Ext(path) == ".java" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			DoOneJava(path, string(data))
		}
		return nil
	})
	if nil != err {
		log.Println(err)
	}
}
