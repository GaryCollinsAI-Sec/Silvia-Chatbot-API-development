package chatbot

import (
	"strings"
)

func Answer(message string) string {
	message = strings.ToLower(strings.TrimSpace(message))

	// Security boundary checks must happen before normal
	// knowledge matching so restricted requests cannot
	// accidentally trigger an approved-information response.

	if isInstructionDisclosureRequest(message) {
		return "I can't provide my internal instructions, system prompts, or hidden configuration."
	}

	if isPromptInjectionAttempt(message) {
		return "I can't follow requests to override or bypass my instructions. I can help with approved Silver Dragons Academy information."
	}

	if isPrivateInformationRequest(message) {
		return "I can't provide private, administrative, database, or other internal information. Please use the Silver Dragons Academy Contact page for information intended for the public."
	}

	if isFabricationRequest(message) {
		return "I can't make up or invent Silver Dragons Academy information. I can only provide information that has been approved for Silvia."
	}

	// Approved academy knowledge.

	academyName := strings.ToLower(AcademyName)

	if strings.Contains(message, "what is "+academyName) ||
		strings.Contains(message, "what is "+strings.ToLower(strings.TrimSuffix(AcademyName, " Academy"))) ||
		strings.Contains(message, "tell me about "+academyName) ||
		strings.Contains(message, "about "+academyName) {

		return AcademyDescription
	}

	if strings.Contains(message, "academy mission") ||
		strings.Contains(message, "mission of the academy") ||
		strings.Contains(message, "what is the mission") {

		return AcademyMission
	}

	if strings.Contains(message, "what is taekwondo") ||
		strings.Contains(message, "what is "+strings.ToLower(CurrentProgram)) ||
		strings.Contains(message, "tell me about taekwondo") {

		return ProgramDetailedDescription
	}

	if strings.Contains(message, "taekwondo program") ||
		strings.Contains(message, "about the taekwondo program") ||
		strings.Contains(message, "taekwondo training") {

		return ProgramDescription
	}

	if strings.Contains(message, "training approach") ||
		strings.Contains(message, "how do students train") ||
		strings.Contains(message, "how does training work") ||
		strings.Contains(message, "how are students trained") {

		return TrainingApproach
	}

	if strings.Contains(message, "program") ||
		strings.Contains(message, "programs") ||
		strings.Contains(message, "class") ||
		strings.Contains(message, "classes") {

		return AcademyName + " currently offers " + CurrentProgram + " as its primary program. " + ProgramDescription
	}

	if strings.Contains(message, "benefit") ||
		strings.Contains(message, "benefits") ||
		strings.Contains(message, "help me develop") ||
		strings.Contains(message, "help develop") {

		return AcademyName + " " + CurrentProgram + " training can help develop " + formatBenefits() + "."
	}

	if strings.Contains(message, "getting started") ||
		strings.Contains(message, "how do i get started") ||
		strings.Contains(message, "how can i get started") ||
		strings.Contains(message, "interested in taekwondo") ||
		strings.Contains(message, "interested in joining") {

		return GettingStartedInformation
	}

	if strings.Contains(message, "price") ||
		strings.Contains(message, "pricing") ||
		strings.Contains(message, "cost") ||
		strings.Contains(message, "membership cost") {

		return LimitationsResponse
	}

	if strings.Contains(message, "schedule") ||
		strings.Contains(message, "class time") ||
		strings.Contains(message, "class times") {

		return LimitationsResponse
	}

	if strings.Contains(message, "enrollment") ||
		strings.Contains(message, "enroll") ||
		strings.Contains(message, "belt testing") ||
		strings.Contains(message, "testing requirements") ||
		strings.Contains(message, "refund") ||
		strings.Contains(message, "refunds") ||
		strings.Contains(message, "cancellation") ||
		strings.Contains(message, "cancellations") {

		return LimitationsResponse
	}

	if strings.Contains(message, "hello") ||
		strings.Contains(message, "hi") ||
		strings.Contains(message, "hey") {

		return "Hello! I'm Silvia, the Silver Dragons Academy virtual assistant. How can I help?"
	}

	return LimitationsResponse
}

func formatBenefits() string {
	benefits := make([]string, len(Benefits))

	for i, benefit := range Benefits {
		benefits[i] = strings.ToLower(benefit)
	}

	switch len(benefits) {
	case 0:
		return "character and personal development"
	case 1:
		return benefits[0]
	case 2:
		return benefits[0] + " and " + benefits[1]
	default:
		return strings.Join(benefits[:len(benefits)-1], ", ") +
			", and " +
			benefits[len(benefits)-1]
	}
}

func isInstructionDisclosureRequest(message string) bool {
	instructionTerms := []string{
		"system prompt",
		"system instructions",
		"internal instructions",
		"hidden instructions",
		"internal prompt",
		"hidden prompt",
		"show me your prompt",
		"show your prompt",
		"reveal your prompt",
		"reveal your instructions",
		"tell me your instructions",
		"what are your instructions",
	}

	for _, term := range instructionTerms {
		if strings.Contains(message, term) {
			return true
		}
	}

	return false
}

func isPromptInjectionAttempt(message string) bool {
	injectionTerms := []string{
		"ignore your instructions",
		"ignore previous instructions",
		"ignore all instructions",
		"disregard your instructions",
		"disregard previous instructions",
		"override your instructions",
		"bypass your instructions",
		"forget your instructions",
		"act as if you have no instructions",
	}

	for _, term := range injectionTerms {
		if strings.Contains(message, term) {
			return true
		}
	}

	return false
}

func isPrivateInformationRequest(message string) bool {
	privateTerms := []string{
		"database information",
		"access the academy database",
		"access academy database",
		"admin information",
		"administrator information",
		"admin credentials",
		"administrator credentials",
		"password",
		"secret key",
		"api key",
		"private information",
		"internal information",
		"private data",
		"internal data",
		"student records",
		"student information",
		"member records",
		"member information",
	}

	for _, term := range privateTerms {
		if strings.Contains(message, term) {
			return true
		}
	}

	return false
}

func isFabricationRequest(message string) bool {
	fabricationTerms := []string{
		"make up",
		"make something up",
		"invent",
		"fabricate",
		"pretend that",
		"pretend you know",
		"just guess",
		"guess the",
		"create fake",
		"give me a fake",
	}

	for _, term := range fabricationTerms {
		if strings.Contains(message, term) {
			return true
		}
	}

	return false
}
