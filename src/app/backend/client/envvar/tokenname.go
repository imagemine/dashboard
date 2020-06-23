package envvar

import "os"

func EnvVariable(envName string, defaultValue string) string {
  if envValue := os.Getenv(envName); envValue != "" {
    return envValue
  }
  return defaultValue
}
