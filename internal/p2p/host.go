package p2p

import (
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

func CreateHost(PORT_HOST string) (host.Host, string, error) {

	PORT := "/ip4/127.0.0.1/tcp/" + PORT_HOST

	h, err := libp2p.New(
		libp2p.ListenAddrStrings(
			PORT,
		),
	)
	if err != nil {
		log.Fatal("Host node unable to create")
		return nil, "", err
	}
	CONNECTION_STRING := PORT + "/p2p/" + h.ID().String()
	log.Println("Hello,my Peer ID is: ", h.ID())
	log.Println("Listening on: ", PORT)
	log.Println("My host Address: ", h.Addrs())

	return h, CONNECTION_STRING, nil
}
