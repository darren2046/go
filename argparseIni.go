package golanglibs

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/Go-ini/ini"
	"github.com/elliotchance/orderedmap"
)

type argparseIniStruct struct {
	cfg         *ini.File
	cfgPath     string
	args        *orderedmap.OrderedMap
	description string
}

func Argparser(description string) *argparseIniStruct {
	var cfg *ini.File
	var err error
	var cfgPath string

	for idx, value := range os.Args {
		if value == "-c" {
			cfgPath = os.Args[idx+1]
		}
	}

	if cfgPath == "" {
		c := pathJoin(getSelfDir(), Os.Path.Basename(os.Args[0])+".ini")
		if pathExists(c) {
			cfgPath = c
		} else {
			c = pathJoin(Os.Path.Basename(os.Args[0]) + ".ini")
			if pathExists(c) {
				cfgPath = c
			}
		}
	}

	if len(cfgPath) != 0 {
		if pathExists(cfgPath) {
			cfg, err = ini.Load(cfgPath)
			panicerr(err)
		} else {
			cfg = ini.Empty()
		}
	} else {
		cfgPath = ""
		cfg = ini.Empty()
	}

	return &argparseIniStruct{cfg: cfg, cfgPath: cfgPath, args: orderedmap.NewOrderedMap(), description: description}
}

func (m *argparseIniStruct) Get(section, key, defaultValue, comment string) (res string) {
	res = m.cfg.Section(section).Key(key).String()
	if res == "" {
		res = defaultValue
		m.cfg.Section(section).Key(key).SetValue(defaultValue)
	}
	m.cfg.Section(section).Key(key).Comment = comment
	if section != "" {
		m.args.Set(section+"."+key, comment)
	} else {
		m.args.Set(key, comment)
	}

	if section != "" {
		if os.Getenv(section+"."+key) != "" {
			res = os.Getenv(section + "." + key)
		}
	} else {
		if os.Getenv(key) != "" {
			res = os.Getenv(key)
		}
	}

	for idx, value := range os.Args {
		if section == "" {
			if "--"+key == value {
				res = os.Args[idx+1]
			}
		} else {
			if "--"+section+"."+key == value {
				res = os.Args[idx+1]
			}
		}

	}
	return
}

func (m *argparseIniStruct) GetInt(section, key, defaultValue, comment string) int {
	return Int(m.Get(section, key, defaultValue, comment))
}

func (m *argparseIniStruct) GetInt64(section, key, defaultValue, comment string) int64 {
	return Int64(m.Get(section, key, defaultValue, comment))
}

func (m *argparseIniStruct) GetFloat64(section, key, defaultValue, comment string) float64 {
	return Float64(m.Get(section, key, defaultValue, comment))
}

func (m *argparseIniStruct) GetBool(section, key, defaultValue, comment string) bool {
	return Bool(m.Get(section, key, defaultValue, comment))
}

func (m *argparseIniStruct) Save(fpath ...string) (exist bool) {
	exist = true
	if len(fpath) != 0 {
		exist = pathExists(fpath[0])
		m.cfg.SaveTo(fpath[0])
	} else if m.cfgPath != "" {
		exist = pathExists(m.cfgPath)
		m.cfg.SaveTo(m.cfgPath)
	}
	return
}

func (m *argparseIniStruct) GetHelpString() (h string) {
	maxLength := 0
	for _, k := range m.args.Keys() {
		if len(Str(k)) > maxLength {
			maxLength = len(Str(k))
		}
	}

	h = "\n    " + m.description + "\n\n"
	h = h + "    -" + fmt.Sprintf("%-"+Str(maxLength+5+1)+"v", "c") + " 配置文件\n"
	h = h + "    -" + fmt.Sprintf("%-"+Str(maxLength+5+1)+"v", "b") + " 在后台运行\n"
	for _, k := range m.args.Keys() {
		v, _ := m.args.Get(k)
		h = h + "    --" + fmt.Sprintf("%-"+Str(maxLength+5)+"v", Str(k)) + " " + Str(v) + "\n"
	}
	return
}

func (m *argparseIniStruct) ParseArgs() *argparseIniStruct {
	if Array(os.Args).Has("-h") || Array(os.Args).Has("--help") {
		fmt.Println(m.GetHelpString())
		Os.Exit(0)
	}
	if Array(os.Args).Has("-b") {
		if runtime.GOOS == "linux" {
			args := os.Args[1:]
			for i := 0; i < len(args); i++ {
				if args[i] == "-b" {
					args[i] = ""
					break
				}
			}
			cmd := exec.Command(os.Args[0], args...)
			cmd.Start()
			os.Exit(0)
		} else {
			fmt.Println("错误: -b 参数只在Linux环境测试过")
			Os.Exit(0)
		}
	}
	if !m.Save() { // 配置文件存在或者未指定则返回true，如果不存在则返回true
		fmt.Println("初始配置文件不存在，已保存模板配置文件。")
		Os.Exit(0)
	}
	return m
}
