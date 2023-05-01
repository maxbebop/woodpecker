package chatservice

/*
func TestNew(t *testing.T) {
	chatBot := mocks.NewChatBot(t)
	want := mockNewServece(chatBot)
	bot := new(chatBot)
	require.Equal(t, bot, want, "creating chat servese")

}

func TestStartChat(t *testing.T) {

	inMsgChannel := make(chan slack.Message)
	chatBot := &mocks.ChatBot{}
	go chatBot.On("GetMessages", inMsgChannel)
	outMsg := slack.OutMessage{Message: slack.Message{}}
	outMsg.Type = slack.Common
	chatBot.On("SendMessage", outMsg)
	close(inMsgChannel)
}

func mockNewServece(chatBot ChatBot) *service {
	return &service{chatBot: chatBot}
}
*/
