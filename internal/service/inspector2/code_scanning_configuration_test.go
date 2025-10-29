// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package inspector2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/inspector2/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfinspector2 "github.com/hashicorp/terraform-provider-aws/internal/service/inspector2"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func testAccCodeScanningConfiguration_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_inspector_code_scanning_configuration.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Inspector2EndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Inspector2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCodeScanningConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCodeScanningConfigurationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodeScanningConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, rName),
					resource.TestCheckResourceAttr(resourceName, "level", string(types.ConfigurationLevelOrganization)),
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.rule_set_categories.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "configuration.0.rule_set_categories.*", "SAST"),
					resource.TestCheckResourceAttrSet(resourceName, names.AttrARN),
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

func testAccCodeScanningConfiguration_periodicScan(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_inspector_code_scanning_configuration.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Inspector2EndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Inspector2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCodeScanningConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCodeScanningConfigurationConfig_periodicScan(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodeScanningConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.periodic_scan_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.periodic_scan_configuration.0.frequency", "WEEKLY"),
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

func testAccCodeScanningConfiguration_continuousIntegrationScan(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_inspector_code_scanning_configuration.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Inspector2EndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Inspector2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCodeScanningConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCodeScanningConfigurationConfig_continuousIntegrationScan(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodeScanningConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.continuous_integration_scan_configuration.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.continuous_integration_scan_configuration.0.supported_events.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "configuration.0.continuous_integration_scan_configuration.0.supported_events.*", "PULL_REQUEST"),
					resource.TestCheckTypeSetElemAttr(resourceName, "configuration.0.continuous_integration_scan_configuration.0.supported_events.*", "PUSH"),
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

func testAccCodeScanningConfiguration_scopeSettings(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_inspector_code_scanning_configuration.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Inspector2EndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Inspector2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCodeScanningConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCodeScanningConfigurationConfig_scopeSettings(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodeScanningConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "scope_settings.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "scope_settings.0.project_select_scope", "ALL"),
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

func testAccCodeScanningConfiguration_tags(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_inspector_code_scanning_configuration.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Inspector2EndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Inspector2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCodeScanningConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCodeScanningConfigurationConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodeScanningConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCodeScanningConfigurationConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodeScanningConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccCodeScanningConfigurationConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodeScanningConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCodeScanningConfiguration_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_inspector_code_scanning_configuration.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.Inspector2EndpointID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.Inspector2ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCodeScanningConfigurationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCodeScanningConfigurationConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCodeScanningConfigurationExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfinspector2.ResourceCodeScanningConfiguration(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckCodeScanningConfigurationExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).Inspector2Client(ctx)

		_, err := tfinspector2.FindCodeScanningConfigurationByARN(ctx, conn, rs.Primary.ID)

		return err
	}
}

func testAccCheckCodeScanningConfigurationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).Inspector2Client(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_inspector_code_scanning_configuration" {
				continue
			}

			_, err := tfinspector2.FindCodeScanningConfigurationByARN(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Inspector Code Scanning Configuration %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCodeScanningConfigurationConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_inspector_code_scanning_configuration" "test" {
  name  = %[1]q
  level = "ORGANIZATION"

  configuration {
    rule_set_categories = ["SAST"]
  }
}
`, rName)
}

func testAccCodeScanningConfigurationConfig_periodicScan(rName string) string {
	return fmt.Sprintf(`
resource "aws_inspector_code_scanning_configuration" "test" {
  name  = %[1]q
  level = "ORGANIZATION"

  configuration {
    rule_set_categories = ["SAST"]

    periodic_scan_configuration {
      frequency = "WEEKLY"
    }
  }
}
`, rName)
}

func testAccCodeScanningConfigurationConfig_continuousIntegrationScan(rName string) string {
	return fmt.Sprintf(`
resource "aws_inspector_code_scanning_configuration" "test" {
  name  = %[1]q
  level = "ORGANIZATION"

  configuration {
    rule_set_categories = ["SAST", "IAC"]

    continuous_integration_scan_configuration {
      supported_events = ["PULL_REQUEST", "PUSH"]
    }
  }
}
`, rName)
}

func testAccCodeScanningConfigurationConfig_scopeSettings(rName string) string {
	return fmt.Sprintf(`
resource "aws_inspector_code_scanning_configuration" "test" {
  name  = %[1]q
  level = "ORGANIZATION"

  configuration {
    rule_set_categories = ["SAST"]
  }

  scope_settings {
    project_select_scope = "ALL"
  }
}
`, rName)
}

func testAccCodeScanningConfigurationConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_inspector_code_scanning_configuration" "test" {
  name  = %[1]q
  level = "ORGANIZATION"

  configuration {
    rule_set_categories = ["SAST"]
  }

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccCodeScanningConfigurationConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_inspector_code_scanning_configuration" "test" {
  name  = %[1]q
  level = "ORGANIZATION"

  configuration {
    rule_set_categories = ["SAST"]
  }

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
