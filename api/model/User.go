package model

type User struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Username   string `json:"username" gorm:"unique"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
	StudentID  string `json:"student_id"`
	Faculty    string `json:"faculty"`
	Specialty  string `json:"specialty"`
	GroupName  string `json:"group_name"`
	Course     int    `json:"course"`
	Photo      string `json:"photo"`
}
