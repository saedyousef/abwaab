// models/users.go

package models

type User struct {
	ID       uint64  `json:"id" gorm:"primary_key"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
