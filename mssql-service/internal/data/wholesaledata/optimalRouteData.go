package wholesaledata

import (
	"context"
	"database/sql"
	"time"
)

type OptimalRouteData struct {
	Class1 int `json:"Class1"`
	Class2 int `json:"Class2"`
	Class3 int `json:"Class3"`
}

type OptimalRouteDataModel struct {
	DB *sql.DB
}

func (r OptimalRouteDataModel) ExecuteGetOptimalRoute(pCallString string) (OptimalRouteData, error) {
	query := `EXEC GetOptimalRoute 
				@pCallString = $1;`

	args := []interface{}{
		pCallString,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := r.DB.QueryRowContext(ctx, query, args...)

	var result OptimalRouteData
	err := row.Scan(&result.Class1, &result.Class2, &result.Class3)
	if err != nil {
		return OptimalRouteData{}, err
	}

	return result, nil
}
