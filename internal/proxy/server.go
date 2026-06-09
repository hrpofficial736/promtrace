package proxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/hrpofficial736/promtrace/internal/certmanager"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
)

type ProxyServer struct {
	certMgr   *certmanager.CertManager
	store     store.Store
	addr      string
	sessionID string
}

func NewServer(certMgr *certmanager.CertManager, store store.Store, addr, sessionID string) *ProxyServer {
	return &ProxyServer{
		certMgr:   certMgr,
		store:     store,
		addr:      addr,
		sessionID: sessionID,
	}
}

func (ps *ProxyServer) StartServer() error {
	server := &http.Server{
		Addr:    ps.addr,
		Handler: http.HandlerFunc(ps.handleRequest),
	}

	return server.ListenAndServe()
}

func (ps *ProxyServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		ps.handleConnect(w, r)
	}
}

func (ps *ProxyServer) handleConnect(w http.ResponseWriter, r *http.Request) {
	host, _, _ := net.SplitHostPort(r.Host)
	logger.Log.Info("host name is " + host)
	// hijacking connection
	hijacker, _ := w.(http.Hijacker)
	clientConn, _, _ := hijacker.Hijack()

	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// getting a fake tls certificate
	hostCert, err := ps.certMgr.GetOrCreateHostCertificate(host)

	if err != nil {
		logger.Log.Error("error while generating host certs", "error", err)
		return
	}

	// tls handshake with the subprocess
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{*hostCert}}
	tlsConn := tls.Server(clientConn, tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		logger.Log.Error("error while handshaking with the subprocess", "error", err)
	}

	// reading decrypted HTTP request from the subprocess
	reader := bufio.NewReader(tlsConn)
	req, err := http.ReadRequest(reader)

	if err != nil {
		logger.Log.Error("error while reading request", "error", err)
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log.Error("error while parsing req body", "error", err)
		return
	}

	req.Body.Close()

	req.Body = io.NopCloser(bytes.NewReader(reqBody))

	logger.Log.Debug("data", "request content", req)

	// opening real tls connection with the actual original real server
	upstreamConn, err := tls.Dial("tcp", r.Host, &tls.Config{})

	if err != nil {
		logger.Log.Error("error while opening real connection with the actual server", "error", err)
	}
	// forwarding request to the real server
	req.Write(upstreamConn)

	// reading response from the real server
	resp, _ := http.ReadResponse(bufio.NewReader(upstreamConn), req)

	if err != nil {
		logger.Log.Error("error while opening real connection with the actual server", "error", err)
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.Log.Error("error while parsing resp body", "error", err)
		return
	}

	resp.Body.Close()

	resp.Body = io.NopCloser(bytes.NewReader(respBody))

	// saving the trace
	logger.Log.Debug("data", "response body", resp)

	ps.store.SaveTrace(
		&store.Trace{
			ID:           util.GenerateID(),
			SessionID:    ps.sessionID,
			Timestamp:    time.Now(),
			Host:         host,
			Method:       req.Method,
			Path:         req.Pattern,
			Model:        "",
			Tokens:       0,
			Cost:         0,
			SystemPrompt: "",
			UserPrompt:   "",
			Response:     "",
			StatusCode:   200,
			LatencyMs:    100,
			CreatedAt:    time.Now(),
		},
	)

	// traces saved to the store

	// forwarding the response back to the subprocess
	resp.Write(tlsConn)

	// closing both connections
	upstreamConn.Close()
	tlsConn.Close()
}

func (ps *ProxyServer) Shutdown() error {
	return nil
}
