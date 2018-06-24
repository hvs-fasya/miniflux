package news

import (
	"net/http"
	"time"

	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/ui/session"
	"github.com/miniflux/miniflux/ui/view"
)

const (
	VisaNewsFilterTitle = "Visa changes"
)

// Visa shows the Visa template
func (c *Controller) Visa(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	sess := session.New(c.store, ctx)
	v := view.New(c.tpl, ctx, sess)
	offset := request.QueryIntParam(r, "offset", 0)
	limit := request.QueryIntParam(r, "limit", NewsEntriesLimit)

	visaFilter, err := c.store.FilterByNameWOUserID(VisaNewsFilterTitle)
	if err != nil || visaFilter == nil {
		logger.Debug("[ERROR] can not get filter by name %s", VisaNewsFilterTitle)
	}

	//visa tab
	visaBuilder := c.store.NewNewsEntryQueryBuilder()
	visaBuilder.WithLimit(limit)
	visaBuilder.WithOffset(offset)
	visaBuilder.WithDirection(DefaultSortingDirection)
	visaBuilder.WithOrder(model.DefaultSortingOrder)
	visaBuilder.WithoutStatus(model.EntryStatusRemoved)
	visaBuilder.WithFilter(visaFilter.Filters)

	visaStartDate := time.Now().AddDate(0, -1, 0)
	visaBuilder.After(&visaStartDate)

	visaEntries, err := visaBuilder.GetEntries()
	if err != nil {
		html.ServerError(w, err)
		return
	}
	visaCount, err := visaBuilder.CountEntries()
	if err != nil {
		html.ServerError(w, err)
		return
	}

	v.Set("visaentries", visaEntries)
	v.Set("visatotal", visaCount)
	v.Set("offset", offset)
	v.Set("limit", NewsEntriesLimit)

	var hasNext bool
	hasNext = (visaCount - offset) > limit
	v.Set("hasNext", hasNext)

	html.OK(w, v.NewsAjaxRender("news_visa"))
}
