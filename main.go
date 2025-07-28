package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/anil1te/lockbox/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("❌ Не указан режим. Используйте:")
		printUsage()
		return
	}

	const key = "secretkey123"

	switch os.Args[1] {
	case "add":

		if len(os.Args) < 5 {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Введите ссылку на сайт: ")
			site, _ := reader.ReadString('\n')
			site = strings.TrimSpace(site)

			fmt.Print("Введите логин: ")
			login, _ := reader.ReadString('\n')
			login = strings.TrimSpace(login)

			fmt.Print("Введите пароль: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)
			utils.AddEntry(site, login, password, key)
			return
		}

		site := os.Args[2]
		login := os.Args[3]
		password := os.Args[4]
		utils.AddEntry(site, login, password, key)

	case "list":
		utils.ListSitesWithCounts()

	case "del":
		if len(os.Args) < 4 {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Введите домен сайта: ")
			site, _ := reader.ReadString('\n')
			site = strings.TrimSpace(site)

			fmt.Print("Введите логин: ")
			login, _ := reader.ReadString('\n')
			login = strings.TrimSpace(login)
			utils.RemoveEntry(site, login, key)
			return
		}
		domain := os.Args[2]
		login := os.Args[3]
		utils.RemoveEntry(domain, login, key)

	case "help":
		printUsage()

	default:
		utils.GetCredentials(os.Args[1], key)
	}
}

func printUsage() {
	fmt.Println(`
📖 Использование:
  add            - Добавить новую запись
  сайт           - Получить данные по сайту
  del            - Удалить запись
  list           - Показать список сайтов и количество записей
💬 Пример:
  lockbox add или
    lockbox add https://google.com login password

  lockbox google

  lockbox list

  lockbox del или
    lockbox del google login
`)
}
