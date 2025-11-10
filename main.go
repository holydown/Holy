package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings" // Importante: Importar el paquete "strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Crear una nueva sesión de Discord usando el token del bot
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Error al crear la sesión de Discord: %v", err)
	}

	// Registrar el controlador de eventos messageCreate
	dg.AddHandler(messageCreate)

	// Especificar las intenciones que necesita el bot
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	// Abrir una conexión websocket con Discord
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error al abrir la conexión de Discord: %v", err)
	}

	// Indicar que el bot está funcionando
	fmt.Println("El bot está funcionando. Presiona Ctrl+C para salir.")

	// Esperar una señal de interrupción para cerrar la conexión
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cerrar la conexión de Discord
	dg.Close()
}

// messageCreate
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignorar los mensajes enviados por el propio bot
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Separar el comando y los argumentos
	args := strings.Split(m.Content, " ")
	command := args[0]

	// Manejar los diferentes comandos
	switch command {
	case "!ayuda":
		ayudaCommand(s, m)
	case "!methods":
		methodsCommand(s, m)
	case "!adduser":
		addUserCommand(s, m, args)
	case "!deleteuser":
		deleteUserCommand(s, m, args)
	case "!ataque":
		ataqueCommand(s, m, args)
	case "!bots":
		// Llamada corregida: botsCommand no espera args según tu definición actual
		botsCommand(s, m)
	}
}
