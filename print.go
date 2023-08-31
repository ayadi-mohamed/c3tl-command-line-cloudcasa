package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func printKubeBackups(kubeBackups []KubeBackup) {
	if len(kubeBackups) == 0 {
		fmt.Println("No Kube Backups found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "Kube Backup ID\tName\tCluster ID\tCC User Email\tSource\tVelero\tStatus\n")
	fmt.Fprintf(w, "----------------\t----\t----------\t--------------\t------\t------\t------\n")
	for _, kubeBackup := range kubeBackups {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			kubeBackup.ID,
			kubeBackup.Name,
			kubeBackup.ClusterID,
			kubeBackup.CCUserEmail,
			formatSource(kubeBackup.Source),
			formatVelero(kubeBackup.Velero),
			formatStatus(kubeBackup.Status))
	}
	w.Flush()
}

func formatSource(source Source) string {
	return fmt.Sprintf("All Namespaces: %v\n"+
		"Snapshot Persistent Volumes: %v\n"+
		"Namespaces: %s\n"+
		"Included Resources: %s\n"+
		"Excluded Namespaces: %s\n"+
		"Snapshot Longhorn: %v\n"+
		"Include Unattached PVCs: %v\n"+
		"CSI Snapshot Timeout: %d",
		source.AllNamespaces,
		source.SnapshotPersistentVols,
		strings.Join(source.Namespaces, ", "),
		strings.Join(source.IncludedResources, ", "),
		strings.Join(source.ExcludedNamespaces, ", "),
		source.SnapshotLonghorn,
		source.IncludeUnattachedPVCs,
		source.CSISnapshotTimeout)
}

func formatVelero(velero Velero) string {
	return fmt.Sprintf("File System Backup: %v\n"+
		"Storage Location: %s\n"+
		"Volume Snapshot Locations: %s\n"+
		"Retention Days: %d",
		velero.FsBackup,
		velero.StorageLocation,
		strings.Join(velero.VolumeSnapshotLocs, ", "),
		velero.RetentionDays)
}

func formatStatus(status Status) string {
	var jobs []string
	for _, job := range status.Jobs {
		jobs = append(jobs, fmt.Sprintf("Job ID: %s, State: %s, Message: %s", job.JobID, job.State, job.Message))
	}

	return fmt.Sprintf("Velero:\n"+
		"  Operation State: %s\n"+
		"  Operation Message: %s\n"+
		"Jobs:\n%s\n"+
		"Last Job Run Time: %d",
		status.Velero.OpState,
		status.Velero.OpMessage,
		strings.Join(jobs, "\n"),
		status.LastJobRunTime)
}
