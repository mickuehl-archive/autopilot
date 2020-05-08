/*
 * Copyright (C) 2018 Josh A. Beam
 * All rights reserved.
 *
 * See https://github.com/joshb/pi-camera-go
 */

package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path"
	"time"
)

func KeyPaths() (string, string, error) {
	keyDir, err := ConfigDir("keys")
	if err != nil {
		return "", "", err
	}

	privateKeyPath := path.Join(keyDir, "private.pem")
	_, err = os.Stat(privateKeyPath)
	privateKeyExists := err == nil

	publicKeyPath := path.Join(keyDir, "public.pem")
	_, err = os.Stat(publicKeyPath)
	publicKeyExists := err == nil

	if !privateKeyExists || !publicKeyExists {
		if err := createKeys(privateKeyPath, publicKeyPath); err != nil {
			return "", "", err
		}
	}

	return privateKeyPath, publicKeyPath, nil
}

func createKeys(privateKeyPath, publicKeyPath string) error {
	serialNumMax := (&big.Int{}).Lsh(big.NewInt(1), 256)
	serialNum, err := rand.Int(rand.Reader, serialNumMax)
	if err != nil {
		return err
	}

	startTime := time.Now()
	endTime := startTime.AddDate(1, 0, 0)

	cert := x509.Certificate{
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		NotBefore:             startTime,
		NotAfter:              endTime,
		SerialNumber:          serialNum,
		Subject:               pkix.Name{Organization: []string{"pi-camera-go"}},
		DNSNames:              []string{"localhost"},
	}

	println("Generating RSA key...")

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	println("Generating certificate...")

	b, err := x509.CreateCertificate(rand.Reader, &cert, &cert, &key.PublicKey, key)
	if err != nil {
		return err
	}

	privateKeyFile, err := os.Create(privateKeyPath)
	if err != nil {
		return err
	}

	defer privateKeyFile.Close()
	block := pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
	if err := pem.Encode(privateKeyFile, &block); err != nil {
		return err
	}

	publicKeyFile, err := os.Create(publicKeyPath)
	if err != nil {
		return err
	}

	defer publicKeyFile.Close()
	block = pem.Block{Type: "CERTIFICATE", Bytes: b}
	if err := pem.Encode(publicKeyFile, &block); err != nil {
		return err
	}

	println("Done generating certificate")
	return nil
}
