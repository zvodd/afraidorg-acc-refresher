package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/viper"
)

func configure() {
	viper.SetConfigType("props")
	fh, err := os.Open("secret.toml")
	if err != nil {
		log.Fatalln("Config file not found")
	}
	viper.ReadConfig(fh)
}

func main() {
	configure()
	username := viper.Get("username").(string)
	password := viper.Get("password").(string)

	client := &http.Client{}

	resp, err := client.Get("https://freedns.afraid.org/profile/")
	if err != nil {
		log.Fatalln("Http Request Error (%s): %v", resp.Request.URL, err)
	}

	from := resp.Request.URL.Query().Get("from")
	// log.Println("form:", from)
	// inspectResp(resp)

	resp2, err := client.PostForm("https://freedns.afraid.org/zc.php?step=2",
		url.Values{
			"username": []string{username},
			"password": []string{password},
			"submit":   []string{"Login"},
			"remote":   []string{""},
			"from":     []string{from},
			"action":   []string{"auth"},
		})
	if err != nil {
		log.Fatalln("Http Request Error (%s): %v", resp.Request.URL, err)
	}

	log.Printf("%#v", resp2)
	if resp2.StatusCode == 200 {
		IFTTT_WebHook_withconf("Afraid.org Success", "Account refresh succeeded")
	} else {
		IFTTT_WebHook_withconf("Afraid.org Failure", "Account refresh failed!")
	}
}

func inspectResp(resp *http.Response) {
	// log.Printf("%#v \n")
	log.Printf("%#v \n\n", resp.Header)
	log.Printf("%#v \n\n", resp.Request)
	log.Printf("UserAgent %v \n\n", resp.Request.UserAgent())
	log.Printf("RequestURI %v \n\n", resp.Request.RequestURI)
	log.Printf("URL %v \n\n", resp.Request.URL.String())

	// from := resp.Request.URL.Query().Get("from")
	// log.Printf("from %v \n\n", from)

	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)

}

func IFTTT_WebHook_withconf(title string, msg string) {
	key := viper.Get("ifttt_key").(string)
	ev := viper.Get("ifttt_event_name").(string)
	err := IFTTT_webhook(ev, key, title, msg)
	if err != nil {
		log.Println("IFTTT Webhook failure?", err)
	}
}

func IFTTT_webhook(eventName string, iftttKey string, v1 string, v2 string) error {
	url := "https://maker.ifttt.com/trigger/" + eventName + "/with/key/" + iftttKey

	vs := map[string]string{
		"value1": v1,
		"value2": v2,
	}
	jsons, err := json.Marshal(vs)
	if err != nil {
		log.Fatalln("JSON Marshal:", err)
	}

	buf := bytes.NewReader(jsons)

	_, err = http.Post(url, "application/json", buf)

	return err
}
