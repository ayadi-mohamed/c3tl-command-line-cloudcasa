package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	//	"k8s.io/client-go/rest"
	//	"k8s.io/client-go/tools/clientcmd"
	//
	// "k8s.io/client-go/tools/clientcmd"
	// "k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string
var authorizationKey string
var ApiserverCloudcasa string

func init() {

	flag.StringVar(&kubeconfig, "kubeconfig", "/home/mayadi/.kube/config", "path to Kubernetes config file")
	flag.Parse()

	file, err := os.Open("settings.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode JSON
	var data map[string]interface{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Access values using keys
	authorizationKey = data["authorizationkey"].(string)
	ApiserverCloudcasa = data["ApiserverCloudcasa"].(string)

	// Print the retrieved values
	fmt.Println("Authorization Key:", authorizationKey)
	fmt.Println("API Server Hostname:", ApiserverCloudcasa)

}

func main() {

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)

	createCmd := flag.NewFlagSet("create", flag.ExitOnError)

	switch os.Args[1] {
	case "list":

		listBackupDefs := listCmd.String("backupdefs", "", "List backup definitions per Cluster")
		listBackupInst := listCmd.Bool("backupinstances", false, "List backup instances")
		listCmd.Parse(os.Args[2:])
		if *listBackupDefs != "" {

			clusterID := listBackupDefs

			kubeBackups, err := listBackupDefinitions(*clusterID)
			if err != nil {
				fmt.Println("Error fetching kubebackups:", err)

			}

			fmt.Println("Kube Backups:")

			printKubeBackups(kubeBackups)

		} else if *listBackupInst {
			listBackupInstances()
		} else {
			fmt.Println("Usage: mycommandline list <subcommand>")
			fmt.Println("Subcommands:")
			fmt.Println("  --backupdefs       List backup definitions Per Cluster")
			fmt.Println("  --backupinstances  List backup instances")
		}
	case "create":
		trigger_type := createCmd.String("trigger-type", "SCHEDULED", "trigger type: ADHOC or SCHEDULED")
		clusterIdCreate := createCmd.String("cluster", "", "cluster: Cluster ID")
		name := createCmd.String("name", "", "name: Backup name")

		namespaces := createCmd.String("namespaces", "all", "namespaces: ALL per default")

		label_selectorCreate := createCmd.String("label-selector", "", "Label: value")

		createCmd.Parse(os.Args[2:])

		fmt.Printf("trigger_type: %s ", *trigger_type)
		fmt.Printf("cluster: %s ", *clusterIdCreate)
		fmt.Printf("name: %s ", *name)
		fmt.Printf("ns: %s ", *namespaces)
		fmt.Printf("labelselector: %s ", *label_selectorCreate)
		label_selector_parsed, err := parseLabelSelector(*label_selectorCreate)

		source := SourceReq{
			Namespaces:    parseNamespaces(*namespaces),
			LabelSelector: label_selector_parsed,
		}
		if err != nil {
			panic(err)
		}

		KubeBackup, err := createBackup(*trigger_type, *clusterIdCreate, *name, source, "https://ui.staging.cloudcasa.io/api/v1/kubebackups")

		config, err := initConfig(kubeconfig)

		kubeBackupCheck(KubeBackup, config, *clusterIdCreate)

	case "check":
		cluster := checkCmd.String("cluster", "", "check backup definitions per Cluster")
		config, err := initConfig(kubeconfig)

		if err != nil {
			panic(err)
		}
		checkCmd.Parse(os.Args[2:])
		if *cluster != "" {

			backups, err := getBackupsCRDList(config, *cluster)
			if err == nil {
				kubeBackups, err := listBackupDefinitions(*cluster)
				if err != nil {
					panic(err)

				}
				compareBackupsList(backups, kubeBackups)
			} else {
				fmt.Printf("%s", err)
			}

		} else {
			fmt.Println("Usage: mycommandline check <subcommand>")
			fmt.Println("Subcommands:")
			fmt.Println("  --cluster       check backup definitions Per Cluster")
		}

	default:
		fmt.Println("Unknown command:", os.Args[1])
		os.Exit(1)
	}
}
