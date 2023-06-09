package server

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"k8s.io/klog/v2"

	"github.com/microyahoo/fio-benchmark/pkg/daemon/client"
	"github.com/microyahoo/fio-benchmark/pkg/util/exec"
)

const (
	WorkersLimit = 32
)

type ServerOptions struct {
	jobFile      string
	cfgFile      string
	chartFile    string
	outputFile   string
	dryrun       bool
	renderFormat string
}

type ServerOption func(*ServerOptions)

func WithJobFile(jobFile string) ServerOption {
	return func(opts *ServerOptions) {
		opts.jobFile = jobFile
	}
}

func WithChartFile(chartFile string) ServerOption {
	return func(opts *ServerOptions) {
		opts.chartFile = chartFile
	}
}

func WithCfgFile(cfgFile string) ServerOption {
	return func(opts *ServerOptions) {
		opts.cfgFile = cfgFile
	}
}

func WithOutputFile(outputFile string) ServerOption {
	return func(opts *ServerOptions) {
		opts.outputFile = outputFile
	}
}

func WithDryrun(dryrun bool) ServerOption {
	return func(opts *ServerOptions) {
		opts.dryrun = dryrun
	}
}

func WithRenderFormat(format string) ServerOption {
	return func(opts *ServerOptions) {
		opts.renderFormat = format
	}
}

type FioServer struct {
	Executor exec.Executor

	ctx        context.Context
	cancelFunc context.CancelFunc

	jobFile    string // TODO: not support
	cfgFile    string
	chartFile  string
	outputFile string

	wg          *sync.WaitGroup
	workerPool  chan *Worker
	jobListener chan *DelayedJob

	dryrun       bool
	renderFormat string

	lock    sync.Mutex
	results []*client.FioResult

	settings *TestSettings
}

func NewFioServer(options ...ServerOption) (*FioServer, error) {
	opts := &ServerOptions{}
	for _, option := range options {
		option(opts)
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	s := &FioServer{
		ctx:          ctx,
		cancelFunc:   cancelFunc,
		Executor:     &exec.CommandExecutor{},
		jobFile:      opts.jobFile,
		cfgFile:      opts.cfgFile,
		chartFile:    opts.chartFile,
		outputFile:   opts.outputFile,
		renderFormat: opts.renderFormat,
		dryrun:       opts.dryrun,
	}
	return s, nil
}

func (s *FioServer) Run(stopCh <-chan struct{}) (err error) {
	_, err = client.FioVersion(s.Executor)
	if err != nil {
		return err
	}
	settings, err := ParseSettings(s.cfgFile)
	if err != nil {
		return err
	}
	s.settings = settings
	err = s.doWork(settings)
	if err != nil {
		return err
	}
	s.printResults(s.outputFile, s.renderFormat)
	s.renderCharts()
	// <-stopCh
	return nil
}

func (s *FioServer) renderCharts() {
	var jobMap = make(map[string]map[string]map[string]map[string][]*client.FioJob) // map[rw][iodepth][bs][numjobs] => []Job
	for _, result := range s.results {
		for _, job := range result.Jobs {
			if _, ok1 := jobMap[job.JobOptions.RW]; !ok1 {
				jobMap[job.JobOptions.RW] = make(map[string]map[string]map[string][]*client.FioJob)
				jobMap[job.JobOptions.RW][job.JobOptions.IODepth] = make(map[string]map[string][]*client.FioJob)
				jobMap[job.JobOptions.RW][job.JobOptions.IODepth][job.JobOptions.BlockSize] = make(map[string][]*client.FioJob)
			} else {
				if _, ok2 := jobMap[job.JobOptions.RW][job.JobOptions.IODepth]; !ok2 {
					jobMap[job.JobOptions.RW][job.JobOptions.IODepth] = make(map[string]map[string][]*client.FioJob)
					jobMap[job.JobOptions.RW][job.JobOptions.IODepth][job.JobOptions.BlockSize] = make(map[string][]*client.FioJob)
				} else {
					if _, ok3 := jobMap[job.JobOptions.RW][job.JobOptions.IODepth][job.JobOptions.BlockSize]; !ok3 {
						jobMap[job.JobOptions.RW][job.JobOptions.IODepth][job.JobOptions.BlockSize] = make(map[string][]*client.FioJob)
					}
				}
			}
			jobMap[job.JobOptions.RW][job.JobOptions.IODepth][job.JobOptions.BlockSize][job.JobOptions.NumJobs] = append(
				jobMap[job.JobOptions.RW][job.JobOptions.IODepth][job.JobOptions.BlockSize][job.JobOptions.NumJobs], job)
		}
	}

	numJobs := s.settings.FioSettings.NumJobs
	sort.Slice(numJobs, func(i, j int) bool { return numJobs[i] < numJobs[j] })

	createLine := func(title string, yaxis string) *charts.Line {
		line := charts.NewLine()
		line.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title: title,
			}),
			charts.WithTooltipOpts(opts.Tooltip{Show: true}),
			charts.WithLegendOpts(opts.Legend{Show: true, Width: "50%", Left: "right"}),
			charts.WithInitializationOpts(opts.Initialization{
				Theme: "shine",
			}),
			charts.WithXAxisOpts(opts.XAxis{
				Name: "num_jobs",
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Name: yaxis,
			}),
		)
		line.SetXAxis(numJobs)
		return line
	}

	type metrics struct {
		readIOPS  float64
		readBw    float64
		readLat   float64
		writeIOPS float64
		writeBw   float64
		writeLat  float64
	}
	generateLines := func(lines []*charts.Line, metricsMap map[string][]*metrics) {
		for filename, metrics := range metricsMap {
			var (
				readIOPSLineData  []opts.LineData
				writeIOPSLineData []opts.LineData
				readBwLineData    []opts.LineData
				writeBwLineData   []opts.LineData
				readLatLineData   []opts.LineData
				writeLatLineData  []opts.LineData
			)
			for _, metric := range metrics {
				readIOPSLineData = append(readIOPSLineData, opts.LineData{Value: metric.readIOPS})
				writeIOPSLineData = append(writeIOPSLineData, opts.LineData{Value: metric.writeIOPS})
				readBwLineData = append(readBwLineData, opts.LineData{Value: metric.readBw})
				writeBwLineData = append(writeBwLineData, opts.LineData{Value: metric.writeBw})
				readLatLineData = append(readLatLineData, opts.LineData{Value: metric.readLat})
				writeLatLineData = append(writeLatLineData, opts.LineData{Value: metric.writeLat})
			}
			lineData := [][]opts.LineData{readIOPSLineData, writeIOPSLineData, readBwLineData, writeBwLineData, readLatLineData, writeLatLineData}
			for i, line := range lines {
				line.AddSeries(filename, lineData[i])
			}
		}
	}
	page := components.NewPage()
	for rw, rwMap := range jobMap {
		for iodepth, iodepthMap := range rwMap {
			for bs, bsMap := range iodepthMap {
				readIOPSLine := createLine(fmt.Sprintf("readiops-%s-%s-%s", rw, bs, iodepth), "iops")
				writeIOPSLine := createLine(fmt.Sprintf("writeiops-%s-%s-%s", rw, bs, iodepth), "iops")
				readBwLine := createLine(fmt.Sprintf("readbw-%s-%s-%s", rw, bs, iodepth), "bandwidth(KiB/s)")
				writeBwLine := createLine(fmt.Sprintf("writebw-%s-%s-%s", rw, bs, iodepth), "bandwidth(KiB/s)")
				readLatLine := createLine(fmt.Sprintf("readlat-%s-%s-%s(ms)", rw, bs, iodepth), "latency(ms)")
				writeLatLine := createLine(fmt.Sprintf("writelat-%s-%s-%s(ms)", rw, bs, iodepth), "latency(ms)")

				var filenameMap = make(map[string][]*metrics) // filename -> numJobs -> metrics
				for _, numJob := range numJobs {
					jobs := bsMap[fmt.Sprintf("%d", numJob)]
					for _, job := range jobs {
						filenameMap[job.JobOptions.FileName] = append(filenameMap[job.JobOptions.FileName], &metrics{
							readIOPS:  job.ReadResult.IOPSMean,
							readBw:    job.ReadResult.BWMean,
							readLat:   job.ReadResult.LatencyNs.Mean / 1000 / 1000, // ms
							writeIOPS: job.WriteResult.IOPSMean,
							writeBw:   job.WriteResult.BWMean,
							writeLat:  job.WriteResult.LatencyNs.Mean / 1000 / 1000, // ms
						})
					}
				}
				generateLines([]*charts.Line{readIOPSLine, writeIOPSLine, readBwLine, writeBwLine, readLatLine, writeLatLine}, filenameMap)
				page.AddCharts(readIOPSLine, writeIOPSLine, readBwLine, writeBwLine, readLatLine, writeLatLine)
			}
		}
	}
	chartFile := s.chartFile
	if chartFile == "" {
		chartFile = fmt.Sprintf("chart-%s.html", uuid.NewString())
	} else if !strings.HasSuffix(chartFile, "html") {
		chartFile = fmt.Sprintf("%s.html", chartFile)
	}
	f, err := os.Create(chartFile)
	if err != nil {
		panic(err)
	}
	err = page.Render(f)
	if err != nil {
		klog.Warningf("Faied to render charts", err)
	}
}

func (s *FioServer) doWork(settings *TestSettings) error {
	klog.Infof("fio test settings: %+v, use_all_disk: %t, workers: %d", settings.FioSettings, settings.UseAllDisks, settings.Workers)
	workQueue, err := NewWorkQueue(settings, s.Executor)
	if err != nil {
		return err
	}
	if len(workQueue.Queue) == 0 {
		klog.Infof("There is no work need to do")
		return nil
	}
	numWorkers := int(settings.Workers)
	if numWorkers > WorkersLimit {
		numWorkers = WorkersLimit
	}
	if numWorkers > len(workQueue.Queue) {
		numWorkers = len(workQueue.Queue)
	}
	wg := &sync.WaitGroup{}
	s.wg = wg
	s.jobListener = make(chan *DelayedJob)
	s.workerPool = make(chan *Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		s.workerPool <- &Worker{wg}
	}
	go func() {
		for job := range s.jobListener {
			time.Sleep(job.delayPeriod)
			worker := <-s.workerPool
			worker.wg.Add(1)
			go func(job Job, worker *Worker) {
				defer worker.wg.Done()
				results, _ := job.Do(s.Executor, s.dryrun)
				s.lock.Lock()
				s.results = append(s.results, results...)
				s.lock.Unlock()

				s.workerPool <- worker // return it back to the worker pool
			}(job, worker)
		}
	}()
	for _, items := range workQueue.Queue {
		s.jobListener <- &DelayedJob{Job: (WorkItems)(items)}
	}
	s.wg.Wait()          // wait for all worker to finish their jobs
	close(s.jobListener) // stop job dispatching loop

	return nil
}

func (s *FioServer) printResults(outputFile, format string) {
	t := table.NewWriter()
	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			klog.Warningf("Failed to open file %s: %s", outputFile, err)
			t.SetOutputMirror(os.Stdout)
		} else {
			defer f.Close()
			t.SetOutputMirror(f)
		}
	} else {
		t.SetOutputMirror(os.Stdout)
	}
	t.AppendHeader(table.Row{ //"job",
		"filename", "rw", "numjobs", "runtime", "direct", "blocksize", "iodepth",
		"read-iops-mean", "read-bw-mean(KiB/s)", "latency-read-min(us)", "latency-read-max(us)", "latency-read-mean(us)", "read-stddev(us)",
		"write-iops-mean", "write-bw-mean(KiB/s)", "latency-write-min(us)", "latency-write-max(us)", "latency-write-mean(us)", "latency-write-stddev(us)",
		"ioengine", "verify"})
	for _, result := range s.results {
		for _, job := range result.Jobs {
			t.AppendRow(table.Row{
				// job.JobName,
				job.JobOptions.FileName,
				job.JobOptions.RW,
				job.JobOptions.NumJobs,
				job.JobOptions.Runtime,
				job.JobOptions.Direct,
				job.JobOptions.BlockSize,
				job.JobOptions.IODepth,
				job.ReadResult.IOPSMean,
				job.ReadResult.BWMean,
				job.ReadResult.LatencyNs.Min / 1000,
				job.ReadResult.LatencyNs.Max / 1000,
				job.ReadResult.LatencyNs.Mean / 1000,
				job.ReadResult.LatencyNs.Stddev / 1000,
				job.WriteResult.IOPSMean,
				job.WriteResult.BWMean,
				job.WriteResult.LatencyNs.Min / 1000,
				job.WriteResult.LatencyNs.Max / 1000,
				job.WriteResult.LatencyNs.Mean / 1000,
				job.WriteResult.LatencyNs.Stddev / 1000,
				job.JobOptions.IOEngine,
				job.JobOptions.Verify,
			})
		}
		t.AppendSeparator()
	}
	t.SortBy([]table.SortBy{
		{
			Name: "filename",
			Mode: table.Asc,
		},
		{
			Name: "numjobs",
			Mode: table.Asc,
		},
		{
			Name: "iodepth",
			Mode: table.Asc,
		},
		{
			Name: "rw",
			Mode: table.Asc,
		},
		{
			Name: "blocksize",
			Mode: table.Asc,
		},
	})
	switch strings.ToLower(format) {
	case "md", "markdown":
		t.RenderMarkdown()
	case "csv":
		t.RenderCSV()
	case "html":
		t.RenderHTML()
	default:
		t.Render()
	}
}

func (s *FioServer) Close() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}

// Worker is registered in a worker pool waiting for jobs and execute them.
type Worker struct {
	wg *sync.WaitGroup
}

// Job defines a task, which is given to a dispatcher to be executed
// by a worker with a separate goroutine
type Job interface {
	Do(context exec.Executor, dryrun bool) ([]*client.FioResult, error)
}

type DelayedJob struct {
	Job
	delayPeriod time.Duration
}
