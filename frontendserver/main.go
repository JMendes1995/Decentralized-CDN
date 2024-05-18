package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var (
	files_list              []File
	mc                      *memcache.Client
	private_ip              = os.Getenv("PRIVATE_IP")
	memcache_local_ip       = os.Getenv("LOCAL_MEMCACHED_IP")
	backend_address         = os.Getenv("BACKEND_ADDRESS")
	serviceAccountKeyFile   = os.Getenv("GCP_SERVICE_ACCOUNT")
	projectID               = os.Getenv("GCP_PROJECT_ID")
	region                  = os.Getenv("GCP_REGION")
	items_cached            []string
	memecached_remote_nodes []string
	mutex                   sync.Mutex
)

type File struct {
	ID          string `json:"id"`
	ContentType string `json:"content_type"`
	FileName    string `json:"file_name"`
	Path        string `json:"path"`
}

func requestBackend(requesturl string) []File {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(requesturl)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	// Decode byte slice into a map
	err = json.NewDecoder(resp.Body).Decode(&files_list)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
	} // Print decoded JSON data
	return files_list
}

func home(c *gin.Context) {
	lt := requestBackend(fmt.Sprintf("https://%s/api", backend_address))
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Files": lt,
	})
}

func publishMessage(message_data string, topicID string) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(serviceAccountKeyFile))
	if err != nil {
		fmt.Errorf("pubsub.NewClient: %w", err)
	}
	defer client.Close()
	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(message_data),
		Attributes: map[string]string{
			"origin":   "golang",
			"username": "gcp",
		},
	})
	// Block until the resul t is returned and a server-generated
	// ID is returned for the published message.
	id, err := result.Get(ctx)
	if err != nil {
		fmt.Errorf("Get: %w", err)
	}
	fmt.Println("published message: msg ID: ", id)
}

func getContent(c *gin.Context) {
	mc = memcache.New(fmt.Sprintf("%s:11211", memcache_local_ip))

	id := c.Param("id")
	content := c.Param("content")
	fmt.Println("/" + content + "/" + id)

	for files := range files_list {
		if files_list[files].ID == id {
			item, err := mc.Get("/" + content + "/" + id)
			if err == nil {
				c.Data(http.StatusOK, "*", item.Value)
			} else {
				fmt.Println("Content not loaded in cached")
				cont, err := http.Get(files_list[files].Path)
				if err != nil {
					c.String(http.StatusInternalServerError, fmt.Sprintf("error getting content '%s' uploaded!", files_list[files].FileName))
				}
				defer cont.Body.Close()
				contentBytes, err := ioutil.ReadAll(cont.Body)
				if err != nil {
					c.String(http.StatusInternalServerError, fmt.Sprintf("error reading content'%s' uploaded!", files_list[files].FileName))
				}
				mc.Set(&memcache.Item{
					Key:        "/" + content + "/" + id,
					Value:      contentBytes,
					Expiration: 3600,
				})

				c.Data(http.StatusOK, cont.Header.Get("Content-Type"), contentBytes)
				mutex.Lock()
				items_cached = append(items_cached, "/"+content+"/"+id)
				defer mutex.Unlock()
				fmt.Println(memecached_remote_nodes)

				for i := range memecached_remote_nodes {
					mc_remote := memcache.New(fmt.Sprintf("%s:11211", memecached_remote_nodes[i]))
					err := mc_remote.Set(&memcache.Item{
						Key:        "/" + content + "/" + id,
						Value:      contentBytes,
						Expiration: 3600,
					})
					if err != nil {
						fmt.Println("error replication the cache to node " + memecached_remote_nodes[i])
					}

				}
			}
		}
	}
}

func webserver() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", home)
	r.GET("/:content/:id", getContent)

	// Run Gin server with TLS
	if err := http.ListenAndServeTLS(":443", "./certificates/server.crt", "./certificates/server.key", r); err != nil {
		panic(err)
	}
}

func initNotificationListener(topicID string) {
	ctx := context.Background()

	client, _ := pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(serviceAccountKeyFile))

	sub := client.Subscription(topicID + "-sub")

	sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", string(msg.Data))
		msg.Ack()

		if string(msg.Data) != memcache_local_ip {

			mc_local := memcache.New(fmt.Sprintf("%s:11211", memcache_local_ip))

			mc_remote := memcache.New(fmt.Sprintf("%s:11211", string(msg.Data)))
			mutex.Lock()
			for item := range items_cached {
				_, err_remote_check := mc_remote.Get(items_cached[item])
				if err_remote_check == nil {
					fmt.Println("item already in cached in the remote node->" + items_cached[item])
				} else {
					content_data, err_check := mc_local.Get(items_cached[item])
					if err_check == nil {
						fmt.Println("sending cached item to new host ->" + items_cached[item])
						msg.Ack()
						// publishMessage(string(msg.Data))
						mc_remote.Set(&memcache.Item{
							Key:        items_cached[item],
							Value:      content_data.Value,
							Expiration: 3600,
						})
					}
				}
			}
			if !isInArray(memecached_remote_nodes, string(msg.Data)) {
				memecached_remote_nodes = append(memecached_remote_nodes, string(msg.Data))
			}
			fmt.Println(memecached_remote_nodes)
			mutex.Unlock()

		} else {
			fmt.Println("got own message retransmiting")
			msg.Ack()
			publishMessage(string(msg.Data), "init_memcache-"+region)

		}
	})
}

func isInArray(arr []string, target string) bool {
	for _, value := range arr {
		if value == target {
			return true
		}
	}
	return false
}

func healthCheck() {
	for {
		publishMessage(memcache_local_ip, "init_memcache-"+region)
		time.Sleep(2 * time.Minute)
	}
}

func main() {
	// go publishMessage(memcache_local_ip, "init_memcache")
	go healthCheck()
	go initNotificationListener("init_memcache-" + region)

	go func() {
		r := gin.Default()
		r.LoadHTMLGlob("templates/*")
		r.GET("/client", home)
		r.GET("/:content/:id", getContent)
		// Run Gin server with TLS
		if err := http.ListenAndServeTLS(":443", "./certificates/server.crt", "./certificates/server.key", r); err != nil {
			panic(err)
		}
	}()

	// go notification_listener()

	select {}
}
