* <font color='#003087' size='4'>API</font>
    * 初始化：`NewClient()`
    * 验签：`client.VerifySign()`
    * 回复通知：`client.WriteNotifyResp()`
    * 通知发货Approve事件：`client.EventCheckoutOrderApproved()`
    * 通知发货Complete事件：`client.EventCheckoutOrderComplete()`
    * 通知退款事件：`client.EventCheckoutOrderRefunded()`

### Webhooks时间处理流程

gin框架写法
```go
rawData, err := context.GetRawData()
if err != nil {
	xlog.Error(err)
    return
}

// 验签
if err = client.WebHookVerifySign(context);err!=nil{
    xlog.Error(err)
    return
}

resp:=&paypal.WebhookNotifyResponse{}
err = json.Unmarshal(rawData, &paypalResponse)
if err != nil {
    xlog.Error(err)
    return
}

if resp.EventType == paypal.WebhookEventCheckoutOrderApproved {
    // 通知发货
    client.EventCheckoutOrderApproved(context, resp)
} else if resp.EventType == paypal.WebhookEventCheckoutOrderComplete {
	// 通知发货
    client.EventCheckoutOrderComplete(context, resp)
} else if resp.EventType == paypal.WebhookEventCheckoutOrderRefunded {
    // 通知退款
    client.EventCheckoutOrderRefunded(context, resp)
}else{
    // do something
}

client.WriteNotifyResp(context)
return
```
其余部分 github gopay存在对应实现 [参见](https://github.com/go-pay/gopay)