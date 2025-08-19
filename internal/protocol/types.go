package protocol

// Transport is a hint for higher layers (not strictly required by Protocol).
type Transport string

const (
	TransportStdio Transport = "stdio"
	TransportHTTP  Transport = "http"
)

func (t Transport) String() string {
	return string(t)
}

type ConnParams struct {
	Host   string
	APIKey string
}
