package models

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Task struct {
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Status        string     `json:"status"`
	Description   *string    `json:"description,omitempty"`
	Priority      string     `json:"priority"`
	ProjectID     string     `json:"project_id"`
	AssigneeID    *string    `json:"assignee_id,omitempty"`
	AssigneeName  *string    `json:"assignee_name,omitempty"`
	AssigneeEmail *string    `json:"assignee_email,omitempty"`
	DueDate       *time.Time `json:"due_date,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ProjectStats struct {
	ByStatus   map[string]int `json:"by_status"`
	ByAssignee map[string]int `json:"by_assignee"`
}
