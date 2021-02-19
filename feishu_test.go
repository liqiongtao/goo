package goo

import "testing"

func TestFeishu_Text(t *testing.T) {
	hookUrl := ""
	FeiShu.Text(hookUrl, "测试")
}
