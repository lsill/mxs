package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"io"
	"math/big"
	"mxs/client/test"
	"mxs/log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main2() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		quicConf := &quic.Config{}
		addr := "localhost:33333"
		server :=  http3.Server{
			Server:   &http.Server{
				Addr:              addr,
				Handler:           setupHandler(addr),
			},
			QuicConfig: quicConf,
		}
		err := server.ListenAndServeTLS(test.GetCertificatePaths())
		if err != nil {
			log.Error("err is %v", err)
		}
		wg.Done()
	}()
	wg.Wait()
}


func setupHandler(www string) http.Handler {
	mux := http.NewServeMux()

	if len(www) > 0 {
		mux.Handle("/", http.FileServer(http.Dir(www)))
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			log.Release("%#v", r)
			const maxSize = 1 << 32	// 1 GB
			num, err := strconv.ParseInt(strings.ReplaceAll(r.RequestURI, "/", ""), 10, 64)
			if err != nil || num <= 0 || num > maxSize {
				w.WriteHeader(400)
				return
			}
			w.Write(generatePRData(int(num)))
		})
	}
	mux.HandleFunc("/mx/pos", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("dasdasd"))
	})
	return mux
}

func generatePRData(l int) []byte {
	res := make([]byte, l)
	seed := uint32(1)
	for i := 0; i < l; i++{
		seed = seed * 48271 %2147483647
		res[i] = byte(seed)
	}
	return res
}

func main3 () {
	addr := "localhost:4242"

	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		log.Error("%v", err)
		return
	}
	go func () {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			log.Error("%v", err)
			return
		}
		stream, err := sess.AcceptStream(context.Background())
		if err != nil {
			log.Error("%v", err)
			panic(err)
		}
		// Echo through the loggingWriter
		_, err = io.Copy(loggingWriter{stream}, stream)
		log.Release("go end")
		return
	}()
	for {

	}

}

func main() {
	load()
	QuicServer("localhost:3333", nil)
	for {

	}
}

var (
	DefaultTLSConfig *tls.Config
)

func load() {
	cert, err := GenCertificate()
	if err != nil {
		panic(err)
	}
	DefaultTLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
}

func QuicServer(addr string, config *QUICConfig) (quic.Listener, error) {
	if config == nil {
		config = &QUICConfig{}
	}
	quicConfig := &quic.Config{
		HandshakeTimeout:                      config.Timeout,
		MaxIdleTimeout:                        config.IdleTimeout,
		KeepAlive:                             config.KeepAlive,
	}
	tlsConfig := config.TLSConfig
	if tlsConfig == nil {
		tlsConfig = DefaultTLSConfig
	}
	var conn net.PacketConn

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	lconn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	conn = lconn
	if config.Key != nil {
		conn = &quicCipherConn{
			UDPConn:lconn,
			key:config.Key,
		}
	}
	ln, err := quic.Listen(conn, tlsConfig, quicConfig)
	if err != nil {
		return nil, err
	}
	l := &quicListener{
		ln:       ln,
		connChan: make(chan quicConn,1024),
		errChan:  make(chan error, 1),
	}
	go l.listenLoop()
	return nil, nil
}

func (l *quicListener) listenLoop () {
	for {
		session, err := l.ln.Accept(context.Background())
		if err != nil {
			log.Error("[quic] accept err %v", err)
			l.errChan <- err
			close(l.errChan)
			return
		}
		go l.sessionLoop(session)
	}
}

func (l *quicListener) sessionLoop(session quic.Session) {
	log.Debug("[quic] begin %s <-> %s", session.RemoteAddr(), session.LocalAddr())
	defer log.Debug("[quic] After %s <-> %s", session.RemoteAddr(), session.LocalAddr())
	for {
		stream, err := session.AcceptStream(context.Background())
		if err !=nil {
			log.Error("[quic] accept steam err %v", err)
			session.CloseWithError(100, fmt.Sprintf("%v", err))
			return
		}

		cc := quicConn{
			Stream: stream,
			laddr:  session.LocalAddr(),
			raddr:  session.RemoteAddr(),
		}
		select {
		case l.connChan <- cc:
		default:
			cc.Close()
			log.Error("[quic] %s<-> %s: connection queue is full", session.RemoteAddr(), session.LocalAddr())
		}
	}
}

type quicConn struct {
	quic.Stream
	laddr net.Addr
	raddr net.Addr
}

type quicListener struct {
	ln       quic.Listener
	connChan chan quicConn
	errChan  chan error
}

type quicCipherConn struct {
	*net.UDPConn
	key []byte
}

// QUICConfig is the config for QUIC client and server
type QUICConfig struct {
	TLSConfig   *tls.Config
	Timeout     time.Duration
	KeepAlive   bool
	IdleTimeout time.Duration
	Key         []byte
}

func echoServer () error {

	addr := "localhost:4242"
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	sess, err := listener.Accept(context.Background())
	if err != nil {
		return err
	}
	stream, err := sess.AcceptStream(context.Background())
	if err != nil {
		panic(err)
	}
	// Echo through the loggingWriter
	_, err = io.Copy(loggingWriter{stream}, stream)
	log.Release("go end")
	return err
}

type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}


func GenCertificate() (cert tls.Certificate, err error) {
	rawCert, rawKey, err := generateKeyPair()
	if err != nil {
		return
	}
	return tls.X509KeyPair(rawCert, rawKey)
}

func generateKeyPair() (rawCert, rawKey []byte, err error) {
	// Create private key and self-signed certificate
	// Adapted from https://golang.org/src/crypto/tls/generate_cert.go

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	validFor := time.Hour * 24 * 365 * 10 // ten years
	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"gost"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return
	}

	rawCert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	rawKey = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return
}
