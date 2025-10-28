package catalog_service

import (
	"github.com/PoojaSrinivasan18/catalog-service/database"
	"github.com/PoojaSrinivasan18/catalog-service/model"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetProductById(c *gin.Context) {
	productId, err := strconv.Atoi(c.Query("productId"))
	if err != nil {
		log.Errorf("Invalid product ID: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	var existingProductDetail model.ProductModel
	database := database.GetDB()

	t := database.Where("product_id=?", productId).First(&existingProductDetail)
	if t.Error != nil {
		log.Errorf("DB query error %v", t.Error)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": t.Error})
		return
	}

	c.IndentedJSON(http.StatusOK, existingProductDetail)
}
func GetAllProducts(c *gin.Context) {
	var products []model.ProductModel
	db := database.GetDB()

	t := db.Find(&products)
	if t.Error != nil {
		log.Errorf("DB query error %v", t.Error)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": t.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, products)
}

func AddProduct(c *gin.Context) {
	var productModel model.ProductModel
	err := c.ShouldBind(&productModel)
	if err != nil {
		log.Errorf("FORM binding error %v", err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	tx := database.GetDB().Create(&productModel)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error adding product"})
		return
	}

	c.IndentedJSON(http.StatusOK, productModel)
}
func DeleteProduct(c *gin.Context) {
	productId, err := strconv.Atoi(c.Query("productId"))
	if err != nil {
		log.Errorf("Invalid product ID: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	var existingProductDetail model.ProductModel
	database := database.GetDB()

	t := database.Where("product_id=?", productId).First(&existingProductDetail)
	if t.Error != nil {
		log.Errorf("DB query error %v", t.Error)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": t.Error})
		return
	}

	tx := database.Model(&existingProductDetail).Delete(existingProductDetail)
	if tx.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error saving product data"})
		return
	}

	c.IndentedJSON(http.StatusOK, "Product deleted successfully")
}
func UpdateProduct(c *gin.Context) {
	// Get product_id from query params
	productID, err := strconv.Atoi(c.Query("productId"))
	if err != nil {
		log.Errorf("Invalid product ID: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid product ID"})
		return
	}

	// Bind incoming JSON body to struct
	var updatedData model.ProductModel
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		log.Errorf("JSON binding error: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	db := database.GetDB()

	// Check if product exists
	var existingProduct model.ProductModel
	if err := db.First(&existingProduct, productID).Error; err != nil {
		log.Errorf("Product not found: %v", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	// Update fields
	existingProduct.Name = updatedData.Name
	existingProduct.Price = updatedData.Price
	existingProduct.Category = updatedData.Category
	existingProduct.IsActive = updatedData.IsActive
	existingProduct.Sku = updatedData.Sku

	// Save the updated record
	if err := db.Save(&existingProduct).Error; err != nil {
		log.Errorf("Error updating product: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Error updating product"})
		return
	}

	c.IndentedJSON(http.StatusOK, existingProduct)
}
