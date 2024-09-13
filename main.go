package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type CounterHandler struct {
	counter int
}

func (ct *CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println(ct.counter)
	ct.counter++
	fmt.Fprintln(w, "Counter:", ct.counter)
}

func HelloHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello, there\n")
}

func StatusHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		w.Write([]byte("404 - not found\n"))
		return
	}
	fmt.Fprintln(w, "home page")
}

func UserAgentHandler(w http.ResponseWriter, r *http.Request) {
	ua := r.Header.Get("User-Agent")
	fmt.Fprintf(w, "User agent: %s\n", ua)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!\n", r.URL.Path[1:])
}

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["name"]
	name := "guest"
	if ok {
		name = keys[0]
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

func Process(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "form.html")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		name := r.FormValue("name")
		occupation := r.FormValue("occupation")
		fmt.Fprintf(w, "%s is a %s\n", name, occupation)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadFile("sid.png")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "image/png")
	w.Write(buf)
}

type User struct {
	Name       string
	Occupation string
}

type Data struct {
	Users []User
}

func main() {

	th := &CounterHandler{counter: 0}
	http.Handle("/count", th)

	http.HandleFunc("/", HelloHandler)

	http.HandleFunc("/status", StatusHandler)

	http.HandleFunc("/", NotFoundHandler)

	http.HandleFunc("/ua", UserAgentHandler)

	http.HandleFunc("/hello", HelloServer)

	http.HandleFunc("/query", QueryHandler)

	http.HandleFunc("/form", Process)

	http.HandleFunc("/image", HandleRequest)

	tmp := template.Must(template.ParseFiles("layout.html"))
	http.HandleFunc("/users", func(w http.ResponseWriter, _ *http.Request) {
		data := Data{
			Users: []User{
				{Name: "John Doe", Occupation: "gardener"},
				{Name: "Roger Roe", Occupation: "driver"},
				{Name: "Peter Smith", Occupation: "teacher"},
			},
		}
		tmp.Execute(w, data)
	})

	fileServer := http.FileServer(http.Dir("./public"))
	http.Handle("/public/", http.StripPrefix("/public", fileServer))

	log.Println("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
