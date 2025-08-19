package pool

import (
	"testing"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	p := New(2, 3*time.Second)
	defer p.Close()

	ok := p.Set("key1", meilisearch.New("http://localhost:7700"))
	assert.False(t, ok)

	v, ok := p.Get("key1")
	assert.True(t, ok)
	assert.NotNil(t, v)

	ok = p.Set("key2", meilisearch.New("http://localhost:7700"))
	assert.False(t, ok)

	v, ok = p.Get("key2")
	assert.True(t, ok)
	assert.NotNil(t, v)

	ok = p.Set("key3", meilisearch.New("http://localhost:7700"))
	assert.True(t, ok)

	length := p.Len()
	assert.Equal(t, 2, length)

	ok = p.Remove("key1")
	assert.False(t, ok)

	time.Sleep(4 * time.Second)

	_, ok = p.Get("key2")
	assert.False(t, ok)
}
