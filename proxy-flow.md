YES — this is much cleaner mentally.

In your case there is:

```txt id="g1"
CLI App (subprocess)
        ↓
Local Proxy (Promtrace)
        ↓
LLM API Server
```

NO browser UI involved.

But internally, the subprocess/app is still acting like an HTTPS client exactly like a browser would.

So replace:

```txt id="g2"
browser
```

with:

```txt id="g3"
subprocess/app/client SDK
```

Now let's rewrite the ENTIRE flow specifically for YOUR architecture.

---

# YOUR REAL ARCHITECTURE

Suppose user runs:

```bash id="g4"
promtrace run claude-cli
```

or:

```bash id="g5"
promtrace wrap my-ai-agent
```

Internally:

```txt id="g6"
Your CLI launches subprocess
```

Example subprocess:

* Claude Code
* OpenAI CLI
* custom AI agent
* LangChain app
* Go/Python SDK app

---

# Then Proxy Variables Are Injected

Example:

```txt id="g7"
HTTPS_PROXY=http://localhost:8080
HTTP_PROXY=http://localhost:8080
```

Now subprocess sends HTTPS traffic THROUGH your proxy.

---

# COMPLETE FLOW (CLI VERSION)

---

# PHASE 1 — SETUP

Usually:

```bash id="g8"
promtrace setup
```

---

# Step 1 — Promtrace Generates Root CA Key Pair

Creates:

```txt id="g9"
CA private key
CA public key
```

---

# Step 2 — Promtrace Creates Root CA Certificate

Certificate contains:

```txt id="g10"
CA public key
```

Promtrace self-signs it using:

```txt id="g11"
CA private key
```

Result:

```txt id="g12"
ca.crt
```

and:

```txt id="g13"
ca.key
```

---

# Step 3 — User Trusts Root CA Certificate

Promtrace installs:

```txt id="g14"
ca.crt
```

into:

* OS trust store
* maybe language runtime trust store

Now machine says:

```txt id="g15"
I trust certificates signed by Promtrace CA.
```

IMPORTANT:

Only the CERTIFICATE is trusted.

NOT the private key.

---

# PHASE 2 — SUBPROCESS STARTS

---

# Step 4 — Promtrace Launches Subprocess

Example:

```bash id="g16"
claude-cli
```

or some AI agent.

---

# Step 5 — Proxy Environment Variables Injected

Subprocess now routes HTTPS through:

```txt id="g17"
localhost:8080
```

where Promtrace proxy is running.

---

# PHASE 3 — HTTPS REQUEST HAPPENS

Suppose subprocess calls:

```txt id="g18"
https://api.openai.com
```

---

# Step 6 — Subprocess Connects To Proxy

NOT directly to OpenAI.

Actual connection:

```txt id="g19"
subprocess
    ↓
Promtrace proxy
```

---

# Step 7 — Proxy Sees Target Host

Proxy sees:

```txt id="g20"
api.openai.com
```

---

# Step 8 — Proxy Generates Fake Host Key Pair

Creates:

```txt id="g21"
fake host private key
fake host public key
```

for:

```txt id="g22"
api.openai.com
```

---

# Step 9 — Proxy Creates Fake Host Certificate

Certificate contains:

```txt id="g23"
domain = api.openai.com
fake host public key
```

---

# Step 10 — Proxy Signs Fake Host Certificate

Using:

```txt id="g24"
CA private key
```

Result:

```txt id="g25"
fake api.openai.com certificate
```

---

# Step 11 — Proxy Sends Fake Host Certificate To Subprocess

Subprocess receives:

```txt id="g26"
fake api.openai.com certificate
```

---

# Step 12 — Subprocess Verifies Fake Certificate

Subprocess uses:

```txt id="g27"
CA public key
(from trusted root CA certificate)
```

to verify:

```txt id="g28"
Was this cert signed by trusted CA?
```

Since user trusted Promtrace CA earlier:

verification succeeds.

---

# Step 13 — Subprocess Challenges Proxy

Subprocess now says:

```txt id="g29"
Prove you own the private key
matching this certificate.
```

---

# Step 14 — Proxy Uses Fake Host Private Key

Proxy uses:

```txt id="g30"
fake host private key
```

to produce cryptographic proof.

---

# Step 15 — Subprocess Verifies Proof

Using:

```txt id="g31"
fake host public key
(from fake host certificate)
```

verification succeeds.

NOW:

```txt id="g32"
TLS connection established
between subprocess and proxy
```

---

# IMPORTANT

At this point:

```txt id="g33"
subprocess thinks proxy IS api.openai.com
```

because TLS checks passed.

---

# PHASE 4 — PROXY CONNECTS TO REAL LLM SERVER

Now proxy becomes HTTPS client itself.

---

# Step 16 — Proxy Connects To Real Server

Connection:

```txt id="g34"
Promtrace proxy
        ↓
api.openai.com
```

---

# Step 17 — Real Server Sends Real Certificate

Contains:

```txt id="g35"
real host public key
```

signed by real internet CA.

---

# Step 18 — Proxy Verifies Real Certificate

Using trusted internet CAs like:

* Let's Encrypt
* DigiCert

---

# Step 19 — Real Server Proves Ownership

Using:

```txt id="g36"
real host private key
```

Proxy verifies using:

```txt id="g37"
real host public key
```

---

# Step 20 — Second TLS Connection Established

Now:

```txt id="g38"
Proxy ←TLS→ Real OpenAI Server
```

---

# FINAL STATE

You now have TWO independent TLS tunnels.

---

# Tunnel 1

```txt id="g39"
Subprocess ←TLS→ Promtrace
```

using:

* fake host cert
* fake host private key

---

# Tunnel 2

```txt id="g40"
Promtrace ←TLS→ Real Server
```

using:

* real host cert
* real host private key

---

# Promtrace Sits In Middle

So traffic becomes:

```txt id="g41"
Encrypted
    ↓
Promtrace decrypts
    ↓
Reads/logs/traces
    ↓
Re-encrypts
    ↓
Sends to real server
```

---

# FINAL VISUAL MODEL

```txt id="g42"
┌──────────────────┐
│ Subprocess/App   │
│ (Claude CLI etc) │
└────────┬─────────┘
         │
         │ TLS using fake cert
         ▼
┌──────────────────┐
│ Promtrace Proxy  │
│ localhost:8080   │
└────────┬─────────┘
         │
         │ TLS using real cert
         ▼
┌──────────────────┐
│ api.openai.com   │
│ Real LLM Server  │
└──────────────────┘
```

# MOST IMPORTANT INSIGHT

Your subprocess NEVER talks directly to real OpenAI server.

It talks to:

```txt id="g43"
Promtrace pretending to be OpenAI
```

BUT:

because user trusted your CA,

the subprocess accepts that identity as valid.

That is the entire mechanism behind HTTPS interception in your CLI.

