---
layout: "influxdb-v2"
page_title: "InfluxDB V2: influxdb-v2_bucket"
sidebar_current: "docs-influxdb-v2-resource-bucket"
description: |-
The influxdb-v2_authorization resource manages influxdb v2 buckets.
---

# influxdb-v2_bucket (Resource)



## Example Usage

```terraform
locals {
  org_id = "example_org_id"
}

resource "influxdb-v2_bucket" "example_bucket" {
  name        = "example_bucket_name"
  description = "Example bucket description"
  org_id      = local.org_id
  retention_rules {
    every_seconds = 0
  }
}
```

<!-- schema generated by tfplugindocs -->
## Argument Reference

The following arguments are supported:

* ``name`` (Required) The name of the bucket.
* ``org_id`` (Required) The organization id to which the bucket is linked.
* ``retention_rules`` (Required) Retention rules that affect the bucket.
    * ``every_seconds`` (Required) How many seconds the rule should be applied.
* ``description`` (Optional) The description of the bucket.
* ``rp`` (Optional) As of now, the influxdb documentation doesn't say what this paramenter is for.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* ``created_at`` - The date the bucket has been created.
* ``updated_at`` - The date the bucket has been updated.
* ``type`` - The type of bucket.
* 
## Import

Import is supported using the following syntax:

```shell
terraform import influxdb-v2_bucket.example_bucket <BUCKET_ID>
```
