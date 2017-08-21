package digitalocean

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	doLiveTest  bool
	doAuthToken string
	doDomain    string
	doIP        string
)

func init() {
	doAuthToken = os.Getenv("DO_AUTH_TOKEN")
	doDomain = os.Getenv("DO_DOMAIN")
	doIP = os.Getenv("DO_IP")
	if len(doAuthToken) > 0 && len(doDomain) > 0 {
		doLiveTest = true
	}
}

func restoreDigitalOceanEnv() {
	os.Setenv("DO_AUTH_TOKEN", doAuthToken)
}

func TestNewDNSProviderValid(t *testing.T) {
	os.Setenv("DO_AUTH_TOKEN", "")
	_, err := NewDNSProviderCredentials(Options{AuthToken: "123"})
	assert.NoError(t, err)
	restoreDigitalOceanEnv()
}

func TestNewDNSProviderValidEnv(t *testing.T) {
	os.Setenv("DO_AUTH_TOKEN", "123")
	_, err := NewDNSProvider()
	assert.NoError(t, err)
	restoreDigitalOceanEnv()
}

func TestNewDNSProviderMissingCredErr(t *testing.T) {
	os.Setenv("DO_AUTH_TOKEN", "")
	_, err := NewDNSProvider()
	assert.EqualError(t, err, "DigitalOcean credentials missing")
	restoreDigitalOceanEnv()
}

func TestDigitalOceanEnsureARecord(t *testing.T) {
	if !doLiveTest {
		t.Skip("skipping live test")
	}

	provider, err := NewDNSProviderCredentials(Options{AuthToken: doAuthToken})
	assert.NoError(t, err)

	err = provider.EnsureARecord(doDomain, doIP)
	assert.NoError(t, err)
}

func TestDigitalOceanDeleteARecords(t *testing.T) {
	if !doLiveTest {
		t.Skip("skipping live test")
	}

	provider, err := NewDNSProviderCredentials(Options{AuthToken: doAuthToken})
	assert.NoError(t, err)

	err = provider.DeleteARecords(doDomain)
	assert.NoError(t, err)
}

func TestDigitalOceanDeleteARecord(t *testing.T) {
	if !doLiveTest {
		t.Skip("skipping live test")
	}

	provider, err := NewDNSProviderCredentials(Options{AuthToken: doAuthToken})
	assert.NoError(t, err)

	err = provider.DeleteARecord(doDomain, doIP)
	assert.NoError(t, err)
}
