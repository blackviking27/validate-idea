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
	stateDone
)

// Data structure to hold the user's answers
type ValidationConfig struct {
	Idea     string
	Provider string
	Model    string
}

// Data definitions for the menus
var aiProviders = []string{"Gemini", "OpenAI", "Ollama"}
var defaultModels = map[string]string{
	"Gemini": "gemini-1.5-flash",
	"OpenAI": "gpt-4o-mini",
	"Ollama": "llama3",
}
var modelsList = map[string][]string{
	"Gemini": {"gemini-1.5-flash", "gemini-1.5-pro", "gemini-1.0-pro"},
	"OpenAI": {"gpt-4o-mini", "gpt-4o", "gpt-3.5-turbo"},
	"Ollama": {"llama3", "mistral", "phi3"},
}

// Styling
var (
	titleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("62")).Bold(true).MarginBottom(1)
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	hintStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)
)

type tuiModel struct {
	state     formState
	ideaInput textinput.Model
	cursor    int
	config    ValidationConfig
	quitting  bool
}

// Initialize the form
func initialModel() tuiModel {
	ti := textinput.New()
	ti.Placeholder = "e.g., A peer-to-peer motorcycle rental app..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60

	return tuiModel{
		state:     stateIdea,
		ideaInput: ti,
		cursor:    0,
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

		// State Machine Logic
		switch m.state {
		case stateIdea:
			if msg.String() == "enter" {
				if strings.TrimSpace(m.ideaInput.Value()) != "" {
					m.config.Idea = m.ideaInput.Value()
					m.state = stateProvider
					m.cursor = 0 // Reset cursor for next menu
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

				// Pre-select the default model index
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
				m.config.Model = currentModels[m.cursor]
				m.state = stateDone
				m.quitting = true
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m tuiModel) View() string {
	if m.quitting && m.state == stateDone {
		return "" // Clear the screen when returning to main execution
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
			// Show default model next to provider
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
	}

	return "\n" + s.String() + "\n"
}

// RunForm is the public function to launch the TUI
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
