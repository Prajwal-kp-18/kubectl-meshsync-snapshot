package meshsync

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Prajwal-kp-18/kubectl-meshsync-snapshot/pkg/kube"
)

// DeployOptions contains options for deploying MeshSync
type DeployOptions struct {
	Namespace string
	Version   string
}

// CaptureOptions contains options for capturing snapshot
type CaptureOptions struct {
	Namespace     string
	AllNamespaces bool
}

// CleanupOptions contains options for cleaning up MeshSync
type CleanupOptions struct {
	Namespace string
	Force     bool
}

// Snapshot represents a MeshSync snapshot
type Snapshot struct {
	APIVersion string                 `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                 `json:"kind" yaml:"kind"`
	Metadata   map[string]interface{} `json:"metadata" yaml:"metadata"`
	Resources  []Resource             `json:"resources" yaml:"resources"`
}

// Resource represents a kubernetes resource in the snapshot
type Resource struct {
	APIVersion string                 `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                 `json:"kind" yaml:"kind"`
	Metadata   map[string]interface{} `json:"metadata" yaml:"metadata"`
	Spec       map[string]interface{} `json:"spec,omitempty" yaml:"spec,omitempty"`
	Status     map[string]interface{} `json:"status,omitempty" yaml:"status,omitempty"`
}

// Deploy deploys MeshSync to the cluster
func Deploy(ctx context.Context, client *kube.Client, opts DeployOptions) error {
	// Create namespace if it doesn't exist
	_, err := client.Clientset.CoreV1().Namespaces().Get(ctx, opts.Namespace, metav1.GetOptions{})
	if err != nil {
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: opts.Namespace,
			},
		}
		_, err = client.Clientset.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create namespace %s: %w", opts.Namespace, err)
		}
	}

	// Create MeshSync deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "meshsync",
			Namespace: opts.Namespace,
			Labels: map[string]string{
				"app": "meshsync",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "meshsync",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "meshsync",
					},
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: "meshsync",
					Containers: []corev1.Container{
						{
							Name:  "meshsync",
							Image: fmt.Sprintf("layer5/meshsync:%s", opts.Version),
							Ports: []corev1.ContainerPort{
								{
									Name:          "api",
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}
	_, err = client.Clientset.AppsV1().Deployments(opts.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create MeshSync deployment: %w", err)
	}

	// Create service account with necessary permissions
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "meshsync",
			Namespace: opts.Namespace,
		},
	}
	_, err = client.Clientset.CoreV1().ServiceAccounts(opts.Namespace).Create(ctx, sa, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service account: %w", err)
	}

	// Create service for meshsync
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "meshsync",
			Namespace: opts.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "meshsync",
			},
			Ports: []corev1.ServicePort{
				{
					Name:     "api",
					Port:     8080,
					Protocol: corev1.ProtocolTCP,
				},
			},
		},
	}
	_, err = client.Clientset.CoreV1().Services(opts.Namespace).Create(ctx, svc, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	// Wait for deployment to be ready
	for start := time.Now(); time.Since(start) < 2*time.Minute; {
		deploy, err := client.Clientset.AppsV1().Deployments(opts.Namespace).Get(ctx, "meshsync", metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to get deployment status: %w", err)
		}
		if deploy.Status.ReadyReplicas > 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}

	return fmt.Errorf("timeout waiting for MeshSync deployment to be ready")
}

// Validate checks if MeshSync is running in the cluster
func Validate(ctx context.Context, client *kube.Client, namespace string) error {
	// Check if MeshSync deployment exists and is ready
	deploy, err := client.Clientset.AppsV1().Deployments(namespace).Get(ctx, "meshsync", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("MeshSync deployment not found: %w", err)
	}

	if deploy.Status.ReadyReplicas == 0 {
		return fmt.Errorf("MeshSync deployment is not ready")
	}

	// Check if MeshSync service exists
	_, err = client.Clientset.CoreV1().Services(namespace).Get(ctx, "meshsync", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("MeshSync service not found: %w", err)
	}

	return nil
}

// CaptureSnapshot captures cluster state using MeshSync
func CaptureSnapshot(ctx context.Context, client *kube.Client, opts CaptureOptions) (*Snapshot, error) {
	// This is a simplified implementation
	// In a real implementation, this would:
	// 1. Connect to the MeshSync API or query its data store
	// 2. Gather all resources based on filters
	// 3. Format them into a structured snapshot

	// For demo purposes, we'll construct a simple snapshot
	snapshot := &Snapshot{
		APIVersion: "meshery.layer5.io/v1alpha1",
		Kind:       "MeshSync",
		Metadata: map[string]interface{}{
			"name":      "kubernetes-snapshot",
			"timestamp": time.Now().Format(time.RFC3339),
		},
		Resources: []Resource{},
	}

	// Get namespaces
	namespaces := []string{opts.Namespace}
	if opts.AllNamespaces {
		nsList, err := client.Clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list namespaces: %w", err)
		}
		namespaces = []string{}
		for _, ns := range nsList.Items {
			namespaces = append(namespaces, ns.Name)
		}
	}

	// For each namespace, get deployments
	for _, ns := range namespaces {
		deployments, err := client.Clientset.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to list deployments in namespace %s: %w", ns, err)
		}

		for _, deploy := range deployments.Items {
			// Convert deployment to map for storage in snapshot
			deployBytes, err := json.Marshal(deploy)
			if err != nil {
				return nil, fmt.Errorf("error marshaling deployment: %w", err)
			}

			var deployMap map[string]interface{}
			err = json.Unmarshal(deployBytes, &deployMap)
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling deployment: %w", err)
			}

			resource := Resource{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Metadata:   deployMap["metadata"].(map[string]interface{}),
				Spec:       deployMap["spec"].(map[string]interface{}),
				Status:     deployMap["status"].(map[string]interface{}),
			}

			snapshot.Resources = append(snapshot.Resources, resource)
		}
	}

	return snapshot, nil
}

// SaveSnapshot saves the snapshot to a file
func SaveSnapshot(snapshot *Snapshot, filePath string, format string) error {
	var data []byte
	var err error

	if format == "json" {
		data, err = json.MarshalIndent(snapshot, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal snapshot to JSON: %w", err)
		}
	} else {
		// Default to YAML
		data, err = yaml.Marshal(snapshot)
		if err != nil {
			return fmt.Errorf("failed to marshal snapshot to YAML: %w", err)
		}
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write snapshot to file: %w", err)
	}

	return nil
}

// Cleanup removes MeshSync resources from the cluster
func Cleanup(ctx context.Context, client *kube.Client, opts CleanupOptions) error {
	// Delete deployment
	err := client.Clientset.AppsV1().Deployments(opts.Namespace).Delete(ctx, "meshsync", metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete MeshSync deployment: %w", err)
	}

	// Delete service
	err = client.Clientset.CoreV1().Services(opts.Namespace).Delete(ctx, "meshsync", metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete MeshSync service: %w", err)
	}

	// Delete service account
	err = client.Clientset.CoreV1().ServiceAccounts(opts.Namespace).Delete(ctx, "meshsync", metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete MeshSync service account: %w", err)
	}

	// Delete namespace if it's empty and force is true
	if opts.Force {
		// Check if namespace is empty
		deployments, err := client.Clientset.AppsV1().Deployments(opts.Namespace).List(ctx, metav1.ListOptions{})
		if err != nil {
			return fmt.Errorf("failed to list deployments in namespace %s: %w", opts.Namespace, err)
		}

		if len(deployments.Items) == 0 {
			err = client.Clientset.CoreV1().Namespaces().Delete(ctx, opts.Namespace, metav1.DeleteOptions{})
			if err != nil {
				return fmt.Errorf("failed to delete namespace %s: %w", opts.Namespace, err)
			}
		}
	}

	return nil
} 