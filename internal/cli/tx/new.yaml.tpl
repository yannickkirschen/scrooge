# --- Create a new transaction ---
# Fill out the form, save and quit the editor.

# Date in ISO-8601 format (e.g. 2000-01-01).
date: "{{`{{ now }}`}}"

# {{ .AccountRefs }}
account: ""

type: "" # receipt, spending or balance

tags: # {{ .Tags }}
#  -

status: "booked" # booked or planned
description: ""
amount: "0.00"
currency: "EUR"
