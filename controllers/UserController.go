package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/theveterandev/htmx-go-template/database"
	"github.com/theveterandev/htmx-go-template/models"
	"github.com/theveterandev/htmx-go-template/routes"
	"github.com/theveterandev/htmx-go-template/services"
)

type UserController struct {
	s *services.UserService
	t *services.TemplateService
}

func InitUserController(db *database.Database) routes.Routes {
	r := routes.Routes{}
	uc := UserController{services.InitUserService(db), services.InitTemplateService()}
	r[routes.Route{Path: "/oauth/v1/list", Method: "GET", Auth: false}] = uc.ListUsers
	r[routes.Route{Path: "/oauth/v1/getByID/:id", Method: "GET", Auth: true}] = uc.GetUserByID
	r[routes.Route{Path: "/oauth/v1/getByUsername/:username", Method: "GET", Auth: true}] = uc.GetUserByUsername
	r[routes.Route{Path: "/oauth/v1/create", Method: "POST", Auth: false}] = uc.CreateUser
	r[routes.Route{Path: "/oauth/v1/delete/:id", Method: "DELETE", Auth: false}] = uc.DeleteUser
	r[routes.Route{Path: "/oauth/v1/signIn", Method: "POST", Auth: false}] = uc.SignIn
	r[routes.Route{Path: "/oauth/v1/signOut", Method: "POST", Auth: false}] = uc.SignOut
	r[routes.Route{Path: "/oauth/v1/signUp", Method: "POST", Auth: false}] = uc.SignUp
	return r
}

func (c UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.s.ListUsers()
	if err != nil {
		http.Error(w, "Error listing users", http.StatusInternalServerError)
		return
	}

	for _, u := range users {
		id := strconv.FormatInt(u.ID, 10)
		delete := fmt.Sprintf("/oauth/v1/delete/%s", id)
		data := map[string]interface{}{"Message": u.Username, "Delete": delete, "ID": u.ID}
		c.t.RenderNotification().Execute(w, data)
	}
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
	w.WriteHeader(http.StatusOK)
}

func (c UserController) SignIn(w http.ResponseWriter, r *http.Request) {
	user, err := c.s.SignIn(r.FormValue("username"), r.FormValue("password"))
	if err != nil {
		http.Error(w, "Invalid username and/or password", http.StatusBadRequest)
		return
	}
	c.t.RenderHomepage().Execute(w, user)
}

func (c UserController) SignOut(w http.ResponseWriter, r *http.Request) {
	c.t.RenderSignIn().Execute(w, nil)
}

func (c UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	if password != confirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	user, err := c.s.SignUp(username, password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	c.t.RenderHomepage().Execute(w, user)
}
