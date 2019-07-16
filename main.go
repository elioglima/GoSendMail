package main

import (
	"GoLibs/logs"
	"GoSendMail/code"
	"time"
)

func main() {

	logs.DebugSucesso = true
	logs.DebugErro = true
	logs.DebugOrigem = true

	// for {
	logs.Sucesso("Iniciando processo de envio..")
	// code.Instalar()
	code.Executa()
	logs.Sucesso("Finalizando processo de envio..")
	time.Sleep(30 * time.Second)
	// }

}
