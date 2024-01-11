package models

type Connection struct {
	ID    string `gorm:"primaryKey"`
	Src   string
	Dst   string
	Port  int
    Count int
}
