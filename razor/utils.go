package razor

// original from terraform-provider-vault
func toStringArray(input []interface{}) []string {
	output := make([]string, len(input))

	for i, item := range input {
		output[i] = item.(string)
	}

	return output
}
