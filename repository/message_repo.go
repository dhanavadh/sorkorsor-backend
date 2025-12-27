package repository

import (
	"fmt"

	"github.com/dhanavadh/sorkorsor-backend/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) FindAll() ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Find(&messages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}
	return messages, nil
}

func (r *MessageRepository) FindByID(id string) (*models.Message, error) {
	var message models.Message
	err := r.db.Where("message_id = ?", id).First(&message).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find message: %w", err)
	}
	return &message, nil
}

func (r *MessageRepository) Create(message *models.Message) error {
	err := r.db.Create(message).Error
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}
	return nil
}
