package applicationrepository

import (
	"context"
	"fmt"
	appStruct "marketplace-api/application-struct"
	pgStruct "marketplace-api/postgre-struct"
	utility "marketplace-api/utility"
	"os"
)

func GetApplicationsFromDB() ([]appStruct.Applications, *pgStruct.ErrorResponse) {
	conn, err := utility.GetPostgresConnection()
	if err != nil {
		fmt.Println(err)
		return nil, &pgStruct.ErrorResponse{
			Message:    "Something went wrong, please try after sometime",
			StatusCode: "500",
		}
	}
	defer conn.Close(context.Background())
	if conn != nil {
		var apps []appStruct.Applications
		rows, err := conn.Query(context.Background(), `select "appName",id,description,"imageName" from "applicationDetails"`)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		}
		defer rows.Close()

		for rows.Next() {
			var appName string
			var id int32
			var description string
			var imageName string
			err = rows.Scan(&appName, &id, &description, &imageName)
			if err != nil {
				fmt.Println(err)
			}
			app := appStruct.Applications{
				AppName:     appName,
				Id:          id,
				Description: description,
				ImageName:   imageName,
			}
			apps = append(apps, app)
		}
		return apps, nil
	} else {
		return nil, &pgStruct.ErrorResponse{
			Message:    "Something went wrong, please try after sometime",
			StatusCode: "500",
		}
	}
}

func GetEnvironmentsFromDB() ([]appStruct.Environment, *pgStruct.ErrorResponse) {
	conn, _ := utility.GetPostgresConnection()
	defer conn.Close(context.Background())
	if conn != nil {
		var envs []appStruct.Environment
		rows, err := conn.Query(context.Background(), `select id, name, description from environment_details`)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		}
		defer rows.Close()

		for rows.Next() {
			var name string
			var id int32
			var description string
			err = rows.Scan(&id, &name, &description)
			if err != nil {
				fmt.Println(err)
			}
			env := appStruct.Environment{
				Name:        name,
				Id:          id,
				Description: description,
			}
			envs = append(envs, env)
		}

		return envs, nil
	} else {
		errorResponse := pgStruct.ErrorResponse{
			Message:    "Something went wrong, please try again later",
			StatusCode: "500",
		}
		return nil, &errorResponse
	}
}

func GetApplicationVersionsFromID(id int) (map[string]string, *pgStruct.ErrorResponse) {
	conn, err := utility.GetPostgresConnection()
	if err != nil {
		fmt.Println(err)
		return nil, &pgStruct.ErrorResponse{Message: "Something went wrong, please try again later", StatusCode: "500"}
	}
	defer conn.Close(context.Background())
	if conn != nil {
		var versions map[string]string
		err = conn.QueryRow(context.Background(), `select versions from "applicationDetails" where id = $1`, id).Scan(&versions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		}
		//fmt.Println("\n\n", versions)
		return versions, nil
	} else {
		return nil, &pgStruct.ErrorResponse{
			Message:    "Something went wrong, please try after sometime",
			StatusCode: "500",
		}
	}
}
