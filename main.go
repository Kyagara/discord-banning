package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	bot, err := discordgo.New("Bot " + "Coloque o Token do bot aqui")

	if err != nil {
		fmt.Println("erro criando sessão,", err)
		return
	}

	bot.AddHandler(handlePresence)

	bot.Identify.Intents = discordgo.IntentsGuildPresences

	err = bot.Open()

	if err != nil {
		fmt.Println("error abrindo conexão,", err)
		return
	}

	fmt.Printf("Bot '%v' ligado, aperte CTRL-C para fechar.\n", bot.State.User.Username)

	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	bot.Close()
}

// Função chamada toda vez que a Presence de um usuário for atualizada
func handlePresence(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	if p.User.Bot {
		return
	}

	// Retornando caso o usuário tenha um ID específico, você pode colocar o seu ID ou deixar do jeito que está para banir qualquer pessoa
	if p.User.ID == "Algum ID de usuário" {
		return
	}

	if len(p.Activities) != 0 {
		if p.Activities[0].Name != "League of Legends" {
			return
		}

		if p.Activities[0].Timestamps.StartTimestamp <= 0 {
			return
		}

		// Tempo que a Presence do Discord foi criada - Timestamp que a partida acabou
		time := p.Activities[0].CreatedAt.UnixMilli() - p.Activities[0].Timestamps.StartTimestamp

		// Se o tempo for menor que 30 minutos (1800000 milisegundos), retornar
		if time < 1800000 {
			return
		}

		// Se passar ou ser igual ao tempo definido, banir o usuário
		s.GuildMemberDelete(p.GuildID, p.User.ID)
	}
}
