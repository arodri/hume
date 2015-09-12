package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"hume/reader"
	"hume/source"
)

var log = logrus.WithFields(logrus.Fields{
	"name": "main",
})

var HumeCmd = &cobra.Command{
	Use: "hume",
}

var Debug bool

func main() {
	var runCmd = &cobra.Command{
		Use:   "validate [file] [config]",
		Short: "Hume commandline client",
		Long:  "Hume client for running jobs from the commandline",
		Run: func(cmd *cobra.Command, args []string) {
			InitializeConfig()
			if len(args) < 2 {
				cmd.Usage()
				logrus.Fatalln("input file and configuration need to be provided")
			}

			inputFile := args[0]
			config := args[1]

			log.Info("Starting")

			sc, err := source.SourceConfigFromFile(config)
			if err != nil {
				log.Fatal(err)
			}
			s := sc.GetSource()
			fr := s.Reader.(*reader.FileReader)
			fr.InputFile = inputFile
			s.Init()

			logrus.Info("Running")
			s.Collect()

			logrus.Info("Evaluating")
			err_cnt, tested := s.Evaluate()
			if err_cnt > 0 {
				logrus.Fatal(fmt.Sprintf("Validation Error: %d/%d tests passed", tested-err_cnt, tested))
			} else {
				logrus.Infof("%d/%d tests passed", tested, tested)
			}

			logrus.Info("Done evaluating")
		},
	}
	runCmd.Flags().BoolVarP(&Debug, "debug", "d", false, "Log with debug statements")

	HumeCmd.AddCommand(runCmd)
	HumeCmd.Execute()
}

func LoadDefaultSettings() {
	viper.SetDefault("debug", false)
}

func InitializeConfig() {
	LoadDefaultSettings()
	if Debug {
		viper.Set("debug", true)
		logrus.SetLevel(logrus.DebugLevel)
	}
}
