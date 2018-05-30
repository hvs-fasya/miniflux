package news

import (
"net/http"
"time"

"github.com/miniflux/miniflux/http/request"
"github.com/miniflux/miniflux/http/response"
"github.com/miniflux/miniflux/http/response/html"
)

// ShowIcon shows the feed icon.
func (c *Controller) ShowIcon(w http.ResponseWriter, r *http.Request) {
	iconID, err := request.IntParam(r, "iconID")
	if err != nil {
		html.BadRequest(w, err)
		return
	}

	icon, err := c.store.IconByID(iconID)
	if err != nil {
		html.ServerError(w, err)
		return
	}

	if icon == nil {
		html.NotFound(w)
		return
	}

	response.Cache(w, r, icon.MimeType, icon.Hash, icon.Content, 72*time.Hour)
}

