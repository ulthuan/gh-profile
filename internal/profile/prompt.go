package profile

import (
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gabe565/gh-profile/internal/github"
	"os"
	"path/filepath"
)

func List() ([]Profile, error) {
	conf := github.ConfigDir()

	files, err := os.ReadDir(filepath.Join(conf, "profiles"))
	if err != nil {
		return []Profile{}, err
	}

	profiles := make([]Profile, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			profiles = append(profiles, Profile{Name: file.Name()})
		}
	}

	return profiles, nil
}

var ErrUnknownProfile = errors.New("unknown profile")

func Select(message string) (Profile, error) {
	profiles, err := List()
	if err != nil {
		return Profile{}, err
	}

	profilesStr := make([]string, 0, len(profiles))
	var defaultName string
	for _, profile := range profiles {
		if profile.IsActive() {
			defaultName = profile.Name
		}
		profilesStr = append(profilesStr, profile.Name)
	}

	var answer string
	if err := survey.AskOne(&survey.Select{
		Message: message,
		Options: profilesStr,
		Default: defaultName,
	}, &answer, survey.WithValidator(survey.Required)); err != nil {
		return Profile{}, err
	}

	for _, profile := range profiles {
		if profile.Name == answer {
			return profile, nil
		}
	}

	return Profile{}, ErrUnknownProfile
}

func PromptNew() (Profile, error) {
	var answer string
	err := survey.AskOne(&survey.Input{
		Message: "Enter new profile name:",
	}, &answer, survey.WithValidator(survey.Required))
	return New(answer), err
}
