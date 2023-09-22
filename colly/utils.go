package colly

import "math/rand"

const TimeFormat = "2006-01-02 15:04:05"

var sourceDomains = []string{
	"https://tkzy1.com",
	"https://tkzy2.com",
	"https://tkzy3.com",
	"https://tkzy4.com",
	"https://tkzy5.com",
	"https://tkzy6.com",
	"https://tkzy7.com",
	"https://tkzy8.com",
	"https://tkzy9.com",
}

func randomSourceDomain() string {
	randInt := rand.Intn(len(sourceDomains))
	return sourceDomains[randInt]
}
