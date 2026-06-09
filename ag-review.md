# Promtrace — Code Review (v2)

## Overall Verdict

**The structure is now correct and the MITM flow in `server.go` is logically right.** You've clearly understood the guide. The separation into `certmanager/`, `proxy/`, `subprocess/`, `truststore/` is exactly right. The CONNECT handler follows the correct sequence: hijack → fake cert → TLS handshake → read request → dial upstream → forward → read response → forward back.

There are bugs to fix before this will actually run, but the architecture is sound. Here's everything I found, grouped by file.

---

## 1. `main.go` — Import path is wrong

```go
import "github.com/hrpofficial736/promtrace/cli"  // ← old top-level cli/ package
```

You moved cli into `internal/cli/`, so this needs to be:
```go
import "github.com/hrpofficial736/promtrace/internal/cli"
```

---

## 2. `certmanager.go` — 7 issues

### 2a. `~/.promtrace/` doesn't expand in Go

[Lines 75, 80, 84](file:///home/harshitrajpandey/cli-projects/promtrace/internal/certmanager/certmanager.go#L75-L84)

`os.Create("~/.promtrace/ca.crt")` won't work. Go does NOT expand `~` — that's a shell feature. You need to resolve the home directory yourself:

```go
homeDir, _ := os.UserHomeDir()
certDir := filepath.Join(homeDir, ".promtrace")
os.MkdirAll(certDir, 0700)
caFile, _ := os.Create(filepath.Join(certDir, "ca.crt"))
```

Also: you accept `certDir` in `NewCertManager()` but then overwrite it with `"~/.promtrace"` on line 84. Use `cm.certDir` that was passed in.

### 2b. Serial number is still hardcoded

[Lines 53, 126](file:///home/harshitrajpandey/cli-projects/promtrace/internal/certmanager/certmanager.go#L53)

`big.NewInt(2026)` — every cert gets serial 2026. Generate a random one:
```go
serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
```

### 2c. Host cert has wrong CommonName

[Line 129](file:///home/harshitrajpandey/cli-projects/promtrace/internal/certmanager/certmanager.go#L129)

```go
CommonName: "my local root promtrace ca"  // ← this is the CA name, not the host
```
Should be:
```go
CommonName: hostname
```

### 2d. Host cert has wrong KeyUsage

[Line 139](file:///home/harshitrajpandey/cli-projects/promtrace/internal/certmanager/certmanager.go#L139)

`KeyUsageCertSign` is for CAs only. A leaf cert should have:
```go
KeyUsage: x509.KeyUsageDigitalSignature
```

Also, `ExtKeyUsage` only needs `ServerAuth` (remove `ClientAuth` — the proxy is impersonating a server, not a client):
```go
ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
```

### 2e. RLock is acquired but never released

[Lines 112-116](file:///home/harshitrajpandey/cli-projects/promtrace/internal/certmanager/certmanager.go#L112-L116)

You call `cm.mu.RLock()` but never `cm.mu.RUnlock()`. Same problem with `cm.mu.Lock()` on line 157 — no `Unlock()`. Use `defer`:

```go
// For the cache read:
cm.mu.RLock()
if cert, ok := cm.cache[hostname]; ok {
    cm.mu.RUnlock()
    return cert, nil
}
cm.mu.RUnlock()

// ... generate cert ...

// For the cache write:
cm.mu.Lock()
cm.cache[hostname] = tlsCert
cm.mu.Unlock()
```

### 2f. Cache is never initialized

[Line 22](file:///home/harshitrajpandey/cli-projects/promtrace/internal/certmanager/certmanager.go#L22)

`cache map[string]*tls.Certificate` — this is `nil` by default. Writing to a nil map panics. Initialize it in `NewCertManager`:

```go
func NewCertManager(certDir string) *CertManager {
    return &CertManager{
        certDir: certDir,
        cache:   make(map[string]*tls.Certificate),
    }
}
```

### 2g. Errors are silently ignored with `_`

Lines 72, 75, 80, 121 — every `_` on an error is a silent failure. Handle them all. At minimum, return the error.

---

## 3. `certmanager.go` — Truststore path mismatch

[Line 103](file:///home/harshitrajpandey/cli-projects/promtrace/internal/certmanager/certmanager.go#L103)

```go
installer.Install("~/.promtrace/rootCA.pem")
```

But you write the cert as `ca.crt`, not `rootCA.pem`. And again, `~` doesn't expand. This should use the resolved path:
```go
installer.Install(filepath.Join(cm.certDir, "ca.crt"))
```

---

## 4. `certmanager.go` — Missing `LoadCA()` method

Your `wrap` command calls `GenerateRootCACertificate()` every time. But setup should generate the CA once, and wrap should **load** the existing CA from disk. You need a `LoadCA()` method that:
1. Reads `ca.crt` and `ca.key` from `cm.certDir`
2. PEM-decodes them
3. Parses them with `x509.ParseCertificate()` and `x509.ParsePKCS1PrivateKey()`
4. Sets `cm.caCert` and `cm.caKey`

Then `setup` calls `GenerateCA()` + truststore install, and `wrap` calls `LoadCA()`.

---

## 5. `server.go` — Host may include port, which breaks TLS dial

[Lines 39, 48, 60](file:///home/harshitrajpandey/cli-projects/promtrace/internal/proxy/server.go#L39-L60)

`r.Host` for a CONNECT request is `api.openai.com:443`. You pass this to `GetOrCreateHostCertificate(host)` — the cert gets `DNSNames: ["api.openai.com:443"]` which is wrong (DNS names don't include ports). You need to strip the port for cert generation:

```go
hostname, _, _ := net.SplitHostPort(r.Host)
// use hostname for cert generation
// use r.Host (with port) for tls.Dial
```

---

## 6. `server.go` — Request body needs buffering before forwarding

[Lines 57-63](file:///home/harshitrajpandey/cli-projects/promtrace/internal/proxy/server.go#L57-L63)

Right now you read the request and forward it, but you never capture the body for logging. The problem: `req.Body` is a stream — once you forward it to upstream (via `req.Write`), the body is consumed and gone. You need to read and buffer it first:

```go
// Read the body
bodyBytes, _ := io.ReadAll(req.Body)
req.Body.Close()

// Restore it for forwarding
req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

// Now bodyBytes has the raw request body for logging
// And req still has a readable body for forwarding
```

Same applies to the response body (lines 66-71) — read it, buffer it for logging, then write a new response back.

---

## 7. `intercept.go` — Dead code

This file contains a standalone `hijackConnection` function that is never called. The hijack logic is already inline in `server.go`'s `handleConnect`. **Delete this file** — it's confusing to have duplicate logic.

---

## 8. `wrap.go` — Incomplete wiring

[Lines 8-16](file:///home/harshitrajpandey/cli-projects/promtrace/internal/cli/wrap.go#L8-L16)

Currently `wrapRun` creates a CertManager and generates a CA, but:
- It doesn't create a `ProxyServer`
- It doesn't start the proxy
- It doesn't launch the subprocess
- The function isn't registered as a cobra command (no `init()` adding it to `rootCmd`)
- It should call `LoadCA()`, not `GenerateRootCACertificate()` (setup does the generation)

The correct flow for `wrapRun` should be:
1. Resolve `certDir` (using `os.UserHomeDir()`)
2. Create CertManager → `LoadCA()`
3. Create ProxyServer → `go ps.StartServer()`
4. Launch subprocess → `subprocess.LaunchChildProcessWithEnvVars(args, addr)`
5. `cmd.Wait()` — block until subprocess exits
6. `ps.Shutdown()`

---

## 9. `truststore.go` — Windows typo

[Line 87](file:///home/harshitrajpandey/cli-projects/promtrace/internal/truststore/truststore.go#L87)

```go
"-addscore"  // ← typo
```
Should be:
```go
"-addstore"
```

---

## What to Do Next (in order)

You have the right skeleton. Now fix things from the inside out:

| Step | What | Why |
|------|------|-----|
| **1** | Fix all the `certmanager.go` bugs above (2a–2g) | Nothing works if certs are broken |
| **2** | Add `LoadCA()` method to CertManager | Wrap needs to load, not regenerate |
| **3** | Wire `setup.go` — call `GenerateCA()` + truststore install | First testable command |
| **4** | Test: run `promtrace setup`, verify `~/.promtrace/ca.crt` exists, inspect it with `openssl x509 -in ~/.promtrace/ca.crt -text -noout` |
| **5** | Fix `server.go` — strip port from host, buffer request/response bodies | Core MITM correctness |
| **6** | Delete `intercept.go` | Dead code |
| **7** | Wire `wrap.go` — full flow (load CA → start proxy → launch subprocess → wait) + register as cobra command |
| **8** | Fix `main.go` import path | Binary won't compile without this |
| **9** | Test end-to-end: `promtrace wrap curl https://api.openai.com/v1/models` — you should see the request pass through and get a response (even if it's a 401 from OpenAI, that means the proxy works) |
| **10** | Add body logging (print the intercepted request/response to stdout for now — SQLite store comes later) |

> [!TIP]
> **Step 4 and Step 9 are the two checkpoints.** If step 4 works, your cert generation is solid. If step 9 works, your entire MITM proxy works end-to-end. Everything after that (store, TUI, diff) is features built on a working foundation.

