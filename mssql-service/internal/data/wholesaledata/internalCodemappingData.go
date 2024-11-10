package wholesaledata

import "database/sql"

type InternalCodemappingData struct {
	ID           int `json:"ID"`
	InternalCode int `json:"InternalCode"`
	OperatorCode int `json:"OperatorCode"`
}

type InternalCodemappingDataModel struct {
	DB *sql.DB
}

func (r InternalCodemappingDataModel) GetAll() ([]InternalCodemappingData, error) {
	query := `SELECT ID, InternalCode, OperatorCode FROM InternalCodeMapping`

	rows, err := r.DB.Query(query)
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	var internalCodemapping InternalCodemappingData
	for rows.Next() {
		err := rows.Scan(&internalCodemapping.ID, &internalCodemapping.InternalCode, &internalCodemapping.OperatorCode)
		if err != nil {
			panic(err.Error())
		}
	}

	err = rows.Err()
	if err != nil {
		panic(err.Error())
	}

	return []InternalCodemappingData{internalCodemapping}, nil
}
