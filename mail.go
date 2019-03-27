package mail

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"
	"time"
)

type Mail interface {
	Auth()
	Send(message Message) error
}

type MailConf struct{
	User		string
	Password	string
	Host 		string
	Port        string
	AuthContent smtp.Auth
}

//附件
type Attachment struct {
	Name    		string
	ContentType 	string
	WithFile    	bool
}

type Message struct {
	From		string
	To          []string
	Cc          []string
	Bcc         []string
	Subject		string
	Body 		string
	ContentType string
	Attachment  Attachment
}

func (mail *MailConf) Auth(){
	mail.AuthContent = smtp.PlainAuth("",mail.User,mail.Password,mail.Host)
}

func (mail *MailConf) Send(message Message) error {
	mail.Auth()
	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary" //边界线
	Header := make(map[string]string)
	Header["From"] = message.From
	Header["To"] = strings.Join(message.To, ";")
	Header["Cc"] = strings.Join(message.Cc, ";")
	Header["Bcc"] = strings.Join(message.Bcc, ";")
	Header["Subject"] = message.Subject
	Header["Content-Type"] = "multipart/mixed;boundary=" + boundary
	Header["Mime-Version"] = "1.0"
	Header["Date"] = time.Now().String()
	mail.writeHeader(buffer, Header)

	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Type:" + message.ContentType + "\r\n"
	body += "\r\n" + message.Body + "\r\n"
	buffer.WriteString(body)

	if message.Attachment.WithFile {
		attachment := "\r\n--" + boundary + "\r\n"
		attachment += "Content-Transfer-Encoding:base64\r\n"
		attachment += "Content-Disposition:attachment\r\n"
		attachment += "Content-Type:" + message.Attachment.ContentType + ";name=\"" + message.Attachment.Name + "\"\r\n"
		buffer.WriteString(attachment)
		defer func() {
			if err := recover(); err != nil {
				log.Fatalln(err)
			}
		}()
		mail.writeFile(buffer, message.Attachment.Name)
	}

	buffer.WriteString("\r\n--" + boundary + "--")
	smtp.SendMail(mail.Host+":"+mail.Port, mail.AuthContent, message.From, message.To, buffer.Bytes())
	return nil
}

//set header
func (mail *MailConf) writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + "\r\n"
	}
	header += "\r\n"
	buffer.WriteString(header)
	return header
}

//set file
func (mail *MailConf) writeFile(buffer *bytes.Buffer, fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	buffer.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}
}

