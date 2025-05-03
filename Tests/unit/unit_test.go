package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"Go/controllers"
	"Go/database"
	"Go/models"
	"Go/scopes"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&models.Category{}, &models.Product{}, &models.Cart{}, &models.CartItem{})
	assert.NoError(t, err)

	database.DB = db
	return db
}

func seedCategoryAndProduct(t *testing.T) (*models.Category, *models.Product) {
	cat := models.Category{Name: "Electronics"}
	assert.NoError(t, database.DB.Create(&cat).Error)

	prod := models.Product{Name: "Laptop", CategoryID: cat.ID, Price: 1500.0}
	assert.NoError(t, database.DB.Create(&prod).Error)
	assert.NoError(t, database.DB.Scopes(scopes.PreloadCategory).First(&prod, prod.ID).Error)

	return &cat, &prod
}

func createTestRequest(t *testing.T, method, url string, body interface{}) (*httptest.ResponseRecorder, echo.Context) {
	var buf bytes.Buffer
	if body != nil {
		assert.NoError(t, json.NewEncoder(&buf).Encode(body))
	}
	req := httptest.NewRequest(method, url, &buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return rec, echo.New().NewContext(req, rec)
}

func decodeJSON(t *testing.T, rec *httptest.ResponseRecorder, out interface{}) {
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), out))
}

func TestCartControllers(t *testing.T) {
	setupTestDB(t)
	_, product := seedCategoryAndProduct(t)

	t.Run("CreateCart", func(t *testing.T) {
		rec, c := createTestRequest(t, http.MethodPost, "/carts", map[string]interface{}{
			"items": []map[string]interface{}{
				{"product_id": product.ID, "amount": 2},
			},
		})
		assert.NoError(t, controllers.CreateCart(c))
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp map[string]interface{}
		decodeJSON(t, rec, &resp)

		assert.NotNil(t, resp["id"])
		items := resp["items"].([]interface{})
		assert.Equal(t, 1, len(items))
		item := items[0].(map[string]interface{})
		assert.Equal(t, float64(product.ID), item["product_id"])
		assert.Equal(t, "Laptop", item["product_name"])
		assert.Equal(t, float64(2), item["amount"])
	})

	t.Run("GetCart", func(t *testing.T) {
		rec, c := createTestRequest(t, http.MethodGet, "/carts/1", nil)
		c.SetParamNames("id")
		c.SetParamValues("1")
		assert.NoError(t, controllers.GetCart(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		decodeJSON(t, rec, &resp)
		assert.Equal(t, float64(1), resp["id"])
		items := resp["items"].([]interface{})
		assert.Equal(t, 1, len(items))
	})

	t.Run("UpdateCart", func(t *testing.T) {
		rec, c := createTestRequest(t, http.MethodPut, "/carts/1", map[string]interface{}{
			"items": []map[string]interface{}{
				{"product_id": product.ID, "amount": 5},
			},
		})
		c.SetParamNames("id")
		c.SetParamValues("1")
		assert.NoError(t, controllers.UpdateCart(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		decodeJSON(t, rec, &resp)
		assert.Equal(t, float64(1), resp["id"])
		items := resp["items"].([]interface{})
		assert.Equal(t, 1, len(items))
		item := items[0].(map[string]interface{})
		assert.Equal(t, float64(product.ID), item["product_id"])
		assert.Equal(t, "Laptop", item["product_name"])
		assert.Equal(t, float64(5), item["amount"])
	})

	t.Run("DeleteCart", func(t *testing.T) {
		rec, c := createTestRequest(t, http.MethodDelete, "/carts/1", nil)
		c.SetParamNames("id")
		c.SetParamValues("1")
		assert.NoError(t, controllers.DeleteCart(c))
		assert.Equal(t, http.StatusNoContent, rec.Code)

		var cart models.Cart
		assert.Error(t, database.DB.First(&cart, 1).Error)
	})
}

func TestCategoryControllers(t *testing.T) {
	setupTestDB(t)

	t.Run("CreateCategory", func(t *testing.T) {
		rec, c := createTestRequest(t, http.MethodPost, "/categories", map[string]interface{}{
			"name": "Books",
		})
		assert.NoError(t, controllers.CreateCategory(c))
		assert.Equal(t, http.StatusCreated, rec.Code)

		var cat models.Category
		decodeJSON(t, rec, &cat)
		assert.Equal(t, "Books", cat.Name)
		assert.NotEqual(t, 0, cat.ID)
	})

	t.Run("GetCategory", func(t *testing.T) {
		var cat models.Category
		database.DB.First(&cat)
		rec, c := createTestRequest(t, http.MethodGet, "/categories/"+strconv.Itoa(int(cat.ID)), nil)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(cat.ID)))
		assert.NoError(t, controllers.GetCategory(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		var fetched models.Category
		decodeJSON(t, rec, &fetched)
		assert.Equal(t, cat.ID, fetched.ID)
		assert.Equal(t, cat.Name, fetched.Name)
	})

	t.Run("UpdateCategory", func(t *testing.T) {
		var cat models.Category
		database.DB.First(&cat)
		rec, c := createTestRequest(t, http.MethodPut, "/categories/"+strconv.Itoa(int(cat.ID)), map[string]interface{}{
			"name": "Novels",
		})
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(cat.ID)))
		assert.NoError(t, controllers.UpdateCategory(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		var updated models.Category
		decodeJSON(t, rec, &updated)
		assert.Equal(t, "Novels", updated.Name)
		assert.Equal(t, cat.ID, updated.ID)
	})

	t.Run("DeleteCategory", func(t *testing.T) {
		var cat models.Category
		database.DB.First(&cat)
		rec, c := createTestRequest(t, http.MethodDelete, "/categories/"+strconv.Itoa(int(cat.ID)), nil)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(cat.ID)))
		assert.NoError(t, controllers.DeleteCategory(c))
		assert.Equal(t, http.StatusNoContent, rec.Code)

		var deleted models.Category
		assert.Error(t, database.DB.First(&deleted, cat.ID).Error)
	})
}

func TestProductControllers(t *testing.T) {
	setupTestDB(t)

	cat := models.Category{Name: "Appliances"}
	assert.NoError(t, database.DB.Create(&cat).Error)

	t.Run("CreateProduct", func(t *testing.T) {
		rec, c := createTestRequest(t, http.MethodPost, "/products", map[string]interface{}{
			"name":        "Washing Machine",
			"category_id": cat.ID,
			"price":       500.0,
		})
		assert.NoError(t, controllers.CreateProduct(c))
		assert.Equal(t, http.StatusCreated, rec.Code)

		var prod models.Product
		decodeJSON(t, rec, &prod)
		assert.Equal(t, "Washing Machine", prod.Name)
		assert.Equal(t, cat.ID, prod.CategoryID)
		assert.Equal(t, 500.0, prod.Price)
		assert.Equal(t, "Appliances", prod.Category.Name)
	})

	t.Run("GetProduct", func(t *testing.T) {
		var prod models.Product
		database.DB.First(&prod)
		rec, c := createTestRequest(t, http.MethodGet, "/products/"+strconv.Itoa(int(prod.ID)), nil)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(prod.ID)))
		assert.NoError(t, controllers.GetProduct(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		var fetched models.Product
		decodeJSON(t, rec, &fetched)
		assert.Equal(t, prod.ID, fetched.ID)
		assert.Equal(t, prod.Name, fetched.Name)
	})

	t.Run("UpdateProduct", func(t *testing.T) {
		var prod models.Product
		database.DB.First(&prod)
		rec, c := createTestRequest(t, http.MethodPut, "/products/"+strconv.Itoa(int(prod.ID)), map[string]interface{}{
			"name":        "Washer",
			"category_id": cat.ID,
			"price":       450.0,
		})
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(prod.ID)))
		assert.NoError(t, controllers.UpdateProduct(c))
		assert.Equal(t, http.StatusOK, rec.Code)

		var updated models.Product
		decodeJSON(t, rec, &updated)
		assert.Equal(t, "Washer", updated.Name)
		assert.Equal(t, cat.ID, updated.CategoryID)
		assert.Equal(t, 450.0, updated.Price)
	})

	t.Run("DeleteProduct", func(t *testing.T) {
		var prod models.Product
		database.DB.First(&prod)
		rec, c := createTestRequest(t, http.MethodDelete, "/products/"+strconv.Itoa(int(prod.ID)), nil)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(int(prod.ID)))
		assert.NoError(t, controllers.DeleteProduct(c))
		assert.Equal(t, http.StatusNoContent, rec.Code)

		var deleted models.Product
		assert.Error(t, database.DB.First(&deleted, prod.ID).Error)
	})
}
