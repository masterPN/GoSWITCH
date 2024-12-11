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
	const query = `SELECT ID, InternalCode, OperatorCode FROM InternalCodeMapping`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var internalCodemappings []InternalCodemappingData

	for rows.Next() {
		var internalCodemapping InternalCodemappingData
		if err := rows.Scan(&internalCodemapping.ID, &internalCodemapping.InternalCode, &internalCodemapping.OperatorCode); err != nil {
			return nil, err
		}
		internalCodemappings = append(internalCodemappings, internalCodemapping)
	}

	return internalCodemappings, nil
}

func (r InternalCodemappingDataModel) Set(input InternalCodemappingData) (InternalCodemappingData, error) {
	query := `
        SELECT ID FROM InternalCodeMapping WHERE InternalCode = $1
    `
	var id int
	err := r.DB.QueryRow(query, input.InternalCode).Scan(&id)
	if err == sql.ErrNoRows {
		// Insert new record
		query = `
            INSERT INTO InternalCodeMapping (InternalCode, OperatorCode)
			OUTPUT INSERTED.ID
            VALUES ($1, $2)
        `
		err = r.DB.QueryRow(query, input.InternalCode, input.OperatorCode).Scan(&id)
		if err != nil {
			return InternalCodemappingData{}, err
		}
	} else if err != nil {
		return InternalCodemappingData{}, err
	} else {
		// Update existing record
		query = `
            UPDATE InternalCodeMapping
            SET InternalCode = $1, OperatorCode = $2
			OUTPUT INSERTED.ID
            WHERE ID = $3
        `
		err = r.DB.QueryRow(query, input.InternalCode, input.OperatorCode, id).Scan(&id)
		if err != nil {
			return InternalCodemappingData{}, err
		}
	}

	input.ID = id
	return input, nil
}

func (r InternalCodemappingDataModel) Delete(internalCode int) error {
	query := `DELETE FROM InternalCodeMapping WHERE InternalCode = $1`
	_, err := r.DB.Exec(query, internalCode)
	return err
}
