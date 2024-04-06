package alipayGlobal

import (
	"crypto/rsa"
	"github.com/lgrisa/library/pay"
	"github.com/lgrisa/library/utils"
	"github.com/lgrisa/library/utils/xrsa"
	"github.com/pkg/errors"
)

type GlobalAliPayClient struct {
	clientId           string
	privateKey         *rsa.PrivateKey //应用私钥(生成签名)
	aliPublicKey       string          //支付宝公钥(验证支付宝返回数据签名)
	isProd             bool
	webNotifyUrl       string
	webRefundNotifyUrl string
	DebugSwitch        pay.DebugSwitch // Debug 开关
}

// NewClient 初始化支付宝客户端
// clientId：应用ID
// privateKey：应用私钥，支持PKCS1和PKCS8
// aliPublicKey：支付宝公钥，支持PKCS1和PKCS8
// isProd：是否是正式环境，沙箱环境请选择新版沙箱应用。
// webNotifyUrl：支付宝支付成功后通知商户的地址
// webRefundNotifyUrl：支付宝退款成功后通知商户的地址
func NewClient(clientId, privateKey, aliPublicKey string, webNotifyUrl, webRefundNotifyUrl string, isProd bool) (*GlobalAliPayClient, error) {
	if clientId == utils.NULL || privateKey == utils.NULL || aliPublicKey == utils.NULL {
		return nil, MissAlipayInitParamErr
	}
	key := xrsa.FormatAlipayPrivateKey(privateKey)
	priKey, err := utils.DecodePrivateKey([]byte(key))
	if err != nil {
		return nil, errors.Errorf("InitAliPay xpem.DecodePrivateKey(%s)：%v", key, err)
	}

	return &GlobalAliPayClient{
		clientId:           clientId,
		privateKey:         priKey,
		aliPublicKey:       xrsa.FormatAlipayPublicKey(aliPublicKey),
		isProd:             isProd,
		webNotifyUrl:       webNotifyUrl,
		webRefundNotifyUrl: webRefundNotifyUrl,
	}, nil
}

func (a *GlobalAliPayClient) isDebugMode() bool {
	return a.DebugSwitch == pay.DebugOn
}
