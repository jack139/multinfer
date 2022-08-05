package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"

	"github.com/jack139/go-infer/cli"
	"github.com/jack139/go-infer/types"

	//"multinfer/models/detpos"
	//"multinfer/models/bert_qa"
	"multinfer/models/keras_qa"
	
)


var (
	rootCmd = &cobra.Command{
		Use:   "multinfer",
		Short: "multi-models infer platform",
	}
)

func init() {
	// 添加模型实例
	//types.ModelList = append(types.ModelList, &detpos.DetPos{})
	//types.ModelList = append(types.ModelList, &bert_qa.BertQA{})
	types.ModelList = append(types.ModelList, &keras_qa.AlbertQA{})

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
