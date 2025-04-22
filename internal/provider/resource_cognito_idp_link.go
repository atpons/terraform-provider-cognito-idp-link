package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	cognitotypes "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &CognitoIdpLinkResource{}
var _ resource.ResourceWithImportState = &CognitoIdpLinkResource{}

type CognitoIdpLinkResource struct {
	client *cognitoidentityprovider.Client
}

type CognitoIdpLinkResourceModel struct {
	UserPoolId      types.String     `tfsdk:"user_pool_id"`
	DestinationUser *DestinationUser `tfsdk:"destination_user"`
	SourceUser      *SourceUser      `tfsdk:"source_user"`
	Id              types.String     `tfsdk:"id"`
}

type DestinationUser struct {
	ProviderName           types.String `tfsdk:"provider_name"`
	ProviderAttributeValue types.String `tfsdk:"provider_attribute_value"`
}

type SourceUser struct {
	ProviderName           types.String `tfsdk:"provider_name"`
	ProviderAttributeName  types.String `tfsdk:"provider_attribute_name"`
	ProviderAttributeValue types.String `tfsdk:"provider_attribute_value"`
}

func NewCognitoIdpLinkResource() resource.Resource {
	return &CognitoIdpLinkResource{}
}

func (r *CognitoIdpLinkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_link"
}

func (r *CognitoIdpLinkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages users to be linked in Cognito User Pools.",
		Attributes: map[string]schema.Attribute{
			"user_pool_id": schema.StringAttribute{
				Description: "Cognito User Pool ID",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"destination_user": schema.SingleNestedAttribute{
				Description: "Information of the destination user to be linked",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"provider_name": schema.StringAttribute{
						Description: "Provider name",
						Required:    true,
					},
					"provider_attribute_value": schema.StringAttribute{
						Description: "Provider attribute value",
						Required:    true,
					},
				},
			},
			"source_user": schema.SingleNestedAttribute{
				Description: "Information of the source user to be linked",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"provider_name": schema.StringAttribute{
						Description: "Provider name",
						Required:    true,
					},
					"provider_attribute_name": schema.StringAttribute{
						Description: "Provider attribute name",
						Required:    true,
					},
					"provider_attribute_value": schema.StringAttribute{
						Description: "Provider attribute value",
						Required:    true,
					},
				},
			},
			"id": schema.StringAttribute{
				Description: "Resource ID",
				Computed:    true,
			},
		},
	}
}

func (r *CognitoIdpLinkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cognitoidentityprovider.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid client type",
			fmt.Sprintf("Expected *cognitoidentityprovider.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *CognitoIdpLinkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CognitoIdpLinkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &cognitoidentityprovider.AdminLinkProviderForUserInput{
		UserPoolId: aws.String(plan.UserPoolId.ValueString()),
		DestinationUser: &cognitotypes.ProviderUserIdentifierType{
			ProviderName:           aws.String(plan.DestinationUser.ProviderName.ValueString()),
			ProviderAttributeValue: aws.String(plan.DestinationUser.ProviderAttributeValue.ValueString()),
		},
		SourceUser: &cognitotypes.ProviderUserIdentifierType{
			ProviderName:           aws.String(plan.SourceUser.ProviderName.ValueString()),
			ProviderAttributeName:  aws.String(plan.SourceUser.ProviderAttributeName.ValueString()),
			ProviderAttributeValue: aws.String(plan.SourceUser.ProviderAttributeValue.ValueString()),
		},
	}

	_, err := r.client.AdminLinkProviderForUser(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to link users",
			fmt.Sprintf("Unable to create link, got error: %s", err),
		)
		return
	}

	plan.Id = types.StringValue(fmt.Sprintf("%s:%s:%s",
		plan.UserPoolId.ValueString(),
		plan.DestinationUser.ProviderAttributeValue.ValueString(),
		plan.SourceUser.ProviderAttributeValue.ValueString(),
	))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *CognitoIdpLinkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CognitoIdpLinkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check existence using AdminListGroupsForUser
	// Since there is no API to get actual link information, we consider it exists if no error is returned
	input := &cognitoidentityprovider.AdminListGroupsForUserInput{
		UserPoolId: aws.String(state.UserPoolId.ValueString()),
		Username:   aws.String(state.DestinationUser.ProviderAttributeValue.ValueString()),
	}

	_, err := r.client.AdminListGroupsForUser(ctx, input)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *CognitoIdpLinkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// This resource does not support update operations
	resp.Diagnostics.AddError(
		"Update operation not supported",
		"This resource does not support update operations. Please recreate the resource.",
	)
}

func (r *CognitoIdpLinkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CognitoIdpLinkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &cognitoidentityprovider.AdminDisableProviderForUserInput{
		UserPoolId: aws.String(state.UserPoolId.ValueString()),
		User: &cognitotypes.ProviderUserIdentifierType{
			ProviderName:           aws.String(state.SourceUser.ProviderName.ValueString()),
			ProviderAttributeName:  aws.String(state.SourceUser.ProviderAttributeName.ValueString()),
			ProviderAttributeValue: aws.String(state.SourceUser.ProviderAttributeValue.ValueString()),
		},
	}

	_, err := r.client.AdminDisableProviderForUser(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to unlink users",
			fmt.Sprintf("Unable to unlink user, got error: %s", err),
		)
		return
	}
}

func (r *CognitoIdpLinkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
