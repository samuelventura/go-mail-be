package main

import (
	"encoding/base64"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"sort"
	"strings"
	"time"
)

//with display name
//User Name <user.name@domain.tld>
//dig google.com MX

func mailSend(dao Dao, id string, from string, to string, subject string, mime string, body string) error {
	mdro := &MessageDro{Mid: id,
		From: from, To: to,
		Subject: subject, Mime: mime,
		Body: body, Created: time.Now()}
	err := dao.AddMessage(mdro)
	if err != nil {
		return err
	}
	fromAddress, err := mail.ParseAddress(from)
	if err != nil {
		return err
	}
	toAddress, err := mail.ParseAddress(to)
	if err != nil {
		return err
	}
	toDomain := strings.Split(toAddress.Address, "@")[1]
	fromDomain := strings.Split(fromAddress.Address, "@")[1]
	fromDomainDro, err := dao.GetDomain(fromDomain)
	if err != nil {
		return err
	}
	email, bodyLength := composeMimeMail(id, toAddress.String(), fromAddress.String(), subject, mime, body)
	err = dkimSign(&email, bodyLength, fromDomain, []byte(fromDomainDro.PrivateKey))
	if err != nil {
		return err
	}
	//log.Println(string(email))
	mxs, err := net.LookupMX(toDomain)
	if err != nil {
		return err
	}
	sort.Slice(mxs, func(i, j int) bool {
		return mxs[i].Pref < mxs[j].Pref
	})
	mxsn := make([]string, 0, len(mxs))
	for _, x := range mxs {
		mxsn = append(mxsn, x.Host)
		addr := fmt.Sprintf("%s:25", x.Host)
		dial := false
		err = smtpSend(addr, fromAddress.Address,
			[]string{toAddress.Address}, email, &dial)
		result := fmt.Sprintf("host:%s dial:%v error:%v", x.Host, dial, err)
		adro := &AttemptDro{Mid: id, Created: time.Now(), Result: result}
		err2 := dao.AddAttempt(adro)
		if err2 != nil {
			return err2
		}
		if dial {
			continue
		}
		return err
	}
	return fmt.Errorf("no working mx %v", mxsn)
}

func escapeHeader(str string) string {
	return mime.QEncoding.Encode("utf-8", str)
}

func composeMimeMail(id string, to string, from string, subject string, mime string, body string) ([]byte, uint) {
	var b strings.Builder
	fmt.Fprintf(&b, "%s: %s\r\n", "Message-Id", id)
	fmt.Fprintf(&b, "%s: %s\r\n", "Date", time.Now().String())
	fmt.Fprintf(&b, "%s: %s\r\n", "From", from)
	fmt.Fprintf(&b, "%s: %s\r\n", "To", to)
	fmt.Fprintf(&b, "%s: %s\r\n", "Subject", escapeHeader(subject))
	fmt.Fprintf(&b, "%s: %s\r\n", "MIME-Version", "1.0")
	fmt.Fprintf(&b, "%s: %s\r\n", "Content-Type", fmt.Sprintf("%s; charset=\"utf-8\"", mime))
	fmt.Fprintf(&b, "%s: %s\r\n", "Content-Transfer-Encoding", "base64")
	bytes := []byte(body)
	b.WriteString("\r\n")
	b.WriteString(base64.StdEncoding.EncodeToString(bytes))
	return []byte(b.String()), uint(len(bytes))
}
