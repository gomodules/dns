[![Go Report Card](https://goreportcard.com/badge/github.com/appscode/go-dns)](https://goreportcard.com/report/github.com/appscode/go-dns)

[Website](https://appscode.com) • [Slack](https://slack.appscode.com) • [Twitter](https://twitter.com/AppsCodeHQ)

# go-dns
Unified DNS API client for GOlang. See here for the documentation of [common provider interface](https://godoc.org/github.com/appscode/go-dns/provider).
```go
type Provider interface {
	EnsureARecord(domain string, ip string) error
	DeleteARecord(domain string, ip string) error
	DeleteARecords(domain string) error
}
```

### Supported DNS Providers
- [x] AWS Route53
- [x] Azure
- [x] Cloudflare
- [x] DigitalOcean
- [x] Google Cloud DNS
- [x] Linode
- [x] Vultr

### Acknowledgement
The initial implementation of this library was forked from https://github.com/xenolf/lego/tree/master/providers/dns
