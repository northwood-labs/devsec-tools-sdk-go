package devsectools

// DomainResponse represents a response from /domain endpoint
type DomainResponse struct {
	Hostname string `json:"hostname"`
}

// HttpResponse represents a response from /http endpoint
type HttpResponse struct {
	Hostname string `json:"hostname"`
	HTTP11   bool   `json:"http11"`
	HTTP2    bool   `json:"http2"`
	HTTP3    bool   `json:"http3"`
}

// TlsResponse represents a response from /tls endpoint
type TlsResponse struct {
	Hostname     string         `json:"hostname"`
	TLSVersions  TLSVersions    `json:"tlsVersions"`
	TLSConn      []TlsConnection `json:"tlsConnections"`
}

// TLSVersions contains TLS support info
type TLSVersions struct {
	TLS10 bool `json:"tls10"`
	TLS11 bool `json:"tls11"`
	TLS12 bool `json:"tls12"`
	TLS13 bool `json:"tls13"`
}

// TlsConnection represents TLS connection details
type TlsConnection struct {
	Version      string        `json:"version"`
	VersionID    int           `json:"versionId"`
	CipherSuites []CipherSuite `json:"cipherSuites"`
}

// CipherSuite represents a single cipher suite
type CipherSuite struct {
	Authentication string `json:"authentication"`
	Encryption     string `json:"encryption"`
	GnuTLSName     string `json:"gnutlsName"`
	Hash           string `json:"hash"`
	IANAName       string `json:"ianaName"`
	IsAEAD         bool   `json:"isAEAD"`
	IsPFS          bool   `json:"isPFS"`
	KeyExchange    string `json:"keyExchange"`
	OpenSSLName    string `json:"opensslName"`
	Strength       string `json:"strength"`
	URL            string `json:"url"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}
