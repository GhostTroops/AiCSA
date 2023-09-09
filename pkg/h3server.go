package pkg

import (
	"crypto/tls"
	"github.com/quic-go/quic-go/http3"
	"net"
	"net/http"
)

type MyServer struct {
	*http3.Server
}

func NewMyServer(Server *http3.Server) *MyServer {
	return &MyServer{Server}
}
func (s *MyServer) ListenServe() error {
	var addr = s.Addr
	//addr, certFile, keyFile string
	var handler http.Handler = s.Handler
	// Load certs
	var err error
	config := s.TLSConfig

	if addr == "" {
		addr = ":https"
	}

	// Open the listeners
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	defer udpConn.Close()

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	tcpConn, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	defer tcpConn.Close()

	tlsConn := tls.NewListener(tcpConn, config)
	defer tlsConn.Close()

	if handler == nil {
		handler = http.DefaultServeMux
	}
	// Start the servers
	httpServer := &http.Server{
		TLSConfig: s.TLSConfig,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.SetQuicHeaders(w.Header())
			handler.ServeHTTP(w, r)
		}),
	}

	hErr := make(chan error)
	qErr := make(chan error)
	go func() {
		hErr <- httpServer.Serve(tlsConn)
	}()
	go func() {
		qErr <- s.Serve(udpConn)
	}()

	select {
	case err := <-hErr:
		s.Close()
		return err
	case err := <-qErr:
		// Cannot close the HTTP server or wait for requests to complete properly :/
		return err
	}
}
