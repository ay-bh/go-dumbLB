package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

type server struct{
	adrs string
	proxy *httputil.ReverseProxy
}

type Server interface{
	Address() string
	isAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

func newServer(adrs string) *server {
	serverUrl, err := url.Parse(adrs)
	handleErr(err)

	return &server{
		adrs: adrs,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}



type loadbalancer struct {
	port string
	roundRobinCount int
	servers []Server
}

func newLoadBalancer(port string, servers []Server) *loadbalancer {
	return &loadbalancer{
		port: port,
		roundRobinCount: 0,
		servers: servers,
	}
}
func handleErr(err error){
	if err!= nil {
		fmt.Printf("error: %v\n",err)
		os.Exit(1)
	}
}

func (s *server) Address() string {return s.adrs}

func (s *server) isAlive() bool {
	client := http.Client{
        Timeout: time.Second * 5,
    }
    resp, err := client.Get(s.adrs)
    if err != nil {
        return false
    }

    defer resp.Body.Close()
    if resp.StatusCode >= 200 && resp.StatusCode < 300 {
        return true
    }
    return false
}

func (s *server) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw,req)
}


func(lb *loadbalancer) getNextAvailableServer() Server{
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.isAlive(){
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}
func(lb *loadbalancer) serverProxy(rw http.ResponseWriter, req *http.Request){
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("Forwarding request to address %q\n",targetServer.Address())
	targetServer.Serve(rw,req)
}

func main(){
	servers := []Server{
		newServer("http://localhost:3001/"),
		newServer("http://localhost:3002/"),
		newServer("http://localhost:3003/"),
		newServer("http://localhost:3004/"),
	}

	lb := newLoadBalancer("8080", servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request){
		lb.serverProxy(rw,req)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("serving requests at localhost:%s\n",lb.port)
	http.ListenAndServe(":"+lb.port,nil)
} 