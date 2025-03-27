package utils

import (
	"cw/logger"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
)

func UserChoice() string {
	mainMenu := []string{
		"1. CexWithdrawer",
		"2. Bridger",
		"3. Cex_Bridger",
		"4. WalletGenerator",
		"5. Сollector",
		"0. Exit",
	}

	var rgx = regexp.MustCompile(`^\d+\.\s*`)

	for {
		selected := promptSelection("Choose module:", mainMenu)
		selected = rgx.ReplaceAllString(selected, "")

		switch selected {
		case "CexWithdrawer", "Bridger", "Cex_Bridger", "WalletGenerator", "Сollector":
			return selected
		case "Exit":
			logger.GlobalLogger.Infof("Exiting program.")
			return ""
		default:
			logger.GlobalLogger.Warnf("Invalid selection: %s", selected)
		}
	}
}

func promptSelection(message string, options []string) string {
	var selected string
	if err := survey.AskOne(&survey.Select{
		Message: message,
		Options: options,
		Default: options[len(options)-1],
	}, &selected); err != nil {
		logger.GlobalLogger.Errorf("Error selecting option: %v", err)
		return ""
	}
	return selected
}

// func handleSubMenu(menuName string, subMenu []string, rgx *regexp.Regexp) string {
// 	for {
// 		selected := promptSelection("Choose "+menuName+" sub-module:", subMenu)
// 		selected = rgx.ReplaceAllString(selected, "")

// 		if selected == "Back" {
// 			return ""
// 		}
// 		return selected
// 	}
// }
