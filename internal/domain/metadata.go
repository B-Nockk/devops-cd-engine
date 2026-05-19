package domain

import "time"

type MetaData struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Note      string
}

func NewMetaData(note string) MetaData {
	now := time.Now()
	return MetaData{
		CreatedAt: now,
		UpdatedAt: now,
		Note:      note,
	}
}

func (m *MetaData) Touch(note string) {
	m.UpdatedAt = time.Now()
	if note != "" {
		m.Note = note
	}
}
