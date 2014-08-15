Golang MPHF library
===================

This is a simple implementation of the Hash and Displace algorithm,
and reduces the size of the table by using a bit vector and rank
query.

Simple usage:

Create a slice of []mphf.KeyValue:

```go
type KeyValue struct {
        Key   []byte
        Value interface{}
}
```

Then simply call `phf = mphf.BuildMPHF(items)`

You can now query the resulting MPHF object with `value, ok := phf.Get(key)`
