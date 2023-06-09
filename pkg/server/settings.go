package server

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type TestSettings struct {
	FioSettings *FioSettings `yaml:"fio_settings"`
	UseAllDisks bool         `yaml:"use_all_disks"` // except root disk
	Workers     int32        `yaml:"workers"`
}

type FioSettings struct {
	NumJobs   []int32  `yaml:"numjobs"`  // 1 2 4 8 16 32 64 128 256 512 1024 2048
	IOEngine  string   `yaml:"ioengine"` // libaio
	Direct    bool     `yaml:"direct"`   // direct io
	Verify    bool     `yaml:"verify"`
	BlockSize []string `yaml:"bs"`       // 4K, 8K, 16K, 32K, 256K, 512K, 1M, 4M
	Runtime   uint64   `yaml:"runtime"`  // seconds
	IODepth   []int32  `yaml:"iodepth"`  // 1, 2, 4, 8, 16, 32, 64, 128
	RW        []string `yaml:"rw"`       // read, write, randread, randwrite, rw, randrw
	FileName  []string `yaml:"filename"` // device name or file name, which can be ignore if specify `use_all_disk`
}

func ParseSettings(cfgFile string) (*TestSettings, error) {
	out, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}
	settings := TestSettings{}
	err = yaml.Unmarshal(out, &settings)
	if err != nil {
		return nil, err
	}
	if !settings.UseAllDisks && len(settings.FioSettings.FileName) == 0 {
		return nil, errors.Errorf("filename or userAllDisks should be specified")
	}
	if settings.FioSettings == nil {
		return nil, errors.Errorf("fio parameters should be specified")
	}
	if settings.Workers <= 0 {
		settings.Workers = 1
	}
	return &settings, nil
}
