package dns

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gomodules.xyz/dns/digitalocean"
)

var (
	apiKey    string
	apiSecret string
)

func init() {
	apiKey = os.Getenv("DO_AUTH_TOKEN")
}

func restoreExoscaleEnv() {
	os.Setenv("DO_AUTH_TOKEN", apiKey)
}

func TestKnownDNSProviderSuccess(t *testing.T) {
	os.Setenv("DO_AUTH_TOKEN", "abc")
	provider, err := Default("digitalocean")
	assert.NoError(t, err)
	assert.NotNil(t, provider)
	if reflect.TypeOf(provider) != reflect.TypeOf(&digitalocean.DNSProvider{}) {
		t.Errorf("Not loaded correct DNS proviver: %v is not *digitalocean.DNSProvider", reflect.TypeOf(provider))
	}
	restoreExoscaleEnv()
}

func TestKnownDNSProviderError(t *testing.T) {
	os.Setenv("DO_AUTH_TOKEN", "")
	_, err := Default("digitalocean")
	assert.Error(t, err)
	restoreExoscaleEnv()
}

func TestUnknownDNSProvider(t *testing.T) {
	_, err := Default("foobar")
	assert.Error(t, err)
}
