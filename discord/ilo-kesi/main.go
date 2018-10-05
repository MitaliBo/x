package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Xe/x/internal"
	"github.com/Xe/x/markov"
	"github.com/Xe/x/web/switchcounter"
	"github.com/bwmarrin/discordgo"
	"github.com/joeshaw/envdecode"
	_ "github.com/joho/godotenv/autoload"
	"github.com/peterh/liner"
)

var (
	repl = flag.Bool("repl", false, "open a bot repl in the console?")
)

// lipuSona is the configuration.
type lipuSona struct {
	DiscordToken            string   `env:"DISCORD_TOKEN,required"` // lipu pi lukin ala
	TokiPonaTokenizerAPIURL string   `env:"TOKI_PONA_TOKENIZER_API_URL,default=https://us-central1-golden-cove-408.cloudfunctions.net/function-1"`
	SwitchCounterWebhook    string   `env:"SWITCH_COUNTER_WEBHOOK,required"`
	IloNimi                 string   `env:"ILO_NIMI,default=Kesi"`
	janLawa                 []string `env:"JAN_LAWA,required"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.Parse()
	internal.HandleLicense()

	var cfg lipuSona
	err := envdecode.StrictDecode(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	cfg.janLawa = append(cfg.janLawa, "console")

	flag.Parse()
	internal.HandleLicense()

	//pretty.Println(cfg)

	sw := switchcounter.NewHTTPClient(http.DefaultClient, cfg.SwitchCounterWebhook)

	line := liner.NewLiner()
	defer line.Close()

	chain := markov.NewChain(3)
	err = chain.Load("cadey.gob")
	if err != nil {
		log.Fatal(err)
	}

	words, err := loadWords("./tokipona.json")
	if err != nil {
		log.Fatal(err)
	}

	i := ilo{
		cfg:   cfg,
		sw:    sw,
		chain: chain,
		words: words,
	}

	line.SetCtrlCAborts(true)

	mc := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		msg := m.ContentWithMentionsReplaced()
		if !i.tokiNiTokiPonaAnuSeme(msg) {
			return
		}

		result, err := i.parse(m.Author.ID, msg)
		if err != nil {
			switch err {
			case ErrJanLawaAla, ErrUnknownAction:
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("mi ken ala la %v", err))
				return
			}

			log.Printf("other error: %s", err)
			return
		}

		s.ChannelMessageSend(m.ChannelID, result.msg)
	}

	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatal(err)
	}

	dg.AddHandler(mc)
	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dg.Close()

	if *repl {
		for {
			if inp, err := line.Prompt("|lipu: "); err == nil {
				if inp == "" {
					return
				}

				line.AppendHistory(inp)

				result, err := i.parse("console", inp)
				if err != nil {
					log.Printf("error: %v", err)
					continue
				}

				fmt.Println(result.msg)
			} else if err == liner.ErrPromptAborted {
				log.Print("Aborted")
				break
			} else {
				log.Print("Error reading line: ", err)
				break
			}
		}

		os.Exit(0)
	} else {
		log.Println("bot is running")
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
	}
}
