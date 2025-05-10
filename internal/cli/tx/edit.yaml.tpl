# --- Edit transaction {{ .Tx.Id }} ---
# Fill out the form, save and quit the editor.

# Date in ISO-8601 format (e.g. 2000-01-01).
date: "{{ .Tx.Date }}"

# {{ .AccountRefs }}
account: {{ .Tx.Account.Id }}

type: {{ .Tx.Type }} # receipt, spending or balance

tags: # {{ .Tags }}
{{- range $tag := .Tx.Tags }}
    - {{ $tag }}
{{- end }}

status: {{ .Tx.Status }} # booked or planned
description: {{ .Tx.Description }}
amount: {{ .Tx.Amount }}
currency: {{ .Tx.Currency }}
