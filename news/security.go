package news

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/reader/scraper"
	"github.com/miniflux/miniflux/ui/session"
	"github.com/miniflux/miniflux/ui/view"
)

var (
	tabs = []string{
		"risk",
		"security",
		"laws",
		"disasters",
	}
)

// Security shows the Security template
func (c *Controller) Security(w http.ResponseWriter, r *http.Request) {
	ctx := context.New(r)

	sess := session.New(c.store, ctx)
	v := view.New(c.tpl, ctx, sess)
	country := request.QueryParam(r, "country", DefaultCountry)
	country = "afghanistan"

	doc, err := scraper.Fetch("https://travel.gc.ca/destinations/"+country, ".tabpanels")
	if err != nil {
		logger.Debug("error: %s", err)
		html.ServerError(w, err)
		return
	}
	cards, err := proceedDoc(doc)
	if err != nil {
		logger.Debug("error: %s", err)
		html.ServerError(w, err)
		return
	}

	//security tab
	for _, tab := range tabs {
		v.Set(tab, cards[tab])
	}

	html.OK(w, v.NewsAjaxRender("news_security"))
}

func proceedDoc(doc string) (map[string]string, error) {
	var err error
	proceeded := make(map[string]string)
	reader := strings.NewReader(doc)

	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return proceeded, err
	}
	for _, tab := range tabs {
		proceeded[tab], err = getDetail(document, tab)
		if err != nil {
			return proceeded, err
		}
	}

	return proceeded, nil
}

func removeNodes(s *goquery.Selection) {
	s.Each(func(i int, s *goquery.Selection) {
		parent := s.Parent()
		if parent.Length() > 0 {
			parent.Get(0).RemoveChild(s.Get(0))
		}
	})
}

func getDetail(doc *goquery.Document, name string) (string, error) {
	selection := doc.Find("details[id=" + name + "]")
	selection.Each(func(i int, s *goquery.Selection) {
		removeNodes(s)
	})
	stringified, err := selection.Html()
	if err != nil {
		return "", err
	}
	return stringified, nil
}
