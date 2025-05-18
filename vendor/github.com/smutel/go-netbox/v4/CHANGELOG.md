# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

### [4.0.1](https://github.com/smutel/go-netbox/compare/v4.0.0...v4.0.1) (2025-05-21)


### Features

* Execute goimports on generated library ([72bf4ac](https://github.com/smutel/go-netbox/commit/72bf4ac660f38ce97f3497423976b47afe537fae))

## [4.0.0](https://github.com/smutel/go-netbox/compare/v3.4.0...v4.0.0) (2025-05-20)


### Features

* Switch to openapi ([de6714a](https://github.com/smutel/go-netbox/commit/de6714a303ba55bda9e48a75d406297ca882be30))

## [3.4.0](https://github.com/smutel/go-netbox/compare/v3.3.0...v3.4.0) (2023-06-12)


### Features

* Update from Netbox 3.4 swagger ([ec98310](https://github.com/smutel/go-netbox/commit/ec98310e2ebcaaf3b43ac6c7f14de8ea112b37e6))

## [3.3.0](https://github.com/smutel/go-netbox/compare/v3.2.3...v3.3.0) (2023-01-09)


### Features

* Update from Netbox 3.3 swagger ([6dba05d](https://github.com/smutel/go-netbox/commit/6dba05dfda4a018f4f2d143c031dc05a31731db6))

### [3.2.3](https://github.com/smutel/go-netbox/compare/v3.2.2...v3.2.3) (2022-10-22)


### Features

* Remove patch for primary IP ([08c16b7](https://github.com/smutel/go-netbox/commit/08c16b7bb58c032bfafba310c1daa999094e07e4))

### [3.2.2](https://github.com/smutel/go-netbox/compare/v3.2.1...v3.2.2) (2022-09-23)


### Bug Fixes

* Library version missing in go.mod ([9d5df44](https://github.com/smutel/go-netbox/commit/9d5df442ac55d1cba69b9b2f2783fbbb8ad6556d))

### [3.2.1](https://github.com/smutel/go-netbox/compare/v3.2.0...v3.2.1) (2022-09-14)


### Features

* Add go.mod and go.sum ([58e500f](https://github.com/smutel/go-netbox/commit/58e500f651fcc42b9260f2c00ec8b2e6c8507366))
* Add patch to remove primary IP ([e02f45b](https://github.com/smutel/go-netbox/commit/e02f45b1aa2bc88e0a84f3ef82e1cf6ebb4be037))
* Allow Create functions for available prefix/ip to be used ([ef3c365](https://github.com/smutel/go-netbox/commit/ef3c36530a1cae76f6cdebd1e60d08bfc69d0c40))
* Call go mod tidy after build ([98a4273](https://github.com/smutel/go-netbox/commit/98a42730ae22a1199dd620265dacfdd051552b30))

## [3.2.0](https://github.com/smutel/go-netbox/compare/v3.1.2...v3.2.0) (2022-08-26)


### Features

* Update from Netbox 3.2 swagger ([74f00f1](https://github.com/smutel/go-netbox/commit/74f00f11c2bc36bf333f977f8074c0121a461560))

### [3.1.2](https://github.com/smutel/go-netbox/compare/v3.1.1...v3.1.2) (2022-06-16)


### Bug Fixes

* Scope as object in vlangroup ([4b65d2a](https://github.com/smutel/go-netbox/commit/4b65d2ad4ef61d51236fb491c33af9fab5a351e2))

### [3.1.1](https://github.com/smutel/go-netbox/compare/v3.1.0...v3.1.1) (2022-06-07)


### Bug Fixes

* Object as JSON objects in contact assignment ([9d7758d](https://github.com/smutel/go-netbox/commit/9d7758d392cb519d67b4aee3058f71d321f6a243))

## [3.1.0](https://github.com/smutel/go-netbox/compare/v3.0.1...v3.1.0) (2022-05-17)


### Features

* Update from Netbox 3.1 swagger ([8882546](https://github.com/smutel/go-netbox/commit/888254651a55451de78cf7c91a1625f99460e831))

### [3.0.1](https://github.com/smutel/go-netbox/compare/v3.0.0...v3.0.1) (2022-04-24)


### Bug Fixes

* local_context_data and config_context as JSON objects ([9d3b627](https://github.com/smutel/go-netbox/commit/9d3b62739a002df6ce6a69c68baebc30023f1be2)), closes [#25](https://github.com/smutel/go-netbox/issues/25)

## [3.0.0](https://github.com/smutel/go-netbox/compare/v2.11.0...v3.0.0) (2022-03-11)


### Features

* Add patch to vcpus required ([41d2380](https://github.com/smutel/go-netbox/commit/41d2380bdfbe7a7771562663f3abfdec0f7b0018))
* Update from Netbox 3.0 swagger ([7744a95](https://github.com/smutel/go-netbox/commit/7744a95bfa8136640b2f4c86531fc43f0370f754))
* Update go-swagger to v0.28.0 ([ea5160b](https://github.com/smutel/go-netbox/commit/ea5160bfbbc26cd780d17f643e2dbaa62ab02f87))


### Bug Fixes

* Rollback on vcpus patch ([8ac0c3c](https://github.com/smutel/go-netbox/commit/8ac0c3c91ff2e66ea159609c883b540365b8406c))

## 2.11.0 (2021-09-23)


### Features

* Add first swagger patch ([b5f1123](https://github.com/smutel/go-netbox/commit/b5f1123f77aa2b1a64539a3a8e8c84723ac70023))
* Add patch to ipaddress assignedobject ([e1b5746](https://github.com/smutel/go-netbox/commit/e1b574625ef04207b417ff816b6b584447568225))
* Init the project ([07fb3e1](https://github.com/smutel/go-netbox/commit/07fb3e1cced9502d91fba7babcaf2361a2779f54))


### Bug Fixes

* Client generation is not done anymore ([4372ce2](https://github.com/smutel/go-netbox/commit/4372ce2c7da78d0751236baa6bf8807d15f1ed1b))
* Fix patch to ipaddress assignedobject ([3b540e3](https://github.com/smutel/go-netbox/commit/3b540e3dac9a3f2b675f360501a23605b515b576))
* Remove additional v from swagger filename ([01b6671](https://github.com/smutel/go-netbox/commit/01b667184a70880c792ca52391940d1826768f55))
* Rename master to main ([b063a64](https://github.com/smutel/go-netbox/commit/b063a6467716970bfe6790790f408c2caeb68b3d))
