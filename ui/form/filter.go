package form

import (
	"net/http"

	"github.com/miniflux/miniflux/errors"
	"github.com/miniflux/miniflux/logger"
	"strings"
)

// FilterForm represents a filter form in the UI
type FilterForm struct {
	FilterName string
	Filters    []string
}

// Validate makes sure the form values are valid.
func (c FilterForm) Validate() error {
	if c.FilterName == "" {
		return errors.NewLocalizedError("The filter name is mandatory.")
	}
	if len(c.Filters) == 0 {
		return errors.NewLocalizedError("The filters set should not be empty")
	}
	return nil
}

//// Merge update the given category fields.
//func (c FilterForm) Merge(filter *model.Filter) *model.Filter {
//	filter.FilterName = c.FilterName
//	return filter
//}

// NewFilterForm returns a new FilterForm.
func NewFilterForm(r *http.Request) *FilterForm {
	_ = r.ParseForm()
	logger.Debug("filter_name: %s, filters: %s", r.Form["filter_name"][0], strings.Join(r.Form["filters[0]"], ", "))
	return &FilterForm{
		FilterName: r.Form["filter_name"][0],
		Filters:    r.Form["filters[]"],
	}
}
