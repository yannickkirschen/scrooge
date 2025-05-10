package editor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"text/template"

	"gopkg.in/yaml.v3"
)

func EditTplTempStructFile[T any](text string, out *T, tplDataNew, tplDataEdit any, funcMap template.FuncMap, save func(T) error) error {
	// 1. Parse template and save to temporary file.
	newTpl, err := template.New("tpl").Funcs(funcMap).Parse(text)
	if err != nil {
		return fmt.Errorf("error parsing template for new temporary file: %s", err)
	}

	w, err := os.CreateTemp("", "*.yaml")
	if err != nil {
		return fmt.Errorf("error opening temporary file for writing: %s", err)
	}

	filename := w.Name()
	if err := newTpl.Execute(w, tplDataNew); err != nil {
		return fmt.Errorf("error executing template for new temporary file %s: %s", filename, err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("error closing temporary file %s: %s", filename, err)
	}

	// 2. Edit the file.
	if err := Edit(GetExe(), filename); err != nil {
		return fmt.Errorf("error editing temporary file %s: %s", filename, err)
	}

	// 3. Read the file and parse to struct.
	r, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening %s for reading: %s", filename, err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("error reading content from %s to byte array: %s", filename, err)
	}

	editedTpl := template.New("tpl").Funcs(funcMap)
	editedTpl, err = editedTpl.Parse(string(b))
	if err != nil {
		return fmt.Errorf("error parsing template edited-transaction: %s", err)
	}

	var edited bytes.Buffer
	if err := editedTpl.Execute(&edited, tplDataEdit); err != nil {
		return fmt.Errorf("error executing template edited-transaction: %s", err)
	}

	if err := yaml.Unmarshal(edited.Bytes(), out); err != nil {
		return fmt.Errorf("error parsing YAML from %s: %s", filename, err)
	}

	// 4. Save the result.
	if err := save(*out); err != nil {
		return fmt.Errorf("error saving data from %s: %s", filename, err)
	}

	os.Remove(filename)
	return nil
}

func EditTempStringFile(text string, save func(string) error) error {
	// 1. Save text in temporary file.
	w, err := os.CreateTemp("", "*.expr")
	if err != nil {
		return fmt.Errorf("error opening temporary file for writing: %s", err)
	}

	w.WriteString(text)

	filename := w.Name()
	if err := w.Close(); err != nil {
		return fmt.Errorf("error closing temporary file %s: %s", filename, err)
	}

	// 2. Edit the file.
	if err := Edit(GetExe(), filename); err != nil {
		return fmt.Errorf("error editing temporary file %s: %s", filename, err)
	}

	// 3. Read the file.
	r, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening %s for reading: %s", filename, err)
	}
	defer r.Close()

	b, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("error reading content from %s to byte array: %s", filename, err)
	}

	// 4. Save the result.
	if err := save(string(b)); err != nil {
		return fmt.Errorf("error saving data from %s: %s", filename, err)
	}

	os.Remove(filename)
	return nil
}

func Edit(exe, filename string) error {
	cmd := exec.Command(exe, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting %s %s: %s", exe, filename, err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("error waiting for %s %s: %s", exe, filename, err)
	}

	return nil
}

func GetExe() string {
	exe, ok := os.LookupEnv("EDITOR")
	if !ok {
		exe = "vim"
	}

	return exe
}
