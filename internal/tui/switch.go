package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type SwitchView struct {
	model       model
	pendingCh   chan string
	succeededCh chan string
}

type step struct {
	name         string
	pendingMsg   string
	completedMsg string
	completed    bool
}

type stage struct {
	completedMsg string
	completed    bool
	steps        []step
}

type model struct {
	parentView      *SwitchView
	progress        progress.Model
	spinner         spinner.Model
	stages          []stage
	currentStep     step
	completedStages []stage
	completedSteps  []step
	percentage      float64
	completedMsg    string
	completed       bool
}

type SwitchModelStepDto struct {
	Name         string
	PendingMsg   string
	CompletedMsg string
}

type SwitchModelStageDto struct {
	CompletedMsg string
	Steps        []SwitchModelStepDto
}

type pendingMsg string
type succeededMsg string
type completedMsg bool
type quitIfCompletedMsg bool

func (sv *SwitchView) pendingCmd() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			return pendingMsg(<-sv.pendingCh)
		},
		func() tea.Msg {
			return succeededMsg(<-sv.succeededCh)
		},
	)
}

func (sv *SwitchView) completedCmd() tea.Cmd {
	return func() tea.Msg {
		return completedMsg(true)
	}
}

func (sv *SwitchView) quitIfCompletedCmd() tea.Cmd {
	return func() tea.Msg {
		return quitIfCompletedMsg(true)
	}
}

func initModel(stages []SwitchModelStageDto, parentView *SwitchView) model {
	prog := progress.New(progress.WithSolidFill(""))
	spnr := spinner.New()
	spnr.Spinner = spinner.Dot

	modelStages := []stage{}

	for _, s := range stages {
		modelSteps := []step{}

		for _, st := range s.Steps {
			modelSteps = append(modelSteps, step{
				name:         st.Name,
				pendingMsg:   st.PendingMsg,
				completedMsg: GetCompletedMsg(st.CompletedMsg),
			})
		}

		modelStages = append(modelStages, stage{
			completedMsg: GetCompletedMsg(s.CompletedMsg),
			steps:        modelSteps,
		})
	}

	return model{
		parentView:      parentView,
		progress:        prog,
		spinner:         spnr,
		completed:       false,
		stages:          modelStages,
		percentage:      0,
		completedMsg:    GetCompletedMsg("Theme switched"),
		currentStep:     step{},
		completedStages: []stage{},
		completedSteps:  []step{},
	}
}

func CreateSwitchView(stages []SwitchModelStageDto) *SwitchView {
	view := &SwitchView{
		model:       model{},
		pendingCh:   make(chan string),
		succeededCh: make(chan string),
	}

	view.model = initModel(stages, view)

	return view
}

func (sv *SwitchView) Launch(launchedCh chan<- bool) error {
	launchedCh <- true
	if _, err := tea.NewProgram(sv.model).Run(); err != nil {
		return err
	}

	return nil
}

func (sv *SwitchView) Pending(stepName string) error {
	sv.pendingCh <- stepName
	return nil
}

func (sv *SwitchView) Proceed(stepName string) error {
	sv.succeededCh <- stepName
	return nil
}

func (sv *SwitchView) completedPercentage(s []stage) float64 {
	completedSteps := 0
	allSteps := 0

	for _, stage := range s {
		allSteps += len(stage.steps)
		for _, step := range stage.steps {
			if !step.completed {
				break
			}
			completedSteps += 1
		}
	}

	return float64(completedSteps) / float64(allSteps)
}

func (sv *SwitchView) getStep(name string, s []stage) step {
	for _, stage := range s {
		for _, step := range stage.steps {
			if step.name == name {
				return step
			}
		}
	}

	return step{}
}

func (sv *SwitchView) completeStep(name string, s []stage) []stage {
	for stageIdx, stage := range s {
		for stepIdx, step := range stage.steps {
			if step.name == name {
				s[stageIdx].steps[stepIdx].completed = true
				return s
			}
		}
	}

	return s
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.parentView.pendingCmd(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case pendingMsg:
		ms := msg.(pendingMsg)
		m.currentStep = m.parentView.getStep(string(ms), m.stages)
		return m, nil

	case succeededMsg:
		ms := msg.(succeededMsg)
		m.stages = m.parentView.completeStep(string(ms), m.stages)
		m.completedSteps = append(m.completedSteps, m.parentView.getStep(string(ms), m.stages))

		for idx, s := range m.stages {
			if s.completed {
				continue
			}

			allCompleted := true
			for _, st := range s.steps {
				if !st.completed {
					allCompleted = false
					break
				}
			}

			if allCompleted {
				m.stages[idx].completed = true
				m.completedStages = append(m.completedStages, m.stages[idx])
			}
		}

		m.percentage = m.parentView.completedPercentage(m.stages)
		cmd := m.progress.SetPercent(m.percentage)
		if m.percentage == 1.00 {
			return m, tea.Batch(cmd, m.parentView.completedCmd())
		}

		return m, tea.Batch(cmd, m.parentView.pendingCmd())

	case completedMsg:
		m.completed = true

		if !m.progress.IsAnimating() {
			return m, tea.Quit
		}

		return m, nil

	case quitIfCompletedMsg:
		if !m.progress.IsAnimating() {
			return m, tea.Quit
		}

		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		if m.completed {
			return m, tea.Batch(cmd, m.parentView.quitIfCompletedCmd())
		}

		return m, cmd

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)

		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	start := "\n"
	end := "\n\n"
	spinnerLine := fmt.Sprintf("%s %s\n", m.spinner.View(), m.currentStep.pendingMsg)
	progressLine := m.progress.View()

	completedStepsLines := ""
	if len(m.completedSteps) > 0 {
		for _, step := range m.completedSteps {
			completedStepsLines += step.completedMsg + "\n"
		}
	}

	completedStagesLines := ""
	if len(m.completedStages) > 0 {
		for _, stage := range m.completedStages {
			completedStagesLines += stage.completedMsg + "\n"
		}
	}

	if m.completed {
		return fmt.Sprintf("%s%s%s%s\n\n%s%s", start, completedStepsLines, completedStagesLines, m.completedMsg, progressLine, end)
	}

	return fmt.Sprintf("%s%s%s%s\n%s%s", start, completedStepsLines, completedStagesLines, spinnerLine, progressLine, end)
}
