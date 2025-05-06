package Functions

import (
	"Wapplyzer/Constant"
	"Wapplyzer/Struct"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func ErrorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type CommentInfo struct {
	Content string
	Type    string
}

func Client() http.Client {
	return http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

func Request(url string) *http.Request {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	ErrorHandler(err)

	request.Header.Set("User-Agent", Constant.UserAgent[rand.Intn(len(Constant.UserAgent)-1)])
	request.Header.Set("Accept", "*/*")
	return request
}

func JSLinkEXtract(body []byte, response http.Response) map[string]string {

	jsLinks := make(map[string]string)

	re := regexp.MustCompile(`<script[^>]*\s*(?:type=["']text/javascript["']\s*)?[^>]*\s*(?:src|href)\s*=\s*["']([^"']+\.js[^"']*)["'][^>]*>`)

	links := re.FindAllStringSubmatch(string(body), -1)
	for _, link := range links {
		link := link[1]
		if link != "" {
			if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
				// harici bir kaynak
				jsLinks[link] = ""
			} else {
				// dahili bir kaynak
				if strings.HasPrefix(link, "//") {
					link = response.Request.URL.Scheme + "://" + strings.Replace(link, "//", "", 1)
				} else if strings.HasPrefix(link, "/") {
					link = response.Request.URL.String() + link
				} else {
					link = response.Request.URL.String() + "/" + link
				}
				jsLinks[link] = ""
			}
		} else {
			continue
		}
	}
	return jsLinks
}

func CSSLinkExtract(body []byte, response http.Response) map[string]string {

	re := regexp.MustCompile(`<link[^>]+href=["']([^"']+\.css[^"']*)["']`)

	cssLinks := make(map[string]string)

	matches := re.FindAllStringSubmatch(string(body), -1)

	for _, match := range matches {
		// href değeri linklerin kısmı
		var newUrl string
		href := strings.TrimSpace(match[1])
		if href != "" {
			if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
				cssLinks[href] = ""
			} else if strings.HasPrefix(href, "//") {
				newHref, _ := strings.CutPrefix(href, "//")
				if strings.HasSuffix(response.Request.URL.String(), "/") {
					newUrl = response.Request.URL.String() + newHref
					cssLinks[newUrl] = ""
				} else {
					newUrl = response.Request.URL.String() + "/" + href
					cssLinks[newUrl] = ""
				}
			} else if strings.HasPrefix(href, "/") {
				// "/assets/asdasdsadad.css
				if strings.HasSuffix(response.Request.URL.String(), "/") {
					newUrl, _ = strings.CutPrefix(href, "/")
					newUrl = response.Request.URL.String() + newUrl
					cssLinks[newUrl] = ""
				} else {
					href := response.Request.URL.String() + href
					cssLinks[href] = ""
				}
			} else if href != "" && !strings.HasPrefix(href, "//") || strings.HasPrefix(href, "/") {
				cssLinks[response.Request.URL.String()+"/"+href] = ""
			}
		}
	}

	return cssLinks
}

func FileRequest(urls map[string]string) map[string]string {

	client := Client()

	versionAndTypeList := make(map[string]string)

	for url, _ := range urls {
		request := Request(url)
		response, err := client.Do(request)
		if err != nil {
			break
		}
		fmt.Println(url)

		responseBody, _ := io.ReadAll(response.Body)

		comments := ExtractComments(responseBody)

		for comment, _ := range comments {
			//fmt.Println("Comment ---> ", strings.TrimSpace(comment))
			if strings.Contains(comment, "/*") || strings.Contains(comment, "/*!") || strings.Contains(comment, "//") {
				//versionAndType := ExtractVersionAndType(comment)
				exportComponentType := ExtractVersionAndType2(comment)

				for _, b := range exportComponentType {
					versionAndTypeList[b] = ""
				}

			}
		}
	}
	return versionAndTypeList
}

func ExtractComments(sourceCode []byte) map[string]string {

	commentList := make(map[string]string)
	re1 := regexp.MustCompile(`\/\/.*`)
	match1 := re1.FindAllString(string(sourceCode), -1)

	re2 := regexp.MustCompile(`\/\*![\s\S]*?\*\/`)
	match2 := re2.FindAllString(string(sourceCode), -1)

	re3 := regexp.MustCompile(`\/\*\*[\s\S]*?\*\/`)
	match3 := re3.FindAllString(string(sourceCode), -1)

	re4 := regexp.MustCompile(`\/\*[\s\S]*?\*\/`) // cogu seyi bulur.
	match4 := re4.FindAllString(string(sourceCode), -1)

	if len(match1) != 0 {
		for _, match := range match1 {
			if strings.TrimSpace(match) != "" {
				commentList[strings.TrimSpace(match)] = ""
			}
		}
	}

	if len(match2) != 0 {
		for _, match := range match2 {
			if strings.TrimSpace(match) != "" {
				commentList[strings.TrimSpace(match)] = ""
			}
		}
	}

	if len(match3) != 0 {
		for _, match := range match3 {
			if strings.TrimSpace(match) != "" {
				commentList[strings.TrimSpace(match)] = ""
			}
		}
	}

	if len(match4) != 0 {
		for _, match := range match4 {
			if strings.TrimSpace(match) != "" {
				commentList[strings.TrimSpace(match)] = ""
			}
		}
	}

	return commentList
}

func ExtractVersionAndType2(comment string) []string {

	var componentAndVersion []string
	re1 := regexp.MustCompile(`\b([A-Za-z0-9 .+\-]+?)\s+v(\d+(?:\.\d+){1,}(?:[a-z0-9\-]*)?)\b`)
	//______ v1.2.3
	match1 := re1.FindAllStringSubmatch(comment, -1)

	for _, match := range match1 {
		if strings.Contains(strings.ToLower(match[1]), "requires") {
			re2 := regexp.MustCompile(`(?m)^\s*\*\s*([^\r\n]+)[\s\S]+?\bversion:\s*([^\r\n]+)`)
			match2 := re2.FindAllStringSubmatch(comment, -1)

			component := strings.TrimSpace(match2[0][1])
			version := strings.TrimSpace(match2[0][2])
			componentWithVersion := component + " " + version
			componentAndVersion = append(componentAndVersion, strings.TrimSpace(componentWithVersion))
		} else {
			componentWithVersion := strings.TrimSpace(match[1]) + " " + strings.TrimSpace(match[2])
			componentAndVersion = append(componentAndVersion, strings.TrimSpace(componentWithVersion))
		}
	}

	re2 := regexp.MustCompile(`(?:\*)?([a-zA-Z\s\.]+)(\d\.[\d\.]+)`)
	// _____ 1.2
	// elde edilen değerler kontrol edilmeli
	match2 := re2.FindAllStringSubmatch(comment, -1)

	for _, match := range match2 {
		component := strings.TrimSpace(match[1])
		version := strings.TrimSpace(match[2])

		if strings.TrimSpace(component) == "" {
			break
		}

		if strings.Contains(strings.ToLower(component), "requires") {
			re2 := regexp.MustCompile(`(?m)^\s*\*\s*([^\r\n]+)[\s\S]+?\bversion:\s*([^\r\n]+)`)
			match2 := re2.FindAllStringSubmatch(comment, -1)

			component := strings.TrimSpace(match2[0][1])
			version := strings.TrimSpace(match2[0][2])
			componentWithVersion := component + " " + version
			componentAndVersion = append(componentAndVersion, componentWithVersion)
		} else {
			if strings.HasSuffix(component, "v") {
				before, _ := strings.CutSuffix(component, "v")

				if strings.TrimSpace(before) == "" {
					break
				} else {
					componentWithVersion := strings.TrimSpace(before) + " " + strings.TrimSpace(version)
					componentAndVersion = append(componentAndVersion, componentWithVersion)

				}
			} else if strings.HasSuffix(strings.TrimSpace(component), "version") {
				before, _ := strings.CutSuffix(strings.TrimSpace(component), "version")

				if strings.TrimSpace(before) == "" {
					break
				}
			} else if strings.TrimSpace(component) != "" {
				if len(strings.TrimSpace(component)) < 3 {
					break
				} else {
					var counter int = 0
					// for döngüsüyle teknolojiler kontrol edilmeli
					for _, tech := range Constant.TechList {
						if strings.ToLower(tech) == strings.ToLower(component) {
							counter++
							componentWithVersion := strings.TrimSpace(component) + " " + strings.TrimSpace(version)
							componentAndVersion = append(componentAndVersion, componentWithVersion)
						}
					}

					if counter == 0 {
						break
					}

				}
			}
		}
	}

	return componentAndVersion
}

//func ExtractVersionAndType(comment string) string {
//	// eğer bir hata çıkarsa - . işaretini regexten sil
//	re1 := regexp.MustCompile(`[a-zA-Z0-9\s-\.]+v(?:ersion)?([\d\.])+`)
//
//	match1 := re1.FindString(comment)
//
//	if len(match1) > 0 {
//		versionNumber := regexp.MustCompile(`v(?:ersion)?([\d.]+)`).FindStringSubmatch(match1)
//
//		componentType := regexp.MustCompile(`v(?:ersion)?([\d.]+)`).Split(match1, -1)
//
//		var compoType string = componentType[0]
//
//		if strings.HasSuffix(strings.TrimSpace(compoType), "-") {
//			compoType = strings.ReplaceAll(strings.TrimSpace(compoType), "-", "")
//		}
//		if strings.Contains(compoType, "Requires") {
//			compoType = strings.ReplaceAll(strings.TrimSpace(compoType), "Requires", "")
//		}
//
//		componentTypeWithVersion := strings.TrimSpace(compoType) + " " + strings.TrimSpace(versionNumber[1])
//		return componentTypeWithVersion
//	} else {
//		re2 := regexp.MustCompile(`(?:V|v)?ersion[\s\:\-]+([\d.])+`)
//		re3 := regexp.MustCompile(`(N|n)ame[:\-\s]+([a-zA-Z.])+`)
//
//		versionMatch := re2.FindString(comment)
//		componentMatch := re3.FindString(comment)
//
//		if versionMatch != "" && componentMatch != "" {
//			versionReg2 := regexp.MustCompile(`[(\d\.)]+`)
//			versionNumeric := versionReg2.FindStringSubmatch(versionMatch)
//
//			componentReg2 := regexp.MustCompile(`[a-zA-Z]+\.(?:css|js)`)
//			componentMatch2 := componentReg2.FindStringSubmatch(componentMatch)
//
//			componentWithVersion := componentMatch2[0] + " " + versionNumeric[0]
//
//			return componentWithVersion
//		}
//		//else {
//		//	re4 := regexp.MustCompile(`[(\n\sa-zA-Z0-9\.)]+\s+v?(?:ersion)?([\d\.])+`)
//		//
//		//	match4 := re4.FindStringSubmatch(comment)
//		//
//		//	fmt.Printcln("MATCH 4 ----------------")
//		//	fmt.Println(match4)
//		//}
//	}
//	return ""
//}

func DetectCMS(url string) {
	client := Client()

	request := Request(url)
	resp, err := client.Do(request)

	if err != nil {
		log.Fatalf("İstek gönderilirken hata: %v", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("HTML parse edilirken hata: %v", err)
	}

	cms := ""
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		if name == "generator" {
			content, _ := s.Attr("content")
			cms = content
		}
	})

	if cms != "" {
		fmt.Printf("Tespit edilen CMS: %s\n", cms)
	} else {
		fmt.Println("CMS bilgisi meta tag'de bulunamadı.")
	}
}

func CreateAndAppendJSONToFile(domain string, content []Struct.ComponentDetail) error {
	file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("dosya açma/oluşturma hatası: %v", err)
	}
	defer file.Close()

	if domain != "" {
		file.Write([]byte("******************* " + domain + " *******************\n"))
	}

	if len(content) > 0 {
		jsonData, err := json.MarshalIndent(content, "", "  ")
		if err != nil {
			return fmt.Errorf("json'a çevirme hatası: %v", err)
		}
		if _, err := file.Write(jsonData); err != nil {
			return fmt.Errorf("yazma hatası: %v", err)
		}

		if _, err := file.WriteString("\n"); err != nil {
			return fmt.Errorf("newline yazma hatası: %v", err)
		}

		if domain == "" {
			file.Write([]byte("**********************************************************************\n"))
		}
	} else {
		// herhangi bir değer bulunamadı
		if domain == "" {
			file.Write([]byte("****************** Body üzerinden bilgi bulunamadı *********************"))
		} else {
			file.Write([]byte("****************** Response headerdan bilgi bulunamadı ****************\n"))
		}
	}

	return nil
}

func HeaderChecker(response http.Response) []Struct.ComponentDetail {
	responseHeaders := response.Header
	headerInfo := []Struct.ComponentDetail{}

	headers := make(map[string]string)

	for key, value := range responseHeaders {
		for _, v := range Constant.TechnologyHeaders {
			if strings.ToLower(key) == strings.ToLower(v) && strings.Join(value, "") != "" {

				if strings.ToLower(key) == "set-cookie" {
					for _, cookies := range value {
						splitCookies := strings.Split(cookies, "; ")
						for _, splitCookie := range splitCookies {
							cookieHeader := strings.Split(splitCookie, "=")[0]
							for constantCookieKey, technology := range Constant.CookieTechMap {
								if strings.ToLower(cookieHeader) == strings.ToLower(constantCookieKey) {
									headerInfo = append(headerInfo, Struct.ComponentDetail{
										Component: strings.TrimSpace(technology),
										Version:   "Unknown",
									})
									continue
								}
							}
						}
					}
				} else {
					if strings.Contains(strings.Join(value, " "), "/") {
						splitValue := strings.Split(strings.Join(value, " "), "/")
						if len(splitValue) == 2 {

							headerInfo = append(headerInfo, Struct.ComponentDetail{
								Component: strings.TrimSpace(splitValue[0]),
								Version:   splitValue[1],
							})
						}
					} else {
						headerInfo = append(headerInfo, Struct.ComponentDetail{
							Component: strings.TrimSpace(strings.Join(value, "")),
							Version:   "Unknown",
						})
					}
				}
				headers[key] = strings.Join(value, " ")
			}
		}
	}

	if len(headerInfo) > 0 {
		return headerInfo
	} else {
		return []Struct.ComponentDetail{}
	}

}

func BodyParser(response http.Response) []Struct.ComponentDetail {

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil
	}

	cssFileLinks := CSSLinkExtract(body, response)
	jsFileLinks := JSLinkEXtract(body, response)

	cssComponents := FileRequest(cssFileLinks) // css alanında kullanılan kütüphaneler
	jsComponents := FileRequest(jsFileLinks)   // js alanında kullanılan kütüphaneler

	totalComponents := make(map[string]string)

	for component, _ := range cssComponents {
		totalComponents[component] = ""
	}

	for component, _ := range jsComponents {
		totalComponents[component] = ""
	}

	var componentDetails (map[string]Struct.ComponentDetail)
	componentDetails = make(map[string]Struct.ComponentDetail)

	for component, _ := range totalComponents {
		splitComponent := strings.Split(strings.TrimSpace(component), " ")
		var component string
		var version string
		for index, splitValue := range splitComponent {
			if len(splitComponent)-1 == index {
				version = splitValue
			} else if splitValue != "-" {
				component = component + " " + splitValue
			}
		}

		componentDetails[strings.TrimSpace(component+" "+version)] = Struct.ComponentDetail{
			Component: strings.TrimSpace(component),
			Version:   strings.TrimSpace(version),
		}
	}

	// body üzerinden çıkarılan unique teknoloji değerlerini liste haline getiriyor.
	// map halinde gönderince yazma fonksiyonunda sıkıntı çıkıyor.
	var componentDetailsList []Struct.ComponentDetail
	for _, componentDetail := range componentDetails {
		componentDetailsList = append(componentDetailsList, componentDetail)
	}

	return componentDetailsList
}

func DetectCDNs(domain string) map[string]string {
	results := make(map[string]string)

	client := Client()
	request := Request(domain)

	response, err := client.Do(request)
	if err != nil {
		return nil
	}

	if response.StatusCode == http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		sourceCode := string(body)

		for cdnName, cdnData := range Constant.CDNPatterns {
			// cdnName --> CND sirket ismi
			for _, pattern := range cdnData.Patterns {

				re := regexp.MustCompile(pattern)
				matches := re.FindAllString(sourceCode, -1)

				if len(matches) > 0 {
					results[cdnName] = ""
				}
			}
		}

		if len(results) > 0 {
			return results
		} else {
			return nil
		}
	}
	return nil
}

func AppendCDN(data map[string]string) error {
	file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("dosya açılamadı: %w", err)
	}
	defer file.Close()

	if len(data) == 0 {
		file.Write([]byte("******************************* CDN bulunamadı ***************************************"))
	} else {
		file.Write([]byte("********************************** CDN ********************************\n"))
		for key := range data {
			if _, err := file.WriteString(key + "\n"); err != nil {
				return fmt.Errorf("dosyaya yazma hatası: %w", err)
			}
		}

		file.Write([]byte("***********************************************************************\n"))
	}
	return nil

}

func SnykVulnCheck(component string, version string) (error bool, url string) {

	baseUrl := "https://security.snyk.io/package/npm/" + strings.TrimSpace(component) + "/" + version

	client := Client()
	request := Request(baseUrl)
	response, err := client.Do(request)

	defer response.Body.Close()

	if err != nil {
		return false, ""
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return false, ""
	}
	var booleanResult bool = false

	doc.Find("tbody").Each(func(i int, s *goquery.Selection) {
		class, exists := s.Attr("class")

		if exists && class == "table__tbody" {
			booleanResult = true
		}
	})

	if booleanResult == true {
		return true, baseUrl
	}

	return false, ""
}

func VulnCheck(domain string, componentInfo []Struct.ComponentDetail) string {

	fmt.Println("***************************** Zafiyet taraması " + domain + "*****************************")

	var counter int = 0

	for _, component := range componentInfo {
		for softOld, softSearch := range Constant.SoftwareType {
			if strings.ToLower(softOld) == strings.ToLower(component.Component) {

				isVulnerable, url := SnykVulnCheck(softSearch, component.Version)
				if isVulnerable == true && url != "" {
					counter++
					fmt.Println("ZAFİYET VAR -----> ", url)
				}
			}
		}
	}
	if counter == 0 {
		fmt.Println("***************** Zafiyetli bileşen bulunamadı (Manuel kontrol ediniz) *****************")
	}
	return ""
}
