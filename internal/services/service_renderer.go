package services

import (
	"github.com/a-h/templ"
	"github.com/quangd42/meal-planner/internal/views"
)

type Renderer struct{}

var privateNavItems = []views.NavItem{
	{
		Name: "Home",
		URL:  templ.URL("/"),
	},
	{
		Name: "Add Recipe",
		URL:  templ.URL("/recipes/add"),
	},
}

var publicNavItems = []views.NavItem{
	{
		Name: "Login",
		URL:  templ.URL("/login"),
	},
	{
		Name: "Register",
		URL:  templ.URL("/register"),
	},
}

func NewRendererService() Renderer {
	return Renderer{}
}

func (rs Renderer) GetNavItems(isLoggedIn bool) []views.NavItem {
	if isLoggedIn {
		return privateNavItems
	} else {
		return publicNavItems
	}
}
