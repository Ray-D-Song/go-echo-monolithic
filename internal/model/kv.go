package model

// Use a kv table to simulate a kv implementation with persistence
// similar to valkey and redis
type KV struct {
	BaseModel
	Key   string `json:"key" gorm:"uniqueIndex;not null"`
	Value string `json:"value"`
}
