package routes

import (
	"net/http"
	"strings"

	"github.com/theveterandev/htmx-go-template/jwt"
)

type Route struct {
	Path   string
	Method string
	Auth   bool
}

type Routes map[Route]http.HandlerFunc

func (routes Routes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routeExists := false

	for route, handler := range routes {
		path := strings.Split(route.Path, "/:")
		pathVarsStr := strings.Replace(r.URL.Path, path[0], "", 1)
		pathVars := removeEmptyStrings(strings.Split(pathVarsStr, "/"))

		validPath := strings.HasPrefix(r.URL.Path, path[0])
		validPathVars := len(pathVars) == len(path[1:])
		validMethod := r.Method == route.Method

		if validPath {
			routeExists = true
		}

		if validPath && validPathVars && validMethod {
			for i, pathVar := range pathVars {
				r.SetPathValue(path[1:][i], pathVar)
			}
			if route.Auth {
				token := r.Header.Get("Authorization")
				claim := jwt.ValidateToken(token)
				if claim == nil || !validateRole(claim.Roles, []string{"admin"}) {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

			}
			handler(w, r)
			return
		}
	}
	if routeExists {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	} else {
		http.NotFound(w, r)
	}
}

func validateRole(roles []string, requiredRoles []string) bool {
	for _, role := range roles {
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				return true
			}
		}
	}
	return false
}

func removeEmptyStrings(slice []string) []string {
	var result []string
	for _, s := range slice {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
