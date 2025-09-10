---
subcategory: "S3 Control"
layout: "aws"
page_title: "AWS: aws_s3control_storage_lens_configuration"
description: |-
  Provides details about a specific S3 Storage Lens configuration.
---

# Data Source: aws_s3control_storage_lens_configuration

Provides details about a specific S3 Storage Lens configuration.

## Example Usage

```terraform
data "aws_s3control_storage_lens_configuration" "example" {
  name = "Organization"
}
```

## Argument Reference

This data source supports the following arguments:

* `name` - (Required) The name of the S3 Storage Lens configuration.
* `account_id` - (Optional) The AWS account ID of the S3 Storage Lens configuration. Defaults to automatically determined account ID of the Terraform AWS provider.

## Attribute Reference

This data source exports the following attributes in addition to the arguments above:

* `arn` - Amazon Resource Name (ARN) of the S3 Storage Lens configuration.
* `config_id` - The ID of the S3 Storage Lens configuration.
* `storage_lens_configuration` - The S3 Storage Lens configuration. See [storage_lens_configuration](#storage_lens_configuration) below.

### storage_lens_configuration

* `enabled` - Whether the S3 Storage Lens configuration is enabled.
* `account_level` - The account-level configurations of the S3 Storage Lens configuration. See [account_level](#account_level) below.
* `aws_org` - The AWS organization configuration. See [aws_org](#aws_org) below.
* `data_export` - Properties of S3 Storage Lens metrics export including the destination, schema and format. See [data_export](#data_export) below.
* `exclude` - What is excluded in this configuration. Conflicts with `include`. See [exclude](#exclude) below.
* `include` - What is included in this configuration. Conflicts with `exclude`. See [include](#include) below.

### account_level

* `activity_metrics` - S3 Storage Lens activity metrics. See [activity_metrics](#activity_metrics) below.
* `advanced_cost_optimization_metrics` - S3 Storage Lens advanced cost optimization metrics. See [advanced_cost_optimization_metrics](#advanced_cost_optimization_metrics) below.
* `advanced_data_protection_metrics` - S3 Storage Lens advanced data protection metrics. See [advanced_data_protection_metrics](#advanced_data_protection_metrics) below.
* `bucket_level` - S3 Storage Lens bucket-level configuration. See [bucket_level](#bucket_level) below.
* `detailed_status_code_metrics` - S3 Storage Lens detailed status code metrics. See [detailed_status_code_metrics](#detailed_status_code_metrics) below.

### activity_metrics

* `enabled` - Whether the activity metrics are enabled.

### advanced_cost_optimization_metrics

* `enabled` - Whether the advanced cost optimization metrics are enabled.

### advanced_data_protection_metrics

* `enabled` - Whether the advanced data protection metrics are enabled.

### bucket_level

* `activity_metrics` - S3 Storage Lens activity metrics. See [activity_metrics](#activity_metrics) above.
* `advanced_cost_optimization_metrics` - S3 Storage Lens advanced cost optimization metrics. See [advanced_cost_optimization_metrics](#advanced_cost_optimization_metrics) above.
* `advanced_data_protection_metrics` - S3 Storage Lens advanced data protection metrics. See [advanced_data_protection_metrics](#advanced_data_protection_metrics) above.
* `detailed_status_code_metrics` - S3 Storage Lens detailed status code metrics. See [detailed_status_code_metrics](#detailed_status_code_metrics) above.
* `prefix_level` - S3 Storage Lens prefix-level configuration. See [prefix_level](#prefix_level) below.

### prefix_level

* `storage_metrics` - S3 Storage Lens prefix-level storage metrics. See [storage_metrics](#storage_metrics) below.

### storage_metrics

* `enabled` - Whether the prefix-level storage metrics are enabled.
* `selection_criteria` - Selection criteria for prefix-level storage metrics. See [selection_criteria](#selection_criteria) below.

### selection_criteria

* `delimiter` - The delimiter of the selection criteria being used.
* `max_depth` - The max depth of the selection criteria.
* `min_storage_bytes_percentage` - The minimum number of storage bytes percentage whose metrics will be aggregated.

### detailed_status_code_metrics

* `enabled` - Whether the detailed status code metrics are enabled.

### aws_org

* `arn` - The Amazon Resource Name (ARN) of the AWS Organization.

### data_export

* `cloud_watch_metrics` - Amazon CloudWatch publishing for S3 Storage Lens metrics. See [cloud_watch_metrics](#cloud_watch_metrics) below.
* `s3_bucket_destination` - The bucket where the S3 Storage Lens metrics export will be located. See [s3_bucket_destination](#s3_bucket_destination) below.

### cloud_watch_metrics

* `enabled` - Whether CloudWatch publishing for S3 Storage Lens metrics is enabled.

### s3_bucket_destination

* `account_id` - The account ID of the owner of the S3 bucket.
* `arn` - The Amazon Resource Name (ARN) of the bucket.
* `encryption` - Encryption of the metrics exports in this bucket. See [encryption](#encryption) below.
* `format` - The export format.
* `output_schema_version` - The schema version of the export file.
* `prefix` - The prefix of the destination bucket where the metrics export will be delivered.

### encryption

* `sse_kms` - Encryption with Key Management Service (KMS) keys (SSE-KMS). See [sse_kms](#sse_kms) below.
* `sse_s3` - Encryption with Amazon S3 managed keys (SSE-S3). See [sse_s3](#sse_s3) below.

### sse_kms

* `key_id` - KMS key ARN.

### sse_s3

An empty configuration block `{}`.

### exclude

* `buckets` - A list of S3 bucket ARNs.
* `regions` - A list of AWS Regions.

### include

* `buckets` - A list of S3 bucket ARNs.
* `regions` - A list of AWS Regions.