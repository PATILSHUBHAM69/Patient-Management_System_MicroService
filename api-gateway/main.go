// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"net/http/httputil"
// 	"net/url"
// )

// var servicePorts = map[string]string{
// 	"service1": "8080",
// 	"service2": "8081",
// 	"service3": "8082",
// 	"service4": "8083",
// }

// func reverseProxy(service string) http.Handler {
// 	targetURL, err := url.Parse("http://localhost:" + servicePorts[service])
// 	if err != nil {
// 		log.Fatal("Error parsing target URL:", err)
// 	}
// 	proxy := httputil.NewSingleHostReverseProxy(targetURL)
// 	return proxy
// }

// func staticHandler() http.Handler {
// 	return http.FileServer(http.Dir("static"))
// }

// func main() {
// 	http.Handle("/service1/", reverseProxy("service1"))
// 	http.Handle("/service2/", reverseProxy("service2"))
// 	http.Handle("/service3/", reverseProxy("service3"))
// 	http.Handle("/service4/", reverseProxy("service4"))

// 	http.Handle("/static/", http.StripPrefix("/static/", staticHandler()))

// 	fmt.Println("API Gateway is running on port 8084")
// 	log.Fatal(http.ListenAndServe(":8084", nil))
// }

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var servicePorts = map[string]string{
	"service1": "8080",
	"service2": "8081",
	"service3": "8082",
	"service4": "8083",
}

func reverseProxy(targetURL *url.URL) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	return proxy
}

func customRouterHandler(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")

	if len(pathSegments) >= 2 {
		service := pathSegments[1]
		if port, exists := servicePorts[service]; exists {
			targetURL, err := url.Parse("http://localhost:" + port)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			proxy := reverseProxy(targetURL)
			proxy.ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func staticHandler() http.Handler {
	return http.FileServer(http.Dir("static"))
}

func main() {
	http.HandleFunc("/", customRouterHandler)
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler()))

	fmt.Println("API Gateway is running on port 8084")
	log.Fatal(http.ListenAndServe(":8084", nil))
}
