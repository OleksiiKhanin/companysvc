package service

import (
	"fmt"
	"github.com/OleksiiKhanin/companysvc/domain"
	retry "github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
	"io"
)

type ipAPIResolver struct {
	client    *retry.Client
	url       string
	l         *log.Logger
	logPrefix string
}

func GetResolverIPAPI(url string, requestAttempt int, l *log.Logger) domain.CountryResolver {
	client := retry.NewClient()
	client.RetryMax = requestAttempt
	client.Logger = l
	return &ipAPIResolver{client: client, url: url, l: l, logPrefix: "countryResolve"}
}

func (i *ipAPIResolver) Resolve(ip string) (string, error) {
	//GET https://ipapi.co/{ip}/country_code/
	uri := fmt.Sprintf("%s/%s/country_code/", i.url, ip)
	resp, err := i.client.Get(uri)
	if err != nil {
		return "", fmt.Errorf("get country from remote api: %w", err)
	}
	defer resp.Body.Close()
	code, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read responce body: %w", err)
	}
	i.l.Tracef("%s: result country code: %s", i.logPrefix, code)
	return string(code), nil
}
