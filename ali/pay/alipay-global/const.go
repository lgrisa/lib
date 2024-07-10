package alipayGlobal

const (
	AliUrlSandbox = "https://open-sea-global.alipay.com"
	AliUrl        = "https://open-na-global.alipay.com"

	payApi        = "/ams/api/v1/payments/pay"
	payApiSandbox = "/ams/sandbox/api/v1/payments/pay"

	inquiryPaymentApi        = "/ams/api/v1/payments/inquiryPayment"
	inquiryPaymentSandboxApi = "/ams/sandbox/api/v1/payments/inquiryPayment"

	refundApi        = "/ams/api/v1/payments/refund"
	refundApiSandbox = "/ams/sandbox/api/v1/payments/refund"
)

//沙盒测试地址：https://global.alipay.com/docs/ac/ref/testwallet

//支持货币类型：https://global.alipay.com/docs/ac/ref/cc#ONkIe

var CurrencyType = map[string]string{
	"USD": "美元",
	"PHP": "菲律宾比索",
	"IDR": "印尼盾",
	"KRW": "韩元",
	"THB": "泰铢",
	"HKD": "港元",
	"MYR": "马来西亚林吉特",
	"CNY": "人民币",
	"BDT": "孟加拉塔卡",
	"PKR": "巴基斯坦卢比",
	"TWD": "新台币",
}
