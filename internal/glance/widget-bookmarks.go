package glance

import (
	"html/template"
)

var bookmarksWidgetTemplate = mustParseTemplate("bookmarks.html", "widget-base.html")

type bookmarksWidget struct {
	widgetBase `yaml:",inline" json:",inline"`
	cachedHTML template.HTML `yaml:"-" json:"-"`
	Groups     []struct {
		Title     string         `yaml:"title" json:"title"`
		Color     *hslColorField `yaml:"color" json:"color"`
		SameTab   bool           `yaml:"same-tab" json:"same-tab"`
		HideArrow bool           `yaml:"hide-arrow" json:"hide-arrow"`
		Target    string         `yaml:"target" json:"target"`
		Links     []struct {
			Title       string          `yaml:"title" json:"title"`
			URL         string          `yaml:"url" json:"url"`
			Description string          `yaml:"description" json:"description"`
			Icon        customIconField `yaml:"icon" json:"icon"`
			// we need a pointer to bool to know whether a value was provided,
			// however there's no way to dereference a pointer in a template so
			// {{ if not .SameTab }} would return true for any non-nil pointer
			// which leaves us with no way of checking if the value is true or
			// false, hence the duplicated fields below
			SameTabRaw   *bool  `yaml:"same-tab" json:"same-tab"`
			SameTab      bool   `yaml:"-" json:"-"`
			HideArrowRaw *bool  `yaml:"hide-arrow" json:"hide-arrow"`
			HideArrow    bool   `yaml:"-" json:"-"`
			Target       string `yaml:"target" json:"target"`
		} `yaml:"links" json:"links"`
	} `yaml:"groups" json:"groups"`
}

func (widget *bookmarksWidget) initialize() error {
	widget.withTitle("Bookmarks").withError(nil)

	for g := range widget.Groups {
		group := &widget.Groups[g]
		for l := range group.Links {
			link := &group.Links[l]
			if link.SameTabRaw == nil {
				link.SameTab = group.SameTab
			} else {
				link.SameTab = *link.SameTabRaw
			}

			if link.HideArrowRaw == nil {
				link.HideArrow = group.HideArrow
			} else {
				link.HideArrow = *link.HideArrowRaw
			}

			if link.Target == "" {
				if group.Target != "" {
					link.Target = group.Target
				} else {
					if link.SameTab {
						link.Target = ""
					} else {
						link.Target = "_blank"
					}
				}
			}
		}
	}

	widget.cachedHTML = widget.renderTemplate(widget, bookmarksWidgetTemplate)

	return nil
}

func (widget *bookmarksWidget) Render() template.HTML {
	return widget.cachedHTML
}
