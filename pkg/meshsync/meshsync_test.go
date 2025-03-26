package meshsync

import (
	"os"
	"testing"
)

func TestSaveSnapshotYAML(t *testing.T) {
	// Create a test snapshot
	snapshot := &Snapshot{
		APIVersion: "meshery.layer5.io/v1alpha1",
		Kind:       "MeshSync",
		Metadata: map[string]interface{}{
			"name": "test-snapshot",
		},
		Resources: []Resource{
			{
				APIVersion: "v1",
				Kind:       "Pod",
				Metadata: map[string]interface{}{
					"name":      "test-pod",
					"namespace": "default",
				},
				Spec: map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "nginx",
							"image": "nginx:latest",
						},
					},
				},
			},
		},
	}

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "snapshot-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Save snapshot to file
	err = SaveSnapshot(snapshot, tmpfile.Name(), "yaml")
	if err != nil {
		t.Fatalf("SaveSnapshot() error = %v", err)
	}

	// Check if file exists and is not empty
	info, err := os.Stat(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("Snapshot file is empty")
	}
}

func TestSaveSnapshotJSON(t *testing.T) {
	// Create a test snapshot
	snapshot := &Snapshot{
		APIVersion: "meshery.layer5.io/v1alpha1",
		Kind:       "MeshSync",
		Metadata: map[string]interface{}{
			"name": "test-snapshot",
		},
		Resources: []Resource{
			{
				APIVersion: "v1",
				Kind:       "Pod",
				Metadata: map[string]interface{}{
					"name":      "test-pod",
					"namespace": "default",
				},
				Spec: map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "nginx",
							"image": "nginx:latest",
						},
					},
				},
			},
		},
	}

	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "snapshot-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Save snapshot to file
	err = SaveSnapshot(snapshot, tmpfile.Name(), "json")
	if err != nil {
		t.Fatalf("SaveSnapshot() error = %v", err)
	}

	// Check if file exists and is not empty
	info, err := os.Stat(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("Snapshot file is empty")
	}
} 