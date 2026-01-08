// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package ec2

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	awstypes "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/retry"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKResource("aws_ec2_capacity_manager_settings", name="Capacity Manager Settings")
// @SingletonIdentity
// @Testing(hasExistsFunction=false)
// @Testing(generator=false)
// @Testing(identityTest=false)
func resourceCapacityManagerSettings() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCapacityManagerSettingsCreate,
		ReadWithoutTimeout:   resourceCapacityManagerSettingsRead,
		UpdateWithoutTimeout: resourceCapacityManagerSettingsUpdate,
		DeleteWithoutTimeout: resourceCapacityManagerSettingsDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			names.AttrEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"organizations_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceCapacityManagerSettingsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Client(ctx)

	enabled := d.Get(names.AttrEnabled).(bool)
	organizationsAccess := d.Get("organizations_access").(bool)

	if enabled {
		input := &ec2.EnableCapacityManagerInput{
			OrganizationsAccess: aws.Bool(organizationsAccess),
		}

		_, err := conn.EnableCapacityManager(ctx, input)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "enabling EC2 Capacity Manager: %s", err)
		}
	}

	d.SetId(meta.(*conns.AWSClient).Region(ctx))

	return append(diags, resourceCapacityManagerSettingsRead(ctx, d, meta)...)
}

func resourceCapacityManagerSettingsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Client(ctx)

	output, err := findCapacityManagerAttributes(ctx, conn)

	if !d.IsNewResource() && retry.NotFound(err) {
		log.Printf("[WARN] EC2 Capacity Manager Settings %s not found, removing from state", d.Id())
		d.SetId("")
		return diags
	}

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading EC2 Capacity Manager Settings: %s", err)
	}

	enabled := output.CapacityManagerStatus == awstypes.CapacityManagerStatusEnabled
	d.Set(names.AttrEnabled, enabled)
	d.Set("organizations_access", output.OrganizationsAccess)

	return diags
}

func resourceCapacityManagerSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Client(ctx)

	if d.HasChange(names.AttrEnabled) {
		enabled := d.Get(names.AttrEnabled).(bool)
		organizationsAccess := d.Get("organizations_access").(bool)

		if enabled {
			input := &ec2.EnableCapacityManagerInput{
				OrganizationsAccess: aws.Bool(organizationsAccess),
			}

			_, err := conn.EnableCapacityManager(ctx, input)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "enabling EC2 Capacity Manager: %s", err)
			}
		} else {
			input := &ec2.DisableCapacityManagerInput{}

			_, err := conn.DisableCapacityManager(ctx, input)
			if err != nil {
				return sdkdiag.AppendErrorf(diags, "disabling EC2 Capacity Manager: %s", err)
			}
		}
	} else if d.HasChange("organizations_access") {
		organizationsAccess := d.Get("organizations_access").(bool)
		input := &ec2.UpdateCapacityManagerOrganizationsAccessInput{
			OrganizationsAccess: aws.Bool(organizationsAccess),
		}

		_, err := conn.UpdateCapacityManagerOrganizationsAccess(ctx, input)
		if err != nil {
			return sdkdiag.AppendErrorf(diags, "updating EC2 Capacity Manager Organizations Access: %s", err)
		}
	}

	return append(diags, resourceCapacityManagerSettingsRead(ctx, d, meta)...)
}

func resourceCapacityManagerSettingsDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Client(ctx)

	// Removing the resource disables Capacity Manager.
	input := &ec2.DisableCapacityManagerInput{}

	_, err := conn.DisableCapacityManager(ctx, input)
	if err != nil {
		return sdkdiag.AppendErrorf(diags, "disabling EC2 Capacity Manager: %s", err)
	}

	return diags
}
