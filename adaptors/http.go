package adaptors

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httputil"

	"github.com/int128/kubelogin/adaptors/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func NewHTTP(i HTTP) adaptors.HTTP {
	return &i
}

type HTTP struct {
	dig.In
	Logger adaptors.Logger
}

func (*HTTP) NewClientConfig() adaptors.HTTPClientConfig {
	return &httpClientConfig{
		certPool: x509.NewCertPool(),
	}
}

func (h *HTTP) NewClient(config adaptors.HTTPClientConfig) (*http.Client, error) {
	return &http.Client{
		Transport: &loggingTransport{
			base: &http.Transport{
				TLSClientConfig: config.TLSConfig(),
				Proxy:           http.ProxyFromEnvironment,
			},
			logger: h.Logger,
		},
	}, nil
}

type loggingTransport struct {
	base   http.RoundTripper
	logger adaptors.Logger
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	const level = 2
	if t.logger.GetDebugLevel() < level {
		return t.base.RoundTrip(req)
	}
	var dumpBody bool
	if t.logger.GetDebugLevel() > level {
		dumpBody = true
	}

	reqDump, err := httputil.DumpRequestOut(req, dumpBody)
	if err != nil {
		t.logger.Debugf(level, "Error: could not dump the request: %s", err)
		return t.base.RoundTrip(req)
	}
	t.logger.Debugf(level, "%s", string(reqDump))
	resp, err := t.base.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	respDump, err := httputil.DumpResponse(resp, dumpBody)
	if err != nil {
		t.logger.Debugf(level, "Error: could not dump the response: %s", err)
		return resp, err
	}
	t.logger.Debugf(level, "%s", string(respDump))
	return resp, err
}

type httpClientConfig struct {
	certPool      *x509.CertPool
	skipTLSVerify bool
}

func (c *httpClientConfig) AddCertificateFromFile(filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "could not read %s", filename)
	}
	if c.certPool.AppendCertsFromPEM(b) != true {
		return errors.Errorf("could not append certificate from %s", filename)
	}
	return nil
}

func (c *httpClientConfig) AddEncodedCertificate(base64String string) error {
	b, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return errors.Wrapf(err, "could not decode base64")
	}
	if c.certPool.AppendCertsFromPEM(b) != true {
		return errors.Errorf("could not append certificate")
	}
	return nil
}

func (c *httpClientConfig) TLSConfig() *tls.Config {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.skipTLSVerify,
	}
	if len(c.certPool.Subjects()) > 0 {
		tlsConfig.RootCAs = c.certPool
	}
	return tlsConfig
}

func (c *httpClientConfig) SetSkipTLSVerify(b bool) {
	c.skipTLSVerify = b
}
