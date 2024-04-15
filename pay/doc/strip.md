* <font color='#003087' size='4'>API</font>
    * 初始化：`NewClient()`
    * 创建订单：`client.CreateOrder()`
    * 验签：`client.WebHookVerify()`
    * 通知发货checkout.session.completed事件：`client.CheckoutSessionCompleted()`
    * 通知发货checkout.session.async_payment_succeeded事件：`client.CheckoutSessionCompleted()`
    * 通知退款事件：`client.EventCheckoutOrderRefunded()`
### Webhooks时间处理流程

gin框架写法
```go
rawData, err := context.GetRawData()
if err != nil {
	xlog.Error(err)
    return
}

stripEvent,err := client.WebHookVerify(context)
if err != nil {
    xlog.Error(err)
    return
}

if stripEvent.Type == stripe.PayCompleted {
    // 通知发货
    client.CheckoutSessionCompleted(context, stripEvent)
} else if stripEvent.Type == stripe.PayAsyncSucceeded {
    // 通知发货
    client.CheckoutSessionCompleted(context, stripEvent)
} else if stripEvent.Type == stripe.REFUNDED {
    // 通知退款
    client.CheckoutSessionRefunded(context, stripEvent)
}else{
    // do something
    context.Status(http.StatusBadRequest)
    return
}

context.Status(http.StatusOK)
```