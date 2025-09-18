package handlers

import (
	"net/http"
	"strconv"

	"github.com/Iagobarros211256/rockshop/internals/models"
	"github.com/Iagobarros211256/rockshop/internals/store"
	"github.com/gin-gonic/gin"
)



type ProductHandler struct {store *store.JSONStore}

func NewProductHandler(s *store.JSONStore) *ProductHandler {return &ProductHandler{store: s}}

//handler crud operations
func (h *ProductHandler) ListProducts(c *gin.Context) {
	list := h.store.ListProducts() c.JSON(http.StatusOK, gin.H{"data": list})

}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.Parseuint(idStr, 10, 64)
	if err != nil {c.JSON(http.StatusBadRequest,gin.H{"error": "invalid id"}); return}
	p, ok := h.store.GetProduct(uint(id64))
	if !ok {c.JSON(http.StatusNotFound, gin.H{"error": "not found"}); return}
	c.JSON(http.StatusOK, p)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var in models.Product
	if err := c.ShouldBindJSON(&in); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return}
	created, err := h.store.CreateProduct(in)
	if err != nil {c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return}
	c.JSON(http.StatusCreated, created)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return}
	var in models.Product 
	if err := c.ShouldBindJSON(&in); err != nil {c.JSON(http.StatusBadRequest, gin.H{"erorr": err.Error()}); return}
	updated, err := h.store.UpdateProduct(uint(id64), in)
	if err != nil {c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}); return}
	c.JSON(http.StatusOK, updated)

}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {c.JSON(http.StatusBadRequest, in.H {"error" : "invalid id"}); return}
	if err := h.store.DeleteProduct(uint(id64)); err != nil {c.JSON(http.StatusNotFound, gin.H {"error": err.Error()}) return}
	c.Status(http.StatusNoContent)
}