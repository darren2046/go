package golanglibs

import "github.com/mozillazg/go-pinyin"

func zh2PinYin(zh string) (ress []string) {
	a := pinyin.NewArgs()
	a.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
	p := pinyin.LazyPinyin(zh, a)

	filterChar := "1234567890qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM"

	var res []string
	for idx, v := range p {
		if len(v) != 1 {
			res = append(res, v)
		} else {
			if idx == 0 {
				res = append(res, v)
			} else {
				if len(p[idx-1]) == 1 {
					if String(p[idx-1]).In(filterChar) {
						res[len(res)-1] = res[len(res)-1] + v
					} else {
						res = append(res, v)
					}
				} else {
					res = append(res, v)
				}
			}
		}
	}

	for _, v := range res {
		vv := String(v).Filter(filterChar).Get()
		if len(vv) != 0 {
			ress = append(ress, vv)
		}
	}
	return
}
