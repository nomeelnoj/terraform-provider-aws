// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package ec2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/retry"
	tfec2 "github.com/hashicorp/terraform-provider-aws/internal/service/ec2"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccEC2CapacityManagerSettings_serial(t *testing.T) {
	t.Parallel()

	testCases := map[string]map[string]func(t *testing.T){
		"Resource": {
			acctest.CtBasic:                testAccEC2CapacityManagerSettings_basic,
			"organizationsAccess":          testAccEC2CapacityManagerSettings_organizationsAccess,
			"disappears":                   testAccEC2CapacityManagerSettings_disappears,
			"updateEnabled":                testAccEC2CapacityManagerSettings_updateEnabled,
			"updateOrganizationsAccess":    testAccEC2CapacityManagerSettings_updateOrganizationsAccess,
		},
	}

	acctest.RunSerialTests2Levels(t, testCases, 0)
}

func testAccEC2CapacityManagerSettings_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_ec2_capacity_manager_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCapacityManagerSettingsDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCapacityManagerSettingsConfig_basic(true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCapacityManagerSettings(ctx, resourceName, true, false),
					resource.TestCheckResourceAttr(resourceName, names.AttrEnabled, acctest.CtTrue),
					resource.TestCheckResourceAttr(resourceName, "organizations_access", acctest.CtFalse),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEC2CapacityManagerSettings_organizationsAccess(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_ec2_capacity_manager_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationManagementAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCapacityManagerSettingsDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCapacityManagerSettingsConfig_basic(true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCapacityManagerSettings(ctx, resourceName, true, true),
					resource.TestCheckResourceAttr(resourceName, names.AttrEnabled, acctest.CtTrue),
					resource.TestCheckResourceAttr(resourceName, "organizations_access", acctest.CtTrue),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEC2CapacityManagerSettings_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_ec2_capacity_manager_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCapacityManagerSettingsDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCapacityManagerSettingsConfig_basic(true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCapacityManagerSettings(ctx, resourceName, true, false),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfec2.ResourceCapacityManagerSettings, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccEC2CapacityManagerSettings_updateEnabled(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_ec2_capacity_manager_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCapacityManagerSettingsDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCapacityManagerSettingsConfig_basic(false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCapacityManagerSettings(ctx, resourceName, false, false),
					resource.TestCheckResourceAttr(resourceName, names.AttrEnabled, acctest.CtFalse),
					resource.TestCheckResourceAttr(resourceName, "organizations_access", acctest.CtFalse),
				),
			},
			{
				Config: testAccCapacityManagerSettingsConfig_basic(true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCapacityManagerSettings(ctx, resourceName, true, false),
					resource.TestCheckResourceAttr(resourceName, names.AttrEnabled, acctest.CtTrue),
					resource.TestCheckResourceAttr(resourceName, "organizations_access", acctest.CtFalse),
				),
			},
			{
				Config: testAccCapacityManagerSettingsConfig_basic(false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCapacityManagerSettings(ctx, resourceName, false, false),
					resource.TestCheckResourceAttr(resourceName, names.AttrEnabled, acctest.CtFalse),
					resource.TestCheckResourceAttr(resourceName, "organizations_access", acctest.CtFalse),
				),
			},
		},
	})
}

func testAccEC2CapacityManagerSettings_updateOrganizationsAccess(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_ec2_capacity_manager_settings.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationManagementAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.EC2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCapacityManagerSettingsDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCapacityManagerSettingsConfig_basic(true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCapacityManagerSettings(ctx, resourceName, true, false),
					resource.TestCheckResourceAttr(resourceName, names.AttrEnabled, acctest.CtTrue),
					resource.TestCheckResourceAttr(resourceName, "organizations_access", acctest.CtFalse),
				),
			},
			{
				Config: testAccCapacityManagerSettingsConfig_basic(true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCapacityManagerSettings(ctx, resourceName, true, true),
					resource.TestCheckResourceAttr(resourceName, names.AttrEnabled, acctest.CtTrue),
					resource.TestCheckResourceAttr(resourceName, "organizations_access", acctest.CtTrue),
				),
			},
			{
				Config: testAccCapacityManagerSettingsConfig_basic(true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCapacityManagerSettings(ctx, resourceName, true, false),
					resource.TestCheckResourceAttr(resourceName, names.AttrEnabled, acctest.CtTrue),
					resource.TestCheckResourceAttr(resourceName, "organizations_access", acctest.CtFalse),
				),
			},
		},
	})
}

func testAccCheckCapacityManagerSettingsDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Client(ctx)

		// Try to get the attributes - if disabled, this should return NotFound
		_, err := tfec2.FindCapacityManagerAttributes(ctx, conn)

		// NotFound error is expected when disabled (destroy succeeded)
		if retry.NotFound(err) {
			return nil
		}

		if err != nil {
			return err
		}

		// If we got here, the resource is still enabled which means destroy failed
		return fmt.Errorf("EC2 Capacity Manager not disabled on resource removal")
	}
}

func testAccCheckCapacityManagerSettings(ctx context.Context, n string, enabled bool, organizationsAccess bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).EC2Client(ctx)

		output, err := tfec2.FindCapacityManagerAttributes(ctx, conn)

		// If disabled, the finder returns NotFound error
		if !enabled && retry.NotFound(err) {
			return nil
		}

		if err != nil {
			return err
		}

		expectedStatus := types.CapacityManagerStatusEnabled
		if output.CapacityManagerStatus != expectedStatus {
			return fmt.Errorf("EC2 Capacity Manager status = %s, expected %s", output.CapacityManagerStatus, expectedStatus)
		}

		if output.OrganizationsAccess == nil {
			return fmt.Errorf("EC2 Capacity Manager OrganizationsAccess is nil")
		}

		if *output.OrganizationsAccess != organizationsAccess {
			return fmt.Errorf("EC2 Capacity Manager OrganizationsAccess = %t, expected %t", *output.OrganizationsAccess, organizationsAccess)
		}

		return nil
	}
}

func testAccCapacityManagerSettingsConfig_basic(enabled, organizationsAccess bool) string {
	return fmt.Sprintf(`
resource "aws_ec2_capacity_manager_settings" "test" {
  enabled              = %[1]t
  organizations_access = %[2]t
}
`, enabled, organizationsAccess)
}
