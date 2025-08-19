package protocol

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/ja7ad/meilisearch-mcp/internal/logger"
	"github.com/ja7ad/meilisearch-mcp/internal/pool"
	"github.com/ja7ad/meilisearch-mcp/internal/util"
	"github.com/meilisearch/meilisearch-go"
)

type Protocol struct {
	transport    Transport
	service      meilisearch.ServiceManager
	host, apiKey string
	pool         *pool.Pool[string, meilisearch.ServiceManager]
	val          *validator.Validate
	logging      *logger.SubLogger

	mu sync.Mutex
}

func New(transport Transport, host, apiKey string, pool *pool.Pool[string, meilisearch.ServiceManager]) *Protocol {
	p := &Protocol{
		pool:      pool,
		host:      host,
		apiKey:    apiKey,
		transport: transport,
		val:       validator.New(),
	}

	p.logging = logger.NewSubLogger("_protocol", p)

	return p
}

func (p *Protocol) String() string {
	return "meilisearch"
}

func (p *Protocol) client(header http.Header) (meilisearch.ServiceManager, error) {
	p.logging.Debug("Getting client based on transport", "transport", p.transport)

	switch p.transport {
	case TransportHTTP:
		if p.pool == nil {
			return nil, ErrNilPool
		}

		host, apiKey, hash := util.MeilisearchHeaders(header)
		if host == "" {
			return nil, ErrMissingHostHeader
		}

		p.logging.Debug("Using host from headers", "host", host)

		if cli, ok := p.pool.Get(hash); ok && cli != nil {
			return cli, nil
		}

		opts := buildOptions(apiKey)
		service, err := meilisearch.Connect(host, opts...)
		if err != nil {
			return nil, fmt.Errorf("connect http service: %w", err)
		}
		p.pool.Set(hash, service)
		return service, nil

	case TransportStdio:
		p.mu.Lock()
		defer p.mu.Unlock()

		if p.service != nil {
			return p.service, nil
		}

		opts := buildOptions(p.apiKey)
		service, err := meilisearch.Connect(p.host, opts...)
		if err != nil {
			return nil, fmt.Errorf("connect stdio service: %w", err)
		}
		p.service = service
		return p.service, nil
	}

	return nil, ErrInvalidTransport
}

func (p *Protocol) validate(field any, tag string) error {
	if err := p.val.Var(field, tag); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			msg := ""

			for _, fieldErr := range verrs {
				switch fieldErr.Tag() {
				case "required":
					msg = "is required"
				case "min":
					msg = fmt.Sprintf("must be at least %s characters", fieldErr.Param())
				case "max":
					msg = fmt.Sprintf("must be at most %s characters", fieldErr.Param())
				case "email":
					msg = "must be a valid email address"
				case "teener":
					msg = "must be a teener (12–18)"
				case "gtcsfield":
					msg = "must be greater than " + fieldErr.Param()
				case "ltfield":
					msg = "must be less than " + fieldErr.Param()
				case "gtfield":
					msg = "must be greater than " + fieldErr.Param()
				case "latitude":
					msg = "must be a valid latitude" + fieldErr.Param()
				case "longitude":
					msg = "must be a valid longitude" + fieldErr.Param()
				default:
					msg = fmt.Sprintf("violates %s constraint", fieldErr.Tag())
				}
			}
			return fmt.Errorf("failed validate field %v: %s", field, msg)
		}
		return err
	}
	return nil
}

func buildOptions(apiKey string) []meilisearch.Option {
	if apiKey == "" {
		return nil
	}
	return []meilisearch.Option{meilisearch.WithAPIKey(apiKey)}
}
