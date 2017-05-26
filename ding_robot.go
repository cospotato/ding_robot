package ding_robot

import (
	"strings"
	"encoding/json"
	"net/http"
	"bytes"
	"io/ioutil"
	"errors"

)

const BaseSendURL = "https://oapi.dingtalk.com/robot/send?access_token={ACCESS_TOKEN}"
const JSONType = "application/json"

type MessageType string
type Orientation string
type AvatarState string

const (
	TypeText MessageType = "text"
	TypeLink MessageType = "link"
	TypeMarkdown MessageType = "markdown"
	TypeActionCard MessageType = "actionCard"
	TypeFeedCard MessageType = "feedCard"

	OrientationVertical Orientation = "0"
	OrientationHorizon Orientation = "1"

	ShowAvatar AvatarState = "0"
	HideAvatar AvatarState = "1"
)

//DingMessage 钉钉机器人消息
type DingMessage struct {
	Type MessageType `json:"msgtype"`
	Text TextElement `json:"text"`
	Link LinkElement `json:"link"`
	Markdown MarkdownElement `json:"markdown"`
	ActionCard ActionCardElement `json:"actionCard"`
	FeedCard FeedCardElement `json:"feedCard"`
	At AtElement `json:"at"`
}

type AtElement struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll bool `json:"isAtAll"`
}

// TextElement 文本元素
type TextElement struct {
	Content string `json:"content"`
}

// LinkElement 链接元素
type LinkElement struct {
	Title string `json:"title"` 		// 消息标题
	Text string `json:"text"`		// 消息内容。如果太长只会部分展示
	MessageURL string `json:"messageUrl"`	// 点击消息跳转的URL
	PictureURL string `json:"picUrl"`	// 图片URL
}

// LinkElement 链接元素
type FeedLinkElement struct {
	Title string `json:"title"` 		// 消息标题
	Text string `json:"text"`		// 消息内容。如果太长只会部分展示
	MessageURL string `json:"messageURL"`	// 点击消息跳转的URL
	PictureURL string `json:"picURL"`	// 图片URL
}

// MarkdownElement Markdown元素
type MarkdownElement struct {
	Title string `json:"title"`		// 首屏会话透出的展示内容
	Text string `json:"text"`		// markdown格式的消息
}

// ActionCardElement ActionCard元素
type ActionCardElement struct {
	Title string `json:"title"`
	Text string `json:"text"`
	SingleTitle string `json:"singleTitle"`
	SingleURL string `json:"singleURL"`
	ButtonOrientation Orientation `json:"btnOrientation"`
	Avatar AvatarState `json:"hideAvatar"`
	Buttons []ButtonElement `json:"btns"`
}

// actionCardBuilder ActionCard构造器
type actionCardBuilder struct {
	actionCard ActionCardElement
}

func NewActionCardBuilder(title string, text string, buttonOrientation Orientation, avatarState AvatarState) *actionCardBuilder {
	return &actionCardBuilder{
		ActionCardElement{
			Title: title,
			Text: text,
			ButtonOrientation: buttonOrientation,
			Avatar: avatarState,
			Buttons: make([]ButtonElement, 0),
		},
	}
}

func (builder *actionCardBuilder) SingleButton(title string, URL string) *actionCardBuilder {
	builder.actionCard.SingleTitle = title
	builder.actionCard.SingleURL = URL
	return builder
}

func (builder *actionCardBuilder) Button(title string, URL string) *actionCardBuilder {
	builder.actionCard.Buttons = append(builder.actionCard.Buttons, ButtonElement{
		Title: title,
		ActionURL: URL,
	})
	return builder
}

func (builder *actionCardBuilder) Build() ActionCardElement {
	return builder.actionCard
}

// ButtonElement 按钮元素
type ButtonElement struct {
	Title string `json:"title"`
	ActionURL string `json:"actionURL"`
}

// FeedCardElement 图文元素
type FeedCardElement struct {
	Links []FeedLinkElement `json:"links"`
}

type feedCardBuilder struct {
	feedCard FeedCardElement
}

func NewFeedCardBuilder() *feedCardBuilder {
	return &feedCardBuilder{
		feedCard: FeedCardElement{
			Links:make([]FeedLinkElement,0),
		},
	}
}

func (builder *feedCardBuilder) Link(title string, messageURL string, pictureURL string) *feedCardBuilder {
	builder.feedCard.Links = append(builder.feedCard.Links, FeedLinkElement{
		Title: title,
		MessageURL: messageURL,
		PictureURL: pictureURL,
	})
	return builder
}

func (builder *feedCardBuilder) Build() FeedCardElement {
	return builder.feedCard
}

type ret struct {
	ErrorCode int `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}

type DingRobot struct {
	AccessToken string
	SendURL string
}

func NewRobot(accessToken string) *DingRobot {
	return &DingRobot{
		AccessToken:accessToken,
		SendURL:strings.Replace(BaseSendURL, "{ACCESS_TOKEN}", accessToken, 1),
	}
}

func (dr DingRobot) SendMessage(msg DingMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	resp, err := http.Post(dr.SendURL, JSONType, bytes.NewReader(body))
	if err != nil {
		return err
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	ret := new(ret)
	err = json.Unmarshal(body, ret)
	if err != nil {
		return err
	}
	if ret.ErrorCode != 0 {
		return errors.New(ret.ErrorMessage)
	}
	return nil
}

type MessageBuilder struct {
	message DingMessage
}

func NewMessageBuilder(msgType MessageType) *MessageBuilder {
	return &MessageBuilder{
		message: DingMessage{
			Type: msgType,
		},
	}
}

func (builder *MessageBuilder) Text(text string) *MessageBuilder {
	builder.message.Text = TextElement{Content: text}
	return builder
}

func (builder *MessageBuilder) Link(title string, text string, messageURL string, pictureURL string) *MessageBuilder {
	builder.message.Link = LinkElement{
		Title: title,
		Text: text,
		MessageURL: messageURL,
		PictureURL: pictureURL,
	}
	return builder
}

func (builder *MessageBuilder) Markdown(title string, text string) *MessageBuilder {
	builder.message.Markdown = MarkdownElement{
		Title: title,
		Text: text,
	}
	return builder
}

func (builder *MessageBuilder) ActionCard(element ActionCardElement) *MessageBuilder {
	builder.message.ActionCard = element
	return builder
}

func (builder *MessageBuilder) FeedCard(element FeedCardElement) *MessageBuilder {
	builder.message.FeedCard = element
	return builder
}

func (builder *MessageBuilder) At(mobiles []string, isAtAll bool) *MessageBuilder {
	builder.message.At = AtElement{
		AtMobiles: mobiles,
		IsAtAll: isAtAll,
	}
	return builder
}

func (builder *MessageBuilder) Build() DingMessage {
	return builder.message
}