# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [1.2.0](https://github.com/smutel/terraform-provider-netbox/compare/v1.1.1...v1.2.0) (2021-01-01)


### Features

* Add ability to return json objects of all endpoints ([dceb7f2](https://github.com/smutel/terraform-provider-netbox/commit/dceb7f2c9ccd6b309c65865eba5798d4ca4bedbf))
* Add custom_fields to IP address resource ([bb962ee](https://github.com/smutel/terraform-provider-netbox/commit/bb962ee52127e611900b03f64f458b81a23d52c1))
* Add custom_fields to tenant resource ([01a000f](https://github.com/smutel/terraform-provider-netbox/commit/01a000ffa21018ea74ca3e354e17a087d465719f))
* Add custom_fields to vlan resource ([1315546](https://github.com/smutel/terraform-provider-netbox/commit/13155460a50b2983745ce1a0f448aa6079bd93e4))
* Add custom_fields to vm resource and update documentation ([0fb146d](https://github.com/smutel/terraform-provider-netbox/commit/0fb146df0b40466c63f795dd5f2c105cc55ee79a))
* Add data source IPAM aggregate ([f33b8d9](https://github.com/smutel/terraform-provider-netbox/commit/f33b8d9908b0639918632203fe878ee092fa333d))
* Add resource IPAM aggregate ([cdaa29f](https://github.com/smutel/terraform-provider-netbox/commit/cdaa29f4b9d93a418c47e28ae603288c91d5ce29))

### [1.1.1](https://github.com/smutel/terraform-provider-netbox/compare/v1.1.0...v1.1.1) (2020-11-25)


### Bug Fixes

* IP address not updating when dns_name removed ([c4cd536](https://github.com/smutel/terraform-provider-netbox/commit/c4cd5362ac383abca5e2409256c6a70fbf329a98))

## [1.1.0](https://github.com/smutel/terraform-provider-netbox/compare/v1.0.0...v1.1.0) (2020-11-01)


### Features

* Add Cluster data ([5fafa79](https://github.com/smutel/terraform-provider-netbox/commit/5fafa7941507eec43b12dacab9903b0c61287693))
* Add Platform data ([7eb4c91](https://github.com/smutel/terraform-provider-netbox/commit/7eb4c912eb1398f5c07cca224f46c7afb2a8da29))
* Add VirtualMachine resource ([b7dd75f](https://github.com/smutel/terraform-provider-netbox/commit/b7dd75f9114649cd9f9d4560967707bf95e0d48e))
* Add VirtualMachineInterface resource ([8798e24](https://github.com/smutel/terraform-provider-netbox/commit/8798e24db012c96e5be73dfd98370ed4ace2f7d4))


### Bug Fixes

* Fix bug in go-netbox ([c238629](https://github.com/smutel/terraform-provider-netbox/commit/c2386296908edb996605b415624e5187d30355b3))
* Unable to remove description for some resources ([f340062](https://github.com/smutel/terraform-provider-netbox/commit/f340062ab2c49766373cec803629d172e4be6d04))


### Enhancements

* Add primary ip4 option for ip addresses ([f7e5ab2](https://github.com/smutel/terraform-provider-netbox/commit/f7e5ab2473e8390b752765b896f21a81d2ae4359))
* Destroy ip address when object_id is changing ([6e87ef7](https://github.com/smutel/terraform-provider-netbox/commit/6e87ef7c69131c94d6c359a70f3affd88b44f5e3))
* Set default value for object_type ([edca694](https://github.com/smutel/terraform-provider-netbox/commit/edca6949c9bb7e7088d2d7fed1f5520925ceb135))

## [1.0.0](https://github.com/smutel/terraform-provider-netbox/compare/v0.2.1...v1.0.0) (2020-10-17)


### Features

* Update provider to work with Netbox 2.9 ([8d5633e](https://github.com/smutel/terraform-provider-netbox/commit/8d5633e7748e25e078862859026e7f0aff387ec9))

### [0.2.1](https://github.com/smutel/terraform-provider-netbox/compare/v0.2.0...v0.2.1) (2020-06-28)


### Bug Fixes

* Data netbox_ipam_vlan return code 400 ([55d1e0b](https://github.com/smutel/terraform-provider-netbox/commit/55d1e0b1503899886dca0ea0c9be698efbadc2d4))

## [0.2.0](https://github.com/smutel/terraform-provider-netbox/compare/v0.1.0...v0.2.0) (2020-06-25)


### Bug Fixes

* Fix issue when using http as scheme ([0171927](https://github.com/smutel/terraform-provider-netbox/commit/0171927))


### Enhancements

* go mod tidy on modules ([1291bb7](https://github.com/smutel/terraform-provider-netbox/commit/1291bb7))


### Features

* Add code to manage IP addresses ([9db82c0](https://github.com/smutel/terraform-provider-netbox/commit/9db82c0))

## 0.1.0 (2020-05-25)


### Features

* Init the project ([c9672cd](https://github.com/smutel/terraform-provider-netbox/commit/c9672cd4f3a2eba67f4482caac321ea83eeffac1))
