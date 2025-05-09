# Changelog

## [[v0.2.0]](https://github.com/mlange-42/ark-serde/compare/v0.1.3...v0.2.0)

### Features

- Adds option `Compress` for in-memory gzip compression (#18, #19)

## [[v0.1.3]](https://github.com/mlange-42/ark-serde/compare/v0.1.2...v0.1.3)

### Documentation

- Adds benchmarks for serialization and de-serialization (#12)

### Performance

- Adds a sub-project to profile (de)-serialization (#13)
- Slightly speeds up deserialization ba using a slice instead of a map for component infos (#13)
- Speed up by up to factor 2 by replacing `encoding/json` with `goccy/go-json` (#15)

## [[v0.1.2]](https://github.com/mlange-42/ark-serde/compare/v0.1.1...v0.1.2)

- Upgrade to Ark v0.4.0 (#9, #10)

## [[v0.1.1]](https://github.com/mlange-42/ark-serde/compare/v0.1.0...v0.1.1)

- Upgrade to Ark v0.3.0 (#7)

## [[v0.1.0]](https://github.com/mlange-42/ark-serde/commits/v0.1.0/)

- Initial release of [ark-serde](https://github.com/mlange-42/ark-serde)
