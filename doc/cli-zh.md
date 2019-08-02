# neon​

命令输出格式：

1. 0 成功

   stdout：json or txt

2. -1 失败

   stderr: error message

## parse

### struct-list

- Command

```go
# neon parse struct-list --src-file="hello.go"
```

- Output

```go
{
	"StructList": ["StructNameA", "StructNameB"]
}
```

## generate

### transfer-struct

- Command

```go
# neon generate service transfer-struct --src-file="hello.go" --src-struct-name="PersonHttp" --dst-struct-name="PersonRpc" --src-var-name="personHttp" --dst-var-name="personRpc"
```

- Output

```go
var personHttp PersonHTTP, personRPC PersonRPC
personHttp.Name = personRPC.Name
personHttp.Age = personRPC.Age
```

