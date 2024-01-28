package data

import "database/sql"

type Models struct {
	Services ServiceModel
	Users    UserModel
	Colleges CollegeModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Services: ServiceModel{DB: db},
		Users:    UserModel{DB: db},
		Colleges: CollegeModel{DB: db},
	}
}
