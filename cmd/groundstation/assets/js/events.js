// events
const E_MEI_1 = 0x00000001;	// main engines ignition
const E_LIFTOFF = 0x00000002;	// lift off
const E_STARTPITCH = 0x00000004;	// pitch maneuver starts
const E_MECO_1 = 0x00000008;	// main engines cutoff 1
const E_STAGESEP = 0x00000010;	// stage separation
const E_SEI_1 = 0x00000020;	// second stage engine ignition
const E_BBURNI = 0x00000040; 	// first stage Boost back burn starts
const E_BBURNO = 0x00000080; 	// first stage Boost back burn stops
const E_EBURNI = 0x00000100; 	// first stage entry back burn starts
const E_EBURNO = 0x00000200; 	// first stage entry back burn starts
const E_SECO_1 = 0x00000400;	// second stage engine cutoff 1
const E_SEI_2 = 0x00000800;	// second stage restart 1  (orbit boost)
const E_SECO_2 = 0x00001000;	// second stage engine cutoff 2
const E_LBURNI = 0x00002000;	// first stage landing burn starts
const E_LBURNO = 0x00004000;	// first stage landing burn stops
const E_THROTTLE_D = 0x00008000; 	// throttle engine down
const E_THROTTLE_U = 0x00010000; 	// throttle engine up

var eventMap = new Map([
    [E_MEI_1, "main engines ignition"],
    [E_LIFTOFF, "lift off"],
    [E_STARTPITCH, "pitch maneuver starts"],
    [E_MECO_1, "main engines cutoff"],
    [E_STAGESEP, "stage separation"],
    [E_SEI_1, "second stage engine ignition"],
    [E_BBURNI, "first stage Boost back burn starts"],
    [E_BBURNO, "first stage Boost back burn stops"],
    [E_EBURNI, "first stage entry burn starts"],
    [E_EBURNO, "first stage entry burn stops"],
    [E_SECO_1, "second stage engine cutoff 1"],
    [E_SEI_2, "second stage restart 1"],
    [E_SECO_2, "second stage engine cutoff 2"],
    [E_LBURNI, "first stage landing burn starts"],
    [E_LBURNO, "first stage landing burn stops"],
    [E_THROTTLE_U, "Throttling up"],
    [E_THROTTLE_D, "Throttling down"]
]);

// datapoint types
const IDVELOCITY = 1;
const IDPOSITION = 2;
const IDTILTANGLE = 3;
const IDTHRUST = 4;
const IDEVENT = 5;
const IDMASSPROPELLANT = 6;
const IDTIME = 7;

var lastEvent = 0;

function outputDate(hh, mm, ss) {
    var htime = hh.toString();
    if (hh <= 9) htime = "0"+htime;
    var mtime = mm.toString();
    if (mm <= 9) mtime = "0" + mtime;
    var stime = ss.toString();
    if (ss <= 9) stime = "0" + stime;
    return htime+":"+mtime+":"+stime;
}