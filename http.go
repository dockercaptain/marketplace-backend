package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type Applications struct {
	Id          int32
	AppName     string
	Description string
	ImageName   string
}

type AppVersionsResponse struct {
	Status      string
	AppVersions []string
}

type Environment struct {
	Id          int32
	Name        string
	Description string
}

type PostgresApp struct {
	Id          int32
	Status      string
	Description string
	ServerName  string
	AdminUser   string
	Password    string
	Version     string
	Environment string
	SizeDisk    string
	StorageType string
	SizeCPU     string
	SizeMemory  string
}

type ErrorResponse struct {
	Message    string
	StatusCode string
}

type SuccessResponse struct {
	Message    string
	StatusCode string
	Status     string
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /marketplace/apps", GetApplications)
	mux.HandleFunc("GET /marketplace/apps/versions/{id}", GetApplicationVersions)
	mux.HandleFunc("GET /marketplace/apps/envs", GetEnvironments)
	mux.HandleFunc("GET /marketplace/apps/{appName}", GetAllAppsDetailsByName)
	mux.HandleFunc("POST /marketplace/apps/{appName}/create", CreateAppBasics)
	http.ListenAndServe(":8080", mux)
	fmt.Println("localhost:8080...")
}

// enable cors
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// 1.
func GetApplications(w http.ResponseWriter, r *http.Request) {
	// enable cors
	enableCors(&w)
	// Dynamically access the path variable
	//fmt.Fprintf(w, "Retrieving item with ID: %s", item)
	applications, err := GetApplicationsFromDB()
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(applications)
	}

}

// 1.
func GetApplicationVersions(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	id := r.PathValue("id")
	fmt.Println(id)
	i, _ := strconv.Atoi(id)
	versions, err := GetApplicationVersionsFromID(i)
	appVersions := BuildResponseForVersions(versions)
	res := AppVersionsResponse{
		Status:      "Success",
		AppVersions: appVersions,
	}
	if err != nil {
		res = AppVersionsResponse{
			Status:      "Failure",
			AppVersions: nil,
		}
	}
	json.NewEncoder(w).Encode(res)

}

func GetEnvironments(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// Dynamically access the path variable
	//fmt.Fprintf(w, "Retrieving item with ID: %s", item)
	environments, err := GetEnvironmentsFromDB()
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(environments)
	}
}

func GetAllAppsDetailsByName(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	appName := r.PathValue("appName")
	var postgresList []PostgresApp
	var err *ErrorResponse
	if appName == "postgres" {
		postgresList, err = GetAllPostgresFromDB()
	}
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(postgresList)
	}
}

func CreateAppBasics(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	appName := r.PathValue("appName")
	if appName == "postgres" {
		var pgApp PostgresApp
		// convert json to golang struct and map to struct reference
		err := json.NewDecoder(r.Body).Decode(&pgApp)
		if err != nil {
			errorResponse := &ErrorResponse{
				Message:    "Please pass valid input",
				StatusCode: "400",
			}
			json.NewEncoder(w).Encode(errorResponse)
		} else {
			// save the record and return response
			success, err := CreateApplicationPostgres(pgApp)
			if err != nil {
				json.NewEncoder(w).Encode(err)
			} else {
				json.NewEncoder(w).Encode(success)
			}
		}
	}
}
func GetApplicationVersionsFromID(id int) (map[string]string, *ErrorResponse) {
	conn, err := GetPostgresConnection()
	if err != nil {
		fmt.Println(err)
		return nil, &ErrorResponse{Message: "Something went wrong, please try again later", StatusCode: "500"}
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
		return nil, &ErrorResponse{
			Message:    "Something went wrong, please try after sometime",
			StatusCode: "500",
		}
	}
}
func CreateApplicationPostgres(pgApp PostgresApp) (*SuccessResponse, *ErrorResponse) {
	conn, err := GetPostgresConnection()
	if err != nil {
		fmt.Println(err)
		return nil, &ErrorResponse{Message: "Something went wrong, please try again later", StatusCode: "500"}
	}
	defer conn.Close(context.Background())
	if conn != nil {
		insertQuery := `INSERT INTO public.installed_postgres_details(
			status, description, "serverName", "adminUser", password, version, environment, "sizeDisk", "storageType", "sizeCPU", "sizeMemory")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`
		conn.Query(context.Background(), insertQuery, pgApp.Status, pgApp.Description, pgApp.ServerName, pgApp.AdminUser, pgApp.Password, pgApp.Version, pgApp.Environment, pgApp.SizeDisk, pgApp.StorageType, pgApp.SizeCPU, pgApp.SizeMemory)
		return &SuccessResponse{Message: "Data saved successfully", StatusCode: "201", Status: "SUCCESS"}, nil
	} else {
		return nil, &ErrorResponse{
			Message:    "Something went wrong, please try after sometime",
			StatusCode: "500",
		}
	}
}

func BuildResponseForVersions(versionsMap map[string]string) []string {
	var appVersions []string
	for _, j := range versionsMap {
		appVersions = append(appVersions, j)
	}
	return appVersions
}

// 2.
func GetApplicationsFromDB() ([]Applications, *ErrorResponse) {
	conn, err := GetPostgresConnection()
	if err != nil {
		fmt.Println(err)
		return nil, &ErrorResponse{
			Message:    "Something went wrong, please try after sometime",
			StatusCode: "500",
		}
	}
	defer conn.Close(context.Background())
	if conn != nil {
		var apps []Applications
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
			app := Applications{
				AppName:     appName,
				Id:          id,
				Description: description,
				ImageName:   imageName,
			}
			apps = append(apps, app)
		}
		return apps, nil
	} else {
		return nil, &ErrorResponse{
			Message:    "Something went wrong, please try after sometime",
			StatusCode: "500",
		}
	}
}

func GetEnvironmentsFromDB() ([]Environment, *ErrorResponse) {
	conn, _ := GetPostgresConnection()
	defer conn.Close(context.Background())
	if conn != nil {
		var envs []Environment
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
			env := Environment{
				Name:        name,
				Id:          id,
				Description: description,
			}
			envs = append(envs, env)
		}

		return envs, nil
	} else {
		errorResponse := ErrorResponse{
			Message:    "Something went wrong, please try again later",
			StatusCode: "500",
		}
		return nil, &errorResponse
	}
}

func GetAllPostgresFromDB() ([]PostgresApp, *ErrorResponse) {
	conn, err := GetPostgresConnection()
	if err != nil {
		return nil, &ErrorResponse{
			Message:    "Something went wrong, please try after sometime",
			StatusCode: "500",
		}
	}
	defer conn.Close(context.Background())
	if conn != nil {
		var postgresList []PostgresApp
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
			err = rows.Scan(&id, &status, &description, &serverName, &adminUser, &password, &version, &environment, &sizeDisk, &storageType, &sizeCPU, &sizeMemory)
			if err != nil {
				fmt.Println(err)
			}
			postgres := PostgresApp{
				Id:          id,
				ServerName:  serverName,
				Description: description,
				Status:      status,
			}
			postgresList = append(postgresList, postgres)
		}
		return postgresList, nil
	} else {
		return nil, &ErrorResponse{
			Message:    "Something went wrong, please try after sometime",
			StatusCode: "500",
		}
	}
}

func GetPostgresConnection() (*pgx.Conn, error) {
	DATABASE_URL := "postgres://postgres:mysecretpassword@192.168.1.41:5432/marketplace"
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, errors.New("unable to connect to database")
	}
	//defer conn.Close(context.Background())
	return conn, nil
}
