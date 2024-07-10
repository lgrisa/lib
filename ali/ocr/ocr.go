package ocr

import (
	"encoding/json"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	ocr_api20210707 "github.com/alibabacloud-go/ocr-api-20210707/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/lgrisa/lib/utils"
	"github.com/pkg/errors"
	"os"
	"strings"
)

// Description:
//
// 使用AK&SK初始化账号Client
//
// @return Client
//
// @throws Exception

type AliOcrClient struct {
	AccessKeyId     string
	AccessKeySecret string

	client *ocr_api20210707.Client
}

func NewAliOcrClient(accessKeyId, accessKeySecret string) (*AliOcrClient, error) {

	if accessKeyId == "" || accessKeySecret == "" {
		return nil, errors.New("accessKeyId or accessKeySecret is empty")
	}

	client, err := createClient(accessKeyId, accessKeySecret)

	if err != nil {
		return nil, errors.Errorf("createClient => %v", err)
	}

	return &AliOcrClient{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		client:          client,
	}, nil
}

func createClient(accessKeyId, accessKeySecret string) (*ocr_api20210707.Client, error) {
	// 工程代码泄露可能会导致 AccessKey 泄露，并威胁账号下所有资源的安全性。以下代码示例仅供参考。
	// 建议使用更安全的 STS 方式，更多鉴权访问方式请参见：https://help.aliyun.com/document_detail/378661.html。
	config := &openapi.Config{
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_ID。
		AccessKeyId: tea.String("LTAI5tRpQYu9tz9gT6hVLPMT"),
		// 必填，请确保代码运行环境设置了环境变量 ALIBABA_CLOUD_ACCESS_KEY_SECRET。
		AccessKeySecret: tea.String("4ifEMUWxBe9JxOI36MmoefS1ssQkqG"),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/ocr-api
	config.Endpoint = tea.String("ocr-api.cn-hangzhou.aliyuncs.com")
	return ocr_api20210707.NewClient(config)
}

func (c *AliOcrClient) OcrVerificationCode(path string) (string, error) {
	// 需要安装额外的依赖库，直接点击下载完整工程即可看到所有依赖。
	//bodyStream := stream.ReadFromFilePath(tea.String("captcha.png"))

	f, err := os.Open(path)

	if err != nil {
		return utils.NULL, errors.Errorf("打开文件失败：%v", err)
	}

	recognizeAdvancedRequest := &ocr_api20210707.RecognizeAdvancedRequest{
		Body: f,
	}
	runtime := &util.RuntimeOptions{}
	resp, tryErr := func() (resp *ocr_api20210707.RecognizeAdvancedResponse, _e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		resp, _e = c.client.RecognizeAdvancedWithOptions(recognizeAdvancedRequest, runtime)
		if _e != nil {
			return nil, _e
		}

		return resp, nil
	}()

	if tryErr != nil {
		errSDK := &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			errSDK = _t
		} else {
			errSDK.Message = tea.String(tryErr.Error())
		}
		// 此处仅做打印展示，请谨慎对待异常处理，在工程项目中切勿直接忽略异常。
		// 错误 message
		fmt.Println(tea.StringValue(errSDK.Message))
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(errSDK.Data)))
		errDecode := d.Decode(&data)
		if errDecode != nil {
			return utils.NULL, errors.Errorf("json decode error: %v", errDecode)
		}

		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println(recommend)
		}

		_, errAssert := util.AssertAsString(errSDK.Message)
		if errAssert != nil {
			return utils.NULL, errors.Errorf("assert message error: %v", errAssert)
		}

		return utils.NULL, errors.Errorf("SDKError => %v", errSDK)
	} else {
		// 复制代码运行请自行打印 API 的返回值
		respData := &ocr_api20210707.RecognizeAllTextResponseBodyData{}

		if err = json.Unmarshal([]byte(tea.StringValue(resp.Body.Data)), respData); err != nil {
			return utils.NULL, errors.Errorf("Data json.Unmarshal => %v", err)
		}

		return *respData.Content, nil
	}
}
