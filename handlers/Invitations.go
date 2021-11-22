package handlers

import (
	"fmt"

	"github.com/SteakBarbare/RPGBot/database"
	"github.com/SteakBarbare/RPGBot/utils"
	"github.com/bwmarrin/discordgo"
)

func inviteCommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot error", "Error getting channel."))
		return
	}

	// Ensure that the command is not being sent from a dm
	if c.Type == discordgo.ChannelTypeDM {
		s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Invalid channel", "Cannot send invites from a DM"))
		return
	}

	// Check if the mention is linked to a player, or if this is a general invitation
	recipients := m.Mentions
	if len(recipients) == 1 {
		duelInvite(s, m, recipients[0])
	} else if m.Content == "-duel Anyone" {
		generalDuelInvite(s, m)
	} else if len(recipients) == 0 {
		s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Invalid Reciepient", "Ensure you are mentioning the player in the format of @<user>. Or, if you are trying to send a general invite leave the user blank."))
		return
	} else if len(recipients) > 1 {
		s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Invalid invite", "Cannot invite multiple players!"))
	}
}

// Send a DM to invite a specific player to a duel
func duelInvite(s *discordgo.Session, m *discordgo.MessageCreate, recipient *discordgo.User) {

	if m.Author.ID == recipient.ID {
		s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Invalid recipient", "Cannot play against yourself!"))
		return
	}

	if recipient.Bot {
		s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Invalid recipient", "Cannot play against bot!"))
		return
	}

	dm, err := s.UserChannelCreate(recipient.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot error", "Error creating direct message."))
		return
	}

	invite, err := s.ChannelMessageSendEmbed(dm.ID, &discordgo.MessageEmbed{
		Title:       utils.FormatUser(m.Author) + "has challenged you to a duel",
		Description: "Click the  ✅  to accept this invitation, or the  ❌  to deny.",
		Color:       0x0099ff,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "duelInvite:" + m.Author.ID,
		},
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot error", "Error sending invite."))
		return
	}

	s.MessageReactionAdd(dm.ID, invite.ID, "✅")
	s.MessageReactionAdd(dm.ID, invite.ID, "❌")

	s.ChannelMessageSend(m.ChannelID, utils.SuccessMessage("Success", "Invite sent to "+utils.FormatUser(recipient)+"!"))
}

// Sends a general invite for any user in the channel to accept
func generalDuelInvite(s *discordgo.Session, m *discordgo.MessageCreate) {
	invite, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Duel invite from " + utils.FormatUser(m.Author),
		Description: "Click the  ✅  to accept this invitation.",
		Color:       0x0099ff,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "generalDuelInvite:" + m.Author.ID,
		},
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, utils.ErrorMessage("Bot error", "Error sending invite."))
		return
	}

	s.MessageReactionAdd(m.ChannelID, invite.ID, "✅")
}

// Handles all invite related reactions
func duelInvitationHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd, m *discordgo.Message, user *discordgo.User, opponentID string, general bool) {

	// If the reaction comes from the sender of the invite(This will only happen in the case of general invites)
	if r.UserID == opponentID {
		return
	}
	sender, err := s.User(opponentID)
	if err != nil || sender == nil {
		return
	}
	opponentDM, _ := s.UserChannelCreate(opponentID)

	// If accepted, send a message to confirm the invitation, then launch a duel
	if r.Emoji.Name == "✅" && (general || !utils.HasOtherReactionsBesides("✅", m.Reactions)) {
		s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, &discordgo.MessageEmbed{
			Title:       "Invite Accepted!",
			Description: "Invite from " + utils.FormatUser(sender) + " accepted!",
			Color:       0x00ff00,
		})

		opponents := []string{user.ID, opponentID}

		// Create a game object

		duelPreparationId := 0

		err = database.DB.QueryRow(`INSERT INTO duelpreparation (selectingplayer, isready, isover, turn) VALUES ($1, $2, 'false', 0) RETURNING id`, opponents[1], 0).Scan(&duelPreparationId)
		if err != nil {
			fmt.Println(err.Error())
		}

		_, err = database.DB.Exec(`INSERT INTO duelPlayers (preparationid, challenger, challenged, challengerchar, challengedchar) VALUES ($1, $2, $3, $4, $5)`, duelPreparationId, opponents[1], opponents[0], nil, nil)

		if err != nil {
			fmt.Println(err.Error())
		}

		s.ChannelMessageSend(r.ChannelID, utils.SuccessMessage("Game on!", utils.FormatUser(user)+" accepted your duel invite ! Now select a character to send in the arena."))
		// selectCharacterBase(s, m.ChannelID, opponents, 0)
		s.AddHandlerOnce(selectCharacter)

		// Send a message to tell the invitation was declined otherwise
	} else if !general && r.Emoji.Name == "❌" && !utils.HasOtherReactionsBesides("❌", m.Reactions) {
		s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, &discordgo.MessageEmbed{
			Title:       "Invite Declined",
			Description: "Invite from " + utils.FormatUser(sender) + " declined.",
			Color:       0xff0000,
		})
		s.ChannelMessageSend(opponentDM.ID, utils.ErrorMessage("Invite declined", utils.FormatUser(user)+" declined your duel invite."))
	}
}
