// commands.go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	v1clientset "github.com/vmware-tanzu/velero/pkg/generated/clientset/versioned/typed/velero/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func listBackupDefinitions(clusterID string) ([]KubeBackup, error) {
	// Fetch and print backup definitions
	kubeBackupsURL := fmt.Sprintf("%skubebackups?where={\"cluster\":\"%s\"}", ApiserverCloudcasa, clusterID)

	client := &http.Client{}
	req, err := http.NewRequest("GET", kubeBackupsURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+authorizationKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var kubeBackupResponse struct {
		Items []KubeBackup `json:"_items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&kubeBackupResponse); err != nil {
		return nil, err
	}

	return kubeBackupResponse.Items, nil
}

func listBackupInstances() {
	// Fetch and print backup instances
	fmt.Println("Listing backup instances...")
	// Implement your logic here
}
func getBackupsCRDList(config *rest.Config, clusterID string) (*v1.BackupList, error) {
	v1.AddToScheme(scheme.Scheme)
	clientset, err := v1clientset.NewForConfig(config)
	if err != nil {
		return nil, err

	}
	backups, err := clientset.Backups("velero").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return backups, nil

}

func compareBackupsList(veleroBackups *v1.BackupList, cloudcasaBackups []KubeBackup) {
	cloudcasaID := make(map[string]bool)
	for _, kubeBackup := range cloudcasaBackups {
		cloudcasaBackup := kubeBackup.UID
		cloudcasaID[cloudcasaBackup] = true

	}

	for _, veleroBackup := range veleroBackups.Items {
		if cloudcasaID[string(veleroBackup.UID)] {

			fmt.Printf("Backup %s UID matches cloudcasa   %s \n", veleroBackup.UID, veleroBackup.UID)
		} else {

			fmt.Printf("Backup %s UID does not matche cloudcasa velero k8s ", veleroBackup.UID)

		}

	}
}

func parseNamespaces(namespacesStr string) []string {
	// Split the input string into individual namespaces using space as the separator
	namespaces := strings.Split(namespacesStr, " ")
	// Remove any empty strings from the resulting slice
	var cleanedNamespaces []string
	for _, ns := range namespaces {
		if ns != "" {
			cleanedNamespaces = append(cleanedNamespaces, ns)
		}
	}
	return cleanedNamespaces
}

func parseLabelSelector(selectorStr string) (LabelSelector, error) {
	parts := strings.Split(selectorStr, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid label selector format")
	}
	return LabelSelector{parts[0]: parts[1]}, nil
}

// Function to create a backup using CloudCasa API
func createBackup(triggerType string, cluster string, name string, source SourceReq, url string) (KubeBackup, error) {

	// Create the backup request payload using the struct
	backupPayload := BackupRequest{
		TriggerType: triggerType,
		Cluster:     cluster,
		Name:        name,
		Source:      source,
	}

	// Convert the payload to JSON
	payloadBytes, err := json.Marshal(backupPayload)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return KubeBackup{}, err
	}
	// Create a POST request with the JSON payload
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return KubeBackup{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authorizationKey)

	// Create an HTTP client and send the request
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return KubeBackup{}, err
	}

	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Request failed with status: %s\n", resp.Status)
		return KubeBackup{}, err
	}

	var kubeBackup KubeBackup

	fmt.Println("Backup request sent successfully!")
	if err := json.NewDecoder(resp.Body).Decode(&kubeBackup); err != nil {
		return KubeBackup{}, err
	}
	return kubeBackup, nil
}

// // Function to initialize Kubernetes configuration
func initConfig(kubeconfig string) (*rest.Config, error) {

	var config *rest.Config
	var err error

	if kubeconfig != "" {
		log.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return config, nil

	}

	log.Printf("using in-cluster configuration")
	config, err = rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return config, nil

}

// Fetching a backup by ID from CloudCasa API
func kubeBackupGetById(id string) (KubeBackup, error) {
	kubeBackupsURL := fmt.Sprintf("%s/kubebackups?where={\"_id\":\"%s\"}", ApiserverCloudcasa, id)

	client := &http.Client{}
	req, err := http.NewRequest("GET", kubeBackupsURL, nil)
	if err != nil {
		return KubeBackup{}, err
	}
	req.Header.Set("Authorization", "Bearer "+authorizationKey)

	resp, err := client.Do(req)
	if err != nil {
		return KubeBackup{}, err
	}
	defer resp.Body.Close()
	var kubeBackupResponse struct {
		Items []KubeBackup `json:"_items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&kubeBackupResponse); err != nil {
		return KubeBackup{}, err
	}
	fmt.Print(kubeBackupResponse.Items)
	return kubeBackupResponse.Items[0], nil
}

// Checking if a backup exists in Kubernetes
func kubeBackupCheck(KubeBackup KubeBackup, config *rest.Config, clusterID string) {
	fmt.Print("\n check check \n ")

	backupsCRDlist, err := getBackupsCRDList(config, clusterID)
	fmt.Print(KubeBackup.ID)
	kubebackupGet, err := kubeBackupGetById(KubeBackup.ID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("UID: %s", kubebackupGet.UID)
	for _, backup := range backupsCRDlist.Items {
		if string(backup.UID) == kubebackupGet.UID {

			fmt.Println("The backup has been launched in kubernetes cluster with uid")
		}

	}
}
