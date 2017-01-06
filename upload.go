package main

import (
    "fmt"
    "time"
	"crypto/md5"
	"io"
	"strconv"
	"os"
	"text/template"
	"net/http"
	"github.com/dchest/uniuri"
	"github.com/disintegration/imaging"	
	"runtime"
	"encoding/json"
	"github.com/juju/errgo"
	"strings"
)

func upload(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
   } else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
	    }
		defer file.Close()
		
		contentType := handler.Header["Content-Type"][0]

		if !isAllowedContentType(contentType) {
			fmt.Println("Only JPEG or PNG files allowed")
		}
		
		filename := generateRandomFilename(contentType)
		
		photoname := strings.Split(filename, ".")
		
		f, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		
		runtime.GOMAXPROCS(runtime.NumCPU())
		
		// load original image
        img, err := imaging.Open(filename)

        if err != nil {
			fmt.Println(err)
            return
        }
		dstimg := imaging.Resize(img, 128, 0, imaging.Box)

        // save resized image
        err = imaging.Save(dstimg, photoname[0]+"_resized.jpg")

		if err != nil {
			fmt.Println(err)
			return
		}
		
		photo := make(map[string]string)
		photo["title"] = handler.Filename
		photo["url"] = "http://localhost/upload/"+photoname[0]+".jpg"
		photo["resized_url"] = "http://localhost/upload/"+photoname[0]+"_resized.jpg"
		
        renderJSON(w, photo, http.StatusCreated)	
	}
}

func main() {
    //http.HandleFunc("/", handler)
	http.HandleFunc("/upload", upload)
    http.ListenAndServe("localhost:8080", nil)
}

func generateRandomFilename(contentType string) string {

	var ext string

	switch contentType {
	case "image/jpeg":

		ext = ".jpg"
	case "image/png":
		ext = ".png"

	case "image/gif":
		ext = ".gif"
	}

	return uniuri.New() + ext
}
var allowedContentTypes = []string{
		"image/png",
		"image/jpeg",
		"image/gif"}
func isAllowedContentType(contentType string) bool {
	for _, value := range allowedContentTypes {
		if contentType == value {
			return true
		}
	}

	return false
}

func renderJSON(w http.ResponseWriter, value interface{}, status int) error {
	body, err := json.Marshal(value)
	if err != nil {
		return errgo.Mask(err)
	}
	return writeBody(w, body, status, "application/json")
}

func writeBody(w http.ResponseWriter, body []byte, status int, contentType string) error {
	w.Header().Set("Content-Type", contentType+"; charset=UTF8")
	w.Header().Set("Content-Length", strconv.Itoa(len(body)))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	_, err := w.Write(body)
	return errgo.Mask(err)
}