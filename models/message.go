package models

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func generateRandomID() string {
	b := make([]byte, 12)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

type Message struct {
	MessageID        string         `gorm:"type:varchar(16);primary_key" json:"message_id"`
	Password         string         `gorm:"type:varchar(255)" json:"password"`
	Recipient        string         `gorm:"type:varchar(255)" json:"recipient"`
	ProfileSticker   int            `gorm:"type:int" json:"profile_sticker"`
	InventorySticker int            `gorm:"type:int" json:"inventory_sticker"`
	Message          string         `gorm:"type:text" json:"message"`
	ImageURL         pq.StringArray `gorm:"type:text[]" json:"image_url"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.MessageID == "" {
		m.MessageID = generateRandomID()
	}
	return nil
}
