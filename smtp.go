package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"sort"
	"strings"
)

//with display name
//Jane Doe <jane.doe@gmail.com>
//dig google.com MX

func sendText(dao Dao, from string, to string, subject string, body string) error {
	return sendEmail(dao, from, to, subject, "text/plain", body)
}

func sendEmail(dao Dao, from string, to string, subject string, mime string, body string) error {
	fromAddress, err := mail.ParseAddress(from)
	if err != nil {
		return err
	}
	toAddress, err := mail.ParseAddress(to)
	if err != nil {
		return err
	}
	// hostname, err := os.Hostname()
	// if err != nil {
	// 	return err
	// }
	toDomain := strings.Split(toAddress.Address, "@")[1]
	fromDomain := strings.Split(fromAddress.Address, "@")[1]
	fromDomainDro, err := dao.GetDomainPrivateKey(fromDomain)
	if err != nil {
		return err
	}
	email, bodyLength := composeMimeMail(toAddress.String(), fromAddress.String(), subject, mime, body)
	err = dkimSign(&email, bodyLength, fromDomain, []byte(fromDomainDro.PrivateKey))
	if err != nil {
		return err
	}
	log.Println(string(email))
	mxs, err := net.LookupMX(toDomain)
	if err != nil {
		return err
	}
	sort.Slice(mxs, func(i, j int) bool {
		return mxs[i].Pref < mxs[j].Pref
	})
	for _, x := range mxs {
		log.Println(x.Host, x.Pref)
	}
	err = nil
	for _, x := range mxs {
		addr := fmt.Sprintf("%s:25", x.Host)
		// conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
		// if err != nil {
		// 	continue
		// }
		// c, err := smtp.NewClient(conn, hostname)
		// if err != nil {
		// 	return "", err
		// }
		// defer c.Quit()
		//FIXME add timeout
		err = smtp.SendMail(addr, nil, fromAddress.Address, []string{toAddress.Address}, email)
		if err == nil {
			break
		}
		//FIXME retry only if network error
		log.Println(x.Host, x.Pref, err)
		break
	}
	return err
}

func sendMail(c smtp.Client, from string, to string, msg []byte) error {
	defer c.Quit()
	err := c.Mail(from)
	if err != nil {
		return err
	}
	err = c.Rcpt(to)
	if err != nil {
		return err
	}
	wc, err := c.Data()
	if err != nil {
		return err
	}
	count, err := wc.Write(msg)
	if err != nil {
		return err
	}
	if count != len(msg) {
		return fmt.Errorf("msg len mismatch %d != %d", count, len(msg))
	}
	return nil
}

func encodeRFC2047(str string) string {
	// use mail's rfc2047 to encode any string
	addr := mail.Address{Name: str}
	return strings.Trim(addr.String(), " <>")
}

func composeMimeMail(to string, from string, subject string, mime string, body string) ([]byte, uint) {
	var b strings.Builder
	fmt.Fprintf(&b, "%s: %s\r\n", "From", from)
	fmt.Fprintf(&b, "%s: %s\r\n", "To", to)
	fmt.Fprintf(&b, "%s: %s\r\n", "Subject", encodeRFC2047(subject))
	fmt.Fprintf(&b, "%s: %s\r\n", "MIME-Version", "1.0")
	fmt.Fprintf(&b, "%s: %s\r\n", "Content-Type", fmt.Sprintf("%s; charset=\"utf-8\"", mime))
	fmt.Fprintf(&b, "%s: %s\r\n", "Content-Transfer-Encoding", "base64")
	bytes := []byte(body)
	b.WriteString("\r\n")
	b.WriteString(base64.StdEncoding.EncodeToString(bytes))
	return []byte(b.String()), uint(len(bytes))
}
