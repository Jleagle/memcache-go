# memcache

![example workflow](https://github.com/Jleagle/memcache-go/actions/workflows/test.yml/badge.svg)

Supports SASL authentication.

```go
func GetData() (data Data, err error) {

	client := memcache.NewClient("localhost:11211")

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
