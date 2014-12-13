package main

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/henrikssn/stored/server"
	"io"
	"log"
	"net"
	"os"
)

var (
	laddr       = flag.String("l", "127.0.0.1:8046", "The address to connect to.")
	showVersion = flag.Bool("v", false, "print doozerd's version string")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

var op = map[string]func([]string){
	"get": get,
	"put": put,
	"nop": nop,
}

var conn net.Conn

func main() {
	flag.Usage = Usage
	flag.Parse()

	args := flag.Args()

	conn, _ = net.Dial("tcp", *laddr)

	op[args[0]](args[1:])

}

func get(args []string) {
	tag := int64(42)
	req := &server.Request{Tag: &tag, Op: server.Operation_GET.Enum(), Key: &args[0]}
	write(req)
	read()
}

func put(args []string) {
	tag := int64(42)
	req := &server.Request{Tag: &tag, Op: server.Operation_PUT.Enum(), Key: &args[0], Value: []byte(args[1])}
	write(req)
	read()
}

func nop(args []string) {
}

func read() *server.Response {
	var size int32
	binary.Read(conn, binary.BigEndian, &size)
	buf := make([]byte, size)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		panic(err)
	}

	resp := new(server.Response)
	proto.Unmarshal(buf, resp)
	log.Println(resp)
	return resp
}

func write(req *server.Request) {
	data, _ := proto.Marshal(req)

	binary.Write(conn, binary.BigEndian, int32(len(data)))
	conn.Write(data)
}
