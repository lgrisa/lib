解析支付宝全球支付的通知

1.解析对应HTTP Header

2.解析对应的Notify 结构体

gin框架写法

付款通知解析示例：
```go
if err :=client.VerifySign(context); err != nil {
    log.Printf("verify sign failed, err:%v", err)
    return
}

notifyBody, err := context.GetRawData()

if err != nil {
    log.Printf("get raw data failed, err:%v", err)
    return
}

notifyReq := &PayNotify{}
err = json.Unmarshal(notifyBody, notifyReq)
if err != nil {
    log.Printf("unmarshal failed, err:%v", err)
    return
}

if notifyReq.NotifyType != "PAYMENT_RESULT" {
	log.Printf("notify type is not PAYMENT_RESULT")
    return
}

transactionId := notifyReq.PaymentId //支付宝交易号
orderId := notifyReq.PaymentRequestId //商户订单号
status := notifyReq.Result.ResultStatus //支付状态

if if status == "" || transactionId == "" || orderId == "" {
	log.Printf("status or transactionId or orderId is empty")
	err = client.WritePayNotifyResp(context,"fail")
	if err != nil {
		log.Printf("write pay notify resp failed, err:%v", err)
	}
    return
}

if status == "S" {
	    //支付成功
} else {
	//其他状态
	return
}

err = client.WritePayNotifyResp(context,"")
if err != nil {
    log.Printf("write pay notify resp failed, err:%v", err)
}
```

退款通知解析示例：
```go
if err := client.VerifySign(ctx); err != nil {
    log.Printf("verify sign failed, err:%v", err)
    return
}

data, err := ctx.GetRawData()
if err != nil {
    log.Printf("get raw data failed, err:%v", err)
    return
}

notify := &RefundNotify{}
if err = json.Unmarshal(data, notify); err != nil {
    log.Printf("unmarshal failed, err:%v", err)
    return
}

if notify.Result.ResultStatus != "S" {
    log.Printf("refund failed")
    return
}

err = client.writeAliRefundResp(ctx, "success", "")
if err != nil {
    log.Printf("write refund resp failed, err:%v", err)
}
```