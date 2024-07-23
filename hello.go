package main

import (
	"bufio"
	"fmt" //Permite o acesso WEB com GO
	"io"
	"io/ioutil"
	"net/http" // permite acesso a sites de http
	"os"       // Apresentar status das requisições para o sistema
	"strconv"
	"strings"
	"time" //adicionar time sleep
)

const monitoramentos = 3
const delay = 5 * time.Second

func main() {
	exibeIntroducao()
	registraLog("site-falso", false)
	for {
		exibeMenu()
		comando := lerComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Abrindo LOGS...")
			imprimeLog()
		case 0:
			fmt.Println("Encerrando APP...")
			os.Exit(0)
		default:
			fmt.Println("Comando inválido")
			os.Exit(-1)
		}
	}

}

// MENSAGEM DE BEM VINDO.
func exibeIntroducao() {
	nome := "Douglas"
	versao := 1.1
	fmt.Println("Olá Sr(a).", nome)
	fmt.Println("O App está na versão", versao)
}

// RECEBE O COMANDO
func lerComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido)

	return comandoLido
}

// ABRE O MENU
func exibeMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("0 - Sair do APP")
}

// STARTA O MONITORAMENTO E PERCORRE TODA A SLICE TESTANDO PELO TEMPO DETERMINADO NA CONST:MONITORAMENTOS E COM O DELAY DETERMINADO NA CONST:DELAY
func iniciarMonitoramento() {
	fmt.Println("Iniciando Monitoramento...")
	sites := lerSitesArquivo()
	for i := 0; i < monitoramentos; i++ {
		for i, site := range sites {
			fmt.Println("Testando site", i, ":", site)
			testaSite(site)
		}
		time.Sleep(delay)
	}
}

// TESTA O STATUS DOS SITES DENTRO DO SLICE
func testaSite(site string) {
	resp, err := http.Get(site)
	if site == "" {
		fmt.Println("URL vazia, ignorando...")
		return
	}
	if err != nil {
		fmt.Println("Ocorreu um erro", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println(site, "Ta potente!")
		registraLog(site, true)
	} else {
		fmt.Println(site, "ixi, deu merda!")
		registraLog(site, false)
	}
}

// FAZER COM QUE O NOSSO MONITORAMENTO LEIA DADOS DE UM ARQUIVO EXTERNO.
func lerSitesArquivo() []string {
	var sites []string
	arquivo, err := os.Open("sites.txt")
	if err != nil {
		fmt.Println("Ocorreu um Erro!", err)
	}

	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		fmt.Println(linha)
		sites = append(sites, linha)

		if err == io.EOF {
			break
		}
	}

	//iremos começar a tratar as mensagens de erro, antes: arquivo, _ := os.Open("sites.txt")
	arquivo.Close()
	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Erro encontrado:", err)
	}

	arquivo.WriteString(time.Now().Format("02/01/2006 15:04:05") + " - " + site + " - Online - " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func imprimeLog() {
	arquivo, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(arquivo))

}
