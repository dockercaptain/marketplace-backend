package main

import (
	"fmt"
	appController "marketplace-api/application-controller"
	pgController "marketplace-api/postgre-app-controller"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fmt.Println("localhost:8080...")
	mux.HandleFunc("GET /marketplace/apps", appController.GetApplications)
	mux.HandleFunc("GET /marketplace/apps/versions/{id}", appController.GetApplicationVersions)
	mux.HandleFunc("GET /marketplace/apps/envs", appController.GetEnvironments)
	mux.HandleFunc("GET /marketplace/apps/{appName}", pgController.GetAllAppsDetailsByName)
	mux.HandleFunc("POST /marketplace/apps/{appName}/create", pgController.CreateAppBasics)
	mux.HandleFunc("GET /marketplace/apps/{appName}/{id}", pgController.GetApplicationById)
	mux.HandleFunc("POST /marketplace/apps/{appName}/helmUpgrade", appController.AppHelmUpgrade)
	mux.HandleFunc("GET /marketplace/apps/release/{appName}", appController.GetAppRelease)
	http.ListenAndServe(":8080", mux)
}
