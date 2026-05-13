package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

//go:embed ascii/*
//go:embed public/*
var content embed.FS

func ServeFile(file_name string) http.HandlerFunc {
	file_data, _ := content.ReadFile(file_name)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write(file_data)
	}
}

func CreateHits(num int) int {
	os.Remove("hits")
	f, _ := os.Create("hits")
	fmt.Fprintf(f, "%d", num)
	f.Close()
	return num
}

func main() {
	var hits = 0

	// read hits file if it exists, create one if it does not exist
	hits_data, err := os.ReadFile("hits")
	if err == nil {
		readhits, e := strconv.Atoi(string(hits_data))
		if e == nil {
			hits = readhits
		} else {
			print("failed to read converted hits data")
			hits = CreateHits(0)
		}
	} else {
		if !os.IsExist(err) {
			hits = CreateHits(0)
		}
	}

	print("read hits data: ", hits, "\n")

	serve_html := func(w http.ResponseWriter, req *http.Request) {
		temp := template.New("ascii.html")
		temp.ParseFS(content, "ascii/ascii.html")

		var decimal_now = (float32(time.Now().Year()) + float32(time.Now().YearDay())/float32(365))
		var rough_age = decimal_now - 2004.626
		if req.URL.EscapedPath() == "/" {
			hits++
			CreateHits(hits)
		}
		data := struct {
			Age  string
			Hits string
		}{Age: fmt.Sprintf("%.3f", rough_age), Hits: fmt.Sprintf("%06d", hits)}

		temp.Execute(w, data)
	}

	serve_icons := func(w http.ResponseWriter, _ *http.Request) {
		icon, _ := content.ReadFile(fmt.Sprintf("public/favicon%d.ico", (hits%5)+1))
		w.Write(icon)
	}

	http.HandleFunc("/script.js", ServeFile("ascii/script.js"))
	http.HandleFunc("/counterbase.png", ServeFile("public/counterbase.png"))
	http.HandleFunc("/manifest.json", ServeFile("public/manifest.json"))
	http.HandleFunc("/logo512.png", ServeFile("public/logo512.png"))
	http.HandleFunc("/robots.txt", ServeFile("public/robots.txt"))
	http.HandleFunc("/favicon.ico", serve_icons)
	http.HandleFunc("/", serve_html)

	var port = 8080
	if len(os.Args) > 1 {
		num, e := strconv.Atoi(os.Args[1])
		if e == nil {
			port = num
		} else {
			print("error reading port, defaulting to 8080: ", e)
		}
	}

	print("starting net server on port ", port, " ...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
