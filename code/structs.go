package code

import (
	libs "GoLibs"
	"time"
)

type MalaDiretaST struct {
	Id           int       `db_autprimary:"true"`
	DataCadastro time.Time `db_notnull:"true"`
	Mensagem     string    `db_notnull:"true" db_comm:"Mensagem de envio"`
	Assunto      string    `db_notnull:"true" db_comm:"Assunto do malling"`

	AG_status int `db_default:"0" db_comm:"0 Aguardando, 1 Executando"`
	AG_dia    int `db_default:"0"`
	AG_mes    int `db_default:"0"`
	AG_hora   int `db_default:"0"`

	SMTP_servidor      string `db_comm:"Servidor SMTP"`
	SMTP_porta         int    `db_comm:"Porta do Servidor SMTP"`
	SMTP_email         string `db_comm:"Email do Servidor SMTP"`
	SMTP_senha         string `db_comm:"Senha do Servidor SMTP"`
	SMTP_retorno_nome  string `db_comm:"Nome do Servidor SMTP"`
	SMTP_retorno_email string `db_comm:"Email do Servidor SMTP"`
}

type MalaDiretaEmailsST struct {
	Id           int       `db_autprimary:"true"`
	DataCadastro time.Time `db_notnull:"true"`
	Email        string    `db_notnull:"true"`
	Descricao    string
}

type ListaEnvioST struct {
	Id              int       `db_autprimary:"true"`
	DataCadastro    time.Time `db_notnull:"true"`
	Email           string    `db_notnull:"true"`
	Descricao       string
	CodeStatus      int `db_default:"0" db_comm:"0 Aguardando, 1 Enviado Sucesso, 2 Erro ao Enviar"`
	EnvioData       time.Time
	EnvioErro       bool `db_default:"0"`
	EnvioTentativas int  `db_default:"0"`
}

type ListaRejeicaoST struct {
	Id           int       `db_autprimary:"true"`
	DataCadastro time.Time `db_notnull:"true"`
	Email        string    `db_notnull:"true"`
	Descricao    string
	Motivo       string
}

type ListaContatosST struct {
	Id            int       `db_autprimary:"true"`
	DataCadastro  time.Time `db_notnull:"true"`
	Email         string    `db_notnull:"true" db_tm1:"100"`
	Cliente       string
	Contato       string
	CategoriaDesc string `db_comm:"Categoria"`
}

func (s *ListaContatosST) InsertSQL() string {
	sSQL := "INSERT INTO ListaContatos (DataCadastro, Email, Cliente, Contato, CategoriaDesc)"
	sSQL += " VALUES ("
	sSQL += libs.Asp(libs.FormatDateTime("yyyy-mm-dd hh:nn:ss", s.DataCadastro))
	sSQL += "," + libs.Asp(s.Email)
	sSQL += "," + libs.Asp(s.Cliente)
	sSQL += "," + libs.Asp(s.Contato)
	sSQL += "," + libs.Asp(s.CategoriaDesc)
	sSQL += " ) "
	return sSQL
}
