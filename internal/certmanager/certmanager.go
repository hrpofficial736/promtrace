package certmanager

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/hrpofficial736/promtrace/internal/config"
	"github.com/hrpofficial736/promtrace/internal/truststore"
	"math/big"
	"os"
	"runtime"
	"sync"
	"time"
)

type CertManager struct {
	caCert *x509.Certificate
	caKey  *rsa.PrivateKey
	cfg    *config.Config
	cache  map[string]*tls.Certificate
	mu     sync.RWMutex
}

func NewCertManager(cfg *config.Config) (*CertManager, error) {
	return &CertManager{
		cfg:   cfg,
		cache: make(map[string]*tls.Certificate),
	}, nil
}

func generateCAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	return privateKey, publicKey, nil
}

func (cm *CertManager) GetCertPath() string {
	return cm.cfg.CACert
}

func (cm *CertManager) GetConfigDir() string {
	return cm.cfg.Dir
}

func (cm *CertManager) GetDBPath() string {
	return cm.cfg.DBPath
}

func (cm *CertManager) GenerateRootCACertificate() error {
	privKey, pubKey, err := generateCAKeyPair()

	if err != nil {
		return err
	}

	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

	// creating certificate template
	caTemplate := &x509.Certificate{
		SerialNumber: serial,

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
	caBytes, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, pubKey, privKey)

	if err != nil {
		return err
	}

	// creating CA certificate file containing public key and other info
	caFile, err := os.Create(cm.cfg.CACert)

	if err != nil {
		return err
	}
	err = pem.Encode(caFile, &pem.Block{Type: "CERTIFICATE", Bytes: caBytes})
	if err != nil {
		return err
	}
	err = caFile.Close()
	if err != nil {
		return err
	}

	// creating CA private key file
	keyFile, err := os.Create(cm.cfg.CAKey)

	if err != nil {
		return err
	}
	err = pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privKey)})
	if err != nil {
		return err
	}
	err = keyFile.Close()
	if err != nil {
		return err
	}

	cm.caCert = caTemplate
	cm.caKey = privKey

	return nil
}

func (cm *CertManager) StoreRootCACertificateInSystemTrustStore() {
	var installer truststore.TrustStoreInstaller

	switch runtime.GOOS {
	case "darwin":
		installer = truststore.DarwinInstaller{}
	case "linux":
		installer = truststore.LinuxInstaller{}
	case "windows":
		installer = truststore.WindowsInstaller{}
	}

	err := installer.Install(cm.cfg.CACert)

	if err != nil {
		panic(err)
	}
}

func (cm *CertManager) LoadCA() error {
	encKeyFile, err := os.ReadFile(cm.cfg.CAKey)
	if err != nil {
		return err
	}

	encCrtFile, err := os.ReadFile(cm.cfg.CACert)
	if err != nil {
		return err
	}

	decKey, _ := pem.Decode(encKeyFile)

	decCrt, _ := pem.Decode(encCrtFile)

	parsedKey, err := x509.ParsePKCS1PrivateKey(decKey.Bytes)
	if err != nil {
		return err
	}
	parsedCrt, err := x509.ParseCertificate(decCrt.Bytes)
	if err != nil {
		return err
	}
	cm.caKey = parsedKey
	cm.caCert = parsedCrt

	return nil

}

func (cm *CertManager) GetOrCreateHostCertificate(hostname string) (*tls.Certificate, error) {

	cm.mu.RLock()
	// returning host certificate from cache if already there
	if cert, ok := cm.cache[hostname]; ok {
		cm.mu.RUnlock()
		return cert, nil
	}

	cm.mu.RUnlock()

	// generating fake private and public key for the host in the original request URL
	// this public key would be used to generate a fake host certificate

	fakePrivKey, _ := rsa.GenerateKey(rand.Reader, 3072)

	fakePubKey := &fakePrivKey.PublicKey

	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	certTemplate := &x509.Certificate{
		SerialNumber: serial,

		Subject: pkix.Name{
			CommonName: hostname,
		},
		DNSNames: []string{hostname},

		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24),
		IsCA:      false,

		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},

		KeyUsage: x509.KeyUsageDigitalSignature,

		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, certTemplate, cm.caCert, fakePubKey, cm.caKey)

	if err != nil {
		return nil, err
	}

	// building tls certificate

	tlsCert := &tls.Certificate{
		Certificate: [][]byte{cert},
		PrivateKey:  fakePrivKey,
	}

	cm.mu.Lock()
	cm.cache[hostname] = tlsCert
	cm.mu.Unlock()
	return tlsCert, nil
}
