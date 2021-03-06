package web

import (
	"context"
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	"github.com/maliceio/engine/api/types/web"
)

func (wr *webRouter) startWeb(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	config := &web.Config{
		Path: r.FormValue("path"),
	}
	result, err := sr.backend.Scan(r.FormValue("path"), config)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, result)
}

func (wr *webRouter) stopWeb(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	config := &web.Config{
		Path: r.FormValue("path"),
	}
	result, err := sr.backend.Scan(r.FormValue("path"), config)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, result)
}

func (wr *webRouter) backUpWeb(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	config := &web.Config{
		Path: r.FormValue("path"),
	}
	result, err := sr.backend.Scan(r.FormValue("path"), config)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, result)
}
