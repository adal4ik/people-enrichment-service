package models

import "time"

type Person struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Patronymic  string    `json:"patronymic"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	Nationality string    `json:"nationality"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
