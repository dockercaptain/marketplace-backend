package postgreappcontroller

import (
	"encoding/json"
	repositorty "marketplace-api/postgre-app-repository"
	pgStruct "marketplace-api/postgre-struct"
	utility "marketplace-api/utility"
	"net/http"
	"strconv"
)

func GetAllAppsDetailsByName(w http.ResponseWriter, r *http.Request) {
	utility.EnableCors(&w)
	appName := r.PathValue("appName")
	var postgresList []pgStruct.PostgresApp
	var err *pgStruct.ErrorResponse
	if appName == "postgres" {
		postgresList, err = repositorty.GetAllPostgresFromDB()
	} else {
		w.WriteHeader(http.StatusNotFound)
		err = &pgStruct.ErrorResponse{
			Message:    "Application not found",
			StatusCode: "404",
		}
		json.NewEncoder(w).Encode(err)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(postgresList)
	}
}

func CreateAppBasics(w http.ResponseWriter, r *http.Request) {
	utility.EnableCors(&w)
	appName := r.PathValue("appName")
	if appName == "postgres" {
		var pgApp pgStruct.PostgresApp
		// convert json to golang struct and map to struct reference
		err := json.NewDecoder(r.Body).Decode(&pgApp)
		if err != nil {
			errorResponse := &pgStruct.ErrorResponse{
				Message:    "Please pass valid input",
				StatusCode: "400",
			}
			json.NewEncoder(w).Encode(errorResponse)
		} else {
			// save the record and return response
			success, err := repositorty.CreateApplicationPostgres(pgApp)
			if err != nil {
				errCode, _ := strconv.Atoi(err.StatusCode)
				w.WriteHeader(errCode)
				json.NewEncoder(w).Encode(err)
			} else {
				json.NewEncoder(w).Encode(success)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		err := &pgStruct.ErrorResponse{
			Message:    "Application not found",
			StatusCode: "404",
		}
		json.NewEncoder(w).Encode(err)
	}
}

func GetApplicationById(w http.ResponseWriter, r *http.Request) {
	utility.EnableCors(&w)
	id := r.PathValue("id")
	i, _ := strconv.Atoi(id)
	appName := r.PathValue("appName")
	if appName == "postgres" {
		pgApp, err := repositorty.GetApplicationPostgresById(i)
		if err != nil {
			errCode, _ := strconv.Atoi(err.StatusCode)
			w.WriteHeader(errCode)
			json.NewEncoder(w).Encode(err)
		} else {
			json.NewEncoder(w).Encode(pgApp)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		err := &pgStruct.ErrorResponse{
			Message:    "Application not found",
			StatusCode: "404",
		}
		json.NewEncoder(w).Encode(err)
	}
}
