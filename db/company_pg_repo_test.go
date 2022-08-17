package db

import (
	"context"
	"github.com/OleksiiKhanin/companysvc/domain"
	log "github.com/sirupsen/logrus"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func buildTestCases() map[string]*domain.FilterOptions {
	limit := 0
	return map[string]*domain.FilterOptions{
		"true": nil,
		" LIMIT $1": &domain.FilterOptions{
			Limit: &limit,
		},
		"name ILIKE $1 LIMIT $2": &domain.FilterOptions{
			Limit: &limit,
			Params: map[string]string{
				"name": "",
			},
		},
		"code ILIKE $1 LIMIT $2": &domain.FilterOptions{
			Limit: &limit,
			Params: map[string]string{
				"code": "",
			},
		},
	}
}

func TestBuildPGRequest(t *testing.T) {
	testCases := buildTestCases()
	for want := range testCases {
		got, _ := buildPGRequest(0, testCases[want])
		if !strings.EqualFold(got, want) {
			t.Errorf("want: %s but get: %s", want, got)
		}
	}
}

func TestCompanyPostgreRepoCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO companies").
		WithArgs("test", "test_code", "UA", "", "").
		WillReturnResult(sqlmock.NewResult(1, 1))
	company := NewCompanyPostgresRepo(db, log.StandardLogger())
	// now we execute our method
	err = company.Create(context.Background(), &domain.Company{Name: "test", Code: "test_code", Country: "UA"})
	if err != nil {
		t.Errorf("error was not expected while create company: %s", err)
	}
}

func TestCompanyPostgreRepoUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE companies SET").
		WithArgs("test", "test_code", "UA", "example.com", "", "old_test", "old_test_code").
		WillReturnResult(sqlmock.NewResult(1, 1))

	company := NewCompanyPostgresRepo(db, log.StandardLogger())
	// now we execute our method
	err = company.Update(context.Background(), "old_test", "old_test_code",
		&domain.Company{Name: "test", Code: "test_code", Country: "UA", Website: "example.com"},
	)
	if err != nil {
		t.Errorf("error was not expected while update company: %s", err)
	}
}

func TestCompanyPostgreRepoDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM companies").
		WithArgs("test", "test_code").
		WillReturnResult(sqlmock.NewResult(1, 1))

	company := NewCompanyPostgresRepo(db, log.StandardLogger())
	// now we execute our method
	err = company.Delete(context.Background(), "test", "test_code")
	if err != nil {
		t.Errorf("error was not expected while delete company: %s", err)
	}
}
