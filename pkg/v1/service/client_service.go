package service

import (
	"context"
	"errors"
	"fmt"
	"secret-manager/pkg/v1/models"
	"strconv"
	"unicode"

	"github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/tools/clientcmd"
)

type IUser interface {
	// CreateClient() (*kubernetes.Clientset, error)
	CreateUserSecret(user models.User) error
	LoginUser(user models.User, token string) error
	GetSecretByLabel(token string) (string, error)
	ListUserSecrets(user string, query string, page string, pageSize string) (*corev1.SecretList, int, error)
}

// Register user
func (s Service) CreateUserSecret(user models.User) error {
	//Check if user already exists
	_, err := s.client.CoreV1().Secrets("default").Get(context.Background(), user.Username, metav1.GetOptions{})
	if err == nil {
		logrus.Error(err)
		return errors.New("user already exists")
	}

	//Check if password is valid
	valid := validatePassword(user.Password)
	if !valid {
		return errors.New("password is not valid")
	}

	//Create new Secret with username as it's name
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: user.Username,
		},
		Data: map[string][]byte{
			"password": []byte(user.Password),
		},
	}

	_, err = s.client.CoreV1().Secrets("default").Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		logrus.Error(err)
		return errors.New("can't create user")
	}

	//Create a namespace with the name username
	// usrNamespace, err := clientset.CoreV1().Namespaces().Get(context.Background(), user.Username, metav1.GetOptions{})
	usrNamespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: user.Username,
		},
	}

	_, err = s.client.CoreV1().Namespaces().Create(context.Background(), usrNamespace, metav1.CreateOptions{})
	if err != nil {
		logrus.Error(err)
		return errors.New("can't create namespace")
	}

	// fmt.Println("Created new secret:")

	return nil
}

// Login user
func (s Service) LoginUser(user models.User, token string) error {
	//Check if user exists
	secret, err := s.client.CoreV1().Secrets("default").Get(context.Background(), user.Username, metav1.GetOptions{})
	if err != nil {
		logrus.Error(err)
		return errors.New("user doesn't exist")
	}

	//Check if password is correct
	pass := secret.Data["password"]
	if string(pass) != user.Password {
		return errors.New("incorrect password")
	}

	//Add token to label
	secret.ObjectMeta.Labels = make(map[string]string)
	secret.ObjectMeta.Labels["token"] = token

	// Update the secret with the new label
	_, err = s.client.CoreV1().Secrets("default").Update(context.Background(), secret, metav1.UpdateOptions{})
	if err != nil {
		return errors.New("failed to set token label in secret")
	}

	//Check if the namespace exists
	_, err = s.client.CoreV1().Namespaces().Get(context.Background(), user.Username, metav1.GetOptions{})

	//if the namespace doesn't exist, error
	if err != nil {
		logrus.Error(err)
		return errors.New("namespace doesn't exist")
	}

	return nil
}

// Get secret by Token
func (s Service) GetSecretByLabel(token string) (string, error) {
	labelSelector := fmt.Sprintf("token=%s", token)
	secret, err := s.client.CoreV1().Secrets("default").List(context.Background(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		return "", errors.New("token not found")
	}

	return secret.Items[0].Name, nil
}

const (
	DefaultPageSize = 10
	DefaultPageNum  = 1
	// maxPageSize     = 10
)

// List secrets
func (s Service) ListUserSecrets(user string, query string, page string, pageSize string) (*corev1.SecretList, int, error) {
	//check if namespace exists
	usrNamespace, err := s.client.CoreV1().Namespaces().Get(context.Background(), user, metav1.GetOptions{})
	if err != nil {
		logrus.Error(err)
		return nil, 0, errors.New("namespace doesn't exist")
	}

	//Calculate total number of secrets in this namespace
	var secretList *corev1.SecretList
	var totalSecrets int

	secretList, err = s.client.CoreV1().Secrets(usrNamespace.Name).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		logrus.Error(err)
		return nil, 0, errors.New("failed to list secrets")
	}

	totalSecrets = len(secretList.Items)

	//Convert page and pagesize to integers
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = DefaultPageNum
	}

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum < 1 || pageSizeNum > DefaultPageSize {
		pageSizeNum = DefaultPageSize
	}

	//Prepare list options for pagination
	listOptions := metav1.ListOptions{
		Limit: int64(pageSizeNum),
	}

	//Check if query is provided
	if query != "" {
		listOptions.FieldSelector = fmt.Sprintf("metadata.name=%s", query)
	}

	//Get secrets with specified options
	var lastContinue string

	for i := 1; i <= pageNum; i++ {
		secretList, err = s.client.CoreV1().Secrets(usrNamespace.Name).List(context.Background(), listOptions)
		if err != nil {
			logrus.Error(err)
			return nil, 0, errors.New("failed to list secrets")
		}

		lastContinue = secretList.Continue
		listOptions.Continue = lastContinue
		if lastContinue == "" {
			break
		}
	}

	return secretList, totalSecrets, nil
}

func validatePassword(password string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 6 {
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
