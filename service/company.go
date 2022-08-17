package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/OleksiiKhanin/companysvc/domain"
	log "github.com/sirupsen/logrus"
	"strings"
)

type companyService struct {
	domain.ICompany

	l                    *log.Logger
	event                domain.Publisher // can be nil
	locationClient       domain.CountryResolver
	allowedCountriesCode []string
	logPrefix            string
	channel              string
}

func NewCompanyService(
	company domain.ICompany,
	publisher domain.Publisher,
	loc domain.CountryResolver,
	channel string,
	l *log.Logger,
	countryCodes ...string,
) domain.ICompany {
	var c = companyService{
		ICompany:             company,
		l:                    l,
		event:                publisher,
		locationClient:       loc,
		allowedCountriesCode: countryCodes,
		logPrefix:            "companyService",
		channel:              channel,
	}
	return &c
}

func (c *companyService) checkUserIP(ctx context.Context) error {
	ip, ok := ctx.Value(domain.CtxUserIPKey).(string)
	if !ok {
		return fmt.Errorf("ip address must be defined")
	}
	code, err := c.locationClient.Resolve(strings.TrimSpace(ip))
	if err != nil {
		return fmt.Errorf("resolve ip %s: %w", ip, err)
	}
	code = strings.ToUpper(code)
	for i := range c.allowedCountriesCode {
		if strings.EqualFold(code, strings.ToUpper(c.allowedCountriesCode[i])) {
			return nil
		}
	}
	return fmt.Errorf("request not allowed")
}

func (c *companyService) Create(ctx context.Context, company *domain.Company) error {
	if err := c.checkUserIP(ctx); err != nil {
		return err
	}
	if err := c.ICompany.Create(ctx, company); err != nil {
		return err
	}
	if c.event != nil {
		data, err := json.Marshal(domain.Event{
			Type:    domain.CreateCompany,
			Subject: *company,
		})
		if err != nil {
			c.l.Infof("%s: create an insert event: %s", c.logPrefix, err.Error())
		} else if err := c.event.Publish(c.channel, data); err != nil {
			c.l.Infof("%s: publish an insert event: %s", c.logPrefix, err.Error())
		}
	}
	return nil
}

func (c *companyService) Delete(ctx context.Context, name, code string) error {
	if err := c.checkUserIP(ctx); err != nil {
		return err
	}
	if err := c.ICompany.Delete(ctx, name, code); err != nil {
		return err
	}
	if c.event != nil {
		data, err := json.Marshal(domain.Event{
			Type:    domain.DeleteCompany,
			OldName: name,
			OldCode: code,
		})
		if err != nil {
			c.l.Infof("%s: create a delete event: %s", c.logPrefix, err.Error())
		} else if err := c.event.Publish(c.channel, data); err != nil {
			c.l.Infof("%s: publish a delete event: %s", c.logPrefix, err.Error())
		}
	}
	return nil
}

func (c *companyService) Update(ctx context.Context, oldName, oldCode string, company *domain.Company) error {
	if err := c.ICompany.Update(ctx, oldName, oldCode, company); err != nil {
		return err
	}
	if c.event != nil {
		data, err := json.Marshal(domain.Event{
			Type:    domain.UpdateCompany,
			Subject: *company,
			OldName: oldName,
			OldCode: oldCode,
		})
		if err != nil {
			c.l.Infof("%s: create an update event: %s", c.logPrefix, err.Error())
		} else if err := c.event.Publish(c.channel, data); err != nil {
			c.l.Infof("%s: publish an update event: %s", c.logPrefix, err.Error())
		}
	}
	return nil
}
