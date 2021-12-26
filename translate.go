package golanglibs

// import (
// 	"encoding/json"
// )

// // 有一定的次数限制, 具体多少未知, 但很少, 到限制就hang住
// // func gen2zh(text string) string {
// // 	s, err := gtranslate.Translate(text, language.English, language.Chinese)
// // 	Panicerr(err)
// // 	return s
// // }

// type bresult struct {
// 	ErrorCode   string `json:"error_code"`
// 	ErrorMsg    string `json:"error_msg"`
// 	From        string `json:"from"`
// 	To          string `json:"to"`
// 	TransResult []struct {
// 		Src string `json:"src"`
// 		Dst string `json:"dst"`
// 	} `json:"trans_result"`
// }

// func baiduTranslateAnyToZH(text string) string {
// 	appID := ""
// 	secretKey := ""
// 	apiURL := "http://api.fanyi.baidu.com/api/trans/vip/translate"

// 	bres := bresult{}

// 	for {
// 		salt := Str(randint(32768, 65536))
// 		preSign := appID + text + salt + secretKey
// 		sign := md5sum(preSign)
// 		params := HttpParam{
// 			"q":     text,
// 			"from":  "auto",
// 			"to":    "zh",
// 			"appid": appID,
// 			"salt":  salt,
// 			"sign":  sign,
// 		}

// 		err := json.Unmarshal([]byte(httpGet(apiURL, params).Content), &bres)
// 		if err == nil {
// 			break
// 		} else {
// 			sleep(3)
// 			continue
// 		}
// 	}

// 	if bres.ErrorCode != "" && len(bres.TransResult) == 0 {
// 		Panicerr("翻译出错, 状态码 " + bres.ErrorCode + ", 原因为 " + bres.ErrorMsg)
// 	}
// 	return bres.TransResult[0].Dst
// }

// func en2zh(text string) (res string) {
// 	res = baiduTranslateAnyToZH(text)
// 	return
// }
