package data

var EventInfoMapping = map[string]eventInfo{
	"MEI":                {E_MEI_1, "Main engines ignition"},
	"LIFT_OFF":           {E_LIFTOFF, "Lift off"},
	"PITCH":              {E_STARTPITCH, "Pitch starts"},
	"MECO":               {E_MECO_1, "Main engine cut off"},
	"STAGE_SEP":          {E_STAGESEP, "Stage Separation"},
	"SEI":                {E_SEI_1, "Second stage ignition"},
	"BOOSTBACK_BURN_ON":  {E_BBURNI, "Boostback burn started"},
	"BOOSTBACK_BURN_OFF": {E_BBURNO, "Boostback burn stopped"},
	"ENTRY_BURN_ON":      {E_EBURNI, "Entry burn started"},
	"ENTRY_BURN_OFF":     {E_EBURNO, "Entry burn stopped"},
	"SECO":               {E_SECO_1, "Second stage engine cut off"},
	"SEI2":               {E_SEI_2, "Second stage 2nd engine ignition"},
	"SECO2":              {E_SECO_2, "Second stage 2nd engine cut off"},
	"LANDING_BURN_ON":    {E_LBURNI, "Landing burn started"},
	"LANDING_BURN_OFF":   {E_LBURNO, "Landing burn stopped"},
	"THROTTLE_DWN":       {E_THROTTLE_D, "Throttling down"},
	"THROTTLE_UP":        {E_THROTTLE_U, "Throttling up"},
}

var EventMappingString = map[string]string{
	"MEI":                "Main engines ignition",
	"LIFT_OFF":           "Lift off",
	"PITCH":              "Pitch starts",
	"MECO":               "Main engine cut off",
	"STAGE_SEP":          "Stage Separation",
	"SEI":                "Second stage ignition",
	"BOOSTBACK_BURN_ON":  "Boostback burn started",
	"BOOSTBACK_BURN_OFF": "Boostback burn stopped",
	"ENTRY_BURN_ON":      "Entry burn started",
	"ENTRY_BURN_OFF":     "Entry burn stopped",
	"SECO":               "Second stage engine cut off",
	"SEI2":               "Second stage 2nd engine ignition",
	"SECO2":              "Second stage 2nd engine cut off",
	"LANDING_BURN_ON":    "Landing burn started",
	"LANDING_BURN_OFF":   "Landing burn stopped",
	"THROTTLE_DWN":       "Throttling down",
	"THROTTLE_UP":        "Throttling up",
}

var EventMapping = map[string]uint32{
	"MEI":                E_MEI_1,
	"LIFT_OFF":           E_LIFTOFF,
	"PITCH":              E_STARTPITCH,
	"MECO":               E_MECO_1,
	"STAGE_SEP":          E_STAGESEP,
	"SEI":                E_SEI_1,
	"BOOSTBACK_BURN_ON":  E_BBURNI,
	"BOOSTBACK_BURN_OFF": E_BBURNO,
	"ENTRY_BURN_ON":      E_EBURNI,
	"ENTRY_BURN_OFF":     E_EBURNO,
	"SECO":               E_SECO_1,
	"SEI2":               E_SEI_2,
	"SECO2":              E_SECO_2,
	"LANDING_BURN_ON":    E_LBURNI,
	"LANDING_BURN_OFF":   E_LBURNO,
	"THROTTLE_DWN":       E_THROTTLE_D,
	"THROTTLE_UP":        E_THROTTLE_U,
}
