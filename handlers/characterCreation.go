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

		authorId, err := strconv.ParseInt(m.Author.ID, 10, 64)

		if err != nil {
			panic(err)
		}

		// Get the character from db
		checkChar := database.DB.QueryRow("SELECT name FROM character_model WHERE player_id=$1 AND name=$2;", authorId, m.Content)
		
		var foundCharName string

		switch err := checkChar.Scan(&foundCharName); err {
		case sql.ErrNoRows:

		case nil:
			s.ChannelMessageSend(m.ChannelID, "A character with the same name already exists for this account, please choose another name")
			s.AddHandlerOnce(NewCharacter)
			return
		}

		character := statsGeneration(m.Content, authorId)

		// Show the new character stats & name
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title: fmt.Sprintln("This is your character, **", m.Content, "** !\n Here are it's starting stats:"),
			Description: fmt.Sprintln(
				"**Precision:** ", strconv.Itoa(character.Precision),
				"\n**Strength:** ", strconv.Itoa(character.Strength),
				"\n**Endurance:** ", strconv.Itoa(character.Endurance),
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

		createCharQuery := `INSERT INTO character_model
		 (name, player_id, precision, strength, endurance, agility, hitpoints)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING character_model_id;`

		var characterModelId int

		err = database.DB.QueryRow(createCharQuery,
			m.Content, authorId, character.Precision, character.Strength, character.Endurance, character.Agility, character.Hitpoints).Scan(&characterModelId)

		if err != nil {
			panic(err)
		}

		createCharInstanceQuery := `INSERT INTO character_instance
		(character_model_id, precision, strength, endurance, agility, hitpoints) 
		VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = database.DB.Exec(createCharInstanceQuery,
		characterModelId, character.Precision, character.Strength, character.Endurance, character.Agility, character.Hitpoints)

		if err != nil {
			panic(err)
		}

	} else {
		s.ChannelMessageSend(m.ChannelID, "Aborting character creation")
	}
}

// Generate the different stats (random)
func statsGeneration(givenName string, author int64) *game.CharacterModel {
	character := game.CharacterModel{
		Name:          givenName,
		PlayerId:       author,
		Precision:     (rand.Intn(20) + 20),
		Strength:      (rand.Intn(20) + 20),
		Endurance:     (rand.Intn(20) + 20),
		Agility:       (rand.Intn(20) + 20),
		Hitpoints:     (rand.Intn(7) + 8)}

	return &character
}
