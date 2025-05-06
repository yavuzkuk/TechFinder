package main

import (
	"Wapplyzer/Functions"
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {

	err := os.WriteFile("output.txt", []byte(""), 0644) // içeriği temizler
	if err != nil {
		log.Fatalf("Dosya temizlenemedi: %v", err)
	}

	var filePath string
	flag.StringVar(&filePath, "f", "", "file path")
	flag.Parse()

	if filePath == "" {
		log.Fatal("file path is required")
	}

	file, err := os.Open(filePath)
	Functions.ErrorHandler(err)

	defer file.Close()

	scanner := bufio.NewScanner(file)
	client := Functions.Client()

	for scanner.Scan() {
		line := scanner.Text()
		var domain string
		if line == "" {
			continue
		} else if strings.HasPrefix(line, "https://") || strings.HasPrefix(line, "http://") {
			domain = line
		} else if !strings.HasPrefix(line, "https://") || !strings.HasPrefix(line, "http://") {
			domain = "https://" + line
		}
		request := Functions.Request(domain)
		fmt.Println("**********" + domain + "***********")

		response, err := client.Do(request)
		if err != nil {
			break
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			headerInfo := Functions.HeaderChecker(*response)
			componentInfo := Functions.BodyParser(*response)

			err := Functions.CreateAndAppendJSONToFile(domain, headerInfo)
			if err != nil {
				log.Fatal(err)
			}

			err = Functions.CreateAndAppendJSONToFile("", componentInfo)
			if err != nil {
				log.Fatal(err)
			}
			Functions.AppendCDN(Functions.DetectCDNs(domain))
			Functions.DetectCMS(domain)
			
			Functions.VulnCheck(domain, componentInfo)
			fmt.Println(response.StatusCode)
		} else if response.StatusCode == http.StatusConflict {
			fmt.Println(response.StatusCode)
		} else {
			fmt.Println(response.StatusCode)
		}
		// 302 yada 301 gibi değer döndürenleri Location değerine göre tekrardan o sayfaya yönlendirebiliriz.
	}
}
