package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var ProcessCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the entire health report from the support packet.",
	Long:  "Generates the entire health report from the support packet, and outputting a pdf file.",
	RunE:  generateCmdF,
}

func init() {
	ProcessCmd.Flags().StringP("packet", "p", "", "the support packet file to process")
	ProcessCmd.Flags().StringP("output", "o", "healthcheck-report.pdf", "the output file name for the PDF.")

	ProcessCmd.Flags().Bool("debug", true, "Whether to show debug logs or not")

	if err := ProcessCmd.MarkFlagRequired("packet"); err != nil {
		panic(err)
	}

	RootCmd.AddCommand(
		ProcessCmd,
	)
}

func generateCmdF(cmd *cobra.Command, args []string) error {
	supportPacketFile, _ := cmd.Flags().GetString("packet")
	outputFileName, _ := cmd.Flags().GetString("output")

	//validate the packet file exists

	_, err := os.Stat(supportPacketFile)

	if err != nil {
		return errors.Wrap(err, "failed to find the support packet file")
	}

	cmdArgs := []string{"compose", "run", "--rm", "mm-healthcheck", "generate", "--packet", supportPacketFile, "--output", outputFileName}

	generate := exec.Command("docker", cmdArgs...)
	stdout, _ := generate.StdoutPipe()
	stderr, _ := generate.StderrPipe()

	err = generate.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start the command")
	}

	go copyOutput(stdout)
	go copyOutput(stderr)

	err = generate.Wait()
	if err != nil {
		return errors.Wrap(err, "failed to wait for the command to finish")
	}

	return nil

}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
