package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

type advertChecker struct {
}

func NewAdvertCheck() *advertChecker {
	return &advertChecker{}
}

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
		"  15 72 &288 [89*77【＋v】做鑫 评抬  ",
		" 15722 ]888*9 77【＋v】做❤评 抬  ",
		"1⃣8⃣1⃣2⃣7⃣8⃣1⃣2⃣8⃣9⃣3⃣，， 视9⃣，语7⃣ 文4毛，，一键招呼， 视频礼物为主，推送，互关也有米的",
		"刈览噐打开𝘺𝘬𝘬6.𝘤𝘤",
		"+lvlv19980402做❤️评抬",
		"吓栽：㊙️觅 注册时,记得把你的妮称 改用：状元。",
		"❤️萍（胎），达歌哆，来找我吧，视舞筷，语彡筷，月稳舞W +👉 peelts",
		"＋18361551070 做❤评苔",
		"商店夏栽👇🏻👇🏻 秘觅 你注册的时候用“皮特”当名字 我好找你， 娜不炎，大袖身材😍，泉咣 我那是诗蓝！",
		"回复率百分之30以上，想做的微xxm八六六九六六",
		"加好友附言：教你如何连 系你的交友对 象，不用再送礼 唔，也可以直接在平 台发好 码不会被屏 蔽",
		"吃饭的(饭)朋友的(友)",
		"虎 刷🎁，礼物是免费的 1⃣️百块钱礼物到手伍十元 。 没流量赔你伍百元＋我喂",
		"18529 110 758欣瓶()不疯浩视屏唔元雨因彡元互关不影响。",
		"互刷🎁，礼物是免费的 壹百块钱礼物到手伍十元 。 没流量赔你伍百元＋我喂 Qᵘᵏ7️⃣⑦8️⃣8️⃣",
		"在这一个月不够一万，了解下新模式，轻松日千",
		"随便看，的平台，你玩不玩",
		"趣聊天",
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

func (c advertChecker) Hit(msg string) bool {
	normalString := c.convertToNormalString(msg)
	prob := c.calcProbability(normalString)

	if prob >= 0.75 {
		return true
	} else {
		return false
	}
}

func (c advertChecker) convertSpecialLetter(code int32) int32 {
	upperLetterStart := []int32{120224, 120276, 120328, 120380, 120432}
	lowLetterStart := []int32{120250, 120302, 120354, 120406, 120458}

	for _, us := range upperLetterStart {
		if us <= code && code < us+26 {
			return code - us + 65 // 'A' = 65
		}
	}

	for _, ls := range lowLetterStart {
		if ls <= code && code < ls+26 {
			return code - ls + 97 // 'a' = 97
		}
	}

	// 各种 + 号转换, 普通加号 ASCII 码 为 43
	switch code {
	case 65291, 10133, 8314, 8315:
		return 43
	}

	// 各种心形状转换， '新' 的编码为 	26032
	if code == 10084 || (128147 <= code && code <= 128159) {
		return 26032 // '新' 的编码为 	26032
	}

	// 各种圆圈内部数字转换
	if 9312 <= code && code <= 9320 {
		return code - 9312 + 49 // 1 的编码为49
	}
	if code == 9450 {
		return 48 // 0 的编码是48
	}

	// 中文 零一二三到九 转成数字
	switch code {
	case 38646: // 零
		return 48 // 0
	case 19968:
		return 49
	case 20108:
		return 50
	case 19977:
		return 51
	case 22235:
		return 52
	case 20116:
		return 53
	case 20845:
		return 54
	case 19971:
		return 55
	case 20843:
		return 56
	case 20061: // 九
		return 57 // 9
	}

	return code
}

func (c advertChecker) convertToNormalString(s string) string {
	var result strings.Builder

	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if !unicode.IsSpace(rs[i]) && !unicode.IsPunct(rs[i]) { // 去除空格和标点符号
			if int32(rs[i]) != 8419 && int32(rs[i]) != 65039 { // 把可见的数字转成正常的数字：过滤掉数字7后边的特殊字符 7⃣
				newCode := c.convertSpecialLetter(int32(rs[i]))
				result.WriteRune(newCode) // 把可见的字母转成正常字母
			}
		}
	}

	return result.String()
}

func (c advertChecker) containPinyin(pinyinMap map[string]int8, pinyins ...string) bool {
	for _, p := range pinyins {
		if _, ok := pinyinMap[p]; !ok {
			return false
		}
	}

	return true
}

func (c advertChecker) levelOnePinyinProb(pinyinMap map[string]int8) (probability float64) {
	levelOneKeywords := [][]string{}
	levelOneKeywords = append(levelOneKeywords, []string{"zuo", "ping", "tai"})         // 做平台 zuo ping tai
	levelOneKeywords = append(levelOneKeywords, []string{"xin", "ping", "tai"})         // 新平台 xin ping tai
	levelOneKeywords = append(levelOneKeywords, []string{"hu", "shua"})                 // 互刷 hu shua
	levelOneKeywords = append(levelOneKeywords, []string{"fan", "you"})                 // 饭友 fan you
	levelOneKeywords = append(levelOneKeywords, []string{"yun", "liao"})                // 韵聊 yun liao
	levelOneKeywords = append(levelOneKeywords, []string{"lan", "qi", "da", "kai"})     // 浏 览器打开 liu lan qi da kai
	levelOneKeywords = append(levelOneKeywords, []string{"qu", "liao", "tian"})         // 趣聊天 qu liao tian
	levelOneKeywords = append(levelOneKeywords, []string{"yi", "jian", "zhao", "hu"})   // 一键招呼 yi jian zhao hu
	levelOneKeywords = append(levelOneKeywords, []string{"qing", "song", "ri", "qian"}) // 轻松日千 qing song ri qian

	for _, keywords := range levelOneKeywords {
		if c.containPinyin(pinyinMap, keywords...) {
			probability += 0.8
		}
	}

	return
}

func (c advertChecker) levelTwoPinyinProb(pinyinMap map[string]int8) (probability float64) {
	levelTwoKeywords := [][]string{}
	levelTwoKeywords = append(levelTwoKeywords, []string{"yi", "jian"})         // 一键 yi jian
	levelTwoKeywords = append(levelTwoKeywords, []string{"zhao", "hu"})         // 招呼 zhao hu
	levelTwoKeywords = append(levelTwoKeywords, []string{"shi", "pin"})         // 视频 shi pin
	levelTwoKeywords = append(levelTwoKeywords, []string{"yu", "yin"})          // 语音 yu yin
	levelTwoKeywords = append(levelTwoKeywords, []string{"li", "wu"})           // 礼物 li wu
	levelTwoKeywords = append(levelTwoKeywords, []string{"tui", "song"})        // 推送 tui song
	levelTwoKeywords = append(levelTwoKeywords, []string{"hu", "guan"})         // 互关 hu guan
	levelTwoKeywords = append(levelTwoKeywords, []string{"you", "mi"})          // 有米 you mi
	levelTwoKeywords = append(levelTwoKeywords, []string{"lan", "qi"})          // 浏 览器 liu lan qi
	levelTwoKeywords = append(levelTwoKeywords, []string{"ping", "tai"})        // 平台 ping tai
	levelTwoKeywords = append(levelTwoKeywords, []string{"xia", "zai"})         // 下载 xia zai
	levelTwoKeywords = append(levelTwoKeywords, []string{"zhu", "ce"})          // 注册 zhu ce
	levelTwoKeywords = append(levelTwoKeywords, []string{"shang", "dian"})      // 商店 shang dian
	levelTwoKeywords = append(levelTwoKeywords, []string{"bu", "yan"})          // 不严 bu yan
	levelTwoKeywords = append(levelTwoKeywords, []string{"da", "xiu"})          // 大秀 da xiu
	levelTwoKeywords = append(levelTwoKeywords, []string{"quan", "guang"})      // 全光 quan guang
	levelTwoKeywords = append(levelTwoKeywords, []string{"hui", "fu", "lv"})    // 回复率 hui fu lv
	levelTwoKeywords = append(levelTwoKeywords, []string{"xiang", "zuo"})       // 想做 xiang zuo
	levelTwoKeywords = append(levelTwoKeywords, []string{"jia", "hao", "you"})  // 加好友 jia hao you
	levelTwoKeywords = append(levelTwoKeywords, []string{"jiao", "ni"})         // 教你 jiao ni
	levelTwoKeywords = append(levelTwoKeywords, []string{"mian", "fei"})        // 免费 mian fei
	levelTwoKeywords = append(levelTwoKeywords, []string{"liu", "liang"})       // 流量 liu liang
	levelTwoKeywords = append(levelTwoKeywords, []string{"pei", "ni"})          // 赔你 pei ni
	levelTwoKeywords = append(levelTwoKeywords, []string{"wo", "wei"})          // 我微 wo wei
	levelTwoKeywords = append(levelTwoKeywords, []string{"xin", "ping"})        // 新平 台 xin ping tai
	levelTwoKeywords = append(levelTwoKeywords, []string{"feng", "hao"})        // 封号 feng hao
	levelTwoKeywords = append(levelTwoKeywords, []string{"dao", "shou"})        // 到手 dao shou
	levelTwoKeywords = append(levelTwoKeywords, []string{"xin", "mo", "shi"})   // 新模式 xin mo shi
	levelTwoKeywords = append(levelTwoKeywords, []string{"sui", "bian", "kan"}) // 随便看 sui bian kan
	levelTwoKeywords = append(levelTwoKeywords, []string{"wan", "bu"})          // 玩不 wan bu
	levelTwoKeywords = append(levelTwoKeywords, []string{"da", "ge", "duo"})    // 大哥多 da ge duo

	for _, keywords := range levelTwoKeywords {
		if c.containPinyin(pinyinMap, keywords...) {
			probability += 0.4
		}
	}

	return
}

func (c advertChecker) levelThreePinyinProb(pinyinMap map[string]int8) (probability float64) {
	levelThreeKeywords := [][]string{}
	levelThreeKeywords = append(levelThreeKeywords, []string{"shi"})        // 视 shi
	levelThreeKeywords = append(levelThreeKeywords, []string{"yu"})         // 语 yu
	levelThreeKeywords = append(levelThreeKeywords, []string{"wen"})        // 文 wen
	levelThreeKeywords = append(levelThreeKeywords, []string{"mao"})        // 毛 mao
	levelThreeKeywords = append(levelThreeKeywords, []string{"yuan"})       // 元 yuan
	levelThreeKeywords = append(levelThreeKeywords, []string{"kuai"})       // 块 kuai
	levelThreeKeywords = append(levelThreeKeywords, []string{"ping", "bi"}) // 屏蔽 ping bi

	for _, keywords := range levelThreeKeywords {
		if c.containPinyin(pinyinMap, keywords...) {
			probability += 0.1
		}
	}

	return
}

func (c advertChecker) specialCharProb(str string) (probability float64) {
	for _, r := range str {
		switch int32(r) {
		case 127873, // 礼物图标 unicode
			128071, // 手 图标
			128073, // 手 图标
			12953,  // 秘 图标
			43,     // +
			86,     // v
			118:    // V
			probability += 0.2
		}
	}

	return
}

func (c advertChecker) serialNumLetterProb(normalString string) (probability float64) {
	matched, err := regexp.MatchString("[a-zA-Z0-9]{3,5}", normalString)
	if err == nil && matched {
		probability += 0.4
	}

	matched, err = regexp.MatchString("[a-zA-Z0-9]{6,}", normalString)
	if err == nil && matched {
		probability += 0.4
	}
	return
}

func (c advertChecker) calcProbability(normalString string) (probability float64) {
	msgLen := len([]rune(normalString))
	if msgLen <= 1 {
		return 0.0
	}

	pinyinArr := pinyin.Pinyin(normalString, pinyin.NewArgs())
	// fmt.Printf("pinyin:%v\n", pinyinArr)
	pinyinMap := map[string]int8{}
	for _, pinyin := range pinyinArr {
		pinyinMap[pinyin[0]] = pinyinMap[pinyin[0]] + 1
	}

	probability += c.levelOnePinyinProb(pinyinMap)
	// fmt.Printf("levelOnePinyinProb: %.2f\n", probability)
	probability += c.levelTwoPinyinProb(pinyinMap)
	// fmt.Printf("levelTwoPinyinProb: %.2f\n", probability)
	probability += c.levelThreePinyinProb(pinyinMap)
	// fmt.Printf("levelThreePinyinProb: %.2f\n", probability)
	probability += c.specialCharProb(normalString)
	// fmt.Printf("specialCharProb: %.2f\n", probability)
	probability += c.serialNumLetterProb(normalString)
	// fmt.Printf("serialNumLetterProb: %.2f\n", probability)

	return
}
