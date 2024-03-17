package certificates

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"os"
	"time"
)

func InitCertificates(pathCert string, pathKey string) error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Printf("Failed to generate private key: %v", err)
		return err

	}

	// serial number from CA
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Printf("Failed to generate serial number: %v", err)
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"nikitads9"},
		},
		DNSNames:  []string{"api.booking-schedule.su"},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(3 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template /* parent*/, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Printf("Failed to create certificate: %v", err)
		return err
	}

	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemCert == nil {
		log.Println("Failed to encode certificate to PEM")
		return errors.New("failed to encode certificate to PEM")
	}
	if err := os.WriteFile(pathCert, pemCert, 0644); err != nil {

		return err
	}

	log.Print("wrote cert.pem\n")

	err = createPrivateKey(privateKey, pathKey)
	if err != nil {
		log.Printf("Failed to create private key")
		return err
	}

	return nil
}

func createPrivateKey(privateKey *ecdsa.PrivateKey, pathKey string) error {
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Printf("Unable to marshal private key: %v", err)
		return err
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if pemKey == nil {
		log.Println("Failed to encode key to PEM")
		return errors.New("failed to encode key to PEM")
	}
	if err := os.WriteFile(pathKey, pemKey, 0600); err != nil {
		return err
	}
	log.Print("wrote key.pem\n")

	return nil
}
