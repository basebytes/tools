package tools

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"path"
	"strings"
)

const (
	mimeVersion="1.0"
	encodingBase64="base64"
	encodingQuotedPrintable="quoted-printable"
	contentTypeText="text/plain; charset=UTF-8"
	contentTypeHtml="text/html; charset=UTF-8"
	contentType8bit="application/octet-stream"

)
var crlf =[]byte{'\r','\n'}

type email struct {
	auth       smtp.Auth
	user       *mail.Address
	server     string
	header     map[string]string
	address    []string
}

func NewEmailSender(user,password,server string) (*email,error){
	u,err:=mail.ParseAddress(user)
	if err!=nil{
		return nil,err
	}
	e:=&email{
		server:server,
		user: u,
		auth: smtp.PlainAuth("",u.Address,password,strings.Split(server,":")[0]),
	}
	return e,nil
}

func UserName(name string) EmailOption{
	return func(email *email) {
		email.header["From"]=(&mail.Address{Name: name,Address: email.user.Address}).String()
	}
}

func ReturnPath(returnPath string) EmailOption{
	return func(email *email) {
		email.header["Return-Path"]=returnPath
	}
}

func ReplyTo(replyTo string) EmailOption{
	return func(email *email) {
		email.header["Reply-To"]=replyTo
	}
}

func Cc(cc string)EmailOption{
	return func(email *email) {
		if ccAddr,err:=mail.ParseAddressList(cc);err==nil{
			email.header["Cc"]= addressSlice(ccAddr).String()
			email.address=append(email.address, addressSlice(ccAddr).Addr()...)
		}
	}
}

func Bcc(bcc string)EmailOption{
	return func(email *email) {
		if bccAddr,err:=mail.ParseAddressList(bcc);err==nil{
			email.address=append(email.address, addressSlice(bccAddr).Addr()...)
		}
	}
}

func Subject(subject string) EmailOption{
	return func(email *email) {
		email.header["Subject"]=subject
	}
}

func HtmlMessage(format string,args ...interface{}) *Message {
	return newMessage(
		contentTypeHtml,
		encodingQuotedPrintable,
		map[string][]string{
			"Content-Disposition":{"inline"},
		},
		fmt.Sprintf(format,args...),
	)
}

func TextMessage(format string,args ...interface{}) *Message {
	return newMessage(
		contentTypeText,
		encodingBase64,
		map[string][]string{
			"Content-Disposition":{"inline"},
		},
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(format,args...))),
	)
}

func AttachmentMessage(name,filePath,contentId string)*Message {
	header:=map[string][]string{
		"Content-Disposition":{"attachment;filename="+name},
	}
	if contentId!=""{
		header["Content-ID"]=[]string{contentId}
	}

	contentType:=contentType8bit
	if t:=mime.TypeByExtension(path.Ext(name));t!=""{
		contentType=t
	}
	content:=""
	if attData, err := ioutil.ReadFile(filePath);err==nil{
		content=base64.StdEncoding.EncodeToString(attData)
	}
	return newMessage(
		contentType,
		encodingBase64,
		header,
		content,
	)
}

func (e *email)Send(to string,messages []*Message,options ...EmailOption) error{
	toAddr,err:=mail.ParseAddressList(to)
	if err!=nil{
		return err
	}
	buffer := bytes.NewBuffer(nil)
	w:=multipart.NewWriter(buffer)
	e.initParams(toAddr,options,w.Boundary())
	e.writeHeader(buffer, e.header)
	for _,message:=range messages{
		iw,err:=w.CreatePart(message.h)
		if err==nil{
			_,err =iw.Write([]byte(message.c))
			_,_=iw.Write(crlf)
		}
		if err!=nil{
			return err
		}
	}
	_=w.Close()
	return smtp.SendMail(e.server,e.auth,e.user.Address,e.address,buffer.Bytes())
}

type Message struct{
	h textproto.MIMEHeader
	c string
}

func (m *Message)header() textproto.MIMEHeader{
	return m.h
}
func (m *Message)content() string{
	return m.c
}

func (e *email) writeHeader(buffer *bytes.Buffer, Header map[string]string) {
	for key, value := range Header {
		_, _ = fmt.Fprintf(buffer, "%s: %s\r\n", key, value)
	}
	buffer.Write(crlf)
}

func (e *email) initParams(to addressSlice,options []EmailOption,boundary string){
	e.header=map[string]string{
		"From":e.user.String(),
		"To":to.String(),
		"Mime-Version":mimeVersion,
		"Content-Type":"multipart/mixed;boundary="+boundary,
		//multipart/alternative
	}
	e.address=to.Addr()
	for _,option:=range options{
		option(e)
	}
}

func newMessage(contentType,encoding string,other map[string][]string,content string) *Message {
	m:=&Message{h: textproto.MIMEHeader{},c:content}
	m.h.Set("Content-Type",contentType)
	m.h.Set("Content-Transfer-Encoding",encoding)
	if other!=nil{
		for k,V:=range other{
			if len(k)==0||len(V)==0{
				continue
			}
			for _,v:=range V{
				if len(v)==0{
					continue
				}
				m.h.Add(k,v)
			}
		}
	}
	return m
}


type EmailOption func(email *email)

type addressSlice []*mail.Address

func (a addressSlice)String() string{
	buffer :=bytes.NewBuffer(nil)
	for _,add:=range a{
		buffer.WriteString(add.Address+",")
	}
	return buffer.String()
}

func (a addressSlice)Addr() []string{
	var addr[]string
	for _,add:=range a{
		addr=append(addr,add.Address)
	}
	return addr
}