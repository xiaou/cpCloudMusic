/// global setting for logs, cmd, env, etc.
//
// cmd: --aa, --bb, --xx-yy etc...
// env: AA, BB, XX_YY ...

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/xiaou/cpCloudMusic/pkg/util"
)

var (
	Version = "0.0.1" // will modify by makefile
)

func setCmd() {
	pflag.Bool("version", false, "show version")

	// set cmd args
	//
	pflag.String("name", "我喜欢的音乐", "歌单名字.")
	pflag.String("out", "~/cpCloudMusic/我喜欢的音乐", "拷贝到目的路径文件夹.")
}

func setEnv() {
	pflag.VisitAll(func(f *pflag.Flag) {
		viper.BindPFlag(f.Name, f) // 让cmd的优先级高于ENV
		viper.SetDefault(f.Name, f.DefValue)
	})
	pflag.Parse()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_")) // XXX_XX
	viper.AutomaticEnv()
}

func BeginMain() {
	setCmd()
	flag.Set("logtostderr", "true")
	rand.Seed(time.Now().UTC().UnixNano())
	util.InitFlags()
	util.InitLogs()
	setEnv()
	if viper.GetBool("version") {
		fmt.Println("version:", Version)
		os.Exit(0)
	}
}

func EndMain() {
	util.FlushLogs()
}
