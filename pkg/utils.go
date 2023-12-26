package pkg

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	util "github.com/hktalent/go-utils"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"log"
	"os"
	"path"
)

var GTls *tls.Config

// A+ https://www.ssllabs.com/ssltest/analyze.html?d=www.51pwn.com
func SetSSL() {
	GTls.NextProtos = []string{"h3", "h2", "http/1.1"}
	GTls.MinVersion = tls.VersionTLS12
	GTls.MaxVersion = tls.VersionTLS13
	GTls.CipherSuites = []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		// 1.3
		tls.TLS_AES_128_GCM_SHA256,
		tls.TLS_AES_256_GCM_SHA384,
		tls.TLS_CHACHA20_POLY1305_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA256,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
	}
}

// https://github.com/quic-go/quic-go/wiki/UDP-Receive-Buffer-Size
func SetUdpReceiveBufferSize() {
	util.DoCmd("sysctl", "-w", "net.core.rmem_max=2500000")
	util.DoCmd("sysctl", "-w", "kern.ipc.maxsockbuf=3014656") // , "2>/dev/null"
}

func RunHttp3(addr string, router *gin.Engine) (err error) {
	var cet, key1 = "./config/ca/server.crt", "./config/ca/server.key"
	if nil == GTls {
		if cert, err := tls.LoadX509KeyPair(cet, key1); nil == err {
			GTls = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
			log.Println("debug util1.GTls set ok")
		}
	}

	go SetUdpReceiveBufferSize()
	//certmagic.HTTPS([]string{"51pwn.com", "exploit-poc.com"}, router)
	SetSSL()
	// tcp
	util.Wg.Add(1)
	go func() {
		defer util.Wg.Done()
		h3s := NewMyServer(&http3.Server{
			Addr:            addr,
			EnableDatagrams: true,
			Handler:         router.Handler(),
			TLSConfig:       GTls,
			QuicConfig:      &quic.Config{EnableDatagrams: true},
		})
		//  && !util.GetValAsBool("devDebug")
		if err := h3s.ListenServe(); nil != err {
			log.Println("HTTP/3.0 ListenAndServe ", err)
		}
		//router.RunTLS(addr, certFile, keyFile)
	}()
	return
}

func GetChildDirs(s string) *[]string {
	files, err := os.ReadDir(s)
	if err != nil {
		log.Println(err)
		return nil
	}

	var dirs = make([]string, 0, 10)
	for _, file := range files {
		if file.IsDir() {
			//if strings.HasPrefix(file.Name(), "apache") {
			dirs = append(dirs, path.Join(s, file.Name()))
			//}
		}
	}
	return &dirs
}
