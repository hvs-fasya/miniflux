// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package daemon

import (
	"net/http"

	"github.com/miniflux/miniflux/api"
	"github.com/miniflux/miniflux/config"
	"github.com/miniflux/miniflux/fever"
	"github.com/miniflux/miniflux/locale"
	"github.com/miniflux/miniflux/middleware"
	"github.com/miniflux/miniflux/news"
	"github.com/miniflux/miniflux/reader/feed"
	"github.com/miniflux/miniflux/scheduler"
	"github.com/miniflux/miniflux/storage"
	"github.com/miniflux/miniflux/template"
	"github.com/miniflux/miniflux/ui"

	"github.com/gorilla/mux"
)

func routes(cfg *config.Config, store *storage.Storage, feedHandler *feed.Handler, pool *scheduler.WorkerPool, translator *locale.Translator) *mux.Router {
	router := mux.NewRouter()
	templateEngine := template.NewEngine(cfg, router, translator)
	templateNewsEngine := template.NewNewsEngine(cfg, router, translator)
	apiController := api.NewController(store, feedHandler)
	feverController := fever.NewController(store)
	uiController := ui.NewController(cfg, store, pool, feedHandler, templateEngine, translator, router)
	newsController := news.NewController(cfg, store, pool, feedHandler, templateNewsEngine, translator, router)
	middleware := middleware.New(cfg, store, router)

	if cfg.BasePath() != "" {
		router = router.PathPrefix(cfg.BasePath()).Subrouter()
	}

	router.Use(middleware.HeaderConfig)
	router.Use(middleware.Logging)
	router.Use(middleware.CommonHeaders)

	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("User-agent: *\nDisallow: /"))
	})

	newsRouter := router.PathPrefix(("/news")).Subrouter()
	newsRouter.HandleFunc("/{name}.css", newsController.Stylesheet).Name("news_stylesheet").Methods("GET")
	newsRouter.HandleFunc("/js", newsController.Javascript).Name("news_javascript").Methods("GET")
	newsRouter.HandleFunc("/mdlselect", newsController.MdlSelect).Name("news_mdlselect").Methods("GET")
	newsRouter.HandleFunc("/dialog-polyfill", newsController.DialogPolyfill).Name("news_dialog-polyfill").Methods("GET")
	newsRouter.HandleFunc("/news", newsController.News).Name("news_news").Methods("GET")
	newsRouter.HandleFunc("/favicon.ico", newsController.Favicon).Name("favicon").Methods("GET")
	//newsRouter.HandleFunc("/icon/{filename}", newsController.AppIcon).Name("appIcon").Methods("GET")
	newsRouter.HandleFunc("/manifest.json", newsController.WebManifest).Name("webManifest").Methods("GET")
	newsRouter.HandleFunc("/", newsController.Home).Name("home").Methods("GET")
	newsRouter.HandleFunc("/media", newsController.Media).Name("news_media").Methods("GET")
	newsRouter.HandleFunc("/travel", newsController.Travel).Name("news_travel").Methods("GET")
	newsRouter.HandleFunc("/official", newsController.Official).Name("news_official").Methods("GET")
	newsRouter.HandleFunc("/security", newsController.Security).Name("news_security").Methods("GET")
	newsRouter.HandleFunc("/visa", newsController.Visa).Name("news_visa").Methods("GET")
	newsRouter.HandleFunc("/sources", newsController.Sources).Name("sources").Methods("GET")
	newsRouter.HandleFunc("/icon/{iconID}", newsController.ShowIcon).Name("feedicon").Methods("GET")
	newsRouter.HandleFunc("/entry/download/{entryID}", newsController.FetchContent).Name("newsFetchContent").Methods("POST")

	feverRouter := router.PathPrefix("/fever").Subrouter()
	feverRouter.Use(middleware.FeverAuth)
	feverRouter.HandleFunc("/", feverController.Handler).Name("feverEndpoint")

	apiRouter := router.PathPrefix("/v1").Subrouter()
	apiRouter.Use(middleware.BasicAuth)
	apiRouter.Methods("OPTIONS").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(""))
			return
		})
	apiRouter.HandleFunc("/users", apiController.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/users", apiController.Users).Methods("GET")
	apiRouter.HandleFunc("/users/{userID:[0-9]+}", apiController.UserByID).Methods("GET")
	apiRouter.HandleFunc("/users/{userID:[0-9]+}", apiController.UpdateUser).Methods("PUT")
	apiRouter.HandleFunc("/users/{userID:[0-9]+}", apiController.RemoveUser).Methods("DELETE")
	apiRouter.HandleFunc("/users/{username}", apiController.UserByUsername).Methods("GET")
	apiRouter.HandleFunc("/categories", apiController.CreateCategory).Methods("POST")
	apiRouter.HandleFunc("/categories", apiController.GetCategories).Methods("GET")
	apiRouter.HandleFunc("/categories/{categoryID}", apiController.UpdateCategory).Methods("PUT")
	apiRouter.HandleFunc("/categories/{categoryID}", apiController.RemoveCategory).Methods("DELETE")
	apiRouter.HandleFunc("/discover", apiController.GetSubscriptions).Methods("POST")
	apiRouter.HandleFunc("/feeds", apiController.CreateFeed).Methods("POST")
	apiRouter.HandleFunc("/feeds", apiController.GetFeeds).Methods("Get")
	apiRouter.HandleFunc("/feeds/{feedID}/refresh", apiController.RefreshFeed).Methods("PUT")
	apiRouter.HandleFunc("/feeds/{feedID}", apiController.GetFeed).Methods("GET")
	apiRouter.HandleFunc("/feeds/{feedID}", apiController.UpdateFeed).Methods("PUT")
	apiRouter.HandleFunc("/feeds/{feedID}", apiController.RemoveFeed).Methods("DELETE")
	apiRouter.HandleFunc("/feeds/{feedID}/icon", apiController.FeedIcon).Methods("GET")
	apiRouter.HandleFunc("/export", apiController.Export).Methods("GET")
	apiRouter.HandleFunc("/import", apiController.Import).Methods("POST")
	apiRouter.HandleFunc("/feeds/{feedID}/entries", apiController.GetFeedEntries).Methods("GET")
	apiRouter.HandleFunc("/feeds/{feedID}/entries/{entryID}", apiController.GetFeedEntry).Methods("GET")
	apiRouter.HandleFunc("/entries", apiController.GetEntries).Methods("GET")
	apiRouter.HandleFunc("/entries", apiController.SetEntryStatus).Methods("PUT")
	apiRouter.HandleFunc("/entries/{entryID}", apiController.GetEntry).Methods("GET")
	apiRouter.HandleFunc("/entries/{entryID}/bookmark", apiController.ToggleBookmark).Methods("PUT")

	apiRouter.HandleFunc("/filters", apiController.CreateFilter).Methods("POST")
	apiRouter.HandleFunc("/filters", apiController.GetFilters).Methods("GET")
	apiRouter.HandleFunc("/filters/{filterID}", apiController.RemoveFilter).Methods("DELETE")
	apiRouter.HandleFunc("/headlines", apiController.CreateHeadline).Methods("POST")
	apiRouter.HandleFunc("/headlines", apiController.HeadlinesFull).Methods("GET")
	apiRouter.HandleFunc("/countries", apiController.GetCountries).Methods("GET")
	apiRouter.HandleFunc("/headlines/{headlineID}/icon", apiController.HeadlineIcon).Methods("GET")
	apiRouter.HandleFunc("/alerts", apiController.AlertsFull).Methods("GET")

	uiRouter := router.NewRoute().Subrouter()
	uiRouter.Use(middleware.AppSession)
	uiRouter.Use(middleware.UserSession)
	uiRouter.HandleFunc("/stylesheets/{name}.css", uiController.Stylesheet).Name("stylesheet").Methods("GET")
	uiRouter.HandleFunc("/js", uiController.Javascript).Name("javascript").Methods("GET")
	uiRouter.HandleFunc("/favicon.ico", uiController.Favicon).Name("favicon").Methods("GET")
	uiRouter.HandleFunc("/icon/{filename}", uiController.AppIcon).Name("appIcon").Methods("GET")
	uiRouter.HandleFunc("/manifest.json", uiController.WebManifest).Name("webManifest").Methods("GET")
	uiRouter.HandleFunc("/subscribe", uiController.AddSubscription).Name("addSubscription").Methods("GET")
	uiRouter.HandleFunc("/subscribe", uiController.SubmitSubscription).Name("submitSubscription").Methods("POST")
	uiRouter.HandleFunc("/subscriptions", uiController.ChooseSubscription).Name("chooseSubscription").Methods("POST")
	uiRouter.HandleFunc("/mark-all-as-read", uiController.MarkAllAsRead).Name("markAllAsRead").Methods("GET")
	uiRouter.HandleFunc("/unread", uiController.ShowUnreadPage).Name("unread").Methods("GET")
	uiRouter.HandleFunc("/history", uiController.ShowHistoryPage).Name("history").Methods("GET")
	uiRouter.HandleFunc("/starred", uiController.ShowStarredPage).Name("starred").Methods("GET")
	uiRouter.HandleFunc("/feed/{feedID}/refresh", uiController.RefreshFeed).Name("refreshFeed").Methods("GET")
	uiRouter.HandleFunc("/feed/{feedID}/edit", uiController.EditFeed).Name("editFeed").Methods("GET")
	uiRouter.HandleFunc("/feed/{feedID}/remove", uiController.RemoveFeed).Name("removeFeed").Methods("POST")
	uiRouter.HandleFunc("/feed/{feedID}/update", uiController.UpdateFeed).Name("updateFeed").Methods("POST")
	uiRouter.HandleFunc("/feed/{feedID}/entries", uiController.ShowFeedEntries).Name("feedEntries").Methods("GET")
	uiRouter.HandleFunc("/feeds", uiController.ShowFeedsPage).Name("feeds").Methods("GET")
	uiRouter.HandleFunc("/feeds/refresh", uiController.RefreshAllFeeds).Name("refreshAllFeeds").Methods("GET")
	uiRouter.HandleFunc("/unread/entry/{entryID}", uiController.ShowUnreadEntry).Name("unreadEntry").Methods("GET")
	uiRouter.HandleFunc("/history/entry/{entryID}", uiController.ShowReadEntry).Name("readEntry").Methods("GET")
	uiRouter.HandleFunc("/history/flush", uiController.FlushHistory).Name("flushHistory").Methods("GET")
	uiRouter.HandleFunc("/feed/{feedID}/entry/{entryID}", uiController.ShowFeedEntry).Name("feedEntry").Methods("GET")
	uiRouter.HandleFunc("/category/{categoryID}/entry/{entryID}", uiController.ShowCategoryEntry).Name("categoryEntry").Methods("GET")
	uiRouter.HandleFunc("/starred/entry/{entryID}", uiController.ShowStarredEntry).Name("starredEntry").Methods("GET")
	uiRouter.HandleFunc("/entry/status", uiController.UpdateEntriesStatus).Name("updateEntriesStatus").Methods("POST")
	uiRouter.HandleFunc("/entry/save/{entryID}", uiController.SaveEntry).Name("saveEntry").Methods("POST")
	uiRouter.HandleFunc("/entry/download/{entryID}", uiController.FetchContent).Name("fetchContent").Methods("POST")
	uiRouter.HandleFunc("/entry/bookmark/{entryID}", uiController.ToggleBookmark).Name("toggleBookmark").Methods("POST")
	uiRouter.HandleFunc("/categories", uiController.CategoryList).Name("categories").Methods("GET")
	uiRouter.HandleFunc("/category/create", uiController.CreateCategory).Name("createCategory").Methods("GET")
	uiRouter.HandleFunc("/category/save", uiController.SaveCategory).Name("saveCategory").Methods("POST")
	uiRouter.HandleFunc("/category/{categoryID}/entries", uiController.CategoryEntries).Name("categoryEntries").Methods("GET")
	uiRouter.HandleFunc("/category/{categoryID}/edit", uiController.EditCategory).Name("editCategory").Methods("GET")
	uiRouter.HandleFunc("/category/{categoryID}/update", uiController.UpdateCategory).Name("updateCategory").Methods("POST")
	uiRouter.HandleFunc("/category/{categoryID}/remove", uiController.RemoveCategory).Name("removeCategory").Methods("POST")
	uiRouter.HandleFunc("/feed/icon/{iconID}", uiController.ShowIcon).Name("icon").Methods("GET")
	uiRouter.HandleFunc("/proxy/{encodedURL}", uiController.ImageProxy).Name("proxy").Methods("GET")
	uiRouter.HandleFunc("/users", uiController.ShowUsers).Name("users").Methods("GET")
	uiRouter.HandleFunc("/user/create", uiController.CreateUser).Name("createUser").Methods("GET")
	uiRouter.HandleFunc("/user/save", uiController.SaveUser).Name("saveUser").Methods("POST")
	uiRouter.HandleFunc("/users/{userID}/edit", uiController.EditUser).Name("editUser").Methods("GET")
	uiRouter.HandleFunc("/users/{userID}/update", uiController.UpdateUser).Name("updateUser").Methods("POST")
	uiRouter.HandleFunc("/users/{userID}/remove", uiController.RemoveUser).Name("removeUser").Methods("POST")
	uiRouter.HandleFunc("/about", uiController.About).Name("about").Methods("GET")
	uiRouter.HandleFunc("/settings", uiController.ShowSettings).Name("settings").Methods("GET")
	uiRouter.HandleFunc("/settings", uiController.UpdateSettings).Name("updateSettings").Methods("POST")
	uiRouter.HandleFunc("/bookmarklet", uiController.Bookmarklet).Name("bookmarklet").Methods("GET")
	uiRouter.HandleFunc("/integrations", uiController.ShowIntegrations).Name("integrations").Methods("GET")
	uiRouter.HandleFunc("/integration", uiController.UpdateIntegration).Name("updateIntegration").Methods("POST")
	uiRouter.HandleFunc("/sessions", uiController.ShowSessions).Name("sessions").Methods("GET")
	uiRouter.HandleFunc("/sessions/{sessionID}/remove", uiController.RemoveSession).Name("removeSession").Methods("POST")
	uiRouter.HandleFunc("/export", uiController.Export).Name("export").Methods("GET")
	uiRouter.HandleFunc("/import", uiController.Import).Name("import").Methods("GET")
	uiRouter.HandleFunc("/upload", uiController.UploadOPML).Name("uploadOPML").Methods("POST")
	uiRouter.HandleFunc("/oauth2/{provider}/unlink", uiController.OAuth2Unlink).Name("oauth2Unlink").Methods("GET")
	uiRouter.HandleFunc("/oauth2/{provider}/redirect", uiController.OAuth2Redirect).Name("oauth2Redirect").Methods("GET")
	uiRouter.HandleFunc("/oauth2/{provider}/callback", uiController.OAuth2Callback).Name("oauth2Callback").Methods("GET")
	uiRouter.HandleFunc("/login", uiController.CheckLogin).Name("checkLogin").Methods("POST")
	uiRouter.HandleFunc("/logout", uiController.Logout).Name("logout").Methods("GET")
	uiRouter.HandleFunc("/", uiController.ShowLoginPage).Name("login").Methods("GET")

	uiRouter.HandleFunc("/filters", uiController.ShowFeedFilters).Name("filters").Methods("GET")
	uiRouter.HandleFunc("/filters/{id}", uiController.ShowFeedFilters).Name("filtersId").Methods("GET")
	uiRouter.HandleFunc("/filters/{id}", uiController.ShowFeedFilters).Queries("months", "{months}").Name("filtersIdMonths").Methods("GET")
	uiRouter.HandleFunc("/filters", uiController.CreateFilter).Name("createFilter").Methods("POST")
	uiRouter.HandleFunc("/filters", uiController.RemoveFilter).Name("removeFilter").Methods("DELETE")
	uiRouter.HandleFunc("/filters/{filterID}/remove", uiController.RemoveFilter).Name("removeFilter").Methods("POST")

	return router
}
