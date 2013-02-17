package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	listenAddr = flag.String("listen", ":2000", "local listen address")
	dataDir    = flag.String("datadir", "/usr/local/share/ponyserver", "path to data directory")
)

type Pony struct {
	Name   string
	Image  string
	Quotes []string
}

var ponies = make(map[string]*Pony)
var ponynamesWithQuotes []string

func main() {
	flag.Parse()
	log.SetPrefix("ponyserver: ")

	preloadPonies()
	if len(ponies) == 0 {
		log.Fatal("No ponies found in datadir ", *dataDir)
	}

	l, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("listen", *listenAddr)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("accept conn", "localAddr.", conn.LocalAddr(), "remoteAddr.", conn.RemoteAddr())
		go handler(conn)
	}
}

func handler(c net.Conn) {
	defer c.Close()

	rb := make([]byte, 512)
	_, err := c.Read(rb)
	if err != nil {
		log.Fatal(err)
	}
	//request := strings.Trim(string(rb), "\r\n"+string(0))
	p := ponies[ponynamesWithQuotes[rand.Intn(len(ponynamesWithQuotes))]]

	regex, err := regexp.Compile(`\$balloon\d+\$`)
	if err != nil {
		log.Fatal(err)
	}
	p.Image = regex.ReplaceAllString(p.Image, drawBalloon(p.Quotes[rand.Intn(len(p.Quotes))]))

	r := strings.NewReplacer("$\\$", "\\")
	_, err = c.Write([]byte(r.Replace(p.Image)))
	if err != nil {
		log.Fatal(err)
	}
}

func drawBalloon(m string) string {
	mess := strings.Split(m, "\n")
	max := 0
	for _, line := range mess {
		if max <= len(line) {
			max = len(line)
		}
	}
	l := max + 2

	balloon := "┌" + strings.Repeat("─", l) + "┐\n"
	for _, line := range mess {
		space := max - len(line)
		if space == max { // empty line
			break
		}
		balloon += "│ " + line + strings.Repeat(" ", space) + " │\n"
	}
	balloon += "└" + strings.Repeat("─", l) + "┘"
	return balloon
}

func preloadPonies() {
	pf, err := filepath.Glob(*dataDir + "/ponies/*.pony")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range pf {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		n := strings.TrimRight(filepath.Base(f), ".pony")
		ponies[n] = &Pony{
			Name:   n,
			Image:  string(b),
			Quotes: preloadQuotes(n),
		}
	}
	for _, p := range ponies {
		if len(p.Quotes) > 0 {
			ponynamesWithQuotes = append(ponynamesWithQuotes, p.Name)
		}
	}
}

func preloadQuotes(n string) []string {
	pf, err := filepath.Glob(*dataDir + "/ponyquotes/" + n + ".*")
	if err != nil {
		log.Fatal(err)
	}
	var quotes []string
	for _, f := range pf {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		quotes = append(quotes, string(b))
	}
	return quotes
}
