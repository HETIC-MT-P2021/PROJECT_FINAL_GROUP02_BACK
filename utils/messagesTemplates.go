package utils

import "github.com/bwmarrin/discordgo"

// Formats the user in a readable format
func FormatUser(u *discordgo.User) string {
	return u.Username + "#" + u.Discriminator
}

// Generic message format for errors
func ErrorMessage(title string, message string) string {
	return "❌  **" + title + "**\n" + message
}

// Generic message format for successful operations
func SuccessMessage(title string, message string) string {
	return "✅  **" + title + "**\n" + message
}

// Good for making sure only 1 reaction is selected to prevent reaction spam
func HasOtherReactionsBesides(allowed string, reactions []*discordgo.MessageReactions) bool {
	for _, r := range reactions {
		// If there is more than one reaction on a reaction that isn't the one allowed
		if r.Count > 1 && allowed != r.Emoji.Name {
			// If it has a count greater than 1
			return true
		}
	}

	return false
}
