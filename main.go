package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"net/http"
	"os"
)

const (
	DEFAULT_PORT = "9000"
	TAB_COLOR    = "49C49F"
)

var (
	index  int
	slides = [...]slack.Attachment{
		slack.Attachment{
			Title:      "Backstory…",
			Text:       `Daily standups @ 9:06am`,
			ImageURL:   "http://static.stthomas.edu/newsroom/news/wp-content/uploads/2015/06/img_pivotal-jp.jpeg",
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title:      "Whiteboard",
			Text:       ``,
			ImageURL:   "http://i.imgur.com/qS1GA34.png",
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Typical Whiteboard usage",
			Text: `
- Works great when accessed using a desktop
- Can utilise web interface`,
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Not perfect mobile experience",
			Text: `
- Require sign in process in order to input and view items
- Website is not mobile friendly`,
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Solution",
			Text:
`- ​*We use Slack*​ as a communication tool
- Slack provides integrations - ​*bots*​
- Simple racquetbot had already been developed
- Why not use Slack to create a bot with ​*conversational UI style*​`,
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Mobile App vs Slack Bot",
			Text:
`Supporting different mobile platforms
- create a new app for ​*iOS, Android*​
- ​*design app twice*​ to match mobile device style
- many screens
- design user interactions / screen flows
- branding
- support for old and new OS versions
- ​*distribute*​ a brand new app to install on their devices
- ​*different languages*​ to build on multiple mobile devices`,
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Mobile App vs Slack Bot",
			Text:
`Slack clients for mobile devices already exist. It IS the UI interface to our app. (And many other bot apps)
- text based interface is ​*simple*​, text as input, text as output. ​*No more building buttons.*​
- very little visual design required (​*no pixel pushing*​, branding etc.)
- ​*conversational UI*​ is a delightful user experience
- markdown, emoji, attachments, links for richer user experience
- ​*fast development cycle*​ providing ​*immediate user value*​`,
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Our Journey",
			Text:
`First day
 - Within mere hours, first fully functioning feature tested, implemented, and delivered
Second day
 - more features: ability to create and update new faces... then ported to events, helps, interestings`,
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Iterated on feature set",
			ImageURL: "https://cdn-images-2.medium.com/max/800/1*a036VX4scwQquEtBOef8tg.jpeg",
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Demo time - Iteration 1",
			Text: `- mirror message, new faces, helps, interesting, events`,
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Demo time - Iteration 2",
			Text:
`- abbreviated commands
- case insentivity`,
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Demo time - Iteration 3",
			Text:
`- combined title into create command
- help/usage
- present standup`,
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Demo time - Iteration 4",
			Text:
`- look & feel improvements`,
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Summary",
			Text:
`- simple UI
- tight development cycles
- immediate user value
- conversational UI style`,
			ImageURL: "https://cdn-images-2.medium.com/max/800/1*B0dX0geQyEYFmCyd5kJIWw.jpeg",
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
		slack.Attachment{
			Title: "Q & A",
			Text:
`--- ​*Thanks!*​ ---`,
			Color:      TAB_COLOR,
			MarkdownIn: []string{"text"},
		},
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
	if err := http.ListenAndServe(":"+getHealthCheckPort(), nil); err != nil {
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
	if matches(event.Msg.Text, "start") {
		postMessage(slides[0], rtm, event)
		index = 1
	}

	if matches(event.Msg.Text, "next") {
		if index < len(slides) {
			postMessage(slides[index], rtm, event)
			index++
		}
	}
}

func matches(keyword string, command string) bool {
	return len(keyword) > 0 && len(keyword) <= len(command) && command[:len(keyword)] == keyword
}

func postMessage(slide slack.Attachment, rtm *slack.RTM, event *slack.MessageEvent) {
	rtm.PostMessage(event.Channel, "", slack.PostMessageParameters{AsUser: true, Markdown: false, Attachments: []slack.Attachment{slide}})
}
