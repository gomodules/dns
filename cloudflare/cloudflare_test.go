package cloudflare

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cflareLiveTest bool
	cflareEmail    string
	cflareAPIKey   string
	cflareDomain   string
	cflareIP       string
)

func init() {
	cflareEmail = os.Getenv("CLOUDFLARE_EMAIL")
	cflareAPIKey = os.Getenv("CLOUDFLARE_API_KEY")
	cflareDomain = os.Getenv("CLOUDFLARE_DOMAIN")
	cflareIP = os.Getenv("CLOUDFLARE_IP")
	if len(cflareEmail) > 0 && len(cflareAPIKey) > 0 && len(cflareDomain) > 0 {
		cflareLiveTest = true
	}
}

func restoreCloudFlareEnv() {
	os.Setenv("CLOUDFLARE_EMAIL", cflareEmail)
	os.Setenv("CLOUDFLARE_API_KEY", cflareAPIKey)
}

func TestNewDNSProviderValid(t *testing.T) {
	os.Setenv("CLOUDFLARE_EMAIL", "")
	os.Setenv("CLOUDFLARE_API_KEY", "")
	_, err := New(Options{
		Email:  "123",
		APIKey: "123",
	})
	assert.NoError(t, err)
	restoreCloudFlareEnv()
}

func TestNewDNSProviderValidEnv(t *testing.T) {
	os.Setenv("CLOUDFLARE_EMAIL", "test@example.com")
	os.Setenv("CLOUDFLARE_API_KEY", "123")
	_, err := Default()
	assert.NoError(t, err)
	restoreCloudFlareEnv()
}

func TestNewDNSProviderMissingCredErr(t *testing.T) {
	os.Setenv("CLOUDFLARE_EMAIL", "")
	os.Setenv("CLOUDFLARE_API_KEY", "")
	_, err := Default()
	assert.EqualError(t, err, "CloudFlare credentials missing")
	restoreCloudFlareEnv()
}

func TestCloudFlareEnsureARecord(t *testing.T) {
	if !cflareLiveTest {
		t.Skip("skipping live test")
	}

	provider, err := New(Options{
		Email:  cflareEmail,
		APIKey: cflareAPIKey,
	})
	assert.NoError(t, err)

	err = provider.EnsureARecord(cflareDomain, cflareIP)
	assert.NoError(t, err)
}

func TestCloudFlareDeleteARecords(t *testing.T) {
	if !cflareLiveTest {
		t.Skip("skipping live test")
	}

	provider, err := New(Options{
		Email:  cflareEmail,
		APIKey: cflareAPIKey,
	})
	assert.NoError(t, err)

	err = provider.DeleteARecords(cflareDomain)
	assert.NoError(t, err)
}

func TestCloudFlareDeleteARecord(t *testing.T) {
	if !cflareLiveTest {
		t.Skip("skipping live test")
	}

	provider, err := New(Options{
		Email:  cflareEmail,
		APIKey: cflareAPIKey,
	})
	assert.NoError(t, err)

	err = provider.DeleteARecord(cflareDomain, cflareIP)
	assert.NoError(t, err)
}
