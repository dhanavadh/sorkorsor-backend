package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Message struct {
	MessageID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"message_id"`
	Password         string         `gorm:"type:varchar(255)" json:"password"`
	Recipient        string         `gorm:"type:varchar(255)" json:"recipient"`
	ProfileSticker   int            `gorm:"type:int" json:"profile_sticker"`
	InventorySticker int            `gorm:"type:int" json:"inventory_sticker"`
	Message          string         `gorm:"type:text" json:"message"`
	ImageURL         pq.StringArray `gorm:"type:text[]" json:"image_url"`
}
