package memory

import (
	"bufio"
	"context"
	"os/exec"
	"time"
)

func Get(filepath string, arg ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, filepath, arg...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	err1 := cmd.Start()
	if err1 != nil {
		return "", err1
	}

	scanner := bufio.NewScanner(out)
	var res string
	for scanner.Scan() {
		res += scanner.Text() + "\n"
	}

	return res, nil
}
