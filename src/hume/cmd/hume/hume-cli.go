package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"hume/reader"
	"hume/source"
	"io/ioutil"
	"os"
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
		Use:   "validate [testing_file] [test_config]",
		Short: "Hume commandline testing client",
		Long:  "Hume client for running validation jobs from the commandline",
		Run: func(cmd *cobra.Command, args []string) {
			InitializeConfig()
			if len(args) < 2 {
				cmd.Usage()
				logrus.Fatalln("input file and configuration need to be provided")
			}

			inputFile := args[0]
			config := args[1]

			logrus.Info("Starting")

			sc, err := source.SourceConfigFromFile(config)
			if err != nil {
				logrus.Fatal(err)
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

	var trainCmd = &cobra.Command{
		Use: "train [training_file] [original_config] [output_config]",
		Short: "Hume commandline training client",
		Long: "Hume client for training an empty configuration file to pass given a training file",
		Run: func(cmd *cobra.Command, args []string) {
			InitializeConfig()
			if len(args) < 3 {
				cmd.Usage()
				logrus.Fatalln("testing file and empty configuration must be provided")
			}

			inputFile := args[0]
			config := args[1]
			outputConfig := args[2]

			f1, _ := os.Stat(config)
			f2, err := os.Stat(outputConfig)
			if err == nil {
				if f1.Name() == f2.Name() && f1.Size() == f2.Size() &&
				f1.Mode() == f2.Mode() && f1.ModTime() == f2.ModTime() {
					logrus.Fatal("Cannot overwrite input configuration file")
				}
			}

			logrus.Info("Starting")

			sc, err := source.SourceConfigFromFile(config)
			if err != nil {
				logrus.Fatal(err)
			}
			s := sc.GetSource()
			fr := s.Reader.(*reader.FileReader)
			fr.InputFile = inputFile
			s.Init()

			logrus.Info("Running")
			s.Collect()

			logrus.Info("Training")
			err_cnt, trained := s.Train()
			if err_cnt > 0 {
				logrus.Fatal(fmt.Sprintf("Training Error: %d/%d evaluators trained", trained-err_cnt, trained))
			} else {
				logrus.Infof("%d/%d evaluators trained", trained, trained)
			}
			// take configured output file, marshal newly configured source to output file
			out, _ := json.MarshalIndent(s,"","\t")
			out = bytes.Replace(out, []byte("\\u003c"), []byte("<"), -1)
			out = bytes.Replace(out, []byte("\\u003e"), []byte(">"), -1)
			out = bytes.Replace(out, []byte("\\u0026"), []byte("&"), -1)
			err = ioutil.WriteFile(outputConfig, out, 0644)
			if err != nil {
				logrus.Fatal(err)
			}

			logrus.Info("Done training")
		},
	}
	trainCmd.Flags().BoolVarP(&Debug, "debug", "d", false, "Log with debug statements")

	HumeCmd.AddCommand(runCmd, trainCmd)
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
