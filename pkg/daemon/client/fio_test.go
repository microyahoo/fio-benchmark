package client

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"k8s.io/klog/v2"

	exectest "github.com/microyahoo/fio-benchmark/pkg/util/exec/test"
)

func TestFioSuite(t *testing.T) {
	suite.Run(t, new(fioTestSuite))
}

type fioTestSuite struct {
	suite.Suite
}

func (s *fioTestSuite) TestFioJobs() {
	output := `
{
  "fio version" : "fio-3.27",
  "timestamp" : 1685782697,
  "timestamp_ms" : 1685782697598,
  "time" : "Sat Jun  3 16:58:17 2023",
  "jobs" : [
    {
      "jobname" : "write_throughput",
      "groupid" : 0,
      "error" : 0,
      "eta" : 0,
      "elapsed" : 101,
      "job options" : {
        "name" : "write_throughput",
        "filename" : "/dev/vdb",
        "numjobs" : "8",
        "runtime" : "100s",
        "ioengine" : "libaio",
        "direct" : "1",
        "verify" : "0",
        "bs" : "4K",
        "iodepth" : "1",
        "rw" : "randwrite",
        "group_reporting" : "1"
      },
      "read" : {
        "io_bytes" : 0,
        "io_kbytes" : 0,
        "bw_bytes" : 0,
        "bw" : 0,
        "iops" : 0.000000,
        "runtime" : 0,
        "total_ios" : 0,
        "short_ios" : 0,
        "drop_ios" : 0,
        "slat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "clat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "lat_ns" : {
          "min" : 10,
          "max" : 30,
          "mean" : 9.000000,
          "stddev" : 3.000000,
          "N" : 0
        },
        "bw_min" : 0,
        "bw_max" : 0,
        "bw_agg" : 0.000000,
        "bw_mean" : 20.000000,
        "bw_dev" : 0.000000,
        "bw_samples" : 0,
        "iops_min" : 0,
        "iops_max" : 0,
        "iops_mean" : 10.000000,
        "iops_stddev" : 0.000000,
        "iops_samples" : 0
      },
      "write" : {
        "io_bytes" : 937209856,
        "io_kbytes" : 915244,
        "bw_bytes" : 9371817,
        "bw" : 9152,
        "iops" : 2288.041359,
        "runtime" : 100003,
        "total_ios" : 228811,
        "short_ios" : 0,
        "drop_ios" : 0,
        "slat_ns" : {
          "min" : 5869,
          "max" : 5414297,
          "mean" : 19195.214041,
          "stddev" : 15994.569417,
          "N" : 228811
        },
        "clat_ns" : {
          "min" : 608761,
          "max" : 68187524,
          "mean" : 3468368.231226,
          "stddev" : 1721234.091105,
          "N" : 228811,
          "percentile" : {
            "1.000000" : 1187840,
            "5.000000" : 1515520,
            "10.000000" : 1728512,
            "20.000000" : 2113536,
            "30.000000" : 2506752,
            "40.000000" : 2801664,
            "50.000000" : 3129344,
            "60.000000" : 3489792,
            "70.000000" : 3948544,
            "80.000000" : 4554752,
            "90.000000" : 5537792,
            "95.000000" : 6520832,
            "99.000000" : 8978432,
            "99.500000" : 10420224,
            "99.900000" : 15138816,
            "99.950000" : 17956864,
            "99.990000" : 27918336
          }
        },
        "lat_ns" : {
          "min" : 624788,
          "max" : 68213304,
          "mean" : 3488600.241282,
          "stddev" : 1721929.434940,
          "N" : 228811
        },
        "bw_min" : 6100,
        "bw_max" : 12824,
        "bw_agg" : 100.000000,
        "bw_mean" : 9157.989950,
        "bw_dev" : 109.078221,
        "bw_samples" : 1592,
        "iops_min" : 1522,
        "iops_max" : 3206,
        "iops_mean" : 2289.386935,
        "iops_stddev" : 27.290962,
        "iops_samples" : 1592
      },
      "trim" : {
        "io_bytes" : 0,
        "io_kbytes" : 0,
        "bw_bytes" : 0,
        "bw" : 0,
        "iops" : 0.000000,
        "runtime" : 0,
        "total_ios" : 0,
        "short_ios" : 0,
        "drop_ios" : 0,
        "slat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "clat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "lat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        },
        "bw_min" : 0,
        "bw_max" : 0,
        "bw_agg" : 0.000000,
        "bw_mean" : 0.000000,
        "bw_dev" : 0.000000,
        "bw_samples" : 0,
        "iops_min" : 0,
        "iops_max" : 0,
        "iops_mean" : 0.000000,
        "iops_stddev" : 0.000000,
        "iops_samples" : 0
      },
      "sync" : {
        "total_ios" : 0,
        "lat_ns" : {
          "min" : 0,
          "max" : 0,
          "mean" : 0.000000,
          "stddev" : 0.000000,
          "N" : 0
        }
      },
      "job_runtime" : 800006,
      "usr_cpu" : 0.395997,
      "sys_cpu" : 0.693620,
      "ctx" : 228874,
      "majf" : 0,
      "minf" : 108,
      "iodepth_level" : {
        "1" : 100.000000,
        "2" : 0.000000,
        "4" : 0.000000,
        "8" : 0.000000,
        "16" : 0.000000,
        "32" : 0.000000,
        ">=64" : 0.000000
      },
      "iodepth_submit" : {
        "0" : 0.000000,
        "4" : 100.000000,
        "8" : 0.000000,
        "16" : 0.000000,
        "32" : 0.000000,
        "64" : 0.000000,
        ">=64" : 0.000000
      },
      "iodepth_complete" : {
        "0" : 0.000000,
        "4" : 100.000000,
        "8" : 0.000000,
        "16" : 0.000000,
        "32" : 0.000000,
        "64" : 0.000000,
        ">=64" : 0.000000
      },
      "latency_ns" : {
        "2" : 0.000000,
        "4" : 0.000000,
        "10" : 0.000000,
        "20" : 0.000000,
        "50" : 0.000000,
        "100" : 0.000000,
        "250" : 0.000000,
        "500" : 0.000000,
        "750" : 0.000000,
        "1000" : 0.000000
      },
      "latency_us" : {
        "2" : 0.000000,
        "4" : 0.000000,
        "10" : 0.000000,
        "20" : 0.000000,
        "50" : 0.000000,
        "100" : 0.000000,
        "250" : 0.000000,
        "500" : 0.000000,
        "750" : 0.018793,
        "1000" : 0.406886
      },
      "latency_ms" : {
        "2" : 15.173222,
        "4" : 55.223307,
        "10" : 28.582105,
        "20" : 0.561162,
        "50" : 0.032778,
        "100" : 0.010000,
        "250" : 0.000000,
        "500" : 0.000000,
        "750" : 0.000000,
        "1000" : 0.000000,
        "2000" : 0.000000,
        ">=2000" : 0.000000
      },
      "latency_depth" : 1,
      "latency_target" : 0,
      "latency_percentile" : 100.000000,
      "latency_window" : 0
    }
  ],
  "disk_util" : [
    {
      "name" : "vdb",
      "read_ios" : 51,
      "write_ios" : 228550,
      "read_merges" : 0,
      "write_merges" : 0,
      "read_ticks" : 111,
      "write_ticks" : 789202,
      "in_queue" : 789313,
      "util" : 100.000000
    }
  ]
}
`
	executor := &exectest.MockExecutor{
		MockExecuteCommandWithOutput: func(command string, args ...string) (string, error) {
			klog.Infof("run command %s %v", command, args)
			return output, nil
		},
	}
	actual, err := FioTest(executor, "/dev/vdb", 8, "4K", 1, "randrw", 120, "libaio", true, true, false)
	s.NoError(err)
	s.Len(actual.Jobs, 1)
	expect := &FioResult{
		Jobs: []*FioJob{
			{
				JobName: "write_throughput",
				JobOptions: &JobOptions{
					Name:      "write_throughput",
					FileName:  "/dev/vdb",
					NumJobs:   "8",
					Runtime:   "100s",
					IOEngine:  "libaio",
					Direct:    "1",
					Verify:    "0",
					BlockSize: "4K",
					IODepth:   "1",
					RW:        "randwrite",
				},
				ReadResult: &ReadResult{
					IOPSMean: 10,
					BWMean:   20,
					LatencyNs: LatencyNs{
						Min:    10,
						Max:    30,
						Mean:   9,
						Stddev: 3,
					},
				},
				WriteResult: &WriteResult{
					IOPSMean: 2289.386935,
					BWMean:   9157.989950,
					LatencyNs: LatencyNs{
						Min:    624788,
						Max:    68213304,
						Mean:   3488600.241282,
						Stddev: 1721929.434940,
					},
				},
			},
		},
	}
	s.Assert().EqualValues(expect, actual)
}
