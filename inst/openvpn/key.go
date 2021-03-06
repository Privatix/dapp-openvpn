package openvpn

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func randomNumber() int64 {
	r, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	return r.Int64()
}

func buildServerCertificate(path string, expired time.Time) error {
	commonName, err := os.Hostname()
	if err != nil {
		return err
	}

	ca := &x509.Certificate{
		SerialNumber: big.NewInt(randomNumber()),
		Subject:      pkix.Name{CommonName: commonName},
		NotBefore:    time.Now(),
		NotAfter:     expired,
		IsCA:         true,
		BasicConstraintsValid: true,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	if err := buildCA(ca, path); err != nil {
		return err
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(randomNumber()),
		Subject:      pkix.Name{CommonName: commonName},
		NotBefore:    time.Now(),
		NotAfter:     expired,
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage: x509.KeyUsageDigitalSignature,
	}

	return buildCertificate(cert, "server", path)
}

func buildCA(ca *x509.Certificate, path string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	bytes, err := x509.CreateCertificate(rand.Reader, ca, ca,
		&priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// Public key.
	certOut, err := os.Create(filepath.Join(path, "ca.crt"))
	if err != nil {
		return err
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE",
		Bytes: bytes}); err != nil {
		return err
	}

	// Private key.
	keyOut, err := os.OpenFile(filepath.Join(path, "ca.key"),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer keyOut.Close()

	return pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv)})
}

func buildCertificate(cert *x509.Certificate, name, path string) error {
	// Load CA.
	catls, err := tls.LoadX509KeyPair(
		filepath.Join(path, "ca.crt"),
		filepath.Join(path, "ca.key"),
	)
	if err != nil {
		return err
	}
	ca, err := x509.ParseCertificate(catls.Certificate[0])
	if err != nil {
		return err
	}

	priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	// Sign the certificate.
	bytes, err := x509.CreateCertificate(rand.Reader, cert, ca,
		&priv.PublicKey, catls.PrivateKey)

	// Public key.
	certOut, err := os.Create(filepath.Join(path, name+".crt"))
	if err != nil {
		return err
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE",
		Bytes: bytes}); err != nil {
		return err
	}

	// Private key.
	keyOut, err := os.OpenFile(filepath.Join(path, name+".key"),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer keyOut.Close()

	return pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv)})
}
