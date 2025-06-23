/*
 *  Copyright (c) 2025, WSO2 LLC. (http://www.wso2.org) All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package inbuiltpolicy

import "regexp"

// Constants for Guardrail policies
const (
	ErrorCode    = "code"
	ErrorType    = "type"
	ErrorMessage = "message"

	GuardrailErrorCode         = 446
	GuardrailAPIMExceptionCode = 900514
	GuardrailErrorType         = "ERROR_TYPE"
	AssessmentAction           = "action"
	InterveningGuardrail       = "interveningGuardrail"
	APIMInternalExceptionCode  = 900967
	AssessmentReason           = "actionReason"
	Direction                  = "direction"
	Assessments                = "assessments"

	RegexGuardrailName         = "RegexGuardrail"
	WordCountGuardrailName     = "WordCountGuardrail"
	SentenceCountGuardrailName = "SentenceCountGuardrail"
	ContentLengthGuardrailName = "ContentLengthGuardrail"

	RegexGuardrailConstant         = "REGEX_GUARDRAIL"
	WordCountGuardrailConstant     = "WORD_COUNT_GUARDRAIL"
	SentenceCountGuardrailConstant = "SENTENCE_COUNT_GUARDRAIL"
	ContentLengthGuardrailConstant = "CONTENT_LENGTH_GUARDRAIL"

	TextCleanRegex     = "^\"|\"$"
	WordSplitRegex     = "\\s+"
	SentenceSplitRegex = "[.!?]"
)

var (
	TextCleanRegexCompiled     = regexp.MustCompile(TextCleanRegex)     // TextCleanRegexCompiled is used to clean text by removing leading and trailing quotes
	WordSplitRegexCompiled     = regexp.MustCompile(WordSplitRegex)     // WordSplitRegexCompiled is used to split text into words based on whitespace
	SentenceSplitRegexCompiled = regexp.MustCompile(SentenceSplitRegex) // SentenceSplitRegexCompiled is used to split text into sentences based on punctuation marks (., !, ?
)
