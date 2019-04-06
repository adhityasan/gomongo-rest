package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strings"

	"github.com/disintegration/gift"
	"github.com/adhityasan/gomongo-rest/imagehandler"
	"github.com/adhityasan/gomongo-rest/pii"
)

// var cuttariImg = `{"url":"http://cdn2.tstatic.net/batam/foto/bank/images/cut-tari-artis-dan-pembawa-acara-televisi.jpg"}`
// var ersamayoriImg = `{"url":"https://cdns.klimg.com/kapanlagi.com/selebriti/Ersa_Mayori/p/ersa-mayori-025.jpg"}`

// IdentifyByAzure identify endpoint using azure face comparation
// also `pii.Save()` for saving data to database
func IdentifyByAzure(w http.ResponseWriter, r *http.Request) {

	decodedPii, errDecode := pii.DecodeFormPost(r)

	if errDecode != nil {
		log.Println(errDecode)
	}

	piiExist, errExist := decodedPii.Exist()

	if errExist != nil {
		log.Println(errExist)
	}

	asyncProc := make(chan interface{}, 1)

	if piiExist == false {
		go func() {
			id, errSave := decodedPii.Save()
			if errSave != nil {
				log.Println(errSave)
			}
			asyncProc <- id
		}()
	}

	adapter := imagehandler.AzureAdapter{}
	go adapter.Compare(decodedPii.PasfotoKTP.Data, decodedPii.FotoSelfie.Data, asyncProc)
	confidence := <-asyncProc

	w.Write([]byte(fmt.Sprintf("%v", confidence)))
}

// Identify as endpoint
func Identify(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("Your submition is completed. Thank you~ ulalaa~")))
	log.Println("Processing Data")
	piiChan := make(chan interface{}, 1)
	compImgChan := make(chan interface{}, 1)

	if decodedPii, errDecode := pii.DecodeFormPost(r); errDecode == nil {
		log.Println("Processing Images")
		adapter := &imagehandler.AwsAdapter{}
		go adapter.Compare(decodedPii.PasfotoKTP.Data, decodedPii.FotoSelfie.Data, compImgChan)
		confidence := <-compImgChan

		log.Println("Check if NIK " + decodedPii.Nik + " exist")
		_, errExist := decodedPii.Exist()
		if errExist != nil {
			log.Println("Get data from dukcapil")
			formatted := new(pii.Pii)
			data := bytes.NewReader(HitDukcapil(decodedPii.Nik))
			json.NewDecoder(data).Decode(&formatted)
			formatted.FotoKTP = decodedPii.FotoKTP
			formatted.FotoSelfie = decodedPii.FotoSelfie
			formatted.PasfotoKTP = decodedPii.PasfotoKTP
			go func() {
				id, errSave := formatted.Save()
				if errSave != nil {
					fmt.Println(errSave)
				}
				log.Println("Data saved")
				piiChan <- id
			}()
			decodedPii = formatted
		}
		savedID := <-piiChan
		log.Println("Processing data finished.\nResult :" + fmt.Sprint(confidence) + "\n" + fmt.Sprint(savedID))
	} else {
		fmt.Println(errDecode)
	}
}

// HitDukcapil function hit Dukcapil simulator API from Docker IP
func HitDukcapil(nik string) (data []byte) {
	dukcapilURI := "http://192.168.99.100:5000"
	param := `{"NIK": "` + nik + `"}`
	var decoded struct {
		Content []map[string]interface{} `json:"content"`
	}
	req, _ := http.NewRequest("POST", dukcapilURI, strings.NewReader(param))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	if resp, err := client.Do(req); err == nil {
		json.NewDecoder(resp.Body).Decode(&decoded)
		if data, err = json.Marshal(decoded.Content[0]); err == nil {
			return data
		}
		return []byte(fmt.Sprintln(err))
	}
	return nil
}

// Aisatsu sample get request for testing purpose
func Aisatsu(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")
	if name == "" {
		name = "Guest"
	}
	log.Printf("こんにちは %s'san\n", name)
	w.Write([]byte(fmt.Sprintf("Hello, %s'san\n", name)))
}

// DoOCR for KTP Image as first endpoint
func DoOCR(w http.ResponseWriter, r *http.Request) {
	// Sementara pakai buffer, next pakai Pii untuk return objectID
	bufKTP := bytes.NewBuffer(nil)
	imgChan := make(chan interface{})
	if imgKTP, headerKTP, err := r.FormFile("foto_ktp"); err == nil {
		fmt.Println("Reading Image " + headerKTP.Filename)
		defer imgKTP.Close()
		img, _ := jpeg.Decode(imgKTP)
		g := gift.New(
			gift.Contrast(20),
			gift.Grayscale(),
		)
		imgEnhance := image.NewRGBA(g.Bounds(img.Bounds()))
		g.Draw(imgEnhance, img)
		err = jpeg.Encode(bufKTP, imgEnhance, nil)
		adapter := &imagehandler.AwsAdapter{}
		go adapter.Read(bufKTP.Bytes(), imgChan)
		ocrRes := <-imgChan
		w.Write([]byte(fmt.Sprintf("%v", ocrRes)))
	} else {
		w.Write([]byte(fmt.Sprintln(err)))
	}
}
