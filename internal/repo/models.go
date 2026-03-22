package repo

import "time"

type Content struct {
	ID           int64     `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	CanvaURL     string    `json:"canva_url" db:"canva_url"`
	Class        int       `json:"class" db:"class"`
	Quarter      int       `json:"quarter" db:"quarter"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	LessonNumber int       `json:"lesson_number" json:"lesson_number"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
