# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [8.0.1](https://github.com/smutel/terraform-provider-netbox/compare/v8.0.0...v8.0.1) (2025-08-09)


### Bug Fixes

* doc generation is broken ([32a8bcf](https://github.com/smutel/terraform-provider-netbox/commit/32a8bcfc05ce2508efd521efcd259016576683fb))

## [8.0.0](https://github.com/smutel/terraform-provider-netbox/compare/v7.1.1...v8.0.0) (2025-06-24)


### Features

* Update to Netbox 4 ([a238b35](https://github.com/smutel/terraform-provider-netbox/commit/a238b35bab18de2782fa268de5e7b6580d41b6a0))

### [7.1.1](https://github.com/smutel/terraform-provider-netbox/compare/v7.1.0...v7.1.1) (2023-10-20)


### Bug Fixes

* Provider is crashing when insecure=true ([40d8df9](https://github.com/smutel/terraform-provider-netbox/commit/40d8df9c8804e5a77f2eb10b18fea0f6a0139014))

## [7.1.0](https://github.com/smutel/terraform-provider-netbox/compare/v7.0.1...v7.1.0) (2023-09-25)


### Features

* Add netbox_dcim_device_role data ([92c401b](https://github.com/smutel/terraform-provider-netbox/commit/92c401b56b6fce60b0beda843c9064b54a1a4aa1))
* Add netbox_dcim_manufacturer data ([5d5c20f](https://github.com/smutel/terraform-provider-netbox/commit/5d5c20f8abf001dbd8c5b88b54b997c383fc22fa))
* Add netbox_ipam_asn data ([f6a1a38](https://github.com/smutel/terraform-provider-netbox/commit/f6a1a38ee3a067236ab48a5c7d4fbb33d221a4cf))
* Add netbox_ipam_ip_range data ([d8f0455](https://github.com/smutel/terraform-provider-netbox/commit/d8f0455e9cc7ec2e1c1db017f3a80aa492152320))
* Add netbox_ipam_prefix data ([42a3222](https://github.com/smutel/terraform-provider-netbox/commit/42a32228563d9130bbd5b2e1a6a9bbd6aad111bc))
* Add netbox_ipam_rir data ([8e83bc2](https://github.com/smutel/terraform-provider-netbox/commit/8e83bc229ed14ffa4f757791032838a6926a985d))
* Add netbox_ipam_route_targets data ([2386105](https://github.com/smutel/terraform-provider-netbox/commit/2386105ae07ef328e7f78970c87cd8447d22d319))
* Add route targets object ([d766780](https://github.com/smutel/terraform-provider-netbox/commit/d766780056ab1beeb2f77c5dc10174877c5465cf))
* Add VRF object ([ab11014](https://github.com/smutel/terraform-provider-netbox/commit/ab11014df641a170b5a4e05662204a7766673ae0))


### Bug Fixes

* Fix some go report issues ([998af89](https://github.com/smutel/terraform-provider-netbox/commit/998af8949a7d056f1b0f8e7c40ef67fe62d8590e))
* Typo in data contact role filename ([008f330](https://github.com/smutel/terraform-provider-netbox/commit/008f330cd5a1bc0eeaa3370a9c715f775c7cc548))


### Enhancements

* Add some missing examples in doc ([0673b06](https://github.com/smutel/terraform-provider-netbox/commit/0673b061c6f8b36262f289697b0d6470137a02b0))
* Remove unused exports csv ([d607ba5](https://github.com/smutel/terraform-provider-netbox/commit/d607ba594c9440bb1127bc20f8645b50bc2990e3))

### [7.0.1](https://github.com/smutel/terraform-provider-netbox/compare/v7.0.0...v7.0.1) (2023-08-11)


### Bug Fixes

* Enable http transport connection reuse ([2fbd144](https://github.com/smutel/terraform-provider-netbox/commit/2fbd14407d0664f4551651a5246f8866f1f9afbc))
* Support FieldName with multiple underscores ([c0787ca](https://github.com/smutel/terraform-provider-netbox/commit/c0787ca39726857098456c071e5786d51611a076))

## [7.0.0](https://github.com/smutel/terraform-provider-netbox/compare/v6.3.0...v7.0.0) (2023-07-03)


### Features

* Improve custom field tests to contain usage test ([2eb5668](https://github.com/smutel/terraform-provider-netbox/commit/2eb56682b20a478ea12eb749906e05d2ca9cdfae))
* Update to work with Netbox 3.4 ([8dbffbf](https://github.com/smutel/terraform-provider-netbox/commit/8dbffbf34c5e62edb77401a6dbd8a5be4305d718))


### Bug Fixes

* Fix choices field for select and multiselect custom field ([e7c882c](https://github.com/smutel/terraform-provider-netbox/commit/e7c882cd53f8f5a63e06a6b232733b156379b27d))
* Value and default for object CF are not lists ([d3b2316](https://github.com/smutel/terraform-provider-netbox/commit/d3b2316cf7ec63e78477aa20c1dd304dde23a496))
* Wrong conversion of filter name with underscore ([b72f202](https://github.com/smutel/terraform-provider-netbox/commit/b72f202fccf0c3dda55772c748b4092f2ce88ca4))

## [6.3.0](https://github.com/smutel/terraform-provider-netbox/compare/v6.2.0...v6.3.0) (2023-04-22)


### Features

* Add netbox_dcim_location data ([e32c18a](https://github.com/smutel/terraform-provider-netbox/commit/e32c18ab82325ea7a4b6e70c1bd5685898c62aee))
* Add netbox_dcim_location resource ([ea33489](https://github.com/smutel/terraform-provider-netbox/commit/ea33489031d91addef212037197dedc479c8c8da))
* Add netbox_dcim_rack data ([6bf7489](https://github.com/smutel/terraform-provider-netbox/commit/6bf748906adde45d7d7c4af5c9d8e0126877d01b))
* Add netbox_dcim_rack resource ([5c8fe86](https://github.com/smutel/terraform-provider-netbox/commit/5c8fe86dbf2a56dcf7467fa9f280cfc9c4378855))
* Add netbox_dcim_rack_role data ([a9932db](https://github.com/smutel/terraform-provider-netbox/commit/a9932dbd2ea7bd4cda38b56fc84b352184610742))
* Add netbox_dcim_rack_role resource ([a0c4808](https://github.com/smutel/terraform-provider-netbox/commit/a0c4808d057db578e80677075c49b48b3c2a3ca7))
* Add netbox_dcim_region data ([80ff164](https://github.com/smutel/terraform-provider-netbox/commit/80ff16494f50179bbb82dec689abc2598c6d51eb))
* Add netbox_dcim_region resource ([6f6970c](https://github.com/smutel/terraform-provider-netbox/commit/6f6970ca377e29ebec3d6d2a144a87399aab5da5))


### Enhancements

* Improve test on resource netbox_dcim_site ([8077515](https://github.com/smutel/terraform-provider-netbox/commit/80775156939cd9984d57a8c50af23f00a80d91b5))
* Replace examples by tests ([3523e3f](https://github.com/smutel/terraform-provider-netbox/commit/3523e3f45e076c191b73eb2bd3242f5747c98b9b))

## [6.2.0](https://github.com/smutel/terraform-provider-netbox/compare/v6.1.1...v6.2.0) (2023-03-14)


### Features

* Add fields and remove forcenew from netbox_virtualization_interface ([89f477c](https://github.com/smutel/terraform-provider-netbox/commit/89f477cbb2eb187dd222e355175c4cade89f86a5))
* Add scope field to VLAN group resource ([67e7628](https://github.com/smutel/terraform-provider-netbox/commit/67e762893a894efe9566959c7c386b462e53b589))
* Add site_id to cluster data source and to vm resource ([a93a3c9](https://github.com/smutel/terraform-provider-netbox/commit/a93a3c92ba51f69aed8a544551491ab62db9ea54))
* Update ipam_ip_addresses ([df80a60](https://github.com/smutel/terraform-provider-netbox/commit/df80a603255703040d392a2be19ffdaea158fa4f)), closes [#196](https://github.com/smutel/terraform-provider-netbox/issues/196)
* Update VM resource for Netbox 3.3 ([2144830](https://github.com/smutel/terraform-provider-netbox/commit/214483026966252a292112399731b3b82e673480)), closes [#192](https://github.com/smutel/terraform-provider-netbox/issues/192)


### Bug Fixes

* Custom fields default values not null ([53382b9](https://github.com/smutel/terraform-provider-netbox/commit/53382b932bf4c3df9aadd7f2c2af8f11fdad3205))
* Export files not working with Netbox 3.3 ([6cd4162](https://github.com/smutel/terraform-provider-netbox/commit/6cd4162a586489aab650f9b45a684e1e48e60374)), closes [#187](https://github.com/smutel/terraform-provider-netbox/issues/187)
* Import does not work for virtualization_vm_primary_ip ([c17026d](https://github.com/smutel/terraform-provider-netbox/commit/c17026d06be926075dc833876606a33417926e13)), closes [#188](https://github.com/smutel/terraform-provider-netbox/issues/188)
* Removing timezone from dcim_site fails ([b2e2106](https://github.com/smutel/terraform-provider-netbox/commit/b2e2106b9e96cdf2049c3754434e7bbafa2127d3))
* Status field missing from cluster resource ([74f5147](https://github.com/smutel/terraform-provider-netbox/commit/74f5147d1e921192feab30dd9de2301e0167c558)), closes [#190](https://github.com/smutel/terraform-provider-netbox/issues/190)
* Trim spaces for vm comment field ([18ec0e8](https://github.com/smutel/terraform-provider-netbox/commit/18ec0e8e8861ea4a639f2a9dfc47b95b7fd4a550))
* Typo for tags field in ipam_rir resource ([dca1f54](https://github.com/smutel/terraform-provider-netbox/commit/dca1f546ff1c064caab1ff96fa6f4d7dbd06c4a3))

### [6.1.1](https://github.com/smutel/terraform-provider-netbox/compare/v6.1.0...v6.1.1) (2023-03-07)


### Bug Fixes

* Increase size of tenant name/desc/slug ([56a4e84](https://github.com/smutel/terraform-provider-netbox/commit/56a4e84925a7e4f2f99a5aef912fe80aaf51551a))

## [6.1.0](https://github.com/smutel/terraform-provider-netbox/compare/v6.0.0...v6.1.0) (2023-02-01)


### Features

* Add filter to JSON data sources ([345f1ad](https://github.com/smutel/terraform-provider-netbox/commit/345f1ad09b067fde167e03cd57c03f3f24f5e2d1))

## [6.0.0](https://github.com/smutel/terraform-provider-netbox/compare/v5.3.1...v6.0.0) (2023-01-13)


### Features

* Update to work with Netbox 3.3 ([bbd12f1](https://github.com/smutel/terraform-provider-netbox/commit/bbd12f1f3d029d484f9b3c37bcb409d115ac2a9c))


### Enhancements

* Update go dependencies ([5f8a611](https://github.com/smutel/terraform-provider-netbox/commit/5f8a61151f13ed55351d3bb2d2d2c3531174bfc5))
* Update version to v6 in go files ([ff5eaa5](https://github.com/smutel/terraform-provider-netbox/commit/ff5eaa563ff4dbb6f6441b3577404d3a92abde46))

### [5.3.1](https://github.com/smutel/terraform-provider-netbox/compare/v5.3.0...v5.3.1) (2022-11-28)


### Bug Fixes

* Workaround in ASN resource for 386 arch ([51e547a](https://github.com/smutel/terraform-provider-netbox/commit/51e547abc23f6cdeb2e73364344220f038d5aa95))

## [5.3.0](https://github.com/smutel/terraform-provider-netbox/compare/v5.2.0...v5.3.0) (2022-11-23)


### Features

* Add ipam rir resource ([2aa2edd](https://github.com/smutel/terraform-provider-netbox/commit/2aa2edd8e84d8cdc3969dfadf8661780fbd71adc))
* Add netbox_dcim_device_role resource ([b250d1e](https://github.com/smutel/terraform-provider-netbox/commit/b250d1e884dcfb54ac03f71b0519ccbde7b424b5))
* Add netbox_dcim_manufacturer ([f214f6a](https://github.com/smutel/terraform-provider-netbox/commit/f214f6aa231ccb77aa0102b760cacc4e0596c6b6))
* Add netbox_dcim_platform resource ([31e5e14](https://github.com/smutel/terraform-provider-netbox/commit/31e5e14f974a5c5b8ebfeb7749fa5c6cfda8473a))
* Add netbox_dcim_site resource ([820a3f1](https://github.com/smutel/terraform-provider-netbox/commit/820a3f130a97cb7d669ccab27e7d758c5f8821f5))
* Add netbox_extras_custom_field resource ([87fbc92](https://github.com/smutel/terraform-provider-netbox/commit/87fbc924c2d31dcfa3f37f63b6ad22adad95c7d2))
* Add netbox_extras_tag resource ([08a4f8a](https://github.com/smutel/terraform-provider-netbox/commit/08a4f8a3ea0115454870d4dde04773701817ecb6))
* Add netbox_ipam_asn resource ([4cc2b9b](https://github.com/smutel/terraform-provider-netbox/commit/4cc2b9bacc75791229836e7a78c2590dd40589d8))
* Add netbox_virtualization_cluster resource ([a292478](https://github.com/smutel/terraform-provider-netbox/commit/a292478f10eefe9534e58a54fea6d09c69e001e8))
* Add netbox_virtualization_cluster_group resource ([c2782ce](https://github.com/smutel/terraform-provider-netbox/commit/c2782ce8868674b0ddd7b2af1e74f17789ed9d7d))
* Add netbox_virtualization_cluster_type resource ([6ac1a1e](https://github.com/smutel/terraform-provider-netbox/commit/6ac1a1e0f3bf68a35525f517e267bd1bd35aff55))
* Add netbox_virtualization_vm_primary_ip resource ([bc2e88a](https://github.com/smutel/terraform-provider-netbox/commit/bc2e88a48b160d7306cb4e89af8f86e62120d794))


### Bug Fixes

* Fix linter deprecation warnings ([47c6ddd](https://github.com/smutel/terraform-provider-netbox/commit/47c6ddd12bfdbebade9a4bfe40d161ac1ea5fcfb))
* Reduce complexity of resourceNetboxVirtualizationVMRead ([c3456ba](https://github.com/smutel/terraform-provider-netbox/commit/c3456bae3bb15e8435515d74d8faae9b18a7e181))
* Some test case are not working correctly ([80a1b00](https://github.com/smutel/terraform-provider-netbox/commit/80a1b00f22e67e3935b5ae5e434e15af1c3faa25))

## [5.2.0](https://github.com/smutel/terraform-provider-netbox/compare/v5.1.0...v5.2.0) (2022-11-10)


### Features

* Add first tests to provider ([2d9cf3b](https://github.com/smutel/terraform-provider-netbox/commit/2d9cf3bbd49d344238be264588597a1b7eaf6a58))
* Add function to modify API requests ([9ed4488](https://github.com/smutel/terraform-provider-netbox/commit/9ed4488c5b3c0287c521d304f290ce30ef525ffc))
* Add suppport for missing customfield types ([c770ce9](https://github.com/smutel/terraform-provider-netbox/commit/c770ce9d5a0a1918e52189bdd6a99ff68a60337e))


### Bug Fixes

* Prevent panic if vcpus is null ([82ad49f](https://github.com/smutel/terraform-provider-netbox/commit/82ad49fa2855885576909d227f12f8885b61eabf))

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
