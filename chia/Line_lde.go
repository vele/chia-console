// Code generated by ldetool --package chia Line.lde. DO NOT EDIT.

package chia

import (
	"bytes"
	"fmt"
	"strconv"
	"unsafe"
)

var constDotsSpace = []byte("... ")
var constFoundSpace = []byte("Found ")
var constHarvesterSpaceChiaDotHarvesterDotHarvesterColonSpaceINFOSpaces = []byte("harvester chia.harvester.harvester: INFO     ")
var constPlotsSpaceWereSpaceEligibleSpaceForSpaceFarmingSpace = []byte("plots were eligible for farming ")
var constProofsDotSpace = []byte("proofs. ")
var constSpace = []byte(" ")
var constSpacePlots = []byte(" plots")
var constSpaceSDot = []byte(" s.")
var constTimeColonSpace = []byte("Time: ")
var constTotalSpace = []byte("Total ")

// filename: line.lde
//2021-05-21T20:36:36.960 harvester chia.harvester.harvester: INFO     4 plots were eligible for farming 8b1b2e42f0... Found 0 proofs. Time: 0.71902 s. Total 1509 plots
type Line struct {
	Rest       []byte
	Time       []byte
	Plots      int
	Block      []byte
	Proofs     int
	ParseTime  []byte
	PlotsCount []byte
}

// Extract ...
func (p *Line) Extract(line []byte) (bool, error) {
	p.Rest = line
	var err error
	var pos int
	var tmp []byte
	var tmpInt int64

	// Take until " " as Time(string)
	pos = bytes.Index(p.Rest, constSpace)
	if pos >= 0 {
		p.Time = p.Rest[:pos]
		p.Rest = p.Rest[pos+len(constSpace):]
	} else {
		return false, nil
	}

	// Checks if the rest starts with `"harvester chia.harvester.harvester: INFO     "` and pass it
	if bytes.HasPrefix(p.Rest, constHarvesterSpaceChiaDotHarvesterDotHarvesterColonSpaceINFOSpaces) {
		p.Rest = p.Rest[len(constHarvesterSpaceChiaDotHarvesterDotHarvesterColonSpaceINFOSpaces):]
	} else {
		return false, nil
	}

	// Take until " " as Plots(int)
	pos = bytes.Index(p.Rest, constSpace)
	if pos >= 0 {
		tmp = p.Rest[:pos]
		p.Rest = p.Rest[pos+len(constSpace):]
	} else {
		return false, nil
	}
	if tmpInt, err = strconv.ParseInt(*(*string)(unsafe.Pointer(&tmp)), 10, 64); err != nil {
		return false, fmt.Errorf("parsing `%s` into field Plots(int): %s", string(*(*string)(unsafe.Pointer(&tmp))), err)
	}
	p.Plots = int(tmpInt)

	// Checks if the rest starts with `"plots were eligible for farming "` and pass it
	if bytes.HasPrefix(p.Rest, constPlotsSpaceWereSpaceEligibleSpaceForSpaceFarmingSpace) {
		p.Rest = p.Rest[len(constPlotsSpaceWereSpaceEligibleSpaceForSpaceFarmingSpace):]
	} else {
		return false, nil
	}

	// Take until "... " as Block(string)
	pos = bytes.Index(p.Rest, constDotsSpace)
	if pos >= 0 {
		p.Block = p.Rest[:pos]
		p.Rest = p.Rest[pos+len(constDotsSpace):]
	} else {
		return false, nil
	}

	// Checks if the rest starts with `"Found "` and pass it
	if bytes.HasPrefix(p.Rest, constFoundSpace) {
		p.Rest = p.Rest[len(constFoundSpace):]
	} else {
		return false, nil
	}

	// Take until " " as Proofs(int)
	pos = bytes.Index(p.Rest, constSpace)
	if pos >= 0 {
		tmp = p.Rest[:pos]
		p.Rest = p.Rest[pos+len(constSpace):]
	} else {
		return false, nil
	}
	if tmpInt, err = strconv.ParseInt(*(*string)(unsafe.Pointer(&tmp)), 10, 64); err != nil {
		return false, fmt.Errorf("parsing `%s` into field Proofs(int): %s", string(*(*string)(unsafe.Pointer(&tmp))), err)
	}
	p.Proofs = int(tmpInt)

	// Checks if the rest starts with `"proofs. "` and pass it
	if bytes.HasPrefix(p.Rest, constProofsDotSpace) {
		p.Rest = p.Rest[len(constProofsDotSpace):]
	} else {
		return false, nil
	}

	// Checks if the rest starts with `"Time: "` and pass it
	if bytes.HasPrefix(p.Rest, constTimeColonSpace) {
		p.Rest = p.Rest[len(constTimeColonSpace):]
	} else {
		return false, nil
	}

	// Take until " " as ParseTime(string)
	pos = bytes.Index(p.Rest, constSpace)
	if pos >= 0 {
		p.ParseTime = p.Rest[:pos]
		p.Rest = p.Rest[pos+len(constSpace):]
		fmt.Println("|321321321|")
	} else {
		return false, nil
	}

	// Checks if the rest starts with `" s."` and pass it
	if bytes.HasPrefix(p.Rest, constSpaceSDot) {
		p.Rest = p.Rest[len(constSpaceSDot):]
	} else {
		return false, nil
	}

	// Checks if the rest starts with `"Total "` and pass it
	if bytes.HasPrefix(p.Rest, constTotalSpace) {
		p.Rest = p.Rest[len(constTotalSpace):]
	} else {
		return false, nil
	}

	// Take until " plots" as PlotsCount(string)
	pos = bytes.Index(p.Rest, constSpacePlots)
	if pos >= 0 {
		p.PlotsCount = p.Rest[:pos]
		p.Rest = p.Rest[pos+len(constSpacePlots):]
	} else {
		return false, nil
	}

	return true, nil
}
