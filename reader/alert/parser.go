package alert

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/storage"
	"github.com/miniflux/miniflux/timer"
)

var (
	errRequestFailed    = "Unable to execute request: %v"
	errServerFailure    = "Unable to fetch feed (statusCode=%d)"
	errDuplicate        = "This feed already exists (%s)"
	errNotFound         = "Feed %d not found"
	errEncoding         = "Unable to normalize encoding: %q"
	errCategoryNotFound = "Category not found for this user"
	errEmptyFeed        = "This feed is empty"
	ADVISORIES          = "https://travel.gc.ca/destinations/"
	urlsMap             = map[string]string{
		"Saint-Barthélemy":            "saint-barthelemy",
		"Côte d'Ivoire (Ivory Coast)": "cote-d-ivoire-ivory-coast",
		"Virgin Islands (U.S.)":       "virgin-islands-u-s",
		"São Tomé and Principe":       "sao-tome-and-principe",
		"Curaçao":                     "curacao",
		"Réunion":                     "reunion",
	}
)

// Handler contains all the logic to create and refresh alert.
type Handler struct {
	store *storage.Storage
}

func Task(store *storage.Storage) error {
	countries, err := store.Countries()
	handler := new(Handler)
	handler.store = store
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(countries))
	//for _, country := range countries {
	//	_, err = handler.handleAlert(country)
	//	if err != nil {
	//		logger.Debug("error: %s", err)
	//		return err
	//	}
	//}
	country, _ := store.CountryByName("Myanmar")
	_, err = handler.handleAlert(country)
	return nil
}

func (h *Handler) handleAlert(country *model.Country) (*model.Alert, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[Handler:handleAlert] countryID=%d", country.ID))

	alert := new(model.Alert)
	alert.CountryID = country.ID

	var countryName = country.Name
	scrapUrl, ok := urlsMap[countryName]
	if !ok {
		scrapUrl = strings.Replace(strings.ToLower(countryName), " ", "-", -1)
		scrapUrl = strings.Replace(scrapUrl, ",", "", -1)
		scrapUrl = strings.Replace(scrapUrl, "(", "", -1)
		scrapUrl = strings.Replace(scrapUrl, ")", "", -1)
		scrapUrl = strings.Replace(scrapUrl, "-&-", "-", -1)
	}
	if countryName == "Saint Vincent & the Grenadines" {
		fmt.Println(scrapUrl)
	}
	scrapUrl = ADVISORIES + scrapUrl

	doc, err := goquery.NewDocument(scrapUrl)
	if err != nil {
		logger.Debug("error: %s", err)
		return alert, nil
	}

	lastUpdated, e := time.Parse("January 02, 2006 15:04", doc.Find("span#Label9").Text())
	if e == nil {
		alert.LastUpdated = lastUpdated
	} else {
		alert.LastUpdated, err = time.Parse("January 2, 2006 15:04", doc.Find("span#Label9").Text())
		if err != nil {
			logger.Debug("error: %s", err)
			return alert, nil
		}
	}

	stillValid, e := time.Parse("January 02, 2006 15:04", doc.Find("span#Label12").Text())
	if e == nil {
		alert.StillValid = stillValid
	} else {
		alert.StillValid, err = time.Parse("January 2, 2006 15:04", doc.Find("span#Label12").Text())
		if err != nil {
			logger.Debug("error: %s", err)
			return alert, nil
		}
	}

	alert.LatestUpdates = doc.Find("span#Label11").Text()
	doc.Find("div.AdvisoryContainer").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			alert.RiskLevel = s.Find("h3").Text()
			alert.RiskDetails = s.Find("p").Text()
		}
	})

	existing, err := h.store.AlertByCountryID(country.ID)
	if err != nil {
		logger.Debug("error: %s", err)
		return nil, err
	}
	if existing == nil {
		err = h.store.CreateAlert(alert)
		if err != nil {
			logger.Debug("error: %s", err)
			return nil, err
		}
		logger.Debug("[Handler:CreateAlert] Alert created for countryID: %d", alert.CountryID)
		return alert, nil
	}

	err = h.store.UpdateAlert(alert)
	if err != nil {
		logger.Debug("error: %s", err)
		return nil, err
	}
	logger.Debug("[Handler:UpdateAlert] Alert updated for countryID: %d", alert.CountryID)

	healthTab := doc.Find("#health")
	panel := healthTab.Find(".panel-body")
	panel.Find("li").Each(func(i int, s *goquery.Selection) {
		//todo: all links not only first
		if i == 0 {
			healthLink := s.Find("a")
			if healthLink.Text() != "" {
				healthTitle := strings.Replace(healthLink.Text(), ": Advice for travellers", "", -1)
				contentLink, _ := healthLink.Attr("href")
				health, err := h.handleHealth(healthTitle, contentLink)
				healthDate, e := time.Parse("January 02, 2006", strings.Replace(s.Text(), healthLink.Text()+" - ", "", -1))
				if e != nil {
					healthDate, err = time.Parse("January 2, 2006", strings.Replace(s.Text(), healthLink.Text()+" - ", "", -1))
					if err != nil {
						logger.Debug("error: %s", err)
					}
				}
				err = h.handleHealthAlert(health, healthDate, countryName)
				if err != nil {
					logger.Debug("error: %s", err)
				}
			}
		}
	})
	if err != nil {
		logger.Debug("error: %s", err)
	}

	return alert, nil
}

func (h *Handler) handleHealth(title string, link string) (*model.Health, error) {
	logger.Debug("[Handler:handleHealth] for health_title: %s, health_link: %s", title, link)
	existing, err := h.store.HealthByTitle(title)
	if err != nil {
		logger.Debug("error: %s", err)
		return nil, err
	}
	if existing != nil && time.Now().Sub(existing.LastUpdated) < time.Duration(1*time.Hour) {
		return existing, nil
	}

	health := model.Health{
		HealthLink:  link,
		HealthTitle: title,
	}
	healthDoc, err := goquery.NewDocument(link)
	if err != nil {
		logger.Debug("error: %s", err)
	}
	health.HealthContent, err = healthDoc.Find("main").Html()
	if err != nil {
		logger.Debug("error: %s", err)
	}

	if existing == nil {

		err = h.store.CreateHealth(&health)
		if err != nil {
			logger.Debug("error: %s", err)
			return nil, err
		}
		logger.Debug("[Handler:CreateHealth] Health created for title: %s", health.HealthTitle)
		return &health, nil
	}

	err = h.store.UpdateHealth(&health)
	if err != nil {
		logger.Debug("error: %s", err)
		return nil, err
	}
	logger.Debug("[Handler:UpdateHealth] Health updated for title: %s", title)

	return &health, nil
}

func (h *Handler) handleHealthAlert(health *model.Health, healthDate time.Time, countryName string) error {
	logger.Debug("[Handler:handleHealthAlert] for health_title: %s, health_date: %s, country_name: %s", health.HealthTitle, healthDate, countryName)
	//ha := &model.HealthAlert{}
	return nil
}
