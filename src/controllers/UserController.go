package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/theveterandev/htmx-go-template/database"
	"github.com/theveterandev/htmx-go-template/src/models"
	"github.com/theveterandev/htmx-go-template/src/routes"
	"github.com/theveterandev/htmx-go-template/src/services"
)

type UserController struct {
	s *services.UserService
}

func InitUserController(db *database.Database) routes.Routes {
	r := routes.Routes{}
	uc := UserController{services.InitUserService(db)}
	r[routes.Route{Path: "/oauth/v1", Method: "GET", Auth: true}] = uc.ListUsers
	r[routes.Route{Path: "/oauth/v1/getByID/:id", Method: "GET", Auth: true}] = uc.GetUserByID
	r[routes.Route{Path: "/oauth/v1/getByUsername/:username", Method: "GET", Auth: true}] = uc.GetUserByUsername
	r[routes.Route{Path: "/oauth/v1/create", Method: "POST", Auth: false}] = uc.CreateUser
	r[routes.Route{Path: "/oauth/v1/:id", Method: "DELETE", Auth: true}] = uc.DeleteUser
	r[routes.Route{Path: "/oauth/v1/signIn", Method: "POST", Auth: false}] = uc.SignIn
	r[routes.Route{Path: "/oauth/v1/signOut", Method: "POST", Auth: false}] = uc.SignOut
	r[routes.Route{Path: "/oauth/v1/signUp", Method: "POST", Auth: false}] = uc.SignUp
	return r
}

func (c UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := c.s.ListUsers()
	if err != nil {
		http.Error(w, "Error listing users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (c UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing id query parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id query parameter. Must be a number.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	user, err := c.s.GetUserByID(id)
	if err != nil {
		http.Error(w, "ID does not exist", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (c UserController) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		http.Error(w, "Missing username query parameter", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	user, err := c.s.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Username does not exist", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (c UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	createdUser, err := c.s.CreateUser(user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(createdUser)
}

func (c UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing user ID.", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID. Must be a number.", http.StatusBadRequest)
		return
	}
	err = c.s.DeleteUser(id)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c UserController) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	user, token, err := c.s.SignIn(u.Username, u.Password)
	if err != nil {
		http.Error(w, "Invalid username and/or password", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"user": user, "token": token})
}

func (c UserController) SignOut(w http.ResponseWriter, r *http.Request) {
	success := c.s.SignOut()
	if !success {
		http.Error(w, "Error signing out user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	u := models.User{}
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	user, token, err := c.s.SignUp(u.Username, u.Password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"user": user, "token": token})
}
