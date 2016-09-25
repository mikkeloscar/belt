package config

import "testing"

// Test parsing valid configuration.
func TestParseConfig(t *testing.T) {
	config := []byte(`
pipeline:
  steps:
  - image: container:latest
  - image: container:{latest,1.0}`)

	_, err := ParseConfig(config)
	if err != nil {
		t.Errorf("should not fail, got %s", err)
	}
}

// Test parsing invalid configuration.
func TestParseConfigInvalid(t *testing.T) {
	config := []byte(`no_pipeline:
matrix:`)

	_, err := ParseConfig(config)
	if err == nil {
		t.Errorf("expected error")
	}

	// test invalid yaml
	config = []byte(`pipeline`)
	_, err = ParseConfig(config)
	if err == nil {
		t.Errorf("expected error")
	}

	// test config with invalid image definition
	config = []byte(`
pipeline:
  steps:
  - image: container:{latest`)

	_, err = ParseConfig(config)
	if err == nil {
		t.Errorf("expected error")
	}
}

// Test parsing container image definition.
func TestParseImages(t *testing.T) {
	images := map[string][]string{
		"image:tag": []string{
			"image:tag",
		},
		"image:{tag}": []string{
			"image:tag",
		},
		"image:{tag1,tag2}": []string{
			"image:tag1",
			"image:tag2",
		},
		"image:{invalidTag": nil,
		"invalid:image:tag": nil,
	}

	for img, result := range images {
		r, err := parseImages(img)
		if err != nil && result != nil {
			t.Errorf("Did not expect image parsing to fail, got: %s", err)
		}

		for i, img := range r {
			if result[i] != img {
				t.Errorf("expected image %s, got %s", result[i], img)
			}
		}
	}
}

// Test valid condition.
func TestValidCondition(t *testing.T) {
	conditions := []struct {
		condition Condition
		task      *Task
		params    *BuildParams
		valid     bool
	}{{
		condition: Condition{
			Branch: "master",
			Env:    map[string]string{"ENV_A": "foo"},
		},
		task: &Task{
			Env: map[string]string{"ENV_A": "foo"},
		},
		params: &BuildParams{
			Branch: "master",
		},
		valid: true,
	}, {
		condition: Condition{
			Branch: "master",
			Env:    map[string]string{"ENV_A": "foo"},
		},
		task: &Task{
			Env: map[string]string{"ENV_B": "foo"},
		},
		params: &BuildParams{
			Branch: "master",
		},
		valid: false,
	}, {
		condition: Condition{
			Branch: "release",
			Env:    map[string]string{"ENV_A": "foo"},
		},
		task: &Task{
			Env: map[string]string{"ENV_A": "foo"},
		},
		params: &BuildParams{
			Branch: "master",
		},
		valid: false,
	}}

	for _, c := range conditions {
		valid := c.condition.Valid(c.task, c.params)
		if c.valid != valid {
			t.Errorf("expected condition to be %t, got %t", c.valid, valid)
		}
	}

}
