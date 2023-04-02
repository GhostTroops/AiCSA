package pkg

import (
	"fmt"
	util "github.com/hktalent/go-utils"
	"github.com/hktalent/go-utils/bigdb/blevExp"
	"reflect"
	"time"
)

const IndexName = "AiCSA"

type IndexData struct {
	JarHash      string    `json:"jar_hash"`
	JavaCodePath string    `json:"java_code_path"`
	SecInfo      string    `json:"sec_info"`
	SecVerify    string    `json:"sec_verify"`
	Date         time.Time `json:"date"`
}

/*
SaveIndexData

	保存结果到数据库
	没有显著的安全隐患,不存在明显的安全风险
	@jarHash: jar包的hash值
	@javaCodePath: java代码的路径
	@secInfo: 安全信息
*/
func SaveIndexData(javaCodePath, secInfo, jarHash string) {
	var x = IndexData{jarHash, javaCodePath, secInfo, "", time.Now()}
	blevExp.SaveIndexDoc(IndexName, util.GetSha1(x), x, blevExp.Nodo, nil)
}

// copy
func Copy2IndexData(m *map[string]interface{}) *IndexData {
	// 定义一个 IndexData 类型的结构体变量
	var indexData IndexData
	// 遍历 map 中的所有键值对，并通过反射机制将值转化为 IndexData 中对应字段的类型
	for k, v := range *m {
		if field, ok := reflect.TypeOf(indexData).FieldByName(k); ok {
			fv := reflect.ValueOf(&indexData).Elem().FieldByName(field.Name)
			if fv.Type().Kind() == reflect.String {
				fv.SetString(fmt.Sprintf("%v", v))
			} else if fv.Type().Kind() == reflect.Struct && fv.Type().Name() == "Time" {
				if t, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%v", v)); err == nil {
					fv.Set(reflect.ValueOf(t))
				}
			}
		}
	}
	return &indexData
}

// 查询数据库
func QueryIndexData(javaCodePath, jarHash string) ([]*IndexData, bool) {
	if x := blevExp.Query4Key(IndexName, fmt.Sprintf("jar_hash:\"%s\" +java_code_path:\"%s\"", jarHash, javaCodePath)); nil != x {
		if 0 < x.Total {
			var a []*IndexData
			for x1 := range x.Hits {
				a = append(a, Copy2IndexData(&x.Hits[x1].Fields))
			}
			return a, true
		}
	}
	return nil, false
}
