package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type Entry struct {
	Site     string `json:"site"`
	URL      string `json:"url"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func extractDomain(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host := parsed.Hostname()
	parts := strings.Split(host, ".")
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	return host
}

func hashKey(key string) []byte {
	sum := sha256.Sum256([]byte(key))
	return sum[:]
}

func encrypt(text string, key string) string {
	block, _ := aes.NewCipher(hashKey(key))
	plaintext := []byte(text)

	// создаём IV (ноль-байты фиксированной длины)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	mode := cipher.NewCFBEncrypter(block, iv)
	mode.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func decrypt(cryptoText string, key string) string {
	ciphertext, _ := base64.StdEncoding.DecodeString(cryptoText)
	block, _ := aes.NewCipher(hashKey(key))
	iv := ciphertext[:aes.BlockSize]
	text := ciphertext[aes.BlockSize:]

	mode := cipher.NewCFBDecrypter(block, iv)
	mode.XORKeyStream(text, text)

	return string(text)
}

func AddEntry(rawURL, login, password, key string) {
	domain := extractDomain(rawURL)
	encLogin := encrypt(login, key)
	encPass := encrypt(password, key)

	entry := Entry{
		Site:     domain,
		URL:      rawURL,
		Login:    encLogin,
		Password: encPass,
	}

	entries := ReadAllEntries()
	entries = append(entries, entry)

	file, err := os.Create("data.json")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(entries)
	fmt.Println("✅ Запись добавлена! Что бы получить данные по сайту, используйте команду: lockbox", extractDomain(domain))

}

func ReadAllEntries() []Entry {
	file, err := os.Open("data.json")
	if err != nil {
		return []Entry{}
	}
	defer file.Close()

	var entries []Entry
	json.NewDecoder(file).Decode(&entries)
	return entries
}

func GetCredentials(domain string, key string) {
	entries := ReadAllEntries()
	found := false

	for _, e := range entries {
		if e.Site == domain {
			login := decrypt(e.Login, key)
			pass := decrypt(e.Password, key)
			fmt.Printf("🔐 Найдено:\n  Login: %s\n  Password: %s\n", login, pass)
			found = true
		}
	}
	if !found {
		fmt.Println("❌ Ничего не найдено для:", domain)
	}
}

func RemoveEntry(domain, login, key string) {
	entries := ReadAllEntries()
	newEntries := make([]Entry, 0)
	found := false

	for _, e := range entries {
		if e.Site == domain && decrypt(e.Login, key) == login {
			found = true
			continue
		}
		newEntries = append(newEntries, e)
	}

	if !found {
		fmt.Println("❌ Запись не найдена")
		return
	}

	file, err := os.Create("data.json")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(newEntries)
	if err != nil {
		fmt.Println("Ошибка при сохранении:", err)
		return
	}

	fmt.Println("✅ Запись успешно удалена")
}

func ListSitesWithCounts() {
	entries := ReadAllEntries()
	if len(entries) == 0 {
		fmt.Println("📭 Нет сохранённых записей.")
		return
	}

	counts := make(map[string]int)
	for _, e := range entries {
		counts[e.Site]++
	}

	fmt.Println("🔐 Сайты и количество записей:")
	for site, count := range counts {
		fmt.Printf("  • %s: %d запись(и)\n", site, count)
	}
}
