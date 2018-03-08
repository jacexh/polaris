package agent

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"github.com/jacexh/polaris/log"
	"go.uber.org/zap"
)

const (
	defaultSnapLen    = 1600
	gatherBuffer      = 100
	connectionTimeout = 3 * time.Minute
	flushInterval     = 1 * time.Minute
	sniffTimeout      = 30 * time.Second
	promiscuous       = false
)

type (
	// Sniffer 端口级嗅探对象
	Sniffer struct {
		iface   string
		snapLen int32
		filter  string
		gather  chan *http.Request
		closed  chan struct{}
	}

	// httpStreamFactory implements tcpassembly.StreamFactory
	httpStreamFactory struct {
		output chan<- *http.Request
	}

	// httpStream will handle the actual decoding of http requests.
	httpStream struct {
		net, transport gopacket.Flow
		r              tcpreader.ReaderStream
	}
)

func newHTTPStreamFactory(output chan<- *http.Request) *httpStreamFactory {
	return &httpStreamFactory{output}
}

func (h *httpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hstream.run(h.output) // Important... we must guarantee that data from the reader stream is read.

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &hstream.r
}

func (h *httpStream) run(c chan<- *http.Request) {
	buf := bufio.NewReader(&h.r)
	for {
		req, err := http.ReadRequest(buf)
		if err == io.EOF {
			// We must read until we see an EOF... very important!
			return
		} else if err != nil {
			log.Logger.Error("Error reading stream",
				zap.String("net", h.net.String()),
				zap.String("transport", h.transport.String()),
				zap.Error(err))
		} else {
			tcpreader.DiscardBytesToEOF(req.Body)
			req.Body.Close()
			log.Logger.Info("Received request from stream",
				zap.String("net", h.net.String()),
				zap.String("transport", h.transport.String()),
			)
			if c != nil {
				c <- req
			}
		}
	}
}

// NewSniffer 根据ip、port实例化一个Sniffer对象
func NewSniffer(ip string, port int) (*Sniffer, error) {
	sn := &Sniffer{
		snapLen: defaultSnapLen,
		filter:  "tcp and dst port " + strconv.Itoa(port),
		gather:  make(chan *http.Request, gatherBuffer),
		closed:  make(chan struct{}, 1),
	}

	ift, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, ifa := range ift {
		addrs, err := ifa.Addrs()
		if err != nil {
			return sn, err
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.String() == ip {
					sn.iface = ifa.Name
					return sn, nil
				}
			case *net.IPAddr:
				if v.IP.String() == ip {
					sn.iface = ifa.Name
					return sn, nil
				}
			}
		}
	}
	return nil, errors.New("cannot found interface by IP " + ip)
}

// Close 停止嗅探
func (sn *Sniffer) Close() {
	sn.closed <- struct{}{}
}

// Run 开始嗅探
func (sn *Sniffer) Run() error {
	log.Logger.Info(fmt.Sprintf("starting capturing on interface %s with BSF filter: %s", sn.iface, sn.filter))

	handle, err := pcap.OpenLive(sn.iface, sn.snapLen, promiscuous, sniffTimeout)
	if err != nil {
		log.Logger.Error("open live stream failed", zap.Error(err))
		return err
	}
	err = handle.SetBPFFilter(sn.filter)
	if err != nil {
		return err
	}

	streamFactory := newHTTPStreamFactory(sn.gather)
	assembler := tcpassembly.NewAssembler(tcpassembly.NewStreamPool(streamFactory))
	log.Logger.Info("reading in packets")

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	ticker := time.Tick(flushInterval)
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates the end of a pcap file.
			if packet == nil {
				return nil
			}

			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Logger.Warn("Unusable packet")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

		case <-ticker:
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(-connectionTimeout))

		case <-sn.closed:
			close(sn.gather) // 在监听停止后，应当关闭采集通道
			return nil
		}

	}
	return nil
}

// Sending 外发提取到的*http.Request对象
func (sn *Sniffer) Sending() <-chan *http.Request {
	return sn.gather
}
