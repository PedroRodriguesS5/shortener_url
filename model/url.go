package model

import "time"

type URLMapping struct {
	URL       string    `json:"url"`
	Custom    string    `json:"custom,omitempty"`
	ExpiresIn int       `json:"expires_in,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"` //minutes
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}
