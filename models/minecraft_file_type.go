package models

type MinecraftFileType int

const (
	Library  MinecraftFileType = 1
	Mod      MinecraftFileType = 2
	Assets   MinecraftFileType = 3
	Config   MinecraftFileType = 4
	World    MinecraftFileType = 5
	Data     MinecraftFileType = 6
	Resource MinecraftFileType = 7
	Texture  MinecraftFileType = 8
	Other    MinecraftFileType = 9
)

func (mineFileType MinecraftFileType) GetTypeName() string {
	switch mineFileType {
	case Library:
		return "library"
	case Mod:
		return "mod"
	case Config:
		return "config"
	case Assets:
		return "assets"
	case World:
		return "world"
	case Data:
		return "data"
	case Resource:
		return "resource"
	case Texture:
		return "texture"
	default:
		return "other"
	}
}

func ParseMinecraftFileType(mineFileType string) MinecraftFileType {
	switch mineFileType {
	case "library":
		return Library
	case "mod":
		return Mod
	case "config":
		return Config
	case "assets":
		return Assets
	case "world":
		return World
	case "data":
		return Data
	case "resource":
		return Resource
	case "texture":
		return Texture
	default:
		return Other
	}
}
