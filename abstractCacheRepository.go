package stdlib

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

// AbstractCacheRepository defines a generic interface for interacting with a Redis-based cache.
// T represents the type of the values stored in the cache.
type AbstractCacheRepository[T any] interface {

	// Get retrieves a value by its key from the cache.
	// If the value is a struct, it will be deserialized from JSON.
	// Returns a pointer to the value or nil if the key does not exist.
	Get(key string) (valueModel T, err error)

	// GetKeysByPatterns retrieves keys by a pattern from the cache.
	// Returns a slice of keys that match the pattern.
	GetKeysByPatterns(pattern string) (keys []string, err error)

	// Set stores a value in the cache with the specified expiration time.
	// If the value is a struct, it will be serialized to JSON.
	// Returns an error if the operation fails.
	Set(key string, value T, expiration time.Duration) error

	// Del deletes a value from the cache by its key.
	// Returns an error if the operation fails.
	Del(key string) error

	// Exists checks if a key exists in the cache.
	// Returns true if the key exists, false otherwise.
	Exists(key string) (bool, error)

	// HGet retrieves a single field value from a hash in Redis.
	// The method returns the value associated with the specified field
	// and nil if the field does not exist or an error occurs.
	HGet(key string, field string) (*any, error)

	// HGetAll retrieves all fields and their associated values from a hash in Redis.
	// The method returns a map of field names to values or an error if the operation fails.
	HGetAll(key string) (map[string]any, error)

	// HScan iterates over fields in a hash by a pattern.
	// Returns the matching fields and their values.
	HScan(key string, pattern string, count int64) (map[string]string, error)

	// HGetFields retrieves specific fields and their associated values from a hash in Redis.
	// The method returns a map of the requested field names to their values.
	// Fields not found in the hash are excluded from the returned map.
	HGetFields(key string, fields ...string) (map[string]any, error)

	// HSet sets a single field in a hash in Redis.
	// This method stores the specified value under the given field name,
	// overwriting any existing value.
	HSet(key string, field string, value any) error

	// HMSet sets multiple fields in a hash in Redis.
	HMSet(key string, fields map[string]any) error

	// HDel deletes a specific field from a hash in Redis.
	// This method removes the field and its value, returning an error if the operation fails.
	HDel(key string, field string) error

	// HExists checks if a specific field exists in a hash in Redis.
	// The method returns true if the field exists, false otherwise.
	HExists(key string, field string) (bool, error)

	// NewPipeline creates a new pipeline, which allows you to chain commands and add options (e.g a TTL) in a convenient way.
	NewPipeline() *CachePipeline
}

type abstractCacheRepositoryImpl[T any] struct {
	client      *redis.Client
	ctx         context.Context
	isPrimitive bool
	self        AbstractCacheRepository[T]
}

// Get implements AbstractCacheRepository.
func (repo *abstractCacheRepositoryImpl[T]) Get(key string) (T, error) {
	var value T

	if repo.self != repo {
		return repo.self.Get(key)
	}
	result, err := repo.client.Get(repo.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return value, nil
		}
		return value, err
	}
	value, err = deserialize[T]([]byte(result), repo.isPrimitive)
	if err != nil {
		return value, err
	}
	return value, nil
}

func (repo *abstractCacheRepositoryImpl[T]) GetKeysByPatterns(pattern string) ([]string, error) {
	if repo.self != repo {
		return repo.self.GetKeysByPatterns(pattern)
	}
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
func (repo *abstractCacheRepositoryImpl[T]) Set(key string, value T, expiration time.Duration) error {
	if repo.self != repo {
		return repo.self.Set(key, value, expiration)
	}
	data, err := serialize(value)
	if err != nil {
		return err
	}
	return repo.client.Set(repo.ctx, key, data, expiration).Err()
}

// Exists implements AbstractCacheRepository.
func (repo *abstractCacheRepositoryImpl[T]) Exists(key string) (bool, error) {
	if repo.self != repo {
		return repo.self.Exists(key)
	}
	count, err := repo.client.Exists(repo.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Del implements Abstrac tCacheRepository.
func (repo *abstractCacheRepositoryImpl[T]) Del(key string) error {
	if repo.self != repo {
		return repo.self.Del(key)
	}
	return repo.client.Del(repo.ctx, key).Err()
}

// HGet retrieves a single field value from a hash.
func (repo *abstractCacheRepositoryImpl[T]) HGet(key string, field string) (*any, error) {
	if repo.self != repo {
		return repo.self.HGet(key, field)
	}
	result, err := repo.client.HGet(repo.ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	var value any
	if err := sonic.Unmarshal([]byte(result), &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func (repo *abstractCacheRepositoryImpl[T]) HScan(key string, pattern string, count int64) (map[string]string, error) {
	var cursor uint64
	result := make(map[string]string)

	for {
		fields, newCursor, err := repo.client.HScan(repo.ctx, key, cursor, pattern, count).Result()
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(fields); i += 2 {
			result[fields[i]] = fields[i+1]
		}

		cursor = newCursor
		if cursor == 0 {
			break
		}
	}

	return result, nil
}

// HGetAll retrieves all fields and values from a hash.
func (repo *abstractCacheRepositoryImpl[T]) HGetAll(key string) (map[string]any, error) {
	result, err := repo.client.HGetAll(repo.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	fields := make(map[string]any, len(result))
	for k, v := range result {
		var value any
		if err := sonic.Unmarshal([]byte(v), &value); err != nil {
			return nil, fmt.Errorf("failed to deserialize field %s: %w", k, err)
		}
		fields[k] = value
	}
	return fields, nil
}

func (repo *abstractCacheRepositoryImpl[T]) HGetFields(key string, fields ...string) (map[string]any, error) {
	if repo.self != repo {
		return repo.self.HGetFields(key, fields...)
	}
	result, err := repo.client.HMGet(repo.ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}
	values := make(map[string]any)
	for i, field := range fields {
		if result[i] != nil {
			var value any
			if err := sonic.Unmarshal([]byte(result[i].(string)), &value); err != nil {
				return nil, fmt.Errorf("failed to deserialize field %s: %w", field, err)
			}
			values[field] = value
		}
	}
	return values, nil
}

// HSet sets a single field in a hash.
func (repo *abstractCacheRepositoryImpl[T]) HSet(key string, field string, value any) error {
	if repo.self != repo {
		return repo.self.HSet(key, field, value)
	}
	if key == "" || field == "" {
		return errors.New("key and field must not be empty")
	}

	var data []byte
	switch v := value.(type) {
	case string:
		data = []byte(v)
	default:
		var err error
		data, err = sonic.Marshal(v)
		if err != nil {
			return err
		}
	}
	return repo.client.HSet(repo.ctx, key, field, data).Err()
}

func (repo *abstractCacheRepositoryImpl[T]) HMSet(key string, fields map[string]any) error {
	if repo.self != repo {
		return repo.self.HMSet(key, fields)
	}

	serializedFields := make(map[string]any, len(fields))
	for field, value := range fields {
		switch v := value.(type) {
		case string:
			serializedFields[field] = v
		default:
			data, err := sonic.Marshal(v)
			if err != nil {
				return fmt.Errorf("failed to serialize field %s: %w", field, err)
			}
			serializedFields[field] = string(data)
		}
	}

	return repo.client.HMSet(repo.ctx, key, serializedFields).Err()
}

// HDel deletes a field from a hash.
func (repo *abstractCacheRepositoryImpl[T]) HDel(key string, field string) error {
	if repo.self != repo {
		return repo.self.HDel(key, field)
	}
	return repo.client.HDel(repo.ctx, key, field).Err()
}

// HExists checks if a field exists in a hash.
func (repo *abstractCacheRepositoryImpl[T]) HExists(key string, field string) (bool, error) {
	if repo.self != repo {
		return repo.self.HExists(key, field)
	}
	return repo.client.HExists(repo.ctx, key, field).Result()
}

func (repo *abstractCacheRepositoryImpl[T]) NewPipeline() *CachePipeline {
	return &CachePipeline{
		pipe: repo.client.Pipeline(),
		ctx:  repo.ctx,
	}
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
func serialize[T any](value T) ([]byte, error) {
	if isPrimitiveType(value) {
		return []byte(fmt.Sprintf("%v", value)), nil
	}
	return sonic.Marshal(value)
}

// Helper function to deserialize data.
func deserialize[T any](data []byte, isPrimitive bool) (T, error) {
	var value T
	if isPrimitive {
		return any(string(data)).(T), nil
	}
	if err := sonic.Unmarshal(data, &value); err != nil {
		return value, fmt.Errorf("failed to deserialize value: %w", err)
	}
	return value, nil
}

func CreateCacheRepository[T any](redisClient *redis.Client, ctx context.Context, self AbstractCacheRepository[T]) *abstractCacheRepositoryImpl[T] {
	repo := &abstractCacheRepositoryImpl[T]{
		client:      redisClient,
		ctx:         ctx,
		isPrimitive: isPrimitiveType(new(T)),
		self:        self,
	}
	return repo
}
