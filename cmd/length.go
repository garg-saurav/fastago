/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"log"

	"github.com/lucblassel/fastago/pkg/seqs"
	"github.com/spf13/cobra"
)

var lengthMode string

// lengthCmd represents the length command
var lengthCmd = &cobra.Command{
	Use:   "length",
	Short: "get length of sequences in fasta file",
	Long: `By default this command outputs the length of each sequence. 
	It is also possible to retreive the minimum, maximum or average length.
	set the -m/--mode flag to either:
		- each: will display the length of each sequence
		- average (or mean) will display the average lengh
		- min (or minimum) will display the minimum length
		- max (or maximum) will display the maximum length`,
	Run: func(cmd *cobra.Command, args []string) {

		records := make(chan seqs.SeqRecord)
		errs := make(chan error)

		go seqs.ReadFastaRecords(inputReader, records, errs)

		switch lengthMode {
		case "each":
			getEach(records, errs, outputWriter)
		case "average":
			getAverage(records, errs, outputWriter)
		case "mean":
			getAverage(records, errs, outputWriter)
		case "min":
			getMin(records, errs, outputWriter)
		case "minimum":
			getMin(records, errs, outputWriter)
		case "max":
			getMax(records, errs, outputWriter)
		case "maximum":
			getMax(records, errs, outputWriter)
		default:
			log.Fatalf(
				"Mode %s not recognized.\n"+
					"The mode must be one of the following values: "+
					"'each' 'average' 'mean' 'min' 'minimum' 'max' 'maximum'", lengthMode)
		}
	},
}

func init() {
	statsCmd.AddCommand(lengthCmd)

	lengthCmd.Flags().StringVarP(&lengthMode, "mode", "m", "each", "How to display lengths")
}

type Record struct {
	Name   string
	Length int
}

func getEach(records chan seqs.SeqRecord, errs chan error, output io.Writer) error {

	for records != nil && errs != nil {
		select {
		case record := <-records:
			fmt.Fprintf(output, "%s\t%d\n", record.Name, record.Seq.Length())
		case err := <-errs:
			return err
		}
	}

	return nil
}

func getAverage(records chan seqs.SeqRecord, errs chan error, output io.Writer) error {
	total, count := 0, 0

	for records != nil && errs != nil {
		select {
		case record := <-records:
			total += record.Seq.Length()
			count++
		case err := <-errs:
			if err != nil {
				return err
			}
			fmt.Fprintln(output, float32(total)/float32(count))
			return nil
		}
	}
	return nil
}

func getMin(records chan seqs.SeqRecord, errs chan error, output io.Writer) error {
	min := -1

	for records != nil && errs != nil {
		select {
		case record := <-records:
			if min < 0 || record.Seq.Length() < min {
				min = record.Seq.Length()
			}
		case err := <-errs:
			if err != nil {
				return err
			}
			fmt.Fprintln(output, min)
			return nil
		}
	}
	return nil
}

func getMax(records chan seqs.SeqRecord, errs chan error, output io.Writer) error {
	max := -1

	for records != nil && errs != nil {
		select {
		case record := <-records:
			if max < 0 || record.Seq.Length() > max {
				max = record.Seq.Length()
			}
		case err := <-errs:
			if err != nil {
				return err
			}
			fmt.Fprintln(output, max)
			return nil
		}
	}
	return nil
}
