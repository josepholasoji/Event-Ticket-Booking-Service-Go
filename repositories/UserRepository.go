package repositories

import (
	"bookingservice/initializations"
	"bookingservice/models"
)

type UserRepositoryInterface interface {
	GetAllUsers() ([]models.User, error)
	FindUserByName(username string) (*models.User, error)
	SaveUser(name string, username string, password string) error
}

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	rows, err := initializations.MySQLDB.Query("SELECT id, name, email, role FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.UserName, &user.Role); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) FindUserByName(username string) (*models.User, error) {
	var user = &models.User{}
	if err := initializations.MySQLDB.QueryRow("SELECT id, name, email, password, role FROM users WHERE email = ?", username).Scan(&user.ID, &user.Name, &user.UserName, &user.Password, &user.Role); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) SaveUser(name string, username string, password string) error {
	if _, err := initializations.MySQLDB.Exec("INSERT INTO users (name, username, password) VALUES (?, ?, ?)", name, username, password); err != nil {
		return err
	}
	return nil
}
