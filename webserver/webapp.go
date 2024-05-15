package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/api/option"
)

var (
	username = os.Getenv("DB_USER") //
	// username = "root"
	host = os.Getenv("DB_ADDR")
	// host = "127.0.0.1:3006"
	db_name = os.Getenv("DB_NAME")
	// db_name = "test"
	db_password = os.Getenv("DB_PWD")
	// db_password = "1234"
	table_name = os.Getenv("DB_TABLE_NAME")
	// table_name = "test_table"
	bucketName = os.Getenv("GCP_CLOUD_STORAGE")

	serviceAccountKeyFile = os.Getenv("GCP_SERVICE_ACCOUNT")
	// bucketName    = "cdncontent12"
	default_query = "SELECT * FROM files"
)

type File struct {
	ID          string `json:"id"`
	ContentType string `json:"content_type"`
	FileName    string `json:"file_name"`
	Path        string `json:"path"`
}

var files_list []File

func Insert(content_type string, file_name string, path_to_file string) {
	db_cred := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, db_password, host, db_name)

	db := connection(db_cred)

	query := fmt.Sprintf("INSERT INTO %s (content_type, file_name, path_to_file) VALUES ('%s', '%s','%s')", table_name, content_type, file_name, path_to_file)
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
	}
	defer stmt.Close()
	stmt.ExecContext(ctx)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
	}
}

func QueryDB() []File {
	files_list := []File{}

	db_cred := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, db_password, host, db_name)

	db := connection(db_cred)

	rows, err := db.Query(default_query)

	for rows.Next() {
		var file File
		var id int
		var content_type string
		var file_name string
		var path_to_file string
		// Add more variables according to your table columns

		err := rows.Scan(&id, &content_type, &file_name, &path_to_file) // Scan the values from the row into variables
		if err != nil {
			panic(err.Error())
		}

		// Do something with the values, for example, print them
		fmt.Println(id, content_type, file_name, path_to_file)

		file.ID = fmt.Sprint(id)
		file.ContentType = content_type
		file.FileName = file_name
		file.Path = path_to_file

		files_list = append(files_list, file)
	}

	if err != nil {
		log.Printf("Error %s getting rows", err)
	}
	defer rows.Close()
	return files_list
}

func connection(db_cred string) *sql.DB {
	db, err := sql.Open("mysql", db_cred)
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
	}
	return db
}

func isImage(contentType string) bool {
	// Check if the content type indicates an image
	return contentType == "image/jpeg" || contentType == "image/png" || contentType == "image/gif"
}

func isVideo(contentType string) bool {
	// Check if the content type indicates a video
	return contentType == "video/mp4" || contentType == "video/mpeg" || contentType == "video/quicktime"
}

func backend(c *gin.Context) {
	c.HTML(http.StatusOK, "backend.html", gin.H{})
}

func home(c *gin.Context) {
	lt := QueryDB()
	files_list = lt
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Files": lt,
	})
}

func insertIntoBucket(filename string, file multipart.File, content_type string) {
	ctx := context.Background()
	cli, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceAccountKeyFile))
	bkt := cli.Bucket(bucketName)
	if err != nil {
		fmt.Errorf("Failed to create client: %v", err)
	}
	defer cli.Close()

	obj := bkt.Object(fmt.Sprintf("media/%s/%s", content_type, filename))
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, file); err != nil {
		fmt.Errorf("Failed to copy to bucket: %v", err)
	}
	// Close, just like writing a file. File appears in GCS after
	if err := w.Close(); err != nil {
		fmt.Errorf("Failed to close: %v", err)
	}
}

func uploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")

	log.Println(header.Filename)
	if err != nil {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"error": "Failed to upload file",
		})
	}
	ct_type := header.Header.Get("Content-Type")

	if isImage(ct_type) {
		insertIntoBucket(header.Filename, file, "image")
		Insert("image", header.Filename, "https://storage.googleapis.com/cdn-content-storage/media/image/"+header.Filename)

	} else if isVideo(ct_type) {
		insertIntoBucket(header.Filename, file, "video")
		Insert("video", header.Filename, "https://storage.googleapis.com/cdn-content-storage/media/video/"+header.Filename)
	} else {
		insertIntoBucket(header.Filename, file, "file")
		Insert("file", header.Filename, "https://storage.googleapis.com/cdn-content-storage/media/file/"+header.Filename)
	}

	c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", header.Filename))
}

func getContent(c *gin.Context) {
	id := c.Param("id")
	for files := range files_list {
		if files_list[files].ID == id {
			c.HTML(http.StatusOK, "image.html", gin.H{
				"contentURL":  files_list[files].Path,
				"contentType": files_list[files].ContentType,
				"ID":          id,
			})
		}
	}
}

func dbDeleteFile(id string) {
	db_cred := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, db_password, host, db_name)

	db := connection(db_cred)

	query := fmt.Sprintf("DELETE FROM %s WHERE id=%s;", table_name, id)
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
	}
	defer stmt.Close()
	stmt.ExecContext(ctx)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
	}
}

func bucketDeleteFile(id string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(serviceAccountKeyFile))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	for files := range files_list {
		if files_list[files].ID == id {
			objectName := fmt.Sprintf("media/%s/%s", files_list[files].ContentType, files_list[files].FileName)
			// Deletes the object.
			err = client.Bucket(bucketName).Object(objectName).Delete(ctx)
			if err != nil {
				log.Fatalf("Failed to delete object: %v", err)
			}
			fmt.Printf("Object %s deleted from bucket %s\n", objectName, bucketName)
		}
	}
}

func deleteContent(c *gin.Context) {
	id := c.Param("id")
	dbDeleteFile(id)
	bucketDeleteFile(id)
	c.String(http.StatusOK, fmt.Sprintf("delete id ===> '%s'", id))
}

func main() {
	r := gin.Default()
	store := persistence.NewInMemoryStore(time.Second)
	r.LoadHTMLGlob("templates/*")

	r.MaxMultipartMemory = 20 << 20 // 20 MiB

	r.GET("/", home)
	r.GET("/backend", backend)
	r.POST("/backend", uploadFile)
	r.GET("/:content/:id", getContent)
	r.POST("/delete-content/:id", deleteContent)

	r.GET("/api", cache.CachePage(store, time.Minute, func(c *gin.Context) {
		lt := QueryDB()
		files_list := lt
		c.IndentedJSON(http.StatusCreated, files_list)
	}))
	if err := http.ListenAndServeTLS(":443", "./certificates/server.crt", "./certificates/server.key", r); err != nil {
		panic(err)
	}
	r.Run() // listen and serve on 0.0.0.0:8080
}
