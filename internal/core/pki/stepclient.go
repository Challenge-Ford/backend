package pki

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-jose/go-jose/v4"
	josejwt "github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
)

type IssuedCertificate struct {
	Certificate  string
	PrivateKey   string
	SerialNumber string
}

type StepCAClient struct {
	caURL       string
	provisioner string
	httpClient  *http.Client
	signingKey  *jose.JSONWebKey
	kid         string
}

// NewStepCAClient creates a StepCA PKI client.
// rootCACertPath: path to the root CA cert PEM file.
func NewStepCAClient(caURL, provisioner, password, rootCACertPath string) (*StepCAClient, error) {
	rootPEM, err := os.ReadFile(rootCACertPath)
	if err != nil {
		return nil, fmt.Errorf("pki: read root ca cert: %w", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(rootPEM) {
		return nil, fmt.Errorf("pki: failed to parse root ca cert")
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: pool},
		},
		Timeout: 15 * time.Second,
	}

	c := &StepCAClient{
		caURL:       caURL,
		provisioner: provisioner,
		httpClient:  httpClient,
	}

	if err := c.loadProvisionerKey(password); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *StepCAClient) loadProvisionerKey(password string) error {
	resp, err := c.httpClient.Get(c.caURL + "/1.0/provisioners")
	if err != nil {
		return fmt.Errorf("pki: fetch provisioners: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		Provisioners []struct {
			Type         string          `json:"type"`
			Name         string          `json:"name"`
			Key          json.RawMessage `json:"key"`
			EncryptedKey string          `json:"encryptedKey"`
		} `json:"provisioners"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("pki: decode provisioners: %w", err)
	}

	for _, p := range result.Provisioners {
		if p.Type == "JWK" && p.Name == c.provisioner {
			var pubKey jose.JSONWebKey
			if err := json.Unmarshal(p.Key, &pubKey); err != nil {
				return fmt.Errorf("pki: parse provisioner public key: %w", err)
			}
			c.kid = pubKey.KeyID

			jwe, err := jose.ParseEncrypted(
				p.EncryptedKey,
				[]jose.KeyAlgorithm{jose.PBES2_HS256_A128KW, jose.PBES2_HS384_A192KW, jose.PBES2_HS512_A256KW},
				[]jose.ContentEncryption{jose.A128CBC_HS256, jose.A192CBC_HS384, jose.A256CBC_HS512, jose.A128GCM, jose.A192GCM, jose.A256GCM},
			)
			if err != nil {
				return fmt.Errorf("pki: parse encrypted key: %w", err)
			}
			decrypted, err := jwe.Decrypt([]byte(password))
			if err != nil {
				return fmt.Errorf("pki: decrypt provisioner key: %w", err)
			}

			var privKey jose.JSONWebKey
			if err := json.Unmarshal(decrypted, &privKey); err != nil {
				return fmt.Errorf("pki: parse decrypted key: %w", err)
			}
			c.signingKey = &privKey
			return nil
		}
	}
	return fmt.Errorf("pki: JWK provisioner %q not found", c.provisioner)
}

func (c *StepCAClient) Issue(ctx context.Context, commonName string) (*IssuedCertificate, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("pki: generate key: %w", err)
	}

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: commonName},
	}, privKey)
	if err != nil {
		return nil, fmt.Errorf("pki: create csr: %w", err)
	}
	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})

	token, err := c.newToken(commonName, c.caURL+"/1.0/sign")
	if err != nil {
		return nil, err
	}

	body, _ := json.Marshal(map[string]any{
		"csr": string(csrPEM),
		"ott": token,
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.caURL+"/1.0/sign", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("pki: sign request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("pki: sign failed (%d): %s", resp.StatusCode, raw)
	}

	var signResp struct {
		CRT string `json:"crt"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&signResp); err != nil {
		return nil, fmt.Errorf("pki: decode sign response: %w", err)
	}

	serial, err := extractSerial(signResp.CRT)
	if err != nil {
		return nil, err
	}

	keyDER, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("pki: marshal key: %w", err)
	}
	keyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}))

	return &IssuedCertificate{
		Certificate:  signResp.CRT,
		PrivateKey:   keyPEM,
		SerialNumber: serial,
	}, nil
}

func (c *StepCAClient) Revoke(ctx context.Context, serialNumber string) error {
	token, err := c.newToken(serialNumber, c.caURL+"/1.0/revoke")
	if err != nil {
		return err
	}

	body, _ := json.Marshal(map[string]any{
		"serial":     serialNumber,
		"ott":        token,
		"reasonCode": 0,
	})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.caURL+"/1.0/revoke", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("pki: revoke request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pki: revoke failed (%d): %s", resp.StatusCode, raw)
	}
	return nil
}

func (c *StepCAClient) newToken(subject, audience string) (string, error) {
	sig, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.ES256, Key: c.signingKey},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", c.kid),
	)
	if err != nil {
		return "", fmt.Errorf("pki: create signer: %w", err)
	}

	now := time.Now()
	claims := struct {
		josejwt.Claims
		SANS []string `json:"sans"`
	}{
		Claims: josejwt.Claims{
			Issuer:   c.provisioner,
			Subject:  subject,
			Audience: josejwt.Audience{audience},
			ID:       uuid.NewString(),
			IssuedAt: josejwt.NewNumericDate(now),
			Expiry:   josejwt.NewNumericDate(now.Add(5 * time.Minute)),
		},
		SANS: []string{subject},
	}

	token, err := josejwt.Signed(sig).Claims(claims).Serialize()
	if err != nil {
		return "", fmt.Errorf("pki: sign token: %w", err)
	}
	return token, nil
}

func extractSerial(certPEM string) (string, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return "", fmt.Errorf("pki: failed to decode certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("pki: parse certificate: %w", err)
	}
	return cert.SerialNumber.String(), nil
}
