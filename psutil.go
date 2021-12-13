package golanglibs

import "runtime"

type progressCPUUsageStruct struct {
	pid  int64
	cpu  float64
	cmd  string
	name string
}

func getSystemProgressCPUUsage() (res map[int64]progressCPUUsageStruct) {
	pg := make(map[string]float64)
	res = make(map[int64]progressCPUUsageStruct)

	var line string
	for _, line = range String(Open("/proc/stat").Read()).Split("\n") {
		if String("cpu ").In(line) {
			break
		}
	}

	var totalCPUSlice1 float64
	for _, i := range String(line).Split()[2:] {
		totalCPUSlice1 = totalCPUSlice1 + Float64(i)
	}

	for _, pid := range listdir("/proc") {
		if !String(pid).Isdigit() {
			continue
		}
		Try(func() {
			a := String(Open("/proc/" + pid + "/stat").Read()).Split()
			totalProcessSlice1 := Float64(Int(a[13]) + Int(a[14]) + Int(a[15]) + Int(a[16]))
			pg[pid] = totalProcessSlice1
		})
	}

	sleep(1)

	for _, line = range String(Open("/proc/stat").Read()).Split("\n") {
		if String("cpu ").In(line) {
			break
		}
	}

	var totalCPUSlice2 float64
	for _, i := range String(line).Split()[2:] {
		totalCPUSlice2 = totalCPUSlice2 + Float64(i)
	}

	for _, pid := range listdir("/proc") {
		if !String(pid).Isdigit() {
			continue
		}
		Try(func() {
			a := String(Open("/proc/" + pid + "/stat").Read()).Split()
			totalProcessSlice2 := Float64(Int(a[13]) + Int(a[14]) + Int(a[15]) + Int(a[16]))
			_, found := pg[pid]
			if found {
				cpuusage := (totalProcessSlice2 - pg[pid]) / (totalCPUSlice2 - totalCPUSlice1) * 100 * Float64(runtime.NumCPU())
				res[Int64(pid)] = progressCPUUsageStruct{
					pid:  Int64(pid),
					cpu:  cpuusage,
					cmd:  String(Open("/proc/"+pid+"/cmdline").Read()).Replace("\x00", " ").Strip().Get(),
					name: String(a[1]).Strip("()").Get(),
				}
			}
		})
	}

	return
}
