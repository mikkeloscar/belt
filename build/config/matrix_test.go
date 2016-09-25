package config

import "testing"

// Test building a matrix from a pipeline config.
func TestBuildMatrix(t *testing.T) {
	pipelines := []struct {
		pipeline    *Pipeline
		buildParams *BuildParams
		tasks       int
	}{{
		&Pipeline{
			Steps: []*CStep{
				&CStep{
					images: []string{"image:latest", "image:1.0"},
				},
			},
		},
		nil,
		2,
	}, {
		&Pipeline{
			Env: []map[string]string{
				map[string]string{
					"ENVA": "foo",
				},
				map[string]string{
					"ENVB": "foo",
				},
			},
			Steps: []*CStep{
				&CStep{
					images: []string{"image:latest", "image:1.0"},
				},
			},
		},
		nil,
		4,
	}}

	for _, p := range pipelines {
		tasks := buildMatrix(p.pipeline, p.buildParams)
		if len(tasks) != p.tasks {
			t.Errorf("expected %d steps, got %d", p.tasks, len(tasks))
		}
	}
}

// Test computing step matrix.
func TestComputeStepMatrix(t *testing.T) {
	stepMatrices := []struct {
		steps       []*CStep
		task        *Task
		buildParams *BuildParams
		numSteps    int
	}{{
		[]*CStep{
			&CStep{
				images: []string{"image:latest"},
			},
		},
		nil,
		nil,
		1,
	}, {
		[]*CStep{
			&CStep{
				images: []string{"image:latest"},
				When: Condition{
					Branch: "release",
				},
			},
		},
		nil,
		&BuildParams{Branch: "master"},
		0,
	}, {
		[]*CStep{
			&CStep{
				images: []string{"image:latest"},
			},
			&CStep{
				images: []string{"image:1.0"},
			},
		},
		nil,
		nil,
		1,
	}, {
		[]*CStep{
			&CStep{
				images: []string{"image:latest", "image:1.0"},
			},
			&CStep{
				images: []string{"image:1.0"},
			},
		},
		nil,
		nil,
		2,
	}}

	for _, m := range stepMatrices {
		steps := computeStepMatrix(m.steps, m.task, m.buildParams)
		if len(steps) != m.numSteps {
			t.Errorf("expected %d steps, got %d", m.numSteps, len(steps))
		}
	}
}

// Test computing steps.
func TestComputeSteps(t *testing.T) {
	csteps := []struct {
		step  *CStep
		steps int
	}{{
		&CStep{
			images: []string{"image:latest"},
			Env:    nil,
		},
		1,
	}, {
		&CStep{
			images: []string{"image:latest"},
			Env: []map[string]string{
				map[string]string{
					"ENVA": "foo",
				},
			},
		},
		1,
	}, {
		&CStep{
			images: []string{"image:latest"},
			Env: []map[string]string{
				map[string]string{
					"ENVA": "foo",
					"ENVB": "bar",
				},
			},
		},
		1,
	}, {
		&CStep{
			images: []string{"image:latest"},
			Env: []map[string]string{
				map[string]string{
					"ENVA": "foo",
				},
				map[string]string{
					"ENVB": "foo",
				},
			},
		},
		2,
	}, {
		&CStep{
			images: []string{"image:latest", "image:1.6"},
			Env: []map[string]string{
				map[string]string{
					"ENVA": "foo",
				},
			},
		},
		2,
	}, {
		&CStep{
			images: []string{"image:latest", "image:1.6"},
			Env: []map[string]string{
				map[string]string{
					"ENVA": "foo",
				},
				map[string]string{
					"ENVB": "bar",
				},
			},
		},
		4,
	}}

	for _, s := range csteps {
		steps := computeSteps(s.step)
		if len(steps) != s.steps {
			t.Errorf("expected %d steps, got %d", s.steps, len(steps))
		}
	}
}
