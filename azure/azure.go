// Package azure implements a DNS provider for solving the DNS-01
// challenge using azure DNS.
// Azure doesn't like trailing dots on domain names, most of the acme code does.
package azure

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2017-09-01/dns"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/xenolf/lego/acme"
	dp "gomodules.xyz/dns/provider"
	"gomodules.xyz/envconfig"
)

// DNSProvider is an implementation of the acme.ChallengeProvider interface
type DNSProvider struct {
	opt Options
}

var _ dp.Provider = &DNSProvider{}

type Options struct {
	TenantId       string `json:"tenant_id" envconfig:"AZURE_TENANT_ID" form:"azure_tenant_id"`
	SubscriptionId string `json:"subscription_id" envconfig:"AZURE_SUBSCRIPTION_ID" form:"azure_subscription_id"`
	ClientId       string `json:"client_id" envconfig:"AZURE_CLIENT_ID" form:"azure_client_id"`
	ClientSecret   string `json:"client_secret" envconfig:"AZURE_CLIENT_SECRET" form:"azure_client_secret"`
	ResourceGroup  string `json:"resource_group" envconfig:"AZURE_RESOURCE_GROUP" form:"azure_resource_group"`
}

// NewDNSProvider returns a DNSProvider instance configured for azure.
// Credentials must be passed in the environment variables: AZURE_CLIENT_ID,
// AZURE_CLIENT_SECRET, AZURE_SUBSCRIPTION_ID, AZURE_TENANT_ID
func Default() (*DNSProvider, error) {
	var opt Options
	if err := envconfig.Process("", &opt); err != nil {
		return nil, err
	}
	return New(opt)
}

// NewDNSProviderCredentials uses the supplied credentials to return a
// DNSProvider instance configured for azure.
func New(opt Options) (*DNSProvider, error) {
	if opt.ClientId == "" || opt.ClientSecret == "" || opt.SubscriptionId == "" || opt.TenantId == "" || opt.ResourceGroup == "" {
		return nil, fmt.Errorf("Azure configuration missing")
	}

	return &DNSProvider{opt: opt}, nil
}

// Timeout returns the timeout and interval to use when checking for DNS
// propagation. Adjusting here to cope with spikes in propagation times.
func (c *DNSProvider) Timeout() (timeout, interval time.Duration) {
	return 120 * time.Second, 2 * time.Second
}

func (c *DNSProvider) EnsureARecord(domain string, ip string) error {
	fqdn := acme.ToFqdn(domain)
	zone, err := c.getHostedZoneID(fqdn)
	if err != nil {
		return err
	}

	spt, err := c.newServicePrincipalTokenFromCredentials(azure.PublicCloud.ResourceManagerEndpoint)
	rsc := dns.NewRecordSetsClient(c.opt.SubscriptionId)
	rsc.Authorizer = autorest.NewBearerAuthorizer(spt)
	relative := toRelativeRecord(fqdn, acme.ToFqdn(zone))
	rs, err := rsc.Get(context.TODO(), c.opt.ResourceGroup, zone, relative, "A")
	found, err := c.checkResourceExistsFromError(err)
	if err != nil {
		return err
	}

	records := make([]dns.ARecord, 0)
	if found {
		records = *rs.ARecords
		for _, record := range records {
			if *record.Ipv4Address == ip {
				log.Println("DNS is already configured. No DNS related change is necessary.")
				return nil
			}
		}
	}
	records = append(records, dns.ARecord{
		Ipv4Address: &ip,
	})

	rec := dns.RecordSet{
		Name: &relative,
		RecordSetProperties: &dns.RecordSetProperties{
			TTL:      to.Int64Ptr(300),
			ARecords: &records,
		},
	}
	_, err = rsc.CreateOrUpdate(context.TODO(), c.opt.ResourceGroup, zone, relative, dns.A, rec, "", "")
	return err
}

func (c *DNSProvider) DeleteARecords(domain string) error {
	fqdn := acme.ToFqdn(domain)
	zone, err := c.getHostedZoneID(fqdn)
	if err != nil {
		return err
	}

	spt, err := c.newServicePrincipalTokenFromCredentials(azure.PublicCloud.ResourceManagerEndpoint)
	rsc := dns.NewRecordSetsClient(c.opt.SubscriptionId)
	rsc.Authorizer = autorest.NewBearerAuthorizer(spt)
	relative := toRelativeRecord(fqdn, acme.ToFqdn(zone))
	_, err = rsc.Delete(context.TODO(), c.opt.ResourceGroup, zone, relative, "A", "")

	//resp, err := rsc.ListByType(c.resourceGroup, zone, "A", nil)
	//if err != nil {
	//	return err
	//}
	//for _, record := range (*resp.Value) {
	//	rsc.Delete(c.resourceGroup, zone, record.Name, "A", "")
	//}
	return err
}

func (c *DNSProvider) DeleteARecord(domain string, ip string) error {
	fqdn := acme.ToFqdn(domain)
	zone, err := c.getHostedZoneID(fqdn)
	if err != nil {
		return err
	}

	spt, err := c.newServicePrincipalTokenFromCredentials(azure.PublicCloud.ResourceManagerEndpoint)
	rsc := dns.NewRecordSetsClient(c.opt.SubscriptionId)
	rsc.Authorizer = autorest.NewBearerAuthorizer(spt)
	relative := toRelativeRecord(fqdn, acme.ToFqdn(zone))

	rs, err := rsc.Get(context.TODO(), c.opt.ResourceGroup, zone, relative, "A")
	if found, err := c.checkResourceExistsFromError(err); err != nil {
		return err
	} else if !found {
		log.Println("No record matched")
		return nil
	}

	records := *rs.ARecords
	updatedRecords := make([]dns.ARecord, 0)
	// make a new list by removing matched records
	for _, item := range records {
		if *item.Ipv4Address != ip {
			updatedRecords = append(updatedRecords, item)
		}
	}

	if len(records) == len(updatedRecords) {
		log.Println("No record matched")
		return nil
	}
	if len(updatedRecords) == 0 { // if all records matched, delete the recordset
		_, err = rsc.Delete(context.TODO(), c.opt.ResourceGroup, zone, relative, "A", "")
	} else { // update the recordset with new list
		rs.ARecords = &updatedRecords
		_, err = rsc.Update(context.TODO(), c.opt.ResourceGroup, zone, relative, "A", rs, "")
	}
	return err
}

// Returns the relative record to the domain
func toRelativeRecord(domain, zone string) string {
	return acme.UnFqdn(strings.TrimSuffix(domain, zone))
}

// Checks that azure has a zone for this domain name.
func (c *DNSProvider) getHostedZoneID(fqdn string) (string, error) {
	// return "learnappscode.com", nil  // only for test

	authZone, err := acme.FindZoneByFqdn(fqdn, acme.RecursiveNameservers)
	if err != nil {
		return "", err
	}

	// Now we want to to Azure and get the zone.
	spt, err := c.newServicePrincipalTokenFromCredentials(azure.PublicCloud.ResourceManagerEndpoint)
	dc := dns.NewZonesClient(c.opt.SubscriptionId)
	dc.Authorizer = autorest.NewBearerAuthorizer(spt)
	zone, err := dc.Get(context.TODO(), c.opt.ResourceGroup, acme.UnFqdn(authZone))

	if err != nil {
		return "", err
	}

	// zone.Name shouldn't have a trailing dot(.)
	return to.String(zone.Name), nil
}

// NewServicePrincipalTokenFromCredentials creates a new ServicePrincipalToken using values of the
// passed credentials map.
func (c *DNSProvider) newServicePrincipalTokenFromCredentials(scope string) (*adal.ServicePrincipalToken, error) {
	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, c.opt.TenantId)
	if err != nil {
		panic(err)
	}
	return adal.NewServicePrincipalToken(*oauthConfig, c.opt.ClientId, c.opt.ClientSecret, scope)
}

// checkExistsFromError inspects an error and returns a true if err is nil,
// false if error is an autorest.Error with StatusCode=404 and will return the
// error back if error is another status code or another type of error.
// ref: https://github.com/kubernetes/kubernetes/blob/54195d590f03a544d78b4449b2fbafaa258fd6df/pkg/cloudprovider/providers/azure/azure_wrap.go#L28
func (c *DNSProvider) checkResourceExistsFromError(err error) (bool, error) {
	if err == nil {
		return true, nil
	}
	v, ok := err.(autorest.DetailedError)
	if ok && v.StatusCode == http.StatusNotFound {
		return false, nil
	}
	return false, v
}
