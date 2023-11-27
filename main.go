package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"

	"github.com/jack139/go-infer/cli"
	"github.com/jack139/go-infer/types"
	"github.com/jack139/go-infer/helper"

	// 推理模型
	"multinfer/models/detpos"
	"multinfer/models/bert_qa"
	"multinfer/models/keras_qa"
	"multinfer/models/ner_pack"
	"multinfer/models/yhface"
	"multinfer/models/talk2ui"
)


var (
	rootCmd = &cobra.Command{
		Use:   "multinfer",
		Short: "multi-models infer platform",
	}
)


func init() {
	// 重载 PersistentPreRunE
	cli.HttpCmd.PersistentPreRunE = preRun
	cli.ServerCmd.PersistentPreRunE = preRun

	// 命令行设置
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(cli.HttpCmd)
	rootCmd.AddCommand(cli.ServerCmd)
}


func preRun(cmd *cobra.Command, args []string) error {
	yaml, _ := cmd.Flags().GetString("yaml")
	helper.InitSettings(yaml)

	// 初始化时根据配置文件加载模型
	if helper.Settings.Customer["Load_Bert_QA"] == "1" {
		types.ModelList = append(types.ModelList, &bert_qa.BertQA{})
	}
	if helper.Settings.Customer["Load_Albert_QA"] == "1" {
		types.ModelList = append(types.ModelList, &keras_qa.AlbertQA{})
	}
	if helper.Settings.Customer["Load_Antigen"] == "1" {
		types.ModelList = append(types.ModelList, &detpos.DetPos{})
	}
	if helper.Settings.Customer["Load_NER_pack"] == "1" {
		types.ModelList = append(types.ModelList, &ner_pack.NER{})
	}
	if helper.Settings.Customer["Load_YHFace"] == "1" {
		types.ModelList = append(types.ModelList, &yhface.FaceLocate{})
		types.ModelList = append(types.ModelList, &yhface.FaceCheck{})
		types.ModelList = append(types.ModelList, &yhface.FaceVerify{})
		types.ModelList = append(types.ModelList, &yhface.FaceSearch{})
		types.ModelList = append(types.ModelList, &yhface.FaceFeatures{})
		types.ModelList = append(types.ModelList, &yhface.FaceMemory{})
	}
	if helper.Settings.Customer["Load_Talk2UI_pack"] == "1" {
		types.ModelList = append(types.ModelList, &talk2ui.Text2Order{})
		types.ModelList = append(types.ModelList, &talk2ui.Wav2Text{})
		types.ModelList = append(types.ModelList, &talk2ui.Wav2Order{})
	}

	return nil
}


func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

