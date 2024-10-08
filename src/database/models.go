package database

import "time"

type User struct {
	ID        string    `json:"id" db:"id"` // UUID
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Verified  bool      `json:"verified" db:"verified"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}

type Profile struct {
	ID        string    `json:"id" db:"id"` // UUID
	FirstName string    `json:"firstName" db:"firstName"`
	LastName  string    `json:"lastName" db:"lastName"`
	Image     string    `json:"image" db:"image"`
	UserID    string    `json:"userId" db:"userId"` // Foreign key to User
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}

type Blog struct {
	ID        string    `json:"id" db:"id"`               // UUID
	ProfileID string    `json:"profileId" db:"profileId"` // Foreign key to User
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	Tags      []string  `json:"tags" db:"tags"` // JSON array for tags
	Published bool      `json:"published" db:"published"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"updatedAt"`
}

type Otp struct {
	ID        string    `json:"id" db:"id"` // UUID
	Email     string    `json:"email" db:"email"`
	Otp       string    `json:"otp" db:"otp"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
}
