// Factory for DNS providers
package dns

import (
	"fmt"

	"github.com/appscode/go-dns/auroradns"
	"github.com/appscode/go-dns/azure"
	"github.com/appscode/go-dns/cloudflare"
	"github.com/appscode/go-dns/digitalocean"
	"github.com/appscode/go-dns/dnsimple"
	"github.com/appscode/go-dns/dnsmadeeasy"
	"github.com/appscode/go-dns/dnspod"
	"github.com/appscode/go-dns/dyn"
	"github.com/appscode/go-dns/exoscale"
	"github.com/appscode/go-dns/gandi"
	"github.com/appscode/go-dns/googlecloud"
	"github.com/appscode/go-dns/linode"
	"github.com/appscode/go-dns/namecheap"
	"github.com/appscode/go-dns/ovh"
	"github.com/appscode/go-dns/pdns"
	"github.com/appscode/go-dns/rackspace"
	"github.com/appscode/go-dns/route53"
	"github.com/appscode/go-dns/vultr"
	"github.com/xenolf/lego/acme"
)

func NewDNSChallengeProviderByName(name string) (acme.ChallengeProvider, error) {
	var err error
	var provider acme.ChallengeProvider
	switch name {
	case "azure":
		provider, err = azure.NewDNSProvider()
	case "auroradns":
		provider, err = auroradns.NewDNSProvider()
	case "cloudflare":
		provider, err = cloudflare.NewDNSProvider()
	case "digitalocean":
		provider, err = digitalocean.NewDNSProvider()
	case "dnsimple":
		provider, err = dnsimple.NewDNSProvider()
	case "dnsmadeeasy":
		provider, err = dnsmadeeasy.NewDNSProvider()
	case "dnspod":
		provider, err = dnspod.NewDNSProvider()
	case "dyn":
		provider, err = dyn.NewDNSProvider()
	case "exoscale":
		provider, err = exoscale.NewDNSProvider()
	case "gandi":
		provider, err = gandi.NewDNSProvider()
	case "gcloud":
		provider, err = googlecloud.NewDNSProvider()
	case "linode":
		provider, err = linode.NewDNSProvider()
	case "manual":
		provider, err = acme.NewDNSProviderManual()
	case "namecheap":
		provider, err = namecheap.NewDNSProvider()
	case "rackspace":
		provider, err = rackspace.NewDNSProvider()
	case "route53":
		provider, err = route53.NewDNSProvider()
	case "vultr":
		provider, err = vultr.NewDNSProvider()
	case "ovh":
		provider, err = ovh.NewDNSProvider()
	case "pdns":
		provider, err = pdns.NewDNSProvider()
	default:
		err = fmt.Errorf("Unrecognised DNS provider: %s", name)
	}
	return provider, err
}
