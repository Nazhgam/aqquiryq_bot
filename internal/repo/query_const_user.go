package repo

const (
	getUserByIDQuery = `
        SELECT id, username, is_admin, created_at
        FROM users
        WHERE id = $1;
    `

	addUserQuery = `
        INSERT INTO users (id, username, is_admin)
        VALUES ($1, $2, FALSE)
        ON CONFLICT (id) DO UPDATE
        SET username = EXCLUDED.username,
            is_admin = FALSE;
    `

	addAdminQuery = `
        INSERT INTO users (id, username, is_admin)
        VALUES ($1, $2, TRUE)
        ON CONFLICT (id) DO UPDATE
        SET username = EXCLUDED.username,
            is_admin = TRUE;
    `

	removeUserQuery = `DELETE FROM users WHERE id = $1;`

	isAdminQuery = `
        SELECT is_admin
        FROM users
        WHERE id = $1;
    `
)
