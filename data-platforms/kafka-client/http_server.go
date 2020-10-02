package main


import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	addr		= flag.String("addr", "localhost:8080", "The address to bind to")
	brokers 	= flag.String("brokers","","csv broker list")
	verbose		= flag.Bool("verbose", false, "If Set true, use sarama logging")
	certFile	= flag.String("certificate", "","")
	keyFile		= flag.String("key file", "","")
	caFile		= flag.String("ca", "","")
	verifySsl	= flag.String("ssl certs", "","")
	topic 		= flag.String("topic name", "simple", "Default Topic is `simple`")
)

func main() {
	flag.Parse()

	if *verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	if *brokers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	brokerList := strings.Split(*brokers, ",")

	server := &Server{

	}

}

type Server struct {
	DataCollector		sarama.SyncProducer
	AccessLogProducer	sarama.AsyncProducer
}

func (s *Server) Close() error {
	if err := s.DataCollector.Close(); err != nil {
		log.Println("Failed to shut down data collector clearly", err)
		return err
	}
	if err := s.AccessLogProducer.Close(); err != nil {
		log.Println("Failed to shut down access log producer clearly", err)
		return err
	}
	log.Println("Closed successfully")
	return nil
}

func (s *Server) Handler() http.Handler {
	return s.withAccessLog(s.collectQueryStringData())
}

func (s *Server) collectQueryStringData() http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// no message key -> randomly distributed over the parts.
		partition, offset, err := s.DataCollector.SendMessage(
			&sarama.ProducerMessage{
				Topic: *topic,
				Value: sarama.StringEncoder(r.URL.RawQuery),
				})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Failed to store your data:, %s", err)
		} else {
			// (topic, partition, offset) 튜플은 unique id로 사용 가능
			fmt.Fprintf(w, "Your data is stored with unique identifier `topic`/%d/%d", partition, offset)
		}
	})
}

type accessLogEntry struct {
	Method 		string	`json:"method"`
	Host		string	`json:"host"`
	Path		string	`json:"path"`
	IP			string	`json:"ip"`
	RespTime 	float64	`json:"response_time"`

	encoded 	[]byte
	err 		error
}

func (ale *accessLogEntry) ensureEncoded() {
	if ale.encoded == nil && ale.err == nil {
		ale.encoded, ale.err = json.Marshal(ale)
	}
}

func (ale *accessLogEntry) Length() int {
	ale.ensureEncoded()
	return len(ale.encoded)
}

func (ale *accessLogEntry) Encode() ([]byte, error) {
	ale.ensureEncoded()
	return ale.encoded, ale.err
}