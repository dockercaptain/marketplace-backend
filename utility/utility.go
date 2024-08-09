package utility

import (
	"context"
	"errors"
	"fmt"
	pgStruct "marketplace-api/postgre-struct"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetPostgresConnection() (*pgx.Conn, error) {
	DATABASE_URL := "postgres://postgres:root@localhost:5432/marketplace"
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, errors.New("unable to connect to database")
	}
	//defer conn.Close(context.Background())
	return conn, nil
}

// enable cors
func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func BuildResponseForVersions(versionsMap map[string]string) []string {
	var appVersions []string
	for _, j := range versionsMap {
		appVersions = append(appVersions, j)
	}
	return appVersions
}

func GetErrorResponse(err error, statusCode string) *pgStruct.ErrorResponse {
	return &pgStruct.ErrorResponse{
		Message:    err.Error(),
		StatusCode: statusCode,
	}
}
