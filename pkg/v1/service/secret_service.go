package service

import (
	"context"
	"errors"
	"secret-manager/pkg/v1/models"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ISecret interface {
	CreateSecret(user string, secretval models.Secret) error
	UpdateSecret(user string, key string, secretval models.UpdateSecret) error
	DeleteSecret(user string, key string) error
	ViewSecret(user string, key string) (map[string][]byte, error)
}

// Create Secret
func (s Service) CreateSecret(user string, secretval models.Secret) error {
	//check if secret with the name secretval.Key already exists
	_, err := s.client.CoreV1().Secrets(user).Get(context.Background(), secretval.Key, metav1.GetOptions{})
	if err == nil {
		logrus.Error(err)
		return errors.New("secret already exists")
	}

	//create new secret
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretval.Key,
		},
		Data: map[string][]byte{
			secretval.Key: []byte(secretval.Value),
		},
	}

	_, err = s.client.CoreV1().Secrets(user).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		logrus.Error(err)
		return errors.New("can't create secret")
	}

	return nil
}

// Update Secret
func (s Service) UpdateSecret(user string, key string, secretval models.UpdateSecret) error {
	//check if secret with the name secretval.Key exists
	secret, err := s.client.CoreV1().Secrets(user).Get(context.Background(), key, metav1.GetOptions{})
	if err != nil {
		logrus.Error(err)
		return errors.New("secret doesn't exist")
	}

	//update the secret
	secret.Data[key] = []byte(secretval.Value)

	_, err = s.client.CoreV1().Secrets(user).Update(context.Background(), secret, metav1.UpdateOptions{})
	if err != nil {
		logrus.Error(err)
		return errors.New("can't update secret")
	}

	return nil
}

// Delete Secret
func (s Service) DeleteSecret(user string, key string) error {
	//check if secret with the name key exists
	_, err := s.client.CoreV1().Secrets(user).Get(context.Background(), key, metav1.GetOptions{})
	if err != nil {
		logrus.Error(err)
		return errors.New("secret doesn't exist")
	}

	//delete the secret
	err = s.client.CoreV1().Secrets(user).Delete(context.Background(), key, metav1.DeleteOptions{})
	if err != nil {
		logrus.Error(err)
		return errors.New("can't delete secret")
	}

	return nil
}

// View Secret
func (s Service) ViewSecret(user string, key string) (map[string][]byte, error) {
	//check if secret with the name key exists
	secret, err := s.client.CoreV1().Secrets(user).Get(context.Background(), key, metav1.GetOptions{})
	if err != nil {
		logrus.Error(err)
		return nil, errors.New("secret doesn't exist")
	}

	return secret.Data, nil
}
