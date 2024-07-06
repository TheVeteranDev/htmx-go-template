package services

import "text/template"

type TemplateService struct{}

func InitTemplateService() *TemplateService {
	return &TemplateService{}
}

func (s TemplateService) RenderIndex() *template.Template {
	return template.Must(template.ParseFiles("templates/index.html"))
}

func (s TemplateService) RenderNavbar() *template.Template {
	return template.Must(template.ParseFiles("templates/fragments/navbar.html"))
}

func (s TemplateService) RenderHomepage() *template.Template {
	return template.Must(template.ParseFiles("templates/fragments/homepage.html"))
}

func (s TemplateService) RenderSignIn() *template.Template {
	return template.Must(template.ParseFiles("templates/fragments/sign-in.html"))
}

func (s TemplateService) RenderSignUp() *template.Template {
	return template.Must(template.ParseFiles("templates/fragments/sign-up.html"))
}
