package repository

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// setup function to initialize Miniredis and Redis client
func setup(t *testing.T) (*miniredis.Miniredis, *redis.Client, RedisRepository) {
	mr := miniredis.RunT(t) // Inicia Miniredis
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	repo := NewRedisRepository(rdb)
	return mr, rdb, repo
}

// teardown function to clean up resources
func teardown(mr *miniredis.Miniredis, rdb *redis.Client) {
	rdb.Close()
	mr.Close()
}

func TestRedisRepository_Get(t *testing.T) {
	t.Run("Test Get with Miniredis", func(t *testing.T) {
		mr, rdb, repo := setup(t)
		defer teardown(mr, rdb)

		tests := []struct {
			name         string
			key          string
			setupFunc    func()
			expected     string
			expectError  bool
		}{
			{
				name: "key exists",
				key:  "existingKey",
				setupFunc: func() {
					mr.Set("existingKey", "testValue")
				},
				expected:    "testValue",
				expectError: false,
			},
			{
				name: "key does not exist",
				key:  "nonExistingKey",
				setupFunc: func() {
					// No key is set in this case
				},
				expected:    "",
				expectError: true,
			},
			{
				name: "key expired",
				key:  "expiredKey",
				setupFunc: func() {
					mr.Set("expiredKey", "testValue")
					mr.SetTTL("expiredKey", time.Minute) // Set a 1 minute TTL
					mr.FastForward(2 * time.Minute) // Simulate Time progresses for the key to expire
				},
				expected:    "",
				expectError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx := context.Background()

				if tt.setupFunc != nil {
					tt.setupFunc() // Execute the specific configuration function of the case
				}

				result, err := repo.Get(ctx, tt.key)

				if tt.expectError {
					assert.Error(t, err)
					assert.Equal(t, "", result)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expected, result)
				}
			})
		}
	})
}

func TestRedisRepository_Set(t *testing.T) {
	t.Run("Test Set with Miniredis", func(t *testing.T) {
		mr, rdb, repo := setup(t)
		defer teardown(mr, rdb)

		tests := []struct {
			name         string
			key          string
			value        string
			expiration   time.Duration
			expected     string
			expectError  bool
		}{
			{
				name:       "set key successfully",
				key:        "newKey",
				value:      "newValue",
				expiration: time.Minute,
				expected:   "newValue",
				expectError: false,
			},
			{
				name:       "set key with no expiration",
				key:        "noExpireKey",
				value:      "noExpireValue",
				expiration: 0, // No expiration
				expected:   "noExpireValue",
				expectError: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx := context.Background()

				err := repo.Set(ctx, tt.key, tt.value, tt.expiration)
				assert.NoError(t, err)

				// Verify if the value was stored correctly in miniredis
				gotValue, _ := mr.Get(tt.key)
				assert.Equal(t, tt.expected, gotValue)

				if tt.expiration > 0 {
					mr.FastForward(tt.expiration)
					gotValue, _ = mr.Get(tt.key)
					assert.Equal(t, "", gotValue) // The key should have expired
				}
			})
		}
	})
}

func TestRedisRepository_Exists(t *testing.T) {
	t.Run("Test Exists with Miniredis", func(t *testing.T) {
		mr, rdb, repo := setup(t)
		defer teardown(mr, rdb)

		tests := []struct {
			name         string
			key          string
			setupFunc    func()
			expected     bool
			expectError  bool
		}{
			{
				name: "key exists",
				key:  "existingKey",
				setupFunc: func() {
					mr.Set("existingKey", "someValue")
				},
				expected:    true,
				expectError: false,
			},
			{
				name:        "key does not exist",
				key:         "nonExistingKey",
				setupFunc:   nil,
				expected:    false,
				expectError: false,
			},
			{
				name: "key expired",
				key:  "expiredKey",
				setupFunc: func() {
					mr.Set("expiredKey", "someValue")
					mr.SetTTL("expiredKey", time.Minute)
					mr.FastForward(2 * time.Minute)
				},
				expected:    false,
				expectError: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx := context.Background()

				if tt.setupFunc != nil {
					tt.setupFunc()
				}

				exists, err := repo.Exists(ctx, tt.key)

				if tt.expectError {
					assert.Error(t, err)
					assert.False(t, exists)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.expected, exists)
					
					if !tt.expected && mr.Exists(tt.key) {
						t.Errorf("expected key '%s' to be expired or non-existent but it still exists", tt.key)
					}
				}
			})
		}
	})
}


func TestRedisRepository_Delete(t *testing.T) {
	t.Run("Test Delete with Miniredis", func(t *testing.T) {
		mr, rdb, repo := setup(t)
		defer teardown(mr, rdb)

		tests := []struct {
			name         string
			key          string
			setupFunc    func()
			expectError  bool
			expectDelete bool
		}{
			{
				name: "delete existing key",
				key:  "existingKey",
				setupFunc: func() {
					mr.Set("existingKey", "valueToDelete")
				},
				expectError: false,
				expectDelete: true,
			},
			{
				name:        "delete non-existing key",
				key:         "nonExistingKey",
				setupFunc:   nil,
				expectError: false,
				expectDelete: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx := context.Background()

				if tt.setupFunc != nil {
					tt.setupFunc()
				}

				err := repo.Delete(ctx, tt.key)

				if tt.expectError {
					assert.Error(t, err)
				} else {

					assert.NoError(t, err)
					exists := mr.Exists(tt.key)

					if tt.expectDelete {
						assert.False(t, exists)
					} else {
						assert.False(t, exists)
					}
				}
			})
		}
	})
}