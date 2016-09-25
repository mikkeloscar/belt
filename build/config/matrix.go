package config

type Build struct {
	Matrix []*Task
}

type Task struct {
	Env      map[string]string
	Services []*Service
	Steps    []*Step
}

type Step struct {
	Name     string
	Image    string
	Env      map[string]string
	Cmds     []string
	Services []*Service
}

func buildMatrix(pipeline *Pipeline, buildParams *BuildParams) []*Task {
	var initTasks []*Task

	if len(pipeline.Env) == 0 {
		task := &Task{
			Services: pipeline.Services,
		}
		initTasks = append(initTasks, task)
	} else {
		for _, e := range pipeline.Env {
			task := &Task{
				Env:      e,
				Services: pipeline.Services,
			}
			initTasks = append(initTasks, task)
			// for _, os := range pipeline.Os {
			// 	task := &Task{}
			// }
		}
	}

	tasks := make([]*Task, 0, len(initTasks))
	for _, t := range initTasks {
		stepMatrix := computeStepMatrix(pipeline.Steps, t, buildParams)
		for _, steps := range stepMatrix {
			task := &Task{
				Env:      t.Env,
				Services: t.Services,
				Steps:    steps,
			}
			tasks = append(tasks, task)
		}
	}

	return tasks
}

func computeStepMatrix(steps []*CStep, task *Task, buildParams *BuildParams) [][]*Step {
	var stepMatrix [][]*Step
	for _, step := range steps {
		if step.When.Valid(task, buildParams) {
			steps := computeSteps(step)
			newStepMatrix := make([][]*Step, 0, len(stepMatrix)*len(steps))
			for _, s := range steps {
				if len(stepMatrix) == 0 {
					newStepMatrix = append(newStepMatrix, []*Step{s})
					continue
				}

				for _, matrixSteps := range stepMatrix {
					newSteps := append(matrixSteps, s)
					newStepMatrix = append(newStepMatrix, newSteps)
				}
			}
			stepMatrix = newStepMatrix
		}
	}

	return stepMatrix
}

func computeSteps(step *CStep) []*Step {
	var stepMatrix []*Step

	for _, image := range step.Images() {
		// if there are no environment variables we act as if there
		// were one to get a len(images) * 1 matrix instead of a
		// len(images) * 0 matrix
		if len(step.Env) == 0 {
			s := &Step{
				Name:     step.Name,
				Image:    image,
				Cmds:     step.Cmds,
				Services: step.Services,
			}
			stepMatrix = append(stepMatrix, s)
			continue
		}

		for _, env := range step.Env {
			s := &Step{
				Name:     step.Name,
				Image:    image,
				Env:      env,
				Cmds:     step.Cmds,
				Services: step.Services,
			}
			stepMatrix = append(stepMatrix, s)
		}
	}

	return stepMatrix
}
