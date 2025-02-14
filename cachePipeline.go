package stdlib

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CachePipeline is a wrapper around redis.Pipeliner that allows you to
// chain commands and add options (e.g a TTL) in a convenient way.
type CachePipeline struct {
	pipe redis.Pipeliner
	ctx  context.Context
	err  error
}

// HSet sets a single field in a Redis hash.
//
// This method adds the field to the hash or updates its value if it already exists.
func (p *CachePipeline) HSet(key, field string, value any) *CachePipeline {
	if p.err != nil {
		return p
	}
	if key == "" || field == "" {
		p.err = errors.New("key and field must not be empty")
	}
	data, err := serialize(value)
	if err != nil {
		p.err = err
		return p
	}
	p.pipe.HSet(p.ctx, key, field, data)
	return p
}

// HMSet sets multiple fields in a Redis hash.
//
// This method allows batch setting of multiple fields in a single Redis operation.
func (p *CachePipeline) HMSet(key string, fields map[string]any) *CachePipeline {
	if p.err != nil {
		return p
	}
	if key == "" || len(fields) == 0 {
		p.err = errors.New("key cannot be empty and fields must not be empty")
		return p
	}
	serializedFields := make(map[string]any, len(fields))
	for field, value := range fields {
		data, err := serialize(value)
		if err != nil {
			p.err = fmt.Errorf("failed to serialize field %s: %w", field, err)
			return p
		}
		serializedFields[field] = data
	}
	p.pipe.HMSet(p.ctx, key, serializedFields)
	return p
}

// HDel removes one or more fields from a Redis hash.
func (p *CachePipeline) HDel(key string, fields ...string) *CachePipeline {
	if p.err != nil {
		return p
	}
	if key == "" || len(fields) == 0 {
		p.err = errors.New("key cannot be empty and at least one field must be specified")
		return p
	}
	p.pipe.HDel(p.ctx, key, fields...)
	return p
}

func (p *CachePipeline) Del(keys ...string) *CachePipeline {
	if p.err != nil {
		return p
	}
	if len(keys) == 0 {
		p.err = errors.New("at least one key must be specified")
		return p
	}
	p.pipe.Del(p.ctx, keys...)
	return p
}

func (p *CachePipeline) Set(key string, value any, expiration time.Duration) *CachePipeline {
	if p.err != nil {
		return p
	}
	if key == "" {
		p.err = errors.New("key cannot be empty")
		return p
	}
	data, err := serialize(value)
	if err != nil {
		p.err = err
		return p
	}
	p.pipe.Set(p.ctx, key, data, expiration)
	return p
}

// Expire sets an expiration time for a Redis hash.
//
// This method ensures that the hash is automatically removed after a specified duration.
func (p *CachePipeline) Expire(key string, expiration time.Duration) *CachePipeline {
	if p.err != nil {
		return p
	}
	if key == "" {
		p.err = errors.New("key cannot be empty")
		return p
	}
	p.pipe.Expire(p.ctx, key, expiration)
	return p
}

// Incr increments the value of a Redis key by amount.
func (p *CachePipeline) IncrBy(key string, amount int64) *CachePipeline {
	if p.err != nil {
		return p
	}
	if key == "" {
		p.err = errors.New("key cannot be empty")
		return p
	}
	p.pipe.IncrBy(p.ctx, key, amount)
	return p
}

// Decr decrements the value of a Redis key by amount.
func (p *CachePipeline) DecrBy(key string, amount int64) *CachePipeline {
	if p.err != nil {
		return p
	}
	if key == "" {
		p.err = errors.New("key cannot be empty")
		return p
	}
	p.pipe.DecrBy(p.ctx, key, amount)
	return p
}

// Exec executes all queued operations in the Redis pipeline and returns results.
func (p *CachePipeline) Exec() ([]redis.Cmder, error) {
	if p.err != nil {
		return nil, p.err
	}
	return p.pipe.Exec(p.ctx)
}

// ExecAndDiscard executes the Redis pipeline but does not return results.
//
// This method is useful for operations where you don't need command results.
func (p *CachePipeline) ExecAndDiscard() error {
	if p.err != nil {
		return p.err
	}
	_, err := p.pipe.Exec(p.ctx)
	return err
}
