package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response/json"
	"github.com/miniflux/miniflux/model"
)

type SecurityOutput struct {
	LastUpdated   time.Time       `json:"last_updated"`
	StillValid    time.Time       `json:"still_valid"`
	LatestUpdates string          `json:"latest_updates"`
	RiskLevel     string          `json:"risk_level"`
	RiskDetails   string          `json:"risk_details"`
	Health        []*model.Health `json:"health"`
}

// AlertsFull is the API handler to get the list of scrapped security alerts +.....
func (c *Controller) AlertsFull(w http.ResponseWriter, r *http.Request) {
	countryName := request.QueryParam(r, "country", "")
	//filter := request.QueryParam(r, "filter", "all")
	var country *model.Country
	var err error
	if countryName != "" {
		country, err = c.store.CountryByName(countryName)
		if err != nil {
			json.ServerError(w, errors.New("Unable to fetch country by name"))
			return
		}
		if country == nil {
			json.BadRequest(w, errors.New("Wrong country name"))
			return
		}
	}

	//securityAlerts
	securityBuilder := c.store.NewSecurityQueryBuilder()
	if country != nil {
		securityBuilder.WithCountryID(country.ID)
	}

	securityAlerts, err := securityBuilder.GetSecurityAlerts()
	if err != nil {
		json.ServerError(w, errors.New("Unable to fetch securityAlerts"))
		return
	}

	var secAlerts = map[string]SecurityOutput{}
	for _, a := range securityAlerts {
		s := SecurityOutput{
			LastUpdated:   a.LastUpdated,
			StillValid:    a.StillValid,
			LatestUpdates: a.LatestUpdates,
			RiskLevel:     a.RiskLevel,
			RiskDetails:   a.RiskDetails,
			Health:        a.Health,
		}
		secAlerts[a.Country.ID] = s
	}

	publishObj := struct {
		Security map[string]SecurityOutput `json:"security"`
	}{
		secAlerts,
	}

	json.OK(w, publishObj)
}
