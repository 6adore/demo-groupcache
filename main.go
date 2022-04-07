// Simple groupcache example: https://github.com/golang/groupcache
// Running 3 instances:
// go run main.go -addr=:8080 -pool=http://127.0.0.1:8080,http://127.0.0.1:8081,http://127.0.0.1:8082
// go run main.go -addr=:8081 -pool=http://127.0.0.1:8081,http://127.0.0.1:8080,http://127.0.0.1:8082
// go run main.go -addr=:8082 -pool=http://127.0.0.1:8082,http://127.0.0.1:8080,http://127.0.0.1:8081
// Testing:
// curl localhost:8080/customer?custid=1

package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	conf "github.com/6adore/demo-groupcach/config"
	"github.com/6adore/demo-groupcach/database"
	"github.com/golang/groupcache"
)

var Group *groupcache.Group
var (
	configFile = flag.String("f", "config/config.yml", "the config file")
	addr       = flag.String("addr", "8080", "server address")
	peers      = flag.String("pool", "http://localhost:8080", "server pool list")
)

func main() {
	flag.Parse()
	conf.MustLoadConfig(*configFile)

	db := database.InitDB()
	if db == nil {
		fmt.Println("DB IS NOT initialized")
		return
	}
	defer db.Close()

	Group = groupcache.NewGroup("foo", 64<<20, groupcache.GetterFunc(
		func(_ context.Context, key string, dest groupcache.Sink) error {
			log.Printf("Looking up from groupcache\n")
			id, _ := strconv.Atoi(key)
			row := db.QueryRow("SELECT cust_name FROM customers WHERE cust_id = ?", id)
			var cust_name string
			if err := row.Scan(&cust_name); err != nil && err == sql.ErrNoRows {
				return errors.New("data not found")
			}
			fmt.Println("name: ", cust_name)
			var v = []byte(cust_name)
			dest.SetBytes(v)
			return nil
		},
	))
	http.HandleFunc("/customer", func(w http.ResponseWriter, r *http.Request) {
		custid := r.FormValue("custid")
		fmt.Println("customer id: ", custid)
		var b []byte
		err := Group.Get(context.Background(), custid, groupcache.AllocatingByteSliceSink(&b))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		w.Write([]byte("\n"))
		w.Write(b)
		w.Write([]byte("\n"))
	})

	p := strings.Split(*peers, ",")
	pool := groupcache.NewHTTPPool(p[0])
	pool.Set(p...)
	http.ListenAndServe(*addr, nil)
}
