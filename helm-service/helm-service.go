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

	helmclient "github.com/mittwald/go-helm-client"
	"github.com/mittwald/go-helm-client/values"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

var chartRepo = repo.Entry{
	Name: "bitnami",
	//URL:  "https://charts.helm.sh/stable",
	URL: "https://charts.bitnami.com/bitnami",
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

	chartSpec := helmclient.ChartSpec{
		ReleaseName:     "postgresql",
		ChartName:       "bitnami/postgresql",
		Namespace:       "test-del-ns",
		UpgradeCRDs:     true,
		Wait:            true,
		Timeout:         300 * time.Second,
		Version:         "15.2.6",
		CreateNamespace: true,
		DryRun:          true,
		ValuesOptions: values.Options{JSONValues: []string{
			`primary.persistence={"size":"11Gi"}`,
			`global.postgresql.auth={"postgresPassword":"NewPassword1"}`,
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

func GetRelease(genericHelmClient helmclient.Client, kubeOption *helmclient.KubeConfClientOptions) {
	kubefileHelmClient, _ := helmclient.NewClientFromKubeConf(kubeOption, helmclient.Burst(1000), helmclient.Timeout(1000*time.Second))
	release, err := helmclient.Client.GetRelease(kubefileHelmClient, "postgresql")
	if err != nil {
		log.Println("unable to get the release", err)
	}
	fmt.Println(release.Name, release.Namespace, release.Info.Status, release.Chart.Name(), release.Chart.AppVersion(), release.Info.LastDeployed.String())
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
