package applicationcontroller

import (
	"encoding/json"
	repositorty "marketplace-api/application-repository"
	appStruct "marketplace-api/application-struct"
	helmservice "marketplace-api/helm-service"
	utility "marketplace-api/utility"
	"net/http"
	"strconv"
)

func GetApplications(w http.ResponseWriter, r *http.Request) {
	// enable cors
	utility.EnableCors(&w)
	// Dynamically access the path variable
	//fmt.Fprintf(w, "Retrieving item with ID: %s", item)
	applications, err := repositorty.GetApplicationsFromDB()
	if err != nil {
		errCode, _ := strconv.Atoi(err.StatusCode)
		w.WriteHeader(errCode)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(applications)
	}

}

func GetApplicationVersions(w http.ResponseWriter, r *http.Request) {
	utility.EnableCors(&w)
	id := r.PathValue("id")
	i, _ := strconv.Atoi(id)
	versions, err := repositorty.GetApplicationVersionsFromID(i)
	appVersions := utility.BuildResponseForVersions(versions)
	res := appStruct.AppVersionsResponse{
		Status:      "Success",
		AppVersions: appVersions,
	}
	if err != nil || appVersions == nil {
		w.WriteHeader(http.StatusNotFound)
		res = appStruct.AppVersionsResponse{
			Status:      "Failure",
			AppVersions: nil,
		}
	}
	json.NewEncoder(w).Encode(res)

}

func GetEnvironments(w http.ResponseWriter, r *http.Request) {
	utility.EnableCors(&w)
	// Dynamically access the path variable
	//fmt.Fprintf(w, "Retrieving item with ID: %s", item)
	environments, err := repositorty.GetEnvironmentsFromDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(environments)
	}
}

func AppHelmUpgrade(w http.ResponseWriter, r *http.Request) {
	utility.EnableCors(&w)
	release, err := helmservice.InstallAndUpgrade()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(release)
	}
}

