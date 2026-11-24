package cli

func init() {
	for _, code := range []Code{"waiver-not-in-review", "waiver-handoff-required", "waiver-not-lower", "waiver-hard-floor"} {
		failureCodeDefinitions[code] = definition(CategoryPrecondition, false, nil, stringsField("paths"))
		registeredCodes = append(registeredCodes, code)
	}
}
