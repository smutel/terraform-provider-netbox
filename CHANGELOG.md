# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [5.1.0](https://github.com/smutel/terraform-provider-netbox/compare/v5.0.0...v5.1.0) (2022-09-24)


### Features

* Add debug config for VS Code ([e71fd8f](https://github.com/smutel/terraform-provider-netbox/commit/e71fd8fca6414a0c2b31637250ec996df9bd82f4))
* Add debug option to main() ([f52ffe6](https://github.com/smutel/terraform-provider-netbox/commit/f52ffe6e12f250c042279f5bcae835f5b9a23719))
* Add IP Range resource type ([35d31a6](https://github.com/smutel/terraform-provider-netbox/commit/35d31a6f878d66c905777167037102982263a6ed))
* Allow prefix and ip to be dynamically assigned ([88cf8e9](https://github.com/smutel/terraform-provider-netbox/commit/88cf8e92d4002c3ae9ebedf5bb1b8eb2e5c6b0a6))
* Change go-netbox major version ([87797d4](https://github.com/smutel/terraform-provider-netbox/commit/87797d471b3d689a49ef33800bcffcbb35af9b68))


### Bug Fixes

* Get all objects with json data sources ([58c2a11](https://github.com/smutel/terraform-provider-netbox/commit/58c2a11313fb24d7767017535bff0e8fe598c2f2))
* Replace deprecated functions with context aware functions ([f8460bf](https://github.com/smutel/terraform-provider-netbox/commit/f8460bfaf38e1dae5a8efedfee553e6120aa041b))


### Enhancements

* Use context aware functions ([b2b8943](https://github.com/smutel/terraform-provider-netbox/commit/b2b89432e7b51d374afb3d9c2dc3f7af4d016752))

## [5.0.0](https://github.com/smutel/terraform-provider-netbox/compare/v4.0.1...v5.0.0) (2022-08-26)


### Features

* Update provider to work with Netbox 3.2 ([ea07af5](https://github.com/smutel/terraform-provider-netbox/commit/ea07af58369ee844b5fb051690e350f400192bc2))


### Enhancements

* Update JSON data/resource with Netbox 3.2 ([adae68a](https://github.com/smutel/terraform-provider-netbox/commit/adae68aed8622c1b541103a1a523a49595e7d025))

### [4.0.1](https://github.com/smutel/terraform-provider-netbox/compare/v4.0.0...v4.0.1) (2022-06-17)


### Bug Fixes

* JSON error when vlangroup contains scope ([875b1ab](https://github.com/smutel/terraform-provider-netbox/commit/875b1abe8b4c627d51f13c6581ac28862ba1803d))


### Enhancements

* Automate the doc generation ([bc82477](https://github.com/smutel/terraform-provider-netbox/commit/bc82477f18f3c670b1293c26bd0172b1d7ec4233))

## [4.0.0](https://github.com/smutel/terraform-provider-netbox/compare/v3.1.0...v4.0.0) (2022-06-09)


### Features

* Add contact assignment resource ([2b413f2](https://github.com/smutel/terraform-provider-netbox/commit/2b413f2245f800b49f2f40e80fdff1c11e270a2a))
* Add contact group resource ([5c1b8cd](https://github.com/smutel/terraform-provider-netbox/commit/5c1b8cda31c4caaff73a89426c5227ac3e10dd31))
* Add contact resource ([43d3089](https://github.com/smutel/terraform-provider-netbox/commit/43d308959e63ae32aa927ccfbe70e48961ce1d25))
* Add contact role resource ([372fa49](https://github.com/smutel/terraform-provider-netbox/commit/372fa492cc25cabd659e215a17e50f2913ab33e1))
* Add tag parameter to tenant group resource ([83f77f6](https://github.com/smutel/terraform-provider-netbox/commit/83f77f64a33021339b3039a06ce8b1610e95d84f))
* Add tag parameter to vlan group resource ([eb88a02](https://github.com/smutel/terraform-provider-netbox/commit/eb88a0230f6afd749b2702cc9c4016d3120351f9))
* Update provider to work with Netbox 3.1 ([169c160](https://github.com/smutel/terraform-provider-netbox/commit/169c160509733d312ad8097a697b2439c982141f))


### Bug Fixes

* Allow ipv6 addresses for data.netbox_ipam_ip_addresses ([c54cab7](https://github.com/smutel/terraform-provider-netbox/commit/c54cab732e4c611548d8b5819e43549ad74f3ccd)), closes [#107](https://github.com/smutel/terraform-provider-netbox/issues/107)


### Enhancements

* Update JSON data/resource with Netbox 3.1 ([973ab72](https://github.com/smutel/terraform-provider-netbox/commit/973ab72db755fbd4187fedda72ac1055aa1de9cd))

## [3.1.0](https://github.com/smutel/terraform-provider-netbox/compare/v3.0.0...v3.1.0) (2022-04-28)


### Bug Fixes

* Allow setting local_context_data ([5ba20c6](https://github.com/smutel/terraform-provider-netbox/commit/5ba20c6b23fa9d8fa7b5b6a16b8f99155c3aedfa)), closes [#59](https://github.com/smutel/terraform-provider-netbox/issues/59)


### Enhancements

* Remove unused private key (vault removed) ([92804a3](https://github.com/smutel/terraform-provider-netbox/commit/92804a3f40c165503e19d5123784459d5acf1992))
* Switch to terraform SDKv2 ([dde04e1](https://github.com/smutel/terraform-provider-netbox/commit/dde04e1a5b7d5c5882866c50d237550bc521678c))

## [3.0.0](https://github.com/smutel/terraform-provider-netbox/compare/v2.2.1...v3.0.0) (2022-03-18)


### Features

* Update provider to work with Netbox 3.0 ([6327700](https://github.com/smutel/terraform-provider-netbox/commit/6327700872152cf16100b60ae2a8643e4f0e8ae4))

### [2.2.1](https://github.com/smutel/terraform-provider-netbox/compare/v2.2.0...v2.2.1) (2022-01-07)


### Bug Fixes

* Store empty static String field as nil in TF state ([00dd09e](https://github.com/smutel/terraform-provider-netbox/commit/00dd09ed12e7a2d8cd577a102c24400554a92103))

## [2.2.0](https://github.com/smutel/terraform-provider-netbox/compare/v2.1.0...v2.2.0) (2021-11-19)


### Bug Fixes

* Handle empty custom fields correctly ([b3d342e](https://github.com/smutel/terraform-provider-netbox/commit/b3d342edcae376a79e78a13b354fe9f5cc05ed9a))
* Netbox Boolean custom fields take a Boolean type value ([7015047](https://github.com/smutel/terraform-provider-netbox/commit/7015047d012bef81b7a487587ed7522579d56391))


### Enhancements

* Manage custom fields in prefix resource ([e8086bc](https://github.com/smutel/terraform-provider-netbox/commit/e8086bc0d347eca4d8aba3809b74f64735cf70eb))

## [2.1.0](https://github.com/smutel/terraform-provider-netbox/compare/v2.0.2...v2.1.0) (2021-10-26)


### Features

* Add custom fields to interface resource ([dc34ae5](https://github.com/smutel/terraform-provider-netbox/commit/dc34ae5e43c16190795db5c88559deb5cc45f2e4))


### Bug Fixes

* Don't display changes when vcpus does not contain .00 ([085050b](https://github.com/smutel/terraform-provider-netbox/commit/085050bb6ee04f7892410fc37e33fbd0c184ffd4))

### [2.0.2](https://github.com/smutel/terraform-provider-netbox/compare/v2.0.1...v2.0.2) (2021-09-30)


### Bug Fixes

* Terraform crashes when ip address not affected to vm/device ([21d283a](https://github.com/smutel/terraform-provider-netbox/commit/21d283a97d2a963945a665f63200b9dea4a27574))

### [2.0.1](https://github.com/smutel/terraform-provider-netbox/compare/v1.3.0...v2.0.1) (2021-09-24)


### Features

* Add an option to initialize the client with the private key ([d638c23](https://github.com/smutel/terraform-provider-netbox/commit/d638c230e124c77f54ac711bfef9df2621367656))
* Add limit param to data resources ([472b9e1](https://github.com/smutel/terraform-provider-netbox/commit/472b9e125e25f5bfe6d3db0c7ec4f303dddfc1fc))
* Update provider to work with Netbox 2.11 ([b4b2512](https://github.com/smutel/terraform-provider-netbox/commit/b4b2512a25cb41ef563d9e5a2c6c7aa834e00c7d))


### Bug Fixes

* Fix CI for release generation ([0e70fbf](https://github.com/smutel/terraform-provider-netbox/commit/0e70fbf7a9462596954db52acbe22c27b0d96bb3))

## [2.0.0](https://github.com/smutel/terraform-provider-netbox/compare/v1.3.0...v2.0.0) (2021-09-24)


### Features

* Add an option to initialize the client with the private key ([d638c23](https://github.com/smutel/terraform-provider-netbox/commit/d638c230e124c77f54ac711bfef9df2621367656))
* Add limit param to data resources ([472b9e1](https://github.com/smutel/terraform-provider-netbox/commit/472b9e125e25f5bfe6d3db0c7ec4f303dddfc1fc))
* Update provider to work with Netbox 2.11 ([b4b2512](https://github.com/smutel/terraform-provider-netbox/commit/b4b2512a25cb41ef563d9e5a2c6c7aa834e00c7d))

## [1.3.0](https://github.com/smutel/terraform-provider-netbox/compare/v1.2.0...v1.3.0) (2021-04-08)


### Features

* Add data IPAM service ([69ae53c](https://github.com/smutel/terraform-provider-netbox/commit/69ae53caf7403fa3ca7170dedd6223b3944d083f))
* Add default importer function to resources ([22c0538](https://github.com/smutel/terraform-provider-netbox/commit/22c0538a12f49b876fbb729ed5b18731d12b362f)), closes [#51](https://github.com/smutel/terraform-provider-netbox/issues/51)
* Add resource IPAM service ([4e3d637](https://github.com/smutel/terraform-provider-netbox/commit/4e3d6374e368e9260457c639ba8bdb3f8ddb03ca))
* Add TLS insecure option ([827cdeb](https://github.com/smutel/terraform-provider-netbox/commit/827cdeb4e91901ee4a02ab3feb3260b6ba49e08b))
* Allow a configurable basepath ([7392901](https://github.com/smutel/terraform-provider-netbox/commit/7392901e99ba49c2e567de9a8023801edbec2f38))


### Bug Fixes

* Make basepath really optional ([37957ea](https://github.com/smutel/terraform-provider-netbox/commit/37957eab201ee0aa8aabd3ba6725ba52f63cac05))
* Workaround for weird terraform behavior to fix boolean custom_fields ([f9fe60c](https://github.com/smutel/terraform-provider-netbox/commit/f9fe60ce9f421ae3234468ef8533e71c0e13a60a))


### Enhancements

* Change ValidateFunc for ipam_ip_addresses to allow IPv6 ([cb13d6a](https://github.com/smutel/terraform-provider-netbox/commit/cb13d6a632803e190b30e7158a1f4c43f1174201))
* Don't commit .terraform.lock.hcl ([77d74b5](https://github.com/smutel/terraform-provider-netbox/commit/77d74b58f592e76148240a0a8f83d88602520883))


### Tests

* Add script to load netbox-docker variables ([e167ef6](https://github.com/smutel/terraform-provider-netbox/commit/e167ef6f270aa1bf25d1a31310e8e7b1da0590cb))

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
