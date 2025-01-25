package discord

var (
	defaultCommands = []Command{
		{
			Name:              "test",
			Description:       "hello biatch!!!",
			Type:              1,
			Integration_types: []int{0, 1},
			Contexts:          []int{0, 1, 2},
		},
	}
)

type Command struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	//TODO: int64?
	Type              int   `json:"type"`
	Integration_types []int `json:"integration_types"`
	Contexts          []int `json:"contexts"`
}

func InstallDefualtCommands(client *Client) {
	client.InstallCommands(defaultCommands)
}
