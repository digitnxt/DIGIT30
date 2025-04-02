package controller

import (
	"net/http"
	"strconv"

	"account/config"
	"account/model"

	"github.com/labstack/echo/v4"
)

// CreateAccount creates a new account.
func CreateAccount(c echo.Context) error {
	acc := new(model.Account)
	db := config.DB()

	// Bind the incoming JSON payload to the Account struct.
	if err := c.Bind(acc); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	// Create the account record.
	account := &model.Account{
		AccountName: acc.AccountName,
		AdminEmail:  acc.AdminEmail,
		AdminPhone:  acc.AdminPhone,
		Config:      acc.Config,
		// CreatedAt is set automatically if using GORM's autoCreateTime tag.
	}

	if err := db.Create(&account).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	response := map[string]interface{}{
		"data": account,
	}

	return c.JSON(http.StatusOK, response)
}

// GetAccount retrieves a single account by its ID.
func GetAccount(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid account ID",
		})
	}

	db := config.DB()
	var account model.Account
	if err := db.First(&account, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Account not found",
		})
	}

	return c.JSON(http.StatusOK, account)
}

// ListAccounts retrieves all accounts.
func ListAccounts(c echo.Context) error {
	db := config.DB()
	var accounts []model.Account
	if err := db.Find(&accounts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, accounts)
}

// UpdateAccount updates an existing account by its ID.
func UpdateAccount(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid account ID",
		})
	}

	db := config.DB()
	var account model.Account
	if err := db.First(&account, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": "Account not found",
		})
	}

	// Bind the new data into the account struct.
	if err := c.Bind(&account); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
		})
	}

	if err := db.Save(&account).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, account)
}

// DeleteAccount removes an account by its ID.
func DeleteAccount(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid account ID",
		})
	}

	db := config.DB()
	if err := db.Delete(&model.Account{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Account deleted successfully",
	})
}
