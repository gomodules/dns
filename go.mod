module gomodules.xyz/dns

go 1.12

require (
	github.com/Azure/azure-sdk-for-go v14.6.0+incompatible
	github.com/Azure/go-autorest v10.6.2+incompatible
	github.com/JamesClonk/vultr v0.0.0-20180805100820-267be0d362b6
	github.com/aws/aws-sdk-go v1.12.7
	github.com/cloudflare/cloudflare-go v0.8.0
	github.com/digitalocean/godo v1.1.0
	github.com/go-ini/ini v1.25.4 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20160202185014-0b12d6b521d8 // indirect
	github.com/juju/ratelimit v0.0.0-20151125201925-77ed1c8a0121 // indirect
	github.com/linode/linodego v0.9.0
	github.com/stretchr/testify v1.5.1
	github.com/tent/http-link-go v0.0.0-20130702225549-ac974c61c2f9 // indirect
	github.com/xenolf/lego v2.6.0+incompatible
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	gomodules.xyz/envconfig v1.3.0
	gomodules.xyz/x v0.0.0-20201105065653-91c568df6331
	google.golang.org/api v0.13.0
	gopkg.in/square/go-jose.v2 v2.3.1 // indirect
	launchpad.net/gocheck v0.0.0-20140225173054-000000000087 // indirect
)

replace github.com/xenolf/lego => github.com/appscode/lego v1.0.2-0.20180815012227-0b6736776fb2
