package imagehandler

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/adhityasan/gomongo-rest/config"
	"github.com/adhityasan/gomongo-rest/imagehandler/aws"
	"github.com/adhityasan/gomongo-rest/imagehandler/azure"
)

// AzureAdapter for Azure IP endpoint setup
type AzureAdapter struct {
	Endpoint *azure.Endpoint
}

// Compare two images using Azure
func (e *AzureAdapter) Compare(img1 interface{}, img2 interface{}, ch chan interface{}) {
	e.Endpoint = &azure.Endpoint{
		URL: config.Of.Azure.Endpoint,
		Key: config.Of.Azure.APIKey,
	}
	faceIDChan := make(chan interface{})
	go e.Endpoint.FaceID(img1, faceIDChan)
	fmt.Println("Processing Image 1")
	go e.Endpoint.FaceID(img2, faceIDChan)
	fmt.Println("Processing Image 2")
	faceID1 := <-faceIDChan
	faceID2 := <-faceIDChan
	jsonParam := `{"faceId1":"` + faceID1.(string) + `","faceId2":"` + faceID2.(string) + `"}`
	if res, err := e.Endpoint.GetConfidence(jsonParam); err == nil {
		ch <- res
	} else {
		fmt.Println(err)
	}
}

// AwsAdapter adapt AWS function
type AwsAdapter struct {
	Gateway *aws.Gateway
}

// Compare of two images using AWS
func (b *AwsAdapter) Compare(img1 []byte, img2 []byte, ch chan interface{}) {
	b.Gateway = &aws.Gateway{
		Region:    config.Of.Aws.Region,
		KeyID:     config.Of.Aws.KeyID,
		SecretKey: config.Of.Aws.SecretKey,
	}
	p := &aws.CompareParam{
		ImgKTP:    img1,
		ImgSelfie: img2,
	}
	if res, err := b.Gateway.Compare(p); err == nil {
		finalres := `{"Confidence": "` + fmt.Sprintf("%f", res) + `"}`
		ch <- finalres
	}
}

func (b *AwsAdapter) Read(img []byte, ch chan interface{}) {
	b.Gateway = &aws.Gateway{
		Region:    config.Of.Aws.Region,
		KeyID:     config.Of.Aws.KeyID,
		SecretKey: config.Of.Aws.SecretKey,
	}
	if res, err := b.Gateway.Read(img); err == nil {
		extracted := extract(res)
		ch <- extracted
	}
}

func extract(str string) interface{} {
	res := make(map[string]string)
	str = strings.TrimSpace(str)
	str = strings.ReplaceAll(str, ":", "")
	reg := regexp.MustCompile(`(?mi)nik|nama|tempat|tgl|lahir|jenis|kelamin|gol.|gol|darah|alainat|alamat|agama|pekerjaan|status|perkawinan`)
	str = reg.ReplaceAllString(str, "")
	strSplit := strings.Split(str, "\n")
	for i, v := range strSplit {
		if v != "" {
			strSplit[i] = strings.TrimSpace(v)
		}
	}
	res["province"] = strSplit[0]
	res["city"] = strSplit[1]
	res["nik"] = strSplit[2]
	res["name"] = strSplit[3]
	fmt.Println(res)
	resJSON, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
	}
	return string(resJSON)
}
