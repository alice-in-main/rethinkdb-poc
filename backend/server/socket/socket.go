package socket

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"schemas"
	"server/queries"
	"server/utils"
	"strings"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", utils.GetenvWithDefault("SERVER_HOST", "localhost:8082"), "http service address")
var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

func Start() {
	flag.Parse()
	http.HandleFunc("/zoom", streamEntitiesInZoom)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func streamEntitiesInZoom(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	zoom, err := parseZoomFromRequest(r)
	if err != nil {
		log.Print("bad zoom:", err)
		return
	}

	stream := queries.GetAllInZoom(zoom)
	for m := range stream {
		c.WriteJSON(m)
	}
}

func parseZoomFromRequest(r *http.Request) (schemas.Zoom, error) {

	rawObj := strings.Split(r.URL.RawQuery, "=")[1]
	rawObj = strings.ReplaceAll(rawObj, "%22", "\"")

	var parsedZoom schemas.Zoom
	err := json.Unmarshal([]byte(rawObj), &parsedZoom)
	if err != nil {
		return schemas.Zoom{}, err
	}

	return parsedZoom, nil
}