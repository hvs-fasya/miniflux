package news


import (
	"net/http"
	"time"

	"github.com/miniflux/miniflux/http/request"
	"github.com/miniflux/miniflux/http/response"
	"github.com/miniflux/miniflux/ui/static"
	"github.com/miniflux/miniflux/http/route"
	"github.com/miniflux/miniflux/http/response/json"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/http/response/html"
	"encoding/base64"
)

// Stylesheet renders the CSS.
func (c *Controller) Stylesheet(w http.ResponseWriter, r *http.Request) {
	stylesheet := request.Param(r, "name", "white")
	body := static.Stylesheets["common"]
	etag := static.StylesheetsChecksums["common"]

	if theme, found := static.Stylesheets[stylesheet]; found {
		body += theme
		etag += static.StylesheetsChecksums[stylesheet]
	}

	response.Cache(w, r, "text/css; charset=utf-8", etag, []byte(body), 48*time.Hour)
}

// Javascript renders application client side code.
func (c *Controller) Javascript(w http.ResponseWriter, r *http.Request) {
	response.Cache(w, r, "text/javascript; charset=utf-8", static.JavascriptChecksums["app"], []byte(static.Javascript["app"]), 48*time.Hour)
}

// WebManifest renders web manifest file.
func (c *Controller) WebManifest(w http.ResponseWriter, r *http.Request) {
	type webManifestIcon struct {
		Source string `json:"src"`
		Sizes  string `json:"sizes"`
		Type   string `json:"type"`
	}

	type webManifest struct {
		Name        string            `json:"name"`
		Description string            `json:"description"`
		ShortName   string            `json:"short_name"`
		StartURL    string            `json:"start_url"`
		Icons       []webManifestIcon `json:"icons"`
		Display     string            `json:"display"`
	}

	manifest := &webManifest{
		Name:        "Miniflux",
		ShortName:   "Miniflux",
		Description: "Minimalist Feed Reader",
		Display:     "minimal-ui",
		StartURL:    route.Path(c.router, "unread"),
		Icons: []webManifestIcon{
			webManifestIcon{Source: route.Path(c.router, "appIcon", "filename", "touch-icon-ipad-retina.png"), Sizes: "144x144", Type: "image/png"},
			webManifestIcon{Source: route.Path(c.router, "appIcon", "filename", "touch-icon-iphone-retina.png"), Sizes: "114x114", Type: "image/png"},
		},
	}

	json.OK(w, manifest)
}

// Favicon renders the application favicon.
func (c *Controller) Favicon(w http.ResponseWriter, r *http.Request) {
	blob, err := base64.StdEncoding.DecodeString(static.Binaries["favicon.ico"])
	if err != nil {
		logger.Error("[Controller:Favicon] %v", err)
		html.NotFound(w)
		return
	}

	response.Cache(w, r, "image/x-icon", static.BinariesChecksums["favicon.ico"], blob, 48*time.Hour)
}

// AppIcon renders application icons.
func (c *Controller) AppIcon(w http.ResponseWriter, r *http.Request) {
	filename := request.Param(r, "filename", "favicon.png")
	encodedBlob, found := static.Binaries[filename]
	if !found {
		logger.Info("[Controller:AppIcon] This icon doesn't exists: %s", filename)
		html.NotFound(w)
		return
	}

	blob, err := base64.StdEncoding.DecodeString(encodedBlob)
	if err != nil {
		logger.Error("[Controller:AppIcon] %v", err)
		html.NotFound(w)
		return
	}

	response.Cache(w, r, "image/png", static.BinariesChecksums[filename], blob, 48*time.Hour)
}
