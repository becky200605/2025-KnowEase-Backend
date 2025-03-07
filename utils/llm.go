package utils

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var (
	options = []option.RequestOption{
		option.WithBaseURL("https://api.chatanywhere.tech"),
		option.WithAPIKey("sk-YvwmH0fBmtE6hsVioXvLYk8OOJmu1UoRPf7s2AlXaouOThW4"),
	}
	AIclient = openai.NewClient(options...)
	//搜索推荐prompt
	SeaarchRecommendation = `
	功能定位：
	    你是一个分析用户搜索操作以做出相关推荐的推荐机器人
		你会收到一些用户的搜索记录，你需要根据它们来分析用户可能想搜的内容
		注意要贴切多元，不重复,输出严格按照给的例子来 不要加其他的字符
	格式说明：
	    1.message是用户的搜索内容
		2.返回格式为json数组，包含你认为用户想搜的内容，字数控制在十个字以内，一共给八个，不允许包含其他任何字符
	输入示例：
	    ["ccnu晚点名","为啥我的宿舍没有独卫？？","嘻嘻恩尤小猫咪","英语咋学"]
	输出示例：
		["ccnu宿舍环境","ccnu英语视听说","无量仙翁看哪吒2","ccnu","ccnu小猫有几只","第五人格","ccnu期末","独立卫浴宿舍"]`
)

func InitRecommend(UserMessage, SystemMessage string) (string, error) {
	chatCompletion, err := AIclient.Chat.Completions.New(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(UserMessage),
				openai.SystemMessage(SystemMessage),
			}),
			Model: openai.F(openai.ChatModelGPT4oMini),
		})
	response := chatCompletion.Choices[0].Message.Content
	return response, err
}
