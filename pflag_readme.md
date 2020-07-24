# pflag

## Description

用于替代 Go 语言的 flag 包，**支持类似于 POSIX/GNU 风格的 flag**。

pflag 包与 flag 包的工作原理甚至是代码实现都是类似的，下面是 pflag 相对 flag 的一些优势：

* 支持更加精细的参数类型：例如，flag 只支持 uint 和 uint64，而 pflag 额外支持 uint8、uint16、int32 等类型。
* 支持更多参数类型：ip、ip mask、ip net、count、以及所有类型的 slice 类型。
* 兼容标准 flag 库的 Flag 和 FlagSet：pflag 更像是对 flag 的扩展。
* 原生支持更丰富的功能：支持 shorthand、deprecated、hidden 等高级功能。

## Usage

pflag 用于替代 Go flag 包，用如下方式代码无修改即可启用：

```Go
  import flag "github.com/spf13/pflag"
```

意下情况 pflag 不能替代 Go flag: 定义了一个 flag 的 shorthand。

把一个 flag 绑定到一个变量的两种方法：

```Go
  var ip *int = flag.Int("flagname", 1234, "help message for flagname")
```

 OR

```Go
  var flagvar int
  func init() {
      flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
  }
```

区别是：方式一定义出来的变量都是指针形式，方式二定义出来的变量都是值形式。
使用时区别如下：

```Go
  fmt.Println("ip has value ", *ip)
  fmt.Println("flagvar has value ", flagvar)
```

定义完毕变量，调用如下函数开始分析：

```Go
  flag.Parse()
```

## Example

```Go
  package main
  
  import flag "github.com/spf13/pflag"
  import (
      "fmt"
      "strings"
  )
  
  // 定义命令行参数对应的变量
  var cliName = flag.StringP("name", "n", "nick", "Input Your Name")
  var cliAge = flag.IntP("age", "a",22, "Input Your Age")
  var cliGender = flag.StringP("gender", "g","male", "Input Your Gender")
  var cliOK = flag.BoolP("ok", "o", false, "Input Are You OK")
  var cliDes = flag.StringP("des-detail", "d", "", "Input Description")
  var cliOldFlag = flag.StringP("badflag", "b", "just for test", "Input badflag")
  
  func wordSepNormalizeFunc(f *flag.FlagSet, name string) flag.NormalizedName {
      from := []string{"-", "_"}
      to := "."
      for _, sep := range from {
          name = strings.Replace(name, sep, to, -1)
      }
      return flag.NormalizedName(name)
  }
  
  func main() {
      // 设置标准化参数名称的函数
      flag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
      
      // 为 age 参数设置 NoOptDefVal
      flag.Lookup("age").NoOptDefVal = "25"
  
      // 把 badflag 参数标记为即将废弃的，请用户使用 des-detail 参数
      flag.CommandLine.MarkDeprecated("badflag", "please use --des-detail instead")
      // 把 badflag 参数的 shorthand 标记为即将废弃的，请用户使用 des-detail 的 shorthand 参数
      flag.CommandLine.MarkShorthandDeprecated("badflag", "please use -d instead")
  
      // 在帮助文档中隐藏参数 gender
      flag.CommandLine.MarkHidden("badflag")
  
      // 把用户传递的命令行参数解析为对应变量的值
      flag.Parse()
  
      fmt.Println("name=", *cliName)
      fmt.Println("age=", *cliAge)
      fmt.Println("gender=", *cliGender)
      fmt.Println("ok=", *cliOK)
      fmt.Println("des=", *cliDes)
  }
```

### 布尔类型的参数

布尔类型的参数有下面几种写法

```Bash
  --flag               // 等同于 --flag=true        
  --flag=value
  --flag value         // 这种写法只有在没有设置默认值时才生效
```

### NoOptDefVal 用法

pflag 包支持通过简便的方式为参数设置默认值之外的值，实现方式为设置参数的 NoOptDefVal 属性：

```Go
  var cliAge = flag.IntP("age", "a", 22, "Input Your Age")
  flag.Lookup("age").NoOptDefVal = "25"
```

下面是传递参数的方式和参数最终的取值：

```Bash
  Parsed Arguments     Resulting Value
  --age=30             cliAge=30
  --age                cliAge=25
  [nothing]            cliAge=22
```

### shorthand

与 flag 包不同，在 pflag 包中，选项名称前面的 -- 和 - 是不一样的。- 表示 shorthand，-- 表示完整的选项名称。
除了最后一个 shorthand，其它的 shorthand 都必须是布尔类型的参数或者是具有默认值的参数。
所以对于布尔类型的参数和设置了 NoOptDefVal 的参数可以写成下面的形式：

```Bash
  -o
  -o=true
  // 注意，下面的写法是不正确的
  -o true
```

非布尔类型的参数和没有设置 NoOptDefVal 的参数的写法如下：

```Bash
  -g female
  -g=female
  -gfemale
```

常的使用中一般会混合上面的两类规则：

```Bash
  -aon "jack"
  -aon="jack"
  -aon"jack"
  -aonjack
  -oa=35
```

注意 -- 后面的参数不会被解析：

```Bash
  -oa=35 -- -gfemale
```

### 标准化参数的名称

如果我们创建了名称为 --des-detail 的参数，但是用户却在传参时写成了 --des_detail 或 --des.detail 会怎么样？默认情况下程序会报错退出，但是我们可以通过 pflag 提供的 SetNormalizeFunc 功能轻松的解决这个问题：

```Go
  func wordSepNormalizeFunc(f *flag.FlagSet, name string) flag.NormalizedName {
    from := []string{"-", "_"}
    to := "."
    for _, sep := range from {
        name = strings.Replace(name, sep, to, -1)
    }
    return flag.NormalizedName(name)
  }
  flag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
```

下面的写法也能正确设置参数了：

```Bash
  --des_detail="person detail"
```

### 把参数标记为即将废弃

在程序的不断升级中添加新的参数和废弃旧的参数都是常见的用例，pflag 包对废弃参数也提供了很好的支持。通过 MarkDeprecated 和 MarkShorthandDeprecated 方法可以分别把参数及其 shorthand 标记为废弃：

```Go
  // 把 badflag 参数标记为即将废弃的，请用户使用 des-detail 参数
  flag.CommandLine.MarkDeprecated("badflag", "please use --des-detail instead")
  // 把 badflag 参数的 shorthand 标记为即将废弃的，请用户使用 des-detail 的 shorthand 参数
  flag.CommandLine.MarkShorthandDeprecated("badflag", "please use -d instead")
```

### 在帮助文档中隐藏参数

pflag 包还支持在参数说明中隐藏参数的功能：

```Go
  // 在帮助文档中隐藏参数 badflag
  flag.CommandLine.MarkHidden("badflag")
```

其实在把参数标记为废弃时，同时也会设置隐藏参数。
