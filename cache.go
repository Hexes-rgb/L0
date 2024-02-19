package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type DataToCache struct {
	OrderUID  string
	OrderData []byte
	CreatedAt time.Time
}

func NewDataToCache(orderId string, orderData []byte, createdAt time.Time) *DataToCache {
	return &DataToCache{
		OrderUID:  orderId,
		OrderData: orderData,
		CreatedAt: createdAt,
	}
}

type CachedData struct {
	OrderData    []byte
	CreatedAt    time.Time
	LastAccessAt time.Time
}

type Cache struct {
	mu      sync.RWMutex
	items   map[string]*CachedData
	maxSize int
}

func NewCache(maxSize int) *Cache {
	return &Cache{
		items:   make(map[string]*CachedData),
		maxSize: maxSize,
	}
}

func (c *Cache) InitializeCache() error {
	query := fmt.Sprintf("SELECT order_uid, order_data, created_at FROM orders ORDER BY created_at DESC LIMIT %d", config.MaxCacheSize)

	rows, err := dbpool.Query(context.Background(), query)
	if err != nil {
		return err
	}
	defer rows.Close()

	var dataToCache []*DataToCache

	for rows.Next() {
		var orderUID string
		var orderData []byte
		var createdAt time.Time
		if err := rows.Scan(&orderUID, &orderData, &createdAt); err != nil {
			return err
		}
		dataToCache = append(dataToCache, NewDataToCache(orderUID, orderData, createdAt))
	}

	if err := rows.Err(); err != nil {
		return err
	}

	for i := len(dataToCache) - 1; i >= 0; i-- {
		cache.Add(dataToCache[i])
	}

	return nil
}

func (c *Cache) Add(data *DataToCache) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.items) >= c.maxSize {
		c.evictOldest()
	}

	cachedData := &CachedData{
		OrderData:    data.OrderData,
		CreatedAt:    data.CreatedAt,
		LastAccessAt: time.Now(),
	}
	c.items[data.OrderUID] = cachedData
}

func (c *Cache) Get(orderUID string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cachedData, ok := c.items[orderUID]; ok {
		cachedData.LastAccessAt = time.Now()
		return cachedData.OrderData, true
	}

	return nil, false
}

func (c *Cache) evictOldest() {
	var oldestOrderUID string
	var oldestCreatedAt time.Time
	var oldestAccessAt time.Time
	for orderUID, cachedData := range c.items {
		if oldestAccessAt.IsZero() ||
			cachedData.LastAccessAt.Before(oldestAccessAt) ||
			(cachedData.LastAccessAt.Equal(oldestAccessAt) && cachedData.CreatedAt.Before(oldestCreatedAt)) {
			oldestOrderUID = orderUID
			oldestAccessAt = cachedData.LastAccessAt
			oldestCreatedAt = cachedData.CreatedAt
		}
	}
	delete(c.items, oldestOrderUID)
}
