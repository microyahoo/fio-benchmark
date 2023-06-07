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
	return len(device.Parents) == 0 && supportedDeviceType(device.Type) && len(device.Partitions) == 0 && device.Filesystem == ""
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

	// ~/go/src/(main ✗) lsblk --all --bytes --pairs --paths --output SIZE,ROTA,RO,TYPE,PKNAME,NAME,KNAME,UUID,WWN,MOUNTPOINT | grep root
	// SIZE="207215394816" ROTA="1" RO="0" TYPE="lvm" PKNAME="/dev/vda2" NAME="/dev/mapper/centos-root" KNAME="/dev/dm-0" UUID="5e322b94-4141-4a15-ae29-4136ae9c2e15" WWN="" MOUNTPOINT="/"
	// SIZE="207215394816" ROTA="1" RO="0" TYPE="lvm" PKNAME="/dev/vda3" NAME="/dev/mapper/centos-root" KNAME="/dev/dm-0" UUID="5e322b94-4141-4a15-ae29-4136ae9c2e15" WWN="" MOUNTPOINT="/"
	// SIZE="207215394816" ROTA="1" RO="0" TYPE="lvm" PKNAME="/dev/vda4" NAME="/dev/mapper/centos-root" KNAME="/dev/dm-0" UUID="5e322b94-4141-4a15-ae29-4136ae9c2e15" WWN="" MOUNTPOINT="/"
	// SIZE="207215394816" ROTA="1" RO="0" TYPE="lvm" PKNAME="/dev/vdd1" NAME="/dev/mapper/centos-root" KNAME="/dev/dm-0" UUID="5e322b94-4141-4a15-ae29-4136ae9c2e15" WWN="" MOUNTPOINT="/"
	deviceProps := make(map[string][]map[string]string) // name -> []map[key]value
	for scanner.Scan() {
		props := parseKeyValuePairString(scanner.Text())
		// NOTE: name maybe prefix with /dev/ or /dev/mapper/
		if len(props) > 0 {
			name := props["NAME"]
			deviceProps[name] = append(deviceProps[name], props)
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
		if d.IsRoot && len(d.Parents) > 0 {
			setRoot(d, disks)
		}
	}
	klog.V(5).Infof("discovered disks are:")
	for _, disk := range disks {
		klog.V(5).Infof("%+v", disk)
	}

	return disks, nil
}

func setRoot(disk *LocalDevice, disks map[string]*LocalDevice) {
	for _, p := range disk.Parents {
		parent := disks[p]
		parent.IsRoot = true
		setRoot(parent, disks)
	}
}

// PopulateDeviceInfo returns the information of the specified block device
func PopulateDeviceInfo(props []map[string]string) (*LocalDevice, error) {
	if len(props) == 0 {
		return nil, errors.New("disk properties is empty")
	}

	var device *LocalDevice
	for i, deviceProps := range props {
		diskType, ok := deviceProps["TYPE"]
		if !ok {
			return nil, errors.New("diskType is empty")
		}
		if !supportedDeviceType(diskType) {
			return nil, fmt.Errorf("unsupported diskType %+s", diskType)
		}
		name := deviceProps["NAME"]
		if i == 0 {
			device = &LocalDevice{Name: name}

			if val, ok := deviceProps["UUID"]; ok {
				device.UUID = val
			}
			if val, ok := deviceProps["MOUNTPOINT"]; ok {
				device.MountPoint = val
				if val == SystemRootPath ||
					os.Getenv(SmdInContainer) == "true" && val == SystemRootfsPath {
					device.IsRoot = true
				}
			}

			if val, ok := deviceProps["TYPE"]; ok {
				device.Type = val
			}
			if val, ok := deviceProps["SIZE"]; ok {
				if size, err := strconv.ParseUint(val, 10, 64); err == nil {
					device.Size = size
				}
			}
			if val, ok := deviceProps["ROTA"]; ok {
				if rotates, err := strconv.ParseBool(val); err == nil {
					device.Rotational = rotates
				}
			}
			if val, ok := deviceProps["RO"]; ok {
				if ro, err := strconv.ParseBool(val); err == nil {
					device.Readonly = ro
				}
			}
			if val, ok := deviceProps["NAME"]; ok {
				device.RealPath = val
			}
			if val, ok := deviceProps["KNAME"]; ok {
				device.KernelName = val
			}
		}
		if val, ok := deviceProps["PKNAME"]; ok {
			if val != "" {
				device.Parents = append(device.Parents, val)
			}
		}
	}
	klog.Infof("device name: %s, parents: %v", device.Name, device.Parents)
	return device, nil
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
