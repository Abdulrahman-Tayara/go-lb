package buckettokens

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBucket_ShouldReset(t *testing.T) {
	type input struct {
		b    *bucket
		rate time.Duration
	}
	tests := []struct {
		name  string
		input input
		want  bool
	}{
		{
			name: "should reset",
			input: input{
				b: &bucket{
					last: time.Now().Add(-1 * time.Second),
				},
				rate: time.Second * 2,
			},
			want: false,
		},
		{
			name: "should not reset",
			input: input{
				b: &bucket{
					last: time.Now().Add(-10 * time.Second),
				},
				rate: time.Second * 2,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.input.b.shouldReset(tt.input.rate))
		})
	}
}

func TestBuckets_Get(t *testing.T) {
	b := newBuckets(&Configs{
		Tokens: 13,
		Rate:   time.Second,
	})

	b.data.Store("some key", &bucket{
		tokens: 10,
	})

	b.data.Store("key to delete", &bucket{
		tokens: 10,
	})

	b.data.Delete("key to delete")

	tests := []struct {
		name  string
		input string // key
		want  *bucket
	}{
		{
			"get not found key",
			"some new key",
			&bucket{
				tokens: b.configs.Tokens,
			},
		},
		{
			"existing key",
			"some key",
			&bucket{
				tokens: 10,
			},
		},
		{
			"deleted key",
			"key to delete",
			&bucket{
				tokens: b.configs.Tokens,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want.tokens, b.get(tt.input).tokens)
		})
	}

}

func TestBuckets_DeleteFunc(t *testing.T) {
	b := newBuckets(&Configs{
		Rate:   time.Second,
		Tokens: 13,
	})

	b.data.Store("k1", &bucket{
		tokens: 10,
		last:   time.Now().Add(-4 * time.Second),
	})
	b.data.Store("k2", &bucket{
		tokens: 10,
		last:   time.Now().Add(-1 * time.Second),
	})

	b.deleteFunc(func(b *bucket) bool {
		return b.shouldReset(2 * time.Second)
	})

	_, ok1 := b.data.Load("k1")
	_, ok2 := b.data.Load("k2")

	assert.False(t, ok1)
	assert.True(t, ok2)
}
