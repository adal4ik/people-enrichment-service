package models

import (
	"time"

	"github.com/google/uuid"
)

type Person struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`        // обязательный
	Surname     string    `json:"surname"`     // обязательный
	Patronymic  *string   `json:"patronymic"`  // необязательный
	Age         *int      `json:"age"`         // необязательный
	Gender      *string   `json:"gender"`      // необязательный
	Nationality *string   `json:"nationality"` // необязательный
	CreatedAt   time.Time `json:"created_at"`  // дата создания
	UpdatedAt   time.Time `json:"updated_at"`  // дата обновления
}
