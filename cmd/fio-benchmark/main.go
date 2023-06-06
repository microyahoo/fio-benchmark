package main

import (
	"math/rand"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/component-base/logs"

	"github.com/microyahoo/fio-benchmark/cmd"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	logs.InitLogs()
	defer logs.FlushLogs()

	rootCmd := cmd.NewFioCommand()
	cobra.CheckErr(rootCmd.Execute())
}
