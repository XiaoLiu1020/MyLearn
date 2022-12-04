- [参考项目](#参考项目)
- [ 可以查看各个test文件学习](#-可以查看各个test文件学习)
	- [微信相关](#微信相关)

# 参考项目
<https://github.com/go-pay/gopay>

里面有各个支付的文档目录:

*   [Alipay](https://github.com/go-pay/gopay/blob/main/doc/alipay.md)
*   [Wechat](https://github.com/go-pay/gopay/blob/main/doc/wechat_v3.md)
*   [QQ](https://github.com/go-pay/gopay/blob/main/doc/qq.md)
*   [Paypal](https://github.com/go-pay/gopay/blob/main/doc/paypal.md)
*   [Apple](https://github.com/go-pay/gopay/blob/main/doc/apple.md)

# &#x20;可以查看各个test文件学习

例如支付宝, 支付宝支付需要支付应用id,应用密钥, 支付的证书,回调地址等

1.  首先初始化`Client`

```go
// FROM: https://github.com/go-pay/gopay/blob/main/doc/alipay.md
func TestMain(m *testing.M) {

	// 初始化支付宝客户端
	//    appid：应用ID
	//    privateKey：应用私钥，支持PKCS1和PKCS8
	//    isProd：是否是正式环境
	client, err = NewClient(cert.Appid, cert.PrivateKey, false)
	if err != nil {
		xlog.Error(err)
		return
	}
	// Debug开关，输出/关闭日志
	client.DebugSwitch = gopay.DebugOff

	// 配置公共参数
	client.SetCharset("utf-8").
		SetSignType(RSA2).
		// SetAppAuthToken("")
		SetReturnUrl("https://www.fmm.ink").
		SetNotifyUrl("https://www.fmm.ink")

	// 自动同步验签（只支持证书模式）
	// 传入 alipayCertPublicKey_RSA2.crt 内容
	client.AutoVerifySign(cert.AlipayPublicContentRSA2)

	// 传入证书内容
	err := client.SetCertSnByContent(cert.AppPublicContent, cert.AlipayRootContent, cert.AlipayPublicContentRSA2)
	// 传入证书文件路径
	//err := client.SetCertSnByPath("cert/appCertPublicKey_2021000117673683.crt", "cert/alipayRootCert.crt", "cert/alipayCertPublicKey_RSA2.crt")
	if err != nil {
		xlog.Debug("SetCertSn:", err)
		return
	}
	os.Exit(m.Run())
}
```

2.  根据业务进行对接接口

*   统一收单交易支付接口（商家扫用户付款码）：client.TradePay()
*   统一收单线下交易预创建（用户扫商品收款码）：client.TradePrecreate()
*   手机网站支付接口2.0（手机网站支付）：client.TradeWapPay()
*   统一收单下单并支付页面接口（电脑网站支付）：client.TradePagePay()
*   统一收单交易创建接口（小程序支付）：client.TradeCreate()

3.  最后监听支付成功后的回调

返回的相关信息可以查看: https://opendocs.alipay.com/open/203/105286

## 微信相关
查看介绍: https://github.com/go-pay/gopay/blob/main/doc/wechat_v3.md