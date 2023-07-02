package checksum

import (
	"bufio"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func InitCheckSum(nameFile string) error {
	cmd := exec.Command("pwd")
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	m, err := os.Open(string(out[:len(out)-1]) + "/" + nameFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err.Error())
	}

	r := bufio.NewReader(m)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	table := crc32.MakeTable(0xD5828281)

	sum := crc32.Checksum(b, table)

	checkSum = strconv.Itoa(int(sum))
	if checkSum == "" {
		return errors.New("failed to get checksum")
	}

	return nil
}

var checkSum string
var ticker = time.NewTicker(time.Second)

func PeriodicCheckSum(nameFile string) {
	for {
		select {
		case <-ticker.C:
			cmd := exec.Command("pwd")
			out, err := cmd.Output()
			if err != nil {
				log.Fatal(err.Error())
			}

			m, err := os.Open(string(out[:len(out)-1]) + "/" + nameFile)
			if err != nil {
				log.Fatal(err.Error())
			}

			r := bufio.NewReader(m)
			b, err := ioutil.ReadAll(r)
			if err != nil {
				log.Fatal(err.Error())
			}

			table := crc32.MakeTable(0xD5828281)

			sum := crc32.Checksum(b, table)
			s := strconv.Itoa(int(sum))
			if s != checkSum {
				log.Fatal("binary file is not legitimate")
			}
		}
	}
}
