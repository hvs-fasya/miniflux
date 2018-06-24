package news

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/miniflux/miniflux/http/context"
	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/html"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/reader/scraper"
	"github.com/miniflux/miniflux/ui/session"
	"github.com/miniflux/miniflux/ui/view"
	"strings"
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
	doc, err = proceedDoc(doc)
	if err != nil {
		logger.Debug("error: %s", err)
		html.ServerError(w, err)
		return
	}

	//security tab
	v.Set("country", country)
	v.Set("doc", doc)

	html.OK(w, v.NewsAjaxRender("news_security"))
}

func proceedDoc(doc string) (string, error) {
	var err error
	reader := strings.NewReader(doc)

	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}

	document, err = removeExtraTabs(document)
	if err != nil {
		return "", err
	}

	proceeded, err := document.Html()
	if err != nil {
		return "", err
	}
	return proceeded, nil
}

func removeExtraTabs(document *goquery.Document) (*goquery.Document, error) {

	document.Find("details[id=entryexit],details[id=health],details[id=assistance]").Each(func(i int, s *goquery.Selection) {
		removeNodes(s)
	})
	logger.Debug("doc: %+v", document)

	return document, nil
}

func removeNodes(s *goquery.Selection) {
	s.Each(func(i int, s *goquery.Selection) {
		parent := s.Parent()
		if parent.Length() > 0 {
			parent.Get(0).RemoveChild(s.Get(0))
		}
	})
}
