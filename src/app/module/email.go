package module

import (
	"fmt"
	"net/smtp"
)

func SendVerification(env *Environments, toAddress string, code string) {
	to := []string{
		toAddress,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte(fmt.Sprintf("From: %s\r\n", env.GMAIL_ID) +
		fmt.Sprintf("To: %s\r\n", toAddress) +
		"Subject: 【京大楽単bot】認証リンクのお知らせ\r\n\r\n" +
		"京大楽単botをご利用いただきありがとうございます。\n\n" +
		"過去問閲覧機能有効化のための認証リンクをお送りします。アクセスしていただいたのち、過去問閲覧機能が有効になります。\n\n\n" +
		fmt.Sprintf("【認証リンク】\nhttps://pf.das82.com/verification?code=%s \n\n\n", code) +
		"----------\n京大楽単bot運営\nお問い合わせ：support@das82.com" +
		"\r\n")

	// Authentication.
	fmt.Println(env.GMAIL_PASSWORD, env.GMAIL_ID)
	auth := smtp.PlainAuth("", env.GMAIL_ID, env.GMAIL_PASSWORD, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, env.GMAIL_ID, to, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}
