package handlers

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/SteakBarbare/RPGBot/connector"
	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/game"
	"github.com/SteakBarbare/RPGBot/utils"
	"github.com/bwmarrin/discordgo"
)

type NewCharacterCommand struct {
	Connector connector.DiscordInterface
	Session		*discordgo.Session
	Message   *discordgo.MessageCreate
}

func NewCharacter(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != "-quit" {

		// Get the character from db
		checkChar := database.DB.QueryRow("SELECT charName FROM characters WHERE player=$1 AND charName=$2;", m.Author.ID, m.Content)
		var foundCharName string
		switch err := checkChar.Scan(&foundCharName); err {
		case sql.ErrNoRows:

		case nil:
			s.ChannelMessageSend(m.ChannelID, "A character with the same name already exists for this account, please choose another name")
			s.AddHandlerOnce(NewCharacter)
			return
		}

		character := statsGeneration(m.Content, m.Author.ID)

		// Show the new character stats & name
		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: fmt.Sprintln("This is your character, **", m.Content, "** !\n Here are it's starting stats:"),
			Description: fmt.Sprintln(
				"**Precision:** ", strconv.Itoa(character.Precision),
				"\n**Strength:** ", strconv.Itoa(character.Strength),
				"\n**Toughness:** ", strconv.Itoa(character.Toughness),
				"\n**Agility:** ", strconv.Itoa(character.Agility),
				"\n**Hitpoints:** ", strconv.Itoa(character.Hitpoints)),
			Color: 0x00ff99,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Player: " + m.Author.ID,
			},
		})

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot error", "Error showing characters."))
			return
		}

		_, err = database.DB.Exec(`INSERT INTO characters(charName, player, precision, strength, toughness, agility, hitpoints, isCharAlive) VALUES ($1, $2, $3, $4, $5, $6, $7, 'true')`,
			m.Content, m.Author.ID, character.Precision, character.Strength, character.Toughness, character.Agility, character.Hitpoints)

		if err != nil {
			panic(err)
		}

		_, err = database.DB.Exec(`INSERT INTO battleChars(charName, player, precision, strength, toughness, agility, hitpoints, isFighting, isDodging, isFleeing) VALUES ($1, $2, $3, $4, $5, $6, $7, 'false', 'false', 'false')`,
			m.Content, m.Author.ID, character.Precision, character.Strength, character.Toughness, character.Agility, character.Hitpoints)

		if err != nil {
			panic(err)
		}

	} else {
		s.ChannelMessageSend(m.ChannelID, "Aborting character creation")
	}
}

// Generate the different stats (random)
func statsGeneration(givenName string, author string) *game.PlayerChar {
	character := game.PlayerChar{
		Name:          givenName,
		Player:        author,
		Precision:   (rand.Intn(20) + 20),
		Strength:      (rand.Intn(20) + 20),
		Toughness:     (rand.Intn(20) + 20),
		Agility:       (rand.Intn(20) + 20),
		Hitpoints:     (rand.Intn(7) + 8)}

	return &character
}
