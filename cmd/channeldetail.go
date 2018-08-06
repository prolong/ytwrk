package cmd

import (
	"fmt"
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/terorie/yt-mango/api"
	"github.com/terorie/yt-mango/net"
	"github.com/terorie/yt-mango/data"
)

var channelDetailCmd = cobra.Command{
	Use: "detail <channel ID>",
	Short: "Get detail about a channel",
	Args: cobra.ExactArgs(1),
	Run: cmdFunc(doChannelDetail),
}

func doChannelDetail(_ *cobra.Command, args []string) error {
	channelID := args[0]

	channelID, err := api.GetChannelID(channelID)
	if err != nil { return err }

	channelReq := api.Main.GrabChannel(channelID)

	res, err := net.Client.Do(channelReq)
	if err != nil { return err }

	var c data.Channel
	api.Main.ParseChannel(&c, res)

	bytes, err := json.MarshalIndent(&c, "", "\t")
	if err != nil { return err }
	fmt.Println(string(bytes))

	return nil
}
