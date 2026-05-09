package main

import "context"

type App struct {
	connections []connect_struct
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.connections = get_connection_data()
}
func (a *App) GetConnections() []ConnectDTO {
	internal := get_connection_data()
	out := make([]ConnectDTO, 0, len(internal))

	for _, c := range internal {
		out = append(out, toConnectDTO(c))
	}
	return out
}

func (a *App) GetLinkData(id int) LinkDataDTO {
	for _, c := range get_connection_data() {
		if c.id == id {
			return toLinkDataDTO(get_link_data(&c))
		}
	}
	return LinkDataDTO{}
}
