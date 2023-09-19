package handles

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/renan5g/go-expert/internal/dto"
	"github.com/renan5g/go-expert/internal/entity"
	"github.com/renan5g/go-expert/internal/infra/database"
	entityPkg "github.com/renan5g/go-expert/pkg/entity"
)

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{ProductDB: db}
}

// Create product godoc
// @Summary 		Create product
// @Description Create products
// @Tags 				products
// @Accept 			json
// @Produce 		json
// @Param 			request body dto.CreateProductInput true "product request"
// @Success 		201
// @Failure 		500 {object} Error
// @Router 			/products [post]
// @Security		ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get product godoc
// @Summary 		Get product
// @Description Get product
// @Tags 				products
// @Accept 			json
// @Produce 		json
// @Param 			id 				path 			string 	true 	"product ID" Format(uuid)
// @Success 		200				{object}	entity.Product
// @Failure 		404
// @Failure 		500 			{object} 	Error
// @Router 			/products/{id} [get]
// @Security		ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := h.ProductDB.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// Update product godoc
// @Summary 		Update product
// @Description Update products
// @Tags 				products
// @Accept 			json
// @Produce 		json
// @Param 			id 				path 			string 									true 	"product ID" Format(uuid)
// @Param 			request 	body 			dto.CreateProductInput 	true 	"product request"
// @Success 		200
// @Failure 		404
// @Failure 		500 {object} Error
// @Router 			/products/{id} [put]
// @Security		ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	product.ID, err = entityPkg.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	if _, err := h.ProductDB.FindById(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Delete product godoc
// @Summary 		Delete product
// @Description Delete products
// @Tags 				products
// @Accept 			json
// @Produce 		json
// @Param 			id 		path 			string	true 	"product ID" Format(uuid)
// @Success 		200
// @Failure 		404
// @Failure 		500 	{object} 	Error
// @Router 			/products/{id} [delete]
// @Security		ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := h.ProductDB.FindById(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	err := h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// List products godoc
// @Summary 		List products
// @Description List all products
// @Tags 				products
// @Accept 			json
// @Produce 		json
// @Param 			page 			query 		string 	false 	"page number"
// @Param 			limit 		query 		string 	false 	"limit"
// @Param 			sort			query 		string 	false 	"sort by asc or desc"
// @Success 		200				{array}		entity.Product
// @Failure 		404
// @Failure 		500 			{object} 	Error
// @Router 			/products [get]
// @Security		ApiKeyAuth
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	products, err := h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
