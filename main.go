package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"

	"github.com/jack139/go-infer/cli"
	"github.com/jack139/go-infer/types"

	"multinfer/models/qa"
	"multinfer/models/detpos"
)


var (
	rootCmd = &cobra.Command{
		Use:   "multinfer",
		Short: "multi-models infer platform",
	}
)

func init() {
	// 添加模型实例
	types.ModelList = append(types.ModelList, &qa.BertQA{})
	types.ModelList = append(types.ModelList, &detpos.DetPos{})

	// 命令行设置
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(cli.HttpCmd)
	rootCmd.AddCommand(cli.ServerCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
