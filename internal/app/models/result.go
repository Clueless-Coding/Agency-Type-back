package models

import "time"

type Result struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	GameMode  string    `json:"game_mode"`
	StartTime time.Time `json:"start_time"`
	Duration  time.Time `json:"duration"`
	Mistakes  int       `json:"misstakes"`
	Accuracy  float64   `json:"accuracy"`
	Words     int       `json:"count_words"`
	WPN       float64   `json:"wpn"`
	CPN       float64   `json:"cpn"`
}
