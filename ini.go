package golanglibs

import (
	"github.com/Go-ini/ini"
)

type iniStruct struct {
	cfg     *ini.File
	cfgPath string
}

func getIni(fpath ...string) *iniStruct {
	var cfg *ini.File
	var err error
	var cfgPath string
	if len(fpath) != 0 {
		cfgPath = fpath[0]
		if pathExists(cfgPath) {
			cfg, err = ini.Load(fpath[0])
			Panicerr(err)
		} else {
			cfg = ini.Empty()
		}
	} else {
		cfgPath = ""
		cfg = ini.Empty()
	}
	return &iniStruct{cfg: cfg, cfgPath: cfgPath}
}

func (m *iniStruct) Get(SectionKeyDefaultComment ...string) (res string) {
	if len(SectionKeyDefaultComment) == 2 {
		res = m.cfg.Section(SectionKeyDefaultComment[0]).Key(SectionKeyDefaultComment[1]).String()
	} else if len(SectionKeyDefaultComment) == 3 {
		res = m.cfg.Section(SectionKeyDefaultComment[0]).Key(SectionKeyDefaultComment[1]).String()
		if res == "" {
			res = SectionKeyDefaultComment[2]
			m.cfg.Section(SectionKeyDefaultComment[0]).Key(SectionKeyDefaultComment[1]).SetValue(SectionKeyDefaultComment[2])
		}
	} else if len(SectionKeyDefaultComment) == 4 {
		res = m.cfg.Section(SectionKeyDefaultComment[0]).Key(SectionKeyDefaultComment[1]).String()
		if res == "" {
			res = SectionKeyDefaultComment[2]
			m.cfg.Section(SectionKeyDefaultComment[0]).Key(SectionKeyDefaultComment[1]).SetValue(SectionKeyDefaultComment[2])
		}
		m.cfg.Section(SectionKeyDefaultComment[0]).Key(SectionKeyDefaultComment[1]).Comment = SectionKeyDefaultComment[3]
	}
	return
}

func (m *iniStruct) GetInt(key ...string) int {
	return Int(m.Get(key...))
}

func (m *iniStruct) GetInt64(key ...string) int64 {
	return Int64(m.Get(key...))
}

func (m *iniStruct) GetFloat64(key ...string) float64 {
	return Float64(m.Get(key...))
}

func (m *iniStruct) Set(SectionKeyValueComment ...string) {
	if len(SectionKeyValueComment) == 3 {
		m.cfg.Section(SectionKeyValueComment[0]).Key(SectionKeyValueComment[1]).SetValue(SectionKeyValueComment[2])
	} else if len(SectionKeyValueComment) == 4 {
		m.cfg.Section(SectionKeyValueComment[0]).Key(SectionKeyValueComment[1]).SetValue(SectionKeyValueComment[2])
		m.cfg.Section(SectionKeyValueComment[0]).Key(SectionKeyValueComment[1]).Comment = SectionKeyValueComment[3]
	} else {
		Panicerr("按顺序指定section, key, value(以及comment)")
	}
}

func (m *iniStruct) Save(fpath ...string) (exist bool) {
	if len(fpath) != 0 {
		exist = pathExists(fpath[0])
		m.cfg.SaveTo(fpath[0])
	} else if m.cfgPath != "" {
		exist = pathExists(m.cfgPath)
		m.cfg.SaveTo(m.cfgPath)
	}
	return
}
