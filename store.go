package main

import (
	"flag"
	"fmt"
	"github.com/henrikssn/stored/endpoint"
	"log"
	"net"
	"os"
)

var (
	addr        = flag.String("l", "127.0.0.1:8080", "The address to connect to.")
	namespace   = flag.String("n", "default_namespace", "The namespace to use.")
	group       = flag.String("g", "default_group", "The group to use.")
	verbose     = flag.Bool("v", false, "Verbose print.")
	showVersion = flag.Bool("V", false, "print store's version string")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] [id] [store_data] \n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

var (
	op = map[string]func([]string){
		"get": get,
		"put": put,
		"del": del,
	}
	client *endpoint.Client
)

var conn net.Conn

func main() {
	flag.Usage = Usage
	flag.Parse()

	args := flag.Args()

	if len(args) < 2 {
		Usage()
		return
	}

	c, err := endpoint.NewClient("http://" + *addr)
	if err != nil {
		log.Printf("Could not create endpoint client: ", err)
	}
	client = c
	op[args[0]](args[1:])

}

func get(args []string) {
	resp, err := client.Get(endpoint.Key{*namespace, *group, args[0]})
	if err != nil {
		log.Printf("An error occured: ", err)
	}
	fmt.Println("HTTP " + client.Response.Status)
	fmt.Printf("%s", resp)
}

func put(args []string) {
	err := client.Put(endpoint.Key{*namespace, *group, args[0]}, []byte(args[1]))
	if err != nil {
		log.Printf("An error occured: ", err)
	}
	fmt.Println("HTTP " + client.Response.Status)
}

func del(args []string) {
	err := client.Delete(endpoint.Key{*namespace, *group, args[0]})
	if err != nil {
		log.Printf("An error occured: ", err)
	}
	fmt.Println("HTTP " + client.Response.Status)
}
