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
	selectCharQuery := `SELECT character_id, name, player_id, precision, strength, endurance, agility, hitpoints, is_alive, is_occupied
	 FROM character
	 WHERE player_id=`

	charRows, err := database.DB.Query(fmt.Sprintln(selectCharQuery, m.Author.ID))
	if err != nil {
		log.Println(err)
	}

	defer charRows.Close()

	// Show the different characters and their stats
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintln("Your characters are: "))
	for charRows.Next() {

		// Check if there is at least one character

		createdCharacter := game.Character{}

		if err := charRows.Scan(&createdCharacter.Id, &createdCharacter.Name, &createdCharacter.PlayerId, &createdCharacter.Precision, &createdCharacter.Strength, &createdCharacter.Endurance, &createdCharacter.Agility, &createdCharacter.Hitpoints, &createdCharacter.IsAlive, &createdCharacter.IsOccupied); err != nil {

			log.Println(err)

		}

		// Send a embed message for each character showing their informations
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: fmt.Sprintln("Name: **", createdCharacter.Name, "**"),
			Description: fmt.Sprintln(
				"**Precision:** ", strconv.Itoa(createdCharacter.Precision),
				"\n**Strength:** ", strconv.Itoa(createdCharacter.Strength),
				"\n**Endurance:** ", strconv.Itoa(createdCharacter.Endurance),
				"\n**Agility:** ", strconv.Itoa(createdCharacter.Agility),
				"\n**Hitpoints:** ", strconv.Itoa(createdCharacter.Hitpoints),
				"\n**Still alive:** ", strconv.FormatBool(createdCharacter.IsAlive),
				"\n**Is in a quest:** ", strconv.FormatBool(createdCharacter.IsOccupied)),
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
