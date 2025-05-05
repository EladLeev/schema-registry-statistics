package main

import (
	"crypto/x509"
	"os"
	"testing"
)

func TestLoadCert(t *testing.T) {
	caCertFile, err := os.CreateTemp("", "cacert*.pem")
	if err != nil {
		t.Fatalf("failed to create temporary file: %s", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			return
		}
	}(caCertFile.Name())

	caCertData := []byte(`
-----BEGIN CERTIFICATE-----
MIIBwjCCAWugAwIBAgIUTvn+Zc80K0mGpMAoGCCqGSM49BAMCMFMxKTAnBgNV
BAMMIDFwcGxlIENlcnRpZmljYXRlIEF1dGhvcml0eTAeFw0yMjAyMTcxMDU3
MTlaFw0yMjAzMTcxMDU3MTlaMFMxKTAnBgNVBAMMIDFwcGxlIENlcnRpZmlj
YXRlIEF1dGhvcml0eTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABLh6FJzh
Kj/ISL9Rbhg0N14O/r5WWqNw4euBJNzPbNybc+n4ebNkMttcV6U9az6POoyG
Ucky6hGz2jBBRGaUuV6jUDBOMB0GA1UdDgQWBBQ2Wpqw4q3iG4nJZc+uM7N/
Y4qr4DAOBgNVHQ8BAf8EBAMCAQYwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsG
AQUFBwMCMCAGA1UdEQQZMBeCFGFwaWtleS5jb22CCnN1Ym1pdC5jb20wCgYI
KoZIzj0EAwIDSAAwRQIgZH1OqzW8NfBvZHXrNmlT0TtIJ0QQs+z7E2N1blSC
X/0CIQDKuZwuAQzS1aA90xSgbbVi/TuV7Yj4l4uV7lRGkW8HvA==
-----END CERTIFICATE-----
`)
	if _, err := caCertFile.Write(caCertData); err != nil {
		t.Fatalf("failed to write to temporary file: %s", err)
	}

	tlsConfig := loadCert(caCertFile.Name())

	if tlsConfig.RootCAs == nil {
		t.Error("RootCAs is nil")
	}

	if !tlsConfig.RootCAs.Equal(x509.NewCertPool()) {
		t.Error("something is wrong with the RootCAs")
	}
}
