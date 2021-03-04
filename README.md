# memcache-go

Memcache wrapper with helpers.

Supports SASL authentication.

```go
func GetData() (data Data, err error) {

	client := NewClient("localhost:11211")

	callback := func() (interface{}, error) {
		// Calculate data
		return Data{Val: 1}, nil
	}

	data := Data{}
	err := client.GetSet("data-key", 10, &data, callback)
	if err != nil {
		return nil, err
	}

	return data, err
}
```
