# OpenAPI Specification
This directory contains the OpenAPI 3.0 specification for the FightPicker ReST API. oapi-codegen is used to generate the Go types for the DTOs.

```bash
make gen-dtos
```
this will generate the DTOs in `internal/http/dtos` based on the OpenAPI spec in this directory.

```bash
make clean-dtos
```
this will delete the generated DTOs if ever needed.

### gRPC
NB proto files are located in `internal/grpc/proto` not here.
see [here](../internal/grpc/proto/README.md) for more details on the gRPC implementation.
