package pool

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/ja7ad/meilisearch-mcp/internal/logger"
	"github.com/meilisearch/meilisearch-go"
)

type Pool[K ~string, T meilisearch.ServiceManager] struct {
	client       *expirable.LRU[K, T]
	evictedCh    chan func()
	stopCh       chan struct{}
	wg           sync.WaitGroup
	shuttingDown atomic.Bool
	logging      *logger.SubLogger
}

func New(size int, ttl time.Duration) *Pool[string, meilisearch.ServiceManager] {
	p := &Pool[string, meilisearch.ServiceManager]{
		evictedCh: make(chan func(), 1024),
		stopCh:    make(chan struct{}),
	}

	if ttl == 0 {
		ttl = 3 * time.Minute // Default TTL if not specified
	}

	if size <= 0 {
		size = 100 // Default size if not specified
	}

	p.client = expirable.NewLRU[string, meilisearch.ServiceManager](
		size,
		func(_ string, v meilisearch.ServiceManager) {
			if v == nil {
				return
			}

			if p.shuttingDown.Load() {
				go v.Close()
				return
			}

			select {
			case p.evictedCh <- v.Close:
			default:
				go v.Close()
			}
		},
		ttl,
	)

	p.wg.Add(1)
	go p.worker()

	p.logging = logger.NewSubLogger("_pool", p)

	return p
}

func (*Pool[K, T]) String() string {
	return "Lru cache"
}

func (p *Pool[K, T]) worker() {
	defer p.wg.Done()
	for {
		select {
		case fn := <-p.evictedCh:
			func() {
				defer func() { _ = recover() }()
				if fn != nil {
					fn()
				}
			}()
		case <-p.stopCh:
			for {
				select {
				case fn := <-p.evictedCh:
					func() {
						defer func() { _ = recover() }()
						if fn != nil {
							fn()
						}
					}()
				default:
					return
				}
			}
		}
	}
}

func (p *Pool[K, T]) Get(key K) (T, bool)     { return p.client.Get(key) }
func (p *Pool[K, T]) Set(key K, value T) bool { return p.client.Add(key, value) }
func (p *Pool[K, T]) Remove(key K) bool       { return p.client.Remove(key) }
func (p *Pool[K, T]) Len() int                { return p.client.Len() }

func (p *Pool[K, T]) Close() {
	p.shuttingDown.Store(true)
	p.client.Purge()
	close(p.stopCh)
	p.wg.Wait()
}
