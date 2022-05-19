---
name: Bug Report
about: Create a bug report to help us improve this project
title: "[Main title for your issue here]"
labels: bug
assignees: smutel
---
<!--- BEFORE CREATING A NEW ISSUE, ENSURE THAT THE PROVIDER VERSION THAT YOUR ARE USING IS -->
<!--- COMPATIBLE WITH YOUR VERSION OF NETBOX. PLEASE CHECK THE COMPATIBILITY MATRIX HERE -->
<!--- https://github.com/smutel/terraform-provider-netbox#compatibility-with-netbox -->

<!--- Verify first that your issue is not already reported on GitHub -->
<!--- Ensure that the latest release is affected by this bug -->
<!--- Complete most of sections below as described -->

## Summary
<!--- Describe here with one sentence the bug encountered -->

## Version

### Netbox version
<!--- Enter below the version of netbox -->
```paste below

```

### Terraform version
<!--- Enter below the result of "terraform -v" -->
```paste below

```

### Provider version
<!--- Enter below the version of terraform-provider-netbox -->
```paste below

```

## Issue details

### Affected Data(s) / Resource(s)
<!--- Give the name of the data(s) or resource(s) affected by this bug -->
* data dcim_site
* resource ipam_prefix

### Terraform Configuration Files
<!-- Copy-paste your Terraform configurations below -->
<!-- For large Terraform configs, please give a link to a https://gist.github.com -->
```hcl

```

### Terraform Output
<!-- Copy-paste the terraform output (only the error) -->
```paste below

```

## Behaviors

### Actual Behavior
<!-- Describe below the actual behavior -->

### Expected Behavior
<!-- Describe below the expected behavior -->

### Steps to Reproduce
<!-- Please list the steps required to reproduce the issue, for example:-->
1. `terraform apply`

