// Factory for DNS providers
package dns

import (
	"fmt"

	"github.com/appscode/go-dns/aws"
	"github.com/appscode/go-dns/azure"
	"github.com/appscode/go-dns/cloudflare"
	"github.com/appscode/go-dns/digitalocean"
	"github.com/appscode/go-dns/googlecloud"
	"github.com/appscode/go-dns/linode"
	dp "github.com/appscode/go-dns/provider"
	"github.com/appscode/go-dns/vultr"
)

func NewDNSProvider(name string) (dp.Provider, error) {
	var err error
	var provider dp.Provider
	switch name {
	case "azure":
		provider, err = azure.NewDNSProvider()
	case "cloudflare":
		provider, err = cloudflare.NewDNSProvider()
	case "digitalocean":
		provider, err = digitalocean.NewDNSProvider()
	case "gcloud":
		provider, err = googlecloud.NewDNSProvider()
	case "linode":
		provider, err = linode.NewDNSProvider()
	case "aws", "route53":
		provider, err = aws.NewDNSProvider()
	case "vultr":
		provider, err = vultr.NewDNSProvider()
	default:
		err = fmt.Errorf("Unrecognised DNS provider: %s", name)
	}
	return provider, err
}
