package query

const (
	InsertUser = `
		INSERT INTO users (email, username, name, password_hash) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at;
	`

	GetUserByEmail = `
		SELECT id, email, username, name, profile_img_url, password_hash, created_at, updated_at 
		FROM users 
		WHERE email = $1;
	`

	GetUserByUsername = `
		SELECT id, email, username, name, profile_img_url, password_hash, created_at, updated_at 
		FROM users 
		WHERE username = $1;
	`
)