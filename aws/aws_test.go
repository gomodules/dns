package aws

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
)

var (
	route53Secret string
	route53Key    string
	route53Region string
	domain        string
	ip            string
	liveTest      bool
)

func init() {
	route53Key = os.Getenv("AWS_ACCESS_KEY_ID")
	route53Secret = os.Getenv("AWS_SECRET_ACCESS_KEY")
	route53Region = os.Getenv("AWS_REGION")

	domain = os.Getenv("AWS_DOMAIN")
	ip = os.Getenv("AWS_IP")
	if len(domain) > 0 {
		liveTest = true
	}
}

func restoreRoute53Env() {
	os.Setenv("AWS_ACCESS_KEY_ID", route53Key)
	os.Setenv("AWS_SECRET_ACCESS_KEY", route53Secret)
	os.Setenv("AWS_REGION", route53Region)
}

func TestCredentialsFromEnv(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", "123")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "123")
	os.Setenv("AWS_REGION", "us-east-1")

	config := &aws.Config{
		CredentialsChainVerboseErrors: aws.Bool(true),
	}

	sess := session.New(config)
	_, err := sess.Config.Credentials.Get()
	assert.NoError(t, err, "Expected credentials to be set from environment")

	restoreRoute53Env()
}

func TestRegionFromEnv(t *testing.T) {
	os.Setenv("AWS_REGION", "us-east-1")

	sess := session.New(aws.NewConfig())
	assert.Equal(t, "us-east-1", *sess.Config.Region, "Expected Region to be set from environment")

	restoreRoute53Env()
}

func TestLiveEnsureARecord(t *testing.T) {
	if !liveTest {
		t.Skip("skipping live test")
	}

	provider, err := NewDNSProvider()
	assert.NoError(t, err)

	err = provider.EnsureARecord(domain, ip)
	assert.NoError(t, err)
}

func TestLiveDeleteARecords(t *testing.T) {
	if !liveTest {
		t.Skip("skipping live test")
	}

	provider, err := NewDNSProvider()
	assert.NoError(t, err)

	err = provider.DeleteARecords(domain)
	assert.NoError(t, err)
}

func TestLiveDeleteARecord(t *testing.T) {
	if !liveTest {
		t.Skip("skipping live test")
	}

	provider, err := NewDNSProvider()
	assert.NoError(t, err)

	err = provider.DeleteARecord(domain, ip)
	assert.NoError(t, err)
}
