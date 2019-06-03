// Factory for DNS providers
package dns

import (
	"fmt"
	"strings"

	"gomodules.xyz/dns/aws"
	"gomodules.xyz/dns/azure"
	"gomodules.xyz/dns/cloudflare"
	"gomodules.xyz/dns/digitalocean"
	"gomodules.xyz/dns/googlecloud"
	"gomodules.xyz/dns/linode"
	dp "gomodules.xyz/dns/provider"
	"gomodules.xyz/dns/vultr"
)

func Default(name string) (dp.Provider, error) {
	var err error
	var provider dp.Provider
	switch strings.ToLower(name) {
	case "azure":
		provider, err = azure.Default()
	case "cloudflare":
		provider, err = cloudflare.Default()
	case "digitalocean":
		provider, err = digitalocean.Default()
	case "gcloud", "googlecloud", "gce", "gke":
		provider, err = googlecloud.Default()
	case "linode":
		provider, err = linode.Default()
	case "aws", "route53":
		provider, err = aws.Default()
	case "vultr":
		provider, err = vultr.Default()
	default:
		err = fmt.Errorf("unrecognised DNS provider: %s", name)
	}
	return provider, err
}
