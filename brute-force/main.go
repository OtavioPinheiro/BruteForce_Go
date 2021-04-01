package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/yeka/zip"
)

const (
	zipPath             = "brute-force\\assets\\lerolero_protected.zip"
	caminhoArquivoSenha = "brute-force\\assets\\rockyou.txt"
	numeroDeThreads     = 10
	linhasPorThreads    = 1500
)

func main() {
	abrirArquivoZip()
}

func bruteForce(zipPath string, listaDeSenhas []string, canal chan<- string) {
	// output := color.New(outputColor)
	arquivo, err := zip.OpenReader(zipPath)

	if err != nil {
		//log.Fatal(err)
		panic("Erro!")
	}
	defer arquivo.Close()

	zipFile := arquivo.File[0]

	for _, value := range listaDeSenhas {
		//	output.Printf("Trying to crack the file with password: %v \n", string(value))

		zipFile.SetPassword(string(value))
		_, err := zipFile.Open()

		if err == nil {
			fmt.Printf("Senha encontrada!\n")

			zipReader, err := zipFile.Open()

			if err != nil {
				// log.Fatal(err)
				panic("Erro")
			}

			buf, err := ioutil.ReadAll(zipReader)
			if err != nil {
				// log.Fatal(err)
				panic("Erro")
			}

			defer zipReader.Close()

			fmt.Printf("Tamanho do %v: %v byte(s)\n", zipFile.Name, len(buf))
			fmt.Println()

			canal <- string(value)
			break
		}
	}
}

func abrirArquivoZip() {
	// fmt.Println("Olá Go!")
	color.Blue("OLÁ GO!!")

	arquivos, err := zip.OpenReader(zipPath)

	if err != nil {
		// fmt.Printf("Um erro aconteceu: %v\n", err)
		panic(fmt.Sprintf("Um erro aconteceu: %v\n", err))
	}

	defer arquivos.Close()

	linhaInicial := 0
	zipFile := arquivos.File[0]

	if zipFile.IsEncrypted() {

		listaSenhas := obterListaDeSenhas(caminhoArquivoSenha)
		fmt.Println("O arquivo está protegido por senha.")

		canal := make(chan string, 1)

		start := time.Now()

		for i := 0; i < numeroDeThreads; i++ {
			linhaFinal := linhasPorThreads * (i + 1)

			// outputColor := RandomOutputColor(i)
			// output := color.New(outputColor)
			fmt.Printf("Começando pela thread %d lendo da linha %d até a linha %d\n", i+1, linhaInicial, linhaFinal)
			go bruteForce(zipPath, listaSenhas[linhaInicial:linhaFinal], canal)

			linhaInicial = linhaFinal + 1
		}

		fmt.Println("----------------------------------------------------")

		color.Yellow("Quebrando a senha...\n")
		fmt.Println()

		select {
		case senha := <-canal:
			// fmt.Printf("\nA senha é:\"%v\"\n", senha)
			// fmt.Printf("A quebra de senha demorou: %v\n", time.Since(start))
			color.Green("SENHA ENCONTRADA!!")

			color.Yellow("                      RELATÓRIO")
			color.Yellow("------------------------------------------------------------")
			color.Yellow("Número de Threads: %v", numeroDeThreads)
			color.Yellow("Linhas por Threads: %v", linhasPorThreads)
			color.Green("SENHA: %v", senha)
			color.Green("TEMPO GASTO: %v", time.Since(start))
		case <-time.After(time.Duration(15) * time.Second):
			//fmt.Printf("Timeout after: %d seconds \n", timeout)
			// fmt.Printf("Senha não encontrada :( \n")
			color.Red("SENHA NÃO ENCONTRADA =(")
		}

	} else {
		fmt.Printf("Sem proteção...\n")
	}
}

func obterListaDeSenhas(caminhoArquivoSenha string) []string {
	arquivo, erro := os.Open(caminhoArquivoSenha)

	if erro != nil {
		panic(fmt.Sprintf("Um erro aconteceu: %v\n", erro))
	}

	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)

	scanner.Split(bufio.ScanLines)
	var senhas []string

	for scanner.Scan() {
		senhas = append(senhas, scanner.Text())
	}

	arquivoStatus, erro := os.Stat(caminhoArquivoSenha)

	if erro != nil {
		panic(fmt.Sprintf("Um erro aconteceu: %v", erro))
	}

	fmt.Printf("Status do arquivo => Nome: %v, Tamanho: %v KB\n", arquivoStatus.Name(), arquivoStatus.Size()/(1024))
	fmt.Printf("Quantidade total de senhas no arquivo %d\n", len(senhas))

	return senhas
}
