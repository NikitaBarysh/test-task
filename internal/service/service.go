package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"test-task/internal/config"
	"test-task/internal/entity"
)

type Service struct {
	client *http.Client
	cfg    *config.Config
}

func NewService(newCfg *config.Config) *Service {
	newClient := &http.Client{
		Timeout: newCfg.ClientTimeOut,
	}
	return &Service{client: newClient, cfg: newCfg}
}

func (s *Service) GetTags(image string) ([]string, error) {
	url := fmt.Sprintf("http://%s:%s/v2/%s/tags/list", s.cfg.DockerRegistryHost, s.cfg.DockerRegistryPort, image)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("err to prepare request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("err to do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var tagsList entity.TagsList
	if err := json.NewDecoder(resp.Body).Decode(&tagsList); err != nil {
		return nil, fmt.Errorf("err to decode response: %w", err)
	}

	return tagsList.Tags, nil
}

func (s *Service) GetManifest(image, tag string) (*entity.Manifest, error) {
	url := fmt.Sprintf("http://%s:%s/v2/%s/manifests/%s", s.cfg.DockerRegistryHost, s.cfg.DockerRegistryPort, image, tag)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("err to prepare request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("err to do request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var manifest entity.Manifest

	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return nil, fmt.Errorf("err to decode response: %w", err)
	}

	return &manifest, nil
}
