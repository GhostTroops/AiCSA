package blevExp

import (
	"flag"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/http"
	"github.com/blevesearch/bleve/v2/index/scorch"
	"github.com/blevesearch/bleve/v2/search"
	"github.com/blevesearch/bleve/v2/search/query"
	indexxx "github.com/blevesearch/bleve_index_api"
	util "github.com/hktalent/go-utils"
	"github.com/shomali11/util/xhashes"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

const (
	LBS    = "LBS"       // 定位
	NetCat = "reverseNc" // reverse shell cmd for nc
)

var (
	DefaulIndexName = "osint"
	Nodo            = func() {}
	SaveThread      chan struct{}
	AllIndex        = strings.Split("osint,cnnvd", ",")
	DataDir         = flag.String("DataDir", "data", "data directory")
)

type IndexData struct {
	Index string
	Id    string
	Doc   interface{}
	FnCbk func()
	FnEnd func()
}
type QueryResult struct {
	ID     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

// 初始化库
func init() {
	util.RegInitFunc(func() {
		//CreateIndex4Name(LBS, NetCat)
		SaveThread = make(chan struct{}, util.GetValAsInt("thread", 16))
	})
}

/*
	if x := pkg.Query4Key("sploitus", "RCE"); nil != x {
		log.Println(x)
	}

IncludeLocations 用于设置是否包含位置信息在搜索结果中。如果设置为 true，搜索结果中会包含每个结果的位置信息。如果设置为 false，搜索结果中将不包含位置信息。
的Explain属性是一个布尔类型，表示是否需要为每条搜索结果返回其匹配原因。当该值为true时，搜索结果会包含一个详细的分析，说明每条结果为什么被认为与查询匹配。这对于理解查询和调试查询非常有用。
{"size":10,"from":0,"explain":true,"highlight":{},"query":{"boost":1,"query":"WebLogicFilterConfig"},"fields":["*"]}
*/
func Query4Key(szIndex string, szQuery string) *bleve.SearchResult {
	return Query4KeyFrom(szIndex, szQuery, 0, 10)
}

func Query4KeyFrom(szIndex string, szQuery string, from, size uint64) *bleve.SearchResult {
	var m1 = &bleve.SearchRequest{}
	if nil == util.Json.Unmarshal([]byte(`{"size":10,"from":0,"explain":false,"highlight":{},"query":{"boost":1,"query":"\"poc\""},"fields":["*"]}`), &m1) {
		if o, ok := m1.Query.(*query.QueryStringQuery); ok {
			o.Query = szQuery
			m1.From = int(from)
			m1.Size = int(size)
			doc := Query4SearchRequest(szIndex, m1)
			if nil == doc {
				return nil
			}
			if 0 < doc.Total {
				return doc
			}
		}
	}
	return nil
}

// 删除匹配条件，且回调返回true的数据
/*
blevExp.Delete4Query(blevExp.DefaulIndexName, "type:uncover", func(m *search.DocumentMatch) bool {
	return -1 < strings.Index(m.ID, "float") || -1 == strings.Index(m.ID, "uncover")
})
go blevExp.Delete4Query(blevExp.DefaulIndexName, "tools:nuclei", func(match *search.DocumentMatch) bool {
				return true
			})
*/
func Delete4Query(szIndex string, szQuery string, fnCbk func(match *search.DocumentMatch) bool) {
	var nPs uint64 = 1000
	var nPs1 uint64 = 0
	var nPush = 10000
	var c01 = make(chan string, nPush+5)
	var fnEnd = func() {
		if nPush <= len(c01) {
			DelIndexDoc4Batch(DefaulIndexName, c01)
		}
	}
	defer fnEnd()
	if o1 := Query4KeyFrom(szIndex, szQuery, 0, nPs); nil != o1 {
		n1 := o1.Total
		for 0 < n1 {
			for _, x := range o1.Hits {
				if fnCbk(x) {
					c01 <- x.ID
					//DelIndexDoc(DefaulIndexName, x.ID)
				} else {
					nPs1++
				}
			}
			nPs1 += uint64(len(o1.Hits))
			if nPush <= len(c01) {
				DelIndexDoc4Batch(DefaulIndexName, c01)
				nPs1 -= uint64(len(c01))
			}
			if nPs1+nPs > n1 {
				break
			}
			o1 = Query4KeyFrom(szIndex, szQuery, nPs1, nPs)
		}
	}
}

func CreateDocChan(n int) chan *IndexData {
	return make(chan *IndexData, n+50)
}

// 从szIndex 移动 到 szIndexd
func Move2(szSrc, szDes, szQuery string, fnCbk func(*search.DocumentMatch) bool) {
	var nPs uint64 = 10000
	var nPs1 uint64 = 0
	if o1 := Query4KeyFrom(szSrc, szQuery, 0, nPs); nil != o1 {
		n1 := o1.Total
		docs1 := CreateDocChan(int(nPs + 50))
		var ddocID = make(chan string, nPs+50)
		for 0 < n1 {
			for _, x := range o1.Hits {
				if fnCbk(x) {
					docs1 <- &IndexData{
						Index: szDes,
						Id:    x.ID,
						Doc:   x.Fields, FnCbk: Nodo, FnEnd: Nodo,
					}
					ddocID <- x.ID
					nPs1--
				}
			}
			SaveIndexDoc4Batch(szDes, docs1, Nodo, Nodo)
			DelIndexDoc4Batch(szSrc, ddocID)
			if int(nPs) > len(o1.Hits) { // 当前页不足 page size
				return
			}
			nPs1 += nPs
			n1 -= nPs
			o1 = Query4KeyFrom(szSrc, szQuery, nPs1, nPs)
		}
	}
}

// 从szIndexs拷贝到szIndexd
func Copy2(szIndexs, szIndexd string) {
	szK := "_all:*"
	if o1 := Query4KeyFrom(szIndexs, szK, 0, 10000); nil != o1 {
		n1 := len(o1.Hits)
		var nPos uint64 = 0
		var docs = make(chan *IndexData, 5005)
		for 0 < n1 {
			for _, x := range o1.Hits {
				nPos++
				if 0 < len(x.Fields) {
					docs <- &IndexData{
						Index: szIndexd,
						Id:    x.ID,
						Doc:   x.Fields,
						FnCbk: func() {},
						FnEnd: func() {},
					}
					n3 := len(docs)
					if 5000 <= n3 {
						SaveIndexDoc4Batch(szIndexd, docs, func() {
							fmt.Printf("copy %s to %s %5d \r", szIndexs, szIndexd, n3)
						}, func() {})
					}
					//SaveIndexDoc(szIndexd, x.ID, x.Fields, func() {
					//	//DelIndexDoc(szIndexs, x.ID)
					//	n1--
					//}, func() {})
				}
			}
			n3 := len(docs)
			if 0 <= n3 {
				SaveIndexDoc4Batch(szIndexd, docs, func() {
					fmt.Printf("copy %s to %s %5d \r", szIndexs, szIndexd, n3)
				}, func() {})
			}
			o1 = Query4KeyFrom(szIndexs, szK, nPos, 10000)
			if nil == o1 {
				break
			}
			n1 = len(o1.Hits)
		}
	}
}

// form index copy to other index
func Copy2_old(szIndexs, szIndexd string) {
	idx := http.IndexByName(szIndexs)
	i, err := idx.Advanced()
	if err != nil {
		fmt.Printf("error getting index: %v", err)
		return
	}
	defer i.Close()
	r, err := i.Reader()
	if err != nil {
		fmt.Printf("error getting index reader: %v", err)
		return
	}
	defer r.Close()

	outChan, err := r.DocIDReaderAll()
	if nil != err {
		fmt.Printf("error DocIDReaderAll: %v", err)
		return
	}
	outC, ok := outChan.(*scorch.IndexSnapshotDocIDReader)
	if !ok {
		outC.Close()
		return
	}
	defer outC.Close()
	var doc = make(chan *IndexData, 5002)
	var fnSaveOk = func(n int) func() {
		return func() {
			log.Printf("move %s to %s,count %d is ok\n", szIndexs, szIndexd, n)
		}
	}
	for {
		if row, err := outC.Next(); nil == err {
			szId, err := r.ExternalID(row)
			if nil != err {
				log.Println("can not get doc", err)
				break
			}
			if oDoc := GetDoc(szIndexs, szId); nil != oDoc && 0 < len(oDoc.Fields) {
				//if oDoc, err := r.Document(szId); nil == err && nil != oDoc {
				doc <- &IndexData{
					Index: szIndexd,
					Id:    szId,
					Doc:   oDoc.Fields,
					FnCbk: func() {}, FnEnd: func() {},
				}
				n3 := len(doc)
				fmt.Printf("%10d cur doc \r", n3)
				if 0 < n3 && 0 == n3%5000 {
					SaveIndexDoc4Batch(szIndexd, doc, fnSaveOk(n3), func() {})
				}
				continue
			}
		}
		break
	}
	n3 := len(doc)
	if 0 < n3 {
		SaveIndexDoc4Batch(szIndexd, doc, fnSaveOk(n3), func() {})
	}
}

// Explain属性是一个布尔类型，表示是否需要为每条搜索结果返回其匹配原因。当该值为true时，搜索结果会包含一个详细的分析，说明每条结果为什么被认为与查询匹配。这对于理解查询和调试查询非常有用。
// IncludeLocations 用于设置是否包含位置信息在搜索结果中。如果设置为 true，搜索结果中会包含每个结果的位置信息。如果设置为 false，搜索结果中将不包含位置信息。
func Query4SearchRequest(szIndex string, searchRequest *bleve.SearchRequest) *bleve.SearchResult {
	index := http.IndexByName(szIndex)
	if index == nil {
		log.Printf("no such index '%s' \n", szIndex)
		return nil
	}
	// validate the query
	if srqv, ok := searchRequest.Query.(query.ValidatableQuery); ok {
		err := srqv.Validate()
		if err != nil {
			log.Printf("error validating query: %v", err)
			return nil
		}
	}
	// execute the query
	searchResponse, err := index.Search(searchRequest)
	if err != nil {
		log.Printf("error validating query: %v", err)
		return nil
	}
	if 0 == searchResponse.Total {
		return nil
	}
	return searchResponse
}

func DelIndexDoc4Batch(szIndex string, docID chan string) bool {
	index := http.IndexByName(szIndex)
	if index == nil {
		log.Printf("no such index '%s' \n", szIndex)
		return false
	}
	badd := index.NewBatch()
	n := len(docID)
	if 0 < n {
		for i := 0; i < n; i++ {
			badd.Delete(<-docID)
		}
	}
	if err := index.Batch(badd); nil != err {
		log.Println("DelIndexDoc4Batch index.Batch ", err)
		return false
	} else {
		log.Println("DelIndexDoc4Batch delete ok: ", n)
	}
	return true
}

func DelIndexDoc(szIndex, docID string) bool {
	index := http.IndexByName(szIndex)
	if index == nil {
		fmt.Sprintf("no such index '%s' \n", szIndex)
		return false
	}
	err := index.Delete(docID)
	if err != nil {
		fmt.Sprintf("error deleting document '%s': %v", docID, err)
		return false
	}
	return true
}

var cBath = make(chan struct{}, 16)

type DocBatch struct {
	MaxBatch int `json:"max_batch"`
	docs     chan *IndexData
	fnCbk    func()
	fnEnd    func()
	Index    string
	lock     *sync.Mutex
}

func CreateDocBatch(szIndex string, nMax int, fnCbk func(), fnEnd func()) *DocBatch {
	return &DocBatch{MaxBatch: nMax, lock: &sync.Mutex{}, Index: szIndex, docs: make(chan *IndexData, nMax), fnCbk: fnCbk, fnEnd: fnEnd}
}
func (r *DocBatch) Len() int {
	return len(r.docs)
}
func (r *DocBatch) Submit(n int) {
	r.lock.Lock()
	defer r.lock.Unlock()
	n1 := len(r.docs)
	if 0 < n1 && n <= n1 {
		dox := r.docs
		SaveIndexDoc4Batch(r.Index, dox, r.fnCbk, func() {
			close(dox)
			log.Println("close dox for SaveIndexDoc4Batch")
		})
		r.docs = make(chan *IndexData, r.MaxBatch)
		if util.GetValAsBool("devDebug") {
			debug.FreeOSMemory()
		}
	}
}
func (r *DocBatch) Push(o interface{}) *DocBatch {
	if x, ok := o.(*IndexData); ok {
		r.docs <- x
	} else {
		szId := fmt.Sprintf("%v", util.GetJson4Query(o, "id"))
		if szId == "" {
			szId = util.GetSha1(o)
		}
		r.docs <- &IndexData{
			Index: r.Index,
			Id:    szId,
			Doc:   o,
			FnCbk: r.fnCbk,
			FnEnd: r.fnEnd,
		}
	}
	r.Submit(r.MaxBatch)
	return r
}

// 批量保存,并发限制16
func SaveIndexDoc4Batch(szIndex string, doc chan *IndexData, fnCbk func(), fnEnd func()) {
	log.Println("SaveIndexDoc4Batch len", len(cBath))
	cBath <- struct{}{}
	go func() {
		defer func() {
			if nil != fnEnd {
				fnEnd()
			}
			<-cBath
		}()
		index := http.IndexByName(szIndex)
		if index == nil || nil == doc {
			fmt.Sprintf("no such index '%s' \n", szIndex)
			return
		}
		badd := index.NewBatch()
		n := len(doc)
		if 0 < n {
			for i := 0; i < n; i++ {
				d1 := <-doc
				badd.Index(d1.Id, d1.Doc)
			}
		}
		if err := index.Batch(badd); nil != err {
			log.Println("index.Batch ", err)
		} else {
			fnCbk()
			log.Printf("save ok, %d", n)
		}
	}()
}

// 并发 默认 16，thread 配置,fnCbk 在成功保存
func SaveIndexDoc(szIndex, docID string, doc interface{}, fnCbk func(), fnEnd func()) {
	fmt.Printf("SaveThread len %d\r", len(SaveThread))
	SaveThread <- struct{}{}
	util.DoSyncFunc(func() {
		defer func() {
			if nil != fnEnd {
				fnEnd()
			}
			<-SaveThread
		}()
		index := http.IndexByName(szIndex)
		if index == nil || nil == doc {
			fmt.Sprintf("no such index '%s' \n", szIndex)
			return
		}
		if err := index.Index(docID, doc); err != nil {
			fmt.Printf("error indexing document '%s': %v\n", docID, err)
		} else if nil != fnCbk {
			fnCbk()
		}
	})
}

func GetUrlId(u string) string {
	return SHA1(u)
}

func SHA1(u string) string {
	return xhashes.SHA1(u)
}

// 在所有索引中查询 url ，内部会将url进行 sha1编码作为id进行查询
func GetAllIndexDoc4Url(szUrl string) interface{} {
	docId := SHA1(strings.TrimSpace(szUrl))
	return GetAllIndexDoc(docId)
}

// 在所有索引中查询 id
func GetAllIndexDoc(docId string) interface{} {
	for _, x := range AllIndex {
		if o := GetDoc(x, docId); nil != o {
			return o
		}
	}
	return nil
}

func CreateIndex4Name(indexName ...string) {
	s1 := `{"default_mapping":{"enabled":true,"display_order":"0"},"type_field":"_type","default_type":"_default","default_analyzer":"standard","default_datetime_parser":"dateTimeOptional","default_field":"_all","byte_array_converter":"json","store_dynamic":true,"index_dynamic":true}`
	for _, x := range indexName {
		CreateIndex(x, s1)
	}
}

// 创建索引
/*
curl 'https://127.0.0.1:8081/api/classInfo' \
-X 'PUT' \
-H 'Accept: application/json, text/plain' \
-H 'Content-Type: application/json;charset=utf-8' \
-H 'Origin: https://127.0.0.1:8081' \
-H 'Cookie: 51pwn4ssh=s%3AEIw15tP5XD_wvn7cnc5E9I0OD8LbDODZ.ffjTn31urO8WPINViwiLu62SpuywsXZltTj4kMA7PGI; ss=MTY2NzA1MDc4MHxEdi1CQkFFQ180SUFBUkFCRUFBQV9nRTNfNElBQVFaemRISnBibWNNQndBRmRHOXJaVzRHYzNSeWFXNW5EUDRCR0FELUFSUmxlVXBvWWtkamFVOXBTa2xWZWswMFRrTkpjMGx1VWpWalEwazJTV3R3V0ZaRFNqa3VaWGxLZWxwWFZtdEphbTlwVGxSVk1FNVVhM2xOUkZFMFRrUm5OVTFFUVhkT1JGVXpUbmxKYzBsdFJqRmFRMGsyU1cxR2EySlhiSFZKYVhkcFdsaG9kMGxxYjNoT2Fsa3pUVlJOTTAxVVozZE1RMHB4WkVkcmFVOXBTWGhKYVhkcFlWZEdNRWxxYjNoT2Fsa3pUVVJWZDA1NlozZE1RMHB3WXpOTmFVOXBTVEZOV0VJelltazFhbUl5TUdsTVEwcDZaRmRKYVU5cFNURk5XRUl6WW10Qk1VMVlRak5pYVRWcVlqSXdhV1pSTGtkdVlsbzBObUpVZGtvNWNGYzJiVEEwUlVvelUwcDRWVEZEWVZsSlRGUmZaVEJCZGtaR1psUk5RV04wWHpoZk5rOW5VVk5NUmxjMmNVRm5TMHBCWXpNPXyHSm4pk58Q6_JDEPGEOjRB1PzKrMn3GOeChY_Zs17IRA==; mysession=MTY2NzA0OTkyNnxEdi1CQkFFQ180SUFBUkFCRUFBQV9nRTNfNElBQVFaemRISnBibWNNQndBRmRHOXJaVzRHYzNSeWFXNW5EUDRCR0FELUFSUmxlVXBvWWtkamFVOXBTa2xWZWswMFRrTkpjMGx1VWpWalEwazJTV3R3V0ZaRFNqa3VaWGxLZWxwWFZtdEphbTlwVDBSSk1VMUVUVEpOVkZrMVQxUk5OVTlFUlhwUFJGRjNUbE5KYzBsdFJqRmFRMGsyU1cxR2EySlhiSFZKYVhkcFdsaG9kMGxxYjNoT2Fsa3pUVlJOTWsxNlNUSk1RMHB4WkVkcmFVOXBTWGhKYVhkcFlWZEdNRWxxYjNoT2Fsa3pUVVJSTlU5VVNUSk1RMHB3WXpOTmFVOXBTVEZOV0VJelltazFhbUl5TUdsTVEwcDZaRmRKYVU5cFNURk5XRUl6WW10Qk1VMVlRak5pYVRWcVlqSXdhV1pSTG01NGVubG5XRXBETVVseWJWWTVSRmN3YUZkc1Yzb3dhVFk1VERkWmIxQkJTbnBUY1ROR1VUQmFSM2RwV0RCMWIwdElRVE54TkZCUk1ERlBabXhtU2pJPXwPXOhkJ6Qdcuvofa8lFQxCof8fKUP0pfJ2OpoIexbGiw==' \
-H 'Content-Length: 273' \
-H 'Host: 127.0.0.1:8081' \
-H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15' \
-H 'Referer: https://127.0.0.1:8081/indexes/_new' \
--data-binary '{"default_mapping":{"enabled":true,"display_order":"0"},"type_field":"_type","default_type":"_default","default_analyzer":"standard","default_datetime_parser":"dateTimeOptional","default_field":"_all","byte_array_converter":"json","store_dynamic":true,"index_dynamic":true}'
*/
func CreateIndex(indexName string, opt interface{}) {
	indexMapping := bleve.NewIndexMapping()
	var err error
	var requestBody []byte
	if s, ok := opt.(string); ok {
		requestBody = []byte(s)
	} else {
		requestBody, err = util.Json.Marshal(opt)
		if nil != err {
			fmt.Sprintf("error parsing index mapping: %v", err)
			return
		}
	}
	err = util.Json.Unmarshal(requestBody, &indexMapping)
	if err != nil {
		fmt.Sprintf("error parsing index mapping: %v", err)
		return
	}
	if newIndex, err := bleve.New(*DataDir+string(os.PathSeparator)+indexName, indexMapping); nil != err {
		fmt.Sprintf("error parsing index mapping: %v", err)
		return
	} else {
		newIndex.SetName(indexName)
		http.RegisterIndexName(indexName, newIndex)
	}
}

// 查询指定索引、指定id的 doc
func GetDoc(szIndex, docID string) *QueryResult {
	index := http.IndexByName(szIndex)
	if index == nil {
		fmt.Sprintf("no such index '%s' \n", szIndex)
		return nil
	}

	doc, err := index.Document(docID)
	if err != nil {
		fmt.Sprintf("error deleting document '%s': %v", docID, err)
		return nil
	}
	if nil == doc {
		return nil
	}
	rv := QueryResult{
		ID:     docID,
		Fields: map[string]interface{}{},
	}

	doc.VisitFields(func(field indexxx.Field) {
		var newval interface{}
		switch field := field.(type) {
		case indexxx.TextField:
			newval = field.Text()
		case indexxx.NumericField:
			n, err := field.Number()
			if err == nil {
				newval = n
			}
		case indexxx.DateTimeField:
			d, _, err := field.DateTime()
			if err == nil {
				newval = d.Format(time.RFC3339Nano)
			}
		}
		existing, existed := rv.Fields[field.Name()]
		if existed {
			switch existing := existing.(type) {
			case []interface{}:
				rv.Fields[field.Name()] = append(existing, newval)
			case interface{}:
				arr := make([]interface{}, 2)
				arr[0] = existing
				arr[1] = newval
				rv.Fields[field.Name()] = arr
			}
		} else {
			rv.Fields[field.Name()] = newval
		}
	})
	return &rv
}
