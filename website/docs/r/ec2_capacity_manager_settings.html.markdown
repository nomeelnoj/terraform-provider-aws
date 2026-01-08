---
subcategory: "EC2 (Elastic Compute Cloud)"
layout: "aws"
page_title: "AWS: aws_ec2_capacity_manager_settings"
description: |-
  Manages EC2 Capacity Manager settings for your AWS account.
---

# Resource: aws_ec2_capacity_manager_settings

Provides a resource to manage EC2 Capacity Manager settings for your AWS account. EC2 Capacity Manager enables customers to monitor, analyze, and manage EC2 capacity across all of their accounts and regions.

~> **NOTE:** Removing this Terraform resource disables EC2 Capacity Manager.

## Example Usage

### Basic Usage

```terraform
resource "aws_ec2_capacity_manager_settings" "example" {
  enabled = true
}
```

### With Organizations Access

```terraform
resource "aws_ec2_capacity_manager_settings" "example" {
  enabled              = true
  organizations_access = true
}
```

### With Delegated Administrator

```terraform
provider "aws" {
  alias  = "org"
  region = "us-west-2"
}

provider "aws" {
  alias  = "delegated"
  region = "us-west-2"
}

data "aws_caller_identity" "delegated" {
  provider = aws.delegated
}

resource "aws_organizations_delegated_administrator" "example" {
  provider          = aws.org
  account_id        = data.aws_caller_identity.delegated.account_id
  service_principal = "ec2.capacitymanager.amazonaws.com"
}

resource "aws_ec2_capacity_manager_settings" "delegated" {
  provider             = aws.delegated
  enabled              = true
  organizations_access = true

  depends_on = [aws_organizations_delegated_administrator.example]
}
```

## Argument Reference

This resource supports the following arguments:

* `enabled` - (Optional) Whether or not EC2 Capacity Manager is enabled. Valid values are `true` or `false`. Defaults to `true`.
* `organizations_access` - (Optional) Whether or not to enable cross-account access for AWS Organizations. When enabled, Capacity Manager can aggregate data from all accounts in your organization. Valid values are `true` or `false`. Defaults to `false`.

## Attribute Reference

This resource exports no additional attributes.

## Import

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import EC2 Capacity Manager settings. For example:

```terraform
import {
  to = aws_ec2_capacity_manager_settings.example
  id = "default"
}
```

Using `terraform import`, import EC2 Capacity Manager settings. For example:

```console
% terraform import aws_ec2_capacity_manager_settings.example default
```
