package config

import (
	"fmt"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type BuildParams struct {
	Branch   string
	Upstream string
}

type Config struct {
	Pipeline *Pipeline
	Matrix   *Matrix
}

func ParseConfig(config []byte) (*Config, error) {
	c := Config{}

	err := yaml.Unmarshal(config, &c)
	if err != nil {
		return nil, err
	}

	if c.Pipeline == nil {
		return nil, fmt.Errorf("Pipeline not defined")
	}

	for _, step := range c.Pipeline.Steps {
		step.images, err = parseImages(step.Image)
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

type Pipeline struct {
	Env      []map[string]string
	Services []*Service
	When     *Condition
	Steps    []*CStep
}

type Service struct {
	Name  string
	Image string
}

type Condition struct {
	Branch string
	Env    map[string]string
}

// Valid checks if a condition is valid based on task and build params
// parameters.
func (c *Condition) Valid(task *Task, params *BuildParams) bool {
	if c.Branch != "" {
		if params.Branch != c.Branch {
			return false
		}
	}

	for k, v := range c.Env {
		if e, ok := task.Env[k]; !ok || e != v {
			return false
		}
	}

	return true
}

type CStep struct {
	Name     string
	Image    string
	Env      []map[string]string
	Cmds     []string
	Services []*Service
	When     Condition
	images   []string
}

func (s *CStep) Images() []string {
	return s.images
}

type Matrix struct {
	Exclude []*CStep
}

var imagePatt = regexp.MustCompile(`^([\w/.]+):({([\w.]+)(\s*,\s*([\w.]+))*}|[\w]+)$`)

// parseImages parses a container image definition of the format:
// imagename:{tag1,tag2} or imagename:tag.
func parseImages(image string) ([]string, error) {
	match := imagePatt.FindStringSubmatch(image)
	if len(match) == 0 {
		return nil, fmt.Errorf("failed to parse image: %s", image)
	}

	name := match[1]
	tag := match[2]
	tags := []string{tag}

	if strings.HasPrefix(tag, "{") && strings.HasSuffix(tag, "}") {
		tags = strings.Split(tag[1:len(tag)-1], ",")
	}

	images := make([]string, 0, len(tags))
	for _, tag := range tags {
		img := fmt.Sprintf("%s:%s", name, tag)
		images = append(images, img)
	}

	return images, nil
}
