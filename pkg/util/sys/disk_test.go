package sys_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"k8s.io/klog/v2"

	exectest "github.com/microyahoo/fio-benchmark/pkg/util/exec/test"
	"github.com/microyahoo/fio-benchmark/pkg/util/sys"
)

func TestDiskSuite(t *testing.T) {
	suite.Run(t, new(diskSuite))
}

type diskSuite struct {
	suite.Suite
}

func (s *diskSuite) TestDiscoverDevices() {
	// lsblk --all --bytes --pairs --output SIZE,ROTA,RO,TYPE,PKNAME,NAME,KNAME,UUID,WWN,MOUNTPOINT
	lsblkOutput := `SIZE="1073741312" ROTA="1" RO="0" TYPE="rom" PKNAME="" NAME="/dev/sr0" KNAME="/dev/sr0" UUID="" WWN="" MOUNTPOINT=""
SIZE="107374182400" ROTA="1" RO="0" TYPE="disk" PKNAME="" NAME="/dev/vda" KNAME="/dev/vda" UUID="" WWN="" MOUNTPOINT=""
SIZE="1073741824" ROTA="1" RO="0" TYPE="part" PKNAME="/dev/vda" NAME="/dev/vda1" KNAME="/dev/vda1" UUID="a080444c-7927-49f7-b94f-e20f823bbc95" WWN="" MOUNTPOINT="/boot"
SIZE="63349719040" ROTA="1" RO="0" TYPE="part" PKNAME="/dev/vda" NAME="/dev/vda2" KNAME="/dev/vda2" UUID="jDjk4o-AaZU-He1S-8t56-4YEY-ujTp-ozFrK5" WWN="" MOUNTPOINT=""
SIZE="99849601024" ROTA="1" RO="0" TYPE="lvm" PKNAME="/dev/vda2" NAME="/dev/mapper/centos-root" KNAME="/dev/dm-0" UUID="5e322b94-4141-4a15-ae29-4136ae9c2e15" WWN="" MOUNTPOINT="/"
SIZE="6442450944" ROTA="1" RO="0" TYPE="lvm" PKNAME="/dev/vda2" NAME="/dev/mapper/centos-swap" KNAME="/dev/dm-1" UUID="d59f7992-9027-407a-84b3-ec69c3dadd4e" WWN="" MOUNTPOINT=""
SIZE="42949672960" ROTA="1" RO="0" TYPE="part" PKNAME="/dev/vda" NAME="/dev/vda3" KNAME="/dev/vda3" UUID="Qn0c4t-Sf93-oIDr-e57o-XQ73-DsyG-pGI8X0" WWN="" MOUNTPOINT=""
SIZE="99849601024" ROTA="1" RO="0" TYPE="lvm" PKNAME="/dev/vda3" NAME="/dev/mapper/centos-root" KNAME="/dev/dm-0" UUID="5e322b94-4141-4a15-ae29-4136ae9c2e15" WWN="" MOUNTPOINT="/"
SIZE="53687091200" ROTA="1" RO="0" TYPE="disk" PKNAME="" NAME="/dev/vdb" KNAME="/dev/vdb" UUID="klSb8f-Uq7t-WCaj-ZAeF-ShgA-mcZB-mojGe5" WWN="" MOUNTPOINT=""
SIZE="53682896896" ROTA="1" RO="0" TYPE="lvm" PKNAME="/dev/vdb" NAME="/dev/mapper/ceph--cfa0aaf9--bd31--401b--8210--6bf0fe67803c-osd--block--2af161f2--cbab--4bf0--a655--8490c8073129" KNAME="/dev/dm-5" UUID="" WWN="" MOUNTPOINT=""
SIZE="53687091200" ROTA="1" RO="0" TYPE="disk" PKNAME="" NAME="/dev/vdc" KNAME="/dev/vdc" UUID="ysYGKD-XKQB-VPTP-iCyX-ldsq-GKEC-Bx9fZX" WWN="" MOUNTPOINT=""
SIZE="53682896896" ROTA="1" RO="0" TYPE="lvm" PKNAME="/dev/vdc" NAME="/dev/mapper/ceph--9ae8c015--ddf8--4acc--944b--b6313fba74aa-osd--block--27180b72--74c8--4967--9a37--8634924236ea" KNAME="/dev/dm-6" UUID="" WWN="" MOUNTPOINT=""
SIZE="10737418240" ROTA="1" RO="0" TYPE="loop" PKNAME="" NAME="/dev/loop0" KNAME="/dev/loop0" UUID="" WWN="" MOUNTPOINT=""
SIZE="10733223936" ROTA="1" RO="0" TYPE="lvm" PKNAME="/dev/loop0" NAME="/dev/mapper/test--rook--vg-test--rook--lv" KNAME="/dev/dm-2" UUID="7acb62e7-ebc8-44f8-b2f0-d1e0a9b62439" WWN="" MOUNTPOINT="/mount/test_vdb"
SIZE="10737418240" ROTA="1" RO="0" TYPE="loop" PKNAME="" NAME="/dev/loop1" KNAME="/dev/loop1" UUID="" WWN="" MOUNTPOINT=""
SIZE="10733223936" ROTA="1" RO="0" TYPE="lvm" PKNAME="/dev/loop1" NAME="/dev/mapper/test--rook--vg1-test--rook--lv1" KNAME="/dev/dm-3" UUID="" WWN="" MOUNTPOINT=""
SIZE="53687091200" ROTA="0" RO="0" TYPE="disk" PKNAME="" NAME="/dev/rbd0" KNAME="/dev/rbd0" UUID="" WWN="" MOUNTPOINT=""
SIZE="53687091200" ROTA="0" RO="0" TYPE="disk" PKNAME="" NAME="rbd1" KNAME="rbd1" UUID="" WWN="" MOUNTPOINT=""
SIZE="53687091200" ROTA="1" RO="0" TYPE="disk" PKNAME="" NAME="/dev/vdd" KNAME="/dev/vdd" UUID="" WWN="" MOUNTPOINT=""
SIZE="53686042624" ROTA="1" RO="0" TYPE="part" PKNAME="/dev/vdd" NAME="/dev/vdd1" KNAME="/dev/vdd1" UUID="0hnEJg-LbJz-1fLe-GWVa-wSpq-WKLZ-UOC3hK" WWN="" MOUNTPOINT=""
SIZE="207215394816" ROTA="1" RO="0" TYPE="lvm" PKNAME="/dev/vdd1" NAME="/dev/mapper/centos-root" KNAME="/dev/dm-0" UUID="5e322b94-4141-4a15-ae29-4136ae9c2e15" WWN="" MOUNTPOINT="/"`

	// udevadm info --query=property <device>
	udevInfo1 := `DEVLINKS=/dev/disk/by-id/virtio-8560782279146-0 /dev/disk/by-path/pci-0000:00:0a.0 /dev/disk/by-path/virtio-pci-0000:00:0a.0
DEVNAME=/dev/vda
DEVPATH=/devices/pci0000:00/0000:00:0a.0/virtio1/block/vda
DEVTYPE=disk
ID_PART_TABLE_TYPE=dos
ID_PATH=pci-0000:00:0a.0
ID_PATH_TAG=pci-0000_00_0a_0
ID_SERIAL=8560782279146-0
MAJOR=252
MINOR=0
MPATH_SBIN_PATH=/sbin
SUBSYSTEM=block
TAGS=:systemd:
USEC_INITIALIZED=60219`
	udevInfo2 := `DEVLINKS=/dev/disk/by-id/lvm-pv-uuid-klSb8f-Uq7t-WCaj-ZAeF-ShgA-mcZB-mojGe5 /dev/disk/by-id/virtio-8560782279146-1 /dev/disk/by-path/pci-0000:00:0b.0 /dev/disk/by-path/virtio-pci-0000:00:0b.0
DEVNAME=/dev/vdb
DEVPATH=/devices/pci0000:00/0000:00:0b.0/virtio2/block/vdb
DEVTYPE=disk
ID_FS_TYPE=LVM2_member
ID_FS_USAGE=raid
ID_FS_UUID=klSb8f-Uq7t-WCaj-ZAeF-ShgA-mcZB-mojGe5
ID_FS_UUID_ENC=klSb8f-Uq7t-WCaj-ZAeF-ShgA-mcZB-mojGe5
ID_FS_VERSION=LVM2 001
ID_PATH=pci-0000:00:0b.0
ID_PATH_TAG=pci-0000_00_0b_0
ID_SERIAL=8560782279146-1
MAJOR=252
MINOR=16
MPATH_SBIN_PATH=/sbin
SUBSYSTEM=block
SYSTEMD_ALIAS=/dev/block/252:16
SYSTEMD_READY=1
SYSTEMD_WANTS=lvm2-pvscan@252:16.service
TAGS=:systemd:
USEC_INITIALIZED=61556`
	udevInfo3 := `DEVLINKS=/dev/disk/by-id/lvm-pv-uuid-ysYGKD-XKQB-VPTP-iCyX-ldsq-GKEC-Bx9fZX /dev/disk/by-id/virtio-8560782279146-3 /dev/disk/by-path/pci-0000:00:0d.0 /dev/disk/by-path/virtio-pci-0000:00:0d.0
DEVNAME=/dev/vdc
DEVPATH=/devices/pci0000:00/0000:00:0d.0/virtio3/block/vdc
DEVTYPE=disk
ID_FS_TYPE=LVM2_member
ID_FS_USAGE=raid
ID_FS_UUID=ysYGKD-XKQB-VPTP-iCyX-ldsq-GKEC-Bx9fZX
ID_FS_UUID_ENC=ysYGKD-XKQB-VPTP-iCyX-ldsq-GKEC-Bx9fZX
ID_FS_VERSION=LVM2 001
ID_PATH=pci-0000:00:0d.0
ID_PATH_TAG=pci-0000_00_0d_0
ID_SERIAL=8560782279146-3
MAJOR=252
MINOR=32
MPATH_SBIN_PATH=/sbin
SUBSYSTEM=block
SYSTEMD_ALIAS=/dev/block/252:32
SYSTEMD_READY=1
SYSTEMD_WANTS=lvm2-pvscan@252:32.service
TAGS=:systemd:
USEC_INITIALIZED=62265`
	udevInfo4 := `DEVLINKS=/dev/disk/by-id/virtio-7076686573460-4 /dev/disk/by-path/pci-0000:00:0e.0 /dev/disk/by-path/virtio-pci-0000:00:0e.0
DEVNAME=/dev/vdd
DEVPATH=/devices/pci0000:00/0000:00:0e.0/virtio5/block/vdd
DEVTYPE=disk
DM_MULTIPATH_TIMESTAMP=1682492340
ID_PATH=pci-0000:00:0e.0
ID_PATH_TAG=pci-0000_00_0e_0
ID_SERIAL=7076686573460-4
MAJOR=252
MINOR=48
MPATH_SBIN_PATH=/sbin
SUBSYSTEM=block
TAGS=:systemd:
USEC_INITIALIZED=29456997`

	// lsblk --noheadings --path --list --output NAME <device>
	deviceChild1 := `/dev/vda
/dev/vda1
/dev/vda2
/dev/mapper/centos-root
/dev/mapper/centos-swap
/dev/vda3
/dev/mapper/centos-root`

	deviceChild2 := `/dev/vdb
/dev/mapper/ceph--cfa0aaf9--bd31--401b--8210--6bf0fe67803c-osd--block--2af161f2--cbab--4bf0--a655--8490c8073129`
	deviceChild3 := `/dev/vdc
/dev/mapper/ceph--9ae8c015--ddf8--4acc--944b--b6313fba74aa-osd--block--27180b72--74c8--4967--9a37--8634924236ea`

	// lsblk <device> --bytes --paths --pairs --output NAME,SIZE,TYPE,PKNAME
	partition1 := `NAME="/dev/vda" SIZE="107374182400" TYPE="disk" PKNAME=""
NAME="/dev/vda1" SIZE="1073741824" TYPE="part" PKNAME="/dev/vda"
NAME="/dev/vda2" SIZE="63349719040" TYPE="part" PKNAME="/dev/vda"
NAME="/dev/mapper/centos-root" SIZE="99849601024" TYPE="lvm" PKNAME="/dev/vda2"
NAME="/dev/mapper/centos-swap" SIZE="6442450944" TYPE="lvm" PKNAME="/dev/vda2"
NAME="/dev/vda3" SIZE="42949672960" TYPE="part" PKNAME="/dev/vda"
NAME="/dev/mapper/centos-root" SIZE="99849601024" TYPE="lvm" PKNAME="/dev/vda3"`
	partition2 := `NAME="/dev/vdb" SIZE="53687091200" TYPE="disk" PKNAME=""
NAME="/dev/mapper/ceph--cfa0aaf9--bd31--401b--8210--6bf0fe67803c-osd--block--2af161f2--cbab--4bf0--a655--8490c8073129" SIZE="53682896896" TYPE="lvm" PKNAME="/dev/vdb"`
	partition3 := `NAME="/dev/vdc" SIZE="53687091200" TYPE="disk" PKNAME=""
NAME="/dev/mapper/ceph--9ae8c015--ddf8--4acc--944b--b6313fba74aa-osd--block--27180b72--74c8--4967--9a37--8634924236ea" SIZE="53682896896" TYPE="lvm" PKNAME="/dev/vdc"`

	udevInfoVda1 := `DEVLINKS=/dev/disk/by-id/virtio-8560782279146-0-part1 /dev/disk/by-path/pci-0000:00:0a.0-part1 /dev/disk/by-path/virtio-pci-0000:00:0a.0-part1 /dev/disk/by-uuid/a080444c-7927-49f7-b94f-e20f823bbc95
DEVNAME=/dev/vda1
DEVPATH=/devices/pci0000:00/0000:00:0a.0/virtio1/block/vda/vda1
DEVTYPE=partition
ID_FS_TYPE=xfs
ID_FS_USAGE=filesystem
ID_PATH=pci-0000:00:0a.0
ID_PATH_TAG=pci-0000_00_0a_0
ID_SERIAL=8560782279146-0
MAJOR=252
MINOR=1
SUBSYSTEM=block
TAGS=:systemd:
USEC_INITIALIZED=60434`
	udevInfoVda2 := `DEVLINKS=/dev/disk/by-id/lvm-pv-uuid-jDjk4o-AaZU-He1S-8t56-4YEY-ujTp-ozFrK5 /dev/disk/by-id/virtio-8560782279146-0-part2 /dev/disk/by-path/pci-0000:00:0a.0-part2 /dev/disk/by-path/virtio-pci-0000:00:0a.0-part2
DEVNAME=/dev/vda2
DEVPATH=/devices/pci0000:00/0000:00:0a.0/virtio1/block/vda/vda2
DEVTYPE=partition
ID_FS_TYPE=LVM2_member
ID_FS_USAGE=raid
ID_PATH=pci-0000:00:0a.0
ID_PATH_TAG=pci-0000_00_0a_0
ID_SERIAL=8560782279146-0
MAJOR=252
MINOR=2
SUBSYSTEM=block
SYSTEMD_ALIAS=/dev/block/252:2
SYSTEMD_READY=1
SYSTEMD_WANTS=lvm2-pvscan@252:2.service
TAGS=:systemd:
USEC_INITIALIZED=60663`
	udevInfoVda3 := `DEVLINKS=/dev/disk/by-id/lvm-pv-uuid-Qn0c4t-Sf93-oIDr-e57o-XQ73-DsyG-pGI8X0 /dev/disk/by-id/virtio-8560782279146-0-part3 /dev/disk/by-path/pci-0000:00:0a.0-part3 /dev/disk/by-path/virtio-pci-0000:00:0a.0-part3
DEVNAME=/dev/vda3
DEVPATH=/devices/pci0000:00/0000:00:0a.0/virtio1/block/vda/vda3
DEVTYPE=partition
ID_FS_TYPE=LVM2_member
ID_FS_USAGE=raid
ID_PATH=pci-0000:00:0a.0
ID_PATH_TAG=pci-0000_00_0a_0
ID_SERIAL=8560782279146-0
MAJOR=252
MINOR=3
SUBSYSTEM=block
SYSTEMD_ALIAS=/dev/block/252:3
SYSTEMD_READY=1
SYSTEMD_WANTS=lvm2-pvscan@252:3.service
TAGS=:systemd:
USEC_INITIALIZED=60881`
	udevInfoVdd1 := `DEVLINKS=/dev/disk/by-id/lvm-pv-uuid-0hnEJg-LbJz-1fLe-GWVa-wSpq-WKLZ-UOC3hK /dev/disk/by-id/virtio-7076686573460-4-part1 /dev/disk/by-path/pci-0000:00:0e.0-part1 /dev/disk/by-path/virtio-pci-0000:00:0e.0-part1
DEVNAME=/dev/vdd1
DEVPATH=/devices/pci0000:00/0000:00:0e.0/virtio5/block/vdd/vdd1
DEVTYPE=partition
ID_FS_TYPE=LVM2_member
ID_FS_USAGE=raid
ID_FS_UUID=0hnEJg-LbJz-1fLe-GWVa-wSpq-WKLZ-UOC3hK
ID_FS_UUID_ENC=0hnEJg-LbJz-1fLe-GWVa-wSpq-WKLZ-UOC3hK
ID_FS_VERSION=LVM2 001
ID_PATH=pci-0000:00:0e.0
ID_PATH_TAG=pci-0000_00_0e_0
ID_SERIAL=7076686573460-4
MAJOR=252
MINOR=49
SUBSYSTEM=block
SYSTEMD_ALIAS=/dev/block/252:49
SYSTEMD_READY=1
SYSTEMD_WANTS=lvm2-pvscan@252:49.service
TAGS=:systemd:
USEC_INITIALIZED=2851184`

	executor := &exectest.MockExecutor{
		MockExecuteCommandWithOutput: func(command string, args ...string) (string, error) {
			klog.Infof("run command %s: %s", command, args)
			if len(args) > 1 && args[0] == "--all" && args[1] == "--bytes" {
				return lsblkOutput, nil
			}
			if len(args) > 2 && args[0] == "info" && args[1] == "--query=property" && args[2] == "/dev/vda" {
				return udevInfo1, nil
			}
			if len(args) > 2 && args[0] == "info" && args[1] == "--query=property" && args[2] == "/dev/vdb" {
				return udevInfo2, nil
			}
			if len(args) > 2 && args[0] == "info" && args[1] == "--query=property" && args[2] == "/dev/vdc" {
				return udevInfo3, nil
			}
			if len(args) > 2 && args[0] == "info" && args[1] == "--query=property" && args[2] == "/dev/vdd" {
				return udevInfo4, nil
			}
			if len(args) > 5 && args[0] == "--noheadings" && args[5] == "/dev/vda" {
				return deviceChild1, nil
			}
			if len(args) > 5 && args[0] == "--noheadings" && args[5] == "/dev/vdb" {
				return deviceChild2, nil
			}
			if len(args) > 5 && args[0] == "--noheadings" && args[5] == "/dev/vdc" {
				return deviceChild3, nil
			}
			if len(args) > 1 && args[0] == "/dev/vda" && args[1] == "--bytes" {
				return partition1, nil
			}
			if len(args) > 1 && args[0] == "/dev/vdb" && args[1] == "--bytes" {
				return partition2, nil
			}
			if len(args) > 1 && args[0] == "/dev/vdc" && args[1] == "--bytes" {
				return partition3, nil
			}
			if len(args) > 2 && args[0] == "info" && args[1] == "--query=property" && args[2] == "/dev/vda1" {
				return udevInfoVda1, nil
			}
			if len(args) > 2 && args[0] == "info" && args[1] == "--query=property" && args[2] == "/dev/vda2" {
				return udevInfoVda2, nil
			}
			if len(args) > 2 && args[0] == "info" && args[1] == "--query=property" && args[2] == "/dev/vda3" {
				return udevInfoVda3, nil
			}
			if len(args) > 2 && args[0] == "info" && args[1] == "--query=property" && args[2] == "/dev/vdd1" {
				return udevInfoVdd1, nil
			}
			return "", errors.New("error")
		},
	}
	deviceInfos, err := sys.DiscoverDevices(executor)
	s.NoError(err)
	expectedInfos := map[string]*sys.LocalDevice{
		"/dev/mapper/centos-root": {
			Name:        "/dev/mapper/centos-root",
			Parents:     []string{"/dev/vda2", "/dev/vda3", "/dev/vdd1"},
			HasChildren: false,
			DevLinks:    "",
			Size:        99849601024,
			UUID:        "5e322b94-4141-4a15-ae29-4136ae9c2e15",
			Type:        "lvm",
			Rotational:  true,
			Readonly:    false,
			Partitions:  nil,
			RealPath:    "/dev/mapper/centos-root",
			KernelName:  "/dev/dm-0",
			Encrypted:   false,
			IsRoot:      true,
			MountPoint:  "/",
			DeviceClass: "",
		},
		"/dev/mapper/centos-swap": {
			Name:        "/dev/mapper/centos-swap",
			Parents:     []string{"/dev/vda2"},
			HasChildren: false,
			Size:        6442450944,
			GUID:        "",
			UUID:        "d59f7992-9027-407a-84b3-ec69c3dadd4e",
			Serial:      "",
			Bus:         "",
			Type:        "lvm",
			Rotational:  true,
			Partitions:  nil,
			RealPath:    "/dev/mapper/centos-swap",
			KernelName:  "/dev/dm-1",
			Encrypted:   false,
			IsRoot:      false,
			MountPoint:  "",
			Empty:       false,
			DeviceClass: "",
		},
		"/dev/mapper/ceph--9ae8c015--ddf8--4acc--944b--b6313fba74aa-osd--block--27180b72--74c8--4967--9a37--8634924236ea": {
			Name:        "/dev/mapper/ceph--9ae8c015--ddf8--4acc--944b--b6313fba74aa-osd--block--27180b72--74c8--4967--9a37--8634924236ea",
			Parents:     []string{"/dev/vdc"},
			HasChildren: false,
			DevLinks:    "",
			Size:        53682896896,
			Type:        "lvm",
			Rotational:  true,
			Partitions:  nil,
			RealPath:    "/dev/mapper/ceph--9ae8c015--ddf8--4acc--944b--b6313fba74aa-osd--block--27180b72--74c8--4967--9a37--8634924236ea",
			KernelName:  "/dev/dm-6",
			IsRoot:      false,
			MountPoint:  "",
			Empty:       false,
			DeviceClass: "",
		},
		"/dev/mapper/ceph--cfa0aaf9--bd31--401b--8210--6bf0fe67803c-osd--block--2af161f2--cbab--4bf0--a655--8490c8073129": {
			Name:        "/dev/mapper/ceph--cfa0aaf9--bd31--401b--8210--6bf0fe67803c-osd--block--2af161f2--cbab--4bf0--a655--8490c8073129",
			Parents:     []string{"/dev/vdb"},
			HasChildren: false,
			DevLinks:    "",
			Size:        53682896896,
			Type:        "lvm",
			Rotational:  true,
			Partitions:  nil,
			RealPath:    "/dev/mapper/ceph--cfa0aaf9--bd31--401b--8210--6bf0fe67803c-osd--block--2af161f2--cbab--4bf0--a655--8490c8073129",
			KernelName:  "/dev/dm-5",
			Encrypted:   false,
			IsRoot:      false,
			MountPoint:  "",
			Empty:       false,
			DeviceClass: "",
		},
		"/dev/mapper/test--rook--vg-test--rook--lv": {
			Name:        "/dev/mapper/test--rook--vg-test--rook--lv",
			Parents:     []string{"/dev/loop0"},
			HasChildren: false,
			DevLinks:    "",
			Size:        10733223936,
			UUID:        "7acb62e7-ebc8-44f8-b2f0-d1e0a9b62439",
			Type:        "lvm",
			Rotational:  true,
			Readonly:    false,
			Partitions:  nil,
			RealPath:    "/dev/mapper/test--rook--vg-test--rook--lv",
			KernelName:  "/dev/dm-2",
			Encrypted:   false,
			IsRoot:      false,
			MountPoint:  "/mount/test_vdb",
			Empty:       false,
			DeviceClass: "",
		},
		"/dev/mapper/test--rook--vg1-test--rook--lv1": {
			Name:        "/dev/mapper/test--rook--vg1-test--rook--lv1",
			Parents:     []string{"/dev/loop1"},
			HasChildren: false,
			DevLinks:    "",
			Size:        10733223936,
			Type:        "lvm",
			Rotational:  true,
			Readonly:    false,
			Partitions:  nil,
			RealPath:    "/dev/mapper/test--rook--vg1-test--rook--lv1",
			KernelName:  "/dev/dm-3",
			IsRoot:      false,
			MountPoint:  "",
			Empty:       false,
			DeviceClass: "",
		},
		"/dev/vda": {
			Name:        "/dev/vda",
			Parents:     nil,
			HasChildren: true,
			DevLinks:    "/dev/disk/by-id/virtio-8560782279146-0 /dev/disk/by-path/pci-0000:00:0a.0 /dev/disk/by-path/virtio-pci-0000:00:0a.0",
			Size:        107374182400,
			Serial:      "8560782279146-0",
			Bus:         "",
			Type:        "disk",
			Rotational:  true,
			Readonly:    false,
			Partitions: []sys.Partition{
				{
					Name:       "/dev/vda1",
					Size:       0x40000000,
					Filesystem: "xfs",
				},
				{
					Name:       "/dev/vda2",
					Size:       0xebff00000,
					Filesystem: "LVM2_member",
				},
				{
					Name:       "/dev/vda3",
					Size:       0xa00000000,
					Filesystem: "LVM2_member",
				},
			},
			PathID:             "pci-0000:00:0a.0",
			WWN:                "",
			WWNVendorExtension: "",
			RealPath:           "/dev/vda",
			KernelName:         "/dev/vda",
			Encrypted:          false,
			IsRoot:             true,
			MountPoint:         "",
			Empty:              false,
			DeviceClass:        "hdd",
		},
		"/dev/vda1": {
			Name:               "/dev/vda1",
			Parents:            []string{"/dev/vda"},
			HasChildren:        false,
			DevLinks:           "/dev/disk/by-id/virtio-8560782279146-0-part1 /dev/disk/by-path/pci-0000:00:0a.0-part1 /dev/disk/by-path/virtio-pci-0000:00:0a.0-part1 /dev/disk/by-uuid/a080444c-7927-49f7-b94f-e20f823bbc95",
			Size:               1073741824,
			GUID:               "",
			UUID:               "a080444c-7927-49f7-b94f-e20f823bbc95",
			Serial:             "8560782279146-0",
			Type:               "part",
			Rotational:         true,
			Readonly:           false,
			Partitions:         nil,
			Filesystem:         "xfs",
			PathID:             "pci-0000:00:0a.0",
			WWN:                "",
			WWNVendorExtension: "",
			RealPath:           "/dev/vda1",
			KernelName:         "/dev/vda1",
			IsRoot:             false,
			MountPoint:         "/boot",
			Empty:              false,
			DeviceClass:        "",
		},
		"/dev/vda2": {
			Name:               "/dev/vda2",
			Parents:            []string{"/dev/vda"},
			HasChildren:        false,
			DevLinks:           "/dev/disk/by-id/lvm-pv-uuid-jDjk4o-AaZU-He1S-8t56-4YEY-ujTp-ozFrK5 /dev/disk/by-id/virtio-8560782279146-0-part2 /dev/disk/by-path/pci-0000:00:0a.0-part2 /dev/disk/by-path/virtio-pci-0000:00:0a.0-part2",
			Size:               63349719040,
			UUID:               "jDjk4o-AaZU-He1S-8t56-4YEY-ujTp-ozFrK5",
			Serial:             "8560782279146-0",
			Bus:                "",
			Type:               "part",
			Rotational:         true,
			Readonly:           false,
			Partitions:         nil,
			Filesystem:         "LVM2_member",
			PathID:             "pci-0000:00:0a.0",
			WWN:                "",
			WWNVendorExtension: "",
			RealPath:           "/dev/vda2",
			KernelName:         "/dev/vda2",
			Encrypted:          false,
			IsRoot:             true,
			MountPoint:         "",
			Empty:              false,
			DeviceClass:        "",
		},
		"/dev/vda3": {
			Name:               "/dev/vda3",
			Parents:            []string{"/dev/vda"},
			HasChildren:        false,
			DevLinks:           "/dev/disk/by-id/lvm-pv-uuid-Qn0c4t-Sf93-oIDr-e57o-XQ73-DsyG-pGI8X0 /dev/disk/by-id/virtio-8560782279146-0-part3 /dev/disk/by-path/pci-0000:00:0a.0-part3 /dev/disk/by-path/virtio-pci-0000:00:0a.0-part3",
			Size:               42949672960,
			GUID:               "",
			UUID:               "Qn0c4t-Sf93-oIDr-e57o-XQ73-DsyG-pGI8X0",
			Serial:             "8560782279146-0",
			Bus:                "",
			Type:               "part",
			Rotational:         true,
			Readonly:           false,
			Partitions:         nil,
			Filesystem:         "LVM2_member",
			Vendor:             "",
			Model:              "",
			PathID:             "pci-0000:00:0a.0",
			WWN:                "",
			WWNVendorExtension: "",
			RealPath:           "/dev/vda3",
			KernelName:         "/dev/vda3",
			Encrypted:          false,
			IsRoot:             true,
			MountPoint:         "",
			Empty:              false,
			DeviceClass:        "",
		},
		"/dev/vdb": {
			Name:        "/dev/vdb",
			HasChildren: true,
			DevLinks:    "/dev/disk/by-id/lvm-pv-uuid-klSb8f-Uq7t-WCaj-ZAeF-ShgA-mcZB-mojGe5 /dev/disk/by-id/virtio-8560782279146-1 /dev/disk/by-path/pci-0000:00:0b.0 /dev/disk/by-path/virtio-pci-0000:00:0b.0",
			Size:        53687091200,
			GUID:        "",
			UUID:        "klSb8f-Uq7t-WCaj-ZAeF-ShgA-mcZB-mojGe5",
			Serial:      "8560782279146-1",
			Type:        "disk",
			Rotational:  true,
			Readonly:    false,
			Partitions: []sys.Partition{
				{
					Name: "/dev/mapper/ceph--cfa0aaf9--bd31--401b--8210--6bf0fe67803c-osd--block--2af161f2--cbab--4bf0--a655--8490c8073129",
				},
			},
			Filesystem:         "LVM2_member",
			PathID:             "pci-0000:00:0b.0",
			WWN:                "",
			WWNVendorExtension: "",
			RealPath:           "/dev/vdb",
			KernelName:         "/dev/vdb",
			Encrypted:          false,
			IsRoot:             false,
			MountPoint:         "",
			Empty:              false,
			DeviceClass:        "hdd",
		},
		"/dev/vdc": {
			Name:        "/dev/vdc",
			HasChildren: true,
			DevLinks:    "/dev/disk/by-id/lvm-pv-uuid-ysYGKD-XKQB-VPTP-iCyX-ldsq-GKEC-Bx9fZX /dev/disk/by-id/virtio-8560782279146-3 /dev/disk/by-path/pci-0000:00:0d.0 /dev/disk/by-path/virtio-pci-0000:00:0d.0",
			Size:        53687091200,
			GUID:        "",
			UUID:        "ysYGKD-XKQB-VPTP-iCyX-ldsq-GKEC-Bx9fZX",
			Serial:      "8560782279146-3",
			Type:        "disk",
			Rotational:  true,
			Readonly:    false,
			Partitions: []sys.Partition{
				{
					Name: "/dev/mapper/ceph--9ae8c015--ddf8--4acc--944b--b6313fba74aa-osd--block--27180b72--74c8--4967--9a37--8634924236ea",
				},
			},
			Filesystem:         "LVM2_member",
			PathID:             "pci-0000:00:0d.0",
			WWN:                "",
			WWNVendorExtension: "",
			RealPath:           "/dev/vdc",
			KernelName:         "/dev/vdc",
			Encrypted:          false,
			IsRoot:             false,
			MountPoint:         "",
			Empty:              false,
			DeviceClass:        "hdd",
		},
		"/dev/vdd1": {
			Name:               "/dev/vdd1",
			HasChildren:        false,
			Parents:            []string{"/dev/vdd"},
			DevLinks:           "/dev/disk/by-id/lvm-pv-uuid-0hnEJg-LbJz-1fLe-GWVa-wSpq-WKLZ-UOC3hK /dev/disk/by-id/virtio-7076686573460-4-part1 /dev/disk/by-path/pci-0000:00:0e.0-part1 /dev/disk/by-path/virtio-pci-0000:00:0e.0-part1",
			Size:               53686042624,
			GUID:               "",
			UUID:               "0hnEJg-LbJz-1fLe-GWVa-wSpq-WKLZ-UOC3hK",
			Serial:             "7076686573460-4",
			Type:               "part",
			Rotational:         true,
			Readonly:           false,
			Filesystem:         "LVM2_member",
			PathID:             "pci-0000:00:0e.0",
			WWN:                "",
			WWNVendorExtension: "",
			RealPath:           "/dev/vdd1",
			KernelName:         "/dev/vdd1",
			Encrypted:          false,
			IsRoot:             true,
			MountPoint:         "",
			Empty:              false,
			DeviceClass:        "",
		},
		"/dev/vdd": {
			Name:               "/dev/vdd",
			HasChildren:        false,
			DevLinks:           "/dev/disk/by-id/virtio-7076686573460-4 /dev/disk/by-path/pci-0000:00:0e.0 /dev/disk/by-path/virtio-pci-0000:00:0e.0",
			Size:               53687091200,
			GUID:               "",
			UUID:               "",
			Serial:             "7076686573460-4",
			Type:               "disk",
			Rotational:         true,
			Readonly:           false,
			Filesystem:         "",
			PathID:             "pci-0000:00:0e.0",
			WWN:                "",
			WWNVendorExtension: "",
			RealPath:           "/dev/vdd",
			KernelName:         "/dev/vdd",
			Encrypted:          false,
			IsRoot:             true,
			MountPoint:         "",
			Empty:              true,
			DeviceClass:        "hdd",
		},
	}
	s.Equal(expectedInfos, deviceInfos)
}
