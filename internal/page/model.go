package page

import (
	"time"
)

type Page struct {
	ID           string    `json:"id"`
	CreateTime   time.Time `json:"create_time"`
	UserName     string    `json:"user_name"`
	URL          string    `json:"url,omitempty"`
	Description  string    `json:"description,omitempty"`
	Category     string    `json:"category,omitempty"`
	Price        int       `json:"price,omitempty"`
	TimeDuration int       `json:"time_duration,omitempty"`
}
