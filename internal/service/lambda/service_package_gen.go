// Code generated by internal/generate/servicepackages/main.go; DO NOT EDIT.

package lambda

import (
	"context"

	aws_sdkv2 "github.com/aws/aws-sdk-go-v2/aws"
	lambda_sdkv2 "github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  dataSourceAlias,
			TypeName: "aws_lambda_alias",
		},
		{
			Factory:  DataSourceCodeSigningConfig,
			TypeName: "aws_lambda_code_signing_config",
		},
		{
			Factory:  DataSourceFunction,
			TypeName: "aws_lambda_function",
		},
		{
			Factory:  DataSourceFunctionURL,
			TypeName: "aws_lambda_function_url",
		},
		{
			Factory:  DataSourceFunctions,
			TypeName: "aws_lambda_functions",
		},
		{
			Factory:  DataSourceInvocation,
			TypeName: "aws_lambda_invocation",
		},
		{
			Factory:  DataSourceLayerVersion,
			TypeName: "aws_lambda_layer_version",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  resourceAlias,
			TypeName: "aws_lambda_alias",
			Name:     "Alias",
		},
		{
			Factory:  resourceCodeSigningConfig,
			TypeName: "aws_lambda_code_signing_config",
			Name:     "Code Signing Config",
		},
		{
			Factory:  ResourceEventSourceMapping,
			TypeName: "aws_lambda_event_source_mapping",
		},
		{
			Factory:  ResourceFunction,
			TypeName: "aws_lambda_function",
			Name:     "Function",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: "arn",
			},
		},
		{
			Factory:  ResourceFunctionEventInvokeConfig,
			TypeName: "aws_lambda_function_event_invoke_config",
		},
		{
			Factory:  ResourceFunctionURL,
			TypeName: "aws_lambda_function_url",
		},
		{
			Factory:  ResourceInvocation,
			TypeName: "aws_lambda_invocation",
		},
		{
			Factory:  ResourceLayerVersion,
			TypeName: "aws_lambda_layer_version",
		},
		{
			Factory:  ResourceLayerVersionPermission,
			TypeName: "aws_lambda_layer_version_permission",
		},
		{
			Factory:  ResourcePermission,
			TypeName: "aws_lambda_permission",
		},
		{
			Factory:  ResourceProvisionedConcurrencyConfig,
			TypeName: "aws_lambda_provisioned_concurrency_config",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.Lambda
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*lambda_sdkv2.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws_sdkv2.Config))

	return lambda_sdkv2.NewFromConfig(cfg, func(o *lambda_sdkv2.Options) {
		if endpoint := config["endpoint"].(string); endpoint != "" {
			o.BaseEndpoint = aws_sdkv2.String(endpoint)
		}
	}), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
