package services

import (
	"strings"

	"github.com/theveterandev/htmx-go-template/database"
	"github.com/theveterandev/htmx-go-template/src/jwt"
	"github.com/theveterandev/htmx-go-template/src/models"
	utils "github.com/theveterandev/htmx-go-template/src/utilities"
)

type UserService struct {
	db *database.Database
}

func InitUserService(db *database.Database) *UserService {
	return &UserService{db}
}

func (s UserService) ListUsers() ([]models.User, error) {
	data, err := s.db.Query("SELECT id, username FROM users")
	if err != nil {
		return nil, err
	}
	defer data.Close()

	users := []models.User{}
	for data.Next() {
		u := models.User{}
		err := data.Scan(&u.ID, &u.Username)
		if err != nil {
			return nil, err
		}
		u.Password = "**********"
		users = append(users, u)
	}

	return users, nil
}

func (s UserService) GetUserByID(id int64) (models.User, error) {
	u := models.User{}
	err := s.db.QueryRow("SELECT id, username FROM users WHERE id = $1", id).Scan(&u.ID, &u.Username)
	u.Password = "**********"
	return u, err
}

func (s UserService) GetUserByUsername(username string) (models.User, error) {
	u := models.User{}
	err := s.db.QueryRow("SELECT id, username FROM users WHERE username = $1", strings.ToLower(username)).Scan(&u.ID, &u.Username)
	return u, err
}

func (s UserService) getUserByUsername(username string) (models.User, error) {
	u := models.User{}
	err := s.db.QueryRow("SELECT * FROM users WHERE username = $1", strings.ToLower(username)).Scan(&u.ID, &u.Username, &u.Password)
	return u, err
}

func (s UserService) CreateUser(u models.User) (models.User, error) {
	pwHash, err := utils.Encrypt(u.Password)
	if err != nil {
		return models.User{}, err
	}
	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING *"
	err = s.db.QueryRow(query, strings.ToLower(u.Username), pwHash).Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (s UserService) DeleteUser(id int64) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (s UserService) SignIn(username string, password string) (models.User, map[string]string, error) {
	u, err := s.getUserByUsername(username)
	if err != nil {
		return u, nil, err
	}

	err = utils.CompareHash(u.Password, password)
	if err != nil {
		return u, nil, err
	}

	token, err := jwt.GenerateToken(u.Username)
	if err != nil {
		return u, nil, err
	}

	u.Password = "**********"

	return u, map[string]string{"token": token}, nil
}

func (s UserService) SignOut() bool {
	return true
}

func (s UserService) SignUp(username string, password string) (models.User, map[string]string, error) {
	u := models.User{Username: username, Password: password}
	createdUser, err := s.CreateUser(u)
	if err != nil {
		return u, nil, err
	}

	token, err := jwt.GenerateToken(createdUser.Username)
	if err != nil {
		return u, nil, err
	}

	createdUser.Password = "**********"

	return createdUser, map[string]string{"token": token}, nil
}
