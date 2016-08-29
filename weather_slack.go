package main

//packageをimport
import (
        "bytes"
        "net/http"
        "net/url"
        "fmt"
        "encoding/xml"
        "io/ioutil"
        "log"
)
//変数宣言
var(
        //SlackAPIトークン
        token  string = ""
        //SlackApiのメッセージリクエストURL
        apiUrl string = "https://slack.com/api/chat.postMessage"
        //livedoor天気予報のRSS(東京)
        FEED_URL string = "http://weather.livedoor.com/forecast/rss/area/130010.xml"

)

//構造体
type WeatherHack struct {
    Title string `xml:"channel>title"`
    Description []string `xml:"channel>item>description"`
}


func main() {
        //getWeather関数にRSSのURLを渡す
        wh, err := getWeather(FEED_URL)

        //error処理
        if err != nil {
            log.Fatalf("Log: %v", err)
            return
        }

        //URLに生成するデータを追加
        data := url.Values{}
        data.Set("token",token)
        data.Add("channel","")
        data.Add("username","天気予報Bot")
        data.Add("icon_url","")
        data.Add("text", wh.Description[2])
        fmt.Println(data)

        //http.clientを使用してGETリクエスト送信
        client := &http.Client{}
        r, _ := http.NewRequest("POST",  fmt.Sprintf("%s",apiUrl), bytes.NewBufferString(data.Encode()))
        r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

        resp, _ := client.Do(r)
        fmt.Println(resp.Status)
}

//天気情報を取得する
func getWeather(feed string) (p *WeatherHack, err error) {

        //http.Getメソッドを利用してGETリクエストを送る
        res, err := http.Get(feed)
        if err != nil {
            return nil, err
        }

        b, err := ioutil.ReadAll(res.Body)
        if err != nil {
            return nil, err
        }
        wh := new(WeatherHack)
        err = xml.Unmarshal(b, &wh)

        return wh, err
}