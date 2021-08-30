package test


import (
	"github.com/bwmarrin/discordgo"
	"github.com/SteakBarbare/RPGBot/connector"
	// "github.com/SteakBarbare/RPGBot/handlers"
	"testing"
)

type discordSessionMock struct{}

func (session *discordSessionMock) ChannelMessageSend(channelID string, message string) (*discordgo.Message, error) {
	return nil, nil
}

func TestNewCharacter(t *testing.T) {
	type fields struct {
		Connector connector.DiscordInterface
		Message   *discordgo.MessageCreate
	}
	
	// Declare params
	sessionMock := discordSessionMock{}
	discordMessage := discordgo.Message {
		ChannelID: "1",
		Content: "-char New",
	}
	discordMessageCreate := discordgo.MessageCreate {
		&discordMessage,
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Empty query",
			fields: fields{
				Connector: &sessionMock,
				Message: &discordgo.MessageCreate {
					&discordgo.Message {
						ChannelID: "1",
						Content: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "true query",
			fields: fields{
				Connector: &sessionMock,
				Message: &discordMessageCreate,
			},
			wantErr: false,
		},
	}
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		command := handlers.NewCharacterCommand{
	// 			Connector: tt.fields.Connector,
	// 			Message:   tt.fields.Message,
	// 		}
	// 		if err := command.Execute(); (err != nil) != tt.wantErr {
	// 			t.Errorf("QueryGoogleCommand.Execute() error = %v, wantErr %v", err, tt.wantErr)
	// 		}
	// 	})
	// }
}