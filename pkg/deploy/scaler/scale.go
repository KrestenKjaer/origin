package scaler

import (
	"time"

	"github.com/golang/glog"
	kapi "k8s.io/kubernetes/pkg/api"
	kerrors "k8s.io/kubernetes/pkg/api/errors"
	kclient "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/kubectl"
	"k8s.io/kubernetes/pkg/util/wait"

	"github.com/openshift/origin/pkg/client"
	"github.com/openshift/origin/pkg/deploy/api"
	"github.com/openshift/origin/pkg/deploy/util"
)

// NewDeploymentConfigScaler returns a new scaler for deploymentConfigs
func NewDeploymentConfigScaler(oc *client.Client, kc *kclient.Client) kubectl.Scaler {
	return &DeploymentConfigScaler{c: NewScalerClient(oc, kc)}
}

// DeploymentConfigScaler is a wrapper for the kubectl Scaler client
type DeploymentConfigScaler struct {
	c kubectl.ScalerClient
}

// Scale updates a replication controller created by the DeploymentConfig with the provided namespace/name,
// to a new size, with optional precondition check (if preconditions is not nil),optional retries (if retry
//  is not nil), and then optionally waits for it's replica count to reach the new value (if wait is not nil).
func (scaler *DeploymentConfigScaler) Scale(namespace, name string, newSize uint, preconditions *kubectl.ScalePrecondition, retry, waitForReplicas *kubectl.RetryParams) error {
	if preconditions == nil {
		preconditions = &kubectl.ScalePrecondition{Size: -1, ResourceVersion: ""}
	}
	if retry == nil {
		// Make it try only once, immediately
		retry = &kubectl.RetryParams{Interval: time.Millisecond, Timeout: time.Millisecond}
	}
	cond := kubectl.ScaleCondition(scaler, preconditions, namespace, name, newSize)
	if err := wait.Poll(retry.Interval, retry.Timeout, cond); err != nil {
		if scaleErr := err.(kubectl.ControllerScaleError); kerrors.IsNotFound(scaleErr.ActualError) {
			glog.Infof("No deployment found for dc/%s. Scaling the deployment configuration template...", name)
			if _, err := scaler.c.(*realScalerClient).UpdateDeploymentConfig(namespace, name, newSize); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	if waitForReplicas != nil {
		rc, err := scaler.c.GetReplicationController(namespace, name)
		if err != nil {
			return err
		}
		return wait.Poll(waitForReplicas.Interval, waitForReplicas.Timeout,
			scaler.c.ControllerHasDesiredReplicas(rc))
	}
	return nil
}

// ScaleSimple does a simple one-shot attempt at scaling - not useful on it's own, but
// a necessary building block for Scale
func (scaler *DeploymentConfigScaler) ScaleSimple(namespace, name string, preconditions *kubectl.ScalePrecondition, newSize uint) (string, error) {
	const scaled = "scaled"
	controller, err := scaler.c.GetReplicationController(namespace, name)
	if err != nil {
		return "", kubectl.ControllerScaleError{FailureType: kubectl.ControllerScaleGetFailure, ResourceVersion: "Unknown", ActualError: err}
	}
	if preconditions != nil {
		if err := preconditions.Validate(controller); err != nil {
			return "", err
		}
	}
	controller.Spec.Replicas = int(newSize)
	// TODO: do retry on 409 errors here?
	if _, err := scaler.c.UpdateReplicationController(namespace, controller); err != nil {
		return "", kubectl.ControllerScaleError{FailureType: kubectl.ControllerScaleUpdateFailure, ResourceVersion: controller.ResourceVersion, ActualError: err}
	}
	// TODO: do a better job of printing objects here.
	return scaled, nil
}

// NewScalerClient returns a new Scaler client bundling both the OpenShift and
// Kubernetes clients
func NewScalerClient(oc client.Interface, kc kclient.Interface) kubectl.ScalerClient {
	return &realScalerClient{oc: oc, kc: kc}
}

// realScalerClient is a ScalerClient which uses an OpenShift and a Kube client.
type realScalerClient struct {
	oc client.Interface
	kc kclient.Interface
}

// GetReplicationController returns the most recent replication controller associated with the deploymentConfig
// with the provided namespace/name combination
func (c *realScalerClient) GetReplicationController(namespace, name string) (*kapi.ReplicationController, error) {
	dc, err := c.oc.DeploymentConfigs(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return c.kc.ReplicationControllers(namespace).Get(util.LatestDeploymentNameForConfig(dc))
}

// UpdateReplicationController updates the provided replication controller
func (c *realScalerClient) UpdateReplicationController(namespace string, rc *kapi.ReplicationController) (*kapi.ReplicationController, error) {
	return c.kc.ReplicationControllers(namespace).Update(rc)
}

// ControllerHasDesiredReplicas checks whether the provided replication controller has the desired replicas
// number set
func (c *realScalerClient) ControllerHasDesiredReplicas(rc *kapi.ReplicationController) wait.ConditionFunc {
	return kclient.ControllerHasDesiredReplicas(c.kc, rc)
}

// UpdateDeploymentConfig tries to get and update the provided deployment configuration template with the
// provided replicas size
func (c *realScalerClient) UpdateDeploymentConfig(namespace, name string, newSize uint) (*api.DeploymentConfig, error) {
	dc, err := c.oc.DeploymentConfigs(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	dc.Template.ControllerTemplate.Replicas = int(newSize)
	return c.oc.DeploymentConfigs(namespace).Update(dc)
}
