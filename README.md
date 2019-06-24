# memcache-go

Memcache helper with retries

```go
func GetWork() (resp WorkResponse, err error) {

	err = helpers.GetMemcache().GetSetInterface("mem-key", 60, &resp, func() (interface{}, error) {
		return DoHeavyWork()
	})

	return resp, err
}
```
