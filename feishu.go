package goo

var FeiShu = new(feishu)

type feishu struct {
}

func (*feishu) Text(hookUrl string, text string) {
	msg := Params{
		"msg_type": "text",
		"content":  Params{"text": text},
	}
	PostJson(hookUrl, msg.Json())
}
