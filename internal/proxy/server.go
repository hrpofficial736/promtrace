package proxy

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/hrpofficial736/promtrace/internal/certmanager"
	"github.com/hrpofficial736/promtrace/internal/logger"
	"github.com/hrpofficial736/promtrace/internal/provider"
	"github.com/hrpofficial736/promtrace/internal/store"
	"github.com/hrpofficial736/promtrace/internal/util"
	"github.com/hrpofficial736/promtrace/pkg/costable"
	"github.com/hrpofficial736/promtrace/pkg/tokenizer"
)

type ProxyServer struct {
	certMgr   *certmanager.CertManager
	store     store.Store
	addr      string
	sessionID string

	tlsConn      *tls.Conn
	upstreamConn *tls.Conn
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
		return
	}

	ps.tlsConn = tlsConn

	// reading decrypted HTTP request from the subprocess
	reader := bufio.NewReader(tlsConn)
	req, err := http.ReadRequest(reader)

	if err != nil {
		logger.Log.Error("error while reading request", "error", err)
		return
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log.Error("error while parsing req body", "error", err)
		return
	}

	req.Body.Close()

	req.Body = io.NopCloser(bytes.NewReader(reqBody))

	// opening real tls connection with the actual original real server
	upstreamConn, err := tls.Dial("tcp", r.Host, &tls.Config{})

	if err != nil {
		logger.Log.Error("error while opening real connection with the actual server", "error", err)
		return
	}
	ps.upstreamConn = upstreamConn
	// forwarding request to the real server
	req.Write(upstreamConn)

	start := time.Now()
	// reading response from the real server
	resp, err := http.ReadResponse(bufio.NewReader(upstreamConn), req)

	if err != nil {
		logger.Log.Error("error while reading response from the actual server", "error", err)
		return
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		logger.Log.Error("error while parsing resp body", "error", err)
		return
	}

	latency := time.Since(start).Milliseconds()

	resp.Body.Close()

	resp.Body = io.NopCloser(bytes.NewReader(respBody))

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gr, err := gzip.NewReader(bytes.NewReader(respBody))
		if err == nil {
			respBody, err = io.ReadAll(gr)
			gr.Close()
		}
	}

	ext := provider.GetExtractor(host)

	if ext == nil {
		logger.Log.Error("invalid host name", "host", host)
		return
	}

	model, sysPrompt, userPrompt := ext.ExtractRequest(reqBody, req.URL.Path)
	response, inTokens, outTokens := ext.ExtractResponse(respBody)

	if inTokens == 0 && outTokens == 0 {
		inTokens = tokenizer.EstimateTokens(sysPrompt + userPrompt)
		outTokens = tokenizer.EstimateTokens(response)
	}

	cost := costable.CalculateCost(model, inTokens, outTokens)

	now := time.Now()

	// saving the trace

	ps.store.SaveTrace(
		&store.Trace{
			ID:           util.GenerateID(),
			SessionID:    ps.sessionID,
			Timestamp:    now,
			Host:         host,
			Method:       req.Method,
			Path:         strings.Split(req.URL.Path, "?")[0],
			Model:        model,
			Tokens:       inTokens + outTokens,
			Cost:         cost,
			SystemPrompt: sysPrompt,
			UserPrompt:   userPrompt,
			RequestBody:  string(reqBody),
			Response:     response,
			StatusCode:   resp.StatusCode,
			LatencyMs:    latency,
			CreatedAt:    now,
		},
	)

	// forwarding the response back to the subprocess
	resp.Write(tlsConn)
}

func (ps *ProxyServer) Shutdown() error {
	if ps.tlsConn != nil {
		ps.tlsConn.Close()
	}

	if ps.upstreamConn != nil {
		ps.upstreamConn.Close()
	}

	return nil
}
