#!/bin/bash

# Calls golint and filters out all the false-positives.

set -eu

uncommented_functions=(
	# Obvious names
	"New\w*"
)

uncommented_types=(
	# Obvious names
	"\w*Options"
)

uncommented_methods=(
	# Obvious names
	"\w*\.Reset"
	"\w*\.Init"

	# Interface implementations
	"\w*\.On\w*"
	"\w*\.BuildFrom\w*"
	"\w*\.BuildBegin\w*"
	"\w*\.BuildEndContainer"
	"\w*\.PrepareFor\w*Contents"
	"\w*\.NotifyChildContainerFinished"
	"RootBuilder.IsContainerOnly"
	"RootBuilder.InitTemplate"
	"RootBuilder.NewInstance"
)

general_ignored=(
	# Umm... no.
	"don't use an underscore in package name"

	# No, it shouldn't.
	"comment on exported .* should be of the form"

	# Yes, it is used. The linter is buggy.
	"imported but not used"
)

cmd="golint"

for ((i = 0; i < ${#general_ignored[@]}; i++)); do
	cmd+=" | grep -v \"${general_ignored[$i]}\""
done

for ((i = 0; i < ${#uncommented_functions[@]}; i++)); do
	cmd+=" | grep -v \"exported function ${uncommented_functions[$i]} should have comment or be unexported\""
done

for ((i = 0; i < ${#uncommented_types[@]}; i++)); do
	cmd+=" | grep -v \"exported type ${uncommented_types[$i]} should have comment or be unexported\""
done

for ((i = 0; i < ${#uncommented_methods[@]}; i++)); do
	cmd+=" | grep -v \"exported method ${uncommented_methods[$i]} should have comment or be unexported\""
done

eval "$cmd"
