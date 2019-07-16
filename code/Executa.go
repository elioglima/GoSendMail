package code

import (
	"GoLibs"
	"GoLibs/logs"
	"GoMysql/GoMysql"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"
)

/*

	Instalador do banco de dados

*/

func Executa() {

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

	for _, Result := range Results {

		p := GoLibs.SendSMTPMailST{}

		smtp_porta := GoMysql.ValueStr(Result, "smtp_porta")
		p.SMTP_Server = GoMysql.ValueStr(Result, "smtp_servidor") + ":" + smtp_porta
		p.SMTP_Mail = GoMysql.ValueStr(Result, "smtp_email")
		p.SMTP_Senha = GoMysql.ValueStr(Result, "smtp_senha")
		smtpRetornoNome := GoMysql.ValueStr(Result, "smtp_retorno_nome")
		smtpRetornoEmail := GoMysql.ValueStr(Result, "smtp_retorno_email")

		p.From = mail.Address{smtpRetornoNome, smtpRetornoEmail}
		p.Subj = GoMysql.ValueStr(Result, "assunto")
		p.Body = GoMysql.ValueStr(Result, "mensagem")

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

		for _, contato := range Rlistacontatos {
			email := GoMysql.ValueStr(contato, "email")
			p.To = mail.Address{email, email}

			code, message, err := GoLibs.SendSMTPMail(p)

			if strings.Contains(err.Error(), "Voce ultrapassou o limite") {
				logs.Rosa("code", code)
				logs.Rosa("message", message)
				logs.Rosa("err", err)
				time.Sleep(20 * time.Second)
				os.Exit(0)
			}

			if err != nil {
				logs.Erro("e-mail", email, err)
				sSQL = "update listaenvio set "
				sSQL += " enviodata = current_date()"
				sSQL += " ,codestatus = 2"
				sSQL += " ,envioerro = 1"
				sSQL += " ,enviotentativas = enviotentativas+1"
				sSQL += " where id = " + strconv.Itoa(GoMysql.ValueInt(contato, "id"))
			} else {
				logs.Atencao("e-mail", email, "enviado")
				sSQL = " delete from listaenvio "
				sSQL += " where id = " + strconv.Itoa(GoMysql.ValueInt(contato, "id"))
			}

			if _, err := Conexao.Execute(sSQL); err != nil {
				logs.Erro(err)
			}

		}

	}

	logs.Sucesso("Processo de envio finalizado com sucesso.")

}
