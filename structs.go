package main

type KubeBackup struct {
	ID          string `json:"_id"`
	Name        string `json:"name"`
	ClusterID   string `json:"cluster"`
	CCUserEmail string `json:"cc_user_email"`
	Source      Source `json:"source"`
	Velero      Velero `json:"velero"`
	Status      Status `json:"status"`
	UID         string `json:"velero_k8s_uid"`
}
type BackupRequest struct {
	TriggerType string    `json:"trigger_type"`
	Cluster     string    `json:"cluster"`
	Name        string    `json:"name"`
	Source      SourceReq `json:"source"`
}

type LabelSelector map[string]string

type SourceReq struct {
	Namespaces    []string      `json:"namespaces"`
	LabelSelector LabelSelector `json:"label_selector"`
}
type Source struct {
	AllNamespaces          bool     `json:"all_namespaces"`
	SnapshotPersistentVols bool     `json:"snapshotPersistentVolumes"`
	Namespaces             []string `json:"namespaces"`
	IncludedResources      []string `json:"included_resources"`
	ExcludedNamespaces     []string `json:"excluded_namespaces"`
	SnapshotLonghorn       bool     `json:"snapshot_longhorn"`
	IncludeUnattachedPVCs  bool     `json:"include_unattached_pvcs"`
	CSISnapshotTimeout     int      `json:"csi_snapshot_timeout"`
}

type Velero struct {
	FsBackup           bool     `json:"fs_backup"`
	StorageLocation    string   `json:"storage_location"`
	VolumeSnapshotLocs []string `json:"volume_snapshot_locations"`
	RetentionDays      int      `json:"retention_days"`
}

type Status struct {
	Velero struct {
		OpState   string `json:"op_state"`
		OpMessage string `json:"op_message"`
	} `json:"velero"`
	Jobs []struct {
		JobID   string `json:"jobid"`
		State   string `json:"state"`
		Message string `json:"message"`
	} `json:"jobs"`
	LastJobRunTime int `json:"last_job_run_time"`
}

type BackupInstance struct {
	ID            string `json:"_id"`
	Name          string `json:"name"`
	ClusterName   string `json:"cluster_name"`
	State         string `json:"state"`
	PVCount       int    `json:"pv_count"`
	TotalSnapData int    `json:"total_snapshot_data"`
	TotalCopyData int    `json:"total_copy_data"`
}
