package golanglibs

import "testing"

func TestJieba(t *testing.T) {
	j := Tools.Jieba()
	Lg.Debug(j.Cut("《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄制作的科幻片."))
	Lg.Debug(j.Cut("Cryptocurrency is a digital payment system that doesn't rely on banks to verify transactions. It’s a peer-to-peer system that can enable anyone anywhere to send and receive payments."))
}
