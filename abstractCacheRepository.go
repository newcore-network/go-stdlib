package stdlib

import (
	"context"
	"fmt"
	"reflect"
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
}

type abstractCacheRepositoryImpl[V any] struct {
	client *redis.Client
	ctx    context.Context
	self   AbstractCacheRepository[V]
}

// Get implements AbstractCacheRepository.
func (repo *abstractCacheRepositoryImpl[V]) Get(key string) (*V, error) {
	var value V

	result, err := repo.client.Get(repo.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	if reflect.TypeOf(value).Kind() == reflect.Struct {
		err := sonic.Unmarshal([]byte(result), &value)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize struct: %w", err)
		}
	} else {
		converted, ok := any(result).(V)
		if !ok {
			return nil, fmt.Errorf("failed to cast value to type %T", value)
		}
		value = converted
	}

	return &value, nil
}

func (repo *abstractCacheRepositoryImpl[V]) GetKeysByPatterns(pattern string) ([]string, error) {
	var cursor uint64
	var keys []string
	for {
		scanKeys, newCursor, err := repo.client.Scan(repo.ctx, cursor, pattern, 10).Result()
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
	switch v := any(value).(type) {
	case string, int, float64, bool:
		return repo.client.Set(repo.ctx, key, fmt.Sprintf("%v", v), expiration).Err()
	case []byte:
		return repo.client.Set(repo.ctx, key, v, expiration).Err()
	default:
		jsonData, err := sonic.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to serialize struct: %w", err)
		}
		return repo.client.Set(repo.ctx, key, jsonData, expiration).Err()
	}
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

func CreateCacheRepository[V any](redisClient *redis.Client, ctx context.Context, self AbstractCacheRepository[V]) *abstractCacheRepositoryImpl[V] {
	repo := &abstractCacheRepositoryImpl[V]{
		client: redisClient,
		ctx:    ctx,
		self:   self,
	}
	return repo
}
