package domain

import "context"

type CompanyWriter interface {
	Create(ctx context.Context, company *Company) error
	Update(ctx context.Context, oldName, oldCode string, company *Company) error
}

type CompanyReader interface {
	Get(ctx context.Context, name, code string) (Company, error)
	GetMany(ctx context.Context, filter *FilterOptions) ([]Company, error)
}

type CompanyDeleter interface {
	Delete(ctx context.Context, name, code string) error
}

// ICompany - main interface for Company data type for all CRUD operations
type ICompany interface {
	CompanyReader
	CompanyWriter
	CompanyDeleter
}

type Publisher interface {
	Publish(subj string, data []byte) error
}

type CountryResolver interface {
	Resolve(ip string) (string, error)
}
