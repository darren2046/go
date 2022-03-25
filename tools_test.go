package golanglibs

import "testing"

func TestNats(t *testing.T) {
	server := "nat://nats.nats.svc.cluster.local"
	subj := Tools.Nats(server).Subject("mysubject")

	go func() {
		for msg := range subj.Subscribe() {
			Lg.Trace(msg)
		}
	}()

	for {
		sleeptime := Random.Int(1, 3)
		Time.Sleep(sleeptime)
		subj.Publish("sleep for " + Str(sleeptime) + " second(s) just now")
	}
}

// func TestJieba(t *testing.T) {
// 	jieba := Tools.Jieba()
// 	Print(jieba.Cut("我来到北京清华大学"))
// 	Print(jieba.Cut("Running tool: /usr/local/bin/go test -timeout 30s -run ^TestJieba$ github.com/ChaunceyShannon/golanglibs"))
// }
