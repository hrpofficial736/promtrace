package proxy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"runtime"
	"time"
)

type TLSInterface interface {
	GenerateRootCAKeyPair() error
}

type tlsStruct struct {
	username string
}

func generateRootCAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	return privateKey, publicKey, nil
}

func (ts *tlsStruct) GenerateRootCACertificate() {
	privKey, pubKey, err := generateRootCAKeyPair()

	if err != nil {

	}

	// creating certificate template
	caTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2026),

		Subject: pkix.Name{
			Organization: []string{"promtrace-ca"},
			CommonName:   "my local root promtrace ca",
		},

		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(10, 0, 0),
		IsCA:      true,

		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},

		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,

		BasicConstraintsValid: true,
	}

	// creating certificate
	caBytes, _ := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, pubKey, privKey)

	// creating CA certificate file containing public key and other info
	caFile, _ := os.Create("ca.crt")
	pem.Encode(caFile, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes})
	caFile.Close()

	// creating CA private key file
	keyFile, _ := os.Create("ca.key")
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})
	keyFile.Close()
}

func StoreRootCACertificateInSystemTrustStore() {
	var installer TrustStoreInstaller

	switch runtime.GOOS {
	case "darwin":
		installer = DarwinInstaller{}
	case "linux":
		installer = LinuxInstaller{}
	case "windows":
		installer = WindowsInstaller{}
	}

	err := installer.Install("./rootCA.pem")

	if err != nil {
		panic(err)
	}
}
