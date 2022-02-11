package golanglibs

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/wunderlist/ttlcache"
)

type fsnotifyFileEventStruct struct {
	Action string
	Path   string
}

func inotify(path string) chan *fsnotifyFileEventStruct {
	ch := make(chan *fsnotifyFileEventStruct)
	watchList := make(map[string]struct{})

	watcher, err := fsnotify.NewWatcher()
	Panicerr(err)

	go func() {
		// 1秒内有同样的事件在同样的路径, 则忽略
		cache := ttlcache.NewCache(time.Second)

		for {
			select {
			case ev := <-watcher.Events:
				_, exists := cache.Get(ev.String() + abspath(ev.Name))
				if !exists {
					cache.Set(ev.String()+abspath(ev.Name), "")
					var action string
					if ev.Op == fsnotify.Create {
						action = "create"
						if pathIsDir(ev.Name) {
							err := watcher.Add(ev.Name)
							Panicerr(err)
							watchList[ev.Name] = struct{}{}
						}
					} else if ev.Op == fsnotify.Chmod {
						action = "chmod"
					} else if ev.Op == fsnotify.Remove {
						action = "remove"
						delete(watchList, ev.Name)
						go func() {
							// 如果action是delete, 那么就会从监控的列表当中移除
							// 如果10秒内再次出现, 就再次添加监听
							// 原因是vim编辑的时候并不是直接编辑, 而是先写入.swp文件, 在保存的时候删除源文件, 然后改名.swp文件, 会生成删除源文件的action
							for range Range(100) {
								if pathExists(ev.Name) {
									err = watcher.Add(ev.Name)
									Panicerr(err)
									watchList[ev.Name] = struct{}{}
									break
								}
								sleep(0.1)
							}
							if !Map(watchList).Has(ev.Name) {
								ch <- &fsnotifyFileEventStruct{
									action: action,
									path:   abspath(ev.Name),
								}
							}
							if len(watchList) == 0 {
								watcher.Close()
								close(ch)
							}

						}()
						continue
					} else if ev.Op == fsnotify.Write {
						action = "write"
					} else if ev.Op == fsnotify.Rename {
						action = "rename"
						delete(watchList, ev.Name)
					}

					ch <- &fsnotifyFileEventStruct{
						action: action,
						path:   abspath(ev.Name),
					}
				}
			case err := <-watcher.Errors:
				Panicerr(err)
			}
		}
	}()

	err = watcher.Add(path)
	Panicerr(err)
	watchList[path] = struct{}{}

	if pathIsDir(path) {
		for p := range walk(path) {
			if pathIsDir(p) {
				err := watcher.Add(p)
				Panicerr(err)
				watchList[p] = struct{}{}
			}
		}
	}

	return ch
}
