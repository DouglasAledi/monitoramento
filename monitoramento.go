package main

import (
	"bufio"
	"encoding/json"

	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	sites := []string{}

	interfaceMenu(sites)
}

func verificarSite(site string, wg *sync.WaitGroup) {
	defer wg.Done()

	start := time.Now()
	resp, err := http.Get(site)
	duracao := time.Since(start)

	if err == nil {
		fmt.Println("✅", site, "está ONLINE -", resp.Status, "- Tempo:", duracao)
	} else {
		fmt.Println("❌", site, "está OFFLINE - Erro:", err, "- Tempo:", duracao)
	}
}

func interfaceMenu(sites []string) {
	var opcao int
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nMenu de opções")
		fmt.Println("1 - Exibir sites")
		fmt.Println("2 - Verificar o status dos sites")
		fmt.Println("3 - Adicionar site")
		fmt.Println("4 - Remover site")
		fmt.Println("5 - Sair")
		fmt.Print("Escolha uma opção: ")

		fmt.Scanln(&opcao)

		switch opcao {
		case 1:
			data, err := os.ReadFile("sites.json")
			if err != nil {
				if os.IsNotExist(err) {
					sites = []string{} // arquivo não existe, slice vazio
				} else {
					fmt.Println("Erro ao ler arquivo:", err)
					return
				}
			} else {
				// arquivo existe, faz o Unmarshal
				err = json.Unmarshal(data, &sites)
				if err != nil {
					fmt.Println("Erro ao converter JSON:", err)
					return
				}
			}

			// mostrar os sites
			for i, site := range sites {
				fmt.Printf("%d - %s\n", i+1, site)
			}

		case 2:
			var wg sync.WaitGroup
			wg.Add(len(sites))

			for _, site := range sites {
				go verificarSite(site, &wg)
			}
			wg.Wait()

		case 3:
			fmt.Print("Insira aqui o site novo: ")
			scanner.Scan()
			novoSite := scanner.Text()
			sites = append(sites, novoSite)

			data, err := json.MarshalIndent(sites, "", "  ")
			if err != nil {
				fmt.Println("Erro ao converter:", err)
				return
			}

			err = os.WriteFile("sites.json", data, 0644)
			if err != nil {
				fmt.Println("Erro ao salvar:", err)
				return
			}

			fmt.Println("✅ Site adicionado com sucesso!")

		case 4:
			fmt.Println("Digite aqui o site que deseja excluir: ")
			scanner.Scan()
			siteExcluido := scanner.Text()

			data, err := os.ReadFile("sites.json")
			if err != nil {
				if os.IsNotExist(err) {
					sites = []string{}
				} else {
					fmt.Println("Erro ao ler arquivo:", err)
					return
				}
			}
			if len(data) > 0 {
				err = json.Unmarshal(data, &sites)
				if err != nil {
					fmt.Println("Erro ao converter JSON:", err)
					return
				}
			}
			indice := -1
			for i, site := range sites {
				if site == siteExcluido {
					indice = i
					break
				}
			}
			if indice != -1 {
				sites = append(sites[:indice], sites[indice+1:]...)
				fmt.Println("✅ Site removido com sucesso!")
			} else {
				fmt.Println("❌ Site não encontrado na lista.")
			}

			data, err = json.MarshalIndent(sites, "", "  ")
			if err != nil {
				fmt.Println("Erro ao converter JSON:", err)
				return
			}

			err = os.WriteFile("sites.json", data, 0644)
			if err != nil {
				fmt.Println("Erro ao salvar arquivo:", err)
				return
			}

		case 5:
			fmt.Println("Encerrando programa...")
			os.Exit(0)

		default:
			fmt.Println("❌ Opção inválida.")
		}
	}
}
