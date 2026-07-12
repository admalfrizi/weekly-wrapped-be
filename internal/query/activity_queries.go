package query

const (
	InsertActivity = `
		INSERT INTO activities (user_id, category_id, value, note, occurred_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id, created_at;
	`

	ListActivities = `
		SELECT a.id, a.user_id, a.category_id, a.value, a.note, a.occurred_at, a.created_at, c.id, c.name, c.icon, c.color_hex
		FROM activities a
		JOIN categories c ON a.category_id = c.id
		WHERE a.user_id = $1 
		ORDER BY a.occurred_at DESC 
		LIMIT $2 OFFSET $3;
	`

	CountActivities = `
		SELECT COUNT(id) FROM activities WHERE user_id = $1;
	`

	GetActivityByID = `
		SELECT 
			a.id, a.user_id, a.category_id, a.value, a.note, a.occurred_at, a.created_at,
			c.id, c.name, c.icon, c.color_hex
		FROM activities a
		JOIN categories c ON a.category_id = c.id
		WHERE a.id = $1 AND a.user_id = $2;
	`

	UpdateActivity = `
		UPDATE activities 
		SET category_id = $1, value = $2, note = $3, occurred_at = $4 
		WHERE id = $5 AND user_id = $6 
		RETURNING id;
	`

	DeleteActivity = `
		DELETE FROM activities WHERE id = $1 AND user_id = $2;
	`
)