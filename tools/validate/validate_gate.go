package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func epicSkippedForArtefactGate(epicsPath, epic string) (skipped bool, status string, err error) {
	epicDir := filepath.Join(epicsPath, epic)
	status, err = readEpicScopeStatus(epicDir)
	if err != nil {
		return false, "", err
	}
	return isSkippedEpicStatus(status), status, nil
}

func runFullValidation(epic string, jsonOut bool) {
	if jsonOut {
		errLog.Printf("Error: --json with default validate gate is not supported; use validate ac|ears|req [EP-XXX] --json\n")
		os.Exit(1)
	}

	ok := acValidationPasses(epic, false)
	if !earsValidationPasses(epic, false) {
		ok = false
	}
	if !reqValidationPasses(epic, false) {
		ok = false
	}
	if !ok {
		os.Exit(1)
	}
}

func acValidationPasses(epic string, jsonOut bool) bool {
	if epic == "all" || epic == "" {
		return validateAllEpics(jsonOut)
	}
	return validateSingleEpic(epic, jsonOut)
}

func earsValidationPasses(epic string, jsonOut bool) bool {
	if epic == "all" || epic == "" {
		return validateAllEARS(jsonOut)
	}
	return validateSingleEpicEARS(epic, jsonOut)
}

func reqValidationPasses(epic string, jsonOut bool) bool {
	if epic == "all" || epic == "" {
		return validateAllREQ(jsonOut)
	}
	return validateSingleEpicREQ(epic, jsonOut)
}

func validateSingleEpicEARS(epic string, jsonOut bool) bool {
	epicsPath, err := epicArtefactsRoot()
	if err != nil {
		errLog.Printf("Error: %v\n", err)
		return false
	}
	skipped, status, err := epicSkippedForArtefactGate(epicsPath, epic)
	if err != nil {
		errLog.Printf("Error validating %s: %v\n", epic, err)
		return false
	}
	if skipped {
		if !jsonOut {
			writeStdout("⏭️  Skipping EARS for %s (status: %s)\n\n", epic, status)
		}
		return true
	}

	reqPath := filepath.Join(epicsPath, epic, "ep-requirements.md")
	result := lintEARSFile(reqPath)
	result.Epic = epic
	if jsonOut {
		data, marshalErr := json.MarshalIndent(result, "", "  ")
		if marshalErr != nil {
			errLog.Printf("Error marshaling JSON: %v\n", marshalErr)
			return false
		}
		writelnStdout(string(data))
	} else {
		printEARSHuman(result)
	}
	return !result.HasGaps
}

func validateSingleEpicREQ(epic string, jsonOut bool) bool {
	epicsPath, err := epicArtefactsRoot()
	if err != nil {
		errLog.Printf("Error: %v\n", err)
		return false
	}
	skipped, status, err := epicSkippedForArtefactGate(epicsPath, epic)
	if err != nil {
		errLog.Printf("Error validating %s: %v\n", epic, err)
		return false
	}
	if skipped {
		if !jsonOut {
			writeStdout("⏭️  Skipping REQ↔AC for %s (status: %s)\n\n", epic, status)
		}
		return true
	}

	result, err := validateREQForEpic(epicsPath, epic)
	if err != nil {
		errLog.Printf("Error validating %s: %v\n", epic, err)
		return false
	}
	if jsonOut {
		data, marshalErr := json.MarshalIndent(result, "", "  ")
		if marshalErr != nil {
			errLog.Printf("Error marshaling JSON: %v\n", marshalErr)
			return false
		}
		writelnStdout(string(data))
	} else {
		printREQACTraceHuman(result)
	}
	return !result.HasGaps
}
