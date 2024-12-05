package discord

//https://discohook.org/

//discordWebhookUrl = "https://discord.com/api/webhooks/1075323981193826354/rCJrCgDxYIV3E-gpuhh6F8zh8smCnev9Tguil9flnMaI2fVMNTwbp2fYEh0yAwcWsDIX"
//discordRobotThread = "1161959914185429053"

type Message struct {
	DiscordWebhookUrl  string
	DiscordRobotThread string

	UserName string

	AuthorName string
	Title      string
	Url        string
	Message    string
	FieldMap   map[string]string
	Content    string

	DiscordFile *File
}

type File struct {
	FilePath string
	Title    string
	Desc     string
}
