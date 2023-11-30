package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cast"
)

const API_KEY = "GKrRZrQdb203y1jvCpG53GNO"
const SECRET_KEY = "a4lT8Xp1yfZKIHPVKW0ODazXSVsaXeP8"

type llmMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const prefixContent string = `上下文：陌生人付费聊天社区上的收益方通过回复消费方的消息或收取礼物获得收益。为了获得更高的收益，部分用户会发送一些不合规的消息用于绕过平台限制，让消费方发送更多消息和礼物。为了更好地维护社区，需要识别出每条消息的消息类型
角色：你是一个非常有经验的细心的聊天社区监管员，可以准确地判断出每条消息的类型
能力：如果用户发送的消息中带有用户的出生年份或年龄，你会判定此消息为年龄消息；如果消息中带有用户的位置信息，你会判定消息为位置消息；你还会准确地区分普通消息，绕过平台限制消息，索要礼物消息等类型，并最终返回消息类型id
消息类型由 id: 名称 - 说明组成，具体如下：
{
0: 普通 - 普通的聊天内容，如果不能确定归为下边的其他类型，归为普通消息类型
1: 位置 - 消息中包含用户的地址，位置信息等，归为位置消息类型
2: 年龄 - 消息中包含用户的出生年份或年龄信息，比如 "我97"，"我03年"，"我26岁"等等，归为年龄消息类型
3: 绕过平台限制 - 向对方提示看不到被平台规则限制的消息，引导对方通过拆分消息发送更多消息，归为绕过平台限制消息类型
4: 索要礼物 - 直接要求对方发送礼物，向对方索要礼物，归为索要礼物消息类型
}
我会提供某个用户发送的多条消息内容，你从这个用户发送的消息中判断消息类型并返回消息类型的id，输入输出为json格式
注意：每条消息输入必须对应一个id输出，id只能是 0，1，2，3，4之一，如果你不确定请输出0
input，analyse, output举例：
input:
[
"小哥哥，在忙什么呢？方便的话聊会儿啊",
"在干嘛呢",
"我96的",
"我27",
"我26",
"我95的",
"我18",
"来个黑桃香槟",
"来个浪漫飞机"
]
analyse:
{
"小哥哥，在忙什么呢？方便的话聊会儿啊":"普通",
"在干嘛呢":"普通",
"我96的":"年龄",
"我27":"年龄",
"我26":"年龄",
"我95的":"年龄",
"我18":"年龄",
"来个黑桃香槟":"索要礼物",
"来个浪漫飞机":"索要礼物"
}
output:
[
0,
0,
2,
2,
2,
2,
2,
4,
4
]
input:
[
"hello，交个朋友",
"你老家哪里的？值得去玩一趟吗",
"我在广东",
"我在广东省深圳市",
"我在广东深圳",
"我在深圳",
"我在广州",
"我在上海"
]
analyse:
{
"hello，交个朋友":"普通",
"你老家哪里的？值得去玩一趟吗":"普通",
"我在广东":"位置",
"我在广东省深圳市":"位置",
"我在广东深圳":"位置",
"我在深圳":"位置",
"我在广州":"位置",
"我在上海":"位置"
}
output:
[
0,
0,
1,
1,
1,
1,
1,
1
]
input:
[
"吃饭了吗？",
"哈喽，互相了解一下吗？",
"你喜欢什么类型的女生呢？",
"一个一个发",
"我是95的",
"可以送我礼物吗",
"我是98的"
]
analyse:
{
"吃法了吗": "普通",
"哈喽，互相了解一下吗": "普通",
"你喜欢什么类型的女生呢？": "普通",
"一个一个发": "绕过平台限制",
"我是95的": "年龄",
"可以送我礼物吗": "索要礼物",
"我是98的": "年龄"
}
output:
[
0,
0,
0,
3,
2,
4,
2
]
input:
[
"我是95的",
"我在梅州",
"看不到你发",
"一个一个发",
"发了吗",
"我在广州",
"我在深圳",
"我在潮汕",
"我是97的",
"我想要黑桃",
"我要戒指",
"刷戒指",
"送个戒指",
"来个爱心",
"送我蝶舞戒指",
"我今年18"
]
analyse:
{
"我是95的":"年龄",
"我在梅州":"位置",
"看不到你发":"绕过平台限制",
"一个一个发":"绕过平台限制",
"发了吗":"绕过平台限制",
"我在广州":"位置",
"我在深圳":"位置",
"我在潮汕":"位置",
"我是97的":"年龄",
"我想要黑桃":"索要礼物",
"我要戒指":"索要礼物",
"刷戒指":"索要礼物",
"送个戒指":"索要礼物",
"来个爱心":"索要礼物",
"送我蝶舞戒指":"索要礼物",
"我今年18":"年龄"
}
output:
[
2,
1,
3,
3,
3,
1,
1,
1,
2,
4,
4,
4,
4,
4,
4,
2
]
input, analyse, output 举例结束，现在输入json数组 input, 你要准确地输出与input数组长度相等的json数组 output
注意：消息输入条数必须等于输出id数量，id只能是 0，1，2，3，4之一，如果你不确定请输出0
`

const llmInput string = `
input:
[
"我26的",
"我97的",
"吃饭了吗？",
"来个娃娃机？",
"我在深圳宝安区",
"一个一个发",
"送我一个礼物吧",
"在干嘛呢？",
"我在宝安区",
"你喜欢什么类型的女生呢？",
"哈喽，互相了解一下吗？",
"送我一个爱心",
"我在广州",
"我01年的",
"看不到",
"发的什么",
"我今年26岁",
"发了吗？",
"我今年18",
"亲亲在干嘛呢",
"上班还可以玩手机",
"我上夜班",
"你需要人陪吗",
"在成都工作",
"我是潼南的",
"刷玫瑰花束",
"来个戒指",
"你是哪里人呀？",
"你有什么兴趣爱好呀？",
"我现在在佛山",
"我是87的，你呢？",
"今年36岁了，不小了",
"可以给送一个钻石吗？",
"送我旋转木马",
"我02的",
"嗨，很高兴认识你~"
]
output:
`

func main() {

	accessToken := GetAccessToken()
	fmt.Printf("access token:%v\n", accessToken)
	if len(accessToken) <= 0 {
		fmt.Printf("fetch baiduyun access token failed")
		return
	}

	url := "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat/completions_pro?access_token=" + accessToken

	reqBody := map[string]interface{}{}
	msgs := []llmMsg{}
	msgs = append(msgs, llmMsg{
		Role:    "user",
		Content: prefixContent + llmInput, // "你是谁？",
	})

	reqBody["temperature"] = float64(0.01)
	//reqBody["top_p"] = float64(0.01)
	// fmt.Printf("msg:%v\n\n", msgs[0])
	reqBody["messages"] = msgs

	reqBodyJSON, _ := json.Marshal(reqBody)
	// fmt.Printf("reqBodyJSON:%v\n\n", string(reqBodyJSON))
	payload := bytes.NewReader(reqBodyJSON)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))

	respObj := map[string]interface{}{}
	err = json.Unmarshal(body, &respObj)
	if err != nil || cast.ToString(respObj["finish_reason"]) != "normal" {
		fmt.Printf("ernie-bot-4 return failed: err %v, body:%v\n", err, string(body))
		return
	}

	resultStr := cast.ToString(respObj["result"])
	fmt.Println(resultStr)

	resultArr := []int64{}
	err = json.Unmarshal([]byte(resultStr), &resultArr)
	if err != nil || len(resultArr) <= 0 {
		fmt.Printf("ernie-bot-4 parse result failed: err %v, body:%v\n", err, string(body))
		return
	}

	fmt.Printf("result:%v", resultArr)
}

/**
 * 使用 AK，SK 生成鉴权签名（Access Token）
 * @return string 鉴权签名信息（Access Token）
 */
func GetAccessToken() string {
	url := "https://aip.baidubce.com/oauth/2.0/token"
	postData := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", API_KEY, SECRET_KEY)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(postData))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// fmt.Printf("access token resp: %v\n", string(body))
	accessTokenObj := map[string]interface{}{}
	err = json.Unmarshal(body, &accessTokenObj)
	if err != nil || len(cast.ToString(accessTokenObj["error"])) > 0 {
		fmt.Printf("get access token failed: err %v %v:%v\n", err, accessTokenObj["error"], accessTokenObj["error_description"])
		return ""
	}

	return cast.ToString(accessTokenObj["access_token"])
}
