package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/provider"
	"github.com/outscale/terraform-provider-outscale/version"
)

func TestProvider(t *testing.T) {
	if err := provider.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	//nolint:staticcheck
	var _ *schema.Provider = provider.Provider()
}

func TestMuxSchema(t *testing.T) {
	upgradedSdkServer, err := tf5to6server.UpgradeServer(
		t.Context(),
		provider.Provider().GRPCProvider,
	)
	if err != nil {
		t.Fatalf("failed to upgrade SDK v2 provider to v6: %s", err)
	}

	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(provider.New(version.GetVersion())),
		func() tfprotov6.ProviderServer {
			return upgradedSdkServer
		},
	}

	_, err = tf6muxserver.NewMuxServer(t.Context(), providers...)
	if err != nil {
		t.Fatalf("mux schema mismatch between SDK and Framework providers: %s", err)
	}
}
