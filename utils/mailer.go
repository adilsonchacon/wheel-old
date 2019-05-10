// utils.MailerSetFrom("", "no-reply@blazesecurity.com.br")
// utils.MailerAddTo("Adilson Chacon", "adilsonchacon@gmail.com")
// utils.MailerAddCc("Eu 2", "achacon@gmail.com")
// utils.MailerAddBcc("Eu 2", "chacon@blazeinfosec.com")
// Send as plain text
// utils.MailerSend("oi de novo com arquivo de configuracao", "mais um teste e q saco", false)
// OR as HTML
// utils.MailerSend("agora com html", "<!DOCTYPE html><html><head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=UTF-8\" /><style type=\"text/css\">h1 {color: red; font-size: 24px; font-weight: bold;}</style></head><body><h1>HI!</h1></body></html>", true)

package utils

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"net/mail"
	"net/smtp"
	"regexp"
	"strconv"
	"strings"
)

type mailerConfig struct {
	User     string
	Name     string
	Password string
	Address  string
	Port     int
}

var (
	from mail.Address
	to   []mail.Address
	cc   []mail.Address
	bcc  []mail.Address
)

func MailerSetFrom(name string, tEmail string) {
	validEmail := regexp.MustCompile(`\A[^@]+@([^@\.]+\.)+[^@\.]+\z`)

	if validEmail.MatchString(tEmail) {
		from = mail.Address{name, tEmail}
	} else {
		LoggerError().Println("email is invalid")
	}
}

func MailerAddTo(name string, tEmail string) {
	validEmail := regexp.MustCompile(`\A[^@]+@([^@\.]+\.)+[^@\.]+\z`)

	if validEmail.MatchString(tEmail) {
		to = append(to, mail.Address{name, tEmail})
	} else {
		LoggerError().Println("email is invalid")
	}
}

func MailerAddCc(name string, tEmail string) {
	validEmail := regexp.MustCompile(`\A[^@]+@([^@\.]+\.)+[^@\.]+\z`)

	if validEmail.MatchString(tEmail) {
		cc = append(cc, mail.Address{name, tEmail})
	} else {
		LoggerError().Println("email is invalid")
	}
}

func MailerAddBcc(name string, tEmail string) {
	validEmail := regexp.MustCompile(`\A[^@]+@([^@\.]+\.)+[^@\.]+\z`)

	if validEmail.MatchString(tEmail) {
		bcc = append(bcc, mail.Address{name, tEmail})
	} else {
		LoggerError().Println("email is invalid")
	}
}

func MailerResetAll() {
	MailerSetFrom("", "")
	to = to[:0]
	cc = cc[:0]
	bcc = bcc[:0]
}

func MailerSend(subject string, body string, html bool) {
	receipts := mailerReciepts()
	headers := make(map[string]string)
	config := mailerLoadConfigFile()

	if from.String() == "<@>" {
		from = mail.Address{config.Name, config.User}
	}

	// Setup header
	headers["From"] = from.String()

	if len(to) > 0 {
		headers["To"] = mailerStringfy("to")
	}

	if len(cc) > 0 {
		headers["Cc"] = mailerStringfy("cc")
	}

	if len(bcc) > 0 {
		headers["Bcc"] = mailerStringfy("bcc")
	}

	if html {
		headers["MIME-version"] = "1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	}

	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := config.Address + ":" + strconv.Itoa(config.Port)
	host, _, _ := net.SplitHostPort(servername)
	auth := smtp.PlainAuth("", config.User, config.Password, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		LoggerFatal().Panic(err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		LoggerFatal().Panic(err)
	}

	// Auth
	if err = client.Auth(auth); err != nil {
		LoggerFatal().Panic(err)
	}

	// To
	if err = client.Mail(from.Address); err != nil {
		LoggerFatal().Panic(err)
	}

	// To, Cc and Bcc
	for i := 0; i < len(receipts); i++ {
		if err = client.Rcpt(receipts[0]); err != nil {
			LoggerFatal().Panic(err)
		}
	}

	// Data
	socket, err := client.Data()
	if err != nil {
		LoggerFatal().Panic(err)
	}

	_, err = socket.Write([]byte(message))
	if err != nil {
		LoggerFatal().Panic(err)
	}

	err = socket.Close()
	if err != nil {
		LoggerFatal().Panic(err)
	}

	client.Quit()
}

func mailerStringfy(tType string) string {
	var isTo = regexp.MustCompile(`(?i)\Ato\z`)
	var isCc = regexp.MustCompile(`(?i)\Acc\z`)
	var isBcc = regexp.MustCompile(`(?i)\Abcc\z`)
	var isFrom = regexp.MustCompile(`(?i)\Afrom\z`)
	var auxEmail []mail.Address
	var auxString []string

	if isTo.MatchString(tType) {
		auxEmail = to
	} else if isCc.MatchString(tType) {
		auxEmail = cc
	} else if isBcc.MatchString(tType) {
		auxEmail = bcc
	} else if isFrom.MatchString(tType) {
		auxEmail = append(auxEmail, from)
	} else {
		LoggerError().Println("Invalid param for mailerStringfy, available params are to, cc, bcc and from")
	}

	for i := 0; i < len(auxEmail); i++ {
		auxString = append(auxString, auxEmail[i].Address)
	}

	return strings.Join(auxString, "; ")
}

func mailerReciepts() []string {
	var reciepts []string

	for i := 0; i < len(to); i++ {
		reciepts = append(reciepts, to[0].Address)
	}

	for i := 0; i < len(cc); i++ {
		reciepts = append(reciepts, cc[0].Address)
	}

	for i := 0; i < len(bcc); i++ {
		reciepts = append(reciepts, bcc[0].Address)
	}

	return reciepts
}

func mailerLoadConfigFile() mailerConfig {
	config := mailerConfig{}

	err := yaml.Unmarshal(mailerReadConfigFile(), &config)
	if err != nil {
		LoggerError().Fatalf("error: %v", err)
	}

	return config
}

func mailerReadConfigFile() []byte {
	data, err := ioutil.ReadFile("./config/email.yml")
	if err != nil {
		LoggerError().Fatal(err)
	}

	return data
}
