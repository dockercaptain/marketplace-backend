package helmservice

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	pgStruct "marketplace-api/postgre-struct"
	"marketplace-api/utility"

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/mittwald/go-helm-client/values"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

var chartRepo = repo.Entry{
	Name: "marketplace-helm",
	//URL:  "https://charts.helm.sh/stable",
	// URL: "https://charts.bitnami.com/bitnami",
	//URL: "https://github.com/dockercaptain/",
	// Username: "dockercaptain",
	// Password: "",
	URL: "https://dockercaptain.github.io/marketplace-helm/",
}

func getkubeconfigListByte() ([]byte, *pgStruct.ErrorResponse) {
	//  read kubeconfig file from golang io operation:-
	file, err := os.Open("localconfig.kubeconfig")
	if err != nil {
		errRes := utility.GetErrorResponse(err, "500")
		return nil, errRes
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		errRes := utility.GetErrorResponse(err, "500")
		return nil, errRes
	}
	kubeconfigListByte := make([]byte, stat.Size())
	return kubeconfigListByte, nil
}
func getKubeOption() (*helmclient.KubeConfClientOptions, *pgStruct.ErrorResponse) {
	kubeconfigListByte, err := getkubeconfigListByte()
	if err != nil {
		return nil, err
	}
	return &helmclient.KubeConfClientOptions{
		Options: &helmclient.Options{
			Namespace:        "test-del-ns", // Change this to the namespace you wish to install the chart in.
			RepositoryCache:  "/tmp/.helmcache",
			RepositoryConfig: "/tmp/.helmrepo",
			Debug:            true,
			Linting:          true, // Change this to false if you don't want linting.
			DebugLog: func(format string, v ...interface{}) {
				// Change this to your own logger. Default is 'log.Printf(format, v...)'.
			},
		},
		KubeContext: "",
		KubeConfig:  kubeconfigListByte,
	}, nil
}

func InstallAndUpgrade() (*release.Release, *pgStruct.ErrorResponse) {
	//  read kubeconfig file from golang io operation:-
	file, _ := os.Open("localconfig.kubeconfig")
	defer file.Close()
	stat, _ := file.Stat()
	kubeconfigListByte := make([]byte, stat.Size())
	_, err := bufio.NewReader(file).Read(kubeconfigListByte)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return nil, &pgStruct.ErrorResponse{
			Message:    err.Error(),
			StatusCode: "500",
		}
	}
	kubeOption := &helmclient.KubeConfClientOptions{
		Options: &helmclient.Options{
			Namespace:        "test-del-ns", // Change this to the namespace you wish to install the chart in.
			RepositoryCache:  "/tmp/.helmcache",
			RepositoryConfig: "/tmp/.helmrepo",
			Debug:            true,
			Linting:          true, // Change this to false if you don't want linting.
			DebugLog: func(format string, v ...interface{}) {
				// Change this to your own logger. Default is 'log.Printf(format, v...)'.
			},
		},
		KubeContext: "",
		KubeConfig:  kubeconfigListByte,
	}

	genericHelmClient, err := helmclient.New(&helmclient.Options{})
	if err != nil {
		log.Println(err)
		return nil, &pgStruct.ErrorResponse{
			Message:    err.Error(),
			StatusCode: "500",
		}
	}

	err = genericHelmClient.AddOrUpdateChartRepo(chartRepo)
	if err != nil {
		log.Println(err)
		return nil, &pgStruct.ErrorResponse{
			Message:    err.Error(),
			StatusCode: "500",
		}
	}
	fmt.Println("Add or update chart")
	// `primary.resources.requests={"memory": "200Mi"}`,
	//`primary.resources.requests={"cpu": "100m"}`,
	// above request should be constant and limit to be
	// provided by user
	chartSpec := helmclient.ChartSpec{
		ReleaseName:     "postgresql",
		ChartName:       "marketplace-helm/postgresql-ha",
		Namespace:       "test-del-ns",
		UpgradeCRDs:     true,
		Wait:            true,
		Timeout:         300 * time.Second,
		Version:         "11.9.4",
		CreateNamespace: true,
		DryRun:          false,
		ValuesOptions: values.Options{JSONValues: []string{
			`persistence={"size":"3Gi"}`,
			`global.postgresql.auth={"postgresPassword":"NewPassword1"}`,
			`global={"defaultStorageClass":"azurefile"}`,
			`postgresql.resources={"limits":{"memory": "1000Mi", "cpu": "200m"},"requests":{"memory": "200Mi", "cpu": "100m"}}`,
			`postgresql.nodeSelector={}`,
		},
		},
	}
	kubefileHelmClient, _ := helmclient.NewClientFromKubeConf(kubeOption, helmclient.Burst(1000), helmclient.Timeout(1000*time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Second)
	defer cancel()
	fmt.Println("updating /installing helm chart")
	release, err := kubefileHelmClient.InstallOrUpgradeChart(ctx, &chartSpec, nil)
	if err != nil {
		log.Println("unable to install helm chart...", err)
		return nil, &pgStruct.ErrorResponse{
			Message:    err.Error(),
			StatusCode: "500",
		}
	}
	fmt.Println("dry-run info or installed chart details:- ")
	return release, nil
}

func GetChart(genericHelmClient helmclient.Client) {
	chart, str, err := genericHelmClient.GetChart("postgresql", &action.ChartPathOptions{
		RepoURL: chartRepo.URL,
		Version: "15.2.6",
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range chart.Files {
		if file.Name == "README.md" {
			//fmt.Println(string(file.Data))
			fmt.Println(string(chart.Name()))
		}
	}
	fmt.Println(str)
}

func GetRelease(name string) (*release.Release, *pgStruct.ErrorResponse) {
	kubeOption, err := getKubeOption()
	if err != nil {
		return nil, err
	}
	kubefileHelmClient, error := helmclient.NewClientFromKubeConf(kubeOption, helmclient.Burst(1000), helmclient.Timeout(1000*time.Second))
	if error != nil {
		errRes := utility.GetErrorResponse(error, "500")
		return nil, errRes
	}
	release, error := helmclient.Client.GetRelease(kubefileHelmClient, name)
	if error != nil {
		errRes := utility.GetErrorResponse(error, "500")
		return nil, errRes
	}
	//	fmt.Println(release.Name, release.Namespace, release.Info.Status, release.Chart.Name(), release.Chart.AppVersion(), release.Info.LastDeployed.String())
	return release, nil
}

func UninstallReleaseByName(name string, kubeOption *helmclient.KubeConfClientOptions) {
	kubefileHelmClient, _ := helmclient.NewClientFromKubeConf(kubeOption, helmclient.Burst(1000), helmclient.Timeout(1000*time.Second))
	err := kubefileHelmClient.UninstallReleaseByName("postgresql")
	if err != nil {
		log.Println("unable to uninstall", err)
	}
}

func RollbackRelease(chartSpec *helmclient.ChartSpec, kubeOption *helmclient.KubeConfClientOptions) {
	kubefileHelmClient, _ := helmclient.NewClientFromKubeConf(kubeOption, helmclient.Burst(1000), helmclient.Timeout(1000*time.Second))
	err := kubefileHelmClient.RollbackRelease(chartSpec)
	if err != nil {
		log.Fatal("unable to Rollback to previous release", err)
	}
}
