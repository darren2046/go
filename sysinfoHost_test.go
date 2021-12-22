package golanglibs

import (
	"testing"
)

// types.HostInfo{
// 	Architecture:  "x86_64",
// 	BootTime:      2021-12-21 09:17:30 Local,
// 	Containerized: (*bool)(nil),
// 	Hostname:      "MacBook-Air.local",
// 	IPs:           []string{
// 	  "127.0.0.1/8",
// 	},
// 	KernelVersion: "21.1.0",
// 	MACs:          []string{
// 	  "d6:32:84:35:39:51",
// 	},
// 	OS: &types.OSInfo{
// 	  Type:     "macos",
// 	  Family:   "darwin",
// 	  Platform: "darwin",
// 	  Name:     "macOS",
// 	  Version:  "12.0.1",
// 	  Major:    12,
// 	  Minor:    0,
// 	  Patch:    1,
// 	  Build:    "21A559",
// 	  Codename: "",
// 	},
// 	Timezone:          "+01",
// 	TimezoneOffsetSec: 86400,
// 	UniqueID:          "FSKEJG-IESF-9fFE-DDFE-0992SKFEJJS",
//   }
func TestHostInfo(t *testing.T) {
	Print(sysinfohoststruct.Info())
}

// &types.HostMemoryInfo{
// 	Total:        0x0000000200000000,    // uint64
// 	Used:         0x000000018ff9a000,
// 	Available:    0x000000009077e000,
// 	Free:         0x0000000014a8e000,
// 	VirtualTotal: 0x0000000080000000,
// 	VirtualUsed:  0x0000000038380000,
// 	VirtualFree:  0x0000000047c80000,
// 	Metrics:      map[string]uint64{
// 	  "external_bytes":       0x000000005e202000,
// 	  "purgeable_bytes":      0x00000000065fb000,
// 	  "reactivated_bytes":    0x0000000cda64a000,
// 	  "speculative_bytes":    0x0000000002ced000,
// 	  "swap_ins_bytes":       0x00000000b45d7000,
// 	  "compressed_bytes":     0x000000007b881000,
// 	  "compressions_bytes":   0x00000011ecb4f000,
// 	  "inactive_bytes":       0x00000000756f5000,
// 	  "purged_bytes":         0x00000002752fc000,
// 	  "throttled_bytes":      0x0000000000000000,
// 	  "active_bytes":         0x0000000078823000,
// 	  "translation_faults":   0x00000000165fb38a,
// 	  "uncompressed_bytes":   0x00000001a68f9000,
// 	  "decompressions_bytes": 0x0000000cceae2000,
// 	  "internal_bytes":       0x0000000092a03000,
// 	  "page_ins_bytes":       0x00000004438bb000,
// 	  "page_outs_bytes":      0x000000002bb8d000,
// 	  "swap_outs_bytes":      0x00000000e9813000,
// 	  "wired_bytes":          0x0000000081d16000,
// 	  "zero_filled_bytes":    0x0000008d39a6e000,
// 	  "copy_on_write_faults": 0x00000000017f87b4,
// 	},
//   }
func TestHostMemory(t *testing.T) {
	Print(sysinfohoststruct.Memory())
}

// types.CPUTimes{
// 	User:    22085450000000,
// 	System:  13220050000000,
// 	Idle:    120236140000000,
// 	IOWait:  0,
// 	IRQ:     0,
// 	Nice:    0,
// 	SoftIRQ: 0,
// 	Steal:   0,
//   }
func TestHostCPUTimes(t *testing.T) {
	Print(sysinfohoststruct.CPUTimes())
}

//   types.ProcessInfo{
// 	Name: "Xorg",
// 	PID:  3428,
// 	PPID: 3426,
// 	CWD:  "/home/jack",
// 	Exe:  "/usr/lib/xorg/Xorg",
// 	Args: []string{
// 	  "/usr/lib/xorg/Xorg",
// 	  "vt2",
// 	  "-displayfd",
// 	  "3",
// 	  "-auth",
// 	  "/run/user/1000/gdm/Xauthority",
// 	  "-background",
// 	  "none",
// 	  "-noreset",
// 	  "-keeptty",
// 	  "-verbose",
// 	  "3",
// 	},
// 	StartTime: 2021-12-16 20:46:30 Local,
//   }
func TestProcessInfo(t *testing.T) {
	Print(getSysinfoProcess(3428).Info())
}

// types.MemoryInfo{
// 	Resident: 0x0000000008bdb000,
// 	Virtual:  0x0000000042156000,
// 	Metrics:  map[string]uint64{},
//   }
func TestProcessMemory(t *testing.T) {
	Print(getSysinfoProcess(3428).Memory())
}

// types.UserInfo{
// 	UID:  "1000",
// 	EUID: "1000",
// 	SUID: "1000",
// 	GID:  "1000",
// 	EGID: "1000",
// 	SGID: "1000",
//   }
func TestProcessUser(t *testing.T) {
	Print(getSysinfoProcess(3428).User())
}

// types.CPUTimes{
// 	User:    2091860000000,
// 	System:  1254070000000,
// 	Idle:    0,
// 	IOWait:  0,
// 	IRQ:     0,
// 	Nice:    0,
// 	SoftIRQ: 0,
// 	Steal:   0,
//   }
func TestProcessCPUTimes(t *testing.T) {
	Print(getSysinfoProcess(3428).CPUTimes())
}
