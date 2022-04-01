package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

type Page struct {
    Title string
    Body  []byte
}

func (p *Page) save() error {
    filename := p.Title + ".txt"
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title + ".txt"
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[len("/view/"):]
    p, _ := loadPage(title)
    fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func attendHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf("ATTEND");
}
func enterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint("ENTER");
}
func leaveHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint("LEAVE");
}

func main() {
    http.HandleFunc("/attend", attendHandler)
    http.HandleFunc("/enter", enterHandler)
    http.HandleFunc("/leave", leaveHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}


