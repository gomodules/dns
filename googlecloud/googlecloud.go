// Package googlecloud implements a DNS provider for solving the DNS-01
// challenge using Google Cloud DNS.
package googlecloud

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gomodules.xyz/x/strings"
	"github.com/xenolf/lego/acme"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	dp "gomodules.xyz/dns/provider"
	"gomodules.xyz/envconfig"
	"google.golang.org/api/dns/v1"
)

// DNSProvider is an implementation of the DNSProvider interface.
type DNSProvider struct {
	project string
	client  *dns.Service
}

type Options struct {
	Project        string `json:"project" envconfig:"GCE_PROJECT"  form:"gcloud_project"`
	CredentialFile string `json:"credential_file" ignore:"true" form:"-"`
	CredentialJson string `json:"-" ignore:"true" form:"gcloud_credential_json"`
	JsonKey        []byte `json:"-" ignore:"true"  form:"-"`
}

var _ dp.Provider = &DNSProvider{}

// NewDNSProvider returns a DNSProvider instance configured for Google Cloud
// DNS. Credentials must be passed in the environment variable: GCE_PROJECT.
func Default() (*DNSProvider, error) {
	var opt Options
	err := envconfig.Process("", &opt)
	if err != nil {
		return nil, err
	}
	return New(opt)
}

// NewDNSProviderCredentials uses the supplied credentials to return a
// DNSProvider instance configured for Google Cloud DNS.
func New(opt Options) (*DNSProvider, error) {
	if opt.Project == "" {
		return nil, fmt.Errorf("Google Cloud project name missing")
	}

	var client *http.Client
	var err error
	if opt.CredentialFile != "" {
		opt.JsonKey, err = os.ReadFile(opt.CredentialFile)
		if err != nil {
			return nil, err
		}
	}
	if opt.JsonKey == nil {
		client, err = google.DefaultClient(context.Background(), dns.NdevClouddnsReadwriteScope)
		if err != nil {
			return nil, fmt.Errorf("Unable to get Google Cloud client: %v", err)
		}
	} else {
		conf, err := google.JWTConfigFromJSON(opt.JsonKey, dns.NdevClouddnsReadwriteScope)
		if err != nil {
			return nil, fmt.Errorf("Unable to load JWT config from Google Service Account file: %v", err)
		}
		client = conf.Client(context.Background())
	}

	svc, err := dns.New(client)
	if err != nil {
		return nil, fmt.Errorf("Unable to create Google Cloud DNS service: %v", err)
	}
	return &DNSProvider{
		project: opt.Project,
		client:  svc,
	}, nil
}

// Timeout customizes the timeout values used by the ACME package for checking
// DNS record validity.
func (c *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return 180 * time.Second, 5 * time.Second
}

func (c *DNSProvider) EnsureARecord(domain string, ip string) error {
	zone, err := c.getHostedZone(domain)
	if err != nil {
		return err
	}

	r1, err := c.client.ResourceRecordSets.List(c.project, zone).
		Name(acme.ToFqdn(domain)).
		Type("A").
		Do()
	log.Println("Retrieved A record", r1, err)
	if err != nil {
		return err
	}

	log.Println("Updating A record for cluster", domain)
	changes := &dns.Change{}
	if len(r1.Rrsets) == 0 || !strings.Contains(r1.Rrsets[0].Rrdatas, ip) {
		ips := []string{ip}
		if len(r1.Rrsets) > 0 {
			ips = append(ips, r1.Rrsets[0].Rrdatas...)
		}
		log.Println("Adding A record ", []string{ip})
		changes.Additions = []*dns.ResourceRecordSet{
			{
				Name:    acme.ToFqdn(domain),
				Type:    "A",
				Ttl:     int64(300),
				Rrdatas: ips,
			},
		}
	}
	if len(r1.Rrsets) == 1 && !strings.Contains(r1.Rrsets[0].Rrdatas, ip) {
		log.Println("Deleting A record ", r1.Rrsets[0].Rrdatas)
		changes.Deletions = []*dns.ResourceRecordSet{
			{
				Name:    acme.ToFqdn(domain),
				Type:    "A",
				Ttl:     r1.Rrsets[0].Ttl,
				Rrdatas: r1.Rrsets[0].Rrdatas,
			},
		}
	}
	if len(changes.Additions)+len(changes.Deletions) == 0 {
		log.Println("DNS is already configured. No DNS related change is necessary.")
		return nil
	}
	r2, err := c.client.Changes.Create(c.project, zone, changes).Do()
	log.Println("Updated A record", r2, err)
	return err
}

func (c *DNSProvider) DeleteARecords(domain string) error {
	zone, err := c.getHostedZone(domain)
	if err != nil {
		return err
	}

	r1, err := c.client.ResourceRecordSets.List(c.project, zone).
		Name(acme.ToFqdn(domain)).
		Type("A").
		Do()
	if err != nil {
		return err
	}
	if len(r1.Rrsets) > 0 {
		changes := &dns.Change{
			Deletions: []*dns.ResourceRecordSet{
				{
					Name:    acme.ToFqdn(domain),
					Type:    "A",
					Ttl:     r1.Rrsets[0].Ttl,
					Rrdatas: r1.Rrsets[0].Rrdatas,
				},
			},
		}
		_, err = c.client.Changes.Create(c.project, zone, changes).Do()
		return err
	}
	return nil
}

func (c *DNSProvider) DeleteARecord(domain string, ip string) error {
	zone, err := c.getHostedZone(domain)
	if err != nil {
		return err
	}

	r1, err := c.client.ResourceRecordSets.List(c.project, zone).
		Name(acme.ToFqdn(domain)).
		Type("A").
		Do()
	if err != nil {
		return err
	}
	if len(r1.Rrsets) == 0 || !strings.Contains(r1.Rrsets[0].Rrdatas, ip) {
		log.Println("No record found")
		return nil
	}

	// create a new list by removing matched ips
	ips := make([]string, 0)
	for _, item := range r1.Rrsets[0].Rrdatas {
		if item != ip {
			ips = append(ips, item)
		}
	}
	// can't update record set, need to delete existing one and insert modified one
	changes := &dns.Change{
		Deletions: []*dns.ResourceRecordSet{
			{
				Name:    acme.ToFqdn(domain),
				Type:    "A",
				Ttl:     r1.Rrsets[0].Ttl,
				Rrdatas: r1.Rrsets[0].Rrdatas,
			},
		},
	}
	if len(ips) > 0 {
		changes.Additions = []*dns.ResourceRecordSet{
			{
				Name:    acme.ToFqdn(domain),
				Type:    "A",
				Ttl:     int64(300),
				Rrdatas: ips,
			},
		}
	}
	_, err = c.client.Changes.Create(c.project, zone, changes).Do()
	return err
}

// getHostedZone returns the managed-zone
func (c *DNSProvider) getHostedZone(domain string) (string, error) {
	authZone, err := acme.FindZoneByFqdn(acme.ToFqdn(domain), acme.RecursiveNameservers)
	if err != nil {
		return "", err
	}

	zones, err := c.client.ManagedZones.
		List(c.project).
		DnsName(authZone).
		Do()
	if err != nil {
		return "", fmt.Errorf("GoogleCloud API call failed: %v", err)
	}

	if len(zones.ManagedZones) == 0 {
		return "", fmt.Errorf("No matching GoogleCloud domain found for domain %s", authZone)
	}

	return zones.ManagedZones[0].Name, nil
}
