package entities

type Result struct {
	Id    string `gorm:"primaryKey"`
	Value string
}
