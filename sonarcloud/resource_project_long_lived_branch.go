package sonarcloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/reinoudk/go-sonarcloud/sonarcloud/settings"
)

const LONGLIVEDBRANCHREGEX = "sonar.branch.longLivedBranch.regex"

type resourceProjectLongLivedBranchType struct{}

func (r resourceProjectLongLivedBranchType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "This resource manage the long lived branch pattern",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"project_key": {
				Type:        types.StringType,
				Required:    true,
				Description: "The key of the project.",
				Validators: []tfsdk.AttributeValidator{
					stringLengthBetween(1, 400),
				},
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
			"value": {
				Type:        types.StringType,
				Required:    true,
				Description: "Value of the long lived branch pattern",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.RequiresReplace(),
				},
			},
		},
	}, nil
}

func (r resourceProjectLongLivedBranchType) Create(ctx context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceProjectLongLivedBranch{
		p: *(p.(*provider)),
	}, nil
}

type resourceProjectLongLivedBranch struct {
	p provider
}

func (r resourceProjectLongLivedBranch) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if r.p.SetDiagErrorIfNotInitialiazed(resp.Diagnostics) {
		return
	}

	var plan ProjectLongLivedBranch
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := settings.SetRequest{
		Component: plan.ProjectKey.Value,
		Key:       LONGLIVEDBRANCHREGEX,
		Value:     plan.Value.Value,
	}

	err := r.p.client.Settings.Set(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not set long lived branch pattern",
			fmt.Sprintf("The Set request returned an error %+v", err),
		)
	}

	result := ProjectLongLivedBranch{
		ID:         types.String{Value: plan.ID.Value},
		ProjectKey: types.String{Value: plan.ProjectKey.Value},
		Value:      types.String{Value: plan.Value.Value},
	}

	diags = resp.State.Set(ctx, result)

	resp.Diagnostics.Append(diags...)
}

// Delete implements tfsdk.Resource
func (resourceProjectLongLivedBranch) Delete(context.Context, tfsdk.DeleteResourceRequest, *tfsdk.DeleteResourceResponse) {
	panic("unimplemented")
}

// Read implements tfsdk.Resource
func (resourceProjectLongLivedBranch) Read(context.Context, tfsdk.ReadResourceRequest, *tfsdk.ReadResourceResponse) {
	panic("unimplemented")
}

// Update implements tfsdk.Resource
func (resourceProjectLongLivedBranch) Update(context.Context, tfsdk.UpdateResourceRequest, *tfsdk.UpdateResourceResponse) {
	panic("unimplemented")
}
