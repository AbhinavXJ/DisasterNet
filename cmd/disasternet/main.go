package main

import (
	"DisasterNet/internal/p2p"
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	// "github.com/libp2p/go-libp2p/core/peer"
	// ma "github.com/multiformats/go-multiaddr"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type IncomingMsg struct {
	Message string `json:"message"`
}

var MessageArr []string
var messageMu sync.Mutex

var cr *p2p.ChatRoom

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")                   // Allow all origins to access
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Allowed HTTP methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // Allowed headers
}

func StoreMessage(msg string) {
	messageMu.Lock()
	defer messageMu.Unlock()
	MessageArr = append(MessageArr, msg)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Only Get Method supported", http.StatusBadRequest)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(MessageArr)
	if err != nil {
		http.Error(w, "failed to encode messages", http.StatusInternalServerError)
		return
	}

}

func PostMessage(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Only POST Method supported", http.StatusBadRequest)
		return

	}

	var msg_post IncomingMsg
	err := json.NewDecoder(r.Body).Decode(&msg_post)
	if err != nil || msg_post.Message == "" {
		http.Error(w, "failed to decode", http.StatusBadRequest)
		return
	}

	err_pub := cr.Publish(msg_post.Message)
	StoreMessage(msg_post.Message)

	if err_pub != nil {
		fmt.Println("Sending message failed trying again...")
		http.Error(w, "failed to publish", http.StatusInternalServerError)

		return
	}
	w.WriteHeader(http.StatusOK)

}

func main() {

	port := flag.String("port", "", "port")
	nickFlag := flag.String("nick", "", "nickname to use in chat. will be generated if empty")
	roomFlag := flag.String("room", "chat-room", "name of chat room to join")
	httpServerRun := flag.Bool("enable-http", false, "run http server on this node")

	// peerAddr := flag.String("peer-address", "", "peer address")
	sameNetworkString := flag.String("same_string", "", "same_string")

	flag.Parse()
	h, _, err1 := p2p.CreateHost(*port)

	if err1 != nil {
		log.Fatal("error creating the host")
	}

	ctx := context.Background()

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}
	peerChan := p2p.InitMDNS(h, *sameNetworkString)

	go func() {

		for {
			peer := <-peerChan // will block until we discover a peer
			if peer.ID > h.ID() {
				// if other end peer id greater than us, don't connect to it, just wait for it to connect us
				fmt.Println("Found peer:", peer, " id is greater than us, wait for it to connect to us")
				continue
			}
			fmt.Println("Discovered new peer via mDNS:", peer.ID, peer.Addrs)

			if err := h.Connect(ctx, peer); err != nil {
				fmt.Println("Connection failed:", err)
				continue
			}

			log.Println("Connection to the peer found through MDNS has been established")
			log.Println("Peer Id:", peer.ID, "Peer Addrs: ", peer.Addrs)

		}
	}()
	//logic for connecting to peer through peer-address if we know it beforehand and pass it through terminal

	// If we received a peer address, we should connect to it.
	// if *peerAddr != "" {
	// 	// Parse the multiaddr string.
	// 	peerMA, err := ma.NewMultiaddr(*peerAddr)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// Connect to the node at the given address.
	// 	if err := h.Connect(context.Background(), *peerAddrInfo); err != nil {
	// 		panic(err)
	// 	}
	// 	log.Println("Connected to", peerAddrInfo.String())
	// }

	// use the nickname from the cli flag, or a default if blank
	nick := *nickFlag
	if len(nick) == 0 {
		nick = "ABHI"
	}

	room := *roomFlag

	// join the chat room
	cr, err = p2p.JoinChatRoom(ctx, ps, h.ID(), nick, room)
	if err != nil {
		panic(err)
	}

	if *httpServerRun {
		go func() {

			http.HandleFunc("/send", PostMessage)

			http.HandleFunc("/messages", GetMessages)
			err := http.ListenAndServe(":3001", nil)
			if err != nil {
				log.Fatal(err)
			}

		}()
	}

	f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("error opening logs.txt")
	}
	// Read incoming messages
	go func() {
		for msg := range cr.Messages {
			text := fmt.Sprintf("Received message at %s from %s: %s\n", time.Now().Local(), msg.SenderNick, msg.Message)
			StoreMessage(text)
			fmt.Printf("Received message at %s from %s: %s\n", time.Now().Local(), msg.SenderNick, msg.Message)
			_, err_log := f.WriteString(text)
			if err_log != nil {
				log.Fatal("error writing logs..")
				continue
			}
		}
	}()

	fmt.Println("Sending test message...")
	reader := bufio.NewReader(os.Stdin)
	err = cr.Publish("Hello from " + h.ID().String())
	if err != nil {
		fmt.Println("Error publishing:", err)
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("Sent message from %s: %s\n", nick, "Hello from "+h.ID().String())

		err_pub := cr.Publish(line)

		if err_pub != nil {
			fmt.Println("Sending message failed trying again...")
			cr.Publish(line)
			continue
		}

	}

}
