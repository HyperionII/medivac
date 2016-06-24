package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/hyperionii/medivac/config"
	"github.com/hyperionii/medivac/handlers/httputils"
	"github.com/hyperionii/medivac/routes"
	"golang.org/x/net/context"
)

// MiddlewareFunc describes a function that takes a ContextHandler and
// returns a ContextHandler.
//
// The idea of a middleware function is to validate/read/modify data before or
// after calling the next middleware function.
type MiddlewareFunc func(httputils.APIHandler) httputils.APIHandler

// extendSessionLifetime determines if the session's lifetime needs to be
// extended. Session's lifetime should be extended only if the session's
// current lifetime is below sessionLifeTime/2. Returns true if the session
// needs to be extended.
func extendSessionLifetime(sessionData *httputils.SessionData, sessionLifeTime time.Duration) bool {
	return sessionData.ExpiresAt.Sub(time.Now()) <= sessionLifeTime/2
}

// ValidateAuth validates that the user cookie is set up before calling the
// handler passed as parameter.
func ValidateAuth(h httputils.APIHandler) httputils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			ctx             = r.Context()
			cookieStore, ok = ctx.Value("cookieStore").(*sessions.CookieStore)
			cfg             = ctx.Value("config").(*config.Config)
		)

		if !ok {
			httputils.WriteError(w, http.StatusInternalServerError, "")
			return fmt.Errorf("validate auth: could not cast value as cookie store: %s", ctx.Value("cookieStore"))
		}

		session, err := cookieStore.Get(r, cfg.SessionCookie.Name)

		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return nil
		}

		sessionData, ok := session.Values["data"].(*httputils.SessionData)

		if !ok || sessionData.IsInvalid() {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return nil
		} else if time.Now().After(sessionData.ExpiresAt) {
			session.Options.MaxAge = -1
			session.Save(r, w)
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

			return nil
		}

		// Extend the session's lifetime.
		cfg, ok = ctx.Value("config").(*config.Config)

		if !ok {
			httputils.WriteError(w, http.StatusInternalServerError, "")
			return fmt.Errorf("validate auth: error casting config object: %s", ctx.Value("config"))
		}

		// Save session only if the session was extended.
		if extendSessionLifetime(sessionData, cfg.SessionCookie.LifeTime) {
			sessionData.ExpiresAt = time.Now().Add(cfg.SessionCookie.LifeTime)
			session.Save(r, w)
		}

		ctx = context.WithValue(ctx, "sessionData", sessionData)
		authenticatedRequest := r.WithContext(ctx)

		return h(w, authenticatedRequest)
	}
}

// GzipContent is a middleware function for handlers to encode content to gzip.
func GzipContent(h httputils.APIHandler) httputils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Add("Vary", "Accept-Encoding")

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			return h(w, r)
		}

		w.Header().Set("Content-Encoding", "gzip")

		gzipResponseWriter := httputils.NewGzipResponseWriter(w)
		defer gzipResponseWriter.Close()

		return h(gzipResponseWriter, r)
	}
}

// Authorize validates privileges for the current user. Each route must have
// an array of privileges that point which users can make a call to it.
//
// Note:
//
// It is assumed that ValidateAuth was called before this function, or at
// least some other session check was done before this.
func Authorize(h httputils.APIHandler) httputils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		var (
			ctx   = r.Context()
			_, ok = ctx.Value("sessionData").(*httputils.SessionData)
		)

		if !ok {
			httputils.WriteError(w, http.StatusInternalServerError, "")
			return fmt.Errorf("authorize: could not cast value as session data: %s", ctx.Value("sessionData"))
		}

		route, ok := ctx.Value("route").(routes.Route)
		if !ok {
			httputils.WriteError(w, http.StatusInternalServerError, "")
			return fmt.Errorf("authorize: could not cast value as route: %s", ctx.Value("route"))
		}

		requiredRoles := route.RequiredRoles()
		if len(requiredRoles) == 0 {
			return h(w, r)
		}

		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return nil
	}
}

// HandleHTTPError sets the appropriate headers to the response if a http
// handler returned an error. This might be used in the future if different
// types of errors are returned.
func HandleHTTPError(h httputils.APIHandler) httputils.APIHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := h(w, r)

		if err != nil {
			httputils.WriteError(w, http.StatusInternalServerError, "")
		}

		return err
	}
}
