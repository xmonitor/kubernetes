# Viper

## What is Viper?

Viper 主要用于处理命令行的 flag 以及相关配置，它支持：

* 设置默认值
* 从 JSON/TOML/YAML/HCL/envfile/Java properties config file 读取配置；
* 从 config files 预读，并关注其变化；
* 读取环境变量；
* 从 etcd/Consul 读取配置，包括关注其变化； 
* 从命令行读取值；
* 从 buffer 中读取值；
* 设置值；


Viper 的目的就是让你不用关注各种奇特的配置。

## Why Viper?

Viper 读取配置值的顺序是：

* explicit call to Set
* flag
* env
* config
* key/value store
* default

Viper 配置 kv 大小写敏感。

## Putting Values into Viper

### Establishing Defaults

设定默认值：

```Go
  viper.SetDefault("ContentDir", "content")
  viper.SetDefault("LayoutDir", "layouts")
  viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
```

## Reading Config Files

Viper 通过代码设定的路径去查询合规的 conf 文件。Viper 自身并没有设定默认的查询路径。

```Go
  viper.SetConfigName("config") // name of config file (without extension)
  viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
  viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
  viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
  viper.AddConfigPath(".")               // optionally look for config in the working directory
  err := viper.ReadInConfig() // Find and read the config file
  if err != nil { // Handle errors reading the config file
  	panic(fmt.Errorf("Fatal error config file: %s \n", err))
  }
```

可以设定没有找到配置文件的错误处理函数：

```Go
  if err := viper.ReadInConfig(); err != nil {
      if _, ok := err.(viper.ConfigFileNotFoundError); ok {
          // Config file not found; ignore error if desired
      } else {
          // Config file was found but another error was produced
      }
  }
```

## Writing Config Files

* WriteConfig 在预定义的路径下把 Viper 配置写进文件内，如果文件存在则会被写掉。如果没有预先设定文件则会写失败。
* SafeWriteConfig 在预定义的路径下把 Viper 配置写进文件内，如果文件存在则不会被写掉，写行为终止。如果没有预先设定文件则会写失败。
* WriteConfigAs 根据设定的路径，把配置写进去。
* SafeWriteConfigAs 根据设定的路径把配置写进去，如果文件存在则终止写行为。

```Go
  viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
  viper.SafeWriteConfig()
  viper.WriteConfigAs("/path/to/my/.config")
  viper.SafeWriteConfigAs("/path/to/my/.config") // will error since it has already been written
  viper.SafeWriteConfigAs("/path/to/my/.other_config")
```

## Watching and re-reading config files

除了读取配置，还可以 watch 配置文件变化并热加载。

```Go
  viper.WatchConfig()
  viper.OnConfigChange(func(e fsnotify.Event) {
  	fmt.Println("Config file changed:", e.Name)
  })
```

## Reading Config from io.Reader

用户可自定义配置文件格式，如下：

```Go
  viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")
  
  // any approach to require this configuration into your program.
  var yamlExample = []byte(`
  Hacker: true
  name: steve
  hobbies:
  - skateboarding
  - snowboarding
  - go
  clothing:
    jacket: leather
    trousers: denim
  age: 35
  eyes : brown
  beard: true
  `)
  
  viper.ReadConfig(bytes.NewBuffer(yamlExample))
  
  viper.Get("name") // this would be "steve"
```

## Setting Overrides

重新在代码中设定相关配置的值：

```go
  viper.Set("Verbose", true)
  viper.Set("LogFile", LogFile)
```

### Registering and Using Aliases

设定配置项 key 的 alias：

```go
  viper.RegisterAlias("loud", "Verbose")
  
  viper.Set("verbose", true) // same result as next line
  viper.Set("loud", true)   // same result as prior line
  
  viper.GetBool("loud") // true
  viper.GetBool("verbose") // true
```

### Working with Environment Variables

如下函数可用于读写环境变量：

* AutomaticEnv()
* BindEnv(string...) : error
* SetEnvPrefix(string)
* SetEnvKeyReplacer(string...) *strings.Replacer
* AllowEmptyEnv(bool)

viper 对环境变量是大小写敏感的。

`SetEnvPrefix` 设定环境变量的前缀，`BindEnv` 和 `AutomaticEnv` 会被其影响。

`BindEnv` 有两个参数：第一个是 key，第二个则是 环境变量。如果用户没有提供环境变量，则 viper 会自动用 ` prefix + "_" + the key name` 到所有  ALL CAPS 中去查询。如果环境变量已经设定，则它不会给环境变量参数前加上 prefix。如第二个参数是 `id`，则 Viper 会读取变量 `ID` 的环境变量的值。

`AutomaticEnv` 用于和 `SetEnvPrefix` 配合使用，当 Viper 调用这个函数后，`viper.Get` 每次都会先读取环境变量的值，且其读取的 key 是其大写形式。

### ENV example

```go
  SetEnvPrefix("spf") // will be uppercased automatically
  BindEnv("id")
  
  os.Setenv("SPF_ID", "13") // typically done outside of the app
  
  id := Get("id") // 13
```

## Working with Flags

Viper能够绑定到flag。

就像BindEnv，在调用绑定方法时，不会设置该值。这意味着您可以尽早绑定，甚至可以在init()函数中绑定 。

对于单个标志，该BindPFlag()方法提供此功能。

```go
  serverCmd.Flags().Int("port", 1138, "Port to run Application server on")
  viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
```

您还可以绑定一组现有的pflags（pflag.FlagSet） 

```go
  pflag.Int("flagname", 1234, "help message for flagname")
 
  pflag.Parse()
  viper.BindPFlags(pflag.CommandLine)
 
  i := viper.GetInt("flagname")
```

在Viper中使用pflag并不排除使用 标准库中使用标志包的其他包。pflag包可以通过导入这些标志来处理为标志包定义的标志。这是通过调用名为AddGoFlagSet（）的pflag包提供的便利函数来实现的。

```go
  package main
  import (
  	"flag"
  	"github.com/spf13/pflag"
  )
  func main() {
   
  	// using standard library "flag" package
  	flag.Int("flagname", 1234, "help message for flagname")
   
  	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
  	pflag.Parse()
  	viper.BindPFlags(pflag.CommandLine)
   
  	i := viper.GetInt("flagname") // retrieve value from viper
   
  	...
  }
```

### Flag interfaces

如果您不使用，Viper提供两个Go接口来绑定其他标志系统Pflags

FlagValue代表一个标志。这是一个关于如何实现此接口的非常简单的示例：

```go
  type myFlag struct {}
  func (f myFlag) HasChanged() bool { return false }
  func (f myFlag) Name() string { return "my-flag-name" }
  func (f myFlag) ValueString() string { return "my-flag-value" }
  func (f myFlag) ValueType() string { return "string" }
```

一旦你的flag实现了这个接口，你可以告诉Viper绑定它：

```Go
  viper.BindFlagValue("my-flag-name", myFlag{})
```

FlagValueSet 代表了一组 flag，下面示例给出了如何实现：

```Go
  type myFlagSet struct {
  	flags []myFlag
  }
  
  func (f myFlagSet) VisitAll(fn func(FlagValue)) {
  	for _, flag := range flags {
  		fn(flag)
  	}
  }
```

viper 绑定：

```Go
  fSet := myFlagSet{
  	flags: []myFlag{myFlag{}, myFlag{}},
  }
  viper.BindFlagValues("my-flags", fSet)
```

### Remote Key/Value Store Support 

要在Viper中启用远程支持，请对viper/remote 包进行空白导入：

```go
  import _ "github.com/spf13/viper/remote"
```

Viper将读取key/value存储（如etcd或Consul）中的路径检索的配置字符串（如JSON，TOML，YAML或HCL）。这些值优先于默认值，但会被从磁盘，标志或环境变量检索的配置值覆盖。

Viper使用crypt从K / V存储中检索配置，这意味着您可以存储加密的配置值，并在拥有正确的gpg密钥环时自动解密。加密是可选的。

您可以将远程配置与本地配置结合使用，也可以独立使用。

crypt有一个命令行帮助程序，您可以使用它来将配置放入K / V存储区。crypt在http://127.0.0.1:4001上默认为etcd 。

```Bash
  $ go get github.com/xordataexchange/crypt/bin/crypt 
  $ crypt set -plaintext /config/hugo.json /Users/hugo/settings/config.json
```

确认您的值已设置：

```Bash
  $ crypt get -plaintext /config/hugo.json
```

## Remote Key/Value Store Example - Unencrypted

> etcd

```Go
  viper.AddRemoteProvider("etcd", "http://127.0.0.1:4001","/config/hugo.json")
  viper.SetConfigType("json") //因为字节流中没有文件扩展名，支持的扩展名是“json”，“toml”，“yaml”，“yml”，“properties”，“props”，“prop”
  err := viper.ReadRemoteConfig()
```

> Consul

```Go
  {
      "port": 8080,
      "hostname": "myhostname.com"
  }

  viper.AddRemoteProvider("consul", "localhost:8500", "MY_CONSUL_KEY")
  viper.SetConfigType("json") // Need to explicitly set this to json
  err := viper.ReadRemoteConfig()
  
  fmt.Println(viper.Get("port")) // 8080
  fmt.Println(viper.Get("hostname")) // myhostname.com
```

> Filestore

```Go
  viper.AddRemoteProvider("firestore", "google-cloud-project-id", "collection/document")
  viper.SetConfigType("json") // Config's format: "json", "toml", "yaml", "yml"
  err := viper.ReadRemoteConfig()
```

## Remote Key/Value Store Example - Encrypted

```Go
viper.AddSecureRemoteProvider("etcd","http://127.0.0.1:4001","/config/hugo.json","/etc/secrets/mykeyring.gpg")
viper.SetConfigType("json") //因为字节流中没有文件扩展名，支持的扩展名是“json”，“toml”，“yaml”，“yml”，“properties”，“props”，“prop” 
err := viper.ReadRemoteConfig()
```

## Watching Changes in etcd - Unencrypted

```Go
//或者，您可以创建一个新的viper实例
var runtime_viper = viper.New()
 
runtime_viper.AddRemoteProvider("etcd", "http://127.0.0.1:4001", "/config/hugo.yml")
runtime_viper.SetConfigType("yaml")
 
// 第一次从远程配置中读取
err := runtime_viper.ReadRemoteConfig()
 
//解密配置
runtime_viper.Unmarshal(&runtime_conf)
 
// 打开一个goroutine来永远监听远程变化
go func(){
	for {
	    time.Sleep(time.Second * 5) // 每次请求后延迟
	    err := runtime_viper.WatchRemoteConfig()
	    if err != nil {
	        log.Errorf("unable to read remote config: %v", err)
	        continue
	    }
 
	    //将新配置解组到我们的运行时配置结构中。您还可以使用通道
        //实现信号以通知系统更改 
	    runtime_viper.Unmarshal(&runtime_conf)
	}
}()
```

## Getting Values From Viper

在Viper中，有几种方法可以根据值的类型获取值。存在以下功能和方法：

```Go
  Get(key string) : interface{}
  GetBool(key string) : bool
  GetFloat64(key string) : float64
  GetInt(key string) : int
  GetString(key string) : string
  GetStringMap(key string) : map[string]interface{}
  GetStringMapString(key string) : map[string]string
  GetStringSlice(key string) : []string
  GetTime(key string) : time.Time
  GetDuration(key string) : time.Duration
  IsSet(key string) : bool
  AllSettings() : map[string]interface{}
```

如果找不到，每个Get函数都将返回零值。IsSet()方法检查给定密钥是否存在。

实例：

```Go
  viper.GetString("logfile") // case-insensitive Setting & Getting
  if viper.GetBool("verbose") {
      fmt.Println("verbose enabled")
  }
```

## Accessing nested keys

访问器方法也接受深层嵌套键的格式化路径。例如，如果加载了以下JSON文件：

```JSON
  {
      "host": {
          "address": "localhost",
          "port": 5799
      },
      "datastore": {
          "metric": {
              "host": "127.0.0.1",
              "port": 3099
          },
          "warehouse": {
              "host": "198.0.0.1",
              "port": 2112
          }
      }
  }
```

Viper可以通过传递.分隔的键路径来访问嵌套字段：

```Go
  GetString("datastore.metric.host") // (returns "127.0.0.1")
```

这符合上面建立的优先规则; 搜索路径将在剩余的配置注册表中级联，直到找到。


例如，给定此配置文件，都datastore.metric.host和 datastore.metric.port已经定义（并且可以被覆盖）。如果另外datastore.metric.protocol在默认值中定义，Viper也会找到它。

但是，如果使用立即值datastore.metric覆盖（通过标志，环境变量，Set()方法，...），则所有子键 datastore.metric变为未定义，它们将被更高优先级的配置级别“遮蔽”。

最后，如果存在与分隔的键路径匹配的键，则将返回其值。例如

```JSON
  {
      "datastore.metric.host": "0.0.0.0",
      "host": {
          "address": "localhost",
          "port": 5799
      },
      "datastore": {
          "metric": {
              "host": "127.0.0.1",
              "port": 3099
          },
          "warehouse": {
              "host": "198.0.0.1",
              "port": 2112
          }
      }
  }
```
 
```Go
  GetString("datastore.metric.host") // returns "0.0.0.0"
```

## Extract sub-tree

例如

```YAML
  app:
    cache1:
      max-items: 100
      item-size: 64
    cache2:
      max-items: 200
      item-size: 80
``` 

执行后
 
```Go
  subv := viper.Sub("app.cache1")
``` 
 
subv 为:
 
```YAML
  max-items：100 
  item-size：64
```

假设

```Go
  func NewCache(cfg *Viper) *Cache {...}
```

它根据格式化为的配置信息创建缓存subv。现在可以轻松地分别创建这两个缓存：

```Go
  cfg1 := viper.Sub("app.cache1")
  cache1 := NewCache(cfg1)
   
  cfg2 := viper.Sub("app.cache2")
  cache2 := NewCache(cfg2)
```

## Unmarshaling

您还可以选择Unmarshaling all或特定值到struct，map等。

有两种方法可以做到这一点：

```Go
  Unmarshal(rawVal interface{}) : error
  UnmarshalKey(key string, rawVal interface{}) : error
```

例如：

```Go
  type config struct {
  	Port int
  	Name string
  	PathMap string `mapstructure:"path_map"`
  }
   
  var C config
   
  err := Unmarshal(&C)
  if err != nil {
  	t.Fatalf("unable to decode into struct, %v", err)
  }
```

## Marshalling to string

您可能需要将viper中保存的所有设置变为字符串，而不是将它们写入文件。您可以使用您喜欢的格式的marshaller和返回的配置AllSettings()。

```
  import (
      yaml "gopkg.in/yaml.v2"
      // ...
  ) 
   
  func yamlStringSettings() string {
      c := viper.AllSettings()
  	bs, err := yaml.Marshal(c)
  	if err != nil {
          t.Fatalf("unable to marshal config to YAML: %v", err)
      }
  	return string(bs)
  }
```

## Viper or Vipers?

Viper随时可以使用。开始使用Viper无需配置或初始化。由于大多数应用程序都希望使用单个中央存储库进行配置，因此viper软件包提供了此功能。它类似于单身人士。

在上面的所有示例中，他们演示了使用viper的单例式方法。

## Working with multiple vipers

您还可以创建许多不同的viper，以便在您的应用程序中使用。每个都有自己独特的配置和价值观。每个都可以从不同的配置文件，键值存储等中读取.viper包支持的所有功能都被镜像为viper上的方法。

例如

```Go
  x := viper.New()
  y := viper.New()
   
  x.SetDefault("ContentDir", "content")
  y.SetDefault("ContentDir", "foobar")
```

使用多viper时，用户可以跟踪不同的viper。
