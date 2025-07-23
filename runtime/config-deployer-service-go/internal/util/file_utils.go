/* Copyright (c) 2025 WSO2 LLC. (http://www.wso2.com) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// RetrieveDefinitionFromUrl retrieves the API definition from the given URL.
func RetrieveDefinitionFromUrl(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error occurred while retrieving the definition from the url: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error occurred while closing the response body: %v\n", err)
		}
	}(response.Body)
	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return "", fmt.Errorf("error occurred while reading the definition from the url: %w", err)
		}
		return string(body), nil
	} else {
		return "", fmt.Errorf("error occurred while retrieving the definition from the url: %s. Status code:"+
			" %d", url, response.StatusCode)
	}
}

// ExtractUploadedArchive extracts the uploaded archive from a byte array and creates the necessary directory structure
func ExtractUploadedArchive(byteArray []byte, importedDirectoryName, apiArchiveLocation, extractLocation string) (string, error) {
	var archiveExtractLocation string

	// Create api import directory structure
	if err := createDirectory(extractLocation); err != nil {
		return "", fmt.Errorf("error creating directory %s: %w", extractLocation, err)
	}

	// Create archive from byte array
	if err := createArchiveFromByteArray(byteArray, apiArchiveLocation); err != nil {
		return "", fmt.Errorf("error creating archive from byte array: %w", err)
	}

	// Extract the archive
	archiveExtractLocation = filepath.Join(extractLocation, importedDirectoryName)
	if _, err := extractArchive(apiArchiveLocation, archiveExtractLocation); err != nil {
		return "", fmt.Errorf("error extracting archive %s: %w", apiArchiveLocation, err)
	}

	return archiveExtractLocation, nil
}

// createDirectory creates a directory at the specified path with the given permissions.
func createDirectory(path string) error {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return fmt.Errorf("error in creating directory at %s: %v\n", path, err)
	}
	return nil
}

// createArchiveFromByteArray creates an archive file from a byte array and saves it to the specified path.
func createArchiveFromByteArray(byteArray []byte, archivePath string) error {
	outFile, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("error in Creating archive from data. %s: %v\n", archivePath, err)
	}
	defer outFile.Close()

	_, err = outFile.Write(byteArray)
	return err
}

// extractArchive extracts the contents of a zip archive to the specified destination directory.
func extractArchive(archiveFilePath, destination string) (string, error) {
	const (
		maxEntryCount = 1024
		sizeLimit     = 0x6400000 // 100MB
		bufferSize    = 512
	)

	var archiveName string
	var entries int
	var total int64

	// Open the zip file
	reader, err := zip.OpenReader(archiveFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	// Process each entry
	for i, file := range reader.File {
		entries++
		if entries > maxEntryCount {
			return "", fmt.Errorf("too many files to unzip")
		}

		// Get archive name from first entry
		if i == 0 && strings.Contains(file.Name, "/") {
			archiveName = strings.Split(file.Name, "/")[0]
		}

		// Create destination file path
		destPath := filepath.Join(destination, file.Name)

		// Security check: prevent zip slip
		if !strings.HasPrefix(destPath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return "", fmt.Errorf("attempt to upload invalid zip archive with file at %s. File path is outside target directory", file.Name)
		}

		// Create directory if needed
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, file.FileInfo().Mode()); err != nil {
				return "", fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			continue
		}

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return "", fmt.Errorf("failed to create parent directory for %s: %w", destPath, err)
		}

		// Open file in zip
		rc, err := file.Open()
		if err != nil {
			return "", fmt.Errorf("failed to open file %s in zip: %w", file.Name, err)
		}

		// Create destination file
		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.FileInfo().Mode())
		if err != nil {
			rc.Close()
			return "", fmt.Errorf("failed to create destination file %s: %w", destPath, err)
		}

		// Copy with size limit
		buffer := make([]byte, bufferSize)
		for {
			if total+bufferSize > sizeLimit {
				rc.Close()
				destFile.Close()
				return "", fmt.Errorf("file being unzipped is too big")
			}

			n, err := rc.Read(buffer)
			if err != nil && err != io.EOF {
				rc.Close()
				destFile.Close()
				return "", fmt.Errorf("failed to read from zip file: %w", err)
			}

			if n == 0 {
				break
			}

			if _, err := destFile.Write(buffer[:n]); err != nil {
				rc.Close()
				destFile.Close()
				return "", fmt.Errorf("failed to write to destination file: %w", err)
			}

			total += int64(n)
		}

		rc.Close()
		destFile.Close()
	}

	return archiveName, nil
}

// DeleteDirectory deletes the specified directory and all its contents.
func DeleteDirectory(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("Failed to delete directory %s: %v\n", path, err)
	}
	return nil
}

// FileExists checks if a file exists at the specified path.
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// MarshalToYAMLWithIndent marshals a struct to YAML with custom indentation
func MarshalToYAMLWithIndent(data interface{}, indent int) ([]byte, error) {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(indent)
	err := encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("error occurred while encoding to YAML: %w", err)
	}
	encoder.Close()
	return buf.Bytes(), nil
}
