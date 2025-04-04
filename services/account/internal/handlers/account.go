package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"account/internal/database"
	"account/internal/models"

	"github.com/labstack/echo/v4"
)

// CreateAccount handles POST /accounts to create a new account.
func CreateAccount(c echo.Context) error {
	account := new(models.Account)
	if err := c.Bind(account); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	query := `INSERT INTO accounts (accountname, admin_email, admin_phone, config)
              VALUES ($1, $2, $3, $4)
              RETURNING id, created_at`
	err := database.DB.QueryRow(query, account.AccountName, account.AdminEmail, account.AdminPhone, account.Config).
		Scan(&account.ID, &account.CreatedAt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, account)
}

// GetAccount handles GET /accounts/:id to fetch a single account.
func GetAccount(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid account ID"})
	}

	account := new(models.Account)
	query := `SELECT id, accountname, admin_email, admin_phone, config, created_at FROM accounts WHERE id = $1`
	err = database.DB.QueryRow(query, id).Scan(
		&account.ID,
		&account.AccountName,
		&account.AdminEmail,
		&account.AdminPhone,
		&account.Config,
		&account.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Account not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, account)
}

// ListAccounts handles GET /accounts to list all accounts.
func ListAccounts(c echo.Context) error {
	rows, err := database.DB.Query(`SELECT id, accountname, admin_email, admin_phone, config, created_at FROM accounts`)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		if err := rows.Scan(&account.ID, &account.AccountName, &account.AdminEmail, &account.AdminPhone, &account.Config, &account.CreatedAt); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		accounts = append(accounts, account)
	}
	return c.JSON(http.StatusOK, accounts)
}

// UpdateAccount handles PUT /accounts/:id to update an account.
func UpdateAccount(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid account ID"})
	}

	account := new(models.Account)
	if err := c.Bind(account); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	query := `UPDATE accounts
              SET accountname = $1, admin_email = $2, admin_phone = $3, config = $4
              WHERE id = $5`
	res, err := database.DB.Exec(query, account.AccountName, account.AdminEmail, account.AdminPhone, account.Config, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	count, err := res.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if count == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Account not found"})
	}
	return GetAccount(c)
}

// DeleteAccount handles DELETE /accounts/:id to remove an account.
func DeleteAccount(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid account ID"})
	}

	query := `DELETE FROM accounts WHERE id = $1`
	res, err := database.DB.Exec(query, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	count, err := res.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if count == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Account not found"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Account deleted"})
}
