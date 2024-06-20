# mocka

> [!WARNING]
> Currently very experimental.

- **Type safety**
- **Generate mock objects** mapping args and results
- **Treat as stub objects** is also possible
- **Configure by yaml**, don't use `//go:generate`
- There is **no mather**
- There is no method chaining loop

## Usage

https://github.com/karamaru-alpha/mocka/tree/main/testdata

### Install
```bash
$ go install
$ mocka -h
Usage of mocka:
  -c string
    	Config yaml file path (default ".mocka.yaml")
```

### Configure

Place `.mocka.yaml` in the root directory.
```yaml
packages:
- github.com/karamaru-alpha/mocka/testdata:
    all: true
- fmt:
    interfaces:
    - Stringer
```

### Generate mocks

Execute `mocka` command on the root directory.
```sh
$ mocka
```

### Test

#### Expect

```go
m := testdata.NewMockHuman(t)
m.MockSay.Expect("Alice").Times(1).Return("Hello, Alice")
```

#### ConditionalExpect

```go
m := testdata.NewMockHuman(t)
m.MockSay.ConditionalExpect(func(name string) bool {
    return name == "Alice" || name == "Bob"
}).Times(1).Return("Hello, Alice or Bob")
```

#### Stabilize

```go
m := testdata.NewMockHuman(t)
m.MockSay.Stabilize(func(name string) string {
    return "Hello, " + name
})
```

## Comparison of other libraries

|                      | mocka | [uber-go/mock][mock] | [vektra/mockery][mockery] | [matryer/moq][moq] | [gojuno/minimock][minimock] |
|----------------------|-------|:--------------------:|:-------------------------:|:------------------:|:---------------------------:|
| Type safety          | ✅     |                      |                           |         ✅          |              ✅              |
| Map args and results | ✅     |          ✅           |             ✅             |                    |              ✅              |
| Treat as a stub      | ✅     |                      |                           |         👑         |              ✅              |
| Configure by yaml    | ✅     |                      |            👑             |                    |                             |
| No matcher           | ✅     |                      |                           |         ✅          |              ✅              |
| No method loop       | ✅     |                      |                           |         ✅          |                             |

[mock]: https://github.com/uber-go/mock
[mockery]: https://github.com/vektra/mockery
[moq]: https://github.com/matryer/moq
[minimock]: https://github.com/gojuno/minimock
[mockey]: https://github.com/bytedance/mockey
