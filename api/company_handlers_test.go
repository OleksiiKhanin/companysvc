package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/OleksiiKhanin/companysvc/domain"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockCompany struct {
	t       *testing.T
	create  func(_ context.Context, _ *domain.Company) error
	update  func(_ context.Context, _ string, _ string, _ *domain.Company) error
	delete  func(_ context.Context, _ string, _ string) error
	get     func(_ context.Context, _ string, _ string) (domain.Company, error)
	getMany func(_ context.Context, _ *domain.FilterOptions) ([]domain.Company, error)
}

func (m *MockCompany) Get(ctx context.Context, name, code string) (domain.Company, error) {
	if m.get != nil {
		return m.get(ctx, name, code)
	}
	m.t.Error("should not be called")
	return domain.Company{}, errors.New("not implemented method")
}

func (m *MockCompany) GetMany(ctx context.Context, filter *domain.FilterOptions) ([]domain.Company, error) {
	if m.getMany != nil {
		return m.getMany(ctx, filter)
	}
	m.t.Error("should not be called")
	return nil, errors.New("not implemented method")
}

func (m *MockCompany) Create(ctx context.Context, company *domain.Company) error {
	if m.create != nil {
		return m.create(ctx, company)
	}
	m.t.Error("should not be called")
	return errors.New("not implemented method")
}

func (m *MockCompany) Update(ctx context.Context, oldName, oldCode string, company *domain.Company) error {
	if m.update != nil {
		return m.update(ctx, oldName, oldCode, company)
	}
	m.t.Error("should not be called")
	return errors.New("not implemented method")
}

func (m *MockCompany) Delete(ctx context.Context, name, code string) error {
	if m.delete != nil {
		return m.delete(ctx, name, code)
	}
	m.t.Error("should not be called")
	return errors.New("not implemented method")
}

func execRequest(req *http.Request, company domain.ICompany) *httptest.ResponseRecorder {
	r := mux.NewRouter()
	InitAPI(r, company, log.StandardLogger())
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func TestGetCompanyHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/company/test/testCode", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := execRequest(
		req,
		&MockCompany{
			get: func(_ context.Context, name, code string) (domain.Company, error) {
				return domain.Company{Name: name, Code: code}, nil
			},
		},
	)
	if rr.Code != http.StatusOK {
		t.Errorf("incorrect status code when try get company: want %d but got %d", http.StatusOK, rr.Code)
	}
	var company domain.Company
	err = json.NewDecoder(rr.Body).Decode(&company)
	if err != nil {
		t.Errorf("Parse error %s", err.Error())
	}
	if company.Name != "test" || company.Code != "testCode" {
		t.Errorf("wait test:testCode but got %s:%s", company.Name, company.Code)
	}
}

func TestGetManyCompaniesHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/companies?limit=3&name=test&code=testCode", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := execRequest(
		req,
		&MockCompany{
			getMany: func(_ context.Context, options *domain.FilterOptions) ([]domain.Company, error) {
				if options.Params["name"] != "test" {
					t.Errorf("want test but get %s", options.Params["name"])
				}
				if options.Params["code"] != "testCode" {
					t.Errorf("want testCode but get %s", options.Params["code"])
				}
				if options.Limit == nil {
					t.Error("limit should be equal 3 but get nil")
				} else if *options.Limit != 3 {
					t.Errorf("want 3 and but get %d", *options.Limit)
				}
				return []domain.Company{}, nil
			},
		},
	)
	if rr.Code != http.StatusOK {
		t.Errorf("incorrect status code when try get companies: want %d but got %d", http.StatusOK, rr.Code)
	}
}

func TestCreateCompanyHandler(t *testing.T) {
	req, err := http.NewRequest(
		"POST",
		"/v1/company",
		strings.NewReader("{\"name\":\"test\",\"code\":\"testCode\"}"),
	)
	if err != nil {
		t.Fatal(err)
	}
	rr := execRequest(
		req,
		&MockCompany{
			create: func(_ context.Context, company *domain.Company) error {
				return nil
			},
		},
	)
	if rr.Code != http.StatusCreated {
		t.Errorf("incorrect status code when try create company: want %d but got %d", http.StatusCreated, rr.Code)
	}
	var company domain.Company
	err = json.NewDecoder(rr.Body).Decode(&company)
	if err != nil {
		t.Errorf("Parse error %s", err.Error())
	}
	if company.Name != "test" || company.Code != "testCode" {
		t.Errorf("wait test:testCode but got %s:%s", company.Name, company.Code)
	}
}

func TestUpdateCompanyHandler(t *testing.T) {
	req, err := http.NewRequest(
		"PUT",
		"/v1/company/oldName/oldCode",
		strings.NewReader("{\"name\":\"test\",\"code\":\"testCode\"}"),
	)
	if err != nil {
		t.Fatal(err)
	}
	rr := execRequest(
		req,
		&MockCompany{
			update: func(_ context.Context, oldName, oldCode string, company *domain.Company) error {
				if !strings.EqualFold(oldName, "oldName") {
					t.Errorf("want oldName but got %s when test update", oldName)
				}
				if !strings.EqualFold(oldCode, "oldCode") {
					t.Errorf("want oldCode but got %s when test update", oldCode)
				}
				return nil
			},
		},
	)
	if rr.Code != http.StatusAccepted {
		t.Errorf("incorrect status code when try update company: want %d but got %d", http.StatusAccepted, rr.Code)
	}
	var company domain.Company
	err = json.NewDecoder(rr.Body).Decode(&company)
	if err != nil {
		t.Errorf("Parse error %s", err.Error())
	}
	if company.Name != "test" || company.Code != "testCode" {
		t.Errorf("wait test:testCode but got %s:%s", company.Name, company.Code)
	}
}

func TestDeleteCompanyHandler(t *testing.T) {
	req, err := http.NewRequest(
		"DELETE",
		"/v1/company/Name/Code",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	rr := execRequest(
		req,
		&MockCompany{
			delete: func(_ context.Context, name, code string) error {
				if !strings.EqualFold(name, "Name") {
					t.Errorf("want Name but got %s when test update", name)
				}
				if !strings.EqualFold(code, "Code") {
					t.Errorf("want Code but got %s when test update", code)
				}
				return nil
			},
		},
	)
	if rr.Code != http.StatusAccepted {
		t.Errorf("incorrect status code when try delete company: want %d but got %d", http.StatusAccepted, rr.Code)
	}
}
