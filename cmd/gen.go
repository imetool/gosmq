package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/imetool/gosmq/internal/gen"
	"github.com/imetool/goutil/util"
	"github.com/spf13/cobra"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "转换赛码表",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(args, len(args))
		if len(args) == 0 {
			fmt.Println("交互模式")
			GenWithSurvey()
			return
		}

		table := Config.Gen()
		path := "dict/" + util.GetFileName(Config.Path) + ".txt"
		gen.Write(table, path)
	},
}

var Config gen.Config

func init() {
	genCmd.PersistentFlags().StringVarP(&Config.Path, "dict", "i", "", "待转换的码表")
	genCmd.PersistentFlags().StringVarP(&Config.Format, "format", "f", "jisu", "待转换码表的格式")
	genCmd.PersistentFlags().StringVarP(&Config.SelectKeys, "select", "k", "_;'", "自定义选重键")
	genCmd.PersistentFlags().IntVarP(&Config.PushStart, "push", "p", 4, "起顶码长")
	genCmd.PersistentFlags().BoolVarP(&Config.SortByWordLen, "sort", "s", true, "按照词长重新排序")
}

func GenWithSurvey() {
	var conf gen.Config
	handle := func(err error) {
		if err != nil {
			if err == terminal.InterruptErr {
				log.Fatal("interrupted")
			}
		}
	}

	err := survey.AskOne(&survey.Input{
		Message: "待转换的码表路径:",
		Suggest: func(toComplete string) []string {
			files, _ := filepath.Glob(toComplete + "*")
			return files
		},
	}, &conf.Path, survey.WithValidator(survey.Required))
	handle(err)

	err = survey.AskOne(&survey.Select{
		Message: "待转换码表的格式:",
		Options: []string{"极速赛码表", "多多(词在前)", "冰凌(编码在前)", "极点"},
	}, &conf.Format)
	handle(err)

	mFormat := make(map[string]string)
	mFormat["极速赛码表"] = "jisu"
	mFormat["多多(词在前)"] = "duoduo"
	mFormat["冰凌(编码在前)"] = "bingling"
	mFormat["极点"] = "jidian"
	conf.Format = mFormat[conf.Format]

	if conf.Format != "jisu" {
		err = survey.AskOne(&survey.Input{
			Message: "起顶码长(码长大于等于此数，首选不会追加空格):",
			Default: "4",
		}, &conf.PushStart)
		handle(err)
	}

	err = survey.AskOne(&survey.Input{
		Message: "自定义选重键:",
		Default: "_;'",
	}, &conf.SelectKeys)
	handle(err)

	if conf.Format != "jisu" {
		err = survey.AskOne(&survey.Confirm{
			Message: "按照词长重新排序",
			Default: true,
		}, &conf.SortByWordLen)
		handle(err)
	}

	fmt.Println()
	fmt.Println(conf)

	table := conf.Gen()
	path := "dict/" + util.GetFileName(conf.Path) + ".txt"
	gen.Write(table, path)
}
