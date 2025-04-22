package provider

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &CognitoIdpLinkProvider{}

type CognitoIdpLinkProvider struct {
	version string
	client  *cognitoidentityprovider.Client
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CognitoIdpLinkProvider{
			version: version,
		}
	}
}

func (p *CognitoIdpLinkProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cognito-idp-link"
	resp.Version = p.version
}

func (p *CognitoIdpLinkProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with AWS Cognito User Pools Identity Provider Link.",
		Attributes:  map[string]schema.Attribute{},
	}
}

func (p *CognitoIdpLinkProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if p.client != nil {
		resp.DataSourceData = p.client
		resp.ResourceData = p.client
		return
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to load AWS configuration",
			"Error occurred while loading AWS configuration: "+err.Error(),
		)
		return
	}

	client := cognitoidentityprovider.NewFromConfig(cfg)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *CognitoIdpLinkProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *CognitoIdpLinkProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCognitoIdpLinkResource,
	}
}

// SetClient sets a custom client for testing
func (p *CognitoIdpLinkProvider) SetClient(client *cognitoidentityprovider.Client) {
	p.client = client
}
