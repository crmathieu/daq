// events
const E_MEI_1 = 0x00000001;	// main engines ignition
const E_LIFTOFF = 0x00000002;	// lift off
const E_STARTPITCH = 0x00000004;	// pitch maneuver starts
const E_MECO_1 = 0x00000008;	// main engines cutoff 1
const E_STAGESEP = 0x00000010;	// stage separation
const E_SEI_1 = 0x00000020;	// second stage engine ignition
const E_MEI_2 = 0x00000040;	// main engines restart 1 	(boost backburn [land touchdown] or entry burn [sea touchdown])
const E_MECO_2 = 0x00000080;	// main engines cutoff 2	
const E_MEI_3 = 0x00000100;	// main engines ignition 3	(entry burn [land touchdown] or landing burn [sea touchdown])
const E_MECO_3 = 0x00000200;	// main engines cutoff 3
const E_SECO_1 = 0x00000400;	// second stage engine cutoff 1
const E_SEI_2 = 0x00001000;	// second stage restart 1  (orbit boost)
const E_SECO_2 = 0x00002000;	// second stage engine cutoff 2
const E_MEI_4 = 0x00004000;	// main engines ignition 4 (landing burn [land touchdown])
const E_MECO_4 = 0x00008000;	// main engines cutoff 4 


var eventMap = new Map([
    [E_MEI_1, "main engines ignition"],
    [E_LIFTOFF, "lift off"],
    [E_STARTPITCH, "pitch maneuver starts"],
    [E_MECO_1, "main engines cutoff 1"],
    [E_STAGESEP, "stage separation"],
    [E_SEI_1, "second stage engine ignition"],
    [E_MEI_2, "main engines restart 1"],
    [E_MECO_2, "main engines cutoff 2"],
    [E_MEI_3, "main engines ignition 3"],
    [E_MECO_3, "main engines cutoff 3"],
    [E_SECO_1, "second stage engine cutoff 1"],
    [E_SEI_2, "second stage restart 1"],
    [E_SECO_2, "second stage engine cutoff 2"],
    [E_MEI_4, "main engines ignition 4"],
    [E_MECO_4, "main engines cutoff 4"]

]);

// datapoint types
const IDVELOCITY = 1;
const IDPOSITION = 2;
const IDTILTANGLE = 3;
const IDTHRUST = 4;
const IDEVENT = 5;
const IDMASSPROPELLANT = 6;

var lastEvent = 0;