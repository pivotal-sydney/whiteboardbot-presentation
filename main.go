package main

import (
        "fmt"
        "github.com/nlopes/slack"
        "net/http"
        "os"
	"strings"
)

const (
        DEFAULT_PORT = "9000"
)

var (
        index int
	slides = [...]string {

`
*WhiteboardBot* - A journey discovering Slack
>Andrew Leung
>aleung@pivotal.io

>Dariusz Lorenc
>dlorenc@pivotal.io
`,

`
*Backstory…*​

- Daily standups @ 9:06am

`,

`http://static.stthomas.edu/newsroom/news/wp-content/uploads/2015/06/img_pivotal-jp.jpeg`,

`http://i.imgur.com/qS1GA34.png`,

`
*Typical Whiteboard usage*
- Works great when accessed using a desktop
- Can utilise web interface
`,

`
*Not perfect mobile experience*
- Require sign in process in order to input and view items
- Website is not mobile friendly
`,

`
*Solution*
- Pivotal uses Slack as a communication tool between offices, teams, technical groups etc.
- Slack provides integrations:
	- bot (web api, realtime api)
	- slash commands
	- webhooks
- Simple racquetbot had already been developed
- Use Slack to create a bot with conversational UI style
`,

`
*Mobile App vs Slack Bot*

Supporting different platforms for doing the same thing
- need to create a new app for iOS, Android
- slow...
- need many screens
- design app twice to match mobile device style
- need to design user interactions
- screen flows
- branding
- support for old and new OS versions
- distribute a brand new app to install on their devices
- learn different languages to build on multiple mobile devices
`,

`
Slack iOS and Android clients already exist.  It IS the UI interface to our app. (And many other bot apps)
 - text based interface is simple, text as input, text as output.  No more building buttons.
 - no to very little visual design required (no pixel pushing, branding etc.)
 - conversational UI is a delightful user experience
 - markdown, emoji, attachments, links for richer user experience
 - very very fast development cycle providing immediate user value
`,

`
*Our Journey*
`,

`
First day
 - Within mere hours, first fully functioning feature tested, implemented, and delivered
 - TDD/BDD setup
 - Continuous Integration (Cicrcle CI)
 - Backlog of stories in Pivotal Tracker
`,

`
Second day
 - more features: ability to create and update new faces... then ported to events, helps, interestings
`,

`
Iterated on feature set
 - abbreviations to reduce number of keystrokes
 - case insensitivity because mobile phones tend to capitalize the first character
 - wb ? for help and usage
 - had Pam help with design for look and feel of the user feedback messages (show iterations through demos)
 - ability to present so we could run standup all using just slack
 - clunky initial pass required users to start an item, then update.
 - combined title into create command
 - we abbreviated commands
 - improved input to be unique for users
 - added ability to register to multiple standups
 - auto-complete the author by pulling data from Slack
`,

`

`,

	}
)

func main() {
        api := slack.New(os.Getenv("BOT_API_TOKEN"))
        rtm := api.NewRTM()
        go rtm.ManageConnection()

        go startHttpServer()

        Loop:
        for {
                select {
                case msg := <-rtm.IncomingEvents:
                        switch ev := msg.Data.(type) {
                        case *slack.MessageEvent:
                                go ParseMessageEvent(rtm, ev)
                        case *slack.InvalidAuthEvent:
                                fmt.Println("Invalid credentials")
                                break Loop
                        default:
                        }
                }
        }
}

func startHttpServer() {
        http.HandleFunc("/", HealthCheckServer)
        if err := http.ListenAndServe(":" + getHealthCheckPort(), nil); err != nil {
                fmt.Printf("ListenAndServe: %v\n", err)
        }
}

func getHealthCheckPort() (port string) {
        if port = os.Getenv("PORT"); len(port) == 0 {
                fmt.Printf("Warning, PORT not set. Defaulting to %+v\n", DEFAULT_PORT)
                port = DEFAULT_PORT
        }
        return
}

func HealthCheckServer(responseWriter http.ResponseWriter, req *http.Request) {
        fmt.Fprintln(responseWriter, "I'm alive")
}

func ParseMessageEvent(rtm *slack.RTM, event *slack.MessageEvent) {
        if event.Msg.Text == "next" {
		if index >= len(slides) {
			postMessage("The End...", rtm, event)
		}
		postMessage(slides[index], rtm, event)
        	index++
        }
}

func postMessage(text string, rtm *slack.RTM, event *slack.MessageEvent) {
	if strings.HasPrefix(text, "http") {
		rtm.PostMessage(event.Channel, text, slack.PostMessageParameters{AsUser: true, Attachments: []slack.Attachment{{ImageURL: text}}})
	} else {
		rtm.PostMessage(event.Channel, text, slack.PostMessageParameters{AsUser: true})
	}
}