// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package template

import (
	"bytes"
	"html/template"
	"time"

	"github.com/miniflux/miniflux/config"
	"github.com/miniflux/miniflux/errors"
	"github.com/miniflux/miniflux/locale"
	"github.com/miniflux/miniflux/logger"

	"github.com/gorilla/mux"
)

// Engine handles the templating system.
type Engine struct {
	templates  map[string]*template.Template
	translator *locale.Translator
	funcMap    *funcMap
}

func (e *Engine) parseAll() {
	commonTemplates := ""
	for _, content := range templateCommonMap {
		commonTemplates += content
	}

	for name, content := range templateViewsMap {
		logger.Debug("[Template] Parsing: %s", name)
		e.templates[name] = template.Must(template.New("main").Funcs(e.funcMap.Map()).Parse(commonTemplates + content))
	}
}

func (e *Engine) parseNewsAll() {
	commonNewsTemplates := ""
	for _, content := range templateNewsCommonMap {
		commonNewsTemplates += content
	}

	for name, content := range templateNewsViewsMap {
		logger.Debug("[Template] Parsing: %s", name)
		e.templates[name] = template.Must(template.New("main").Funcs(e.funcMap.Map()).Parse(commonNewsTemplates + content))
	}
}

// Render process a template and write the ouput.
func (e *Engine) Render(name, language string, data interface{}) []byte {
	tpl, ok := e.templates[name]
	if !ok {
		logger.Fatal("[Template] The template %s does not exists", name)
	}

	lang := e.translator.GetLanguage(language)
	tpl.Funcs(template.FuncMap{
		"elapsed": func(timezone string, t time.Time) string {
			return elapsedTime(lang, timezone, t)
		},
		"t": func(key interface{}, args ...interface{}) string {
			switch key.(type) {
			case string:
				return lang.Get(key.(string), args...)
			case errors.LocalizedError:
				return key.(errors.LocalizedError).Localize(lang)
			case *errors.LocalizedError:
				return key.(*errors.LocalizedError).Localize(lang)
			case error:
				return key.(error).Error()
			default:
				return ""
			}
		},
		"plural": func(key string, n int, args ...interface{}) string {
			return lang.Plural(key, n, args...)
		},
	})

	var b bytes.Buffer
	err := tpl.ExecuteTemplate(&b, "base", data)
	if err != nil {
		logger.Fatal("[Template] Unable to render template: %v", err)
	}

	return b.Bytes()
}

// NewsRender process a template and write the ouput.
func (e *Engine) NewsRender(name, language string, data interface{}) []byte {
	tpl, ok := e.templates[name]
	if !ok {
		logger.Fatal("[Template] The template %s does not exists", name)
	}

	lang := e.translator.GetLanguage(language)
	tpl.Funcs(template.FuncMap{
		"elapsed": func(timezone string, t time.Time) string {
			return elapsedTime(lang, timezone, t)
		},
		"t": func(key interface{}, args ...interface{}) string {
			switch key.(type) {
			case string:
				return lang.Get(key.(string), args...)
			case errors.LocalizedError:
				return key.(errors.LocalizedError).Localize(lang)
			case *errors.LocalizedError:
				return key.(*errors.LocalizedError).Localize(lang)
			case error:
				return key.(error).Error()
			default:
				return ""
			}
		},
		"plural": func(key string, n int, args ...interface{}) string {
			return lang.Plural(key, n, args...)
		},
	})

	var b bytes.Buffer
	err := tpl.ExecuteTemplate(&b, "news_base", data)
	if err != nil {
		logger.Fatal("[Template] Unable to render template: %v", err)
	}

	return b.Bytes()
}

// NewEngine returns a new template engine.
func NewEngine(cfg *config.Config, router *mux.Router, translator *locale.Translator) *Engine {
	tpl := &Engine{
		templates:  make(map[string]*template.Template),
		translator: translator,
		funcMap:    newFuncMap(cfg, router),
	}

	tpl.parseAll()
	return tpl
}

// NewNewsEngine returns a new template engine.
func NewNewsEngine(cfg *config.Config, router *mux.Router, translator *locale.Translator) *Engine {
	tpl := &Engine{
		templates:  make(map[string]*template.Template),
		translator: translator,
		funcMap:    newFuncMap(cfg, router),
	}

	tpl.parseNewsAll()
	return tpl
}
