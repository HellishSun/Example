package main

import (
	"fmt"
	"html/template"
	// "log"
	"math/rand"
	"net/http"
	"net/url"
	"unicode/utf8"
	"github.com/gorilla/mux"
)


//ВАЖНО! СДЕЛАТЬ ПРОВЕРКУ НА СУЩЕСТВОВАНИЕ СОКРАЩЕННОЙ ССЫЛКИ

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

var MapKeyAdderss map[string]string // карта для хранения ключей ввида: map[короткая_ссылка:длинная_ссылка] [ключ]ссылка

var MapUrlAddress map[string]string // карта для харнения длинных ссылок. Имеет ввид map[длинная_ссылка:короткая_ссылка] ключ:значение
// короткая ссылка имеет ввид: to/<сгенерированный_случайныйКод>

type Result struct {
	Link   string //отвечает за URL, который поступил на форму
	Code   string //это сформированная строка, которую мы сохраним в MAPe !КЛЮЧ!
	Status string //будет заполняться в соответствии с  тем, какой результат будет.
}

func shorting() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return "to/" + string(b) // возвращаем сформированную строку, по которой мы будем определять ссылку !КЛЮЧ!
}

func isValidUrl(token string) bool { //проверяет, является ли введённый url корректным, путём простых проверок.
	_, err := url.ParseRequestURI(token)
	if err != nil {
		return false
	}
	u, err := url.Parse(token)
	if err != nil || u.Host == "" {
		return false
	}
	return true
}


func indexPage(w http.ResponseWriter, r *http.Request) {
	
	templ, _ := template.ParseFiles("./templates/index.html")
	result := Result{}

	if r.Method == "POST" {
		if LenStr(r.FormValue("s")){ // true, если длинна < 356 рун
			result.Link = r.FormValue("s")
			if !isValidUrl(result.Link){ 
				result.Status="Ссылка имеет не правильный формат"
				result.Link=""
			}else if isValidUrl(result.Link){
			// пройдя основные проверки мы получаем нормальную ссылку

				// провекра на существование в карте ссылок
				if codeFromUrl, BeUrl := MapUrlAddress[result.Link]; BeUrl{
					// вытаскиваем по ссылке ключ из карты Ссылок и ищем его в карте ключей
					if linkFromKey, BeCode := MapKeyAdderss[codeFromUrl]; BeCode{
						// делаем проверку на соответсвие карт
						if linkFromKey == result.Link{
							fmt.Println("Успешное соответсвие")
							result.Status = "Такая ссылка существует"
							result.Code = MapUrlAddress[result.Link]
						}
					}else{
						fmt.Println("Если ключа не существет, то не может сущестовать ")
						result.Status = "Ошибка: ссылка существует, ключ нет"
						result.Code=""
					}
				
				}else if !BeUrl{//если ссылки не существует. Мы записываем ее
					result.Code = shorting()// генерируем ключ

					if _, be := MapKeyAdderss[result.Code]; be{
						fmt.Println("Если ссылки не существет, то не может сущестовать ключа")
						result.Status = "Ошибка: ключ существует, а ссылка нет"
						result.Code = ""
					} else if !be {
						//тут уже будет происходить добавление новой ссылки и ключа
						MapUrlAddress[result.Link] = result.Code
						MapKeyAdderss[result.Code] = result.Link

						//вывод данных на страницу
						result.Status = "Ссылка с ключем добавленна"

					}
				}
			}   
			
		}else{
			result.Status="Длинна ссылки превышает 355 символов"
			result.Code=""
		}
	}

	templ.Execute(w, result)
}


func redirectTo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	fmt.Printf("Это vars: '%s'\n", vars)
	fmt.Printf("Это MapUrlAddress: '%s'\n", MapUrlAddress)
	fmt.Printf("Это MapKeyAdderss: '%s'\n", MapKeyAdderss)

	key := string(vars["key"]) // ключ по которому мы будем осуществлять поиск длинной ссылки в МАPe
	fmt.Printf("this is KEY: '%s'\n ", key)

	if link, be := MapKeyAdderss["to/"+key]; be{ // ключ в MapKeyAdderss имеет вид: to/key
		fmt.Fprintf(w, "<script>location='%s';</script>", link)
	}else{
		link = "http://localhost:8000/404"
		fmt.Fprintf(w, "<script>location='%s';</script>", link)
	}
}

func LenStr(url string) bool{// return true if length string less than 350 symbol
	len := utf8.RuneCountInString(url) // выводим кол-ов рун в строке	
	if len < 356{
		return true
	}else{
		return false
	}
}

func NotFound(w http.ResponseWriter, r *http.Request){
	tmpl, _ := template.ParseFiles("./templates/404.html")
	tmpl.Execute(w, nil)
}

func main() {

	MapUrlAddress = make(map[string]string) // инициализация МАРы. Без этого она будет пустой и будет выдавать ошибки
	MapKeyAdderss = make(map[string]string)
	router := mux.NewRouter()


	router.HandleFunc("/404", NotFound)
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/to/{key}", redirectTo) // отправляем карту ввида map["key":]
	
	http.ListenAndServe(":8000", router)


}
