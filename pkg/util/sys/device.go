package sys

import (
	"fmt"
	"strconv"
	"strings"

	"k8s.io/klog/v2"

	"github.com/microyahoo/fio-benchmark/pkg/util/exec"
)

const (
	// DiskType is a disk type
	DiskType = "disk"
	// SSDType is an ssd type
	SSDType = "ssd"
	// PartType is a partition type
	PartType = "part"
	// CryptType is an encrypted type
	CryptType = "crypt"
	// LVMType is an LVM type
	LVMType = "lvm"
	// MultiPath is for multipath devices
	MultiPath = "mpath"
	// LinearType is a linear type
	LinearType = "linear"

	sgdiskCmd = "sgdisk"

	// CephLVPrefix is the prefix of a LV owned by ceph-volume
	CephLVPrefix = "ceph--"
	// DeviceMapperPrefix is the prefix of a LV from the device mapper interface
	DeviceMapperPrefix = "dm-"

	ErrMsgSgdiskNotFound = "sgdiskNotFound"

	SystemRootPath   = "/"
	SystemRootfsPath = "/rootfs"

	DiskBusUsb  = "usb"
	DiskBusScsi = "scsi"
	DiskBusAta  = "ata"
)

// Partition represents a partition metadata
type Partition struct {
	Name       string `json:"name"`
	Size       uint64 `json:"size"`
	Label      string `json:"label"`
	Filesystem string `json:"filesystem"`
}

// LocalDevice contains information about an unformatted block device
type LocalDevice struct {
	// Name is the device name
	Name string `json:"name"`
	// Parent is the device parent's name
	Parent string `json:"parent"`
	// HasChildren is whether the device has a children device
	HasChildren bool `json:"has_children"`
	// DevLinks is the persistent device path on the host
	DevLinks string `json:"dev_links"`
	// Size is the device capacity in byte
	Size uint64 `json:"size"`
	// GUID(Globally Disk Identifier) is used to locate disk drive
	GUID string `json:"guid"`
	// UUID(Universally Unique Identifier) is used by /dev/disk/by-uuid
	UUID string `json:"uuid"`
	// Serial is the disk serial used by /dev/disk/by-id
	Serial string `json:"serial"`
	// Bus is the bus type of disk
	Bus string `json:"bus"`
	// Type is device type
	Type string `json:"type"`
	// Rotational is the boolean whether the device is rotational: true for hdd, false for ssd and nvme
	Rotational bool `json:"rotational"`
	// ReadOnly is the boolean whether the device is readonly
	Readonly bool `json:"read_only"`
	// Partitions is a partition slice
	Partitions []Partition `json:"partitions"`
	// Filesystem is the filesystem currently on the device
	Filesystem string `json:"filesystem"`
	// Vendor is the device vendor
	Vendor string `json:"vendor"`
	// Model is the device model
	Model string `json:"model"`
	// PathID is the path id of the device
	PathID string `json:"path_id"`
	// WWN is the world wide name of the device
	WWN string `json:"wwn"`
	// WWNVendorExtension is the WWN_VENDOR_EXTENSION from udev info
	WWNVendorExtension string `json:"wwn_vendor_extension"`
	// RealPath is the device pathname behind the PVC, behind /mnt/<pvc>/name
	RealPath string `json:"real_path,omitempty"`
	// KernelName is the kernel name of the device
	KernelName string `json:"kernel_name,omitempty"`
	// Whether this device should be encrypted
	Encrypted bool `json:"encrypted,omitempty"`
	// Whether system root
	IsRoot bool `json:"is_root"`
	// mount point
	MountPoint string `json:"mount_point"`
	// Empty checks whether the device is completely empty
	Empty bool `json:"empty"`
	// DeviceClass is the device class of device. (hdd, ssd, nvme)
	DeviceClass string `json:"device_class"`
}

// GetDevicePartitions gets partitions on a given device
func GetDevicePartitions(executor exec.Executor, device string) (partitions []Partition, unusedSpace uint64, err error) {

	var devicePath string
	splitDevicePath := strings.Split(device, "/")
	if len(splitDevicePath) == 1 {
		devicePath = fmt.Sprintf("/dev/%s", device) //device path for OSD on devices.
	} else {
		devicePath = device //use the exact device path (like /mnt/<pvc-name>) in case of PVC block device
	}

	output, err := executor.ExecuteCommandWithOutput("lsblk", devicePath,
		"--bytes", "--paths", "--pairs", "--output", "NAME,SIZE,TYPE,PKNAME")
	klog.Infof("Output: %+v", output)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get device %s partitions. %+v", device, err)
	}
	partInfo := strings.Split(output, "\n")
	var deviceSize uint64
	var totalPartitionSize uint64
	for _, info := range partInfo {
		props := parseKeyValuePairString(info)
		name := props["NAME"]
		if name == device {
			// found the main device
			klog.Infof("Device found - %s", name)
			deviceSize, err = strconv.ParseUint(props["SIZE"], 10, 64)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to get device %s size. %+v", device, err)
			}
		} else if props["PKNAME"] == device && props["TYPE"] == PartType {
			// found a partition
			p := Partition{Name: name}
			p.Size, err = strconv.ParseUint(props["SIZE"], 10, 64)
			if err != nil {
				return nil, 0, fmt.Errorf("failed to get partition %s size. %+v", name, err)
			}
			totalPartitionSize += p.Size

			info, err := GetUdevInfo(executor, name)
			if err != nil {
				return nil, 0, err
			}
			if v, ok := info["PARTNAME"]; ok {
				p.Label = v
			}
			if v, ok := info["ID_PART_ENTRY_NAME"]; ok {
				p.Label = v
			}
			if v, ok := info["ID_FS_TYPE"]; ok {
				p.Filesystem = v
			}

			partitions = append(partitions, p)
		} else if props["TYPE"] == LVMType && hasCephLvmPrefix(name) {
			p := Partition{Name: name}
			partitions = append(partitions, p)
		}
	}

	if deviceSize > 0 {
		unusedSpace = deviceSize - totalPartitionSize
	}
	return partitions, unusedSpace, nil
}

func hasCephLvmPrefix(deviceName string) bool {
	if strings.HasPrefix(deviceName, CephLVPrefix) ||
		strings.HasPrefix(strings.TrimPrefix(deviceName, "/dev/mapper/"), CephLVPrefix) ||
		strings.HasPrefix(strings.TrimPrefix(deviceName, "/dev/"), CephLVPrefix) {
		return true
	}
	return false
}

// GetDeviceProperties gets device properties
func GetDeviceProperties(executor exec.Executor, device string) (map[string]string, error) {
	// As we are mounting the block mode PVs on /mnt we use the entire path,
	// e.g., if the device path is /mnt/example-pvc then its taken completely
	// else if its just vdb then the following is used
	devicePath := strings.Split(device, "/")
	if len(devicePath) == 1 {
		device = fmt.Sprintf("/dev/%s", device)
	}
	return GetDevicePropertiesFromPath(executor, device)
}

// GetDevicePropertiesFromPath gets a device property from a path
func GetDevicePropertiesFromPath(executor exec.Executor, devicePath string) (map[string]string, error) {
	output, err := executor.ExecuteCommandWithOutput("lsblk", devicePath,
		"--bytes", "--nodeps", "--pairs", "--paths", "--output", "SIZE,ROTA,RO,TYPE,PKNAME,NAME,KNAME,UUID")
	if err != nil {
		klog.Errorf("failed to execute lsblk. output: %s", output)
		return nil, err
	}

	return parseKeyValuePairString(output), nil
}

// IsLV returns if a device is owned by LVM, is a logical volume
func IsLV(executor exec.Executor, devicePath string) (bool, error) {
	devProps, err := GetDevicePropertiesFromPath(executor, devicePath)
	if err != nil {
		return false, fmt.Errorf("failed to get device properties for %q: %+v", devicePath, err)
	}
	diskType, ok := devProps["TYPE"]
	if !ok {
		return false, fmt.Errorf("TYPE property is not found for %q", devicePath)
	}
	return diskType == LVMType, nil
}

// GetUdevInfo gets udev information
func GetUdevInfo(executor exec.Executor, device string) (map[string]string, error) {
	devicePath := strings.Split(device, "/")
	if len(devicePath) == 1 {
		device = fmt.Sprintf("/dev/%s", device)
	}
	output, err := executor.ExecuteCommandWithOutput("udevadm", "info", "--query=property", device)
	if err != nil {
		return nil, err
	}

	return parseUdevInfo(output), nil
}

// GetDeviceFilesystems get the file systems available
func GetDeviceFilesystems(executor exec.Executor, device string) (string, error) {
	devicePath := strings.Split(device, "/")
	if len(devicePath) == 1 {
		device = fmt.Sprintf("/dev/%s", device)
	}
	output, err := executor.ExecuteCommandWithOutput("udevadm", "info", "--query=property", device)
	if err != nil {
		return "", err
	}

	return parseFS(output), nil
}

func GetDiskDeviceClass(disk *LocalDevice) string {
	if disk.Rotational {
		return "hdd"
	}
	if strings.Contains(disk.RealPath, "nvme") {
		return "nvme"
	}
	return "ssd"
}

// GetLVName returns the LV name of the device in the form of "VG/LV".
func GetLVName(executor exec.Executor, devicePath string) (string, error) {
	devInfo, err := executor.ExecuteCommandWithOutput("dmsetup", "info", "-c", "--noheadings", "-o", "name", devicePath)
	if err != nil {
		return "", fmt.Errorf("failed to execute dmsetup info for %q. %v", devicePath, err)
	}
	out, err := executor.ExecuteCommandWithOutput("dmsetup", "splitname", "--noheadings", devInfo)
	if err != nil {
		return "", fmt.Errorf("failed to execute dmsetup splitname for %q. %v", devInfo, err)
	}
	split := strings.Split(out, ":")
	if len(split) < 2 {
		return "", fmt.Errorf("dmsetup splitname returned unexpected result for %q. output: %q", devInfo, out)
	}
	return fmt.Sprintf("%s/%s", split[0], split[1]), nil
}

// converts a raw key value pair string into a map of key value pairs
// example raw string of `foo="0" bar="1" baz="biz"` is returned as:
// map[string]string{"foo":"0", "bar":"1", "baz":"biz"}
func parseKeyValuePairString(propsRaw string) map[string]string {
	// first split the single raw string on spaces and initialize a map of
	// a length equal to the number of pairs
	props := strings.Split(propsRaw, " ")
	propMap := make(map[string]string, len(props))

	for _, kvpRaw := range props {
		// split each individual key value pair on the equals sign
		kvp := strings.Split(kvpRaw, "=")
		if len(kvp) == 2 {
			// first element is the final key, second element is the final value
			// (don't forget to remove surrounding quotes from the value)
			propMap[kvp[0]] = strings.Replace(kvp[1], `"`, "", -1)
		}
	}

	return propMap
}

// find fs from udevadm info
func parseFS(output string) string {
	m := parseUdevInfo(output)
	if v, ok := m["ID_FS_TYPE"]; ok {
		return v
	}
	return ""
}

func parseUdevInfo(output string) map[string]string {
	lines := strings.Split(output, "\n")
	result := make(map[string]string, len(lines))
	for _, v := range lines {
		pairs := strings.Split(v, "=")
		if len(pairs) > 1 {
			result[pairs[0]] = pairs[1]
		}
	}
	return result
}

// ListDevicesChild list all child available on a device
// For an encrypted device, it will return the encrypted device like so:
// lsblk --noheadings --output NAME --path --list /dev/sdd
// /dev/sdd
// /dev/mapper/ocs-deviceset-thin-1-data-0hmfgp-block-dmcrypt
func ListDevicesChild(executor exec.Executor, device string) ([]string, error) {
	devicePath := strings.Split(device, "/")
	if len(devicePath) == 1 {
		device = fmt.Sprintf("/dev/%s", device)
	}
	childListRaw, err := executor.ExecuteCommandWithOutput("lsblk", "--noheadings", "--path", "--list", "--output", "NAME", device)
	if err != nil {
		return []string{}, fmt.Errorf("failed to list child devices of %q. %v", device, err)
	}

	return strings.Split(childListRaw, "\n"), nil
}

// IsDeviceEncrypted returns whether the disk has a "crypt" label on it
func IsDeviceEncrypted(executor exec.Executor, device string) (bool, error) {
	deviceType, err := executor.ExecuteCommandWithOutput("lsblk", "--noheadings", "--output", "TYPE", device)
	if err != nil {
		return false, fmt.Errorf("failed to get devices type of %q. %v", device, err)
	}

	return deviceType == "crypt", nil
}
