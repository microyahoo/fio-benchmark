package sys

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	perrors "github.com/pkg/errors"
	"k8s.io/klog/v2"

	"github.com/microyahoo/fio-benchmark/pkg/util/exec"
)

const (
	DiskByPath     = "/dev/disk/by-path/"
	SmdInContainer = "SMD_IN_CONTAINER"
)

var (
	isRBD = regexp.MustCompile("^(?:/dev/)?rbd[0-9]+p?[0-9]{0,}$")
)

func DevDiskByPath(path string) string {
	return filepath.Join(DiskByPath, path)
}

func supportedDeviceType(device string) bool {
	return device == DiskType ||
		device == SSDType ||
		device == CryptType ||
		device == LVMType ||
		device == MultiPath ||
		device == PartType ||
		device == LinearType
}

// GetDeviceEmpty check whether a device is completely empty
func GetDeviceEmpty(device *LocalDevice) bool {
	return device.Parent == "" && supportedDeviceType(device.Type) && len(device.Partitions) == 0 && device.Filesystem == ""
}

func ignoreDevice(d string) bool {
	return isRBD.MatchString(d)
}

// DiscoverDevices returns all the details of devices available on the local node
func DiscoverDevices(executor exec.Executor) (map[string]*LocalDevice, error) {
	output, err := executor.ExecuteCommandWithOutput("lsblk", "--all", "--bytes", "--pairs",
		"--paths", "--output", "SIZE,ROTA,RO,TYPE,PKNAME,NAME,KNAME,UUID,WWN,MOUNTPOINT")
	if err != nil {
		klog.Errorf("failed to execute lsblk with error: %s output: %s", err, output)
		return nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	deviceProps := make(map[string]map[string]string) // name -> map[key]value
	for scanner.Scan() {
		props := parseKeyValuePairString(scanner.Text())
		// NOTE: name maybe prefix with /dev/ or /dev/mapper/
		if len(props) > 0 {
			deviceProps[props["NAME"]] = props
		}
	}
	if scanner.Err() != nil {
		return nil, perrors.Wrapf(scanner.Err(), "failed to scan through lsblk")
	}

	var disks = make(map[string]*LocalDevice)

	for name, d := range deviceProps {
		// Ignore RBD device
		if ignoreDevice(name) {
			// skip device
			klog.Warningf("skipping rbd device %q", name)
			continue
		}

		// Populate device information coming from lsblk
		disk, err := PopulateDeviceInfo(d)
		if err != nil {
			klog.Warningf("skipping device %q. %v", name, err)
			continue
		}

		// Populate udev information coming from udev
		disk, err = PopulateDeviceUdevInfo(executor, name, disk)
		if err != nil {
			// go on without udev info
			// not ideal for our filesystem check later but we can't really fail either...
			klog.Warningf("failed to get udev info for device %q. %v", name, err)
		}

		if disk.Type == DiskType {
			deviceChild, err := ListDevicesChild(executor, name)
			if err != nil {
				klog.Warningf("failed to detect child devices for device %q, assuming they are none. %v", name, err)
			}
			// lsblk will output at least 2 lines if they are partitions, one for the parent
			// and N for the child
			if len(deviceChild) > 1 {
				disk.HasChildren = true
			}
			partitions, _, err := GetDevicePartitions(executor, name)
			if err != nil {
				klog.Warningf("failed to detect child partitions for device %q, assuming they are none. %v", name, err)
			}
			if len(partitions) > 0 {
				disk.Partitions = partitions
			}
			disk.DeviceClass = GetDiskDeviceClass(disk)
		}
		disk.Empty = GetDeviceEmpty(disk)

		disks[name] = disk
	}
	for _, d := range disks {
		for d.IsRoot && d.Parent != "" {
			d = disks[d.Parent]
			d.IsRoot = true
		}
	}
	klog.V(5).Infof("discovered disks are:")
	for _, disk := range disks {
		klog.V(5).Infof("%+v", disk)
	}

	return disks, nil
}

// PopulateDeviceInfo returns the information of the specified block device
func PopulateDeviceInfo(diskProps map[string]string) (*LocalDevice, error) {
	if diskProps == nil {
		return nil, errors.New("disk properties is empty")
	}

	diskType, ok := diskProps["TYPE"]
	if !ok {
		return nil, errors.New("diskType is empty")
	}
	if !supportedDeviceType(diskType) {
		return nil, fmt.Errorf("unsupported diskType %+s", diskType)
	}

	disk := &LocalDevice{Name: diskProps["NAME"]}

	if val, ok := diskProps["UUID"]; ok {
		disk.UUID = val
	}
	if val, ok := diskProps["MOUNTPOINT"]; ok {
		disk.MountPoint = val
		if val == SystemRootPath ||
			os.Getenv(SmdInContainer) == "true" && val == SystemRootfsPath {
			disk.IsRoot = true
		}
	}

	if val, ok := diskProps["TYPE"]; ok {
		disk.Type = val
	}
	if val, ok := diskProps["SIZE"]; ok {
		if size, err := strconv.ParseUint(val, 10, 64); err == nil {
			disk.Size = size
		}
	}
	if val, ok := diskProps["ROTA"]; ok {
		if rotates, err := strconv.ParseBool(val); err == nil {
			disk.Rotational = rotates
		}
	}
	if val, ok := diskProps["RO"]; ok {
		if ro, err := strconv.ParseBool(val); err == nil {
			disk.Readonly = ro
		}
	}
	if val, ok := diskProps["PKNAME"]; ok {
		if val != "" {
			disk.Parent = val
		}
	}
	if val, ok := diskProps["NAME"]; ok {
		disk.RealPath = val
	}
	if val, ok := diskProps["KNAME"]; ok {
		disk.KernelName = val
	}

	return disk, nil
}

// PopulateDeviceUdevInfo fills the udev info into the block device information
func PopulateDeviceUdevInfo(executor exec.Executor, device string, disk *LocalDevice) (*LocalDevice, error) {
	udevInfo, err := GetUdevInfo(executor, device)
	if err != nil {
		return disk, err
	}
	// parse udev info output
	if val, ok := udevInfo["DEVLINKS"]; ok {
		disk.DevLinks = val
	}
	if val, ok := udevInfo["ID_FS_TYPE"]; ok {
		disk.Filesystem = val
	}
	if val, ok := udevInfo["ID_SERIAL"]; ok {
		disk.Serial = val
	}
	if val, ok := udevInfo["ID_BUS"]; ok {
		disk.Bus = val
	}
	if val, ok := udevInfo["ID_VENDOR"]; ok {
		disk.Vendor = val
	}
	if val, ok := udevInfo["ID_MODEL"]; ok {
		disk.Model = val
	}
	if val, ok := udevInfo["ID_WWN_WITH_EXTENSION"]; ok {
		disk.WWNVendorExtension = val
	}
	if val, ok := udevInfo["ID_WWN"]; ok {
		disk.WWN = val
	}
	if val, ok := udevInfo["ID_PATH"]; ok {
		disk.PathID = val
	}
	return disk, nil
}
