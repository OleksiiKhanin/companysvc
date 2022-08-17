package api

import (
	"context"
	"github.com/OleksiiKhanin/companysvc/domain"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type API struct {
	iCompany  domain.ICompany
	l         *log.Logger
	logPrefix string
}

func middlewareSetUserIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		addr := strings.Split(r.RemoteAddr, ":")
		ctx := context.WithValue(r.Context(), domain.CtxUserIPKey, addr[0])
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func getRecoveryMiddleware(l *log.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					l.Errorf("PANIC recovered: %v", err)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// InitAPI - init all CRUD operation
func InitAPI(r *mux.Router, company domain.ICompany, l *log.Logger) {
	api := API{iCompany: company, l: l, logPrefix: "API"}

	r.Use(getRecoveryMiddleware(l), middlewareSetUserIP)
	r.Methods(http.MethodGet).Path("/v1/companies").HandlerFunc(api.getCompaniesHandler)
	r.Methods(http.MethodGet).Path("/v1/company/{name}/{code}").HandlerFunc(api.getCompanyHandler)
	r.Methods(http.MethodDelete).Path("/v1/company/{name}/{code}").HandlerFunc(api.deleteCompanyHandler)
	r.Methods(http.MethodPut).Path("/v1/company/{name}/{code}").HandlerFunc(api.updateCompanyHandler)
	r.Methods(http.MethodPost).Path("/v1/company").HandlerFunc(api.createCompanyHandler)
}
