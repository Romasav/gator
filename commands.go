package main

import "fmt"

type commands struct {
	handlers map[string]func(*state, command) error
}

func newCommands() *commands {
	return &commands{
		handlers: map[string]func(*state, command) error{},
	}
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.Name]
	if !exists {
		return fmt.Errorf("command with a name %v dosent exists", cmd.Name)
	}
	err := handler(s, cmd)
	if err != nil {
		return fmt.Errorf("filed to run the command: %w", err)
	}
	return nil
}

type command struct {
	Name      string
	Arguments []string
}

func newCommand(name string, args []string) *command {
	return &command{
		Name:      name,
		Arguments: args,
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("login requires exactly 1 argument (username), found %v arguments", cmd.Arguments)
	}
	username := cmd.Arguments[0]

	err := s.config.SetUpUser(username)
	if err != nil {
		return fmt.Errorf("failed to set up user: %w", err)
	}
	fmt.Println("The user has been set!")
	return nil
}
