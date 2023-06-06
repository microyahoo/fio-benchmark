package server

import (
	"k8s.io/klog/v2"

	"github.com/microyahoo/fio-benchmark/pkg/daemon/client"
	"github.com/microyahoo/fio-benchmark/pkg/util/exec"
	"github.com/microyahoo/fio-benchmark/pkg/util/sys"
)

type WorkItem struct {
	FileName  string
	NumJobs   int32
	BlockSize string
	IODepth   int32
	RW        string
	Runtime   uint64
	Verify    bool
	Direct    bool
	IOEngine  string
}

type WorkItems []*WorkItem

func (wis WorkItems) Do(executor exec.Executor, dryrun bool) ([]*client.FioResult, error) {
	var results []*client.FioResult
	for _, wi := range wis {
		result, err := client.FioTest(executor, wi.FileName, wi.NumJobs, wi.BlockSize, wi.IODepth, wi.RW, wi.Runtime, wi.IOEngine, wi.Verify, wi.Direct, dryrun)
		if err != nil {
			klog.Warningf("Failed to do fio test: %v", err)
			continue
		}
		if result != nil {
			results = append(results, result)
		}
	}
	return results, nil
}

type WorkQueue struct {
	Queue map[string][]*WorkItem // filename -> items
}

func NewWorkQueue(s *TestSettings, executor exec.Executor) (*WorkQueue, error) {
	var queue map[string][]*WorkItem
	fs := s.FioSettings
	if fs.FileName != "" {
		queue = make(map[string][]*WorkItem, 1)
		var items []*WorkItem
		for _, job := range fs.NumJobs {
			for _, bs := range fs.BlockSize {
				for _, depth := range fs.IODepth {
					for _, rw := range fs.RW {
						item := &WorkItem{
							FileName:  fs.FileName,
							NumJobs:   job,
							BlockSize: bs,
							IODepth:   depth,
							RW:        rw,
							Runtime:   fs.Runtime,
							Verify:    fs.Verify,
							Direct:    fs.Direct,
							IOEngine:  fs.IOEngine,
						}
						items = append(items, item)
					}
				}
			}
		}
		queue[fs.FileName] = items
		return &WorkQueue{queue}, nil
	}
	if !s.UseAllDisks {
		return nil, nil
	}
	devices, err := sys.DiscoverDevices(executor)
	if err != nil {
		return nil, err
	}
	queue = make(map[string][]*WorkItem, len(devices)-1) // exclude root disk
	for _, d := range devices {
		if d.Type != sys.DiskType || d.Bus == sys.DiskBusUsb || d.IsRoot {
			continue
		}
		var items []*WorkItem
		for _, job := range fs.NumJobs {
			for _, bs := range fs.BlockSize {
				for _, depth := range fs.IODepth {
					for _, rw := range fs.RW {
						item := &WorkItem{
							FileName:  d.RealPath,
							NumJobs:   job,
							BlockSize: bs,
							IODepth:   depth,
							RW:        rw,
							Runtime:   fs.Runtime,
							Verify:    fs.Verify,
							Direct:    fs.Direct,
							IOEngine:  fs.IOEngine,
						}
						items = append(items, item)
					}
				}
			}
		}
		queue[d.RealPath] = items
	}
	return &WorkQueue{queue}, nil
}
