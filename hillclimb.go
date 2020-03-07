package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/emedvedev/enigma"
	// "bufio"
	// "text/template"
	// "github.com/mkideal/cli"
)

//YOUR PLUGBOARD IS WRONG!! CHANGE THAT ANSWER WRITTEN QUESTION

var enigmaSettings = struct {
	Reflector string
	Rings     []int
	Positions []string
	Rotors    []string
}{
	Reflector: "C-thin",
	Rings:     []int{1, 1, 1, 16},
	Positions: []string{"A", "A", "B", "Q"},
	Rotors:    []string{"Beta", "II", "IV", "III"},
}

// var defaultAplhabets = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
var defaultAplhabets = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var allRotors = []string{"I", "II", "V", "VI", "Beta", "Gamma"}
var trigramScores = make(map[string]float64)
var totalIOCScore = 0.0
var averageIOCScore = 0.0
var totalNumbers = 1.0

func setRotor1(Rotor1Value string) {
	enigmaSettings.Rotors[0] = Rotor1Value
}

func setRotor2(Rotor2Value string) {
	enigmaSettings.Rotors[1] = Rotor2Value
}

func setPosition1(Position1Value string) {
	enigmaSettings.Positions[0] = Position1Value
}

func setPosition2(Position2Value string) {
	enigmaSettings.Positions[1] = Position2Value
}

func getRotor1() string {
	return string(enigmaSettings.Rotors[0])
}

func getRotor2() string {
	return string(enigmaSettings.Rotors[1])
}

func getPosition1() string {
	return string(enigmaSettings.Positions[0])
}

func getPosition2() string {
	return string(enigmaSettings.Positions[1])
}

func createAndEncode(plaintext string, plugboardvalue []string) string {

	config := make([]enigma.RotorConfig, len(enigmaSettings.Rotors))

	for index, rotor := range enigmaSettings.Rotors {
		ring := enigmaSettings.Rings[index]
		value := enigmaSettings.Positions[index][0]
		config[index] = enigma.RotorConfig{ID: rotor, Start: value, Ring: ring}
	}
	e := enigma.NewEnigma(config, enigmaSettings.Reflector, plugboardvalue)
	encoded := e.EncodeString(plaintext)

	return string(encoded)
}

func getIOC(encoded string) float64 {

	ioc := 0.0
	n := float64(len(encoded))

	for letterASCII := 65; letterASCII < 91; letterASCII++ {
		i := string(letterASCII)
		fi := int(strings.Count(encoded, i))
		ioc += float64(fi * (fi - 1))
	}

	ioc = float64(ioc / (n * (n - 1)))

	return ioc
}

func swapCharacters(old string, new string, tempPlugboard string) string {

	// fmt.Println(tempPlugboard, old, new)
	tempPlugboard = strings.ReplaceAll(tempPlugboard, old, "1")
	tempPlugboard = strings.ReplaceAll(tempPlugboard, new, old)
	tempPlugboard = strings.ReplaceAll(tempPlugboard, "1", new)

	return tempPlugboard
}

func createEnigmaPlugboard(currentPlugboardSettings string) []string {

	var i int
	var enigmaPlugboardSetting []string
	var tempDefaultAlphabets string
	tempDefaultAlphabets = defaultAplhabets
	for i = 0; i < len(currentPlugboardSettings); i++ {
		if currentPlugboardSettings[i] != tempDefaultAlphabets[i] && string(currentPlugboardSettings[i]) != "-" && string(tempDefaultAlphabets[i]) != "-" {
			enigmaPlugboardSetting = append(enigmaPlugboardSetting, string(currentPlugboardSettings[i])+string(tempDefaultAlphabets[i]))
			var x = string(currentPlugboardSettings[i])
			var y = string(tempDefaultAlphabets[i])
			currentPlugboardSettings = strings.ReplaceAll(currentPlugboardSettings, x, "-")
			currentPlugboardSettings = strings.ReplaceAll(currentPlugboardSettings, y, "-")
			tempDefaultAlphabets = strings.ReplaceAll(tempDefaultAlphabets, x, "-")
			tempDefaultAlphabets = strings.ReplaceAll(tempDefaultAlphabets, y, "-")
		}
	}
	//fmt.Println(enigmaPlugboardSetting)
	return enigmaPlugboardSetting
}

func iocScoringSystem(data string) string {

	var i int
	var j int
	var tempPlugboardSettings string
	var tempPlugboardSettings1 string
	var tempPlugboardSettings2 string
	var tempPlugboardSettings3 string
	var tempPlugboardSettings4 string
	var currentPlugboardSettings string
	var bestPlugboardSettings string
	var currentTempPlugboardSettings string
	var localIOC float64
	var maximumIOC float64
	var enigmaPlugboard []string
	var decrypted string

	bestPlugboardSettings = defaultAplhabets
	currentPlugboardSettings = defaultAplhabets
	maximumIOC = 0.0

	for i = 0; i < 26; i++ {

		currentPlugboardSettings = bestPlugboardSettings
		localIOC = 0.0
		for j = i + 1; j < 26; j++ {

			// fmt.Println(currentPlugboardSettings[j], defaultAplhabets[j])
			if string(currentPlugboardSettings[j]) != string(defaultAplhabets[j]) {

				tempPlugboardSettings = swapCharacters(string(defaultAplhabets[j]), string(currentPlugboardSettings[j]), currentPlugboardSettings)
				tempPlugboardSettings = swapCharacters(string(defaultAplhabets[i]), string(currentPlugboardSettings[i]), tempPlugboardSettings)

				tempPlugboardSettings1 = swapCharacters(string(defaultAplhabets[i]), string(currentPlugboardSettings[i]), tempPlugboardSettings)
				tempPlugboardSettings2 = swapCharacters(string(defaultAplhabets[j]), string(currentPlugboardSettings[j]), tempPlugboardSettings)
				tempPlugboardSettings3 = swapCharacters(string(defaultAplhabets[i]), string(currentPlugboardSettings[i]), tempPlugboardSettings)
				tempPlugboardSettings4 = swapCharacters(string(defaultAplhabets[j]), string(currentPlugboardSettings[j]), tempPlugboardSettings)

				// fmt.Println("temp1:" + tempPlugboardSettings1)
				// fmt.Println("temp2:" + tempPlugboardSettings2)
				// fmt.Println("temp3:" + tempPlugboardSettings3)
				// fmt.Println("temp4:" + tempPlugboardSettings4)

				var enigmaPlugboard1 = createEnigmaPlugboard(tempPlugboardSettings1)
				var decrypted1 = createAndEncode(data, enigmaPlugboard1)
				var localIOC1 = getIOC(decrypted1)

				var enigmaPlugboard2 = createEnigmaPlugboard(tempPlugboardSettings2)
				var decrypted2 = createAndEncode(data, enigmaPlugboard2)
				var localIOC2 = getIOC(decrypted2)

				var enigmaPlugboard3 = createEnigmaPlugboard(tempPlugboardSettings3)
				var decrypted3 = createAndEncode(data, enigmaPlugboard3)
				var localIOC3 = getIOC(decrypted3)

				var enigmaPlugboard4 = createEnigmaPlugboard(tempPlugboardSettings4)
				var decrypted4 = createAndEncode(data, enigmaPlugboard4)
				var localIOC4 = getIOC(decrypted4)

				// fmt.Println(localIOC1)
				// fmt.Println(localIOC2)
				// fmt.Println(localIOC3)
				// fmt.Println(localIOC4)

				if localIOC1 > localIOC2 && localIOC1 > localIOC3 && localIOC1 > localIOC4 {

					localIOC = localIOC1
					tempPlugboardSettings = tempPlugboardSettings1
				} else if localIOC2 > localIOC1 && localIOC2 > localIOC3 && localIOC2 > localIOC4 {

					localIOC = localIOC2
					tempPlugboardSettings = tempPlugboardSettings2
				} else if localIOC3 > localIOC2 && localIOC3 > localIOC1 && localIOC3 > localIOC4 {
					localIOC = localIOC3
					tempPlugboardSettings = tempPlugboardSettings3
				} else {
					localIOC = localIOC4
					tempPlugboardSettings = tempPlugboardSettings4
				}

			} else {

				tempPlugboardSettings = swapCharacters(string(defaultAplhabets[i]), string(currentPlugboardSettings[j]), currentPlugboardSettings)
				enigmaPlugboard = createEnigmaPlugboard(tempPlugboardSettings)
				decrypted = createAndEncode(data, enigmaPlugboard)
				localIOC = getIOC(decrypted)
			}

			if localIOC > maximumIOC {

				maximumIOC = localIOC
				currentTempPlugboardSettings = tempPlugboardSettings

				// fmt.Println(localIOC, tempPlugboardSettings)
			}

		}

		bestPlugboardSettings = currentTempPlugboardSettings

	}

	return bestPlugboardSettings
}

func setTrigramMap() {

	var i int

	file, err := os.Open("english_trigrams.txt")
	if err != nil {
		log.Fatal(err)
	}

	dataBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var wholeText = string(dataBytes)
	var textSet = strings.Split(wholeText, "\n")
	var sum float64

	for i = 0; i < len(textSet)-1; i++ {
		var k = strings.Split(textSet[i], " ")
		var value, err = strconv.Atoi(k[1])
		if err != nil {
			log.Fatal(err)
		}
		trigramScores[k[0]] = float64(value)
		sum += float64(value)
	}

	for i = 0; i < len(textSet)-1; i++ {
		var k = strings.Split(textSet[i], " ")
		trigramScores[k[0]] = math.Log(trigramScores[k[0]] / sum)
	}
}

func trigramScoringSystem(data, iocPlugboardSetting string) float64 {

	var enigmaPlugboard []string
	var decrypted string
	var sum = 0.0
	var i int
	enigmaPlugboard = createEnigmaPlugboard(iocPlugboardSetting)
	decrypted = createAndEncode(data, enigmaPlugboard)

	setTrigramMap()
	for i = 0; i < len(decrypted)-3; i++ {
		sum += float64(trigramScores[decrypted[i:i+3]])
	}

	return sum
}

func scoreSettings(data string) (string, float64) {

	var iocPlugboardSetting string
	var trigramPlugboardScore float64
	// var localIOCScore float64
	// var enigmaPlugboard []string
	// var decrypted string

	// enigmaPlugboard = createEnigmaPlugboard(defaultAplhabets)
	// decrypted = createAndEncode(data, enigmaPlugboard)
	// localIOCScore = getIOC(decrypted)

	// averageIOCScore = totalIOCScore / totalNumbers

	// fmt.Println(averageIOCScore)
	// fmt.Println(localIOCScore)

	// if localIOCScore > averageIOCScore {

	// 	totalIOCScore += localIOCScore
	// 	totalNumbers++

	// }

	iocPlugboardSetting = iocScoringSystem(data)
	trigramPlugboardScore = trigramScoringSystem(data, iocPlugboardSetting)

	// fmt.Println(trigramPlugboardScore)
	return iocPlugboardSetting, trigramPlugboardScore

}

func main() {

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	dataBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var setupScore float64
	var setupPlugboard string
	var maxScore = -50000.0
	var bestRotor1 string
	var bestRotor2 string
	var bestPosition1 string
	var bestPosition2 string
	var bestPlugboard string
	var bestEnigmaPlugboard []string
	var text string

	plaintext := string(dataBytes)
	var i, j, m, n int

	// fmt.Println(text)
	// text = createAndEncode(text, enigmaPlugboard)
	// fmt.Println(text)

	for i = 0; i < len(allRotors); i++ {

		setRotor1(allRotors[i])
		for j = 0; j < len(allRotors); j++ {

			setRotor2(allRotors[j])
			fmt.Println("Iteration: "+getRotor1(), getRotor2())
			for m = 0; m < 26; m++ {

				setPosition1(string(m + 65))
				for n = 0; n < 26; n++ {

					setPosition2(string(n + 65))

					setupPlugboard, setupScore = scoreSettings(plaintext)
					// fmt.Println(getPosition1(), getPosition2(), setupScore)

					if setupScore > maxScore {

						maxScore = setupScore
						bestPlugboard = setupPlugboard
						bestRotor1 = getRotor1()
						bestRotor2 = getRotor2()
						bestPosition1 = getPosition1()
						bestPosition2 = getPosition2()
						// fmt.Println(bestRotor1, bestRotor2, bestPosition1, bestPosition2, bestPlugboard)
					}

				}
			}
		}
	}

	fmt.Println("Enigma Setup:")
	fmt.Println(bestRotor1, bestRotor2, bestPosition1, bestPosition2, bestPlugboard)

	setRotor1(bestRotor1)
	setRotor2(bestRotor2)
	setPosition1(bestPosition1)
	setPosition2(bestPosition2)
	bestEnigmaPlugboard = createEnigmaPlugboard(bestPlugboard)
	text = createAndEncode(plaintext, bestEnigmaPlugboard)
	fmt.Println("Plain Text:")
	fmt.Println(text)

}
