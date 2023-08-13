package commands

type TextCommandManager struct {
	commandSlice []TextCommand
	commandMap   map[string]TextCommand
}

func NewTextCommandManager() *TextCommandManager {
	return &TextCommandManager{
		commandSlice: make([]TextCommand, 0),
		commandMap:   make(map[string]TextCommand),
	}
}

func (m *TextCommandManager) Register(command TextCommand) {
	m.commandMap[command.Name] = command
	for _, alias := range command.Aliases {
		m.commandMap[alias] = command
	}
	m.commandSlice = append(m.commandSlice, command)
}

func (m *TextCommandManager) Get(name string) (TextCommand, bool) {
	command, ok := m.commandMap[name]
	return command, ok
}

func (m *TextCommandManager) GetAll() []TextCommand {
	return m.commandSlice
}
