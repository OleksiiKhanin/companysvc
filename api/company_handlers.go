package api

import (
	"encoding/json"
	"errors"
	"github.com/OleksiiKhanin/companysvc/domain"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type parameters struct {
	name string
	code string
}

func (a *API) parseParameters(r *http.Request) (*parameters, error) {
	vars := mux.Vars(r)
	name := strings.TrimSpace(vars["name"])
	if name == "" {
		return nil, errors.New("name parameter can not be empty")
	}

	code := strings.TrimSpace(vars["code"])
	if code == "" {
		return nil, errors.New("code parameter can not be empty")
	}

	return &parameters{name: name, code: code}, nil
}

//getCompanyHandler get exact one company with name and code parameters
func (a *API) getCompanyHandler(w http.ResponseWriter, r *http.Request) {
	p, err := a.parseParameters(r)
	if err != nil {
		a.l.Infof("%s:Parse query parameters: %s", a.logPrefix, err.Error())
		a.handleError(w, httpError{code: http.StatusBadRequest, message: err.Error()})
		return
	}
	company, err := a.iCompany.Get(r.Context(), p.name, p.code)
	if err != nil {
		a.l.Warnf("%s:Get company: %s", a.logPrefix, err.Error())
		a.handleError(w, httpError{code: http.StatusInternalServerError, message: "Can not get company"})
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(company)
}

// getCompaniesHandler - return a list of companies by filter in query parameters like name, code, website etc
func (a *API) getCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	var companyParams = []string{"name", "code", "country", "website", "phone"}
	query := r.URL.Query()
	filter := domain.FilterOptions{
		Params: make(map[string]string),
	}
	if limit, err := strconv.Atoi(query.Get("limit")); err == nil {
		filter.Limit = &limit
	}
	for _, param := range companyParams {
		if query.Has(param) {
			filter.Params[param] = strings.TrimSpace(query.Get(param))
		}
	}
	companies, err := a.iCompany.GetMany(r.Context(), &filter)
	if err != nil {
		a.l.Warnf("%s:Get many companies: %s", a.logPrefix, err.Error())
		a.handleError(w, httpError{code: http.StatusInternalServerError, message: "Can not get companies"})

		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(companies)
}

func (a *API) createCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var newCompany domain.Company
	if err := json.NewDecoder(r.Body).Decode(&newCompany); err != nil {
		a.l.Infof("%s:Parse company: %s", a.logPrefix, err.Error())
		a.handleError(w, httpError{code: http.StatusBadRequest, message: err.Error()})
		return
	}
	if err := a.iCompany.Create(r.Context(), &newCompany); err != nil {
		a.l.Warnf("%s:Create company: %s", a.logPrefix, err.Error())
		a.handleError(w, httpError{code: http.StatusInternalServerError, message: "Can not create company"})
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCompany)
}

func (a *API) updateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	p, err := a.parseParameters(r)
	if err != nil {
		a.l.Infof("%s:Parse query parameters: %s", a.logPrefix, err.Error())
		a.handleError(w, httpError{code: http.StatusBadRequest, message: err.Error()})
		return
	}
	var company domain.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		a.l.Infof("%s:Parse company: %s", a.logPrefix, err.Error())
		a.handleError(w, httpError{code: http.StatusBadRequest, message: err.Error()})
		return
	}
	if err := a.iCompany.Update(r.Context(), p.name, p.code, &company); err != nil {
		a.l.Warnf("%s:Update company: %s", a.logPrefix, err.Error())
		a.handleError(w, httpError{code: http.StatusInternalServerError, message: "Can not update company"})
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(company)
}

func (a *API) deleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	p, err := a.parseParameters(r)
	if err != nil {
		a.l.Infof("%s:Parse query parameters: %s", a.logPrefix, err.Error())
		a.handleError(w, httpError{code: http.StatusBadRequest, message: err.Error()})
		return
	}
	if err = a.iCompany.Delete(r.Context(), p.name, p.code); err != nil {
		a.l.Warnf("%s:Update company: %s", a.logPrefix, err.Error())
		a.handleError(w, httpError{code: http.StatusInternalServerError, message: "Can not delete company"})
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
