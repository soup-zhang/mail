package mail

import (
	"testing"
)

func TestMAIL(t *testing.T){
	var mc MailConf
	mc.User = "39xxxx758@qq.com"
	mc.Password = "ibbhxxxxxxbgjg"
	mc.Host = "smtp.qq.com"
	mc.Port = "587"

	var message Message
	message.From = "39xxxx758@qq.com"
	message.To = []string{"soup-zhang@126.com.com"}
	//message.To = []string{"soup-zhang1@126.com.com","soup-zhang2@126.com.com"}
	//message.Cc = []string{"soup-zhang3@126.com.com"}
	message.Subject = "测试主题"
	message.ContentType = "text/plain;charset=utf-8"  //"text/html;charset=utf-8"
	message.Body = "测试邮件内容"

	//附件信息
	//message.Attachment.Name = "static/image/test.jpg"
	//message.Attachment.ContentType = "image/jpg"
	//message.Attachment.WithFile = true

	mc.Send(message)
}