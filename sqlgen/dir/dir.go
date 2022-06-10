package dir

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

func CheckDirExists(outputDir string) (bool, error) {
	override := true
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		prompt := survey.Confirm{
			Message: fmt.Sprintf("%s folder exists, existing files will be replaced. continue?", outputDir),
		}
		e := survey.AskOne(&prompt, &override)

		if e != nil {
			return false, e
		}
	}
	return override, nil
}
