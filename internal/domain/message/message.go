package message

import "time"

// MessageStatus represents the status of a message
type MessageStatus string

const (
	MessageStatusQueued     MessageStatus = "queued"     // Kuyruğa alındı, gönderilmeyi bekliyor
	MessageStatusProcessing MessageStatus = "processing" // Şu an gönderiliyor
	MessageStatusSent       MessageStatus = "sent"       // Provider'a başarıyla iletildi
	MessageStatusDelivered  MessageStatus = "delivered"  // Alıcıya ulaştı (SMS/Email için)
	MessageStatusFailed     MessageStatus = "failed"     // Gönderim başarısız
	MessageStatusCancelled  MessageStatus = "cancelled"  // İptal edildi
)

// Message represents a message entity in the domain
type Message struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	To        string        `gorm:"not null" json:"to"`
	Content   string        `gorm:"not null" json:"content"`
	Status    MessageStatus `gorm:"type:varchar(20);default:'queued'" json:"status"`
	MessageID string        `gorm:"type:varchar(255)" json:"message_id,omitempty"`

	// Retry tracking
	RetryCount int `gorm:"default:0" json:"retry_count,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName specifies the table name for Message
func (Message) TableName() string {
	return "messages"
}

// IsValidContent checks if message content is within character limit
func (m *Message) IsValidContent(maxLength int) bool {
	return len(m.Content) > 0 && len(m.Content) <= maxLength
}
