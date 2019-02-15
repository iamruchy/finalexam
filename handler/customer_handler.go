package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iamruchy/finalexam/database"
)

type Customer struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func CreateCustomerHandler(c *gin.Context) {
	cust := Customer{}
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}
	id, err := database.InsertCustomer(cust.Name, cust.Email, cust.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	row, err := database.GetCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cust)

}

func GetCustomerHandler(c *gin.Context) {
	custs := []Customer{}
	rows, err := database.GetCustomer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	for rows.Next() {
		cust := Customer{}
		if err := rows.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
			return
		}
		custs = append(custs, cust)
	}
	c.JSON(http.StatusOK, custs)
}

func GetCustomerByIDHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Id must be number"})
	}

	row, err := database.GetCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	cust := Customer{}
	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		c.JSON(http.StatusNoContent, cust)
		return
	}

	c.JSON(http.StatusOK, cust)
}

func UpdateCustomerHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Id must be number"})
	}

	cust := Customer{}
	if err := c.ShouldBindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}
	if err = database.UpdateCustomer(id, cust.Name, cust.Email, cust.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	row, err := database.GetCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cust)

}

func DeleteCustomerHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Id must be number"})
	}
	err = database.DeleteCustomer(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}

func loginMiddleware(c *gin.Context) {
	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Unauthorized",
		})
		c.Abort()
		return
	}
	c.Next()
}

func NewRouter() *gin.Engine {
	database.CreateTable()
	r := gin.Default()
	r.Use(loginMiddleware)
	r.GET("/customers/:id", GetCustomerByIDHandler)
	r.GET("/customers", GetCustomerHandler)
	r.POST("/customers", CreateCustomerHandler)
	r.PUT("/customers/:id", UpdateCustomerHandler)
	r.DELETE("/customers/:id", DeleteCustomerHandler)

	return r
}
