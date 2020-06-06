package data

const (
	E_MEI_1 	= 0x00000001	// main engines ignition
	E_LIFTOFF 	= 0x00000002	// lift off
	E_STARTPITCH= 0x00000004	// pitch maneuver starts
	E_MECO_1	= 0x00000008	// main engines cutoff 1
	E_STAGESEP	= 0x00000010	// stage separation
	E_SEI_1		= 0x00000020	// second stage engine ignition
	E_BBURNI	= 0x00000040 	// first stage Boost back burn starts
	E_BBURNO	= 0x00000080 	// first stage Boost back burn stops
	E_EBURNI	= 0x00000100 	// first stage entry back burn starts
	E_EBURNO	= 0x00000200 	// first stage entry back burn starts
	E_SECO_1	= 0x00000400	// second stage engine cutoff 1
	E_SEI_2		= 0x00000800	// second stage restart 1  (orbit boost)
	E_SECO_2	= 0x00001000	// second stage engine cutoff 2
	E_LBURNI	= 0x00002000	// first stage landing burn starts
	E_LBURNO	= 0x00004000	// first stage landing burn stops
	E_THROTTLE_D  = 0x00008000 	// throttle engine down
	E_THROTTLE_U  = 0x00010000 	// throttle engine down
)

var EventMapping = map[string]uint32{
	"MEI": E_MEI_1,
	"LIFT_OFF": 	E_LIFTOFF,
	"PITCH": E_STARTPITCH,
	"MECO": E_MECO_1,
	"STAGE_SEP": E_STAGESEP,
	"SEI": E_SEI_1,
	"BOOSTBACK_BURN_ON": E_BBURNI,
	"BOOSTBACK_BURN_OFF": E_BBURNO,
	"ENTRY_BURN_ON": E_EBURNI,
	"ENTRY_BURN_OFF": E_EBURNO,
	"SECO": E_SECO_1,
	"SEI2": E_SEI_2,
	"SECO2":E_SECO_2,
	"LANDING_BURN_ON": E_LBURNI,
	"LANDING_BURN_OFF":E_LBURNO,
	"THROTTLE_DWN": E_THROTTLE_D,
	"THROTTLE_UP": E_THROTTLE_U,
} 

var EventMappingString = map[string]string{
	"MEI": "Main engines ignition",
	"LIFT_OFF": "Lift off",
	"PITCH": "Pitch starts",
	"MECO": "Main engine cut off",
	"STAGE_SEP": "Stage Separation",
	"SEI": "Second stage ignition",
	"BOOSTBACK_BURN_ON": "Boostback burn started",
	"BOOSTBACK_BURN_OFF": "Boostback burn stopped",
	"ENTRY_BURN_ON": "Entry burn started",
	"ENTRY_BURN_OFF": "ENtry burn stopped",
	"SECO": "Second stage engine cut off",
	"SEI2": "Second stage 2nd engine ignition",
	"SECO2": "Second stage 2nd engine cut off",
	"LANDING_BURN_ON": "Landing burn started",
	"LANDING_BURN_OFF": "Landing burn stopped",
	"THROTTLE_DWN": "Throttling down",
	"THROTTLE_UP": "Throttling up",
} 

// note:
// 	E_BBSTART and E_BBSTOP can be used interchangeably with E_MEI_2 and E_MECO_2
//  E_EBSTART and E_EBSTOP can be used interchangeably with E_MEI_3 and E_MECO_3
//  E_LBSTART and E_LBSTOP can be used interchangeably with E_MEI_4 and E_MECO_4

