package main

import (
	"net/http"
	"regexp"
	"text/template"
	"github.com/gernest/alien"
	"path"
	"mime"
	"io"
)

var (
	_DEFAULT_METHODS = []string{http.MethodGet}
	_RE = regexp.MustCompile("[\\s,]+")

	_FUNCS = template.FuncMap{
		"pathVar": func(r *http.Request, n string)string {
			p := alien.GetParams(r)
			return p.Get(n)
		},
	}
)

type Endpoint Action

func (a *Endpoint) GetMethods() []string {
	if a.Method == "" {
		return _DEFAULT_METHODS
	}
	return _RE.Split(a.Method, -1)
}

func (a *Endpoint) GetUri() string {
	return a.Uri
}

func (a *Endpoint) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// redirect
		if a.RedirectBody != "" {
			http.Redirect(w, r, a.RedirectBody, http.StatusMovedPermanently)
			return
		}

		h := w.Header()
		if ServiceConf.DefaultContentType != "" {
			h.Set("Content-Type", ServiceConf.DefaultContentType)
		}
		a.setHeaders(h)
		a.setCookies(h)

		// file body
		if a.FileBody != "" {
			a.showFile(w, r)
			return
		}

		if a.Status != 0 {
			w.WriteHeader(a.Status)
		}

		// template
		if a.TmplBody != "" {
			a.execTemplate(w, r)
			return
		}

		// average body
		if a.Body != "" {
			io.WriteString(w, a.Body)
		}
	}
}

func CreateEndpoints() []*Endpoint {
	actions := ServiceConf.Actions
	es := make([]*Endpoint, len(ServiceConf.Actions))
	for i := range actions {
		es[i] = (*Endpoint)(&actions[i])
	}
	return es
}

func (a *Endpoint) setHeaders(h http.Header) {
	if a.Headers == nil || len(a.Headers) == 0 {
		return
	}
	for k,v := range a.Headers {
		h.Set(k, v)
	}
}

func (a *Endpoint) setCookies(h http.Header) {
	if a.Cookies == nil || len(a.Cookies) == 0 {
		return
	}
	for n,v := range a.Cookies {
		c := http.Cookie{Name:n, Value:v}
		h.Add("Set-Cookie", c.String())
	}
}

func (a *Endpoint) showFile(w http.ResponseWriter, r *http.Request) {
	ext := path.Ext(a.FileBody)
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
	}
	fn := path.Join(ServiceConf.Root, a.FileBody)
	if a.Status != 0 {
		w.WriteHeader(a.Status)
	}
	http.ServeFile(w, r, fn)
}

func (a *Endpoint) execTemplate(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("webpage").Funcs(_FUNCS).Parse(a.TmplBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, map[string]interface{}{"r": r})
}

