package pkg

import (
	"context"
	"fmt"
	util "github.com/hktalent/go-utils"
	ratelimit "github.com/projectdiscovery/ratelimit"
	"github.com/sashabaranov/go-openai"
	"time"

	"net/http"
	"net/url"
	"strings"
)

var (
	GptApi  *openai.Client
	Prefix  string
	Limiter *ratelimit.Limiter
)

func init() {
	util.RegInitFunc(func() {
		Limiter = ratelimit.New(util.Ctx_global, uint(util.GetValAsInt("LimitPerMinute", 20)), time.Minute)
		Prefix = util.GetVal("Prefix")
		szProxy := util.GetVal("proxy")
		chatGptKey := util.GetVal("api_key")
		if szProxy == "" {
			GptApi = openai.NewClient(chatGptKey)
		} else {
			config := openai.DefaultConfig(chatGptKey)
			proxyUrl, err := url.Parse(szProxy)
			if err != nil {
				panic(err)
			}
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
			config.HTTPClient = &http.Client{
				Transport: transport,
			}
			GptApi = openai.NewClientWithConfig(config)
		}
	})
}

/*
https://github.com/hktalent/PostExploitation/blob/b945899b24bf719df3c502fa315929c6bddabc57/test/src/ysoserial/payloads/util/ExpJndi.java#L4
 1. Model：除了 openai.GPT3Ada（GPT-3 Ada 模型）之外，可以选择其他的 GPT-3 模型，包括以下选项：
    • openai.GPT3: GPT-3 模型
    • openai.GPT3_Babbage: GPT-3 Babbage 模型
    • openai.GPT3_Curie: GPT-3 Curie 模型
    • openai.GPT3_Davinci: GPT-3 Davinci 模型

openai.GPT3Ada 是 OpenAI 最新发布的 GPT-3 模型，它与 GPT-3 模型在结构上基本相同，但在训练算法上有所不同。openai.GPT3Ada 集成了一种名为 AdaBelief Optimizer 的训练算法，它可有效应对 GPT-3 模型中存在的过拟合和样本不平衡问题。

具体来说，相比于 GPT-3，openai.GPT3Ada 的训练算法更加稳定和可靠，能够在处理大规模数据时表现出更好的训练效果。在实际应用中，openai.GPT3Ada 可以生成更加流畅、连贯、富于创造力的文本内容，并且具有更高的准确性和可靠性。

但需要注意的是，openai.GPT3Ada 目前仅在 OpenAI 的研究论文中发布，尚未完全开放给开发者使用。预计在未来，随着相关技术的不断进步和完善，openai.GPT3Ada 将有望成为自然语言处理领域的一个重要创新。

	同时，也可以选择使用 OpenAI 的其他文本生成模型，这需要注册并获取访问 API 的密钥。这些模型包括 GPT-2、DialoGPT、Codex、等等。
	2. MaxTokens：除了指定生成的最大令牌数之外，还可以选择以下选项：
		• 0：不限制生成令牌的数量
		• 1-2048：限制生成令牌的数量，范围从1到2048之间。
	3. Prompt：可以根据具体需求添加前缀或者上下文信息，以帮助模型更好地生成文本。同时，还可以选择以下选项：
		• 将提示信息设置为无，即不提供上下文信息
		• 使用多个语句组合作为提示信息
	4. Temperature：控制生成的文本的创造性和多样性。可以选择以下选项：
		• 0-1：限制生成的文本更倾向于重复已有的内容，更接近于人类的写作风格。
		• 1：文本生成的中等水平，尝试平衡创造性和准确性。
		• 1+：鼓励生成更创新、更多样化的文本。
	5. FrequencyPenalty：控制生成文本的惩罚程度。如果生成的文本包含重复的模式、词语等，则可以通过惩罚来强制模型更广泛地探索语言空间。可以选择以下选项：
		• 0-1：惩罚力度较弱，可以生成一些重复的文本。
		• 1+：惩罚力度较强，可以生成更多非常规和多样化的文本。
	6. PresencePenalty：控制生成文本的惩罚程度，以确保模型在生成文本时遵循给定上下文背景。可以选择以下选项：
		• 0-1：惩罚力度较弱，文本可能与上下文不太相关。
		• 1+：惩罚力度较强，文本应该能够更好地与上下文配合。

	还有一些其他的参数，如 Stop、N、Stream、LogProbs、Echo 等等，这些参数可以根据具体需求进行调整。
*/
func GptNew(s string) (string, error) {
	Limiter.Take()
	ctx := context.Background()
	resp, err := GptApi.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: s,
				},
			},
		},
	)

	/*
			error, status code: 429,
			message: Rate limit reached for default-gpt-3.5-turbo in organization org-xDa56WiVjCMJiVef9SrQOFPW on requests per min.
			Limit: 20 / min. Please try again in 3s.
			Contact support@openai.com if you continue to have issues.
			Please add a payment method to your account to increase your rate limit.
			Visit https://platform.openai.com/account/billing to add a payment method.

		Completion error: error, status code: 400,
		message: This model's maximum context length is 4097 tokens.
		However, your messages resulted in 4424 tokens.
		Please reduce the length of the messages.
	*/
	if err != nil {
		fmt.Printf("%s\nCompletion error: %v\n", s, err)
		return "", err
	}
	fmt.Println(len(resp.Choices), resp.Choices[0].Message.Content)
	return strings.TrimSpace(resp.Choices[0].Message.Content), err
}
