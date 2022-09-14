package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jessevdk/go-flags"
)

type options struct {
	Server   string `short:"s" long:"server" description:"DNS Server" required:"false" default:"8.8.8.8"`
	Port     int    `short:"p" long:"port" description:"query num" required:"false" default:"53"`
	Domain   string `short:"d" long:"domain" description:"Domain Name" required:"true"`
	Type     string `short:"t" long:"type" description:"Record Type" required:"false"`
	Count    int    `short:"n" long:"count" description:"query num" required:"false" default:"3"`
	Timeout  int    `long:"timeout" description:"deadline time" required:"false"`
	Protocol string `long:"protocol" description:"Record Type" required:"false" default:"udp"`
	Thread   int    `long:"threads" description:"thread num" required:"false" default:"1"`
	Version  bool   `long:"version" description:"show version"`
	Debug    bool   `long:"debug" description:""`
	Verbose  bool   `long:"verbose" description:""`
}

type DNSQuery struct {
	ID           uint16
	QR           bool
	Opcode       uint8
	AA           bool
	TC           bool
	RD           bool
	RA           bool
	Z            uint8
	ResponseCode uint8
	QDCount      uint16
	ANCount      uint16
	NSCount      uint16
	ARCount      uint16
	Questions    []DNSQuestion
}

type DNSQuestion struct {
	Domain string
	Type   uint16
	Class  uint16
}

func (q DNSQuery) encode() []byte {

	q.QDCount = uint16(len(q.Questions))

	var buffer bytes.Buffer

	binary.Write(&buffer, binary.BigEndian, q.ID)

	b2i := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}

	queryParams1 := byte(b2i(q.QR)<<7 | int(q.Opcode)<<3 | b2i(q.AA)<<1 | b2i(q.RD))
	queryParams2 := byte(b2i(q.RA)<<7 | int(q.Z)<<4)

	binary.Write(&buffer, binary.BigEndian, queryParams1)
	binary.Write(&buffer, binary.BigEndian, queryParams2)
	binary.Write(&buffer, binary.BigEndian, q.QDCount)
	binary.Write(&buffer, binary.BigEndian, q.ANCount)
	binary.Write(&buffer, binary.BigEndian, q.NSCount)
	binary.Write(&buffer, binary.BigEndian, q.ARCount)

	for _, question := range q.Questions {
		buffer.Write(question.encode())
	}

	return buffer.Bytes()
}

func (q DNSQuestion) encode() []byte {
	var buffer bytes.Buffer

	domainParts := strings.Split(q.Domain, ".")
	for _, part := range domainParts {
		if err := binary.Write(&buffer, binary.BigEndian, byte(len(part))); err != nil {
			log.Fatalf("Error binary.Write(..) for '%s': '%s'", part, err)
		}

		for _, c := range part {
			if err := binary.Write(&buffer, binary.BigEndian, uint8(c)); err != nil {
				log.Fatalf("Error binary.Write(..) for '%s'; '%c': '%s'", part, c, err)
			}
		}
	}

	binary.Write(&buffer, binary.BigEndian, uint8(0))
	binary.Write(&buffer, binary.BigEndian, q.Type)
	binary.Write(&buffer, binary.BigEndian, q.Class)

	return buffer.Bytes()
}

func Do(conn net.Conn, encodedQuery []byte, cnt int) (res []int64) {

	for i := 0; i < cnt; i++ {
		start := time.Now()
		conn.Write(encodedQuery)

		encodedAnswer := make([]byte, len(encodedQuery))
		if _, err := bufio.NewReader(conn).Read(encodedAnswer); err != nil {
			log.Fatal(err)
		}
		res = append(res, time.Since(start).Milliseconds())
	}

	return res
}

func Calc(wg *sync.WaitGroup, c chan<- []int64, conn net.Conn, encodedQuery []byte, cnt int) {
	defer wg.Done()
	res := Do(conn, encodedQuery, cnt)
	c <- res
}

func resultShow(tmp []int64) {
	fmt.Println(tmp)
}

func main() {
	var opts options
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	q := DNSQuestion{
		Domain: opts.Domain,
		Type:   0x1, // A record
		Class:  0x1, // Internet
	}

	query := DNSQuery{
		ID:        0xAAAA,
		RD:        true,
		Questions: []DNSQuestion{q},
	}

	// Setup a UDP connection
	conn, err := net.Dial(opts.Protocol, fmt.Sprintf("%s:%d", opts.Server, opts.Port))
	if err != nil {
		log.Fatal("failed to connect:", err)
	}
	defer conn.Close()

	if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		log.Fatal("failed to set deadline: ", err)
	}

	encodedQuery := query.encode()

	var wg sync.WaitGroup
	c := make(chan []int64, opts.Thread)
	for i := 0; i < opts.Thread; i++ {
		wg.Add(1)
		go Calc(&wg, c, conn, encodedQuery, opts.Count)
	}
	wg.Wait()
	close(c)

	tmp := []int64{}
	for i := range c {
		tmp = append(tmp, i...)
	}

	resultShow(tmp)

}
