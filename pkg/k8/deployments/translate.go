package deployments

import (
	"encoding/json"
	"fmt"

	"github.com/okteto/cnd/pkg/model"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
)

const (
	initSyncImageTag = "okteto/init-syncthing:0.4.0"
	syncImageTag     = "okteto/syncthing:0.4.0"
)

var (
	devReplicas                      int32 = 1
	devTerminationGracePeriodSeconds int64
)

func translateToDevModeDeployment(d *appsv1.Deployment, dev *model.Dev) ([]*model.Dev, error) {

	d.Status = appsv1.DeploymentStatus{}
	manifest, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	setAnnotation(d.GetObjectMeta(), model.CNDDeploymentAnnotation, string(manifest))
	setLabel(d.GetObjectMeta(), model.CNDLabel, d.Name)
	setLabel(d.Spec.Template.GetObjectMeta(), model.CNDLabel, d.Name)
	cndManifest, err := json.Marshal(dev)
	if err != nil {
		return nil, err
	}
	setAnnotation(d.GetObjectMeta(), dev.GetCNDManifestAnnotation(), string(cndManifest))
	d.Spec.Strategy = appsv1.DeploymentStrategy{Type: appsv1.RecreateDeploymentStrategyType}
	d.Spec.Template.Spec.TerminationGracePeriodSeconds = &devTerminationGracePeriodSeconds

	fmt.Printf("Annotated deployment: %+v\n", d)

	cndManifests, err := getCNDManifests(d)
	if err != nil {
		return nil, err
	}
	fmt.Printf("All CND manifests: %+v\n", cndManifests)

	for _, cndManifest := range cndManifests {
		for i, c := range d.Spec.Template.Spec.Containers {
			if c.Name == cndManifest.Swap.Deployment.Container {
				updateCndContainer(&d.Spec.Template.Spec.Containers[i], cndManifest)
				break
			}
		}
		createInitSyncthingContainer(d, cndManifest)
		createSyncthingVolume(d, cndManifest)
	}

	createSyncthingContainer(d, cndManifests)

	fmt.Printf("Final deployment: %+v\n", d)
	if *(d.Spec.Replicas) != devReplicas {
		log.Info("cnd only supports running with 1 replica")
		d.Spec.Replicas = &devReplicas
	}
	return cndManifests, nil
}

func updateCndContainer(c *apiv1.Container, dev *model.Dev) {
	if dev.Swap.Deployment.Image != "" {
		c.Image = dev.Swap.Deployment.Image
	}

	if len(dev.Swap.Deployment.Command) > 0 {
		c.Command = dev.Swap.Deployment.Command
	}
	if len(dev.Swap.Deployment.Args) > 0 {
		c.Args = dev.Swap.Deployment.Args
	}

	c.WorkingDir = dev.Mount.Target
	c.ReadinessProbe = nil
	c.LivenessProbe = nil

	if c.VolumeMounts == nil {
		c.VolumeMounts = []apiv1.VolumeMount{}
	}

	volumeMount := apiv1.VolumeMount{
		Name:      dev.GetCNDSyncVolume(),
		MountPath: dev.Mount.Target,
	}

	c.VolumeMounts = append(
		c.VolumeMounts,
		volumeMount,
	)

	c.Resources = apiv1.ResourceRequirements{}
}

func createInitSyncthingContainer(d *appsv1.Deployment, dev *model.Dev) {
	initSyncthingContainer := apiv1.Container{
		Name:  dev.GetCNDInitSyncContainer(),
		Image: initSyncImageTag,
		VolumeMounts: []apiv1.VolumeMount{
			apiv1.VolumeMount{
				Name:      dev.GetCNDSyncVolume(),
				MountPath: "/src",
			},
		},
	}

	if d.Spec.Template.Spec.InitContainers == nil {
		d.Spec.Template.Spec.InitContainers = []apiv1.Container{}
	}

	d.Spec.Template.Spec.InitContainers = append(d.Spec.Template.Spec.InitContainers, initSyncthingContainer)
}

func createSyncthingContainer(d *appsv1.Deployment, cndManifests []*model.Dev) {
	syncthingContainer := apiv1.Container{
		Name:            model.CNDSyncContainer,
		Image:           syncImageTag,
		ImagePullPolicy: apiv1.PullAlways,
		VolumeMounts:    []apiv1.VolumeMount{},
		Ports: []apiv1.ContainerPort{
			apiv1.ContainerPort{
				ContainerPort: 8384,
			},
			apiv1.ContainerPort{
				ContainerPort: 22000,
			},
		},
	}
	for _, dev := range cndManifests {
		syncthingContainer.VolumeMounts = append(
			syncthingContainer.VolumeMounts,
			apiv1.VolumeMount{
				Name:      dev.GetCNDSyncVolume(),
				MountPath: dev.GetCNDSyncMount(),
			},
		)
	}
	d.Spec.Template.Spec.Containers = append(d.Spec.Template.Spec.Containers, syncthingContainer)
}

func createSyncthingVolume(d *appsv1.Deployment, dev *model.Dev) {
	if d.Spec.Template.Spec.Volumes == nil {
		d.Spec.Template.Spec.Volumes = []apiv1.Volume{}
	}

	syncVolume := apiv1.Volume{Name: dev.GetCNDSyncVolume()}

	d.Spec.Template.Spec.Volumes = append(
		d.Spec.Template.Spec.Volumes,
		syncVolume,
	)
}
