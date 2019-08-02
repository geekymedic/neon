# validator

`validaor`可以借助于`neon-cli`和`gogenerator`客户端工具，为一个结构体生成相应的校验函数，比如：

```go
//go:generate neon generate validate --src-file=book.go --struct-name=Book --package-name=controll
type Book struct{
  	Name string	`json:"string"`
  	Price float64 
  	Fact Factory	`json:"factory"`
}

type Factory struct {
  	Address string `json:"address"`
}
```



`Book`被生成如下的一个参数检查对象：

```go
package controll

import . "github.com/geekymedic/neon/utils/validator"

// define new struct
type BookValidator struct {
    Name *Field
    Price *Field
    FactAddress *Field
}

func NewBookValidator(book *Book) *BookValidator {
    return &BookValidator{
        Name: &Field{
            Tag: "Book.Name",
            Value: book.Name,
        }, Price: &Field{
            Tag: "Book.Price",
            Value: book.Price,
        }, FactAddress: &Field{
            Tag: "Book.Fact.Address",
            Value: book.Fact.Address,
        }, 
    }
}
```

这个文件默认保存在当前工作目录下。通过*BrookValidator*可以做一些参数检查工作，如：

```go
val := NewValidaotor()
valBook := NewBrookValidator(&Brook{})
err := val.Validate(valBook.Name, Min("1"), Max("10")).Validate(valBook.Price, Gt("10.9")).Validate(valBook.FactAddress, Contains("HelloWord")).Err()
if err != nil {
  	return err
}

.....
```

`Min`、`Max`、`Contains`为校验函数，目前已适配`gopkg.in/go-playground/validator.v8`所提供的大部分函数，`v9`版本接下来会继续适配。

- [:link:validator-v8](https://github.com/go-playground/validator/tree/v8)

- [:link:validator-v9](https://github.com/go-playground/validator/tree/v9)

## V8

### 逻辑运算
#### Eq(相等)

`Eq`根据`Field.Value`的数据类型，采用的不同的语义：

- 数字类型

比较数值

```go
NewValidator().(&Field{"age", 90}, Eq("90")) == true
NewValidator().(&Field{"age", 90.1}, Eq("90.1")) == true
```

- 字符串

```go
NewValidator().(&Field{"address", "ShenZhen", Eq("ShenZhen")) == true
```

- 数组、字典类型

比较长度

```go
NewValidator().(&Field{"address", []int{100, 200}, Eq("2")) == true
NewValidator().(&Field{"address", map[string]int{"Lucy": 100, "Lisi": 200}, Eq("2")) == true
```

- 时间

和当前的时间比较

#### Max(最大值)、Min(最小值)

`Max、Min`根据`Field.Value`的数据类型，采用的不同的语义：

- 数字类型

比较取值范围

```go
NewValidator().(&Field{"age", 100, Max("4")) == false
NewValidator().(&Field{"age", 3, Max("4")) == true
NewValidator().(&Field{"age", 3, Min("4")) == false              
NewValidator().(&Field{"age", 5, Min("4")) == true

NewValidator().(&Field{"age", 90, Min("4"), Max("70")) == false
NewValidator().(&Field{"age", 5, Min("4"), Max("70")) == true
```

- 字符串、数组、字典类型

比较长度

```go
NewValidator().(&Field{"names", []string{"a", "b", "c"}, Min("1"), Max("2")) == false
NewValidator().(&Field{"names", []string{"a", "b", "c"}, Min("1"), Max("10")) == true
```

- 时间

和当前的时间比较

#### Gt（大于）、Gte（大于等于）、Lt（小于）、Lte(小于等于)

- 数字类型

数值比较

```go
NewValidator().(&Field{"age", 100, Gt("4")) == true
NewValidator().(&Field{"age", 2, Gt("4")) == false
NewValidator().(&Field{"age", 4, Gte("4")) == true
NewValidator().(&Field{"age", 3, Gte("4")) == false
NewValidator().(&Field{"age", 3, Lt("4")) == true
NewValidator().(&Field{"age", 4, Lt("4")) == false
NewValidator().(&Field{"age", 4, Lte("4")) == true
NewValidator().(&Field{"age", 5, Lte("4")) == false
```

- 字符串、数组、字典

比较长度

```go
NewValidator().(&Field{"md5", "3364847f5601fc1bc2f04853197985bd", Gt("32")) == true
NewValidator().(&Field{"address_list", []string{"Beijing", "ShangHai"}, Gte("2")) == true                     
NewValidator().(&Field{"info", map[string]int{"Beijing": 90, "ShangHai": 70}, Lt("3")) == true
NewValidator().(&Field{"info", map[string]int{"Beijing": 90, "ShangHai": 70}, Lte("2")) == true                       
```



### 字符串比较

#### Contains、ContainsAny

> 只支持字符串类型

```go
NewValidator().(&Field{"md5", "3364847f5601fc1bc2f04853197985bd", Contains("3364847f5601fc1bc2f04853197985bd")) == true
NewValidator().(&Field{"md5", "3364847f5601fc1bc2f04853197985bd", Contains("3364847f5601fc1bc2f04853197985bc")) == false
       
NewValidator().(&Field{"name", "ShenZhen", ContainsAny("Saddss")) == true
NewValidator().(&Field{"name", "ShenZhen", ContainsAny("Good")) == false
```



#### ContainsRune

切割为`Rune(4个字节)`比较

```go
NewValidator().(&Field{"国家", "中", ContainsAny("日本")) == false
NewValidator().(&Field{"国家", "中国", ContainsAny("中国好")) == true
```

#### Excludes(不包含)、ExcludesAll(不包含任意一个)

```go
NewValidator().(&Field{"name", "abc", Excludes("J")) == true
NewValidator().(&Field{"name", "abc", ExcludesAll("J")) == true
NewValidator().(&Field{"name", "abc", ExcludesAll("cd")) ==false                      
```

#### ExcludesRune

切割为`Rune(4个字节)`比较

```go
NewValidator().(&Field{"国家", "中", ExcludesRune("日本")) == true
NewValidator().(&Field{"国家", "中国", ExcludesRune("中国好")) == false
```

### 其他类型的比较

#### IP相关

-  IsIP

```go
NewValidator().(&Field{"ip", "127.0.0.1.1", IsIP()) == false
NewValidator().(&Field{"ip", "127.0.0.1", IsIP()) == true
```

- IsIPv4

```go
NewValidator().(&Field{"ipv4", "127.0.0.1.1", IsIPv4()) == false
NewValidator().(&Field{"ipv4", "2001:da8:8000:1::81", IsIPv4()) == false 
NewValidator().(&Field{"ipv4", "183.14.132.105", IsIPv4()) == true  
```

- IsIPv6

```go
NewValidator().(&Field{"ipv4", "127.0.0.1.1", IsIPv6()) == false
NewValidator().(&Field{"ipv4", "2001:da8:8000:1::81", IsIPv6()) == true 
NewValidator().(&Field{"ipv4", "183.14.132.105", IsIPv6()) == false  
```

- IsIPAddrResolvable

```go
NewValidator().(&Field{"IsIPAddrResolvable", "103.123.3.1.0", IsIPAddrResolvable()) == false
NewValidator().(&Field{"IsIPAddrResolvable", "::1", IsIPAddrResolvable()) == true
```

- IsIP4AddrResolvable

```go
NewValidator().(&Field{"IsIP4AddrResolvable", "2001:da8:8000:1::81", IsIP4AddrResolvable()) == false
NewValidator().(&Field{"IsIP4AddrResolvable", "103.123.3.1", IsIP4AddrResolvable()) == true
```

- IsIP6AddrResolvable

```go
NewValidator().(&Field{"IsIP6AddrResolvable", "2001:da8:8000:1::81", IsIP6AddrResolvable()) == true
NewValidator().(&Field{"IsIP6AddrResolvable", "103.123.3.1", IsIP6AddrResolvable()) == false
```

- IsTCPAddrResolvable

```go
NewValidator().(&Field{"IsTCPAddrResolvable", "[::1]:80", IsTCPAddrResolvable()) == true
NewValidator().(&Field{"IsTCPAddrResolvable", "256.0.0.0:1", IsTCPAddrResolvable()) == false                      
```

- IsTCP4AddrResolvable

```go
NewValidator().(&Field{"IsTCP4AddrResolvable", "127.0.0.1:80", IsTCP4AddrResolvable()) == true
NewValidator().(&Field{"IsTCP4AddrResolvable", "[::1]:80", IsTCP4AddrResolvable()) == false  
```

- IsTCP6AddrResolvable

```go
NewValidator().(&Field{"IsTCP6AddrResolvable", "[::1]:80", IsTCP6AddrResolvable()) == true
NewValidator().(&Field{"IsTCP6AddrResolvable", "127.0.0.1:80", IsTCP6AddrResolvable()) == false 
```

- IsUDPAddrResolvable

```go
NewValidator().(&Field{"IsTCP6AddrResolvable", "[::1]:80", IsTCP6AddrResolvable()) == true
NewValidator().(&Field{"IsTCP6AddrResolvable", "127.0.0.1:80", IsTCP6AddrResolvable()) == false 
```

- IsUDP4AddrResolvable

```go
NewValidator().(&Field{"IsUDP4AddrResolvable", "127.0.0.1:80", IsTCP4AddrResolvable()) == true
NewValidator().(&Field{"IsUDP4AddrResolvable", "[::1]:80", IsTCP4AddrResolvable()) == false  
```

- IsUDP6AddrResolvable

```go
NewValidator().(&Field{"IsUDP6AddrResolvable", "[::1]:80", IsUDP6AddrResolvable()) == true
NewValidator().(&Field{"IsUDP6AddrResolvable", "127.0.0.1:80", IsUDP6AddrResolvable()) == false 
```

- IsUnixAddrResolvable

```go
NewValidator().(&Field{"IsEmail", "v.sock", IsUnixAddrResolvable()) == false
```

#### Email

```go
NewValidator().(&Field{"Email", "helloword@gmail.com", IsEmail()) == true
```

#### Base64

```go
NewValidator().(&Field{"Base64", "aGVsbG93b3Jk", IsBase64()) == true
```

#### Mac

```go
NewValidator().(&Field{"Mac", "44-45-53-54-00-00", IsMac()) == true
```

#### 字母(Alpha)

```go
NewValidator().(&Field{"IsAlpha", "abcdxcEdffHG", IsAlpha()) == true
```

> 只支持字符串类型

#### 数字(IsNunmber)

```go
NewValidator().(&Field{"IsNunmber", "3881", IsNunmber()) == true
NewValidator().(&Field{"IsNunmber", "a41", IsNunmber()) == false
```

> 只支持字符串类型

#### Numeric

```go
NewValidator().(&Field{"IsNunmber", "381", IsNumeric()) == true
NewValidator().(&Field{"IsNunmber", "+381", IsNumeric()) == true
NewValidator().(&Field{"IsNunmber", "-381", IsNumeric()) == true
NewValidator().(&Field{"IsNunmber", "+3.81", IsNumeric()) == true
NewValidator().(&Field{"IsNunmber", "-3.81", IsNumeric()) == true
NewValidator().(&Field{"IsNunmber", "+3.", IsNumeric()) == false
```

> 只支持字符串类型

#### Hexadecimal

`16进制数`

```go
NewValidator().(&Field{"IsHexadecimal", "ff0044", IsHexadecimal()) == true
NewValidator().(&Field{"IsHexadecimal", "abcdefg", IsHexadecimal()) == true
NewValidator().(&Field{"IsHexadecimal", "h", IsHexadecimal()) == false
```

> 只支持字符串类型

#### 字母或者数字(IsAlphanum)

```go
NewValidator().(&Field{"IsAlphanum", "abcdxcEdffHG99120aaa99", IsAlphanum()) == true
```

> 只支持字符串类型

#### 时间

- IsTime

```go
NewValidator().(&Field{"IsTime", "2019-02-03 00:00:00", IsTime()) == true
```

- TimeAfter

```go
NewValidator().(&Field{"TimeAfter", time.Now().Add(time.Second), TimeAfter(time.Now())) == true
NewValidator().(&Field{"TimeAfter", time.Now(), TimeAfter("2080-02-03 00:00:00")) == false
```

- TimeBefore

```go
NewValidator().(&Field{"TimeBefore", time.Now(), TimeAfter(time.Now().Add(time.Second))) == true
NewValidator().(&Field{"TimeBefore", time.Now(), TimeAfter("2010-02-03 00:00:00")) == false
```

> 1：时区使用系统设置的时区
>
> 2：时间除了可以使用时间对象外， 还可以使用字符串作为参数传递，内部会对参数进行` time.ParseInLocation("2006-01-02 15:04:05", typ, time.Local)`转为当地时间

#### 经纬度

- 经度:`IsLatitude`

取值范围：`[-90, +90]`

```go
NewValidator().(&Field{"Coordinates", "80"}, IsLatitude()) == true
```
> 只支持字符串类型

- 维度:`IsLongitude`

取值范围：`[-180, -90), (90, 180]`

```go
NewValidator().(&Field{"Coordinates", "91"}, IsLongitude()) == true
```

> 只支持字符串类型

#### URL、URI

```go
NewValidator().(&Field{"IsURI", "http://foo.bar#com", IsURI()) == true
NewValidator().(&Field{"IsURI", "http://foobar.com", IsURI()) == true          
NewValidator().(&Field{"IsURI", "http://foobar.org:8080/", IsURI()) == true  
NewValidator().(&Field{"IsURI", "foobar.com", IsURI()) == false
                       
NewValidator().(&Field{"IsURL", "http://foo.bar#com", IsURL()) == false
NewValidator().(&Field{"IsURL", "http://foobar.coffee/"}, IsURL()) == true
NewValidator().(&Field{"IsURL", "foobar.com"}, IsURL()) == false
```

- [:link:https://tools.ietf.org/html/rfc3986#section-2.1](https://tools.ietf.org/html/rfc3986#section-2.1)

#### ASCII

```go
NewValidator().(&Field{"IsASCII", "foobar"}, IsASCII()) == true
NewValidator().(&Field{"IsASCII", "ｘｙｚ０９８"}, IsASCII()) == false
NewValidator().(&Field{"IsASCII", "ｶﾀｶﾅ"}, IsASCII()) == false
```

#### UUID、UUID3、UUID4、UUID5

```go
NewValidator().(&Field{"IsUUID5", "987fbc97-4bed-5078-9f07-9141ba07c9f3"}, IsUUID5()) == true
NewValidator().(&Field{"IsUUID4", "625e63f3-58f5-40b7-83a1-a72ad31acffb"}, IsUUID4()) == true
NewValidator().(&Field{"IsUUID3", "a987fbc9-4bed-3078-cf07-9141ba07c9f3"}, IsUUID3()) == true
NewValidator().(&Field{"IsUUID", "a987fbc9-4bed-3078-cf07-9141ba07c9f3"}, IsUUID()) == true
```

#### CIDR、CIDRv4、CIDRv6

```go
NewValidator().(&Field{"IsCIDR", "192.168.255.254/24"}, IsUUID()) == true
NewValidator().(&Field{"IsCIDR", "192.168.255.254/48"}, IsCIDR()) == false
NewValidator().(&Field{"IsCIDRv4", "172.16.0.1"}, IsUUIDv4()) == true
NewValidator().(&Field{"IsCIDRv6", "2001:cdba:0000:0000:0000:0000:3257:9652"}, IsUUIDv6()) == true
```

#### DataURI

#### SSN

#### 颜色

- HexColor
- RGB
- RGBA
- HSL
- HSLA

#### ISBN

- ISBN
- ISBN10
- ISBN13
