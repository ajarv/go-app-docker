package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

 	"gopkg.in/yaml.v2"
	"github.com/gorilla/mux"
	"github.com/go-redis/redis"
)

func logRequest(req *http.Request) {
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("--Request :%v\n------------\n", string(requestDump))
}

func getDebugData(req *http.Request) map[string]interface{}  {
	var hostname, err = os.Hostname()
	if err != nil {
		hostname = "unknown host"
	}

	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
        items := make(map[string]string)
        for _, item := range data {
            key, val := getkeyval(item)
            items[key] = val
        }
        return items
    }
    environment := getenvironment(os.Environ(), func(item string) (key, val string) {
        splits := strings.Split(item, "=")
        key = splits[0]
        val = splits[1]
        return
    })

	data := map[string]interface{} {
		"Host": hostname,
		"ApiVersion":appVersion,
		"AppName":appName,
		"Request": map[string]interface{}{ "Headers": req.Header},
		"Environment": environment}
	return data
}

func writeData(w http.ResponseWriter, r *http.Request,data map[string]interface{}){
	if strings.Contains(r.Header["Accept"][0], "html") {
		w.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(w, data)
		if err != nil {
			w.Write([]byte(`{"result":"Error"}`))
		}
		return
	} 

	if strings.Contains(r.Header["Accept"][0], "json") {
		data:= getDebugData(r)
		w.Header().Set("Content-Type", "application/json")
		b,err := json.Marshal(&data)
		if err != nil {
			w.Write([]byte(`{"result":"Error"}`))
		return
		}
		w.Write(b)
		return
	} 
	
	w.Header().Set("Content-Type", "application/yaml")
	b,err := yaml.Marshal(&data)
	if err != nil {
		w.Write([]byte(`{"result":"Error"}`))
		return
	}
	w.Write(b)

}

func handler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	data:= getDebugData(r)
	writeData(w,r,data)
}

func killHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		log.Printf("Dieing now .. good bye")
		time.Sleep(4 * time.Second)
		os.Exit(3)
	}()

	logRequest(r)
	data:= getDebugData(r)
	data["message"] ="Will terminate myself on your request in a few .. good bye !!"
	writeData(w,r,data)
}
func redisHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	data:= getDebugData(r)
	redisdb := redis.NewClient(&redis.Options{
		Addr:        redisHost + ":" + redisPort,
		Password:    redisPassword, // no password set
		DB:          0,  // use default DB
		DialTimeout: time.Second ,
	})

	count, err := redisdb.Incr("hits.go").Result()
	if err != nil {
		data["warning"] = fmt.Sprintf("Unable to connect redis server at %s:%s", redisHost,redisPort)
	}else{
		data["message"] = fmt.Sprintf("Hit Count from Redis key hits.go : %v", count)
	}
	writeData(w,r,data)
}

var redisHost = getEnv("REDIS_SERVICE_HOST", "localhost")
var redisPort = getEnv("REDIS_SERVICE_PORT", "6379")
var redisPassword = getEnv("REDIS_SERVICE_PASSWORD", "")
var appVersion = getEnv("APP_VERSION","1.0.0")
var appName = getEnv("APP_NAME","GO_WEB")

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseCFSettings(){
	if value, ok := os.LookupEnv("VCAP_SERVICES"); ok {
		// var serviceData map[string]interface{}

		// if err := json.Unmarshal([]byte (value), &serviceData); err == nil {
		// 	  redisHost = serviceData["rediscloud"][0]["credentials"]["hostname"]  
		// 	  redisPort = serviceData["rediscloud"][0]["credentials"]["port"]  
		// 	  redisPassword = serviceData["rediscloud"][0]["credentials"]["password"]  
		// }
		fmt.Println(value)

	}

}


var tmpl = template.Must(template.ParseFiles("templates/layout.html"))

func main() {
	var host string
	var dir string
	var port string

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.StringVar(&host, "host", "0.0.0.0", "listen host")
	flag.StringVar(&port, "port", "8080", "listen port")
	flag.Parse()

	parseCFSettings();



	r := mux.NewRouter()
	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir(dir)))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	r.HandleFunc("/", handler)
	r.HandleFunc("/die", killHandler)
	r.HandleFunc("/redis", redisHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    host + ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	fmt.Fprintf(os.Stdout, "Server listening %s:%s\n", host, port)
	log.Fatal(srv.ListenAndServe())
}
