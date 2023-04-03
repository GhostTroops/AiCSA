package pkg

import (
	"fmt"
	util "github.com/hktalent/go-utils"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 记录该文件，便于下一次重复再次尝试
func ReDoCode(s, data, szHashCode string) {
	util.PutAny[string](szHashCode+"_|_"+s, data)
}

var (
	regM1 = regexp.MustCompile("\\s*\n\\s*")
	c001  = make(chan struct{}, 20) // 并发控制
)

// 拆分大小
const MaxSplit = 3500

/*
DoOneJava
s 文件名
data 文件内容
szHashCode 项目hash
*/
func DoOneJava(s, data, szHashCode string) {
	data = regM1.ReplaceAllString(data, "")
	if a, ok := QueryIndexData(s, szHashCode); ok {
		log.Println("文件"+s+",已经处理过:", len(a))
		//for _, x := range a {
		//	fmt.Println(x.SecInfo)
		//}
	} else {
		Limiter.Take()
		c001 <- struct{}{}
		util.Wg.Add(1)
		util.DefaultPool.Submit(func() {
			defer util.Wg.Done()
			defer func() {
				<-c001
			}()
			szOldData := data
			var a = []string{data}
			if MaxSplit < len(data) {
				n1 := len(data) / MaxSplit
				if 0 < len(data)%MaxSplit {
					n1++
				}
				a = make([]string, n1)
				x0 := 0
				for {
					if MaxSplit <= len(data) {
						a[x0] = data[0:MaxSplit]
					} else {
						a[x0] = data
						break
					}
					data = data[MaxSplit:]
					if 0 == len(data) || MaxSplit > len(a[x0]) {
						break
					}
					x0++
				}
			}
			var aR []string
			sz11 := fmt.Sprintf("项目:%s 文件:%s,", szHashCode, s)
			isOne := 1 == len(a)
			if isOne {
				if s1, err := GptNew(fmt.Sprintf("%s,代码:%s", sz11, data+"\n"+Prefix)); nil == err {
					if 0 < len(s1) {
						aR = append(aR, s1)
						fmt.Println(s1)
					}
				} else { // 记录该文件，便于下一次重复再次尝试
					ReDoCode(s, szOldData, szHashCode)
				}
			} else {
				for i, x1 := range a {
					if s1, err := GptNew(fmt.Sprintf("先不分析等我发送完该文件内容,%s,第%d段代码:%s", sz11, i+1, x1)); nil == err {
						if 0 < len(s1) {
							//aR = append(aR, s1)
							fmt.Println(s1)
						}
					} else {
						log.Println(err)
					}
				}

				if s1, err := GptNew(fmt.Sprintf(Prefix, "前面多个拆分发送的"+sz11)); nil == err {
					if 0 < len(s1) {
						aR = append(aR, s1)
						fmt.Println(s1)
					}
				} else { // 记录该文件，便于下一次重复再次尝试
					ReDoCode(s, szOldData, szHashCode)
				}
			}
			if 0 < len(aR) {
				SaveIndexData(s, strings.Join(aR, "\n"), szHashCode)
			}
		})
	}
}

func getLastDir(s string) string {
	// 使用 filepath.Split 函数获取到最后一层非空目录
	_, lastDir := filepath.Split(s)
	// 使用 Strings.TrimRight 函数去除可能存在的目录分隔符号
	lastDir = strings.TrimRight(lastDir, string(filepath.Separator))
	return lastDir
}

func WalkDir(s string) {
	var szHashCode = getLastDir(s)
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
			DoOneJava(path, string(data), szHashCode)
		}
		return nil
	})
	if nil != err {
		log.Println(err)
	}
}
