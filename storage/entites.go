package storage

import "database/sql"

type AnimalResponse struct {
	AnimalId sql.NullInt32  `db:"animal_id"`
	Animal   sql.NullString `db:"animal"`
	Price    sql.NullInt32  `db:"price"`
}
