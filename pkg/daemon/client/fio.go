package client

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/google/uuid"
	"k8s.io/klog/v2"

	"github.com/microyahoo/fio-benchmark/pkg/util/exec"
)

// fio --name=write_throughput --filename=/dev/vdb --numjobs=8 --time_based --runtime=100s --ioengine=libaio --direct=1 --verify=0 --bs=4K --iodepth=1 --rw=randwrite --group_reporting=1
func FioTest(executor exec.Executor, filename string, numJobs int32, bs string, iodepth int32, rw string, runtime uint64, ioengine string, verify, direct, dryrun bool) (*FioResult, error) {
	name := fmt.Sprintf("%s-%s", rw, uuid.NewString())
	if ioengine == "" {
		ioengine = "libaio"
	}
	d := "1"
	if !direct {
		d = "0"
	}
	args := []string{
		"--name", name,
		"--filename", filename,
		"--numjobs", fmt.Sprintf("%d", numJobs),
		"--time_based",
		"--ioengine", ioengine,
		"--bs", bs,
		"--rw", rw,
		"--direct", d,
		"--group_reporting",
		"--iodepth", fmt.Sprintf("%d", iodepth),
		"--runtime", fmt.Sprintf("%ds", runtime),
		"--output-format", "json"}
	if !verify {
		args = append(args, "--verify", "0")
	}
	if dryrun {
		klog.Infof("Running command: %s %s", FioTool, strings.Join(args, " "))
		return nil, nil
	}
	output, err := executor.ExecuteCommandWithOutput(FioTool, args...)
	if err != nil {
		return nil, err
	}
	var r *FioResult
	err = json.Unmarshal([]byte(output), &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func DropCaches(executor exec.Executor) error {
	// echo	3 > /proc/sys/vm/drop_caches to clear PageCache, dentries and inodes
	data := []byte("3")
	err := os.WriteFile("/proc/sys/vm/drop_caches", data, 0)
	if err != nil {
		return err
	}
	return nil
}

//	{
//	  "fio version" : "fio-3.27",
//	  "timestamp" : 1685782697,
//	  "timestamp_ms" : 1685782697598,
//	  "time" : "Sat Jun  3 16:58:17 2023",
//	  "jobs" : [
//	    {
//	      "jobname" : "write_throughput",
//	      "groupid" : 0,
//	      "error" : 0,
//	      "eta" : 0,
//	      "elapsed" : 101,
//	      "job options" : {
//	        "name" : "write_throughput",
//	        "filename" : "/dev/vdb",
//	        "numjobs" : "8",
//	        "runtime" : "100s",
//	        "ioengine" : "libaio",
//	        "direct" : "1",
//	        "verify" : "0",
//	        "bs" : "4K",
//	        "iodepth" : "1",
//	        "rw" : "randwrite",
//	        "group_reporting" : "1"
//	      },
//	      "read" : {
//	        "io_bytes" : 0,
//	        "io_kbytes" : 0,
//	        "bw_bytes" : 0,
//	        "bw" : 0,
//	        "iops" : 0.000000,
//	        "runtime" : 0,
//	        "total_ios" : 0,
//	        "short_ios" : 0,
//	        "drop_ios" : 0,
//	        "slat_ns" : {
//	          "min" : 0,
//	          "max" : 0,
//	          "mean" : 0.000000,
//	          "stddev" : 0.000000,
//	          "N" : 0
//	        },
//	        "clat_ns" : {
//	          "min" : 0,
//	          "max" : 0,
//	          "mean" : 0.000000,
//	          "stddev" : 0.000000,
//	          "N" : 0
//	        },
//	        "lat_ns" : {
//	          "min" : 0,
//	          "max" : 0,
//	          "mean" : 0.000000,
//	          "stddev" : 0.000000,
//	          "N" : 0
//	        },
//	        "bw_min" : 0,
//	        "bw_max" : 0,
//	        "bw_agg" : 0.000000,
//	        "bw_mean" : 0.000000,
//	        "bw_dev" : 0.000000,
//	        "bw_samples" : 0,
//	        "iops_min" : 0,
//	        "iops_max" : 0,
//	        "iops_mean" : 0.000000,
//	        "iops_stddev" : 0.000000,
//	        "iops_samples" : 0
//	      },
//	      "write" : {
//	        "io_bytes" : 937209856,
//	        "io_kbytes" : 915244,
//	        "bw_bytes" : 9371817,
//	        "bw" : 9152,
//	        "iops" : 2288.041359,
//	        "runtime" : 100003,
//	        "total_ios" : 228811,
//	        "short_ios" : 0,
//	        "drop_ios" : 0,
//	        "slat_ns" : {
//	          "min" : 5869,
//	          "max" : 5414297,
//	          "mean" : 19195.214041,
//	          "stddev" : 15994.569417,
//	          "N" : 228811
//	        },
//	        "clat_ns" : {
//	          "min" : 608761,
//	          "max" : 68187524,
//	          "mean" : 3468368.231226,
//	          "stddev" : 1721234.091105,
//	          "N" : 228811,
//	          "percentile" : {
//	            "1.000000" : 1187840,
//	            "5.000000" : 1515520,
//	            "10.000000" : 1728512,
//	            "20.000000" : 2113536,
//	            "30.000000" : 2506752,
//	            "40.000000" : 2801664,
//	            "50.000000" : 3129344,
//	            "60.000000" : 3489792,
//	            "70.000000" : 3948544,
//	            "80.000000" : 4554752,
//	            "90.000000" : 5537792,
//	            "95.000000" : 6520832,
//	            "99.000000" : 8978432,
//	            "99.500000" : 10420224,
//	            "99.900000" : 15138816,
//	            "99.950000" : 17956864,
//	            "99.990000" : 27918336
//	          }
//	        },
//	        "lat_ns" : {
//	          "min" : 624788,
//	          "max" : 68213304,
//	          "mean" : 3488600.241282,
//	          "stddev" : 1721929.434940,
//	          "N" : 228811
//	        },
//	        "bw_min" : 6100,
//	        "bw_max" : 12824,
//	        "bw_agg" : 100.000000,
//	        "bw_mean" : 9157.989950,
//	        "bw_dev" : 109.078221,
//	        "bw_samples" : 1592,
//	        "iops_min" : 1522,
//	        "iops_max" : 3206,
//	        "iops_mean" : 2289.386935,
//	        "iops_stddev" : 27.290962,
//	        "iops_samples" : 1592
//	      },
//	      "trim" : {
//	        "io_bytes" : 0,
//	        "io_kbytes" : 0,
//	        "bw_bytes" : 0,
//	        "bw" : 0,
//	        "iops" : 0.000000,
//	        "runtime" : 0,
//	        "total_ios" : 0,
//	        "short_ios" : 0,
//	        "drop_ios" : 0,
//	        "slat_ns" : {
//	          "min" : 0,
//	          "max" : 0,
//	          "mean" : 0.000000,
//	          "stddev" : 0.000000,
//	          "N" : 0
//	        },
//	        "clat_ns" : {
//	          "min" : 0,
//	          "max" : 0,
//	          "mean" : 0.000000,
//	          "stddev" : 0.000000,
//	          "N" : 0
//	        },
//	        "lat_ns" : {
//	          "min" : 0,
//	          "max" : 0,
//	          "mean" : 0.000000,
//	          "stddev" : 0.000000,
//	          "N" : 0
//	        },
//	        "bw_min" : 0,
//	        "bw_max" : 0,
//	        "bw_agg" : 0.000000,
//	        "bw_mean" : 0.000000,
//	        "bw_dev" : 0.000000,
//	        "bw_samples" : 0,
//	        "iops_min" : 0,
//	        "iops_max" : 0,
//	        "iops_mean" : 0.000000,
//	        "iops_stddev" : 0.000000,
//	        "iops_samples" : 0
//	      },
//	      "sync" : {
//	        "total_ios" : 0,
//	        "lat_ns" : {
//	          "min" : 0,
//	          "max" : 0,
//	          "mean" : 0.000000,
//	          "stddev" : 0.000000,
//	          "N" : 0
//	        }
//	      },
//	      "job_runtime" : 800006,
//	      "usr_cpu" : 0.395997,
//	      "sys_cpu" : 0.693620,
//	      "ctx" : 228874,
//	      "majf" : 0,
//	      "minf" : 108,
//	      "iodepth_level" : {
//	        "1" : 100.000000,
//	        "2" : 0.000000,
//	        "4" : 0.000000,
//	        "8" : 0.000000,
//	        "16" : 0.000000,
//	        "32" : 0.000000,
//	        ">=64" : 0.000000
//	      },
//	      "iodepth_submit" : {
//	        "0" : 0.000000,
//	        "4" : 100.000000,
//	        "8" : 0.000000,
//	        "16" : 0.000000,
//	        "32" : 0.000000,
//	        "64" : 0.000000,
//	        ">=64" : 0.000000
//	      },
//	      "iodepth_complete" : {
//	        "0" : 0.000000,
//	        "4" : 100.000000,
//	        "8" : 0.000000,
//	        "16" : 0.000000,
//	        "32" : 0.000000,
//	        "64" : 0.000000,
//	        ">=64" : 0.000000
//	      },
//	      "latency_ns" : {
//	        "2" : 0.000000,
//	        "4" : 0.000000,
//	        "10" : 0.000000,
//	        "20" : 0.000000,
//	        "50" : 0.000000,
//	        "100" : 0.000000,
//	        "250" : 0.000000,
//	        "500" : 0.000000,
//	        "750" : 0.000000,
//	        "1000" : 0.000000
//	      },
//	      "latency_us" : {
//	        "2" : 0.000000,
//	        "4" : 0.000000,
//	        "10" : 0.000000,
//	        "20" : 0.000000,
//	        "50" : 0.000000,
//	        "100" : 0.000000,
//	        "250" : 0.000000,
//	        "500" : 0.000000,
//	        "750" : 0.018793,
//	        "1000" : 0.406886
//	      },
//	      "latency_ms" : {
//	        "2" : 15.173222,
//	        "4" : 55.223307,
//	        "10" : 28.582105,
//	        "20" : 0.561162,
//	        "50" : 0.032778,
//	        "100" : 0.010000,
//	        "250" : 0.000000,
//	        "500" : 0.000000,
//	        "750" : 0.000000,
//	        "1000" : 0.000000,
//	        "2000" : 0.000000,
//	        ">=2000" : 0.000000
//	      },
//	      "latency_depth" : 1,
//	      "latency_target" : 0,
//	      "latency_percentile" : 100.000000,
//	      "latency_window" : 0
//	    }
//	  ],
//	  "disk_util" : [
//	    {
//	      "name" : "vdb",
//	      "read_ios" : 51,
//	      "write_ios" : 228550,
//	      "read_merges" : 0,
//	      "write_merges" : 0,
//	      "read_ticks" : 111,
//	      "write_ticks" : 789202,
//	      "in_queue" : 789313,
//	      "util" : 100.000000
//	    }
//	  ]
//	}
type FioResult struct {
	Jobs []*FioJob `json:"jobs"`
}

type FioJob struct {
	JobName     string       `json:"jobname"`
	JobOptions  *JobOptions  `json:"job options"`
	ReadResult  *ReadResult  `json:"read"`
	WriteResult *WriteResult `json:"write"`
}

type JobOptions struct {
	Name      string `json:"name"`
	FileName  string `json:"filename"`
	NumJobs   string `json:"numjobs"`
	Runtime   string `json:"runtime"`
	IOEngine  string `json:"ioengine"`
	Direct    string `json:"direct"`
	Verify    string `json:"verify"`
	BlockSize string `json:"bs"`
	IODepth   string `json:"iodepth"`
	RW        string `json:"rw"`
}

type ReadResult struct {
	IOPSMean  float64   `json:"iops_mean"`
	BWMean    float64   `json:"bw_mean"`
	LatencyNs LatencyNs `json:"lat_ns"`
	// IOKBytes  uint64
	// BWBytes   uint64
	// IOPS      uint64
}

type LatencyNs struct {
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Mean   float64 `json:"mean"`
	Stddev float64 `json:"stddev"`
}

type WriteResult struct {
	IOPSMean  float64   `json:"iops_mean"`
	BWMean    float64   `json:"bw_mean"`
	LatencyNs LatencyNs `json:"lat_ns"`
}

func RenderCharts(results []*FioResult, numJobs []int32, chartFile string) error {
	var jobMap = make(map[string]map[string]map[string]map[string][]*FioJob) // map[rw][iodepth][bs][numjobs] => []Job
	for _, result := range results {
		for _, job := range result.Jobs {
			if _, ok1 := jobMap[job.JobOptions.RW]; !ok1 {
				jobMap[job.JobOptions.RW] = make(map[string]map[string]map[string][]*FioJob)
				jobMap[job.JobOptions.RW][job.JobOptions.IODepth] = make(map[string]map[string][]*FioJob)
				jobMap[job.JobOptions.RW][job.JobOptions.IODepth][job.JobOptions.BlockSize] = make(map[string][]*FioJob)
			} else {
				if _, ok2 := jobMap[job.JobOptions.RW][job.JobOptions.IODepth]; !ok2 {
					jobMap[job.JobOptions.RW][job.JobOptions.IODepth] = make(map[string]map[string][]*FioJob)
					jobMap[job.JobOptions.RW][job.JobOptions.IODepth][job.JobOptions.BlockSize] = make(map[string][]*FioJob)
				} else {
					if _, ok3 := jobMap[job.JobOptions.RW][job.JobOptions.IODepth][job.JobOptions.BlockSize]; !ok3 {
						jobMap[job.JobOptions.RW][job.JobOptions.IODepth][job.JobOptions.BlockSize] = make(map[string][]*FioJob)
					}
				}
			}
			jobMap[job.JobOptions.RW][job.JobOptions.IODepth][job.JobOptions.BlockSize][job.JobOptions.NumJobs] = append(
				jobMap[job.JobOptions.RW][job.JobOptions.IODepth][job.JobOptions.BlockSize][job.JobOptions.NumJobs], job)
		}
	}

	sort.Slice(numJobs, func(i, j int) bool { return numJobs[i] < numJobs[j] })

	createLine := func(title string, yaxis string) *charts.Line {
		line := charts.NewLine()
		line.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title: title,
			}),
			charts.WithTooltipOpts(opts.Tooltip{Show: true, TriggerOn: "mousemove|click"}),
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
				readLatLine := createLine(fmt.Sprintf("readlat-%s-%s-%s", rw, bs, iodepth), "latency(ms)")
				writeLatLine := createLine(fmt.Sprintf("writelat-%s-%s-%s", rw, bs, iodepth), "latency(ms)")

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
	if chartFile == "" {
		chartFile = fmt.Sprintf("chart-%s.html", uuid.NewString())
	} else if !strings.HasSuffix(chartFile, "html") {
		chartFile = fmt.Sprintf("%s.html", chartFile)
	}
	f, err := os.Create(chartFile)
	if err != nil {
		return err
	}
	err = page.Render(f)
	if err != nil {
		return err
	}
	return nil
}

func Render3DCharts(results []*FioResult, chartFile string) error {
	var jobsMap = make(map[string]map[string]map[string]map[string]map[string]*FioJob) // map[rw][bs][iodepth][numjobs][filename]=> Job
	for _, result := range results {
		for _, job := range result.Jobs {
			if _, ok1 := jobsMap[job.JobOptions.RW]; !ok1 {
				jobsMap[job.JobOptions.RW] = make(map[string]map[string]map[string]map[string]*FioJob)
				jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize] = make(map[string]map[string]map[string]*FioJob)
				jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize][job.JobOptions.IODepth] = make(map[string]map[string]*FioJob)
				jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize][job.JobOptions.IODepth][job.JobOptions.NumJobs] = make(map[string]*FioJob)
			} else {
				if _, ok2 := jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize]; !ok2 {
					jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize] = make(map[string]map[string]map[string]*FioJob)
					jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize][job.JobOptions.IODepth] = make(map[string]map[string]*FioJob)
					jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize][job.JobOptions.IODepth][job.JobOptions.NumJobs] = make(map[string]*FioJob)
				} else {
					if _, ok3 := jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize][job.JobOptions.IODepth]; !ok3 {
						jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize][job.JobOptions.IODepth] = make(map[string]map[string]*FioJob)
						jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize][job.JobOptions.IODepth][job.JobOptions.NumJobs] = make(map[string]*FioJob)
					} else {
						if _, ok4 := jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize][job.JobOptions.IODepth][job.JobOptions.NumJobs]; !ok4 {
							jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize][job.JobOptions.IODepth][job.JobOptions.NumJobs] = make(map[string]*FioJob)
						}
					}
				}
			}
			jobsMap[job.JobOptions.RW][job.JobOptions.BlockSize][job.JobOptions.IODepth][job.JobOptions.NumJobs][job.JobOptions.FileName] = job
		}
	}

	createLine := func(title string, zaxis string) *charts.Line3D {
		line := charts.NewLine3D()
		line.SetGlobalOptions(
			charts.WithTitleOpts(opts.Title{
				Title: title,
			}),
			charts.WithTooltipOpts(opts.Tooltip{Show: true, TriggerOn: "mousemove|click"}),
			charts.WithLegendOpts(opts.Legend{Show: true, Width: "50%", Left: "right"}),
			charts.WithInitializationOpts(opts.Initialization{
				Theme: "shine",
			}),
			charts.WithVisualMapOpts(opts.VisualMap{
				Calculable: true,
			}),
			charts.WithXAxisOpts(opts.XAxis{
				Name: "num_jobs",
			}),
			charts.WithYAxisOpts(opts.YAxis{
				Name: "iodepth",
			}),
			// charts.WithZAxisOpts(opts.ZAxis{
			// 	Name: "iodepth",
			// }),
		)
		// line.SetXAxis(numJobs)
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
	generateLines := func(lines []*charts.Line3D, metricsMap map[string]map[string]map[string]*metrics) { // filename -> iodepth -> numJobs -> metrics
		for filename, iodepthMap := range metricsMap {
			var (
				readIOPSLineData  []opts.Chart3DData
				writeIOPSLineData []opts.Chart3DData
				readBwLineData    []opts.Chart3DData
				writeBwLineData   []opts.Chart3DData
				readLatLineData   []opts.Chart3DData
				writeLatLineData  []opts.Chart3DData
				lineData          [][]opts.Chart3DData
			)
			for iodepth, numJobsMap := range iodepthMap {
				for numJobs, metrics := range numJobsMap {
					readIOPSLineData = append(readIOPSLineData, opts.Chart3DData{Value: []interface{}{numJobs, iodepth, metrics.readIOPS}})
					writeIOPSLineData = append(writeIOPSLineData, opts.Chart3DData{Value: []interface{}{numJobs, iodepth, metrics.writeIOPS}})
					readBwLineData = append(readBwLineData, opts.Chart3DData{Value: []interface{}{numJobs, iodepth, metrics.readBw}})
					writeBwLineData = append(writeBwLineData, opts.Chart3DData{Value: []interface{}{numJobs, iodepth, metrics.writeBw}})
					readLatLineData = append(readLatLineData, opts.Chart3DData{Value: []interface{}{numJobs, iodepth, metrics.readLat}})
					writeLatLineData = append(writeLatLineData, opts.Chart3DData{Value: []interface{}{numJobs, iodepth, metrics.writeLat}})

				}
			}
			lineData = [][]opts.Chart3DData{readIOPSLineData, writeIOPSLineData, readBwLineData, writeBwLineData, readLatLineData, writeLatLineData}
			for i, line := range lines {
				line.AddSeries(filename, lineData[i])
			}
		}
	}
	page := components.NewPage()
	for rw, bsMap := range jobsMap {
		for bs, iodepthMap := range bsMap {
			readIOPSLine := createLine(fmt.Sprintf("readIOPS-%s-%s", rw, bs), "iops")
			writeIOPSLine := createLine(fmt.Sprintf("writeIOPS-%s-%s", rw, bs), "iops")
			readBwLine := createLine(fmt.Sprintf("readBW-%s-%s", rw, bs), "bandwidth(KiB/s)")
			writeBwLine := createLine(fmt.Sprintf("writeBW-%s-%s", rw, bs), "bandwidth(KiB/s)")
			readLatLine := createLine(fmt.Sprintf("readLat-%s-%s", rw, bs), "latency(ms)")
			writeLatLine := createLine(fmt.Sprintf("writeLat-%s-%s", rw, bs), "latency(ms)")

			var filenameMap = make(map[string]map[string]map[string]*metrics) // filename -> iodepth -> numJobs -> metrics
			for iodepth, numJobsMap := range iodepthMap {
				for numJobs, jobMap := range numJobsMap {
					for filename, job := range jobMap {
						if _, ok := filenameMap[filename]; !ok {
							filenameMap[filename] = make(map[string]map[string]*metrics)
							filenameMap[filename][iodepth] = make(map[string]*metrics)
						} else {
							if _, ok2 := filenameMap[filename][iodepth]; !ok2 {
								filenameMap[filename][iodepth] = make(map[string]*metrics)
							}
						}
						filenameMap[filename][iodepth][numJobs] = &metrics{
							readIOPS:  job.ReadResult.IOPSMean,
							readBw:    job.ReadResult.BWMean,
							readLat:   job.ReadResult.LatencyNs.Mean / 1000 / 1000, // ms
							writeIOPS: job.WriteResult.IOPSMean,
							writeBw:   job.WriteResult.BWMean,
							writeLat:  job.WriteResult.LatencyNs.Mean / 1000 / 1000, // ms
						}
					}
				}
			}
			generateLines([]*charts.Line3D{readIOPSLine, writeIOPSLine, readBwLine, writeBwLine, readLatLine, writeLatLine}, filenameMap)
			page.AddCharts(readIOPSLine, writeIOPSLine, readBwLine, writeBwLine, readLatLine, writeLatLine)
		}
	}
	if chartFile == "" {
		chartFile = fmt.Sprintf("3dchart-%s.html", uuid.NewString())
	} else if !strings.HasSuffix(chartFile, "html") {
		chartFile = fmt.Sprintf("%s.html", chartFile)
	}
	f, err := os.Create(chartFile)
	if err != nil {
		return err
	}
	err = page.Render(f)
	if err != nil {
		return err
	}
	return nil
}
