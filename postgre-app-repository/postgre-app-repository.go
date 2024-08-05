package postgreapprepository

import (
	"context"
	"fmt"
	pgStruct "marketplace-api/postgre-struct"
	utility "marketplace-api/utility"
	"os"
)

func CreateApplicationPostgres(pgApp pgStruct.PostgresApp) (*pgStruct.SuccessResponse, *pgStruct.ErrorResponse) {
	conn, _ := utility.GetPostgresConnection()
	tx, err := conn.Begin(context.Background())
	errResponse := &pgStruct.ErrorResponse{Message: "Something went wrong, please try again later", StatusCode: "500"}
	if err != nil {
		fmt.Println(err)
		return nil, errResponse
	}
	defer tx.Rollback(context.Background())

	insertQuery := `INSERT INTO public.installed_postgres_details(
			status, description, "serverName", "adminUser", password, version, environment, "sizeDisk", "storageType", "sizeCPU", "sizeMemory", "issues")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`
	_, err = tx.Exec(context.Background(), insertQuery, pgApp.Status, pgApp.Description, pgApp.ServerName, pgApp.AdminUser, pgApp.Password, pgApp.Version, pgApp.Environment, pgApp.SizeDisk, pgApp.StorageType, pgApp.SizeCPU, pgApp.SizeMemory, pgApp.Issues)
	if err != nil {
		fmt.Println(err)
		return nil, errResponse
	}
	err = tx.Commit(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, errResponse
	}
	return &pgStruct.SuccessResponse{Message: "Data saved successfully", StatusCode: "201", Status: "SUCCESS"}, nil

}

func GetApplicationPostgresById(pgId int) (*pgStruct.PostgresApp, *pgStruct.ErrorResponse) {
	conn, _ := utility.GetPostgresConnection()
	errResponse := &pgStruct.ErrorResponse{Message: "Something went wrong, please try again later", StatusCode: "500"}
	defer conn.Close(context.Background())
	var pgApp *pgStruct.PostgresApp
	selectQuery := `SELECT * FROM public.installed_postgres_details where id=$1;`
	var id int32
	var status string
	var description string
	var serverName string
	var adminUser string
	var password string
	var version string
	var environment string
	var sizeDisk string
	var storageType string
	var sizeCPU string
	var sizeMemory string
	var issues string
	err := conn.QueryRow(context.Background(), selectQuery, pgId).Scan(&id, &status, &description, &serverName, &adminUser, &password, &version, &environment, &sizeDisk, &storageType, &sizeCPU, &sizeMemory, &issues)
	if err != nil {
		errResponse.StatusCode = "404"
		errResponse.Message = `Postgres app doesn't exist for id`
		return nil, errResponse
	}
	pgApp = &pgStruct.PostgresApp{
		Id:          id,
		Status:      status,
		Description: description,
		ServerName:  serverName,
		AdminUser:   adminUser,
		Password:    password,
		Version:     version,
		Environment: environment,
		SizeDisk:    sizeDisk,
		StorageType: storageType,
		SizeCPU:     sizeCPU,
		SizeMemory:  sizeMemory,
		Issues: issues,
	}
	return pgApp, nil
}

func GetAllPostgresFromDB() ([]pgStruct.PostgresApp, *pgStruct.ErrorResponse) {
	conn, err := utility.GetPostgresConnection()
	if err != nil {
		return nil, &pgStruct.ErrorResponse{
			Message:    "Something went wrong, please try after sometime",
			StatusCode: "500",
		}
	}
	defer conn.Close(context.Background())
	if conn != nil {
		var postgresList []pgStruct.PostgresApp
		rows, err := conn.Query(context.Background(), `select * from installed_postgres_details`)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		}
		defer rows.Close()

		for rows.Next() {
			var id int32
			var status string
			var description string
			var serverName string
			var adminUser string
			var password string
			var version string
			var environment string
			var sizeDisk string
			var storageType string
			var sizeCPU string
			var sizeMemory string
			var issues string
			err = rows.Scan(&id, &status, &description, &serverName, &adminUser, &password, &version, &environment, &sizeDisk, &storageType, &sizeCPU, &sizeMemory, &issues)
			if err != nil {
				fmt.Println(err)
			}
			postgres := pgStruct.PostgresApp{
				Id:          id,
				ServerName:  serverName,
				Description: description,
				Status:      status,
				AdminUser:   adminUser,
				Password:    password,
				Version:     version,
				Environment: environment,
				SizeDisk:    sizeDisk,
				StorageType: storageType,
				SizeCPU:     sizeCPU,
				SizeMemory:  sizeMemory,
				Issues: issues,
			}
			postgresList = append(postgresList, postgres)
		}
		return postgresList, nil
	} else {
		return nil, &pgStruct.ErrorResponse{
			Message:    "Something went wrong, please try after sometime",
			StatusCode: "500",
		}
	}
}
