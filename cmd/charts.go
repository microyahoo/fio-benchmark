package cmd

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/microyahoo/fio-benchmark/pkg/daemon/client"
)

var chartsCmd = &cobra.Command{
	Use:   "generate-charts",
	Short: "Generate charts based on specified CSV file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return generate(cmd, args)
	},
	TraverseChildren: true,
}

var (
	csvFile   string
	chartFile string
)

func generate(cmd *cobra.Command, args []string) error {
	if csvFile == "" {
		return errors.New("CSV file should be specified")
	}
	f, err := os.Open(csvFile)
	if err != nil {
		return err
	}
	// filename,rw,numjobs,runtime,direct,blocksize,iodepth,read-iops-mean,read-bw-mean(KiB/s),latency-read-min(us),latency-read-max(us),latency-read-mean(us),read-stddev(us),write-iops-mean,write-bw-mean(KiB/s),latency-write-min(us),latency-write-max(us),latency-write-mean(us),latency-write-stddev(us),ioengine,verify
	// /dev/nvme0n1,randread,1,120s,1,4K,1,13369.941423,53479.774059,54,6039,74.393800365,11.089774646,0,0,0,0,0,0,libaio,
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return err
	}
	var jobs []*client.FioJob
	var numJobsMap = make(map[int32]struct{})
	for i, record := range records {
		if i == 0 {
			// skip first line
			continue
		}
		if len(record) != 21 {
			return errors.Errorf("Invalid csv format")
		}
		filename := record[0]
		rw := record[1]
		numjobs := record[2]
		j, _ := strconv.ParseInt(numjobs, 10, 32)
		numJobsMap[int32(j)] = struct{}{}
		runtime := record[3]
		direct := record[4]
		bs := record[5]
		iodepth := record[6]
		readIOPS, _ := strconv.ParseFloat(record[7], 64)
		readBw, _ := strconv.ParseFloat(record[8], 64)
		readLatMin, _ := strconv.ParseFloat(record[9], 64)
		readLatMin *= 1000
		readLatMax, _ := strconv.ParseFloat(record[10], 64)
		readLatMax *= 1000
		readLatMean, _ := strconv.ParseFloat(record[11], 64)
		readLatMean *= 1000
		readLatStddev, _ := strconv.ParseFloat(record[12], 64)
		readLatStddev *= 1000
		writeIOPS, _ := strconv.ParseFloat(record[13], 64)
		writeBw, _ := strconv.ParseFloat(record[14], 64)
		writeLatMin, _ := strconv.ParseFloat(record[15], 64)
		writeLatMin *= 1000
		writeLatMax, _ := strconv.ParseFloat(record[16], 64)
		writeLatMax *= 1000
		writeLatMean, _ := strconv.ParseFloat(record[17], 64)
		writeLatMean *= 1000
		writeLatStddev, _ := strconv.ParseFloat(record[18], 64)
		writeLatStddev *= 1000
		ioengine := record[19]
		verify := record[20]
		job := &client.FioJob{
			JobName: filename,
			JobOptions: &client.JobOptions{
				FileName:  filename,
				NumJobs:   numjobs,
				Runtime:   runtime,
				IOEngine:  ioengine,
				Direct:    direct,
				Verify:    verify,
				BlockSize: bs,
				IODepth:   iodepth,
				RW:        rw,
			},
			ReadResult: &client.ReadResult{
				IOPSMean: readIOPS,
				BWMean:   readBw,
				LatencyNs: client.LatencyNs{
					Min:    readLatMin,
					Max:    readLatMax,
					Mean:   readLatMean,
					Stddev: readLatStddev,
				},
			},
			WriteResult: &client.WriteResult{
				IOPSMean: writeIOPS,
				BWMean:   writeBw,
				LatencyNs: client.LatencyNs{
					Min:    writeLatMin,
					Max:    writeLatMax,
					Mean:   writeLatMean,
					Stddev: writeLatStddev,
				},
			},
		}
		jobs = append(jobs, job)
	}
	var results = []*client.FioResult{
		{
			Jobs: jobs,
		},
	}
	var numJobs []int32
	for j := range numJobsMap {
		numJobs = append(numJobs, j)
	}
	err = client.RenderCharts(results, numJobs, chartFile)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	chartsCmd.Flags().StringVar(&csvFile, "csv-file", "", "CSV file you want to generate chart")
	chartsCmd.Flags().StringVar(&chartFile, "chart-file", "", "chart file you want to generate")
}
