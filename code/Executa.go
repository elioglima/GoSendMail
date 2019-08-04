package code

import (
	"GoLibs"
	"GoLibs/logs"
	"GoMysql"
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"

)

var (
	smtpClient *smtp.Client
)

/*

	Instalador do banco de dados

*/

func ConectarServidorSMTP(p *GoLibs.SendSMTPMailDadosST) error {

	var err error

	// Connect to the SMTP Server
	// servername := "smtp.perfectvision.kinghost.net:587"
	// "atendimento@perfectvision.kinghost.net"
	logs.Atencao("Conectando SMTP", p.SMTP_Server)
	host, _, _ := net.SplitHostPort(p.SMTP_Server)
	auth := smtp.PlainAuth("", p.SMTP_Mail, p.SMTP_Senha, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	smtpClient, err = smtp.Dial(p.SMTP_Server)
	if err != nil {
		logs.Erro("e-mail")
		return err
	}

	smtpClient.StartTLS(tlsconfig)

	// Auth
	if err = smtpClient.Auth(auth); err != nil {
		logs.Erro("e-mail")
		return err
	}

	return nil
}

func Executa() {
	// 8425
	logs.DebugOrigem = false
	logs.Atencao("Iniciando Processo de execução.")

	Params := GoMysql.ParamsConexaoST{}
	Params.IP = "localhost"
	Params.PORTA = 3306
	Params.BANCO = "xpressapi"
	Params.USUARIO = "root"
	Params.SENHA = "AB@102030"

	logs.Atencao("Iniciando biblioteca de conexão")
	Conexao := GoMysql.NewConexao(Params)

	logs.Atencao("Efetuando conexão")
	if err := Conexao.Conectar(); err != nil {
		logs.Erro(err)
		return
	}

	logs.Atencao("Verificando banco de dados")

	sSQL := "select * from maladireta limit 0,1"
	RecordCount, Results, err := Conexao.Query(sSQL)
	if err != nil {
		logs.Erro(err)
		return
	}

	if RecordCount == 0 {
		logs.Erro("Nenhuma tarefa localizada")
		return
	}
	logs.Atencao("Selecionando tarefas")

	for _, Result := range Results {

		p := &GoLibs.SendSMTPMailDadosST{}

		smtp_porta := GoMysql.ValueStr(Result, "smtp_porta")
		p.SMTP_Server = GoMysql.ValueStr(Result, "smtp_servidor") + ":" + smtp_porta
		p.SMTP_Mail = GoMysql.ValueStr(Result, "smtp_email")
		p.SMTP_Senha = GoMysql.ValueStr(Result, "smtp_senha")
		smtpRetornoNome := GoMysql.ValueStr(Result, "smtp_retorno_nome")
		smtpRetornoEmail := GoMysql.ValueStr(Result, "smtp_retorno_email")

		p.From = mail.Address{smtpRetornoNome, smtpRetornoEmail}
		p.Subj = GoMysql.ValueStr(Result, "assunto")
		p.Body = GoMysql.ValueStr(Result, "mensagem")

		logs.Atencao("Listando emails")

		sSQL := "select * from listaenvio where codestatus = 0 limit 0,50"
		RecordCount, Rlistacontatos, err := Conexao.Query(sSQL)
		if err != nil {
			logs.Erro(err)
			return
		}

		if RecordCount == 0 {
			logs.Erro("Nenhuma tarefa localizada")
			return
		}

		// Setup headers
		headers := make(map[string]string)
		headers["From"] = p.From.String()
		headers["To"] = p.To.String()
		headers["Subject"] = p.Subj

		// Setup message
		msg := ""
		for k, v := range headers {
			msg += fmt.Sprintf("%s: %s\r\n", k, v)
		}

		msg += "Content-Type: text/html; charset=\"utf-8\"\r\n"
		msg += "Content-Transfer-Encoding: 7bit\r\n"
		msg += fmt.Sprintf("\r\n%s", p.Body+"\r\n")

		if err := ConectarServidorSMTP(p); err != nil {
			logs.Erro("ConectarServidorSMTP(", err)
			return
		}

		// defer c.Quit()

		logs.Atencao("Preparando para envio")
		for _, contato := range Rlistacontatos {

			logs.Cyan("Iniciando envio")

			// logs.Atencao("c.Mail(p.From.Address)", p.From.Address)
			if err = smtpClient.Mail(p.From.Address); err != nil {
				logs.Erro("c.Mail(", err)
				if err := ConectarServidorSMTP(p); err != nil {
					logs.Erro("c.Mail(", err)
					time.Sleep(3 * time.Second)
					os.Exit(0)
				}

				time.Sleep(3 * time.Second)
				if err = smtpClient.Mail(p.From.Address); err != nil {
					logs.Erro("c.Mail(", err)
					os.Exit(0)
					continue
				}
			}

			email := strings.TrimSpace(strings.ToLower(GoMysql.ValueStr(contato, "email")))
			// email = "diretoria@maxtime.info"
			p.To = mail.Address{email, email}

			if strings.Contains(p.To.Address, " ") {
				sSQL := " delete from listaenvio "
				sSQL += " where id = " + strconv.Itoa(GoMysql.ValueInt(contato, "id"))
				if _, err := Conexao.Execute(sSQL); err != nil {
					logs.Erro("Conexao.Execute(", err)
				}
				continue
			}

			logs.Atencao("c.Rcpt(p.To.Address)", p.To.Address)
			if err = smtpClient.Rcpt(p.To.Address); err != nil {
				logs.Erro("Rcpt", err)

				if strings.Contains(strings.ToLower(err.Error()), "ultrapassou o limite") {
					logs.Atencao("ultrapassou o limite")
					time.Sleep(3 * time.Second)
					os.Exit(0)
				} else if strings.Contains(strings.ToLower(err.Error()), "bad recipient a") {
					sSQL := " delete from listaenvio "
					sSQL += " where id = " + strconv.Itoa(GoMysql.ValueInt(contato, "id"))
					if _, err := Conexao.Execute(sSQL); err != nil {
						logs.Erro(err)
					}

				} else if strings.Contains(strings.ToLower(err.Error()), "domain not found") {
					sSQL := " delete from listaenvio "
					sSQL += " where id = " + strconv.Itoa(GoMysql.ValueInt(contato, "id"))
					if _, err := Conexao.Execute(sSQL); err != nil {
						logs.Erro("delete from listaenvio", err)
					}
					logs.Atencao("Registro deletado")
				}
				continue
			}

			logs.Atencao("Deletando DB")
			sSQL := " delete from listaenvio "
			sSQL += " where id = " + strconv.Itoa(GoMysql.ValueInt(contato, "id"))
			if _, err := Conexao.Execute(sSQL); err != nil {
				logs.Erro(err)
				continue

			}

			logs.Atencao("Data")
			w, err := smtpClient.Data()
			if err != nil {
				logs.Erro(err)
				continue

			}

			logs.Atencao("w.Write([]byte(msg))")
			_, err = w.Write([]byte(msg))
			if err != nil {
				logs.Erro(err)

			}

			w.Close()

			// logs.Atencao("c.Text.ReadResponse")
			// code, message, err := c.Text.ReadResponse(0)
			// logs.Atencao("code, message, err", code, message, err)
			// if strings.Contains(err.Error(), "EOF") {
			// 	logs.Rosa("code", code)
			// 	logs.Rosa("message", message)
			// 	logs.Rosa("err", err)
			// 	continue
			// }

			// if err != nil {
			// 	logs.Erro("e-mail", email, err)
			// 	sSQL = "update listaenvio set "
			// 	sSQL += " enviodata = current_date()"
			// 	sSQL += " ,codestatus = 2"
			// 	sSQL += " ,envioerro = 1"
			// 	sSQL += " ,enviotentativas = enviotentativas+1"
			// 	sSQL += " where id = " + strconv.Itoa(GoMysql.ValueInt(contato, "id"))
			// } else {
			// 	logs.Sucesso("e-mail", email, "enviado")

			// }

		}

	}

	logs.Sucesso("Processo de envio finalizado com sucesso.")
	os.Exit(0)

}
