package code

import (
	"GoLibs/logs"
	"GoMysql/GoMysql"
	"bufio"
	"os"
	"strings"
	"time"
)

/*

	Instalador do banco de dados

*/

func Instalar() {

	logs.DebugSucesso = true
	logs.DebugErro = true
	logs.DebugOrigem = true

	logs.Atencao("Processo de instalação iniciando.")

	Params := GoMysql.ParamsConexaoST{}
	Params.IP = "localhost"
	Params.PORTA = 3306
	Params.BANCO = "xpressapi"
	Params.USUARIO = "root"
	Params.SENHA = "AB@102030"

	logs.Atencao("Iniciando biblioteca de conexão")
	Conexao := GoMysql.NewConexao(Params)

	logs.Atencao("Efetuando conexão")
	if err := Conexao.ConectarSystem(); err != nil {
		logs.Erro(err)
		return
	}

	logs.Atencao("Verificando banco de dados")
	if err := Conexao.CreateDB(); err != nil {
		logs.Erro(err)
		return
	}

	logs.Atencao("Verificando tabela de ListaContatos")
	ListaContatos := &ListaContatosST{}
	if err := Conexao.CreateTable(ListaContatos); err != nil {
		logs.Erro(err)
		return
	}

	logs.Atencao("Verificando tabela de MalaDireta")
	MalaDireta := &MalaDiretaST{}
	if err := Conexao.CreateTable(MalaDireta); err != nil {
		logs.Erro(err)
		return
	}

	logs.Atencao("Verificando tabela de MalaDiretaEmails")
	MalaDiretaEmails := &MalaDiretaEmailsST{}
	if err := Conexao.CreateTable(MalaDiretaEmails); err != nil {
		logs.Erro(err)
		return
	}

	logs.Atencao("Verificando tabela de Lista de envios")
	ListaEnvio := &ListaEnvioST{}
	if err := Conexao.CreateTable(ListaEnvio); err != nil {
		logs.Erro(err)
		return
	}

	logs.Atencao("Verificando tabela de Lista de rejeição")
	ListaRejeicao := &ListaRejeicaoST{}
	if err := Conexao.CreateTable(ListaRejeicao); err != nil {
		logs.Erro(err)
		return
	}

	logs.Atencao("Inserindo dados")

	file, err := os.Open("Contatos.txt")
	if err != nil {
		logs.Erro(err)
		return
	}
	defer file.Close()

	// _, err = Conexao.Execute("delete from ListaContatos")
	// if err != nil {
	// 	logs.Erro(err)
	// 	return
	// }

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linha := scanner.Text()
		campos := strings.Split(linha, ";")
		ListaContatos := &ListaContatosST{}

		DataCadastro := strings.Replace(campos[0], " ", "T", -1) + ".000Z"
		layout := "2006-01-02T15:04:05.000Z"
		DataCadastroNew, err := time.Parse(layout, DataCadastro)

		if err != nil {
			ListaContatos.DataCadastro = time.Now()
		} else {
			ListaContatos.DataCadastro = DataCadastroNew
		}

		ListaContatos.Email = campos[1]
		ListaContatos.Cliente = strings.Replace(strings.Replace(campos[2], "\"", "", -1), "'", "", -1)
		ListaContatos.Contato = strings.Replace(strings.Replace(campos[3], "\"", "", -1), "'", "", -1)
		ListaContatos.CategoriaDesc = strings.Replace(strings.Replace(campos[4], "\"", "", -1), "'", "", -1)

		sSQL := ListaContatos.InsertSQL()
		_, err = Conexao.Execute(sSQL)
		if err != nil {
			logs.Erro(err, sSQL)
			return
		}

	}

	if err := scanner.Err(); err != nil {
		logs.Erro(err)
		return
	}

	logs.Sucesso("Processo de instalação com sucesso.")

}
