package main

import (
	"cli/version"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func main() {
	command := &cobra.Command{
		Use:     version.Name,
		Short:   version.ShortName,
		Version: version.Version,
	}
	command.SetFlagErrorFunc(func(_ *cobra.Command, _ error) error {
		command.Println("unknown command,use -h or --help")
		return nil
	})
	command.AddCommand(ps())
	command.AddCommand(start())
	command.AddCommand(stop())
	command.AddCommand(rm())
	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func ps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps",
		Short: "show docker running containers",
	}
	all := false
	cmd.Flags().BoolVarP(&all, "all", "a",
		false, "all containers")
	cmd.Run = func(cmd *cobra.Command, _ []string) {
		err := psCmd(all)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
	}
	return cmd
}

func start() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start containers for number",
	}
	cmd.Run = func(cmd *cobra.Command, ss []string) {
		err := batchCmd("start", ss)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
	}
	return cmd
}

func stop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "stop containers for number",
	}
	cmd.Run = func(cmd *cobra.Command, ss []string) {
		err := batchCmd("stop", ss)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
	}
	return cmd
}

func rm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm",
		Short: "remove containers for number",
	}
	cmd.Run = func(cmd *cobra.Command, ss []string) {
		err := batchCmd("rm", ss)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
	}
	return cmd
}

func batchCmd(st string, ss []string) error {
	if len(ss) == 0 {
		return errors.New("miss container number")
	}
	s1 := ss[0]
	split := strings.Split(s1, "...")
	ind := 0
	end := 0
	ints := make([]int, 0)
	if len(split) == 2 {
		ind, _ = strconv.Atoi(split[0])
		end, _ = strconv.Atoi(split[1])
		for i := ind - 1; i < end; i++ {
			ints = append(ints, i)
		}
	} else {
		for _, s := range ss {
			atoi, _ := strconv.Atoi(s)
			ints = append(ints, atoi-1)
		}
	}
	open, err := os.Open(fn)
	if err != nil {
		return errors.New("not run ps")
	}
	decoder := gob.NewDecoder(open)
	cs := make([]string, 0)
	err = decoder.Decode(&cs)
	if err != nil {
		return err
	}
	ln := len(cs)
	var wg sync.WaitGroup
	for _, i := range ints {
		if i < ln {
			cn := cs[i]
			wg.Add(1)
			go func() {
				cm := exec.Command("docker", st, cn)
				_, err := pipeCommands(cm)
				wg.Done()
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Println(cn)
			}()
		}
	}
	wg.Wait()
	return nil
}

var fn = "bass.ps"

func psCmd(all bool) error {
	var psc *exec.Cmd
	if all {
		psc = exec.Command("docker", "ps", "-a")
	} else {
		psc = exec.Command("docker", "ps")
	}
	grep := exec.Command("grep", "-v", "k8s")
	awk := exec.Command("awk", "{print $NF}")
	tail := exec.Command("tail", "+2")
	output, err := pipeCommands(psc, grep, awk, tail)
	if err != nil {
		return err
	}
	res := string(output)
	ts := strings.TrimSpace(res)
	if len(ts) == 0 {
		os.RemoveAll(fn)
		return errors.New("empty")
	}
	split := strings.Split(ts, "\n")
	for i, s := range split {
		fmt.Printf("%d.%s\n", i+1, s)
	}
	fl, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(fl)
	err = encoder.Encode(split)
	if err != nil {
		return err
	}
	return nil
}

func pipeCommands(commands ...*exec.Cmd) ([]byte, error) {
	ln := len(commands)
	if ln == 0 {
		return nil, nil
	}
	if ln == 1 {
		return commands[0].Output()
	}
	for i, command := range commands[:ln-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			return nil, err
		}
		err = command.Start()
		if err != nil {
			return nil, err
		}
		commands[i+1].Stdin = out
	}
	final, err := commands[ln-1].Output()
	if err != nil {
		return nil, err
	}
	return final, nil
}
