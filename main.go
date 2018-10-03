package main

import (
  "regexp"
  "log"
  "fmt"
  "net/http"
	"html/template"
  "github.com/ndvo/roman"
)
// Marca em IDs e Classes a estrutura de uma lei construída em observância à 
// LEI COMPLEMENTAR Nº 95, DE 26 DE FEVEREIRO DE 1998


// O termo ‘dispositivo’ mencionado nesta Lei refere-se a artigos, parágrafos, incisos, alíneas ou itens.

func textToHTML(t string) (h string){
  if m, _ := regexp.MatchString(".", t); m {
    h = marcaAgrupamentos(
        marcaListas(
        marcaAlineas(
        marcaIncisos(
        marcaEpigrafe(
        marcaParagrafos(
        marcaArtigos(t)))))))

  }else{
    h = "nope"
  }
  return 
}

// Art. 4º A epígrafe, grafada em caracteres maiúsculos, propiciará identificação numérica singular à lei e será formada pelo título designativo da espécie normativa, pelo número respectivo e pelo ano de promulgação.
func marcaEpigrafe(t string)(tm string){
  r, _ := regexp.Compile(`(LEI|DECRETO|PORTARIA)\sN[oªº°]\s+((?:[0-9]{1,3}\.?)+),\s+DE\s+([1-3]?[0-9])\sDE\s([A-Z]+)\sDE\s([0-9]{4})`)
  tm = r.ReplaceAllString(t, "<h1 id=\"epigrafe\">$1 Nº $2, DE $3 DE $4 DE $5</h1>" )
  return
}


//Art. 5º A ementa será grafada por meio de caracteres que a realcem e explicitará, de modo conciso e sob a forma de título, o objeto da lei.
func marcaEmenta(t string)(tm string){
  //TODO: Identificar de forma mais específica a Ementa
  r, _ := regexp.Compile(`\s*\p{Lu}\p{Ll}*[\p{Zs}\p{P}\p{L}]*\p{Po}\p{Z}*\`)
  tm = r.ReplaceAllString(t, "<p class=\"ementa\">$0</p>")
  return
}

// Art. 6º O preâmbulo indicará o órgão ou instituição competente para a prática do ato e sua base legal.

// Art. 10
// V - o agrupamento de artigos poderá constituir Subseções; o de Subseções, a Seção; o de Seções, o Capítulo; o de Capítulos, o Título; o de Títulos, o Livro e o de Livros, a Parte;
// VI - os Capítulos, Títulos, Livros e Partes serão grafados em letras maiúsculas e identificados por algarismos romanos, podendo estas últimas desdobrar-se em Parte Geral e Parte Especial ou ser subdivididas em partes expressas em numeral ordinal, por extenso;
// VII - as Subseções e Seções serão identificadas em algarismos romanos, grafadas em letras minúsculas e postas em negrito ou caracteres que as coloquem em realce;
// VIII - a composição prevista no inciso V poderá também compreender agrupamentos em Disposições Preliminares, Gerais, Finais ou Transitórias, conforme necessário.

// marcaAgrupamentos identifica Subsecoes, Secoes, Capitulo, Titulo, Livro, Parte
func marcaAgrupamentos(t string) (tm string){
  r, _ := regexp.Compile(`(?m:^\s*(Parte|Livro|Título|Capítulo|Seção|Subseção)\s+([IXVLC]+)\s*\n\s*([\p{L} ]+)\s*\n)`)
  tm = r.ReplaceAllString(t, "<div class=\"$1\">\n\t<p class=\"agrupamento-tipo $1-$2\">$1 $2</p>\n\t<p class=\"agrupamento-nome\">$3</p>\n</div>\n")
  return
}


// Art. 10 I - a unidade básica de articulação será o artigo, indicado pela abreviatura "Art.", seguida de numeração ordinal até o nono e cardinal a partir deste;
// b) é vedada, mesmo quando recomendável, qualquer renumeração de artigos e de unidades superiores ao artigo, referidas no inciso V do art. 10, devendo ser utilizado o mesmo número do artigo ou unidade imediatamente anterior, seguido de letras maiúsculas, em ordem alfabética, tantas quantas forem suficientes para identificar os acréscimos;                    (Redação dada pela Lei Complementar nº 107, de 26.4.2001)
func marcaArtigos(t string) (tm string){
  r, _ := regexp.Compile(`(?m:^\s*[Aa]rt[\.\s]\s*(\d+)[oªº°]?[\.\s]\s*(.*)$)`)
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
// IV - os incisos serão representados por algarismos romanos, as alíneas por letras minúsculas e os itens por algarismos arábicos;
//  d) é admissível a reordenação interna das unidades em que se desdobra o artigo, identificando-se o artigo assim modificado por alteração de redação, supressão ou acréscimo com as letras ‘NR’ maiúsculas, entre parênteses, uma única vez ao seu final, obedecidas, quando for o caso, as prescrições da alínea "c".                         (Redação dada pela Lei Complementar nº 107, de 26.4.2001)
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
  r, _ := regexp.Compile(`(?m:^\s*(<li[^>]*>[^<]*</\s*li>[\s\n]*)+)[\s\n]*`)
  tm = r.ReplaceAllString(t, "<ol>$0</ol>\r\r")
  return
}

// marcaParagrafos marca parágrafos com classes e id
func marcaParagrafos(t string) (tm string){
  r, _ := regexp.Compile(`(?m:^\s*§\s*(\d+)[ºoª]\.?\s+(.*))`)
  tm = r.ReplaceAllString(t, "<p class=\"paragrafo\" id=\"\">§$1 $2</p> ")
  return
}

func api(w http.ResponseWriter, r *http.Request) {
  //w.Header().Set("Content-Type", "application/txt; charset=utf-8")
	switch r.Method {
	case "GET":
    t, _ := template.ParseFiles("templates/norma.html")
    m:= "Conversor txt para html"
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
