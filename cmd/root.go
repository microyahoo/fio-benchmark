package cmd

import (
	"flag"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/microyahoo/fio-benchmark/pkg/server"
	genericServer "github.com/microyahoo/fio-benchmark/pkg/server"
)

// fioBenchmarkOptions defines the options of fio benchmark
type fioBenchmarkOptions struct {
	jobFile      string // TODO
	cfgFile      string
	outputFile   string
	chartFile    string
	dryrun       bool
	renderFormat string
}

func newFioBenchmarkOptions() *fioBenchmarkOptions {
	return &fioBenchmarkOptions{
		// TODO
	}
}

func NewFioCommand() *cobra.Command {
	o := newFioBenchmarkOptions()
	cmds := &cobra.Command{
		Use: "fio-benchmark",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(o.Run(genericServer.SetupSignalHandler()))
		},
	}
	cmds.Flags().SortFlags = false
	klog.InitFlags(nil)
	// Make cobra aware of select glog flags
	// Enabling all flags causes unwanted deprecation warnings
	// from glog to always print in plugin mode
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("v"))
	// pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("logtostderr"))
	pflag.CommandLine.Set("logtostderr", "true")

	// cmds.Flags().StringVar(&o.jobFile, "job-file", "", "fio benchmark job file") // TODO(zhengliang)
	cmds.Flags().StringVar(&o.outputFile, "output-file", "", "redirect fio benchmark result to output file")
	cmds.Flags().StringVar(&o.renderFormat, "render-format", "", "redirect fio benchmark result to output file with rendered format, eg. table, html, markdown, csv")
	cmds.Flags().StringVar(&o.cfgFile, "config-file", "", "fio benchmark config file, which will be ignored if job file is specified")
	cmds.Flags().StringVar(&o.chartFile, "chart-file", "", "echarts file for fio benchmark result")
	cmds.Flags().BoolVar(&o.dryrun, "dryrun", true, "dry-run")

	cmds.AddCommand(versionCmd, chartsCmd)

	return cmds
}

func (o *fioBenchmarkOptions) Run(stopCh <-chan struct{}) error {
	klog.Info("Starting fio benchmark")
	klog.V(4).Infof("fio benchmark options(job-file: %s, config-file: %s)",
		o.jobFile, o.cfgFile)

	server, err := server.NewFioServer(
		server.WithJobFile(o.jobFile),
		server.WithCfgFile(o.cfgFile),
		server.WithChartFile(o.chartFile),
		server.WithOutputFile(o.outputFile),
		server.WithRenderFormat(o.renderFormat),
		server.WithDryrun(o.dryrun))
	if err != nil {
		return err
	}
	genericServer.RegisterInterruptHandler(server.Close)

	err = server.Run(stopCh)
	if err != nil {
		return err
	}
	return nil
}
