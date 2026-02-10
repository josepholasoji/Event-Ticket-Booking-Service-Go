package models

type User struct {
	ID       uint   `db:"id"`
	Name     string `db:"not null"`
	UserName string `db:"email;not null"`
	Password string `db:"password_hash;not null"`
	Role     string `db:"role;not null"`
}
