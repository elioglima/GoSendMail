package code

import (
	"GoLibs"
	"GoLibs/logs"
	"GoMysql/GoMysql"
	"net/mail"
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

		logs.Atencao("Enviando email.")
		p.To = mail.Address{"Diretoria", "diretoria@maxtime.info"}
		GoLibs.SendSMTPMail(p)

	}

	logs.Sucesso("Processo de instalação com sucesso.")

}
