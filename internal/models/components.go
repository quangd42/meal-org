package models

type Link struct {
	Name       string
	URL        string
	IsExternal bool
}

type NavItem struct {
	Link
	IsButton      bool
	IsCurrent     bool
	IsPostRequest bool
}
