package model

// Rider represents a rider in the system.
type Rider struct {
	ID       int    `json:"id" gorm:"primaryKey autoIncrement not null"`
	Name     string `json:"name" gorm:"not null"`
	Assigned bool   `json:"assigned" gorm:"not null default false"`
}
