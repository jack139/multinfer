package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"

	"github.com/jack139/go-infer/cli"
	"github.com/jack139/go-infer/types"
	"github.com/jack139/go-infer/helper"

	"multinfer/models/detpos"
	"multinfer/models/bert_qa"
	"multinfer/models/keras_qa"
	
)


var (
	rootCmd = &cobra.Command{
		Use:   "multinfer",
		Short: "multi-models infer platform",
	}
)

func init() {
	// 命令行设置
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(cli.HttpCmd)
	rootCmd.AddCommand(cli.ServerCmd)
}

func main() {

	// 根据配置文件 加载模型
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		yaml, _ := cmd.Flags().GetString("yaml")
		helper.InitSettings(yaml)

		// 添加模型实例
		if helper.Settings.Customer["Load_Bert_QA"] == "1" {
			types.ModelList = append(types.ModelList, &bert_qa.BertQA{})
		}
		if helper.Settings.Customer["Load_Albert_QA"] == "1" {
			types.ModelList = append(types.ModelList, &keras_qa.AlbertQA{})
		}
		if helper.Settings.Customer["Load_Antigen"] == "1" {
			types.ModelList = append(types.ModelList, &detpos.DetPos{})
		}

		return nil
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
