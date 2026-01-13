package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"strapi": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("STRAPI_ENDPOINT"); v == "" {
		t.Fatal("STRAPI_ENDPOINT must be set for acceptance tests")
	}
	if v := os.Getenv("STRAPI_API_TOKEN"); v == "" {
		t.Fatal("STRAPI_API_TOKEN must be set for acceptance tests")
	}
}
