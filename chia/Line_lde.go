// Code generated by ldetool --package chia Line.lde. DO NOT EDIT.

package chia

import (
	"bytes"
)

var constFoundSpace = []byte("Found ")
var constHarvesterSpaceChiaDotHarvesterDotHarvesterColonSpaceINFOSpaces = []byte("harvester chia.harvester.harvester: INFO     ")
var constSpace = []byte(" ")
var constSpacePlots = []byte(" plots")
var constSpacePlotsSpaceWereSpaceEligibleSpaceForSpaceFarming = []byte(" plots were eligible for farming")
var constSpaceProofs = []byte(" proofs")
var constSpaceSDot = []byte(" s.")
var constTimeColonSpace = []byte("Time: ")
var constTotalSpace = []byte("Total ")

// filename: line.lde
//2021-05-21T20:36:36.960 harvester chia.harvester.harvester: INFO     4 plots were eligible for farming 8b1b2e42f0... Found 0 proofs. Time: 0.71902 s. Total 1509 plots
type Line struct {
	Rest       []byte
	Time       []byte
	Plots      []byte
	Proofs     []byte
	ParseTime  []byte
	PlotsCount []byte
}

// Extract ...
func (p *Line) Extract(line []byte) (bool, error) {
	p.Rest = line
	var pos int

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

	// Take until " plots were eligible for farming" as Plots(string)
	pos = bytes.Index(p.Rest, constSpacePlotsSpaceWereSpaceEligibleSpaceForSpaceFarming)
	if pos >= 0 {
		p.Plots = p.Rest[:pos]
		p.Rest = p.Rest[pos+len(constSpacePlotsSpaceWereSpaceEligibleSpaceForSpaceFarming):]
	} else {
		return false, nil
	}

	// Checks if the rest starts with `"Found "` and pass it
	if bytes.HasPrefix(p.Rest, constFoundSpace) {
		p.Rest = p.Rest[len(constFoundSpace):]
	} else {
		return false, nil
	}

	// Take until " proofs" as Proofs(string)
	pos = bytes.Index(p.Rest, constSpaceProofs)
	if pos >= 0 {
		p.Proofs = p.Rest[:pos]
		p.Rest = p.Rest[pos+len(constSpaceProofs):]
	} else {
		return false, nil
	}

	// Checks if the rest starts with `"Time: "` and pass it
	if bytes.HasPrefix(p.Rest, constTimeColonSpace) {
		p.Rest = p.Rest[len(constTimeColonSpace):]
	} else {
		return false, nil
	}

	// Take until " s." as ParseTime(string)
	pos = bytes.Index(p.Rest, constSpaceSDot)
	if pos >= 0 {
		p.ParseTime = p.Rest[:pos]
		p.Rest = p.Rest[pos+len(constSpaceSDot):]
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
