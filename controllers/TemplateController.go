package controllers

import (
	"net/http"

	"github.com/theveterandev/htmx-go-template/routes"
	"github.com/theveterandev/htmx-go-template/services"
)

type TemplateController struct {
	s *services.TemplateService
}

func InitTemplateController() routes.Routes {
	r := routes.Routes{}
	tc := TemplateController{services.InitTemplateService()}
	r[routes.Route{Path: "/", Method: "GET", Auth: false}] = tc.RenderIndex
	r[routes.Route{Path: "/navbar", Method: "GET", Auth: false}] = tc.RenderNavbar
	r[routes.Route{Path: "/homepage", Method: "GET", Auth: false}] = tc.RenderHomepage
	r[routes.Route{Path: "/sign-in", Method: "GET", Auth: false}] = tc.RenderSignIn
	r[routes.Route{Path: "/sign-up", Method: "GET", Auth: false}] = tc.RenderSignUp
	return r
}

func (c TemplateController) RenderIndex(w http.ResponseWriter, r *http.Request) {
	c.s.RenderIndex().Execute(w, nil)
}

func (c TemplateController) RenderNavbar(w http.ResponseWriter, r *http.Request) {
	c.s.RenderNavbar().Execute(w, nil)
}

func (c TemplateController) RenderHomepage(w http.ResponseWriter, r *http.Request) {
	c.s.RenderHomepage().Execute(w, nil)
}

func (c TemplateController) RenderSignIn(w http.ResponseWriter, r *http.Request) {
	c.s.RenderSignIn().Execute(w, nil)
}

func (c TemplateController) RenderSignUp(w http.ResponseWriter, r *http.Request) {
	c.s.RenderSignUp().Execute(w, nil)
}
