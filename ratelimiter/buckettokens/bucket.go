package buckettokens

import (
	"sync"
	"time"
)

type bucket struct {
	tokens int
	last   time.Time
}

func (b *bucket) take() {
	b.tokens--
}

func (b *bucket) hasTokens() bool {
	return b.tokens > 0
}

func (b *bucket) shouldReset(rate time.Duration) bool {
	return b.last.Add(rate).After(time.Now())
}

type buckets struct {
	configs *Configs
	data    *sync.Map
}

func newBuckets(cfg *Configs) *buckets {
	return &buckets{
		data:    &sync.Map{},
		configs: cfg,
	}
}

func (b *buckets) get(key string) *bucket {
	bu, _ := b.data.LoadOrStore(key, &bucket{
		tokens: b.configs.Tokens,
		last:   time.Now(),
	})
	return bu.(*bucket)
}

func (b *buckets) delete(key string) {
	b.data.Delete(key)
}

func (b *buckets) deleteFunc(predicate func(*bucket) bool) {
	b.data.Range(func(key, value any) bool {
		if predicate(value.(*bucket)) {
			b.data.Delete(key)
		}
		return true
	})
}
