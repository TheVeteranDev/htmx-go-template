package controllers

import (
	"net/http"

	"github.com/theveterandev/htmx-go-template/src/routes"
	"github.com/theveterandev/htmx-go-template/src/services"
)

type TemplateController struct {
	s *services.TemplateService
}

func InitTemplateController() routes.Routes {
	r := routes.Routes{}
	tc := TemplateController{services.InitTemplateService()}
	r[routes.Route{Path: "/", Method: "GET", Auth: false}] = tc.RenderIndex
	r[routes.Route{Path: "/navbar", Method: "GET", Auth: false}] = tc.RenderNavbar
	return r
}

func (c TemplateController) RenderIndex(w http.ResponseWriter, r *http.Request) {
	c.s.RenderIndex().Execute(w, nil)
}

func (c TemplateController) RenderNavbar(w http.ResponseWriter, r *http.Request) {
	c.s.RenderNavbar().Execute(w, nil)
}
