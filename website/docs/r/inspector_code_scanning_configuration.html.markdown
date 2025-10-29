---
subcategory: "Inspector"
layout: "aws"
page_title: "AWS: aws_inspector_code_scanning_configuration"
description: |-
  Terraform resource for managing an AWS Inspector Code Scanning Configuration.
---

# Resource: aws_inspector_code_scanning_configuration

Terraform resource for managing an AWS Inspector Code Scanning Configuration.

## Example Usage

### Basic Usage

```terraform
resource "aws_inspector_code_scanning_configuration" "example" {
  name  = "example"
  level = "ORGANIZATION"

  configuration {
    rule_set_categories = ["SAST"]
  }
}
```

### With Periodic Scan Configuration

```terraform
resource "aws_inspector_code_scanning_configuration" "example" {
  name  = "example"
  level = "ORGANIZATION"

  configuration {
    rule_set_categories = ["SAST", "IAC", "SCA"]

    periodic_scan_configuration {
      frequency            = "WEEKLY"
      frequency_expression = "cron(0 0 * * MON *)"
    }
  }
}
```

### With Continuous Integration Scan Configuration

```terraform
resource "aws_inspector_code_scanning_configuration" "example" {
  name  = "example"
  level = "ORGANIZATION"

  configuration {
    rule_set_categories = ["SAST"]

    continuous_integration_scan_configuration {
      supported_events = ["PULL_REQUEST", "PUSH"]
    }
  }
}
```

### With Scope Settings

```terraform
resource "aws_inspector_code_scanning_configuration" "example" {
  name  = "example"
  level = "ORGANIZATION"

  configuration {
    rule_set_categories = ["SAST"]
  }

  scope_settings {
    project_select_scope = "ALL"
  }

  tags = {
    Name = "example"
  }
}
```

## Argument Reference

This resource supports the following arguments:

* `name` - (Required) The name of the scan configuration. Must be between 1 and 128 characters.
* `level` - (Required) The security level for the scan configuration. Valid values: `ORGANIZATION`.
* `configuration` - (Required) The configuration settings for the code security scan. See [Configuration](#configuration) below.
* `scope_settings` - (Optional) The scope settings that define which repositories will be scanned. See [Scope Settings](#scope-settings) below.
* `tags` - (Optional) A map of tags to assign to the resource. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

### Configuration

* `rule_set_categories` - (Required) The categories of security rules to be applied during the scan. Valid values: `SAST`, `IAC`, `SCA`.
* `continuous_integration_scan_configuration` - (Optional) Configuration settings for continuous integration scans. See [Continuous Integration Scan Configuration](#continuous-integration-scan-configuration) below.
* `periodic_scan_configuration` - (Optional) Configuration settings for periodic scans. See [Periodic Scan Configuration](#periodic-scan-configuration) below.

### Continuous Integration Scan Configuration

* `supported_events` - (Required) The repository events that trigger continuous integration scans. Valid values: `PULL_REQUEST`, `PUSH`.

### Periodic Scan Configuration

* `frequency` - (Required) The frequency at which periodic scans are performed. Valid values: `WEEKLY`, `MONTHLY`, `NEVER`.
* `frequency_expression` - (Optional) The schedule expression for periodic scans, in cron format.

### Scope Settings

* `project_select_scope` - (Required) The scope of projects to be selected for scanning. Setting to `ALL` applies the scope settings to all existing and future projects imported into Amazon Inspector.

## Attribute Reference

This resource exports the following attributes in addition to the arguments above:

* `arn` - ARN of the code scanning configuration.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Inspector Code Scanning Configuration using the `arn`. For example:

```terraform
import {
  to = aws_inspector_code_scanning_configuration.example
  id = "arn:aws:inspector2:us-east-1:123456789012:configuration/scan/abcd1234-5678-90ab-cdef-EXAMPLE11111"
}
```

Using `terraform import`, import Inspector Code Scanning Configuration using the `arn`. For example:

```console
% terraform import aws_inspector_code_scanning_configuration.example arn:aws:inspector2:us-east-1:123456789012:configuration/scan/abcd1234-5678-90ab-cdef-EXAMPLE11111
```
