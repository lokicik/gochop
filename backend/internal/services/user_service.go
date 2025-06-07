package services

import (
	"context"
	"gochop/backend/internal/db"
	"time"
)

// User represents a user in the system
type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Image    string    `json:"image,omitempty"`
	IsAdmin  bool      `json:"is_admin"`
	Created  time.Time `json:"created_at"`
	Updated  time.Time `json:"updated_at"`
}

// CreateUserInput represents the input for creating a new user
type CreateUserInput struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image,omitempty"`
}

// UpdateUserInput represents the input for updating a user
type UpdateUserInput struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Image string `json:"image,omitempty"`
}

// UserService handles user-related operations
type UserService struct{}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{}
}

// CreateUser creates a new user in the database
func (s *UserService) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	query := `
		INSERT INTO users (id, name, email, image, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			image = EXCLUDED.image,
			updated_at = NOW()
		RETURNING id, name, email, image, created_at, updated_at
	`

	var user User
	err := db.DB.QueryRow(ctx, query, input.ID, input.Name, input.Email, input.Image).Scan(
		&user.ID, &user.Name, &user.Email, &user.Image, &user.Created, &user.Updated,
	)

	if err != nil {
		return nil, err
	}

	// Check if user is admin (for simplicity, first user or specific email)
	user.IsAdmin = s.isUserAdmin(user.Email)

	return &user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
	query := `
		SELECT id, name, email, COALESCE(image, '') as image, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	var user User
	err := db.DB.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.Name, &user.Email, &user.Image, &user.Created, &user.Updated,
	)

	if err != nil {
		return nil, err
	}

	// Check if user is admin
	user.IsAdmin = s.isUserAdmin(user.Email)

	return &user, nil
}

// UpdateUser updates a user's information
func (s *UserService) UpdateUser(ctx context.Context, userID string, input UpdateUserInput) (*User, error) {
	query := `
		UPDATE users 
		SET name = COALESCE(NULLIF($2, ''), name),
			email = COALESCE(NULLIF($3, ''), email),
			image = COALESCE(NULLIF($4, ''), image),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, email, COALESCE(image, '') as image, created_at, updated_at
	`

	var user User
	err := db.DB.QueryRow(ctx, query, userID, input.Name, input.Email, input.Image).Scan(
		&user.ID, &user.Name, &user.Email, &user.Image, &user.Created, &user.Updated,
	)

	if err != nil {
		return nil, err
	}

	// Check if user is admin
	user.IsAdmin = s.isUserAdmin(user.Email)

	return &user, nil
}

// GetUserStats retrieves statistics for a user
func (s *UserService) GetUserStats(ctx context.Context, userID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Get total links created by user
	var totalLinks int
	linkQuery := "SELECT COUNT(*) FROM links WHERE user_id = $1"
	err := db.DB.QueryRow(ctx, linkQuery, userID).Scan(&totalLinks)
	if err != nil {
		totalLinks = 0
	}
	stats["total_links"] = totalLinks

	// Get total clicks for user's links
	var totalClicks int
	clickQuery := `
		SELECT COALESCE(SUM(click_count), 0) 
		FROM (
			SELECT COUNT(*) as click_count 
			FROM analytics a 
			JOIN links l ON a.short_code = l.short_code 
			WHERE l.user_id = $1
		) as user_clicks
	`
	err = db.DB.QueryRow(ctx, clickQuery, userID).Scan(&totalClicks)
	if err != nil {
		totalClicks = 0
	}
	stats["total_clicks"] = totalClicks

	// Get active links (non-expired)
	var activeLinks int
	activeQuery := "SELECT COUNT(*) FROM links WHERE user_id = $1 AND expires_at > NOW()"
	err = db.DB.QueryRow(ctx, activeQuery, userID).Scan(&activeLinks)
	if err != nil {
		activeLinks = 0
	}
	stats["active_links"] = activeLinks

	// Get most clicked link
	var mostClickedLink string
	var mostClicks int
	mostClickedQuery := `
		SELECT l.short_code, COUNT(a.id) as clicks
		FROM links l
		LEFT JOIN analytics a ON l.short_code = a.short_code
		WHERE l.user_id = $1
		GROUP BY l.short_code
		ORDER BY clicks DESC
		LIMIT 1
	`
	err = db.DB.QueryRow(ctx, mostClickedQuery, userID).Scan(&mostClickedLink, &mostClicks)
	if err == nil && mostClickedLink != "" {
		stats["most_clicked_link"] = map[string]interface{}{
			"short_code": mostClickedLink,
			"clicks":     mostClicks,
		}
	}

	return stats, nil
}

// isUserAdmin checks if a user should have admin privileges
func (s *UserService) isUserAdmin(email string) bool {
	// For now, simple admin check - you can enhance this logic
	// Could check against environment variable, database table, etc.
	adminEmails := []string{
		"admin@gochop.io",
		"lokman@example.com", // Add your email here
	}

	for _, adminEmail := range adminEmails {
		if email == adminEmail {
			return true
		}
	}

	return false
}

// ListUsers returns all users (admin only)
func (s *UserService) ListUsers(ctx context.Context) ([]User, error) {
	query := `
		SELECT id, name, email, COALESCE(image, '') as image, created_at, updated_at
		FROM users 
		ORDER BY created_at DESC
	`

	rows, err := db.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Image, &user.Created, &user.Updated)
		if err != nil {
			continue
		}
		
		// Check if user is admin
		user.IsAdmin = s.isUserAdmin(user.Email)
		users = append(users, user)
	}

	return users, nil
} 