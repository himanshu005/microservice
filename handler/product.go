package handler

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/himanshu005/microservice/data"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.GetProducts(w, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}

	if r.Method == http.MethodPut {
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("invalid id", idString, id)
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		p.l.Println("valid id", id)
		p.updateProduct(id, w, r)
		return
	}

	//Method are not supported
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) GetProducts(w http.ResponseWriter, r *http.Request) {
	p.l.Println("GET Data method")
	lp := data.GetProducts()

	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal data", http.StatusInternalServerError)
	}
}

func (p *Product) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Post Data method")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unnable to unmarshal JSON", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p *Product) updateProduct(id int, w http.ResponseWriter, r *http.Request) {
	p.l.Println("Update Data method")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Unnable to unmarshal JSON", http.StatusBadRequest)
	}
	p.l.Printf("data %#v", prod)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
	}

	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
	}

}
