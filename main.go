package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var (
	Token      string
	userID     string
	re         = regexp.MustCompile("(discord.com/gifts/|discordapp.com/gifts/|discord.gift/)([a-zA-Z0-9]+)")
	rePrivnote = regexp.MustCompile("https://privnote.com/.*")
	reGiveaway = regexp.MustCompile("You won the \\*\\*(.*)\\*\\*")
	magenta    = color.New(color.FgMagenta)
	green      = color.New(color.FgGreen)
	red        = color.New(color.FgRed)
	strPost    = []byte("POST")
	strGet     = []byte("GET")
)

func init() {
	file, err := ioutil.ReadFile("token.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed read file: %s\n", err)
		os.Exit(1)
	}

	var f interface{}
	err = json.Unmarshal(file, &f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse JSON: %s\n", err)
		os.Exit(1)
	}

	m := f.(map[string]interface{})

	str := fmt.Sprintf("%v", m["token"])

	flag.StringVar(&Token, "t", str, "Token")
	flag.Parse()
}

func main() {
	c := exec.Command("clear")

	c.Stdout = os.Stdout
	c.Run()
	color.Red(`Imma rape your discord
	
	Made with <3 by Daddie0

`)
	dg, err := discordgo.New(Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	t := time.Now()
	color.Cyan("Sniping Discord Nitro on " + strconv.Itoa(len(dg.State.Guilds)) + " Servers 🔫\n\n")

	magenta.Print(t.Format("15:04:05 "))
	fmt.Println("[+] Bot is ready")
	userID = dg.State.User.ID

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if re.Match([]byte(m.Content)) {

		code := re.FindStringSubmatch(m.Content)

		if len(code) < 2 {
			return
		}

		if len(code[2]) < 16 {
			magenta.Print(time.Now().Format("15:04:05 "))
			red.Print("[=] Auto-detected a fake code: ")
			red.Print(code[2])
			println(" from " + m.Author.String())
			return
		}

		var strRequestURI = []byte("https://discordapp.com/api/v6/entitlements/gift-codes/" + code[2] + "/redeem")
		req := fasthttp.AcquireRequest()
		req.Header.SetContentType("application/json")
		req.Header.Set("authorization", Token)
		req.SetBody([]byte(`{"channel_id":` + m.ChannelID + "}"))
		req.Header.SetMethodBytes(strPost)
		req.SetRequestURIBytes(strRequestURI)
		res := fasthttp.AcquireResponse()

		if err := fasthttp.Do(req, res); err != nil {
			panic("handle error")
		}

		fasthttp.ReleaseRequest(req)

		body := res.Body()

		bodyString := string(body)
		magenta := color.New(color.FgMagenta)
		magenta.Print(time.Now().Format("15:04:05 "))
		green.Print("[-] Sniped code: ")
		red.Print(code[2])
		println(" from " + m.Author.String())
		magenta.Print(time.Now().Format("15:04:05 "))
		if strings.Contains(bodyString, "This gift has been redeemed already.") {
			color.Yellow("[-] Code has been already redeemed")
		}
		if strings.Contains(bodyString, "nitro") {
			green.Println("[+] Code applied")
		}
		if strings.Contains(bodyString, "Unknown Gift Code") {
			red.Println("[x] Invalid Code")
		}
		fasthttp.ReleaseResponse(res)

	} else if strings.Contains(strings.ToLower(m.Content), "**giveaway**") || (strings.Contains(strings.ToLower(m.Content), "react with") && strings.Contains(strings.ToLower(m.Content), "giveaway")) {
		time.Sleep(time.Minute)
		magenta.Print(time.Now().Format("15:04:05 "))
		color.Yellow("[-] Enter Giveaway ")
		s.MessageReactionAdd(m.ChannelID, m.ID, "🎉")

	} else if (strings.Contains(strings.ToLower(m.Content), "giveaway") || strings.Contains(strings.ToLower(m.Content), "win") || strings.Contains(strings.ToLower(m.Content), "won")) && strings.Contains(m.Content, userID) {
		var won = reGiveaway.FindStringSubmatch(m.Content)
		if len(won) < 2 {
			return
		}
		magenta.Print(time.Now().Format("15:04:05 "))
		green.Print("[+] Won Giveaway: ")
		color.Magenta(won[1])
	}

}
