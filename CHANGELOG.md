# Changelog

All notable changes to this project are documented here, newest first.

Entries are generated from [Conventional Commits](https://www.conventionalcommits.org)
and grouped by change type. This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### CI & Build

- Bump care to v0.8.1 by [@iberflow](https://github.com/iberflow) in [3fe3e4e](https://github.com/toaweme/structs/commit/3fe3e4edef4fc1257e0e96ac3f29e1a30ecd9b4a).
- Use stable go for release gate to avoid old-go.mod tool-install failures by [@iberflow](https://github.com/iberflow) in [ab6a3ad](https://github.com/toaweme/structs/commit/ab6a3ad0c436e5b424efa699c4dce93f5b883980).

## [0.3.1] - 2026-07-01

### CI & Build

- Bump care to v0.8.0 by [@iberflow](https://github.com/iberflow) in [670cdc1](https://github.com/toaweme/structs/commit/670cdc109295c0f851bbbb1fea8f5a2c7e4b1d62).
- Bump care to v0.7.1 and pin to commit sha by [@iberflow](https://github.com/iberflow) in [639147b](https://github.com/toaweme/structs/commit/639147bcad9b6052e1da37a996efbe26d0a463c7).
- Bump care to v0.6.0 and fix card-svg dark/light wiring by [@iberflow](https://github.com/iberflow) in [804ac1b](https://github.com/toaweme/structs/commit/804ac1b58c2db68d6445c8c65d690a7ffd3bb887).
- Pin care version as an input by [@iberflow](https://github.com/iberflow) in [e17dc14](https://github.com/toaweme/structs/commit/e17dc14d1f5933c5177f620ea8c5df1f1848c7ab).
- Bump care to v0.5.0 and pin to commit sha by [@iberflow](https://github.com/iberflow) in [f86d8d0](https://github.com/toaweme/structs/commit/f86d8d0c5780bd3a058df24e958892ec20581f7e).

## [0.3.0] - 2026-06-29

### Fixes

- Replace deprecated reflect.Ptr with reflect.Pointer by [@iberflow](https://github.com/iberflow) in [9633e63](https://github.com/toaweme/structs/commit/9633e63cba2e716c053936d433c8a1a6bfb062e7).

### Documentation

- Update README by [@iberflow](https://github.com/iberflow) in [0360e78](https://github.com/toaweme/structs/commit/0360e78bbaa9f7c877a851f8e57f91805da682d9).

### Chores & Other

- Align README, CHANGELOG, and quality workflow with org standards by [@iberflow](https://github.com/iberflow) in [09708e0](https://github.com/toaweme/structs/commit/09708e03666d3d4dacc5a7ab70238088f0e72b73).
- Update readme by [@iberflow](https://github.com/iberflow) in [ed76ac7](https://github.com/toaweme/structs/commit/ed76ac7c2ac5995681b2ad93dcf35dc98d6315bf).
- Update readme by [@iberflow](https://github.com/iberflow) in [8341525](https://github.com/toaweme/structs/commit/8341525d2eaf66b8cfb3e481c907fd6ac361d5ee).

## [0.2.0] - 2026-06-12

### Features

- Promote fields of unexported embedded structs by [@iberflow](https://github.com/iberflow) in [e29336f](https://github.com/toaweme/structs/commit/e29336fae9d7836f02a69d7ab5f7f569f1943cb8).
- Ci workflow by [@iberflow](https://github.com/iberflow) in [be3233e](https://github.com/toaweme/structs/commit/be3233e148d4ee3cd1d9c5d30453dcceb59efa41).

### Chores & Other

- Separate struct nesting and embedding tests by [@iberflow](https://github.com/iberflow) in [13efc66](https://github.com/toaweme/structs/commit/13efc6648fe7fca18bd591cf7fd0c9054cef558e).
- Cleanup validation and readme by [@iberflow](https://github.com/iberflow) in [105c682](https://github.com/toaweme/structs/commit/105c682d5a6e78f4078a3e6efc197278539cd3da).
- Tidy up by [@iberflow](https://github.com/iberflow) in [dbcd5a9](https://github.com/toaweme/structs/commit/dbcd5a9dc7dcede25a9c5fae267203da974d6ad3).

## [0.1.0] - 2026-06-11

### Features

- Initial commit by [@iberflow](https://github.com/iberflow) in [475eb77](https://github.com/toaweme/structs/commit/475eb777d8447694e65d1266a4df6f4dd26ec1e2).
- Support using *nested* map[string]any for setting struct field values instead of just fqn strings like parent.child.key by [@iberflow](https://github.com/iberflow) in [6d4e07f](https://github.com/toaweme/structs/commit/6d4e07f394169ceffa99490c65dfd061ad9bd6df).
- Set slices of structs from maps by [@iberflow](https://github.com/iberflow) in [a578e5f](https://github.com/toaweme/structs/commit/a578e5ff9b8ddb1d2ada63b44b73ca566bab4670).
- Configurable encoding tags by [@iberflow](https://github.com/iberflow) in [5a128e8](https://github.com/toaweme/structs/commit/5a128e86c4c0ae278b6b1a9b2e42c59021ff223c).
- Scalar-slice splitting via sep tag and oneof validation rule by [@iberflow](https://github.com/iberflow) in [a2ef59b](https://github.com/toaweme/structs/commit/a2ef59becd58568e90535bfc1b51cd0008bac244).
- MIT license by [@iberflow](https://github.com/iberflow) in [20b699e](https://github.com/toaweme/structs/commit/20b699ee3484e276a19098f7169d08a6f2ed89f9).

### Fixes

- Make ToAnySlice public by [@iberflow](https://github.com/iberflow) in [a83a31d](https://github.com/toaweme/structs/commit/a83a31d9ed9944155eb3ffc0b5435fb7b2c1225d).
- Manager and readme by [@iberflow](https://github.com/iberflow) in [15ffa3a](https://github.com/toaweme/structs/commit/15ffa3a47465be9e6dc34289044a79a8261bbd95).
- Example by [@iberflow](https://github.com/iberflow) in [335e741](https://github.com/toaweme/structs/commit/335e741caa5a93cb8f78a927d7368d75ecf5726b).
- Externalize rule funcs by [@iberflow](https://github.com/iberflow) in [5472b38](https://github.com/toaweme/structs/commit/5472b3851891d0ccd7d859c1c3724d00ee13f759).
- Manager by [@iberflow](https://github.com/iberflow) in [e802cb2](https://github.com/toaweme/structs/commit/e802cb278efee7fa7df0836a3ecc83a705ec833e).
- Add env vars by [@iberflow](https://github.com/iberflow) in [365abb9](https://github.com/toaweme/structs/commit/365abb97dfc37a13fc3d132040dc42c531c9df31).
- Unused var by [@iberflow](https://github.com/iberflow) in [918574e](https://github.com/toaweme/structs/commit/918574e6dceb5ab2765b822b26c4e1fe3636455a).
- Field type and tag parsing by [@iberflow](https://github.com/iberflow) in [bc119af](https://github.com/toaweme/structs/commit/bc119af84086be6665e55758fc431d85ede56f79).
- Default value handling by [@iberflow](https://github.com/iberflow) in [37be449](https://github.com/toaweme/structs/commit/37be44902ee3d1c509238fd45ac8db6be69c0d15).
- Setting nested structs by [@iberflow](https://github.com/iberflow) in [c2cdb3f](https://github.com/toaweme/structs/commit/c2cdb3f0394a06e8f62acc5fa2481fc1d4f133fc).
- Set and get struct fields + simplified logic with expanded Field by [@iberflow](https://github.com/iberflow) in [475af3d](https://github.com/toaweme/structs/commit/475af3df1b9869c59c58ad501ac353d18d15e174).
- Nested values by [@iberflow](https://github.com/iberflow) in [a7bff9b](https://github.com/toaweme/structs/commit/a7bff9b1f8229c192e9cd7ca75e57f2a593d0064).
- Set go.mod go version to 1.18 by [@iberflow](https://github.com/iberflow) in [fa5a91d](https://github.com/toaweme/structs/commit/fa5a91ddedec8a3170f5f931c022b9273eba02b6).
- Setting default values by [@iberflow](https://github.com/iberflow) in [7c023a0](https://github.com/toaweme/structs/commit/7c023a0260f113113bb9f9b1c67ff10be544e254).
- Setting slices by [@iberflow](https://github.com/iberflow) in [df41a18](https://github.com/toaweme/structs/commit/df41a18f8a9e1e31fe313ffa6798329994cd157d).
- Nil setting on any by [@iberflow](https://github.com/iberflow) in [9dcaaa6](https://github.com/toaweme/structs/commit/9dcaaa6b2627bfa77fb8d38b32446183230060b4).
- Slices by [@iberflow](https://github.com/iberflow) in [8d969f8](https://github.com/toaweme/structs/commit/8d969f827fd430da96e79a45157e56d979d7cb57).
- Struct field validation when field is set by [@iberflow](https://github.com/iberflow) in [f6cdf0a](https://github.com/toaweme/structs/commit/f6cdf0a20754d2fc8e2d2c73fe2a5e8d489ad81b).
- Strip stdlib tag options like ",omitempty" from parsed tag values by [@iberflow](https://github.com/iberflow) in [c374097](https://github.com/toaweme/structs/commit/c3740970ce0647d96a57f8b46c2c4ad6e427c92d).
- Handle numeric types in ToInt and ToFloat by [@iberflow](https://github.com/iberflow) in [dd5c9e0](https://github.com/toaweme/structs/commit/dd5c9e01e1449e6a0e6390fdf5b51099284083d1).
- Embedded struct field inlining by [@iberflow](https://github.com/iberflow) in [a39025a](https://github.com/toaweme/structs/commit/a39025a5c5522f68d68e7f571d3b63e1ebba768f).

### Chores & Other

- Cleanup by [@iberflow](https://github.com/iberflow) in [82b7cc3](https://github.com/toaweme/structs/commit/82b7cc35db38d221f6be7a787cdb27bded4a74da).
- Mod tidy by [@iberflow](https://github.com/iberflow) in [6f91ff9](https://github.com/toaweme/structs/commit/6f91ff9b7ffa9d002e32166f0f9803fd50a9ca96).
- Remove logging by [@iberflow](https://github.com/iberflow) in [12b1a7c](https://github.com/toaweme/structs/commit/12b1a7c12787848e19791c49e456e2ecb788f35d).
- Move to awee-ai org by [@iberflow](https://github.com/iberflow) in [11ce854](https://github.com/toaweme/structs/commit/11ce85456bdf4c7d3562773a67911d79ca768047).
- Bump deps by [@iberflow](https://github.com/iberflow) in [0ebd4e2](https://github.com/toaweme/structs/commit/0ebd4e2c76914a4780bf81b56b97c50259ee1312).
- Comment out log by [@iberflow](https://github.com/iberflow) in [a191447](https://github.com/toaweme/structs/commit/a191447dfecfd750066b3c96cdf3ed041071c3b6).
- Bump deps by [@iberflow](https://github.com/iberflow) in [2651ffc](https://github.com/toaweme/structs/commit/2651ffcf1b2e0504657cd56161dc3e0ad8b535a0).
- Move org by [@iberflow](https://github.com/iberflow) in [c5ce88d](https://github.com/toaweme/structs/commit/c5ce88d7ad1d3bcf8c19b8bb4aa1a49edc077152).
- Fmt by [@iberflow](https://github.com/iberflow) in [5d6a27d](https://github.com/toaweme/structs/commit/5d6a27df8e75e8bee1610dcc4c16ba1718f5b2e2).
- Migrate from testify to stdlib assertions by [@iberflow](https://github.com/iberflow) in [f436793](https://github.com/toaweme/structs/commit/f436793c26e91dbb9914d982c3c41f3b341c78ae).
- Fix lint issues by [@iberflow](https://github.com/iberflow) in [252a8ec](https://github.com/toaweme/structs/commit/252a8ecbd57e1d95aca5e0d3176e76095688ebaf).
- Tidy up comments by [@iberflow](https://github.com/iberflow) in [4ff0438](https://github.com/toaweme/structs/commit/4ff043858e88c959666f88dbd1153da3af2a5ca4).
- Downgrade go.mod go directive to 1.22 by [@iberflow](https://github.com/iberflow) in [e5bdde5](https://github.com/toaweme/structs/commit/e5bdde5910a8b1e3c69b8579e3dd2636e04a05e5).
- Cleanup sdk, examples, readme by [@iberflow](https://github.com/iberflow) in [8d45ed5](https://github.com/toaweme/structs/commit/8d45ed5132d226fcd8de35b7800902fa18d8f0be).

[Unreleased]: https://github.com/toaweme/structs/compare/v0.3.1...HEAD
[0.3.1]: https://github.com/toaweme/structs/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/toaweme/structs/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/toaweme/structs/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/toaweme/structs/releases/tag/v0.1.0
