// Package linode implements a DNS provider for solving the DNS-01 challenge
// using Linode DNS.
package linode

import (
	"errors"
	"log"
	"strings"
	"time"

	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/linode/linodego"
	"github.com/xenolf/lego/acme"
	"golang.org/x/oauth2"
	dp "gomodules.xyz/dns/provider"
	"gomodules.xyz/envconfig"
)

const (
	dnsMinTTLSecs      = 300
	dnsUpdateFreqMins  = 15
	dnsUpdateFudgeSecs = 120
)

type hostedZoneInfo struct {
	domainId     int
	resourceName string
}

// DNSProvider implements the acme.ChallengeProvider interface.
type DNSProvider struct {
	linode *linodego.Client
}

type Options struct {
	ApiKey string `json:"api_key" envconfig:"LINODE_API_KEY" form:"linode_api_key"`
}

var _ dp.Provider = &DNSProvider{}

// NewDNSProvider returns a DNSProvider instance configured for Linode.
// Credentials must be passed in the environment variable: LINODE_API_KEY.
func Default() (*DNSProvider, error) {
	var opt Options
	err := envconfig.Process("", &opt)
	if err != nil {
		return nil, err
	}
	return New(opt)
}

// NewDNSProviderCredentials uses the supplied credentials to return a
// DNSProvider instance configured for Linode.
func New(opt Options) (*DNSProvider, error) {
	if len(opt.ApiKey) == 0 {
		return nil, errors.New("Linode credentials missing")
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: opt.ApiKey,
	})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}
	client := linodego.NewClient(oauth2Client)
	return &DNSProvider{
		linode: &client,
	}, nil
}

// Timeout returns the timeout and interval to use when checking for DNS
// propagation.  Adjusting here to cope with spikes in propagation times.
func (p *DNSProvider) Timeout() (timeout, interval time.Duration) {
	// Since Linode only updates their zone files every X minutes, we need
	// to figure out how many minutes we have to wait until we hit the next
	// interval of X.  We then wait another couple of minutes, just to be
	// safe.  Hopefully at some point during all of this, the record will
	// have propagated throughout Linode's network.
	minsRemaining := dnsUpdateFreqMins - (time.Now().Minute() % dnsUpdateFreqMins)

	timeout = (time.Duration(minsRemaining) * time.Minute) +
		(dnsMinTTLSecs * time.Second) +
		(dnsUpdateFudgeSecs * time.Second)
	interval = 15 * time.Second
	return
}

func (p *DNSProvider) EnsureARecord(domain string, ip string) error {
	zone, err := p.getHostedZoneInfo(acme.ToFqdn(domain))
	if err != nil {
		return err
	}

	jsonFilter, err := json.Marshal(map[string]string{"type": "A"})
	if err != nil {
		return err
	}

	records, err := p.linode.ListDomainRecords(context.TODO(), zone.domainId, linodego.NewListOptions(0, string(jsonFilter)))
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Type == "A" && record.Name == zone.resourceName && record.Target == ip {
			log.Println("DNS is already configured. No DNS related change is necessary.")
			return nil
		}
	}
	_, err = p.linode.CreateDomainRecord(context.TODO(), zone.domainId, linodego.DomainRecordCreateOptions{
		Type:   "A",
		Target: ip,
		Name:   zone.resourceName,
		TTLSec: 300,
	})
	return err
}

func (p *DNSProvider) DeleteARecords(domain string) error {
	zone, err := p.getHostedZoneInfo(acme.ToFqdn(domain))
	if err != nil {
		return err
	}

	jsonFilter, err := json.Marshal(map[string]string{"type": "A"})
	if err != nil {
		return err
	}

	records, err := p.linode.ListDomainRecords(context.TODO(), zone.domainId, linodego.NewListOptions(0, string(jsonFilter)))
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Type == "A" && record.Name == zone.resourceName {
			err = p.linode.DeleteDomainRecord(context.TODO(), zone.domainId, record.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *DNSProvider) DeleteARecord(domain string, ip string) error {
	zone, err := p.getHostedZoneInfo(acme.ToFqdn(domain))
	if err != nil {
		return err
	}

	jsonFilter, err := json.Marshal(map[string]string{"type": "A"})
	if err != nil {
		return err
	}

	records, err := p.linode.ListDomainRecords(context.TODO(), zone.domainId, linodego.NewListOptions(0, string(jsonFilter)))
	if err != nil {
		return err
	}
	for _, record := range records {
		if record.Type == "A" && record.Name == zone.resourceName && record.Target == ip {
			err = p.linode.DeleteDomainRecord(context.TODO(), zone.domainId, record.ID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *DNSProvider) getHostedZoneInfo(fqdn string) (*hostedZoneInfo, error) {
	// Lookup the zone that handles the specified FQDN.
	authZone, err := acme.FindZoneByFqdn(fqdn, acme.RecursiveNameservers)
	if err != nil {
		return nil, err
	}
	resourceName := strings.TrimSuffix(fqdn, "."+authZone)

	jsonFilter, err := json.Marshal(map[string]string{"domain": acme.UnFqdn(authZone)})
	if err != nil {
		return nil, err
	}
	// Query the authority zone.
	domains, err := p.linode.ListDomains(context.TODO(), linodego.NewListOptions(0, string(jsonFilter)))
	if err != nil {
		return nil, err
	}
	if len(domains) == 0 {
		return nil, fmt.Errorf("domain %s not found", fqdn)
	}

	return &hostedZoneInfo{
		domainId:     domains[0].ID,
		resourceName: resourceName,
	}, nil
}
