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

	var line *stringStruct
	for _, line = range String(Open("/proc/stat").Read().S).Split("\n") {
		if String("cpu ").In(line.S) {
			break
		}
	}

	var totalCPUSlice1 float64
	for _, i := range String(line.S).Split()[2:] {
		totalCPUSlice1 = totalCPUSlice1 + Float64(i)
	}

	for _, pid := range listdir("/proc") {
		if !String(pid).Isdigit() {
			continue
		}
		Try(func() {
			a := String(Open("/proc/" + pid + "/stat").Read().S).Split()
			totalProcessSlice1 := Float64(Int(a[13]) + Int(a[14]) + Int(a[15]) + Int(a[16]))
			pg[pid] = totalProcessSlice1
		})
	}

	sleep(1)

	for _, line = range String(Open("/proc/stat").Read().S).Split("\n") {
		if String("cpu ").In(line.S) {
			break
		}
	}

	var totalCPUSlice2 float64
	for _, i := range String(line.S).Split()[2:] {
		totalCPUSlice2 = totalCPUSlice2 + Float64(i)
	}

	for _, pid := range listdir("/proc") {
		if !String(pid).Isdigit() {
			continue
		}
		Try(func() {
			a := String(Open("/proc/" + pid + "/stat").Read().S).Split()
			totalProcessSlice2 := Float64(Int(a[13]) + Int(a[14]) + Int(a[15]) + Int(a[16]))
			_, found := pg[pid]
			if found {
				cpuusage := (totalProcessSlice2 - pg[pid]) / (totalCPUSlice2 - totalCPUSlice1) * 100 * Float64(runtime.NumCPU())
				res[Int64(pid)] = progressCPUUsageStruct{
					pid:  Int64(pid),
					cpu:  cpuusage,
					cmd:  String(Open("/proc/"+pid+"/cmdline").Read().S).Replace("\x00", " ").Strip().Get(),
					name: a[1].Strip("()").Get(),
				}
			}
		})
	}

	return
}
