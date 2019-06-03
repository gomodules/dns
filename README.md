[![Go Report Card](https://goreportcard.com/badge/gomodules.xyz/dns)](https://goreportcard.com/report/gomodules.xyz/dns)
[![Build Status](https://travis-ci.org/gomodules/dns.svg?branch=master)](https://travis-ci.org/gomodules/dns)

# go-dns
Unified DNS API client for GOlang. See here for the documentation of [common provider interface](https://godoc.org/gomodules.xyz/dns/provider).
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
