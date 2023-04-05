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
	c001  chan struct{}
	regM2 *regexp.Regexp
)

// 拆分大小
const MaxSplit = 3500

func init() {
	util.RegInitFunc(func() {
		regM2 = regexp.MustCompile(util.GetVal("sourceExt"))
	})
}

/*
DoOneJava
s 文件名
data 文件内容
szHashCode 项目hash
*/
func DoOneJava(s, data, szHashCode string) {
	data = regM1.ReplaceAllString(data, "")
	if a, ok := QueryIndexData(s, szHashCode); ok {
		fmt.Printf("文件 %s,已经处理过:%d\r", s, len(a))
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
				n09 := len(a)
				var szFirst = ""

				for i, x1 := range a {
					if 0 == i {
						szFirst = fmt.Sprintf("接下来我将 %s 拆分为 %d 部分发给你，直到全部给你发完才开始分析", sz11, n09)
					} else {
						szFirst = ""
					}
					if s1, err := GptNew(fmt.Sprintf("%s %s,第 %d/%d 段代码:%s", szFirst, sz11, i+1, n09, x1)); nil == err {
						if 0 < len(s1) {
							//aR = append(aR, s1)
							fmt.Println(s1)
						}
					} else {
						log.Println(err)
					}
				}

				if s1, err := GptNew(fmt.Sprintf("分析前面多个拆分发送的 %s java代码存在哪些安全风险、易受到攻击的脆弱代码,如何验证、确认他们", sz11)); nil == err {
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

		if !info.IsDir() && regM2.Match([]byte(path)) {
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
