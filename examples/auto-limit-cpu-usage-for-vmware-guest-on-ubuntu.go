package main

import (
	. "github.com/ChaunceyShannon/golanglibs"
)

var currentlimiting []*PexpectStruct
var lastlimit float64

func main() {
	vw := Cmd.GetOutput("xdotool search --name 'VMware Workstation'").Strip().S
	Lg.Trace("Window of VMware Workstation:", vw)

	for {
		cw := Cmd.GetOutput("xdotool getmouselocation").
			Split(" ")[3].
			Split(":")[1].S
		Lg.Trace("Current window on:", cw)

		if cw != vw {
			if Time.Now()-lastlimit > 60 {
				Lg.Trace("Mouse is now out of VMware Workstation, starting to limit cpu usage to 10% for each vmware-vmx procress")
				for _, pid := range Cmd.GetOutputWithShell("ps -ef | grep -v grep | grep vmware-vmx | awk '{print $2}'").Strip().Splitlines(true) {
					// If limited
					if func() bool {
						for _, p := range currentlimiting {
							if Int(pid.S) == p.Pid {
								return true
							}
						}
						return false
					}() {
						Lg.Trace("Pid already limited, skip:", pid.S)
						continue
					}
					Lg.Trace("Limiting pid:", pid.S)
					go func(pid string) {
						p := Tools.Pexpect("cpulimit -l 10 -p " + pid)
						currentlimiting = append(currentlimiting, p)
					}(pid.S)
				}
				lastlimit = Time.Now()
			}
		} else {
			Lg.Trace("Mouse is inside the VMWare Workstation")
			if len(currentlimiting) != 0 {
				Lg.Trace("Remove the liminations")
				for _, p := range currentlimiting {
					Lg.Trace("Terminate cpulimit:", p.Pid)
					p.Close()
				}
			}
			currentlimiting = []*PexpectStruct{}
		}

		Time.Sleep(3)
	}
}
