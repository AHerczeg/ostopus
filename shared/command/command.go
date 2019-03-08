package command

type Command struct {
	Name string `json:"name"`
	Body string `json:"name,omitempty"`
}
