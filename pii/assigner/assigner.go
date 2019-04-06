package assigner

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/adhityasan/gomongo-rest/config"
	"github.com/adhityasan/gomongo-rest/pii"
)

// Assigner asdasdasd
func Assigner(nik string, piitostore *pii.Pii) error {
	requstBody, err := json.Marshal(map[string]string{
		"NIK": nik,
	})

	req, err := http.NewRequest("POST", config.Of.Dukcapil.Endpoint, bytes.NewBuffer(requstBody))
	if err != nil {
		log.Fatalln(err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	resp, errDo := client.Do(req)
	if errDo != nil {
		log.Fatalln(err)
		return err
	}

	defer resp.Body.Close()

	var decoded struct {
		Content []interface{}
	}

	errDecode := json.NewDecoder(resp.Body).Decode(&decoded)
	if errDecode != nil {
		log.Println("errdecode", errDecode)
	}

	thecontents, errMarshalContent := json.Marshal(decoded.Content[0])
	if errMarshalContent != nil {
		log.Println(errMarshalContent)
	}

	json.Unmarshal(thecontents, &piitostore)

	return nil
}
