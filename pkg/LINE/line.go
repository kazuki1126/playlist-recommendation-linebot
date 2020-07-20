package line

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	spotify "github.com/kazuki1126/playlist-recommendation-linebot/pkg/spotify"

	"github.com/line/line-bot-sdk-go/linebot"
)

type lineAuthorizaiton struct {
	Secret string
	Token  string
}

var LineAuth = lineAuthorizaiton{
	Secret: os.Getenv("LINE_SECRET"),
	Token:  os.Getenv("LINE_TOKEN"),
}

type btnTemplateArgs struct {
	actions  []linebot.TemplateAction
	title    string
	text     string
	altText  string
	imageURL string
}

var action1 = linebot.NewPostbackAction("チルな気分", spotify.ChillMusic, "", "おっさんチルしたいよ")
var action2 = linebot.NewPostbackAction("パーティーしたい気分", spotify.PartyMusic, "", "おっさんパーティーしたいよ")
var action3 = linebot.NewPostbackAction("おっさんの時代の洋楽", spotify.OldMusic, "", "おっさんの時代の曲聞かせてよ")
var action4 = linebot.NewPostbackAction("なんでもいい", spotify.AnyMusic, "", "おっさん、なんでもいいから聞かせてよ")
var action5 = linebot.NewPostbackAction("フェス系 (EDM)がいい", spotify.EDM, "", "おっさんフェス系聞いてぶち上がりたい")
var action6 = linebot.NewPostbackAction("ヒップホップ", spotify.Hiphop, "", "おっさんラップ聞きたいよ")

var oyajiNoSerifu = "どうじゃこいつら聞いてみい"
var oyajiErr error = errors.New("すまんがわしは疲れた。また別の日にしてくれ。")

func SendReply(w http.ResponseWriter, req *http.Request) {
	bot, err := linebot.New(LineAuth.Secret, LineAuth.Token)
	if err != nil {
		log.Fatalln(err)
	}
	events, err := bot.ParseRequest(req)
	if err != nil {
		log.Fatalln(err)
	}
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			args := btnTemplateArgs{
				actions:  []linebot.TemplateAction{action1, action2, action3, action4},
				title:    "どんな感じの曲がいいか言ってみい",
				text:     "俺がプレイリストおすすめしたるわい",
				altText:  "今使ってるデバイスじゃ見れんぞ、スマホを見るんじゃ",
				imageURL: "https://i1.wp.com/liveforlivemusic.com/wp-content/uploads/2016/01/musicbrain.jpg?resize=610%2C390&ssl=1",
			}
			tplMessage := createBtnTpl(args)
			if _, err = bot.ReplyMessage(event.ReplyToken, tplMessage).Do(); err != nil {
				log.Print(err)
			}
		case linebot.EventTypePostback:
			switch event.Postback.Data {
			case spotify.PartyMusic:
				args := btnTemplateArgs{
					actions:  []linebot.TemplateAction{action5, action6},
					title:    "もうちょい詳しく言ってみい",
					text:     "いいやつおすすめしたるわい",
					altText:  "今使ってるデバイスじゃ見れんぞ、スマホを見るんじゃ",
					imageURL: "https://festivalsherpa-wpengine.netdna-ssl.com/wp-content/uploads/2016/06/Dad-at-music-festival-by-vice.com_.jpg",
				}
				tplMessage := createBtnTpl(args)
				if _, err = bot.ReplyMessage(event.ReplyToken, tplMessage).Do(); err != nil {
					log.Print(err)
				}
			default:
				fmt.Println(event.Postback.Data)
				messages, err := returnPlayListsAsMsg(event.Postback.Data)
				if err != nil {
					fmt.Println(err)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(oyajiErr.Error())).Do(); err != nil {
						log.Print(err)
					}
				} else {
					if _, err = bot.ReplyMessage(event.ReplyToken, messages...).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}

func createBtnTpl(args btnTemplateArgs) *linebot.TemplateMessage {
	btnTemplate := linebot.NewButtonsTemplate(args.imageURL, args.title, args.text, args.actions...)
	btnTemplate = btnTemplate.WithImageOptions(linebot.ImageAspectRatioTypeRectangle, linebot.ImageSizeTypeCover, "#FFFFFF")
	tplMessage := linebot.NewTemplateMessage(args.altText, btnTemplate)
	return tplMessage
}

func returnPlayListsAsMsg(musicType string) ([]linebot.SendingMessage, error) {
	if musicType == spotify.AnyMusic {
		randomCategory := spotify.GetRandomCategory(spotify.PlaylistCategories)
		playlistURLs, err := spotify.GetPlayLists(randomCategory)
		if err != nil {
			return nil, err
		}
		var messages = []linebot.SendingMessage{}
		messages = append(messages, linebot.NewTextMessage(oyajiNoSerifu))
		for _, playlistURL := range playlistURLs {
			messages = append(messages, linebot.NewTextMessage(playlistURL))
		}
		return messages, nil
	}
	playlistURLs, err := spotify.GetPlayLists(musicType)
	if err != nil {
		return nil, err
	}
	var messages = []linebot.SendingMessage{}
	messages = append(messages, linebot.NewTextMessage(oyajiNoSerifu))
	for _, playlistURL := range playlistURLs {
		messages = append(messages, linebot.NewTextMessage(playlistURL))
	}
	return messages, nil
}
