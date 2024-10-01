package repository

import (
	"fmt"
	"io/fs"
	"strings"

	"gopkg.in/yaml.v3"
)

type Resource struct {
	group     string
	version   string
	kind      string
	namespace string
	name      string
	content   map[string]any
}

// Group returns the API group of the resource.
func (r Resource) Group() string {
	return r.group
}

// Version returns the API version of the resource.
func (r Resource) Version() string {
	return r.version
}

// Kind returns the kind of the resource.
func (r Resource) Kind() string {
	return r.kind
}

// Namespace returns the namespace of the resource.
func (r Resource) Namespace() string {
	return r.namespace
}

// Name returns the name of the resource.
func (r Resource) Name() string {
	return r.name
}

// Content returns the content of the resource.
// Do not alter the returned map or the resource could become inconsistent.
func (r Resource) Content() map[string]any {
	return r.content
}

// MakeResource creates a new Resource from a map of the resource's content.
func MakeResource(content map[string]any) (Resource, error) {
	apiVersionRaw, ok := content["apiVersion"]

	if !ok {
		return Resource{}, fmt.Errorf("missing apiVersion")
	}

	apiVersion, ok := apiVersionRaw.(string)

	if !ok {
		return Resource{}, fmt.Errorf("apiVersion '%#v' must be a string but it is %T", apiVersionRaw, apiVersionRaw)
	}

	var group, version string

	if apiVersion == "v1" {
		group = "core"
		version = "v1"
	} else {
		var found bool
		group, version, found = strings.Cut(apiVersion, "/")

		if !found {
			return Resource{}, fmt.Errorf("missing '/' separator in apiVersion'%s'", apiVersion)
		}
		if len(group) == 0 {
			return Resource{}, fmt.Errorf("missing group in apiVersion '%s'", apiVersion)
		}
		if len(version) == 0 {
			return Resource{}, fmt.Errorf("missing version in apiVersion '%s'", apiVersion)
		}
	}

	kindRaw, ok := content["kind"]

	if !ok {
		return Resource{}, fmt.Errorf("missing kind")
	}

	kind, ok := kindRaw.(string)

	if !ok {
		return Resource{}, fmt.Errorf("kind '%#v' must be a string but it is %T", kindRaw, kindRaw)
	}

	metadataRaw, ok := content["metadata"]

	if !ok {
		return Resource{}, fmt.Errorf("missing metadata")
	}

	metadata, ok := metadataRaw.(map[string]any)

	if !ok {
		return Resource{}, fmt.Errorf("metadata '%#v' must be a map but it is %T", metadataRaw, metadataRaw)
	}

	namespaceRaw, ok := metadata["namespace"]

	if !ok {
		namespaceRaw = ""
	}

	namespace, ok := namespaceRaw.(string)

	if !ok {
		return Resource{}, fmt.Errorf("metadata.namespace '%#v' must be a string but it is %T", namespaceRaw, namespaceRaw)
	}

	nameRaw, ok := metadata["name"]

	if !ok {
		return Resource{}, fmt.Errorf("missing name")
	}

	name, ok := nameRaw.(string)

	if !ok {
		return Resource{}, fmt.Errorf("metadata.name '%#v' must be a string but it is %T", nameRaw, nameRaw)
	}

	res := Resource{
		group:     group,
		version:   version,
		kind:      kind,
		namespace: namespace,
		name:      name,
		content:   content,
	}

	return res, nil
}

type Repository interface {
	Query(query string) (ResultSet, error)
	Close() error
	Upsert(resource Resource) error
}

type ResultSet struct {
	Names   []string
	Results []map[string]*Resource
}

// LoadYAMLFromDirectoryTree loads all YAML files from a directory tree into a repository.
func LoadYAMLFromDirectoryTree(filesystem fs.FS, directory string, repo Repository) error {
	return fs.WalkDir(filesystem, directory, func(path string, d fs.DirEntry, pathErr error) error {
		if pathErr != nil {
			return pathErr
		}

		if d.IsDir() {
			return nil
		}

		data, err := fs.ReadFile(filesystem, path)
		if err != nil {
			return fmt.Errorf("file '%s': %w", path, err)
		}

		parsed := make(map[string]any)

		err = yaml.Unmarshal(data, &parsed)
		if err != nil {
			return fmt.Errorf("file '%s': %w", path, err)
		}

		res, err := MakeResource(parsed)
		if err != nil {
			return fmt.Errorf("file '%s': %w", path, err)
		}

		return repo.Upsert(res)
	})
}
