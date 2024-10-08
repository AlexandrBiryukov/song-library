// models/song.go
package models

import "gorm.io/gorm"

type Song struct {
	gorm.Model
	GroupName   string
	Song        string
	ReleaseDate string
	Text        string
	Link        string
}
