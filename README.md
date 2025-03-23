# Ark Serde

[![Test status](https://img.shields.io/github/actions/workflow/status/mlange-42/ark-serde/tests.yml?branch=main&label=Tests&logo=github)](https://github.com/mlange-42/ark-serde/actions/workflows/tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/mlange-42/ark-serde/badge.svg?branch=main)](https://coveralls.io/github/mlange-42/ark-serde?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/mlange-42/ark-serde)](https://goreportcard.com/report/github.com/mlange-42/ark-serde)
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/ark-serde.svg)](https://pkg.go.dev/github.com/mlange-42/ark-serde)
[![GitHub](https://img.shields.io/badge/github-repo-blue?logo=github)](https://github.com/mlange-42/ark-serde)
[![MIT license](https://img.shields.io/github/license/mlange-42/ark-serde)](https://github.com/mlange-42/ark-serde/blob/main/LICENSE)

*Ark Serde* provides JSON serialization and deserialization for the [Ark](https://github.com/mlange-42/ark) Entity Component System (ECS).

<div align="center">

<a href="https://github.com/mlange-42/ark">
<img src="https://github.com/user-attachments/assets/4bbe57c6-2e16-43be-ad5e-0cf26c220f21" alt="Ark (logo)" width="500px" />
</a>

</div>

## Features

- Serialize/deserialize an entire Ark world in one line.
- Proper serialization of entity relations, as well as of entities stored in components.
- Skip arbitrary components and resources when serializing or deserializing.
- Optional in-memory GZIP compression for vast reduction of file sizes.

## Installation

```
go get github.com/mlange-42/ark-serde
```

## Usage

See the [API docs](https://pkg.go.dev/github.com/mlange-42/ark-serde) for more details and examples.  
[![Go Reference](https://pkg.go.dev/badge/github.com/mlange-42/ark-serde.svg)](https://pkg.go.dev/github.com/mlange-42/ark-serde)

Serialize a world:

```go
jsonData, err := arkserde.Serialize(&world)
if err != nil {
    // handle the error
}
```

Deserialize a world:

```go
err = arkserde.Deserialize(jsonData, &world)
if err != nil {
    // handle the error
}
```

## License

This project is distributed under the [MIT licence](./LICENSE).
