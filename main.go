/*

rtop-bot - remote system monitoring bot

Copyright (c) 2015 RapidLoop

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ivanfoo/rtop-bot/bot"
)

const (
	VERSION = "0.3"
)

//----------------------------------------------------------------------------

func usage() {
	fmt.Printf(
		`rtop-bot %s - (c) 2015 RapidLoop - http://www.rtop-monitor.org/rtop-bot
rtop-bot is a Slack bot that can do remote system monitoring over SSH

Usage:
    rtop-bot -s slackBotToken [-u user] [-i identity_file]

where:
    -s slackBotToken is the API token for the Slack bot
    -u user
    -i identity_file
`, VERSION)
	os.Exit(1)
}

func main() {
	var SlackToken = flag.String("s", "", "create Slack bot")
	var Username = flag.String("u", "", "ssh user to use")
	var SSHKeyPath = flag.String("i", "", "private key to use")

	flag.Parse()

	if *SlackToken == "" {
		usage()
	}

	bot := bot.NewBot(bot.BotOptions{
		Username:   *Username,
		SSHKeyPath: *SSHKeyPath,
		SlackToken: *SlackToken,
	})

	bot.DoSlack()
	/*
				if *slackFlag == "" {
					usage()
				}
		SlackToken
				log.SetPrefix("rtop-bot: ")
				log.SetFlags(0)

				// get username for SSH connections
				if *userFlag != "" {
					sysUser, err = user.Lookup(*userFlag)
				} else {
					sysUser, err = user.Current()
				}

				if err != nil {
					log.Print(err)
					os.Exit(1)
				}

				sshUsername = sysUser.Username

				// expand ~/.ssh/id_rsa and check if it exists
				if *keyFlag != "" {
					idRsaPath = *keyFlag
				} else {
					idRsaPath = filepath.Join(sysUser.HomeDir, ".ssh", "id_rsa")
				}

				if _, err := os.Stat(idRsaPath); os.IsNotExist(err) {
					idRsaPath = ""
				}

				// expand ~/.ssh/config and parse if it exists
				sshConfig := filepath.Join(sysUser.HomeDir, ".ssh", "config")
				if _, err := os.Stat(sshConfig); err == nil {
					parseSshConfig(sshConfig)
				}

				bot.doSlack(*slackFlag)
	*/
}
