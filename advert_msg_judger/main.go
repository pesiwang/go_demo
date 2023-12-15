package main

import (
	"fmt"
)

func main() {
	// var builder strings.Builder

	// for i := 120320; i <= 120400; i++ {
	// 	builder.WriteRune(rune(i))
	// }
	// fmt.Printf("build rune: [%v]\n", builder.String())

	// builder.Reset()
	// for i := 0; i <= 9; i++ {
	// 	builder.WriteRune(rune(i))
	// }
	// fmt.Printf("build rune: [%v]\n", builder.String())

	// str := "1⃣8⃣1⃣2⃣7⃣8⃣1⃣2⃣8⃣9⃣3⃣"
	// str := "Qᵃᵏ7️⃣⑦8️⃣8️⃣"
	// str := "刈览噐打开𝘺𝘬𝘬6.𝘤𝘤"
	// fmt.Printf("str: %v\n", []rune(str))
	// builder.Reset()
	// for _, s := range str {
	// 	if 48 <= s && s <= 57 {
	// 		builder.WriteRune(s)
	// 	}
	// }
	// fmt.Printf("only num str:%v\n", builder.String())
	msgs := []string{
		// "  15 72 &288 [89*77【＋v】做鑫 评抬  ",
		// " 15722 ]888*9 77【＋v】做❤评 抬  ",
		// "1⃣8⃣1⃣2⃣7⃣8⃣1⃣2⃣8⃣9⃣3⃣，， 视9⃣，语7⃣ 文4毛，，一键招呼， 视频礼物为主，推送，互关也有米的",
		// "刈览噐打开𝘺𝘬𝘬6.𝘤𝘤",
		// "+lvlv19980402做❤️评抬",
		// "吓栽：㊙️觅 注册时,记得把你的妮称 改用：状元。",
		// "❤️萍（胎），达歌哆，来找我吧，视舞筷，语彡筷，月稳舞W +👉 peelts",
		// "＋18361551070 做❤评苔",
		// "商店夏栽👇🏻👇🏻 秘觅 你注册的时候用“皮特”当名字 我好找你， 娜不炎，大袖身材😍，泉咣 我那是诗蓝！",
		// "回复率百分之30以上，想做的微xxm八六六九六六",
		// "加好友附言：教你如何连 系你的交友对 象，不用再送礼 唔，也可以直接在平 台发好 码不会被屏 蔽",
		// "吃饭的(饭)朋友的(友)",
		// "虎 刷🎁，礼物是免费的 1⃣️百块钱礼物到手伍十元 。 没流量赔你伍百元＋我喂",
		// "18529 110 758欣瓶()不疯浩视屏唔元雨因彡元互关不影响。",
		// "互刷🎁，礼物是免费的 壹百块钱礼物到手伍十元 。 没流量赔你伍百元＋我喂 Qᵘᵏ7️⃣⑦8️⃣8️⃣",
		// "在这一个月不够一万，了解下新模式，轻松日千",
		// "随便看，的平台，你玩不玩",
		// "趣聊天",
		"我选对我好的呀。姐姐，你起床了吗？你会喜欢会做饭的男生吗？我也尽心来找个女朋友，我想过年带回家。",
	}

	// for _, m := range msgs {
	// 	rawMsg := removeSpacesAndPunctuation(m)
	// 	fmt.Println(rawMsg)
	// 	for _, c := range rawMsg {
	// 		fmt.Printf("%v,", c)
	// 	}
	// 	fmt.Println()
	// }

	advertChecker := NewAdvertCheck()
	for _, m := range msgs {
		if advertChecker.Hit(m) {
			fmt.Printf("yes: %v\n", m)
		} else {
			fmt.Printf("no: %v\n", m)
		}
	}
}
