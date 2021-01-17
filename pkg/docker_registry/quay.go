package docker_registry

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
)

const QuayImplementationName = "quay"

type QuayNotFoundError apiError

var quayPatterns = []string{"^quay\\.io", "^quay\\..*\\.com"}

type quay struct {
	*defaultImplementation
	quayApi
	quayCredentials
}

type quayOptions struct {
	defaultImplementationOptions
	quayCredentials
}

type quayCredentials struct {
	token string
}

func newQuay(options quayOptions) (*quay, error) {
	d, err := newDefaultImplementation(options.defaultImplementationOptions)
	if err != nil {
		return nil, err
	}

	quay := &quay{
		defaultImplementation: d,
		quayApi:               newQuayApi(),
		quayCredentials:       options.quayCredentials,
	}

	return quay, nil
}

func (r *quay) DeleteRepo(ctx context.Context, reference string) error {
	return r.deleteRepo(ctx, reference)
}

func (r *quay) deleteRepo(ctx context.Context, reference string) error {
	hostname, namespace, repository, err := r.parseReference(reference)
	if err != nil {
		return err
	}

	resp, err := r.quayApi.DeleteRepository(ctx, hostname, namespace, repository, r.quayCredentials.token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return QuayNotFoundError{error: err}
	}

	return nil
}

func (r *quay) String() string {
	return QuayImplementationName
}

func (r *quay) parseReference(reference string) (string, string, string, error) {
	parsedReference, err := name.NewRepository(reference)
	if err != nil {
		return "", "", "", err
	}

	hostname := parsedReference.RegistryStr()
	repositoryStr := parsedReference.RepositoryStr()

	var namespace, repository string
	switch len(strings.Split(repositoryStr, "/")) {
	case 1:
		namespace = repositoryStr
	case 2:
		repository = path.Base(repositoryStr)
		namespace = path.Base(strings.TrimSuffix(repositoryStr, repository))
	default:
		return "", "", "", fmt.Errorf("unexpected reference %s", reference)
	}

	return hostname, namespace, repository, nil
}
