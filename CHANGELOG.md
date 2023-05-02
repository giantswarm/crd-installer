# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



## [Unreleased]

## [0.2.0] - 2023-05-02

### Changed

- Changed again path where to look for crds in calico source zip (made it work for all releases from 3.21 to 3.25).

## [0.1.2] - 2022-01-14

### Changed

- Changed path where to look for crds in calico source zip (changed in release 3.21).

## [0.1.1] - 2021-04-28

### Fixed

- Initialize scheme with apiextensions v1 types to fix unknown kind error.

### Added

- Add logging.
- Check CRD status to ensure that it has been accepted by Kubernetes.

## [0.1.0] - 2021-04-26

### Added

- Initial implementation.

[Unreleased]: https://github.com/giantswarm/crd-installer/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/giantswarm/crd-installer/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/giantswarm/crd-installer/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/giantswarm/crd-installer/releases/tag/v0.1.0
