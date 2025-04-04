package handlers

import (
	"account/internal/database"
	"account/internal/models"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v13"
	"github.com/labstack/echo/v4"
)

// CreateAccount handles POST /accounts to create a new account.
func CreateAccount(c echo.Context) error {
	account := new(models.Account)
	if err := c.Bind(account); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Extract admin password from the request payload.
	adminPassword := account.AdminPassword
	if adminPassword == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Admin password is required"})
	}

	// Insert account into the DB.
	query := `INSERT INTO accounts (accountname, admin_email, admin_phone, config)
              VALUES ($1, $2, $3, $4)
              RETURNING id, created_at`
	err := database.DB.QueryRow(query, account.AccountName, account.AdminEmail, account.AdminPhone, account.Config).
		Scan(&account.ID, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "Account not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Create a Keycloak realm and related resources for this account,
	// passing in the admin email and password from the request.
	if err := createKeycloakRealm(account.AccountName, account.AdminEmail, adminPassword); err != nil {
		log.Printf("Account created but failed to create Keycloak realm and resources: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Account created but failed to create Keycloak realm and resources: " + err.Error(),
		})
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

// createKeycloakRealm creates a new realm, an admin user with the realm-admin role,
// and a public client in Keycloak.
// It now takes adminEmail and adminPassword from the request payload.
func createKeycloakRealm(realmName, adminEmail, adminPassword string) error {
	// Use environment variable for Keycloak base URL (e.g., "http://localhost:8090")
	keycloakURL := os.Getenv("KEYCLOAK_BASE_URL")
	client := gocloak.NewClient(keycloakURL)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Login to Keycloak as an admin.
	token, err := client.LoginAdmin(ctx, os.Getenv("KEYCLOAK_USERNAME"), os.Getenv("KEYCLOAK_PASSWORD"), os.Getenv("KEYCLOAK_REALM"))
	if err != nil {
		log.Printf("Failed to login to Keycloak: %v", err)
		return err
	}

	// Create a new realm.
	realmRep := gocloak.RealmRepresentation{
		Realm:   &realmName,
		Enabled: gocloak.BoolP(true),
	}
	if _, err := client.CreateRealm(ctx, token.AccessToken, realmRep); err != nil {
		log.Printf("Failed to create realm %s: %v", realmName, err)
		return err
	}
	log.Printf("Successfully created realm %s", realmName)

	// Create an admin user for the new realm using the provided adminEmail.
	adminUser := gocloak.User{
		Username: gocloak.StringP("admin"),
		Email:    gocloak.StringP(adminEmail),
		Enabled:  gocloak.BoolP(true),
	}
	userID, err := client.CreateUser(ctx, token.AccessToken, realmName, adminUser)
	if err != nil {
		log.Printf("Failed to create admin user for realm %s: %v", realmName, err)
		return err
	}
	log.Printf("Admin user created with ID: %s", userID)

	// Set password for the admin user using the provided adminPassword.
	if err := client.SetPassword(ctx, token.AccessToken, userID, realmName, adminPassword, false); err != nil {
		log.Printf("Failed to set password for admin user: %v", err)
		return err
	}

	// Retrieve the "realm-admin" role from the realm.
	realmAdminRole, err := client.GetRealmRole(ctx, token.AccessToken, realmName, "realm-admin")
	if err != nil {
		// If the role is not found, create it.
		if strings.Contains(err.Error(), "404") {
			newRole := gocloak.Role{
				Name: gocloak.StringP("realm-admin"),
			}
			if _, err := client.CreateRealmRole(ctx, token.AccessToken, realmName, newRole); err != nil {
				log.Printf("Failed to create realm-admin role: %v", err)
				return err
			}
			// Retrieve the newly created role.
			realmAdminRole, err = client.GetRealmRole(ctx, token.AccessToken, realmName, "realm-admin")
			if err != nil {
				log.Printf("Failed to get realm-admin role after creating it for realm %s: %v", realmName, err)
				return err
			}
		} else {
			log.Printf("Failed to get realm-admin role for realm %s: %v", realmName, err)
			return err
		}
	}

	// Assign the "realm-admin" role to the admin user.
	err = client.AddRealmRoleToUser(ctx, token.AccessToken, realmName, userID, []gocloak.Role{*realmAdminRole})
	if err != nil {
		log.Printf("Failed to assign realm-admin role to admin user: %v", err)
		return err
	}
	log.Printf("Realm-admin role assigned to admin user")

	// Create a public client in the new realm.
	publicClient := gocloak.Client{
		ClientID:                  gocloak.StringP("public-client"),
		Name:                      gocloak.StringP("Public Client"),
		PublicClient:              gocloak.BoolP(true),
		RedirectURIs:              &[]string{"*"},
		StandardFlowEnabled:       gocloak.BoolP(true),
		ImplicitFlowEnabled:       gocloak.BoolP(true),
		DirectAccessGrantsEnabled: gocloak.BoolP(true),
	}
	clientID, err := client.CreateClient(ctx, token.AccessToken, realmName, publicClient)
	if err != nil {
		log.Printf("Failed to create public client in realm %s: %v", realmName, err)
		return err
	}
	log.Printf("Public client created with ID: %s", clientID)

	return nil
}
