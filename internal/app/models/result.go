package models

import "time"

type Result struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	GameMode  string    `json:"game_mode"`
	StartTime time.Time `json:"start_time"`
	Duration  float64   `json:"duration"`
	Mistakes  int       `json:"mistakes"`
	Accuracy  float64   `json:"accuracy"`
	Words     int       `json:"count_words"`
	WPM       float64   `json:"wpm"`
	CPM       float64   `json:"cpm"`
}
