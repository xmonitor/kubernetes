# cobra readme

## Overview

* 兼容 POSIX flags；
* 可以有子命令；
* 可以自动产生 -h 和 -v；
* 可以在 bash/zsh 上自动产生命令提示；
* 可以函数 command alias；

## Concepts

cobra 的理念是命令行应该像句子一样易读，基本的概念有 command、argument 和 flag。command 是 action，args 是 对象，flags 则是 action 的修饰。

一个正常的句子是 `APPNAME VERB NOUN --ADJECTIVE`，对应的命令行应该是  `APPNAME COMMAND ARG --FLAG`。

如：

```Bash
hugo server --port=1313
```

则 server 是 command，port 是 flag。

## Commands

[Command](https://godoc.org/github.com/spf13/cobra#Command) 是一个行为的中心，其 Run 函数执行顺序是：

```Go
    //   * PersistentPreRun()
    //   * PreRun()
    //   * Run()
    //   * PostRun()
    //   * PersistentPostRun()
    // All functions get the same args, the arguments after the command name.
```
如注释所说，所有函数共用共同的参数。

还有一个重要的函数是 Context：

```Go
// Context returns underlying command context. If command wasn't executed with ExecuteContext Context returns Background context.
func (c *Command) Context() context.Context
```

## Flags

Cobra 既支持 POSIX flag，也支持 go 语言的 flag。如果一个 flag 是 persist 类型的 flag，则子命令可获取到这个 flag，如果是 local 类型则子命令获取不到。

Cobra 对 flag 的支持基于 pflag 库。

## Getting Started

Cobra  应用的目录结果一般如下：

```Go
  ▾ appName/
    ▾ cmd/
        add.go
        your.go
        commands.go
        here.go
      main.go
```

Cobra 应用的 main.go 一般有固定模式，其主要作用就是初始化 Cobra：

```Go
package main

import (
  "{pathToYourApp}/cmd"
)

func main() {
  cmd.Execute()
}
```

## Using the Cobra Library

Cobra 应用必须提供一个 root command 的实现，然后可在 rootCommand 上添加子命令扩展。

### Create rootCmd

```Go
  // app/cmd/root.go
  var rootCmd = &cobra.Command{
    Use:   "hugo",
    Short: "Hugo is a very fast static site generator",
    Long: `A Fast and Flexible Static Site Generator built with
                  love by spf13 and friends in Go.
                  Complete documentation is available at http://hugo.spf13.com`,
    Run: func(cmd *cobra.Command, args []string) {
      // Do Stuff Here
    },
  }
  
  func Execute() {
    if err := rootCmd.Execute(); err != nil {
      fmt.Println(err)
      os.Exit(1)
    }
  }
```

可以在 init 函数里添加 flags 和 配置处理，如：

```Go
  // app/cmd/root.go
  package cmd

  import (
  	"fmt"
  	"os"
  
  	homedir "github.com/mitchellh/go-homedir"
  	"github.com/spf13/cobra"
  	"github.com/spf13/viper"
  )
  
  var (
  	// Used for flags.
  	cfgFile     string
  	userLicense string
  
  	rootCmd = &cobra.Command{
  		Use:   "cobra",
  		Short: "A generator for Cobra based Applications",
  		Long: `Cobra is a CLI library for Go that empowers applications.
  This application is a tool to generate the needed files
  to quickly create a Cobra application.`,
  	}
  )
  
  // Execute executes the root command.
  func Execute() error {
  	return rootCmd.Execute()
  }
  
  func init() {
  	cobra.OnInitialize(initConfig)
  
  	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
  	rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
  	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
  	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
  	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
  	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
  	viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
  	viper.SetDefault("license", "apache")
  
  	rootCmd.AddCommand(addCmd)
  	rootCmd.AddCommand(initCmd)
  }
  
  func er(msg interface{}) {
  	fmt.Println("Error:", msg)
  	os.Exit(1)
  }
  
  func initConfig() {
  	if cfgFile != "" {
  		// Use config file from the flag.
  		viper.SetConfigFile(cfgFile)
  	} else {
  		// Find home directory.
  		home, err := homedir.Dir()
  		if err != nil {
  			er(err)
  		}
  
  		// Search config in home directory with name ".cobra" (without extension).
  		viper.AddConfigPath(home)
  		viper.SetConfigName(".cobra")
  	}
  
  	viper.AutomaticEnv()
  
  	if err := viper.ReadInConfig(); err == nil {
  		fmt.Println("Using config file:", viper.ConfigFileUsed())
  	}
  }
```

### Create additional commands

```Go
package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of Hugo",
  Long:  `All software has versions. This is Hugo's`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
  },
}
```

### Returning and handling errors

当需要对命令调用者返回错误时，command 的 RunE 会被调用：

```Go
package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(tryCmd)
}

var tryCmd = &cobra.Command{
  Use:   "try",
  Short: "Try and possibly fail at something",
  RunE: func(cmd *cobra.Command, args []string) error {
    err := someFunc()
    if err := nil {
	return err
    }
  },
}
```

## Working with Flags

rootCmd 的 flag 应该是一个 全局变量：

```Go
rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
```

### Local Flag on Parent Commands

默认情况下，子命令执行的时候只分析其 local flag，parent command 的 local flag 会被忽略，但如下情况除外：

```Go
command := cobra.Command{
  Use: "print [OPTIONS] [COMMANDS]",
  TraverseChildren: true,
}
```

TraverseChildren 为 true 时，其每个子命令执行时，所有子命令的 local flag 都会被分析。

### Bind Flags with Config

也可以通过 viper 库设定一个 command 的 flag:

```Go
var author string

func init() {
  rootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
  viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
}
```

### Required flags

默认情况下 flags 都是 optional 的，如果想把它变成 required，则需要显示声明：

```Go
rootCmd.Flags().StringVarP(&Region, "region", "r", "", "AWS region (required)")
rootCmd.MarkFlagRequired("region")
```

## Positional and Custom Arguments

Cobra 给出了一些其实现好的 参数 validation 函数：

* NoArgs - 没有参数则报错
* ArbitraryArgs - 可以接收任何参数 
* OnlyValidArgs - 如果收到的参数不再 ValidArgs 中则报错
* MinimumNArgs(int) - 参数数目少于 N 则报错
* MaximumNArgs(int) - 参数数目多于 N 则报错
* ExactArgs(int) - 参数数目不等于 N 则报错
* ExactValidArgs(int) - 如果没有 N 个参数或者有不在 ValidArgs 中的参数，则报错
* RangeArgs(min, max) - 如果参数数目不在 (min, max) 之间则报错

也可以自定义 validator：

```Go
var cmd = &cobra.Command{
  Short: "hello",
  Args: func(cmd *cobra.Command, args []string) error {
    if len(args) < 1 {
      return errors.New("requires a color argument")
    }
    if myapp.IsValidColor(args[0]) {
      return nil
    }
    return fmt.Errorf("invalid color specified: %s", args[0])
  },
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("Hello, World!")
  },
}
```

### PreRun and PostRun Hooks

command 的 Run 函数将以如下顺序执行：

* PersistentPreRun
* PreRun
* Run
* PostRun
* PersistentPostRun

如果子命令没有定义 Persistent*Run 系列函数，则子命令会继承父命令的函数。如下代码示例，subCmd 会执行 parentCmd 的 PersistentPreRun，但是 PersistentPostRun 则是执行自己定义的函数。

```Go
package main

import (
  "fmt"

  "github.com/spf13/cobra"
)

func main() {

  var rootCmd = &cobra.Command{
    Use:   "root [sub]",
    Short: "My root command",
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside rootCmd PersistentPreRun with args: %v\n", args)
    },
    PreRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside rootCmd PreRun with args: %v\n", args)
    },
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside rootCmd Run with args: %v\n", args)
    },
    PostRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside rootCmd PostRun with args: %v\n", args)
    },
    PersistentPostRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside rootCmd PersistentPostRun with args: %v\n", args)
    },
  }

  var subCmd = &cobra.Command{
    Use:   "sub [no options!]",
    Short: "My subcommand",
    PreRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside subCmd PreRun with args: %v\n", args)
    },
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside subCmd Run with args: %v\n", args)
    },
    PostRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside subCmd PostRun with args: %v\n", args)
    },
    PersistentPostRun: func(cmd *cobra.Command, args []string) {
      fmt.Printf("Inside subCmd PersistentPostRun with args: %v\n", args)
    },
  }

  rootCmd.AddCommand(subCmd)

  rootCmd.SetArgs([]string{""})
  rootCmd.Execute()
  fmt.Println()
  rootCmd.SetArgs([]string{"sub", "arg1", "arg2"})
  rootCmd.Execute()
}
```

output

```Go
Inside rootCmd PersistentPreRun with args: []
Inside rootCmd PreRun with args: []
Inside rootCmd Run with args: []
Inside rootCmd PostRun with args: []
Inside rootCmd PersistentPostRun with args: []

Inside rootCmd PersistentPreRun with args: [arg1 arg2]
Inside subCmd PreRun with args: [arg1 arg2]
Inside subCmd Run with args: [arg1 arg2]
Inside subCmd PostRun with args: [arg1 arg2]
Inside subCmd PersistentPostRun with args: [arg1 arg2]
```
