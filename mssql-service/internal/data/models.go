package data

import "database/sql"

type Models struct {
	CallLogs interface {
		Insert(callLog *CallLog) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		CallLogs: CallLogModel{DB: db},
	}
}
