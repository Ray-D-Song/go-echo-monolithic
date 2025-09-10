package repository

import (
	"errors"

	"github.com/ray-d-song/go-echo-monolithic/internal/model"
	"gorm.io/gorm"
)

// KVRepository handles key-value data operations
type KVRepository struct {
	db *gorm.DB
}

// NewKVRepository creates a new KV repository
func NewKVRepository(db *gorm.DB) *KVRepository {
	return &KVRepository{db: db}
}

// Set stores a key-value pair
func (r *KVRepository) Set(key, value string) error {
	kv := &model.KV{
		Key:   key,
		Value: value,
	}
	
	// Use GORM's Upsert functionality
	return r.db.Save(kv).Error
}

// Get retrieves a value by key
func (r *KVRepository) Get(key string) (string, error) {
	var kv model.KV
	if err := r.db.Where("key = ?", key).First(&kv).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil // Return empty string for non-existent keys
		}
		return "", err
	}
	return kv.Value, nil
}

// Delete removes a key-value pair
func (r *KVRepository) Delete(key string) error {
	result := r.db.Where("key = ?", key).Delete(&model.KV{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Exists checks if a key exists
func (r *KVRepository) Exists(key string) (bool, error) {
	var count int64
	err := r.db.Model(&model.KV{}).Where("key = ?", key).Count(&count).Error
	return count > 0, err
}

// GetAll retrieves all key-value pairs
func (r *KVRepository) GetAll() ([]model.KV, error) {
	var kvs []model.KV
	err := r.db.Find(&kvs).Error
	return kvs, err
}

// GetKeys retrieves all keys
func (r *KVRepository) GetKeys() ([]string, error) {
	var keys []string
	err := r.db.Model(&model.KV{}).Pluck("key", &keys).Error
	return keys, err
}

// Clear removes all key-value pairs
func (r *KVRepository) Clear() error {
	return r.db.Where("1 = 1").Delete(&model.KV{}).Error
}