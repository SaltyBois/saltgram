package data

import (
	"crypto/rand"
	"encoding/binary"

	"gorm.io/gorm"
)

// Custom GORM base
type Identifiable struct {
	ID uint64 `gorm:"primaryKey;type:numeric" json:"id"`
}

func (i *Identifiable) BeforeCreate(tx *gorm.DB) error {
	i.ID = generateUint64()
	return nil
}

func generateUint64() uint64 {
	buf := make([]byte, 4)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}