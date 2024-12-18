package stdlib

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

// AbstractCacheRepository defines a generic interface for interacting with a Redis-based cache.
// V represents the type of the values stored in the cache.
type AbstractCacheRepository[V any] interface {

	// Get retrieves a value by its key from the cache.
	// If the value is a struct, it will be deserialized from JSON.
	// Returns a pointer to the value or nil if the key does not exist.
	Get(key string) (valueModel *V, err error)

	// GetKeysByPatterns retrieves keys by a pattern from the cache.
	// Returns a slice of keys that match the pattern.
	GetKeysByPatterns(pattern string) (keys []string, err error)

	// Set stores a value in the cache with the specified expiration time.
	// If the value is a struct, it will be serialized to JSON.
	// Returns an error if the operation fails.
	Set(key string, value V, expiration time.Duration) error

	// Del deletes a value from the cache by its key.
	// Returns an error if the operation fails.
	Del(key string) error

	// Exists checks if a key exists in the cache.
	// Returns true if the key exists, false otherwise.
	Exists(key string) (bool, error)

	// HGet retrieves a single field value from a hash in Redis.
	// The method returns the value associated with the specified field
	// and nil if the field does not exist or an error occurs.
	HGet(key string, field string) (valueModel *V, err error)

	// HGetAll retrieves all fields and their associated values from a hash in Redis.
	// The method returns a map of field names to values or an error if the operation fails.
	HGetAll(key string) (map[string]V, error)

	// HGetFields retrieves specific fields and their associated values from a hash in Redis.
	// The method returns a map of the requested field names to their values.
	// Fields not found in the hash are excluded from the returned map.
	HGetFields(key string, fields ...string) (map[string]V, error)

	// HSet sets a single field in a hash in Redis.
	// This method stores the specified value under the given field name,
	// overwriting any existing value.
	HSet(key string, field string, value V) error

	// HDel deletes a specific field from a hash in Redis.
	// This method removes the field and its value, returning an error if the operation fails.
	HDel(key string, field string) error

	// HExists checks if a specific field exists in a hash in Redis.
	// The method returns true if the field exists, false otherwise.
	HExists(key string, field string) (bool, error)
}

type abstractCacheRepositoryImpl[V any] struct {
	client      *redis.Client
	ctx         context.Context
	isPrimitive bool
	self        AbstractCacheRepository[V]
}

// Get implements AbstractCacheRepository.
func (repo *abstractCacheRepositoryImpl[V]) Get(key string) (*V, error) {
	result, err := repo.client.Get(repo.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	value, err := deserialize[V]([]byte(result), repo.isPrimitive)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func (repo *abstractCacheRepositoryImpl[V]) GetKeysByPatterns(pattern string) ([]string, error) {
	var cursor uint64
	var keys []string
	for {
		scanKeys, newCursor, err := repo.client.Scan(repo.ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}
		keys = append(keys, scanKeys...)
		cursor = newCursor
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}

// Set implements AbstractCacheRepository.
func (repo *abstractCacheRepositoryImpl[V]) Set(key string, value V, expiration time.Duration) error {
	data, err := serialize(value)
	if err != nil {
		return err
	}
	return repo.client.Set(repo.ctx, key, data, expiration).Err()
}

// Exists implements AbstractCacheRepository.
func (repo *abstractCacheRepositoryImpl[V]) Exists(key string) (bool, error) {
	count, err := repo.client.Exists(repo.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Del implements Abstrac tCacheRepository.
func (repo *abstractCacheRepositoryImpl[V]) Del(key string) error {
	return repo.client.Del(repo.ctx, key).Err()
}

// HGet retrieves a single field value from a hash.
func (repo *abstractCacheRepositoryImpl[V]) HGet(key string, field string) (*V, error) {
	result, err := repo.client.HGet(repo.ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	value, err := deserialize[V]([]byte(result), repo.isPrimitive)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

// HGetAll retrieves all fields and values from a hash.
func (repo *abstractCacheRepositoryImpl[V]) HGetAll(key string) (map[string]V, error) {
	result, err := repo.client.HGetAll(repo.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	fields := make(map[string]V, len(result))
	for k, v := range result {
		value, err := deserialize[V]([]byte(v), repo.isPrimitive)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize field %s: %w", k, err)
		}
		fields[k] = value
	}
	return fields, nil
}

// HSet sets a single field in a hash.
func (repo *abstractCacheRepositoryImpl[V]) HSet(key string, field string, value V) error {
	data, err := serialize(value)
	if err != nil {
		return err
	}
	return repo.client.HSet(repo.ctx, key, field, data).Err()
}

// HDel deletes a field from a hash.
func (repo *abstractCacheRepositoryImpl[V]) HDel(key string, field string) error {
	return repo.client.HDel(repo.ctx, key, field).Err()
}

// HExists checks if a field exists in a hash.
func (repo *abstractCacheRepositoryImpl[V]) HExists(key string, field string) (bool, error) {
	return repo.client.HExists(repo.ctx, key, field).Result()
}

// Helper function to determine if a type is primitive.
func isPrimitiveType(value any) bool {
	switch value.(type) {
	case string, int, float64, bool, []byte:
		return true
	default:
		return false
	}
}

// Helper function to serialize a value.
func serialize[V any](value V) ([]byte, error) {
	if isPrimitiveType(value) {
		return []byte(fmt.Sprintf("%v", value)), nil
	}
	return sonic.Marshal(value)
}

// Helper function to deserialize data.
func deserialize[V any](data []byte, isPrimitive bool) (V, error) {
	var value V
	if isPrimitive {
		return any(string(data)).(V), nil
	}
	if err := sonic.Unmarshal(data, &value); err != nil {
		return value, fmt.Errorf("failed to deserialize value: %w", err)
	}
	return value, nil
}

func CreateCacheRepository[V any](redisClient *redis.Client, ctx context.Context, self AbstractCacheRepository[V]) *abstractCacheRepositoryImpl[V] {
	repo := &abstractCacheRepositoryImpl[V]{
		client:      redisClient,
		ctx:         ctx,
		isPrimitive: isPrimitiveType(new(V)),
		self:        self,
	}
	return repo
}
