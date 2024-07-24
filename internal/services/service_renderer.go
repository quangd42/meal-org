package services

import (
	"github.com/a-h/templ"
	"github.com/quangd42/meal-planner/internal/views"
)

type Renderer struct{}

var privateNavItems = []views.NavItem{
	{
		Name: "Home",
		URL:  "/",
	},
	{
		Name: "Add Recipe",
		URL:  "/recipes/add",
	},
	{
		Name: "Logout",
		URL:  "#",
	},
}

var publicNavItems = []views.NavItem{
	{
		Name: "Login",
		URL:  "/login",
	},
	{
		Name: "Register",
		URL:  "/register",
	},
}

func NewRendererService() Renderer {
	return Renderer{}
}

func (rs Renderer) GetNavItems(isLoggedIn bool, url string) []views.NavItem {
	var srcItems []views.NavItem
	var items []views.NavItem

	if isLoggedIn {
		srcItems = privateNavItems
	} else {
		srcItems = publicNavItems
	}
	for _, i := range srcItems {
		if templ.URL(i.URL) == templ.URL(url) {
			i.IsCurrent = true
		}
		items = append(items, i)
	}

	return items
}
