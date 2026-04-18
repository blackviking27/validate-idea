package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Define the steps of our form
type formState int

const (
	stateIdea formState = iota
	stateProvider
	stateModel
	stateCustomModel // 1. Added a new state for typing the custom model
	stateDone
)

type ValidationConfig struct {
	Idea     string
	Provider string
	Model    string
}

// Data definitions for the menus
var aiProviders = []string{"Gemini", "Ollama"}
var defaultModels = map[string]string{
	"Gemini": "gemini-1.5-flash",
	"Ollama": "llama3",
}

// 2. Added "Enter Custom Model..." to the end of each list
var modelsList = map[string][]string{
	"Gemini": {"gemini-1.5-flash", "gemini-1.5-pro", "gemini-1.0-pro", "Enter Custom Model..."},
	"Ollama": {"llama3", "mistral", "phi3", "Enter Custom Model..."},
}

var (
	titleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Bold(true).MarginBottom(1)
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	hintStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
)

type tuiModel struct {
	state      formState
	ideaInput  textinput.Model
	modelInput textinput.Model // 3. Added a second text input for the model
	cursor     int
	config     ValidationConfig
	quitting   bool
}

func initialModel() tuiModel {
	// Setup Idea Input
	ideaTi := textinput.New()
	ideaTi.Placeholder = "e.g., A peer-to-peer motorcycle rental app..."
	ideaTi.Focus()
	ideaTi.CharLimit = 256
	ideaTi.Width = 60

	// Setup Custom Model Input
	modelTi := textinput.New()
	modelTi.Placeholder = "e.g., llama3:8b-instruct-fp16"
	modelTi.CharLimit = 64
	modelTi.Width = 40

	return tuiModel{
		state:      stateIdea,
		ideaInput:  ideaTi,
		modelInput: modelTi,
		cursor:     0,
	}
}

func (m tuiModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}

		switch m.state {
		case stateIdea:
			if msg.String() == "enter" {
				if strings.TrimSpace(m.ideaInput.Value()) != "" {
					m.config.Idea = m.ideaInput.Value()
					m.state = stateProvider
					m.cursor = 0
				}
				return m, nil
			}
			var cmd tea.Cmd
			m.ideaInput, cmd = m.ideaInput.Update(msg)
			return m, cmd

		case stateProvider:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(aiProviders)-1 {
					m.cursor++
				}
			case "enter":
				m.config.Provider = aiProviders[m.cursor]
				m.state = stateModel
				m.cursor = 0

				defModel := defaultModels[m.config.Provider]
				for i, mod := range modelsList[m.config.Provider] {
					if mod == defModel {
						m.cursor = i
						break
					}
				}
			}

		case stateModel:
			currentModels := modelsList[m.config.Provider]
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(currentModels)-1 {
					m.cursor++
				}
			case "enter":
				selection := currentModels[m.cursor]

				// 4. Logic routing: Are we using a preset or going custom?
				if selection == "Enter Custom Model..." {
					m.state = stateCustomModel
					m.modelInput.Focus()
					return m, textinput.Blink
				}

				m.config.Model = selection
				m.state = stateDone
				m.quitting = true
				return m, tea.Quit
			}

		// 5. Handle the new custom text input state
		case stateCustomModel:
			if msg.String() == "enter" {
				if strings.TrimSpace(m.modelInput.Value()) != "" {
					m.config.Model = m.modelInput.Value()
					m.state = stateDone
					m.quitting = true
					return m, tea.Quit
				}
			}
			var cmd tea.Cmd
			m.modelInput, cmd = m.modelInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m tuiModel) View() string {
	if m.quitting && m.state == stateDone {
		return ""
	}
	if m.quitting {
		return "Form aborted.\n"
	}

	var s strings.Builder

	switch m.state {
	case stateIdea:
		s.WriteString(titleStyle.Render("💡 What is your startup idea?"))
		s.WriteString("\n")
		s.WriteString(m.ideaInput.View())
		s.WriteString("\n\n" + hintStyle.Render("(Press Enter to continue)"))

	case stateProvider:
		s.WriteString(titleStyle.Render("🤖 Choose your AI Provider"))
		s.WriteString("\n")
		for i, choice := range aiProviders {
			cursor := "  "
			if m.cursor == i {
				cursor = cursorStyle.Render("❯ ")
			}
			defaultHint := hintStyle.Render(fmt.Sprintf("(Default: %s)", defaultModels[choice]))
			s.WriteString(fmt.Sprintf("%s%s %s\n", cursor, choice, defaultHint))
		}
		s.WriteString("\n" + hintStyle.Render("(Use arrow keys and Enter)"))

	case stateModel:
		s.WriteString(titleStyle.Render(fmt.Sprintf("⚙️  Choose a model for %s", m.config.Provider)))
		s.WriteString("\n")
		for i, choice := range modelsList[m.config.Provider] {
			cursor := "  "
			if m.cursor == i {
				cursor = cursorStyle.Render("❯ ")
			}
			s.WriteString(fmt.Sprintf("%s%s\n", cursor, choice))
		}
		s.WriteString("\n" + hintStyle.Render("(Use arrow keys and Enter)"))

	case stateCustomModel:
		s.WriteString(titleStyle.Render("⌨️  Type the exact model name:"))
		s.WriteString("\n")
		s.WriteString(m.modelInput.View())
		s.WriteString("\n\n" + hintStyle.Render("(Press Enter to save)"))
	}

	return "\n" + s.String() + "\n"
}

func RunForm() (ValidationConfig, error) {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		return ValidationConfig{}, err
	}

	finalModel := m.(tuiModel)
	if finalModel.state != stateDone {
		return ValidationConfig{}, fmt.Errorf("user aborted setup")
	}

	return finalModel.config, nil
}
