package linode

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	apiKey     string
	domain     string
	ip         string
	isTestLive bool
)

func init() {
	apiKey = os.Getenv("LINODE_API_KEY")
	domain = os.Getenv("LINODE_DOMAIN")
	ip = os.Getenv("LINODE_IP")
	if len(apiKey) > 0 && len(domain) > 0 {
		isTestLive = true
	}
}

func restoreEnv() {
	os.Setenv("LINODE_API_KEY", apiKey)
}

func TestNewDNSProviderWithEnv(t *testing.T) {
	os.Setenv("LINODE_API_KEY", "testing")
	defer restoreEnv()
	_, err := NewDNSProvider()
	assert.NoError(t, err)
}

func TestNewDNSProviderWithoutEnv(t *testing.T) {
	os.Setenv("LINODE_API_KEY", "")
	defer restoreEnv()
	_, err := NewDNSProvider()
	assert.EqualError(t, err, "Linode credentials missing")
}

func TestNewDNSProviderCredentialsWithKey(t *testing.T) {
	_, err := NewDNSProviderCredentials(Options{ApiKey: "testing"})
	assert.NoError(t, err)
}

func TestNewDNSProviderCredentialsWithoutKey(t *testing.T) {
	_, err := NewDNSProviderCredentials(Options{})
	assert.EqualError(t, err, "Linode credentials missing")
}

func TestLiveEnsureARecord(t *testing.T) {
	if !isTestLive {
		t.Skip("skipping live test")
	}

	provider, err := NewDNSProvider()
	assert.NoError(t, err)

	err = provider.EnsureARecord(domain, ip)
	assert.NoError(t, err)
}

func TestLiveDeleteARecords(t *testing.T) {
	if !isTestLive {
		t.Skip("skipping live test")
	}

	provider, err := NewDNSProvider()
	assert.NoError(t, err)

	err = provider.DeleteARecords(domain)
	assert.NoError(t, err)
}
