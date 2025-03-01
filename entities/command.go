package entities

type CliCommand struct {
	Name        string
	Description string
	Callback    func() error
}
