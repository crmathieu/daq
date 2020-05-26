package data

const (
	E_MEI_1 	= 0x00000001	// main engines ignition
	E_LIFTOFF 	= 0x00000002	// lift off
	E_STARTPITCH= 0x00000004	// pitch maneuver starts
	E_MECO_1	= 0x00000008	// main engines cutoff 1
	E_STAGESEP	= 0x00000010	// stage separation
	E_SEI_1		= 0x00000020	// second stage engine ignition
	E_BBURNI	= 0x00000040 	// first stage Boost back burn starts
	E_BBURNCO	= 0x00000080 	// first stage Boost back burn starts
	E_MEI_2		= 0x00000040	// main engines restart 1 	(boost backburn [land touchdown] or entry burn [sea touchdown])
	E_MECO_2	= 0x00000080	// main engines cutoff 2	
	E_EBURNI	= 0x00000100 	// first stage entry back burn starts
	E_EBURNCO	= 0x00000200 	// first stage entry back burn starts
	E_MEI_3		= 0x00000100	// main engines ignition 3	(entry burn [land touchdown] or landing burn [sea touchdown])
	E_MECO_3	= 0x00000200	// main engines cutoff 3
	E_SECO_1	= 0x00000400	// second stage engine cutoff 1
	E_SEI_2		= 0x00001000	// second stage restart 1  (orbit boost)
	E_SECO_2	= 0x00002000	// second stage engine cutoff 2
	E_LBURNI	= 0x00004000	// first stage landing burn starts
	E_LBURNCO	= 0x00008000	// first stage landing burn stops
	E_MEI_4		= 0x00004000	// main engines ignition 4 (landing burn [land touchdown])
	E_MECO_4	= 0x00008000	// main engines cutoff 4 
)

// note:
// 	E_BBSTART and E_BBSTOP can be used interchangeably with E_MEI_2 and E_MECO_2
//  E_EBSTART and E_EBSTOP can be used interchangeably with E_MEI_3 and E_MECO_3
//  E_LBSTART and E_LBSTOP can be used interchangeably with E_MEI_4 and E_MECO_4

