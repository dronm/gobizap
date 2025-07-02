package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dronm/gobizap/md"
	"github.com/hoisie/mustache"
)

const (
	TEMPL_EXT string = ".tmpl" // template extension folder or file for new project initialization
)

func FileExists(fileName string) (bool, error) {
	if _, err := os.Stat(fileName); err != nil && errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}

// helper functions
func readText(prompt, defValue string, obligatory bool, termColors TerminalColors) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	var text string
	var err error
	for obligatory || text == "" {
		if defValue != "" {
			termColors.PromptColor.Printf("%s", prompt)
			termColors.CommentColor.Printf(" (default: %s)-> ", defValue)
		} else {
			termColors.PromptColor.Printf("%s-> ", prompt)
		}
		text, err = reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		text = strings.Replace(text, "\n", "", -1)
		if text == "" && defValue != "" {
			text = defValue
		}
		if obligatory && text == "" {
			termColors.ErrorColor.Println("This parameter can not be empty!")
		} else {
			//any answer will do
			break
		}
	}
	return text, nil
}

func readBool(prompt, defValue string, obligatory bool, termColors TerminalColors) (bool, error) {
	t, err := readText(prompt, defValue, obligatory, termColors)
	if err != nil {
		return false, err
	}
	if strings.ToUpper(t) == "Y" || strings.ToUpper(t) == "YES" {
		return true, nil
	}
	return false, nil
}

// copyProjectFiles recursively creates all project files (folders, files, symbolic links)
// from templates in sourceDir.
func copyProjectFiles(sourceDir, destDir string, params map[string]interface{}) error {
	return filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		rel_path := strings.TrimPrefix(path, sourceDir)
		rel_path = strings.ReplaceAll(rel_path, TEMPL_EXT, "")

		name := d.Name()
		is_template := (strings.HasSuffix(name, TEMPL_EXT) && rel_path != filepath.Join(md.BUILD_DIR, md.TMPL_DIR))
		if strings.Contains(name, "{{") && strings.Contains(name, "}}") {
			//may also have template in name
			rel_path = mustache.Render(rel_path, params)
		}
		new_file := filepath.Join(destDir, rel_path)
		//find out if it is a symbolic link folder and copy its contents if so
		if d.IsDir() {
			//directory
			if _, err := os.Stat(new_file); os.IsNotExist(err) {
				if err := os.Mkdir(new_file, FILE_PERMISSION); err != nil {
					return fmt.Errorf("os.Mkdir() failed: %v", err)
				}
			}

		} else if d.Type()&os.ModeSymlink != 0 {
			//symlink
			exists, err := FileExists(new_file)
			if err != nil {
				return fmt.Errorf("FileExists() for symlink %s failed: %v", new_file, err)
			}
			if !exists {
				target, err := os.Readlink(path)
				if err != nil {
					return fmt.Errorf("os.Readlink() failed: %v", err)
				}
				if err = os.Symlink(target, new_file); err != nil {
					return fmt.Errorf("os.Symlink() failed: %v", err)
				}
			}
		} else {
			//ordinary file
			exists, err := FileExists(new_file)
			if err != nil {
				return fmt.Errorf("FileExists() for normal file %s failed: %v", new_file, err)
			}
			if !exists {
				var new_file_data []byte
				if is_template {
					transformed_data := mustache.RenderFile(path, params)
					new_file_data = []byte(transformed_data)
				} else {
					var err error
					new_file_data, err = os.ReadFile(path)
					if err != nil {
						return fmt.Errorf("os.ReadFile() failed: %v", err)
					}
				}
				err := os.WriteFile(new_file, new_file_data, FILE_PERMISSION)
				if err != nil {
					return fmt.Errorf("os.WriteFile() failed: %v", err)
				}
			}
			//fmt.Println("Writing to file:", new_file)
		}
		return nil
	})
}
