package util

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Nerzal/gocloak/v13"
)

// LoginToKeycloak logs in as admin and returns the Keycloak client and token.
func LoginToKeycloak(ctx context.Context) (*gocloak.GoCloak, *gocloak.JWT, error) {
	keycloakURL := os.Getenv("KEYCLOAK_BASE_URL")
	// Do not take the address of the value returned by NewClient.
	client := gocloak.NewClient(keycloakURL)

	token, err := client.LoginAdmin(ctx,
		os.Getenv("KEYCLOAK_USERNAME"),
		os.Getenv("KEYCLOAK_PASSWORD"),
		os.Getenv("KEYCLOAK_REALM"),
	)
	if err != nil {
		// Cast nil to gocloak.GoCloak to satisfy the return type.
		return nil, nil, err
	}
	return client, token, nil
}

// CreateRealm creates a new realm with the given name.
func CreateRealm(ctx context.Context, client gocloak.GoCloak, token *gocloak.JWT, realmName string) error {
	realmRep := gocloak.RealmRepresentation{
		Realm:   &realmName,
		Enabled: gocloak.BoolP(true),
	}
	if _, err := client.CreateRealm(ctx, token.AccessToken, realmRep); err != nil {
		log.Printf("Failed to create realm %s: %v", realmName, err)
		return err
	}
	return nil
}

// CreateAdminUser creates an admin user in the specified realm using the provided adminEmail.
func CreateAdminUser(ctx context.Context, client gocloak.GoCloak, token *gocloak.JWT, realmName, adminEmail string) (string, error) {
	adminUser := gocloak.User{
		Username: gocloak.StringP("admin"),
		Email:    gocloak.StringP(adminEmail),
		Enabled:  gocloak.BoolP(true),
	}
	userID, err := client.CreateUser(ctx, token.AccessToken, realmName, adminUser)
	if err != nil {
		log.Printf("Failed to create admin user for realm %s: %v", realmName, err)
		return "", err
	}
	return userID, nil
}

// SetAdminPassword sets the password for the given admin user.
func SetAdminPassword(ctx context.Context, client gocloak.GoCloak, token *gocloak.JWT, userID, realmName string) error {
	// In production, use a secure password or load it from environment variables.
	adminPassword := "admin123"
	if err := client.SetPassword(ctx, token.AccessToken, userID, realmName, adminPassword, false); err != nil {
		log.Printf("Failed to set password for admin user: %v", err)
		return err
	}
	return nil
}

// GetOrCreateRealmAdminRole retrieves the "realm-admin" role, or creates it if not found.
func GetOrCreateRealmAdminRole(ctx context.Context, client gocloak.GoCloak, token *gocloak.JWT, realmName string) (*gocloak.Role, error) {
	role, err := client.GetRealmRole(ctx, token.AccessToken, realmName, "realm-admin")
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			newRole := gocloak.Role{
				Name: gocloak.StringP("realm-admin"),
			}
			if _, err := client.CreateRealmRole(ctx, token.AccessToken, realmName, newRole); err != nil {
				log.Printf("Failed to create realm-admin role: %v", err)
				return nil, err
			}
			role, err = client.GetRealmRole(ctx, token.AccessToken, realmName, "realm-admin")
			if err != nil {
				log.Printf("Failed to get realm-admin role after creating it for realm %s: %v", realmName, err)
				return nil, err
			}
		} else {
			log.Printf("Failed to get realm-admin role for realm %s: %v", realmName, err)
			return nil, err
		}
	}
	return role, nil
}

// AssignRealmAdminRole assigns the specified role to the admin user.
func AssignRealmAdminRole(ctx context.Context, client gocloak.GoCloak, token *gocloak.JWT, realmName, userID string, role *gocloak.Role) error {
	if err := client.AddRealmRoleToUser(ctx, token.AccessToken, realmName, userID, []gocloak.Role{*role}); err != nil {
		log.Printf("Failed to assign realm-admin role to admin user: %v", err)
		return err
	}
	return nil
}

// CreatePublicClient creates a public client in the specified realm.
func CreatePublicClient(ctx context.Context, client gocloak.GoCloak, token *gocloak.JWT, realmName string) (string, error) {
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
		return "", err
	}
	return clientID, nil
}

// CreateKeycloakRealm orchestrates the creation of a new realm, admin user, role assignment, and public client.
func CreateKeycloakRealm(realmName, adminEmail string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, token, err := LoginToKeycloak(ctx)
	if err != nil {
		return err
	}

	// Create the realm.
	if err := CreateRealm(ctx, client, token, realmName); err != nil {
		return err
	}
	log.Printf("Successfully created realm %s", realmName)

	// Create the admin user.
	userID, err := CreateAdminUser(ctx, client, token, realmName, adminEmail)
	if err != nil {
		return err
	}
	log.Printf("Admin user created with ID: %s", userID)

	// Set the admin user's password.
	if err := SetAdminPassword(ctx, client, token, userID, realmName); err != nil {
		return err
	}

	// Get (or create) the realm-admin role.
	role, err := GetOrCreateRealmAdminRole(ctx, client, token, realmName)
	if err != nil {
		return err
	}

	// Assign the realm-admin role to the admin user.
	if err := AssignRealmAdminRole(ctx, client, token, realmName, userID, role); err != nil {
		return err
	}
	log.Printf("Realm-admin role assigned to admin user")

	// Create the public client.
	clientID, err := CreatePublicClient(ctx, client, token, realmName)
	if err != nil {
		return err
	}
	log.Printf("Public client created with ID: %s", clientID)

	return nil
}
