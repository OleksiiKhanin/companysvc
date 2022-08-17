package service

import (
	"context"
	"errors"
	"github.com/OleksiiKhanin/companysvc/domain"
	log "github.com/sirupsen/logrus"
	"strings"
	"testing"
)

type CountryResolverMock func(ip string) (string, error)

func (c CountryResolverMock) Resolve(ip string) (string, error) {
	return c(ip)
}

type PublisherMock func(subj string, data []byte) error

func (pub PublisherMock) Publish(subj string, data []byte) error {
	return pub(subj, data)
}

type MockICompanyDB struct {
	storage []*domain.Company
}

func (m *MockICompanyDB) find(name, code string) int {
	for i := range m.storage {
		if m.storage[i].Name == name && m.storage[i].Code == code {
			return i
		}
	}
	return -1
}

func (m *MockICompanyDB) Get(ctx context.Context, name, code string) (domain.Company, error) {
	i := m.find(name, code)
	if i < 0 {
		return domain.Company{}, errors.New("companies not find")
	}
	return *m.storage[i], nil
}

func (m *MockICompanyDB) GetMany(ctx context.Context, filter *domain.FilterOptions) ([]domain.Company, error) {
	res := make([]domain.Company, 0)
	for i := range m.storage {
		if filter.Limit != nil && *filter.Limit > 0 && len(res) >= *filter.Limit {
			return res, nil
		}
		if filter.Params["name"] != "" && strings.Contains(m.storage[i].Name, filter.Params["name"]) {
			res = append(res, *m.storage[i])
		}
		if filter.Params["code"] != "" && strings.Contains(m.storage[i].Code, filter.Params["code"]) {
			res = append(res, *m.storage[i])
		}
		if filter.Params["country"] != "" && strings.Contains(m.storage[i].Country, filter.Params["country"]) {
			res = append(res, *m.storage[i])
		}
		if filter.Params["website"] != "" && strings.Contains(m.storage[i].Website, filter.Params["website"]) {
			res = append(res, *m.storage[i])
		}
		if filter.Params["phone"] != "" && strings.Contains(m.storage[i].Phone, filter.Params["phone"]) {
			res = append(res, *m.storage[i])
		}
	}
	return res, nil
}

func (m *MockICompanyDB) Delete(_ context.Context, name, code string) error {
	i := m.find(name, code)
	if i < 0 {
		return errors.New("company does not exist")
	}
	if i >= 0 {
		m.storage = m.storage[:i]
		m.storage = append(m.storage, m.storage[i:]...)
	}
	return nil
}

func (m *MockICompanyDB) Create(_ context.Context, company *domain.Company) error {
	if m.find(company.Name, company.Code) >= 0 {
		return errors.New("company already exist")
	}
	m.storage = append(m.storage, company)
	return nil
}

func (m *MockICompanyDB) Update(_ context.Context, oldName, oldCode string, company *domain.Company) error {
	i := m.find(oldName, oldCode)
	if i < 0 {
		return errors.New("company not found")
	}
	m.storage[i] = company
	return nil
}

const (
	countrySuccess = "UA"
	countryFail    = "RU"
	pubChannel     = "test_channel"
)

type testCase struct {
	c       *domain.Company
	country string
	success bool
}

func TestCompanyServiceCreate(t *testing.T) {
	createCases := []testCase{
		{
			c:       &domain.Company{Name: "1", Code: "1"},
			country: countrySuccess,
			success: true,
		},
		{
			c:       &domain.Company{Name: "3", Code: "3"},
			country: countryFail,
			success: false,
		},
		{
			c:       &domain.Company{Name: "1", Code: "1"},
			country: countrySuccess,
			success: false,
		},
	}

	var company = companyService{
		ICompany:             &MockICompanyDB{},
		l:                    log.StandardLogger(),
		channel:              pubChannel,
		allowedCountriesCode: []string{countrySuccess},
	}

	ctx := context.WithValue(context.Background(), domain.CtxUserIPKey, "1.1.1.1")
	for i := range createCases {
		company.locationClient = CountryResolverMock(func(ip string) (string, error) {
			return createCases[i].country, nil
		})
		if !createCases[i].success {
			company.event = PublisherMock(func(_ string, _ []byte) error {
				t.Errorf("should not be called for case: %s:%s and country %s",
					createCases[i].c.Name,
					createCases[i].c.Code,
					createCases[i].country)
				return nil
			})
		} else {
			company.event = nil
		}
		err := company.Create(ctx, createCases[i].c)
		if createCases[i].success && err != nil {
			t.Error(err.Error())
		}
		if !createCases[i].success && err == nil {
			t.Errorf("case %s should be fail but not", createCases[i].country)
		}
	}
}

func TestCompanyServiceDelete(t *testing.T) {
	createCases := []testCase{
		{
			c:       &domain.Company{Name: "1", Code: "1"},
			country: countrySuccess,
			success: true,
		},
		{
			c:       &domain.Company{Name: "3", Code: "3"},
			country: countryFail,
			success: false,
		},
	}

	var company = companyService{
		ICompany:             &MockICompanyDB{},
		l:                    log.StandardLogger(),
		channel:              pubChannel,
		allowedCountriesCode: []string{countrySuccess},
	}

	ctx := context.WithValue(context.Background(), domain.CtxUserIPKey, "1.1.1.1")
	for i := range createCases {
		company.locationClient = CountryResolverMock(func(ip string) (string, error) {
			return createCases[i].country, nil
		})
		if !createCases[i].success {
			company.event = PublisherMock(func(_ string, _ []byte) error {
				t.Errorf("should not be called for case: %s:%s",
					createCases[i].c.Name,
					createCases[i].c.Code)
				return nil
			})
		} else {
			company.event = nil
		}
		err := company.Create(ctx, createCases[i].c)
		if createCases[i].success && err != nil {
			t.Error(err.Error())
		}
		if !createCases[i].success && err == nil {
			t.Errorf("case %s should be fail but not", createCases[i].country)
		}
	}
}

func TestCompanyServiceUpdate(t *testing.T) {
	updateCases := []testCase{
		{
			c:       &domain.Company{Name: "not found", Code: "not found"},
			success: false,
		},
		{
			c:       &domain.Company{Name: "1", Code: "1"},
			success: true,
		},
	}

	var company = companyService{
		ICompany: &MockICompanyDB{storage: []*domain.Company{
			&domain.Company{Name: "1", Code: "1"},
		}},
		l: log.StandardLogger(),
		locationClient: CountryResolverMock(func(ip string) (string, error) {
			t.Errorf("should not be called in update method")
			return "", nil
		}),
		channel:              pubChannel,
		allowedCountriesCode: []string{countrySuccess},
	}

	ctx := context.WithValue(context.Background(), domain.CtxUserIPKey, "1.1.1.1")
	for i := range updateCases {
		if !updateCases[i].success {
			company.event = PublisherMock(func(_ string, _ []byte) error {
				t.Errorf("should not be called for case: %s:%s",
					updateCases[i].c.Name,
					updateCases[i].c.Code)
				return nil
			})
		} else {
			company.event = nil
		}
		err := company.Update(ctx, updateCases[i].c.Name, updateCases[i].c.Code, updateCases[i].c)
		if updateCases[i].success && err != nil {
			t.Error(err.Error())
		}
		if !updateCases[i].success && err == nil {
			t.Errorf("case %s should be fail but not", updateCases[i].country)
		}
	}
}
