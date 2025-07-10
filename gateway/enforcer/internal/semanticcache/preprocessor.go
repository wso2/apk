package semanticcache

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
)

// - This file contains the preprocessing functions. 
var stopwordSet map[string]struct{}

func textToLowercase(text string) string {
	return strings.ToLower(text)
}

func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func wordTokenizer(text string) []string {
	// Use regex to split on whitespace and punctuation
	re := regexp.MustCompile(`\w+`)
	tokens := re.FindAllString(text, -1)
	var result []string
	for _, token := range tokens {
		if len(strings.TrimSpace(token)) > 0 {
			result = append(result, token)
		}
	}
	return result
}

func sentTokenize(text string) []string {
	var sentenceEndRegex = regexp.MustCompile(`(?m)([^.!?]+[.!?])`)
	matches := sentenceEndRegex.FindAllString(text, -1)
	for i := range matches {
		matches[i] = strings.TrimSpace(matches[i])
	}
	return matches
}

// ConvertPromptToChunks converts a prompt into chunks of sentences based on the specified chunk size.
func ConvertPromptToChunks(prompt string, chunkSize int) []string{
	sentencesList := sentTokenize(prompt)
	if chunkSize <= 0 {
		fmt.Printf("Unable to perform chunking. Invalid Chunksize provided: %d", chunkSize)
		return nil
	}

	var chunks []string
	for i := 0; i < len(sentencesList); i += chunkSize {
		end := i + chunkSize
		if end > len(sentencesList) {
			end = len(sentencesList)
		}

		var builder strings.Builder
		for _, sentence := range sentencesList[i:end] {
			clean := normalizeWhitespace(sentence)
			if clean != "" {
				if builder.Len() > 0 {
					builder.WriteString(" ")
				}
				builder.WriteString(clean)
			}
		}
		chunks = append(chunks, builder.String())
	}
	return chunks
}

func removePunctuation(token string) string {
	var result strings.Builder
	for _, char := range token {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// InitStopwordList initialized the stop-words map which will be by removeStopWords method
func InitStopwordList() {
	words := []string{
		"a","about","above","after","again","against","ain","all","am","an","and","any",
		"are","aren","aren't","as","at","be","because","been","before","being","below",
		"between","both","but","by","can","couldn","couldn't","d","did","didn","didn't",
		"do","does","doesn","doesn't","doing","don","don't","down","during","each","few",
		"for","from","further","had","hadn","hadn't","has","hasn","hasn't","have","haven",
		"haven't","having","he","her","here","hers","herself","him","himself","his","how",
		"i","if","in","into","is","isn","isn't","it","it's","its","itself","just","ll","m",
		"ma","me","mightn","mightn't","more","most","mustn","mustn't","my","myself","needn",
		"needn't","no","nor","not","now","o","of","off","on","once","only","or","other","our",
		"ours","ourselves","out","over","own","re","s","same","shan","shan't","she","she's",
		"should","should've","shouldn","shouldn't","so","some","such","t","than","that","that'll",
		"the","their","theirs","them","themselves","then","there","these","they","this","those",
		"through","to","too","under","until","up","ve","very","was","wasn","wasn't","we","were",
		"weren","weren't","what","when","where","which","while","who","whom","why","will","with",
		"won","won't","wouldn","wouldn't","y","you","you'd","you'll","you're","you've","your","yours",
		"yourself","yourselves","could","he'd","he'll","he's","here's","how's","i'd","i'll","i'm","i've",
		"let's","ought","she'd","she'll","that's","there's","they'd","they'll","they're","they've","we'd",
		"we'll","we're","we've","what's","when's","where's","who's","why's","would","able","abst","accordance",
		"according","accordingly","across","act","actually","added","adj","affected","affecting","affects",
		"afterwards","ah","almost","alone","along","already","also","although","always","among","amongst",
		"announce","another","anybody","anyhow","anymore","anyone","anything","anyway","anyways","anywhere",
		"apparently","approximately","arent","arise","around","aside","ask","asking","auth","available","away",
		"awfully","b","back","became","become","becomes","becoming","beforehand","begin","beginning","beginnings",
		"begins","behind","believe","beside","besides","beyond","biol","brief","briefly","c","ca","came","cannot",
		"can't","cause","causes","certain","certainly","co","com","come","comes","contain","containing","contains",
		"couldnt","date","different","done","downwards","due","e","ed","edu","effect","eg","eight","eighty","either",
		"else","elsewhere","end","ending","enough","especially","et","etc","even","ever","every","everybody","everyone",
		"everything","everywhere","ex","except","f","far","ff","fifth","first","five","fix","followed","following","follows",
		"former","formerly","forth","found","four","furthermore","g","gave","get","gets","getting","give","given","gives",
		"giving","go","goes","gone","got","gotten","h","happens","hardly","hed","hence","hereafter","hereby","herein",
		"heres","hereupon","hes","hi","hid","hither","home","howbeit","however","hundred","id","ie","im","immediate",
		"immediately","importance","important","inc","indeed","index","information","instead","invention","inward","itd",
		"it'll","j","k","keep","keeps","kept","kg","km","know","known","knows","l","largely","last","lately","later",
		"latter","latterly","least","less","lest","let","lets","like","liked","likely","line","little","'ll","look",
		"looking","looks","ltd","made","mainly","make","makes","many","may","maybe","mean","means","meantime","meanwhile",
		"merely","mg","might","million","miss","ml","moreover","mostly","mr","mrs","much","mug","must","n","na","name",
		"namely","nay","nd","near","nearly","necessarily","necessary","need","needs","neither","never","nevertheless","new",
		"next","nine","ninety","nobody","non","none","nonetheless","noone","normally","nos","noted","nothing","nowhere","obtain",
		"obtained","obviously","often","oh","ok","okay","old","omitted","one","ones","onto","ord","others","otherwise","outside",
		"overall","owing","p","page","pages","part","particular","particularly","past","per","perhaps","placed","please","plus",
		"poorly","possible","possibly","potentially","pp","predominantly","present","previously","primarily","probably","promptly",
		"proud","provides","put","q","que","quickly","quite","qv","r","ran","rather","rd","readily","really","recent","recently",
		"ref","refs","regarding","regardless","regards","related","relatively","research","respectively","resulted","resulting",
		"results","right","run","said","saw","say","saying","says","sec","section","see","seeing","seem","seemed","seeming","seems",
		"seen","self","selves","sent","seven","several","shall","shed","shes","show","showed","shown","showns","shows","significant",
		"significantly","similar","similarly","since","six","slightly","somebody","somehow","someone","somethan","something","sometime",
		"sometimes","somewhat","somewhere","soon","sorry","specifically","specified","specify","specifying","still","stop","strongly","sub",
		"substantially","successfully","sufficiently","suggest","sup","sure","take","taken","taking","tell","tends","th","thank","thanks",
		"thanx","thats","that've","thence","thereafter","thereby","thered","therefore","therein","there'll","thereof","therere","theres",
		"thereto","thereupon","there've","theyd","theyre","think","thou","though","thoughh","thousand","throug","throughout","thru","thus",
		"til","tip","together","took","toward","towards","tried","tries","truly","try","trying","ts","twice","two","u","un","unfortunately",
		"unless","unlike","unlikely","unto","upon","ups","us","use","used","useful","usefully","usefulness","uses","using","usually","v","value",
		"various","'ve","via","viz","vol","vols","vs","w","want","wants","wasnt","way","wed","welcome","went","werent","whatever","what'll","whats",
		"whence","whenever","whereafter","whereas","whereby","wherein","wheres","whereupon","wherever","whether","whim","whither","whod","whoever",
		"whole","who'll","whomever","whos","whose","widely","willing","wish","within","without","wont","words","world","wouldnt","www","x","yes","yet",
		"youd","youre","z","zero","a's","ain't","allow","allows","apart","appear","appreciate","appropriate","associated","best","better","c'mon","c's",
		"cant","changes","clearly","concerning","consequently","consider","considering","corresponding","course","currently","definitely","described",
		"despite","entirely","exactly","example","going","greetings","hello","help","hopefully","ignored","inasmuch","indicate","indicated","indicates",
		"inner","insofar","it'd","keep","keeps","novel","presumably","reasonably","second","secondly","sensible","serious","seriously","sure","t's","third",
		"thorough","thoroughly","three","well","wonder",
	}

	stopwordSet = make(map[string]struct{}, len(words))
	for _, word := range words {
		stopwordSet[word] = struct{}{}
	}
}

func removeStopWords(tokens []string) []string {
	var filtered []string
	filtered = make([]string, 0, len(tokens))

	for _, token := range tokens {
		word := strings.ToLower(token)
		if _, isStopWord := stopwordSet[word]; !isStopWord {
			filtered = append(filtered, token)
		}
	}
	return filtered
}

func stemPrompt(tokens []string) []string {
	var stemmed []string
	stemmed = make([]string, 0, len(tokens))

	for _, token := range tokens {
		word, err := snowball.Stem(token, "english", true)
		if err != nil {
			fmt.Printf("Unable to stem the token %s", token)
		}
		stemmed = append(stemmed, word)
	}
	return stemmed
}


func combineTokens(tokens []string) string {
	return strings.Join(tokens, " ")
}

func estimateTokens(text string) int {
	words := strings.Fields(text)
	return int(float64(len(words)) * 1.25)
}

// ChunkSentencesByTokenLimit takes a list of sentences and chunks them into groups that do not exceed the maxTokens limit.
func ChunkSentencesByTokenLimit(sentences []string, maxTokens int) []string {
	var chunks []string
	var currentChunk []string
	currentTokens := 0

	for _, sentence := range sentences {
		normalized := normalizeWhitespace(sentence)
		tokens := estimateTokens(normalized)

		if currentTokens+tokens > maxTokens {
			if len(currentChunk) > 0 {
				chunks = append(chunks, strings.Join(currentChunk, " "))
				currentChunk = []string{}
				currentTokens = 0
			}
		}
		currentChunk = append(currentChunk, normalized)
		currentTokens += tokens
	}
	if len(currentChunk) > 0 {
		chunks = append(chunks, strings.Join(currentChunk, " "))
	}
	return chunks
}


// PreprocessText function will convert the text to lowercase, remove punctuations, tokenize to words,
// remove the stop words, stem the tokens and form the final cleaned prompt with reduced token count
func PreprocessText(prompt string) string {
	lowercasedPrompt := textToLowercase(prompt)
	tokenizedPrompt := wordTokenizer(lowercasedPrompt)
	promptWithoutStopWords := removeStopWords(tokenizedPrompt)
	stemmedPrompt := stemPrompt(promptWithoutStopWords)
	combinedCleanedPrompt := combineTokens(stemmedPrompt)
	return combinedCleanedPrompt
}

// ValidateVectorStoreConfigProps validates the properties of the vector store configuration.
func ValidateVectorStoreConfigProps(config VectorDBProviderConfig) error {
	if config.VectorStoreProvider != "REDIS"  && config.VectorStoreProvider != "MILVUS" {
		return fmt.Errorf("invalid vector store provider found in the vector store configuration")
	}
	if config.EmbeddingDimention == "" {
		return fmt.Errorf("missing embedding dimension in the vector store configuration")
	}
	if config.Threshold == "" {
		return fmt.Errorf("missing threshold in the vector store configuration")
	}
	if config.DBHost == "" {
		return fmt.Errorf("missing database host in the vector store configuration")
	}
	if config.DBPort == 0 || config.DBPort < 0{
		return fmt.Errorf("missing/invalid database port in the vector store configuration")
	}
	if config.Username == "" {
		return fmt.Errorf("missing DB username in the vector store configuration")
	}
	if config.Password == "" {
		return fmt.Errorf("missing DB password in the vector store configuration")
	}
	if config.DatabaseName == "" {
		return fmt.Errorf("missing database name in the vector store configuration")
	}
	return nil
}

// ValidateEmbeddingProviderConfigProps validates the properties of the embedding provider configuration.
func ValidateEmbeddingProviderConfigProps(config EmbeddingProviderConfig) error {
	if config.AuthHeaderName == "" {
		return fmt.Errorf("missing auth header name in the embedding provider configuration")
	}
	if config.APIKey == "" {
		return fmt.Errorf("missing API key in the embedding provider configuration")
	}
	if config.EmbeddingEndpoint == "" {
		return fmt.Errorf("missing embedding endpoint in the embedding provider configuration")
	}
	if config.EmbeddingProvider != "MISTRAL_AI" && config.EmbeddingProvider != "AZURE_OPENAI" {
		return fmt.Errorf("missing/Invalid embedding provider found in the embedding provider configuration")
	}
	if config.EmbeddingModel == "" {
		return fmt.Errorf("missing embedding model in the embedding provider configuration")
	}
	return nil
}