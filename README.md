# mus-examples-go
Contains several examples of using the [mus-go](https://github.com/mus-format/mus-go)
serializer (each package is one example):
- `unsafe`: explains how the `unsafe` package can be used.
- `protobuf`: shows how to implement Protobuf encoding using mus-go.
- `dts`: demonstrates how [mus-dts-go](https://github.com/mus-format/mus-dts-go) 
  can be used.
- `versioning`: demonstrates data versioning.
- `marshal_func`: demonstrates how to use `MarshalMUS` function.
- `generic_marshal`: demonstrates how to implement generic marshal function.
- `oneof`: shows how to serialize an interface.
- `pm`: demonstrates how to use the `pm` package to serialize a graph or linked 
  list.
- `out_or_order`: shows how to deserialize values out of order.

More information can be found in the corresponding `main.go` files.