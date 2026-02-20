package repo

const (
	getAvailableClassesQuery = `
        SELECT DISTINCT class
        FROM contents
        WHERE is_active = TRUE
        ORDER BY class;
    `

	getQuartersByClassQuery = `
        SELECT DISTINCT quarter
        FROM contents
        WHERE class = $1
          AND is_active = TRUE
        ORDER BY quarter;
    `

	getByClassAndQuarterQuery = `
        SELECT id, title, canva_url, class, quarter, lesson_number
        FROM contents
        WHERE class = $1
          AND quarter = $2
          AND is_active = TRUE
        ORDER BY lesson_number;
    `

	getByClassQuery = `
        SELECT id, title, canva_url, class, quarter, lesson_number
        FROM contents
        WHERE class = $1
          AND is_active = TRUE
        ORDER BY lesson_number;
    `

	getContentByIDQuery = `
        SELECT id, title, canva_url, class, quarter, lesson_number
        FROM contents
        WHERE id = $1
          AND is_active = TRUE;
    `

	addContentQuery = `
        INSERT INTO contents (title, canva_url, class, quarter, lesson_number)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id;
    `

	deleteContentQuery = `
        UPDATE contents
        SET is_active = FALSE,
            updated_at = NOW()
        WHERE id = $1;
    `
)
