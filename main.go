package main

import (
  "regexp"
  "log"
  "fmt"
  "net/http"
	"html/template"
  "github.com/ndvo/roman"
)


func textToHTML(t string) (h string){
  if m, _ := regexp.MatchString(".", t); m {
    h = marcaListas(
        marcaAlineas(
        marcaIncisos(
        marcaArtigos(t))))

  }else{
    h = "nope"
  }
  return 
}


func marcaArtigos(t string) (tm string){
  r, _ := regexp.Compile(`(?m:^\s*[Aa]rt[\.\s]\s*(\d+)[oªº]?[\.\s]\s*(.*)$)`)
  tm = r.ReplaceAllString(t, "<p class=\"art\" id=\"art$1\" >Art. ${1} ${2}</p>")
  return
}

// replaceIncisos converte o texto de um inciso em html, aplicando a conversão
// de números romanos para arábicos.
func replaceIncisos(match string) string {
  r, _ := regexp.Compile(`(?m:^\s*([IXVLC]+)\s*[-\s\.]\s*(.*)$)`)
  submatches := r.FindStringSubmatch(match)
  arabic, err := roman.ToInt(submatches[1])
  if err == nil {
    println("roman is ", submatches[1], "arabic is", arabic)
    return  fmt.Sprintf(
      `<li value="%d"type="I" class="inciso" id="inciso%d" >%s</li>`,
        arabic,
        arabic,
        submatches[2])
  }else{
    println("failed to convert ", submatches[1], "to int")
    return  match
  }
}


// marcaIncisos identifica os incisos a partir de regex e os substitui por HTML li do tipo I
func marcaIncisos(t string) (tm string){
  r, _ := regexp.Compile(`(?m:^\s*([IXVLC]+)\s*[-\s\.]\s*(.*)$)`)
  //tm = r.ReplaceAllStringFunc(t, "<li type=\"I\" class=\"inciso\" id=\"inciso$1\" >${2}</li>")
  tm = r.ReplaceAllStringFunc(t, replaceIncisos)
  return
}

// marcaIncisos identifica as alíneas a partir de regex e os substitui por HTML li do tipo a
func marcaAlineas(t string) (tm string){
  r, _ := regexp.Compile(`(?m:^\s*([a-z])\s*[)\.-]\s*(.*)$)`)
  tm = r.ReplaceAllString(t, "<li type=\"a\" class=\"alinea\" id=\"alinea$1\" >${2}</li>")
  return
}

// marcaListas identifica listas e as cerca com a tag ol
// presume-se que não há listas não ordenadas na legislação
func marcaListas(t string) (tm string){
  // Multiple lines group
  // Set a group from the first <li> until the first non <li> element
  // Surround the group with <ol>
  r, _ := regexp.Compile(`(?m:^\s*(<li[^>]*>[^<]*</\s*li>[\s\n]*)+)`)
  tm = r.ReplaceAllString(t, "<ol>$0</ol>")
  return
}


func api(w http.ResponseWriter, r *http.Request) {
  //w.Header().Set("Content-Type", "application/txt; charset=utf-8")
	switch r.Method {
	case "GET":
    t, _ := template.ParseFiles("templates/norma.html")
    m:= "Convert txt to html"
    t.Execute(w, m)

	case "POST":
		r.ParseForm()
		original := r.Form["texto"][0]

    fmt.Fprint(w, textToHTML(original))
  }
}

func main() {
	http.HandleFunc("/", api)
	log.Fatal(http.ListenAndServe(":8080", nil))
	println("Server is up and running.")
}
