package main

import (
	"context"
	"flag"
	"github.com/golang/glog"
	"os"

	logrus "github.com/sirupsen/logrus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	gw "github.com/xjayleex/idl/protos/grpc-gateway-test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/http"
	"path"
	"strings"
)

const (
	grpcPort = ":10000"
	gwPort = ":8080"
	gwCert = "/Users/ijaehyeon/keys/server.crt"
	gwKey ="/Users/ijaehyeon/keys/server.pem"
	grpcServerCert = "/Users/ijaehyeon/keys/server.crt"
)
var (
	getEndpoint  = flag.String("get", "localhost"+ grpcPort, "endpoint of YourService")
	postEndpoint = flag.String("post", "localhost"+ grpcPort, "endpoint of YourService")
	swaggerPath = flag.String("swagger_path", "proto", "path which contains swagger definitions.")
	logger = logrus.WithFields(logrus.Fields{
		"Name": "gRPC-Gateway",
	})
)

func NewGateway(ctx context.Context, opts ...runtime.ServeMuxOption) (http.Handler, error) {
	mux := runtime.NewServeMux(opts...) // + runtime.WithMarshalerOption()
	grpcDialOpts := []grpc.DialOption{}

	if cred, err := credentials.NewClientTLSFromFile(grpcServerCert,""); err == nil {
		grpcDialOpts = append(grpcDialOpts, grpc.WithTransportCredentials(cred))
	} else {
		logger.Info(err)
		return nil, err
	}


	err := gw.RegisterEnrollmentHandlerFromEndpoint(ctx, mux, *getEndpoint, grpcDialOpts)
	if err != nil {
		logger.Info(err)
		return nil, errors.New("grpc Gateway : `GET` error")
	}

	err = gw.RegisterEnrollmentHandlerFromEndpoint(ctx, mux, *postEndpoint,grpcDialOpts)
	if err != nil {
		logger.Info(err)
		return nil, errors.New("grpc Gateway : `POST error")
	}

	return mux, nil
}

func swaggerHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, ".swagger.json") {
		glog.Errorf("Swagger not found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	glog.Infof("Serving %s", r.URL.Path)
	p := strings.TrimPrefix(r.URL.Path, "/swagger/")
	p = path.Join(*swaggerPath, p)
	http.ServeFile(w, r, p)
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflight := func (w http.ResponseWriter, r *http.Request) {
					headers := []string{"Content-Type", "Accept"}
					w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
					methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
					w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
					glog.Infof("preflight request for %s", r.URL.Path)
				}
				preflight(w,r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func Serve(address string, opts ...runtime.ServeMuxOption) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", swaggerHandler)
	//opts = append(opts, runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler))
	gwHandler, err := NewGateway(ctx, opts...)
	if err != nil {
		return err
	}

	mux.Handle("/", gwHandler)

	return http.ListenAndServeTLS(address, gwCert, gwKey, allowCORS(mux))
}

func init() {
	flag.Parse()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main(){
	if err := Serve(gwPort); err != nil {
		logger.Fatalf("Cannot serve gateway on port %s", gwPort)
	}
}
