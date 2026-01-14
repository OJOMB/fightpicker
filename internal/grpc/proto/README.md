# Proto Specifications
This directory contains the Protocol Buffers (proto) specifications used for gRPC services in the FightPicker application. The proto files define the structure of the messages and the services that can be called remotely.

```bash
make gen-proto
```

This will generate the Go code for the proto files in `internal/grpc/gen` based on the proto specifications in this directory.

```bash
make clean-proto
```

This will delete the generated proto Go code if ever needed.