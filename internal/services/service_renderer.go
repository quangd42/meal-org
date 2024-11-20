package services

import (
	"github.com/a-h/templ"
	"github.com/quangd42/meal-planner/internal/models"
)

type Renderer struct{}

var privateNavItems = []models.NavItem{
	{
		Link: models.Link{
			Name: "Home",
			URL:  "/",
		},
	},
	{
		Link: models.Link{
			Name: "All Recipes",
			URL:  "/recipes",
		},
	},
	{
		Link: models.Link{
			Name: "Add Recipe",
			URL:  "/recipes/add",
		},
	},
	{
		Link: models.Link{
			Name: "Logout",
			URL:  "/logout",
		},
		IsPostRequest: true,
	},
}

var publicNavItems = []models.NavItem{
	{
		Link: models.Link{
			Name: "Login",
			URL:  "/login",
		},
	},
	{
		Link: models.Link{
			Name: "Register",
			URL:  "/register",
		},
	},
}

func NewRendererService() Renderer {
	return Renderer{}
}

func (rs Renderer) GetNavItems(isLoggedIn bool, url string) []models.NavItem {
	var srcItems []models.NavItem
	var items []models.NavItem

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
