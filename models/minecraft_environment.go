package models

type MinecraftEnvironment int

const (
	Server MinecraftEnvironment = 1
	Client MinecraftEnvironment = 2
)

func (env MinecraftEnvironment) GetFolderName() string {
	switch env {
	case Client:
		return "client_files"
	default:
		return "server_files"
	}
}

func ParseMinecraftEnvironment(env string) MinecraftEnvironment {
	switch env {
	case "client":
		return Client
	default:
		return Server
	}
}
