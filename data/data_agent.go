package data

import (
	"math/rand"
	"strconv"

	"github.com/davidroman0O/gogog/types"
)

func getRandomBrowser() string {
	browsers := []string{"Chrome", "Firefox", "Safari", "Edge"}
	return browsers[rand.Intn(len(browsers))]
}

func getRandomVersion() string {
	major := rand.Intn(10) + 1 // for example, between 1 and 10
	minor := rand.Intn(10)
	patch := rand.Intn(10)
	return strconv.Itoa(major) + "." + strconv.Itoa(minor) + "." + strconv.Itoa(patch)
}

func getRandomPlatform() string {
	platforms := []string{"Windows NT 10.0; Win64; x64", "Macintosh; Intel Mac OS X 10_15_7", "X11; Linux x86_64"}
	return platforms[rand.Intn(len(platforms))]
}

func generateUserAgent() types.UserAgent {
	browser := getRandomBrowser()
	version := getRandomVersion()
	platform := getRandomPlatform()
	return types.UserAgent("Mozilla/5.0 (" + platform + ") " + browser + "/" + version)
}
