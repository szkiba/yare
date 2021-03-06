# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<a name="unreleased"></a>
## [Unreleased]


<a name="v1.3.1"></a>
## [v1.3.1] - 2021-02-06
### Chore
- Improve changelog generation.

### Feat
- Convert http header names to canonical form in response.


<a name="v1.3.0"></a>
## [v1.3.0] - 2021-02-05
### Build
- Changelog generation moved to separated make target.

### Feat
- Support optinal parsers.

### Fix
- Ignore JWT parsing error in Authorization header.
- Fix handling response without body.


<a name="v1.2.1"></a>
## [v1.2.1] - 2021-02-04
### Build
- Added changelog generation during build.

### Docs
- Added 'go get' as install method.

### Fix
- Fix http 400 status for request without body and Content-Type.


<a name="v1.2.0"></a>
## [v1.2.0] - 2021-02-03
### Fix
- Yet another fix in package import path.


<a name="v1.1.0"></a>
## [v1.1.0] - 2021-02-03
### Fix
- Removed /v1 from yare package import path.


<a name="v1.0.1"></a>
## [v1.0.1] - 2021-02-03
### Docs
- Fix releases link.
- Fix package comment first line.


<a name="v1.0.0"></a>
## v1.0.0 - 2021-02-03

[Unreleased]: https://github.com/szkiba/yare/compare/v1.3.1...HEAD
[v1.3.1]: https://github.com/szkiba/yare/compare/v1.3.0...v1.3.1
[v1.3.0]: https://github.com/szkiba/yare/compare/v1.2.1...v1.3.0
[v1.2.1]: https://github.com/szkiba/yare/compare/v1.2.0...v1.2.1
[v1.2.0]: https://github.com/szkiba/yare/compare/v1.1.0...v1.2.0
[v1.1.0]: https://github.com/szkiba/yare/compare/v1.0.1...v1.1.0
[v1.0.1]: https://github.com/szkiba/yare/compare/v1.0.0...v1.0.1
