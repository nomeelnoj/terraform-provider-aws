// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package inspector2

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/inspector2"
	awstypes "github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_inspector_code_scanning_configuration", name="Code Scanning Configuration")
// @Tags(identifierAttribute="arn")
func resourceCodeScanningConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCodeScanningConfigurationCreate,
		ReadWithoutTimeout:   resourceCodeScanningConfigurationRead,
		UpdateWithoutTimeout: resourceCodeScanningConfigurationUpdate,
		DeleteWithoutTimeout: resourceCodeScanningConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			names.AttrARN: {
				Type:     schema.TypeString,
				Computed: true,
			},
			"configuration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"continuous_integration_scan_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"supported_events": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: enum.Validate[awstypes.ContinuousIntegrationScanEvent](),
										},
									},
								},
							},
						},
						"periodic_scan_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frequency": {
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: enum.Validate[awstypes.PeriodicScanFrequency](),
									},
									"frequency_expression": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"rule_set_categories": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: enum.Validate[awstypes.RuleSetCategory](),
							},
						},
					},
				},
			},
			"level": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: enum.Validate[awstypes.ConfigurationLevel](),
			},
			names.AttrName: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"scope_settings": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_select_scope": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: enum.Validate[awstypes.ProjectSelectionScope](),
						},
					},
				},
			},
			names.AttrTags:    tftags.TagsSchema(),
			names.AttrTagsAll: tftags.TagsSchemaComputed(),
		},

		CustomizeDiff: verify.SetTagsDiff,
	}
}

func resourceCodeScanningConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).Inspector2Client(ctx)

	name := d.Get(names.AttrName).(string)
	input := &inspector2.CreateCodeSecurityScanConfigurationInput{
		Configuration: expandCodeSecurityScanConfiguration(d.Get("configuration").([]interface{})),
		Level:         awstypes.ConfigurationLevel(d.Get("level").(string)),
		Name:          aws.String(name),
		Tags:          getTagsIn(ctx),
	}

	if v, ok := d.GetOk("scope_settings"); ok && len(v.([]interface{})) > 0 && v.([]interface{})[0] != nil {
		input.ScopeSettings = expandScopeSettings(v.([]interface{}))
	}

	output, err := conn.CreateCodeSecurityScanConfiguration(ctx, input)

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "creating Inspector Code Scanning Configuration (%s): %s", name, err)
	}

	d.SetId(aws.ToString(output.ScanConfigurationArn))

	return append(diags, resourceCodeScanningConfigurationRead(ctx, d, meta)...)
}

func resourceCodeScanningConfigurationRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).Inspector2Client(ctx)

	output, err := findCodeScanningConfigurationByARN(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] Inspector Code Scanning Configuration (%s) not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Inspector Code Scanning Configuration (%s): %s", d.Id(), err)
	}

	d.Set(names.AttrARN, output.ScanConfigurationArn)
	if err := d.Set("configuration", flattenCodeSecurityScanConfiguration(output.Configuration)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting configuration: %s", err)
	}
	d.Set("level", output.Level)
	d.Set(names.AttrName, output.Name)
	if err := d.Set("scope_settings", flattenScopeSettings(output.ScopeSettings)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting scope_settings: %s", err)
	}

	setTagsOut(ctx, output.Tags)

	return diags
}

func resourceCodeScanningConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).Inspector2Client(ctx)

	if d.HasChangesExcept(names.AttrTags, names.AttrTagsAll) {
		input := &inspector2.UpdateCodeSecurityScanConfigurationInput{
			Configuration:        expandCodeSecurityScanConfiguration(d.Get("configuration").([]interface{})),
			ScanConfigurationArn: aws.String(d.Id()),
		}

		_, err := conn.UpdateCodeSecurityScanConfiguration(ctx, input)

		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating Inspector Code Scanning Configuration (%s): %s", d.Id(), err)
		}
	}

	return append(diags, resourceCodeScanningConfigurationRead(ctx, d, meta)...)
}

func resourceCodeScanningConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).Inspector2Client(ctx)

	log.Printf("[DEBUG] Deleting Inspector Code Scanning Configuration: %s", d.Id())
	_, err := conn.DeleteCodeSecurityScanConfiguration(ctx, &inspector2.DeleteCodeSecurityScanConfigurationInput{
		ScanConfigurationArn: aws.String(d.Id()),
	})

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "deleting Inspector Code Scanning Configuration (%s): %s", d.Id(), err)
	}

	return diags
}

func findCodeScanningConfigurationByARN(ctx context.Context, conn *inspector2.Client, arn string) (*inspector2.GetCodeSecurityScanConfigurationOutput, error) {
	input := &inspector2.GetCodeSecurityScanConfigurationInput{
		ScanConfigurationArn: aws.String(arn),
	}

	output, err := conn.GetCodeSecurityScanConfiguration(ctx, input)

	if errs.IsA[*awstypes.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output, nil
}

func expandCodeSecurityScanConfiguration(tfList []interface{}) *awstypes.CodeSecurityScanConfiguration {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}

	tfMap := tfList[0].(map[string]interface{})
	config := &awstypes.CodeSecurityScanConfiguration{
		RuleSetCategories: flex.ExpandStringyValueSet[awstypes.RuleSetCategory](tfMap["rule_set_categories"].(*schema.Set)),
	}

	if v, ok := tfMap["continuous_integration_scan_configuration"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		config.ContinuousIntegrationScanConfiguration = expandContinuousIntegrationScanConfiguration(v)
	}

	if v, ok := tfMap["periodic_scan_configuration"].([]interface{}); ok && len(v) > 0 && v[0] != nil {
		config.PeriodicScanConfiguration = expandPeriodicScanConfiguration(v)
	}

	return config
}

func expandContinuousIntegrationScanConfiguration(tfList []interface{}) *awstypes.ContinuousIntegrationScanConfiguration {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}

	tfMap := tfList[0].(map[string]interface{})
	return &awstypes.ContinuousIntegrationScanConfiguration{
		SupportedEvents: flex.ExpandStringyValueSet[awstypes.ContinuousIntegrationScanEvent](tfMap["supported_events"].(*schema.Set)),
	}
}

func expandPeriodicScanConfiguration(tfList []interface{}) *awstypes.PeriodicScanConfiguration {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}

	tfMap := tfList[0].(map[string]interface{})
	config := &awstypes.PeriodicScanConfiguration{
		Frequency: awstypes.PeriodicScanFrequency(tfMap["frequency"].(string)),
	}

	if v, ok := tfMap["frequency_expression"].(string); ok && v != "" {
		config.FrequencyExpression = aws.String(v)
	}

	return config
}

func expandScopeSettings(tfList []interface{}) *awstypes.ScopeSettings {
	if len(tfList) == 0 || tfList[0] == nil {
		return nil
	}

	tfMap := tfList[0].(map[string]interface{})
	return &awstypes.ScopeSettings{
		ProjectSelectionScope: awstypes.ProjectSelectionScope(tfMap["project_select_scope"].(string)),
	}
}

func flattenCodeSecurityScanConfiguration(config *awstypes.CodeSecurityScanConfiguration) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	tfMap := map[string]interface{}{
		"rule_set_categories": flex.FlattenStringyValueSet(config.RuleSetCategories),
	}

	if config.ContinuousIntegrationScanConfiguration != nil {
		tfMap["continuous_integration_scan_configuration"] = flattenContinuousIntegrationScanConfiguration(config.ContinuousIntegrationScanConfiguration)
	}

	if config.PeriodicScanConfiguration != nil {
		tfMap["periodic_scan_configuration"] = flattenPeriodicScanConfiguration(config.PeriodicScanConfiguration)
	}

	return []interface{}{tfMap}
}

func flattenContinuousIntegrationScanConfiguration(config *awstypes.ContinuousIntegrationScanConfiguration) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"supported_events": flex.FlattenStringyValueSet(config.SupportedEvents),
		},
	}
}

func flattenPeriodicScanConfiguration(config *awstypes.PeriodicScanConfiguration) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	tfMap := map[string]interface{}{
		"frequency": config.Frequency,
	}

	if config.FrequencyExpression != nil {
		tfMap["frequency_expression"] = aws.ToString(config.FrequencyExpression)
	}

	return []interface{}{tfMap}
}

func flattenScopeSettings(settings *awstypes.ScopeSettings) []interface{} {
	if settings == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"project_select_scope": settings.ProjectSelectionScope,
		},
	}
}
