package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	Token   string
	OwnerID string // Reemplaza con tu ID de usuario de Discord
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se encontró el archivo .env. Se usarán variables de entorno del sistema.")
	}

	// Obtener el token del bot de la variable de entorno
	Token = os.Getenv("DISCORD_TOKEN")
	if Token == "" {
		log.Fatal("Error: DISCORD_TOKEN no está definido en las variables de entorno.")
		return
	}
	OwnerID = os.Getenv("1422676828161703956")
	if OwnerID == "" {
		log.Println("Advertencia: OWNER_ID no está definido, algunas funciones estarán restringidas.")
	}

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error al crear la sesión de Discord:", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify = discordgo.Identify{
		Token: Token,
		Intents: discordgo.IntentsGuilds |
			discordgo.IntentsGuildMessages |
			discordgo.IntentsDirectMessages,
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("Error al abrir la conexión:", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefix := "!"

	if strings.HasPrefix(m.Content, prefix) {
		command := strings.TrimPrefix(m.Content, prefix)
		command = strings.SplitN(command, " ", 2)[0] // Obtener solo el primer comando
		args := strings.SplitN(m.Content, " ", 3)

		switch command {
		case "ayuda":
			ayudaCommand(s, m)
		case "methods":
			methodsCommand(s, m)
		case "adduser":
			addUserCommand(s, m, args)
		case "deleteuser":
			deleteUserCommand(s, m, args)
		case "ataque":
			ataqueCommand(s, m, args)
		case "bots":
			botsCommand(s, m)
		default:
			s.ChannelMessageSend(m.ChannelID, "Comando desconocido. Usa !ayuda para ver la lista de comandos.")
		}
	}
}

