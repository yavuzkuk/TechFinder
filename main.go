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

// https://www.cybersecurity-help.cz/vdb/twbs/bootstrap/3.3.2/
// https://www.cybersecurity-help.cz/vdb/jquery/  bu adres üzerinden zafiyet var mı diye kontrol edilir.
// farklı farklı paketleri var.
// buradan araştırılmalı

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
		statusCode, location := Functions.CheckRedirect(line)

		for {
			if statusCode >= 300 && statusCode < 400 && strings.TrimSpace(location) != "" {
				line = location
				statusCode, location = Functions.CheckRedirect(strings.TrimSpace(location))
			} else {
				break
			}
		}

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
			fmt.Println("!!!!!!!!!!!!!!! 10 saniye cevap alınamadı !!!!!!!!!!!!!!!\n")
			//fmt.Println("İstek hatası:", err)
			continue
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

			vulnLink := Functions.VulnCheck(domain, componentInfo)
			if len(vulnLink) != 0 {
				Functions.AppendVulnLink(vulnLink)
			}
		} else if response.StatusCode == http.StatusConflict {
			fmt.Println(response.StatusCode)
		}
	}
}
