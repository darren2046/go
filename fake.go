package golanglibs

func fakeNameChinese() string {
	for {
		length := Random.Int(2, 3)
		name := ""
		for idx := range Range(Int(length)) {
			py := Str(Random.Choice(pychinesenamemap[idx]))
			char := Str(Random.Choice(chinesenamechar[py][idx]))
			name += char
		}
		if String(name).Utf8Len() >= 2 {
			return name
		}
	}
}

func fakeNameEnglish() string {
	return String(" ").Join(
		[]string{
			Str(Random.Choice(englishFirstName)),
			Str(Random.Choice(englishLastName)),
		}).S
}
