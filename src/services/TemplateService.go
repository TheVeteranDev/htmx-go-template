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
