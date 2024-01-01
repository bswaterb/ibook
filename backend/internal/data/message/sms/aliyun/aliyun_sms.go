package aliyun

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"ibook/internal/conf"
	"ibook/internal/service/message/sms"
)

type smsRepo struct {
	client       *dysmsapi20170525.Client
	aliyunConfig *conf.AliyunSMS
}

func NewAliyunSMSRepo(aliyunConfig *conf.AliyunSMS) sms.SMSRepo {
	return &smsRepo{
		client:       createDYSMSClient(aliyunConfig.AccessKeyId, aliyunConfig.AccessKeySecret),
		aliyunConfig: aliyunConfig,
	}
}

func createDYSMSClient(AccessKeyId, AccessKeySecret string) *dysmsapi20170525.Client {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(AccessKeyId),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(AccessKeySecret),
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	result, err := dysmsapi20170525.NewClient(config)
	if err != nil {
		return nil
	}
	return result
}

func (s smsRepo) SendMessage(ctx *gin.Context, tplId string, phoneNumbers []string, args []sms.MsgArgs) error {
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers: tea.String(phoneNumbers[0]),
		SignName:     tea.String(s.aliyunConfig.SignName),
	}
	_, err := s.client.SendSmsWithOptions(sendSmsRequest, &util.RuntimeOptions{})
	if err != nil {
		return err
	}
	return nil
}
