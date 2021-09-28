package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/SteakBarbare/RPGBot/game"
	"github.com/SteakBarbare/RPGBot/utils"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/bwmarrin/discordgo"
)

func ShowCharacters(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Get the characters info from db
	charRows, err := database.DB.Query(fmt.Sprintln("SELECT * FROM Characters WHERE player=", m.Author.ID))
	if err != nil {
		log.Fatal(err)
	}

	defer charRows.Close()

	// Show the different characters and their stats
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintln("Your characters are: "))
	for charRows.Next() {

		// Check if there is at least one character

		createdCharacter := game.PlayerChar{}

		if err := charRows.Scan(&createdCharacter.Id, &createdCharacter.Name, &createdCharacter.Player, &createdCharacter.Precision, &createdCharacter.Strength, &createdCharacter.Toughness, &createdCharacter.Agility, &createdCharacter.Hitpoints, &createdCharacter.IsCharAlive); err != nil {

			log.Fatal(err)

		}

		// Send a embed message for each character showing their informations
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: fmt.Sprintln("Name: **", createdCharacter.Name, "**"),
			Description: fmt.Sprintln(
				"**Precision:** ", strconv.Itoa(createdCharacter.Precision),
				"\n**Strength:** ", strconv.Itoa(createdCharacter.Strength),
				"\n**Toughness:** ", strconv.Itoa(createdCharacter.Toughness),
				"\n**Agility:** ", strconv.Itoa(createdCharacter.Agility),
				"\n**Hitpoints:** ", strconv.Itoa(createdCharacter.Hitpoints),
				"\n**Still Alive:** ", createdCharacter.IsCharAlive),
			Color: 0x0099ff,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Player: " + m.Author.ID,
			},
		})

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot error", "Error showing characters."))
			return
		}

	}

}
