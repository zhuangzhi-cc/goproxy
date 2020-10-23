package goproxy

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func rebuildZip(filename, from, to string) string {
	fmt.Println("rebuild zip:", filename, from, to)
	if from != "" {
		f, err := ReArchive(filename, from, to)
		if err != nil {
			log.Println("rebuild zip error:", err)
		} else {
			filename = f
		}
	}
	return filename
}

func rebuildMod(filename, from, to string) string {
	fmt.Println("rebuild mod:", filename, from, to)
	if from != "" {
		f, err := RebuildMod(filename, from, to)
		if err != nil {
			log.Println("rebuild mod error:", err)
		} else {
			filename = f
		}
	}
	return filename
}

func rename(originPath string, renames map[string]string) (newPath, from, to string) {
	for from, to := range renames {
		if strings.Index(originPath, from) >= 0 {
			return strings.ReplaceAll(originPath, from, to), from, to
		}
	}
	return originPath, "", ""
}

func replacePath(modulePath, from, to string) string {
	return strings.ReplaceAll(modulePath, from, to)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func dirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func replaceInFile(file, old, new string) error {
	return Execute("sed", "-i", "-e", fmt.Sprintf("s+%s+%s+g", old, new),
		file)
}

func ReArchive(file, from, to string) (string, error) {
	if !strings.HasSuffix(file, ".zip") {
		return "", fmt.Errorf("expect a zip file, but: %s", file)
	}
	target := strings.ReplaceAll(file, from, to)

	if fileExists(target) {
		return target, nil
	}

	dir, err := ioutil.TempDir("/tmp", "rearchive")
	if err != nil {
		return "", err
	}

	out, err := Unzip(file, dir)
	if err != nil {
		return "", err
	}

	mod := FindModFile(out, dir)

	if mod != "" {
		err = replaceInFile(mod, from, to)
		if err != nil {
			return "", err
		}
	}
	err = MoveFiles(dir, from, to)
	if err != nil {
		return "", err
	}

	infos, err := ioutil.ReadDir(dir + "/" + to)
	if err != nil {
		return "", err
	}
	os.Remove(target)

	targetDir := filepath.Dir(target)
	fmt.Println("targetDir:", targetDir)
	if !dirExists(targetDir) {
		Execute("mkdir", "-p", targetDir)
	}

	err = Zip(target, dir, to+"/"+infos[0].Name())
	os.RemoveAll(dir)
	return target, err
}

func ExecuteAt(dir, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Dir = dir
	out, err := c.Output()
	if err != nil {
		fmt.Printf("command:\n%s\n", c)
		fmt.Printf("output:\n%s\n", string(out))
		fmt.Println("error:\n", err)
	}
	return err
}

func Execute(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	out, err := c.Output()
	if err != nil {
		fmt.Printf("%s\n%s\n", c, string(out))
		fmt.Println(string(err.(*exec.ExitError).Stderr))
	}
	return err
}

func MoveFiles(dir, from, to string) error {
	todir := to
	if id := strings.LastIndex(todir, "/"); id > 0 {
		todir = todir[0:id]
	}
	err := Execute("mkdir", "-p", fmt.Sprintf("%s/%s", dir, todir))
	if err != nil {
		return err
	}

	err = Execute("mv", "-v", fmt.Sprintf("%s/%s", dir, from), fmt.Sprintf("%s/%s", dir, to))
	if err != nil {
		return err
	}

	return nil
}

func Zip(zip, folder string, args ...string) error {
	fmt.Println("zip:", zip, folder)
	return ExecuteAt(folder, "zip", append([]string{"-r", zip}, args...)...)
}

func Unzip(src, dest string) (string, error) {
	cmd := exec.Command("unzip", src, "-d", dest)
	out, err := cmd.Output()
	return string(out), err
}

func FindModFile(x, dir string) string {
	scanner := bufio.NewScanner(strings.NewReader(x))
	for scanner.Scan() {
		l := scanner.Text()
		if id := strings.Index(l, "go.mod"); id > 0 {
			f := strings.Index(l, dir)
			if f > 0 {
				return l[f : id+6]
			}
		}
	}
	return ""
}

func FindFileName(x, dir string) string {
	scanner := bufio.NewScanner(strings.NewReader(x))
	for scanner.Scan() {
		l := scanner.Text()
		if id := strings.Index(l, "go.mod"); id > 0 {
			f := strings.Index(l, dir)
			if f > 0 {
				return l[f : id+6]
			}
		}
	}
	return ""
}

func RebuildMod(file, from, to string) (string, error) {
	if !strings.HasSuffix(file, ".mod") {
		return "", fmt.Errorf("expect a mod file, but: %s", file)
	}
	target := strings.ReplaceAll(file, from, to)

	if fileExists(target) {
		return target, nil
	}

	dir := filepath.Dir(target)

	if !dirExists(dir) {
		Execute("mkdir", "-p", dir)
	}

	err := Execute("cp", file, target)
	if err != nil {
		fmt.Println("cp failed!!")
		return "", err
	}
	err = replaceInFile(target, from, to)
	return target, err
}
