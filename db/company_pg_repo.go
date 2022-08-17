package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/OleksiiKhanin/companysvc/domain"
	log "github.com/sirupsen/logrus"
	"strings"
)

func buildPGRequest(start int, options *domain.FilterOptions) (string, []any) {
	if options == nil || (len(options.Params) == 0 && options.Limit == nil) {
		return "true", []any{}
	}
	var limit string
	q := make([]string, 0, len(options.Params))
	values := make([]any, 0, len(options.Params))
	for k := range options.Params {
		start++
		q = append(q, fmt.Sprintf("%s ILIKE $%d", k, start))
		values = append(values, "%"+options.Params[k]+"%")
	}
	if options.Limit != nil {
		start++
		limit = fmt.Sprintf(" LIMIT $%d", start)
		values = append(values, *options.Limit)
	}
	return strings.Join(q, " AND ") + limit, values
}

type companyPostgreRepo struct {
	storage   *sql.DB
	l         *log.Logger
	logPrefix string
}

func NewCompanyPostgresRepo(storage *sql.DB, l *log.Logger) domain.ICompany {
	return &companyPostgreRepo{storage: storage, l: l, logPrefix: "Repository"}
}

func (c *companyPostgreRepo) Get(ctx context.Context, name, code string) (domain.Company, error) {
	query := "SELECT name, code, country, website, phone FROM companies WHERE name=$1 AND code=$2"
	c.l.Tracef("%s:Try execute: %s", c.logPrefix, query)
	var company domain.Company
	row := c.storage.QueryRowContext(ctx, query, name, code)
	err := row.Scan(
		&company.Name,
		&company.Code,
		&company.Country,
		&company.Website,
		&company.Phone,
	)
	if err != nil {
		return company, fmt.Errorf("get company from storage %w", err)
	}
	return company, nil
}

func (c *companyPostgreRepo) GetMany(ctx context.Context, options *domain.FilterOptions) ([]domain.Company, error) {
	whereStmt, values := buildPGRequest(0, options)
	query := fmt.Sprintf("SELECT name, code, country, website, phone FROM companies WHERE %s", whereStmt)
	c.l.Tracef("%s:Try execute: %s", c.logPrefix, query)
	rows, err := c.storage.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("get list of companies %w", err)
	}
	defer rows.Close()
	var companies []domain.Company
	for rows.Next() {
		var company domain.Company
		if err = rows.Scan(
			&company.Name,
			&company.Code,
			&company.Country,
			&company.Website,
			&company.Phone,
		); err == nil {
			companies = append(companies, company)
		}
	}
	return companies, nil
}

func (c *companyPostgreRepo) Create(ctx context.Context, company *domain.Company) error {
	query := "INSERT INTO companies (name, code, country, website, phone) VALUES ($1, $2, $3, $4, $5)"
	c.l.Tracef("%s:Try execute: %s", c.logPrefix, query)
	_, err := c.storage.ExecContext(ctx,
		query,
		company.Name,
		company.Code,
		company.Country,
		company.Website,
		company.Phone,
	)
	if err != nil {
		return fmt.Errorf("create company in storage: %w", err)
	}
	return nil
}

func (c *companyPostgreRepo) Update(ctx context.Context, oldName, oldCode string, company *domain.Company) error {
	query := "UPDATE companies SET name=$1, code=$2, country=$3, website=$4, phone=$5 WHERE name=$6 and code=$7"
	c.l.Tracef("%s:Try execute: %s", c.logPrefix, query)
	_, err := c.storage.ExecContext(ctx,
		query,
		company.Name,
		company.Code,
		company.Country,
		company.Website,
		company.Phone,
		oldName,
		oldCode,
	)
	if err != nil {
		return fmt.Errorf("update company in storage: %w", err)
	}
	return nil
}

func (c *companyPostgreRepo) Delete(ctx context.Context, name, code string) error {
	query := "DELETE FROM companies WHERE name=$1 and code=$2 LIMIT 1"
	c.l.Tracef("%s:Try execute: %s", c.logPrefix, query)
	_, err := c.storage.ExecContext(ctx, query, name, code)
	if err != nil {
		return fmt.Errorf("delete company vrom storage: %w", err)
	}
	return nil
}
