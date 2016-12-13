// Factory for DNS providers
package dns

import (
	"fmt"

	"github.com/appscode/go-dns/cloudflare"
	"github.com/appscode/go-dns/googlecloud"
	dp "github.com/appscode/go-dns/provider"
)

func NewDNSChallengeProviderByName(name string) (dp.Provider, error) {
	var err error
	var provider dp.Provider
	switch name {
	//case "azure":
	//	provider, err = azure.NewDNSProvider()
	//case "auroradns":
	//	provider, err = auroradns.NewDNSProvider()
	case "cloudflare":
		provider, err = cloudflare.NewDNSProvider()
	//case "digitalocean":
	//	provider, err = digitalocean.NewDNSProvider()
	//case "dnsimple":
	//	provider, err = dnsimple.NewDNSProvider()
	//case "dnsmadeeasy":
	//	provider, err = dnsmadeeasy.NewDNSProvider()
	//case "dnspod":
	//	provider, err = dnspod.NewDNSProvider()
	//case "dyn":
	//	provider, err = dyn.NewDNSProvider()
	//case "exoscale":
	//	provider, err = exoscale.NewDNSProvider()
	//case "gandi":
	//	provider, err = gandi.NewDNSProvider()
	case "gcloud":
		provider, err = googlecloud.NewDNSProvider()
	//case "linode":
	//	provider, err = linode.NewDNSProvider()
	//case "manual":
	//	provider, err = acme.NewDNSProviderManual()
	//case "namecheap":
	//	provider, err = namecheap.NewDNSProvider()
	//case "rackspace":
	//	provider, err = rackspace.NewDNSProvider()
	//case "route53":
	//	provider, err = route53.NewDNSProvider()
	//case "vultr":
	//	provider, err = vultr.NewDNSProvider()
	//case "ovh":
	//	provider, err = ovh.NewDNSProvider()
	//case "pdns":
	//	provider, err = pdns.NewDNSProvider()
	default:
		err = fmt.Errorf("Unrecognised DNS provider: %s", name)
	}
	return provider, err
}
