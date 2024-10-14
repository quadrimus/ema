package tool

import (
	"bytes"
	"sort"
)

// Help returns text with manual for user.
func Help() string {
	list := registeredFormats()

	sort.Slice(list, func(i, j int) bool {
		if list[i].input == list[j].input {
			return list[i].output < list[j].output
		}
		return list[i].input < list[j].input
	})

	b := &bytes.Buffer{}
	b.WriteString("Usage: ema <input> <output>\n\n")
	b.WriteString("Arguments:\n")
	b.WriteString("  path/to/file.ext           to use file\n")
	b.WriteString("  path/to/file.ext#<format>  to use file as format\n")
	b.WriteString("  -#<format>                 to use stdin/stdout as format\n")
	b.WriteByte('\n')
	b.WriteString("Formats:\n")
	for _, v := range list {
		b.WriteString("  ")
		b.WriteString(v.input)
		b.WriteString(" -> ")
		b.WriteString(v.output)
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
	b.WriteString("Examples:\n")
	b.WriteString("  ema file.ema file.json\n")
	b.WriteString("  ema file.ext#ema file.json\n")
	b.WriteString("  ema file.ema -#json\n")
	b.WriteString("  ema -#ema -#json\n")
	b.WriteByte('\n')
	return b.String()
}

type formats struct {
	input  string
	output string
}

func registeredFormats() []formats {
	registeredWorkersMutex.Lock()
	defer registeredWorkersMutex.Unlock()

	list := make([]formats, 0, len(registeredWorkers))
	for key := range registeredWorkers {
		list = append(list, formats{
			input:  key[0],
			output: key[1],
		})
	}
	return list
}
