package initializers

import "github.com/joho/godotenv"

func LoadEnvVariables() error {
	err := godotenv.Load()

	if err != nil {
		return nil
	}

	return nil
}
