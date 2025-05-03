package controllers_test

import (
	"Go/models"
	"Go/routes"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func doRequest(e *echo.Echo, method, path string, body interface{}) (*httptest.ResponseRecorder, *http.Request) {
	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		assert.NoError(nil, err, "Kod powinien poprawnie serializować dane wejściowe")
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	return rec, req
}

func TestCartAPI(t *testing.T) {
	setupTestDB(t)
	e := echo.New()
	routes.CartRoutes(e)
	_, prod := seedCategoryAndProduct(t)

	t.Run("POST /carts positive", func(t *testing.T) {
		payload := map[string]interface{}{
			"items": []map[string]interface{}{
				{"product_id": prod.ID, "amount": 3},
			},
		}
		rec, req := doRequest(e, http.MethodPost, "/carts", payload)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code, "Tworzenie koszyka – pozytywny scenariusz")

		var response map[string]interface{}
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
		assert.NotNil(t, response["id"], "Odpowiedź powinna zawierać pole 'id'")
		items := response["items"].([]interface{})
		assert.Equal(t, 1, len(items), "Koszyk powinien zawierać 1 element")
	})

	t.Run("POST /carts negative - invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/carts", bytes.NewBufferString("invalid-json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Przy niepoprawnym JSON oczekujemy 400")
	})

	t.Run("GET /carts/:id positive", func(t *testing.T) {
		payload := map[string]interface{}{
			"items": []map[string]interface{}{
				{"product_id": prod.ID, "amount": 1},
			},
		}
		recCreate, reqCreate := doRequest(e, http.MethodPost, "/carts", payload)
		e.ServeHTTP(recCreate, reqCreate)
		var createResp map[string]interface{}
		assert.NoError(t, json.Unmarshal(recCreate.Body.Bytes(), &createResp))
		idStr := strconv.Itoa(int(createResp["id"].(float64)))

		rec, req := doRequest(e, http.MethodGet, "/carts/"+idStr, nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "Pobranie koszyka – pozytywny scenariusz")
	})

	t.Run("GET /carts/:id negative", func(t *testing.T) {
		rec, req := doRequest(e, http.MethodGet, "/carts/9999", nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code, "Pobranie koszyka o nieistniejącym ID powinno zwrócić 404")
	})

	t.Run("PUT /carts/:id positive", func(t *testing.T) {
		payload := map[string]interface{}{
			"items": []map[string]interface{}{
				{"product_id": prod.ID, "amount": 2},
			},
		}
		recCreate, reqCreate := doRequest(e, http.MethodPost, "/carts", payload)
		e.ServeHTTP(recCreate, reqCreate)
		var createResp map[string]interface{}
		assert.NoError(t, json.Unmarshal(recCreate.Body.Bytes(), &createResp))
		idStr := strconv.Itoa(int(createResp["id"].(float64)))

		updatePayload := map[string]interface{}{
			"items": []map[string]interface{}{
				{"product_id": prod.ID, "amount": 7},
			},
		}
		rec, req := doRequest(e, http.MethodPut, "/carts/"+idStr, updatePayload)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "Aktualizacja koszyka – pozytywny scenariusz")

		var updateResp map[string]interface{}
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &updateResp))
		items := updateResp["items"].([]interface{})
		item := items[0].(map[string]interface{})
		assert.Equal(t, float64(7), item["amount"], "Zaktualizowana ilość powinna wynosić 7")
	})

	t.Run("PUT /carts/:id negative", func(t *testing.T) {
		updatePayload := map[string]interface{}{
			"items": []map[string]interface{}{
				{"product_id": prod.ID, "amount": 4},
			},
		}
		rec, req := doRequest(e, http.MethodPut, "/carts/9999", updatePayload)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code, "Aktualizacja nieistniejącego koszyka powinna zwrócić 404")
	})

	t.Run("DELETE /carts/:id positive", func(t *testing.T) {
		payload := map[string]interface{}{
			"items": []map[string]interface{}{
				{"product_id": prod.ID, "amount": 1},
			},
		}
		recCreate, reqCreate := doRequest(e, http.MethodPost, "/carts", payload)
		e.ServeHTTP(recCreate, reqCreate)
		var createResp map[string]interface{}
		assert.NoError(t, json.Unmarshal(recCreate.Body.Bytes(), &createResp))
		idStr := strconv.Itoa(int(createResp["id"].(float64)))

		rec, req := doRequest(e, http.MethodDelete, "/carts/"+idStr, nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code, "Usunięcie koszyka – pozytywny scenariusz")
	})

	t.Run("DELETE /carts/:id negative", func(t *testing.T) {
		rec, req := doRequest(e, http.MethodDelete, "/carts/9999", nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code, "Usunięcie nieistniejącego koszyka powinno zwrócić 404")
	})
}

func TestCategoryAPI(t *testing.T) {
	setupTestDB(t)
	e := echo.New()
	routes.CategoryRoutes(e)

	t.Run("POST /categories positive", func(t *testing.T) {
		payload := map[string]interface{}{
			"name": "Books",
		}
		rec, req := doRequest(e, http.MethodPost, "/categories", payload)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code, "Tworzenie kategorii – pozytywny scenariusz")

		var cat models.Category
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &cat))
		assert.Equal(t, "Books", cat.Name)
		assert.NotEqual(t, 0, cat.ID, "ID powinno być ustawione")
	})

	t.Run("POST /categories negative", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBufferString("invalid-json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Niepoprawny JSON powinien dać 400")
	})

	t.Run("GET /categories/:id positive", func(t *testing.T) {
		payload := map[string]interface{}{"name": "Music"}
		recCreate, reqCreate := doRequest(e, http.MethodPost, "/categories", payload)
		e.ServeHTTP(recCreate, reqCreate)
		var cat models.Category
		assert.NoError(t, json.Unmarshal(recCreate.Body.Bytes(), &cat))
		idStr := strconv.Itoa(int(cat.ID))

		rec, req := doRequest(e, http.MethodGet, "/categories/"+idStr, nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "Pobranie istniejacej kategorii – pozytywny scenariusz")
	})

	t.Run("GET /categories/:id negative", func(t *testing.T) {
		rec, req := doRequest(e, http.MethodGet, "/categories/9999", nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code, "Pobranie kategorii o nieistniejącym ID powinno zwrócić 404")
	})

	t.Run("PUT /categories/:id positive", func(t *testing.T) {
		payload := map[string]interface{}{"name": "Clothes"}
		recCreate, reqCreate := doRequest(e, http.MethodPost, "/categories", payload)
		e.ServeHTTP(recCreate, reqCreate)
		var cat models.Category
		assert.NoError(t, json.Unmarshal(recCreate.Body.Bytes(), &cat))
		idStr := strconv.Itoa(int(cat.ID))

		updatePayload := map[string]interface{}{"name": "Fashion"}
		rec, req := doRequest(e, http.MethodPut, "/categories/"+idStr, updatePayload)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "Aktualizacja kategorii – pozytywny scenariusz")
		var updatedCat models.Category
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &updatedCat))
		assert.Equal(t, "Fashion", updatedCat.Name)
	})

	t.Run("PUT /categories/:id negative", func(t *testing.T) {
		updatePayload := map[string]interface{}{"name": "Nonexistent"}
		rec, req := doRequest(e, http.MethodPut, "/categories/9999", updatePayload)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code, "Aktualizacja nieistniejącej kategorii powinna zwrócić 404")
	})

	t.Run("DELETE /categories/:id positive", func(t *testing.T) {
		// Utwórz, potem usuń kategorię.
		payload := map[string]interface{}{"name": "Toys"}
		recCreate, reqCreate := doRequest(e, http.MethodPost, "/categories", payload)
		e.ServeHTTP(recCreate, reqCreate)
		var cat models.Category
		assert.NoError(t, json.Unmarshal(recCreate.Body.Bytes(), &cat))
		idStr := strconv.Itoa(int(cat.ID))
		rec, req := doRequest(e, http.MethodDelete, "/categories/"+idStr, nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code, "Usunięcie kategorii – pozytywny scenariusz")
	})

	t.Run("DELETE /categories/:id negative", func(t *testing.T) {
		rec, req := doRequest(e, http.MethodDelete, "/categories/9999", nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code, "Usunięcie nieistniejącej kategorii powinno zwrócić 404")
	})
}

func TestProductAPI(t *testing.T) {
	setupTestDB(t)
	e := echo.New()
	routes.ProductRoutes(e)
	routes.CategoryRoutes(e)
	recCat, reqCat := doRequest(e, http.MethodPost, "/categories", map[string]interface{}{"name": "Appliances"})
	e.ServeHTTP(recCat, reqCat)
	var cat models.Category
	assert.NoError(t, json.Unmarshal(recCat.Body.Bytes(), &cat))
	catID := cat.ID

	t.Run("POST /products positive", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":        "Refrigerator",
			"category_id": catID,
			"price":       999.99,
		}
		rec, req := doRequest(e, http.MethodPost, "/products", payload)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code, "Tworzenie produktu – pozytywny scenariusz")
		var prod models.Product
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &prod))
		assert.Equal(t, "Refrigerator", prod.Name)
		assert.Equal(t, catID, prod.CategoryID)
		assert.Equal(t, 999.99, prod.Price)
		assert.Equal(t, "Appliances", prod.Category.Name)
	})

	t.Run("POST /products negative", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString("invalid-json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code, "Niepoprawny JSON powinien skutkować 400")
	})

	t.Run("GET /products/:id positive", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":        "Microwave",
			"category_id": catID,
			"price":       200.0,
		}
		recCreate, reqCreate := doRequest(e, http.MethodPost, "/products", payload)
		e.ServeHTTP(recCreate, reqCreate)
		var prod models.Product
		assert.NoError(t, json.Unmarshal(recCreate.Body.Bytes(), &prod))
		idStr := strconv.Itoa(int(prod.ID))
		rec, req := doRequest(e, http.MethodGet, "/products/"+idStr, nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "Pobranie produktu – pozytywny scenariusz")
	})

	t.Run("GET /products/:id negative", func(t *testing.T) {
		rec, req := doRequest(e, http.MethodGet, "/products/9999", nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code, "Pobranie nieistniejącego produktu powinno zwrócić 404")
	})

	t.Run("PUT /products/:id positive", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":        "Dishwasher",
			"category_id": catID,
			"price":       600.0,
		}
		recCreate, reqCreate := doRequest(e, http.MethodPost, "/products", payload)
		e.ServeHTTP(recCreate, reqCreate)
		var prod models.Product
		assert.NoError(t, json.Unmarshal(recCreate.Body.Bytes(), &prod))
		idStr := strconv.Itoa(int(prod.ID))
		updatePayload := map[string]interface{}{
			"name":        "Dishwasher Pro",
			"category_id": catID,
			"price":       750.0,
		}
		rec, req := doRequest(e, http.MethodPut, "/products/"+idStr, updatePayload)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "Aktualizacja produktu – pozytywny scenariusz")
		var updatedProd models.Product
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &updatedProd))
		assert.Equal(t, "Dishwasher Pro", updatedProd.Name)
		assert.Equal(t, 750.0, updatedProd.Price)
	})

	t.Run("PUT /products/:id negative", func(t *testing.T) {
		updatePayload := map[string]interface{}{
			"name":        "Nonexistent Product",
			"category_id": catID,
			"price":       100.0,
		}
		rec, req := doRequest(e, http.MethodPut, "/products/9999", updatePayload)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code, "Aktualizacja nieistniejącego produktu powinna zwrócić 404")
	})

	t.Run("DELETE /products/:id positive", func(t *testing.T) {
		// Utwórz produkt, a następnie usuń
		payload := map[string]interface{}{
			"name":        "Oven",
			"category_id": catID,
			"price":       300.0,
		}
		recCreate, reqCreate := doRequest(e, http.MethodPost, "/products", payload)
		e.ServeHTTP(recCreate, reqCreate)
		var prod models.Product
		assert.NoError(t, json.Unmarshal(recCreate.Body.Bytes(), &prod))
		idStr := strconv.Itoa(int(prod.ID))
		rec, req := doRequest(e, http.MethodDelete, "/products/"+idStr, nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNoContent, rec.Code, "Usunięcie produktu – pozytywny scenariusz")
	})

	t.Run("DELETE /products/:id negative", func(t *testing.T) {
		rec, req := doRequest(e, http.MethodDelete, "/products/9999", nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code, "Usunięcie nieistniejącego produktu powinno zwrócić 404")
	})

	t.Run("GET /products positive (lista)", func(t *testing.T) {
		rec, req := doRequest(e, http.MethodGet, "/products", nil)
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "Pobranie listy produktów – pozytywny scenariusz")
		var prods []models.Product
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &prods))
		assert.GreaterOrEqual(t, len(prods), 0)
	})
}
