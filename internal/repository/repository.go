package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewRedisClient(url string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}

func NewPostgresDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"uniqueIndex;size:255" json:"email"`
	Password  string    `gorm:"size:255" json:"-"`
	Name      string    `gorm:"size:100" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type APIKey struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Key       string    `gorm:"uniqueIndex;size:64" json:"key"`
	Name      string    `gorm:"size:100" json:"name"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Quota     int64     `gorm:"default:0" json:"quota"`
	ExpiresAt *time.Time `json:"expires_at"`
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type Item struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:255" json:"name"`
	Description string  `gorm:"text" json:"description"`
	UserID    uint      `gorm:"index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type KeyStore interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id uint) (*User, error)
	CreateAPIKey(key *APIKey) error
	ValidateAPIKey(key string) (*APIKey, error)
	GetUserItems(userID uint) ([]Item, error)
	CreateItem(item *Item) error
	GetItem(id uint) (*Item, error)
	UpdateItem(item *Item) error
	DeleteItem(id uint) error
}

type postgresKeyStore struct {
	db *gorm.DB
}

func NewKeyStore(db *gorm.DB) KeyStore {
	return &postgresKeyStore{db: db}
}

func (k *postgresKeyStore) CreateUser(user *User) error {
	return k.db.Create(user).Error
}

func (k *postgresKeyStore) GetUserByEmail(email string) (*User, error) {
	var user User
	err := k.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (k *postgresKeyStore) GetUserByID(id uint) (*User, error) {
	var user User
	err := k.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (k *postgresKeyStore) CreateAPIKey(key *APIKey) error {
	return k.db.Create(key).Error
}

func (k *postgresKeyStore) ValidateAPIKey(keyStr string) (*APIKey, error) {
	var key APIKey
	err := k.db.Where("key = ? AND active = true", keyStr).First(&key).Error
	if err != nil {
		return nil, err
	}
	if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("key expired")
	}
	return &key, nil
}

func (k *postgresKeyStore) GetUserItems(userID uint) ([]Item, error) {
	var items []Item
	err := k.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error
	return items, err
}

func (k *postgresKeyStore) CreateItem(item *Item) error {
	return k.db.Create(item).Error
}

func (k *postgresKeyStore) GetItem(id uint) (*Item, error) {
	var item Item
	err := k.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (k *postgresKeyStore) UpdateItem(item *Item) error {
	return k.db.Save(item).Error
}

func (k *postgresKeyStore) DeleteItem(id uint) error {
	return k.db.Delete(&Item{}, id).Error
}
