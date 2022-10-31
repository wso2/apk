package base

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

//Generate random strings with given length
func GenerateRandomName(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// Execute : Run apictl command
//
func Execute(t *testing.T, args ...string) (string, error) {
	path := filepath.Join(RelativeBinaryPath, BinaryName)
	cmd := exec.Command(path, args...)

	t.Log("base.Execute() - apkctl command:", cmd.String())
	// run command
	output, err := cmd.Output()

	t.Log("base.Execute() - apkctl command output:", string(output))
	return string(output), err
}

func ExecuteKubernetesCommands(args ...string) (string, error) {
	cmd := exec.Command(K8sBinaryName, args...)

	// // run command
	// output, err := cmd.Output()

	var errBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)

	output, err := cmd.Output()
	return string(output), err

	// return string(output), err
}

// func GetExportedPathFromOutput(output string) string {
// 	//Check directory path to omit changes due to OS differences
// 	if strings.Contains(output, ":\\") {
// 		arrayOutput := []rune(output)
// 		extractedPath := string(arrayOutput[strings.Index(output, ":\\")-1:])
// 		return strings.ReplaceAll(strings.ReplaceAll(extractedPath, "\n", ""), " ", "")
// 	} else {
// 		return strings.ReplaceAll(strings.ReplaceAll(output[strings.Index(output, string(os.PathSeparator)):], "\n", ""), " ", "")
// 	}
// }

func GetExportedPathFromOutput(output string) string {
	out := strings.ReplaceAll(output, " ", "")
	arrayOutput := []rune(out)
	extractedPath := string(arrayOutput[strings.Index(out, ":")+1:])
	return strings.ReplaceAll(strings.ReplaceAll(extractedPath, "\n", ""), " ", "")
}

// IsFileAvailable checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func IsFileAvailable(t *testing.T, filepath string) bool {
	t.Log("base.IsFileAvailable() - API file path:", filepath)

	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
