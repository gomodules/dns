package googlecloud

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
)

var (
	gcloudLiveTest bool
	gcloudProject  string
	gcloudDomain   string
	gcloudIP       string
)

func init() {
	gcloudProject = os.Getenv("GCE_PROJECT")
	gcloudDomain = os.Getenv("GCE_DOMAIN")
	gcloudIP = os.Getenv("GCE_IP")
	_, err := google.DefaultClient(context.Background(), dns.NdevClouddnsReadwriteScope)
	if err == nil && len(gcloudProject) > 0 && len(gcloudDomain) > 0 {
		gcloudLiveTest = true
	}
}

func restoreGCloudEnv() {
	os.Setenv("GCE_PROJECT", gcloudProject)
}

func TestNewDNSProviderValid(t *testing.T) {
	if !gcloudLiveTest {
		t.Skip("skipping live test (requires credentials)")
	}
	os.Setenv("GCE_PROJECT", "")
	_, err := NewDNSProviderCredentials(Options{Project: "my-project"})
	assert.NoError(t, err)
	restoreGCloudEnv()
}

func TestNewDNSProviderValidEnv(t *testing.T) {
	if !gcloudLiveTest {
		t.Skip("skipping live test (requires credentials)")
	}
	os.Setenv("GCE_PROJECT", "my-project")
	_, err := NewDNSProvider()
	assert.NoError(t, err)
	restoreGCloudEnv()
}

func TestNewDNSProviderMissingCredErr(t *testing.T) {
	os.Setenv("GCE_PROJECT", "")
	_, err := NewDNSProvider()
	assert.EqualError(t, err, "Google Cloud project name missing")
	restoreGCloudEnv()
}

func TestLiveGoogleCloudEnsureARecord(t *testing.T) {
	if !gcloudLiveTest {
		t.Skip("skipping live test")
	}

	provider, err := NewDNSProviderCredentials(Options{Project: gcloudProject})
	assert.NoError(t, err)

	err = provider.EnsureARecord(gcloudDomain, gcloudIP)
	assert.NoError(t, err)
}

func TestLiveGoogleCloudDeleteARecords(t *testing.T) {
	if !gcloudLiveTest {
		t.Skip("skipping live test")
	}

	provider, err := NewDNSProviderCredentials(Options{Project: gcloudProject})
	assert.NoError(t, err)

	err = provider.DeleteARecords(gcloudDomain)
	assert.NoError(t, err)
}
