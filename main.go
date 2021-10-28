package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var productList []Product

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productId, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	product, listItemIndex := findProductById(productId)
	if product == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		productJson, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productJson)
	case http.MethodPut:
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedProduct.Id != productId {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		product = &updatedProduct
		productList[listItemIndex] = *product
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func findProductById(id int) (*Product, int) {
	for i, product := range productList {
		if product.Id == id {
			return &product, i
		}
	}
	return nil, 0
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productsJson, err := json.Marshal(productList)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJson)
	case http.MethodPost:
		var newProduct Product

		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyBytes, &newProduct)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newProduct.Id != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newProduct.Id = getNextId()

		productList = append(productList, newProduct)

		w.WriteHeader(http.StatusCreated)

		return
	}
}

func getNextId() int {
	highestId := 1

	for _, product := range productList {
		if highestId < product.Id {
			highestId = product.Id
		}
	}

	return highestId + 1
}

func init() {
	productsJson := `[
		{
			"id": 1,
			"name": "phone"
		},
		{
			"id": 2,
			"name": "laptop"
		}
	]`

	if err := json.Unmarshal([]byte(productsJson), &productList); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/products/", productHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
