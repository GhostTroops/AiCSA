package pkg

import (
	util "github.com/hktalent/go-utils"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type GptApiPool struct {
	GptApi []*openai.Client
	Look   *sync.Mutex
	Pos    int
}

func (this *GptApiPool) GetGptApi() *openai.Client {
	this.Look.Lock()
	defer this.Look.Unlock()
	x := this.GptApi[this.Pos]
	this.Pos++
	if this.Pos >= len(this.GptApi) {
		this.Pos = 0
	}
	return x
}

func (this *GptApiPool) initGptApi(chatGptKey string) *openai.Client {
	this.Look.Lock()
	defer this.Look.Unlock()
	szProxy := util.GetVal("proxy")
	var GptApi *openai.Client
	if szProxy == "" {
		GptApi = openai.NewClient(chatGptKey)
	} else {
		config := openai.DefaultConfig(chatGptKey)
		proxyUrl, err := url.Parse(szProxy)
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		config.HTTPClient = &http.Client{
			Transport: transport,
		}
		GptApi = openai.NewClientWithConfig(config)
	}
	return GptApi
}

func NewGptApiPool() *GptApiPool {
	x := &GptApiPool{
		Look: new(sync.Mutex),
		Pos:  0,
	}
	chatGptKey := strings.Split(util.GetVal("api_key"), ",")
	x.GptApi = make([]*openai.Client, len(chatGptKey))
	for i := 0; i < len(chatGptKey); i++ {
		x.GptApi[i] = x.initGptApi(chatGptKey[i])
	}
	return x
}
