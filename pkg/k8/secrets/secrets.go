package secrets

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/okteto/cnd/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//Create creates the syncthing config secret
func Create(d *appsv1.Deployment, cndManifests []*model.Dev, c *kubernetes.Clientset) error {
	s, err := c.Core().Secrets(d.Namespace).Get(cndManifests[0].GetCNDSyncSecret(), metav1.GetOptions{})
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return fmt.Errorf("Error getting kubernetes secret: %s", err)
	}
	config, err := getConfigXML(cndManifests)
	if err != nil {
		return fmt.Errorf("Error generating syncthing configuration: %s", err)
	}
	data := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: cndManifests[0].GetCNDSyncSecret()},
		Type:       v1.SecretTypeOpaque,
		Data: map[string][]byte{
			"config.xml": config,
			"cert.pem":   []byte(certPEM),
			"key.pem":    []byte(keyPEM),
		},
	}
	if s.Name == "" {
		_, err := c.Core().Secrets(d.Namespace).Create(data)
		if err != nil {
			return fmt.Errorf("Error creating kubernetes sync secret: %s", err)
		}
		log.Info("Created syncthing secret '%s'.", cndManifests[0].GetCNDSyncSecret())
	} else {
		_, err := c.Core().Secrets(d.Namespace).Update(data)
		if err != nil {
			return fmt.Errorf("Error updating kubernetes sync secret: %s", err)
		}
		log.Info("Sync secret '%s' was updated.", cndManifests[0].GetCNDSyncSecret())
	}
	return nil
}

//Delete deletes the syncthing config secret
func Delete(d *appsv1.Deployment, dev *model.Dev, c *kubernetes.Clientset) error {
	err := c.Core().Secrets(d.Namespace).Delete(dev.GetCNDSyncSecret(), &metav1.DeleteOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil
		}
		return fmt.Errorf("Error deleting kubernetes sync secret: %s", err)
	}
	return nil
}
