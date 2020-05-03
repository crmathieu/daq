//
//  Created by Charles Mathieu on 11/16/15.
//  Copyright © 2015 Charles.Mathieu. All rights reserved.
//
package main

import (
    "fmt"
    "strconv"
    "math"
    "strings"
)

const (
    R_SGC = 287.05               // specific gas constant in J / kg.K
    GAMMA_AIR = 1.4
    EARTHRADIUS = float64(6356766.0)
    EARTHMASS = float64(5.972e24)
    UNIVERSALG = 6.67384e-11
    GRAVITY_ACC = 9.80665       // accel. of GRAVITY_ACC
    MOL_WT  = 28.9644           // molecular weight of air
    R_GAS   = 8.31432           // gas constant

    // standard sealevel
    STD_SEALEVEL_TEMP = 288.15  // in Kelvin
    STD_SEALEVEL_DENSITY = 1.2250
    STD_SEALEVEL_PRESSURE = 101325.0
)

type layer struct {
    name           string       
    base_altitude  float64      // in km
    top_altitude   float64      // in km
    base_temp      float64      // in kelvin
    base_pressure  float64      // in pa
    pressRatio     float64      // in Bars
    lapseRate      float64      // slope in linear temp variation for this layer
}

type isa76 struct {
    layers              []*layer
    params              AtmoParams

    seaLevelTemp        float64
    layersInitialized   bool
    layersDisplay       string
    custAltitude        string
    lastAltitude        float64
    lastUnit            string
    
    // hydrostatic constant
    GMR                 float64 
}

type AtmoParams struct {
    layerName               string
    geopotentialAltMeters   float64
    geopotentialAltFeet     float64
    geometricAltMeters      float64
    geometricAltFeet        float64
    temperatureK            float64
    temperatureC            float64
    temperatureF            float64
    temperatureR            float64
    pressure                float64
    density                 float64
    soundSpeed              float64
    pressureRatio           float64
    tempRatio               float64
    densityRatio            float64
}        


var stdPressure = make(map[float64]float64)
var custPressure = make(map[float64]float64)

func (isa *isa76) set() {
    isa.seaLevelTemp = 0.0
    isa.layersInitialized = false
    isa.layersDisplay = ""
    isa.custAltitude = ""
    isa.lastAltitude = 0.0001
    isa.lastUnit = ""

    stdPressure[-5.004] = 1.7543
    stdPressure[0.0] = 1.0
    stdPressure[11.0] = 2.233611E-1
    stdPressure[20.0] = 5.403295E-2
    stdPressure[32.0] = 8.5666784E-3
    stdPressure[47.0] = 1.0945601E-3
    stdPressure[51.0] = 6.6063531E-4
    stdPressure[71.0] = 3.9046834E-5
    stdPressure[84.852] = 3.68501E-6


    isa.GMR = GRAVITY_ACC * (MOL_WT / R_GAS) // hydrostatic constant
        
    isa.seaLevelTemp = STD_SEALEVEL_TEMP
    
    isa.layers = append(isa.layers, &layer{name: "Negative Troposphere", base_altitude:-5.004, top_altitude:0.0,     base_pressure: STD_SEALEVEL_PRESSURE * 1.7543,         base_temp: 320.676,              pressRatio: 1.7543,       lapseRate:-6.5 })
    isa.layers = append(isa.layers, &layer{name: "Troposphere",          base_altitude:0.0,    top_altitude:11.0,    base_pressure: STD_SEALEVEL_PRESSURE,                  base_temp: STD_SEALEVEL_TEMP,    pressRatio: 1.0,          lapseRate:-6.5 })
    isa.layers = append(isa.layers, &layer{name: "Tropopause",           base_altitude:11.0,   top_altitude:20.0,    base_pressure: STD_SEALEVEL_PRESSURE * 2.233611E-1,    base_temp: 216.65,               pressRatio: 2.233611E-1,  lapseRate: 0.0 })
    isa.layers = append(isa.layers, &layer{name: "Low Stratosphere",     base_altitude:20.0,   top_altitude:32.0,    base_pressure: STD_SEALEVEL_PRESSURE * 5.403295E-2,    base_temp: 216.65,               pressRatio: 5.403295E-2,  lapseRate: 1.0 })
    isa.layers = append(isa.layers, &layer{name: "High Stratosphere",    base_altitude:32.0,   top_altitude:47.0,    base_pressure: STD_SEALEVEL_PRESSURE * 8.5666784E-3,   base_temp: 228.65,               pressRatio: 8.5666784E-3, lapseRate: 2.8 })
    isa.layers = append(isa.layers, &layer{name: "Stratopause",          base_altitude:47.0,   top_altitude:51.0,    base_pressure: STD_SEALEVEL_PRESSURE * 1.0945601E-3,   base_temp: 270.65,               pressRatio: 1.0945601E-3, lapseRate: 0.0 })
    isa.layers = append(isa.layers, &layer{name: "Low Mesophere",        base_altitude:51.0,   top_altitude:71.0,    base_pressure: STD_SEALEVEL_PRESSURE * 6.6063531E-4,   base_temp: 270.65,               pressRatio: 6.6063531E-4, lapseRate:-2.8 })
    isa.layers = append(isa.layers, &layer{name: "High Mesophere",       base_altitude:71.0,   top_altitude:84.852,  base_pressure: STD_SEALEVEL_PRESSURE * 3.9046834E-5,   base_temp: 214.65,               pressRatio: 3.9046834E-5, lapseRate:-2.0 })
    isa.layers = append(isa.layers, &layer{name: "Mesopause",            base_altitude:84.852, top_altitude:86.0,    base_pressure: STD_SEALEVEL_PRESSURE * 3.68501E-6,     base_temp: 186.946,              pressRatio: 3.68501E-6,   lapseRate: 0.0 })
}


func (isa *isa76) setSealevelTemp(temp float64) {
    
    // temperature must be in kelvin
    if temp < 0.0 {
        temp = 0.0
    }
    isa.seaLevelTemp = temp
    localTemp := isa.seaLevelTemp

    // recalculate temps and pressure ratio for each layer
    for k, layer := range isa.layers {
        delta_h := layer.top_altitude - layer.base_altitude
        if k == 0 {
           layer.base_temp = localTemp - (layer.lapseRate * delta_h)
        } else {
            layer.base_temp = localTemp
            localTemp = layer.base_temp + (layer.lapseRate * delta_h) 
        }
        fmt.Println("---------------------")
        fmt.Println(layer.base_temp)
        fmt.Println("---------------------")
        custPressure[layer.base_altitude] = isa.getRelativePressureRatio(layer, k, localTemp, delta_h)
    }
    for _, layer := range isa.layers {
        fmt.Printf("%#v\n\n", *layer)
    }
}

func (isa *isa76) swap() {
    for _, layer := range isa.layers {
        layer.pressRatio = custPressure[layer.base_altitude]
    }
}

func (isa *isa76) getRelativePressureRatio(l *layer, index int, localTemp, delta_h float64) float64 {
    switch index {
    case 0: return 1
    case 1: return 1
    default:
        if l.lapseRate == 0.0 {
            // isothermal
            return math.Exp(-isa.GMR*(delta_h/l.base_temp))   // -> aka delta
        } else {
            return math.Pow(l.base_temp/localTemp, isa.GMR/l.lapseRate) // -> aka delta
        }
    }
}

func (isa *isa76) customAltitude(alt, unit string) string {
    const METRICS_MAX_ALTITUDE = 86000
    const METRICS_MIN_ALTITUDE = -5000

    altitude, _ := strconv.ParseFloat(alt, 64)
    mulfactor := 1.0
    cmul := 3.2808
    m1 := "m"
    m2 := "meters"
    cm1 := "ft"
    cm2 := "feet"
    if unit == "f" {
        // convert altitude in metrics
        altitude = altitude / 3.2808
        mulfactor = 3.2808
        cmul = 1/3.2808
        m1 = "ft"
        m2 = "feet"
        cm1 = "m"
        cm2 = "meters"
    }

    // calculate geopotential altitude and perform sanity check
    var geopoAltitude = altitude * (EARTHRADIUS/(EARTHRADIUS+altitude))
    if geopoAltitude > METRICS_MAX_ALTITUDE {
        altitude = METRICS_MAX_ALTITUDE
        geopoAltitude = altitude * (EARTHRADIUS/(EARTHRADIUS+altitude))
    } else if geopoAltitude < (METRICS_MIN_ALTITUDE - 4) {
        altitude = METRICS_MIN_ALTITUDE
        geopoAltitude = altitude * (EARTHRADIUS/(EARTHRADIUS+altitude))
    }
    
    // convert in km
    altitude = altitude / 1000
    geopoAltitude = geopoAltitude / 1000
    
    if geopoAltitude != isa.lastAltitude {
        isa.lastAltitude = geopoAltitude
        isa.getParameters(altitude, geopoAltitude)
    }
    localValue := formatTitle("Geometric Altitude")
    postProcessed := normalizeNumber("%7.4f", altitude*1000*mulfactor) //math.Round(isa.params.temperatureK*1000)/1000)
    localValue = localValue + addLine(postProcessed, m1, m2)

    postProcessed = normalizeNumber("%7.4f", altitude*1000*cmul) //math.Round(isa.params.temperatureK*1000)/1000)
    localValue = localValue + addLine(postProcessed, cm1, cm2)

    localValue = localValue + "\n" + formatTitle("Geopotential Altitude")

    postProcessed = normalizeNumber("%7.4f", geopoAltitude*1000*mulfactor) //math.Round(isa.params.temperatureK*1000)/1000)
    localValue = localValue + addLine(postProcessed, m1, m2)

    postProcessed = normalizeNumber("%7.4f", geopoAltitude*1000*cmul) //math.Round(isa.params.temperatureK*1000)/1000)
    localValue = localValue + addLine(postProcessed, cm1, cm2)
    return localValue
}
    
//func (isa *isa76) getParameters(geometricAltitude float64, geoPotentialAlt geopoAltitude:Double) {
func (isa *isa76) getParameters(geometricAltitude float64, geopoAltitude float64) {
    
    for _, layer := range isa.layers {
        
        if geopoAltitude >= layer.base_altitude && geopoAltitude <= layer.top_altitude {

            delta_h := geopoAltitude - layer.base_altitude
            localTemp := layer.base_temp + (layer.lapseRate * delta_h)
            
            isa.params.tempRatio = localTemp / isa.seaLevelTemp // -> aka theta
            isa.params.geopotentialAltMeters = geopoAltitude * 1000
            isa.params.geopotentialAltFeet = isa.params.geopotentialAltMeters * 3.28084
            
            isa.params.geometricAltMeters = geometricAltitude * 1000
            isa.params.geometricAltFeet = isa.params.geometricAltMeters * 3.28084
            
            if layer.lapseRate == 0.0 {
                // isothermal
                isa.params.pressureRatio = layer.pressRatio * math.Exp(-isa.GMR*(delta_h/layer.base_temp))   // -> aka delta
            } else {
                isa.params.pressureRatio = layer.pressRatio * math.Pow(layer.base_temp/localTemp, isa.GMR/layer.lapseRate) // -> aka delta
            }
            
            //isa.params.densityRatio = isa.params.pressureRatio / isa.params.tempRatio // -> aka sigma

            isa.params.layerName = layer.name
            isa.params.densityRatio = isa.params.pressureRatio / isa.params.tempRatio   // -> aka sigma
            isa.params.temperatureK = isa.params.tempRatio * isa.seaLevelTemp
            isa.params.temperatureC = isa.params.temperatureK - 273.15
            isa.params.temperatureF = ((isa.params.temperatureK - 273.15) * 1.8) + 32
            isa.params.temperatureR = (isa.params.temperatureK) * 1.8
            
            isa.params.pressure     = isa.params.pressureRatio * STD_SEALEVEL_PRESSURE
            isa.params.density      = isa.params.densityRatio * STD_SEALEVEL_DENSITY
            isa.params.soundSpeed   = math.Pow((GAMMA_AIR * isa.params.temperatureK * R_SGC), 0.5)
            break
        }
    }
}

func (isa *isa76) showEarthAcceleration() string {
        
        //         Re^2         where Re is the radius of the earth
        // g = -------------    hg is the geometric altitude above sea level
        //     go (Re + hg)2
        localValue := formatTitle("Acceleration")
        glocal := GRAVITY_ACC * math.Pow(EARTHRADIUS, 2) / math.Pow(EARTHRADIUS + isa.params.geometricAltMeters, 2)
        
        //postProcessed := output.normalizeNumber(String(format:"%2.6f",round(glocal*1000000)/1000000))
        //return postProcessed
        return localValue + addLine(normalizeNumber("%7.6f",math.Round(glocal*1000000)/1000000), "m/s²", "Accel. of gravity")

}
 
/*
    private func formatTitle(title:String, abrvFrom: String) -> String {
    
        let additional = "<font size=-2>(Standard sea level)</font>"

        if abrvFrom == "f" {
            return output.startTable(additional, title: title+" at<br>"+output.normalizeNumber("\(round(params.geometricAltFeet))")+" ft ("+output.normalizeNumber("\(round(params.geometricAltMeters))")+" m)<br>")
        } else {
            return output.startTable(additional, title: title+" at<br>"+output.normalizeNumber("\(round(params.geometricAltMeters))")+" m ("+output.normalizeNumber("\(round(params.geometricAltFeet))")+" ft)<br>")
        }

    }
*/    
func formatTitle(label string) string {
    return fmt.Sprintf("** %v **\n", label)
}

func addLine(value, unit, label string) string {
    return fmt.Sprintf("   %v %v (%v)\n", value, unit, label)
}

func addHeadline(layerName, unit, label string) string {
    return fmt.Sprintf("   %v %v (%v)\n", layerName, unit, label)
}

func normalizeNumber(formatstr string, value float64) string {
    return fmt.Sprintf(formatstr, value)
}

func endTable() string {
    return fmt.Sprintf("\n")
}

func (isa *isa76) showTemperature(abrvFrom string) string {
    var postProcessed = ""
    var localValue = ""
    localValue = formatTitle("Temperature")
    
//    postProcessed = normalizeNumber(fmt.String(format:"%4.2f",round(isa.params.temperatureK*1000)/1000))
    postProcessed = normalizeNumber("%4.2f", math.Round(isa.params.temperatureK*1000)/1000)
    localValue = localValue + addLine(postProcessed, "K", "Kelvin")
    
    postProcessed = normalizeNumber("%4.2f", math.Round(isa.params.temperatureC*1000)/1000)
    localValue = localValue + addLine(postProcessed, "C", "Celcius")
    
    postProcessed = normalizeNumber("%4.2f", math.Round(isa.params.temperatureF*1000)/1000)
    localValue = localValue + addLine(postProcessed, "F", "Fahrenheit")
    
    postProcessed = normalizeNumber("%4.2f", math.Round(isa.params.temperatureR*1000)/1000)
    localValue = localValue + addLine(postProcessed, "R", "Rankine")
    
    postProcessed = normalizeNumber("%4.2f", math.Round(isa.params.tempRatio*1000)/1000)
    localValue = localValue + addLine(postProcessed, "", "[θ] (Ratio t/t0)")
    //localValue += endTable()
    
    return localValue
}
  
    
func (isa *isa76) showPressure(abrvFrom string) string {
    var postProcessed = ""
    var localValue = ""
    
    localValue = formatTitle("Pressure")

    
    postProcessed = normalizeNumber("%7.6f", math.Round(isa.params.pressure*100000)/100000)
    localValue = localValue + addLine(postProcessed, "pa", "pascals")

    postProcessed = normalizeNumber("%7.6f", math.Round((isa.params.pressure * 1E-03)*1E+9)/1E+9)
    localValue = localValue + addLine(postProcessed, "kpa", "kilopascals")

    postProcessed = normalizeNumber("%7.6f", math.Round((isa.params.pressure * 1E-05)*1E+9)/1E+9)
    localValue = localValue + addLine(postProcessed, "bar", "bars")
    
    postProcessed = normalizeNumber("%7.6f", math.Round((isa.params.pressure * 1E-02)*1E+6)/1E+6)
    localValue = localValue + addLine(postProcessed, "mbar", "millibars")
    
    postProcessed = normalizeNumber("%7.6f", math.Round((isa.params.pressure/6.8947572798677E+03)*1E+5)/1E+5)
    localValue = localValue + addLine(postProcessed, "lb/in²", "pound/inch²")
    
    postProcessed = normalizeNumber("%7.6f", math.Round((isa.params.pressure/0.47880258888E+02)*1E+5)/1E+5)
    localValue = localValue + addLine(postProcessed, "lb/ft²", "pound/foot²")
    
    postProcessed = normalizeNumber("%7.6f", math.Round((isa.params.pressure/1.013250E+05)*1E+9)/1E+9)
    localValue = localValue + addLine(postProcessed, "atm", "atmosphere")
    
    postProcessed = normalizeNumber("%7.6f", math.Round((isa.params.pressure/1.3332239E+02)*1E+5)/1E+5)
    localValue = localValue + addLine(postProcessed, "mmhg", "mm of mercury")
    
    postProcessed = normalizeNumber("%7.6f", math.Round((isa.params.pressureRatio)*1E+8)/1E+8)
    localValue = localValue + addLine(postProcessed, "", "[δ] (Ratio p/p0)")
    
    //localValue = localValue + endTable()
    
    return localValue
}
   

func (isa *isa76) showSoundSpeed(abrvFrom string) string {
    var postProcessed = ""
    var localValue = ""
    
    localValue = formatTitle("Speed of Sound")

    
    metricValue := math.Pow((GAMMA_AIR * isa.params.temperatureK * R_SGC), 0.5)
    
    postProcessed = normalizeNumber("%4.3f", math.Round(metricValue*1000)/1000)
    localValue = localValue + addLine(postProcessed, "m/s", "meters/second")

    postProcessed = normalizeNumber("%4.3f", math.Round((metricValue/340.37713655487805)*1000)/1000)
    localValue = localValue + addLine(postProcessed, "mach", "mach number")
    
    postProcessed = normalizeNumber("%4.3f", math.Round((metricValue*3.281)*1000)/1000)
    localValue = localValue + addLine(postProcessed, "ft/s", "feet/second")
    
    postProcessed = normalizeNumber("%4.3f", math.Round((metricValue*3.6)*1000)/1000)
    localValue = localValue + addLine(postProcessed, "km/h", "kilometers/hour")

    postProcessed = normalizeNumber("%4.3f", math.Round((metricValue*(0.000621371*3600))*1000)/1000)
    localValue = localValue + addLine(postProcessed, "mi/h", "miles/hour")
    
    postProcessed = normalizeNumber("%4.3f", math.Round((metricValue*1.94384)*1000)/1000)
    localValue = localValue + addLine(postProcessed, "kts", "knots")
    //localValue = localValue + endTable()

    return localValue
}
    
    
func (isa *isa76) showDensity(abrvFrom string) string {
    var postProcessed = ""
    var localValue = ""
    
    localValue = formatTitle("Density")

    
    postProcessed = normalizeNumber("%3.9f", math.Round(isa.params.density*1000000000)/1000000000)
    localValue = localValue + addLine(postProcessed, "kg/m³", "kilograms/meter³")
    
    postProcessed = normalizeNumber("%3.9f", math.Round((isa.params.density * 1.94024058983E-03)*1000000000)/1000000000)
    localValue = localValue + addLine(postProcessed, "slug/ft³", "slug/foot³")
    
    postProcessed = normalizeNumber("%3.9f", math.Round((isa.params.density/1000)*1000000000)/1000000000)
    localValue = localValue + addLine(postProcessed, "g/cm³", "gram/centimeter³")
    
    postProcessed = normalizeNumber("%3.9f", math.Round((isa.params.density * 1.94024058983E-03/53.7056928034)*1000000000)/1000000000)
    localValue = localValue + addLine(postProcessed, "lb/in³", "pound/inch³")

    postProcessed = normalizeNumber("%3.9f", math.Round((isa.params.density * 1.94024058983E-03/3.10848616723E-02)*1000000000)/1000000000)
    localValue = localValue + addLine(postProcessed, "lb/ft³", "pound/foot³")
    
    postProcessed = normalizeNumber("%3.9f", math.Round((isa.params.densityRatio)*1000000000)/1000000000)
    localValue = localValue + addLine(postProcessed, "", "[σ] (Ratio ρ/ρ0)")
   // localValue = localValue + endTable()

    return localValue
}
    
func (isa *isa76) showAllMetrics(abrvFrom string) string {
    var postProcessed = ""
    var localValue = ""
    
    localValue = formatTitle("All Parameters")

    // layer name
    localValue = localValue + addHeadline(isa.params.layerName, "", "Layer Name")
    
    // geometric
    postProcessed = normalizeNumber("%5.2f", math.Round(isa.params.geometricAltMeters*100)/100)
    localValue = localValue + addLine(postProcessed, "m", "Geometric Alt.")
    
    // geopotential
    postProcessed = normalizeNumber("%5.2f", math.Round(isa.params.geopotentialAltMeters*100)/100)
    localValue = localValue + addLine(postProcessed, "m", "Geopotential Alt.")
    
    // density
    postProcessed = normalizeNumber("%3.9f", math.Round(isa.params.density*1000000000)/1000000000)
    localValue = localValue + addLine(postProcessed, "kg/m³", "Density")
    
    // speed of sound
    metricValue := math.Pow((GAMMA_AIR * isa.params.temperatureK * R_SGC), 0.5)
    postProcessed = normalizeNumber("%4.3f", math.Round(metricValue*1000)/1000)
    localValue = localValue + addLine(postProcessed, "m/s", "Speed of sound")
    
    // pressure
    postProcessed = normalizeNumber("%7.6f", math.Round(isa.params.pressure*1000)/1000)
    localValue = localValue + addLine(postProcessed, "pa", "Pressure")

    // temperature
    postProcessed = normalizeNumber("%4.2f", math.Round(isa.params.temperatureK*1000)/1000)
    localValue = localValue + addLine(postProcessed, "K", "Temperature")
    
    // gravity acceleration
    // temperature
    localValue = localValue + addLine(isa.showEarthAcceleration(), "m/s²", "Accel. of gravity")

    //localValue += output.endTable()
    return localValue

}
    
func (isa *isa76) showCustomAltitude(fromValue string) string {
    var postProcessed = ""
    if isa.lastAltitude > -610 {
        
        isa.custAltitude = strings.ToUpper("Geometric Alt:") + "\n"
        postProcessed = normalizeNumber("%7.6f", math.Round(isa.params.geometricAltMeters*1000)/1000)
        isa.custAltitude = isa.custAltitude + postProcessed + " m ("
        postProcessed = normalizeNumber("%7.6f", math.Round(isa.params.geometricAltFeet*1000)/1000)
        isa.custAltitude = isa.custAltitude + postProcessed + " ft)\n"
        
        isa.custAltitude = isa.custAltitude + strings.ToUpper("\nGeopotential Alt:") + "\n"
        postProcessed  = normalizeNumber("%7.6f", math.Round(isa.params.geopotentialAltMeters*1000)/1000)
        isa.custAltitude = isa.custAltitude + postProcessed + " m ("
        postProcessed  = normalizeNumber("%7.6f", math.Round(isa.params.geopotentialAltFeet*1000)/1000)
        isa.custAltitude = isa.custAltitude + postProcessed + " ft)\n"
        
        isa.custAltitude = isa.custAltitude + strings.ToUpper("\nTemperature:") + "\n"
        postProcessed = normalizeNumber("%7.6f", math.Round(isa.params.temperatureK*1000)/1000)
        isa.custAltitude = isa.custAltitude + postProcessed + " Kelvin\n"
        
        postProcessed = normalizeNumber("%7.6f", math.Round(isa.params.temperatureC*1000)/1000)
        isa.custAltitude = isa.custAltitude + postProcessed + " Celcius\n"
        postProcessed = normalizeNumber("%7.6f", math.Round(isa.params.temperatureF*1000)/1000)
        isa.custAltitude = isa.custAltitude + postProcessed + " Fahrenheit\n"
        
        
        
        isa.custAltitude = isa.custAltitude + strings.ToUpper("\nPressure:") + "\n"
        postProcessed = normalizeNumber("%7.6f", math.Round(isa.params.pressure*10000000)/10000000)
        isa.custAltitude = isa.custAltitude + postProcessed + " pa\n"
        
        isa.custAltitude = isa.custAltitude + strings.ToUpper("\nDensity:") + "\n"
        postProcessed = normalizeNumber("%7.6f", math.Round(isa.params.density*100000000)/100000000)
        isa.custAltitude = isa.custAltitude + postProcessed + " Kg/m³\n"
        
        isa.custAltitude = isa.custAltitude + strings.ToUpper("\nSpeed of sound:") + "\n"
        postProcessed = normalizeNumber("%7.6f", math.Round(math.Pow((GAMMA_AIR * isa.params.temperatureK * R_SGC), 0.5)*1000)/1000)
        isa.custAltitude = isa.custAltitude + postProcessed + " m/s\n"
        
        isa.custAltitude = isa.custAltitude + strings.ToUpper("\nLayer Name:") + "\n" + isa.params.layerName + "\n"
    }
    return isa.custAltitude
}
    
func (isa *isa76) showAbout() string {
    //let path = NSBundle.mainBundle().pathForResource("Comparison_US_standard_atmosphere_1962", ofType: "svg")
    //let baseURL = NSURL(fileURLWithPath: path!)
    
    //let myUrl = NSURL(string: "Comparison_US_standard_atmosphere_1962.svg");

//   let help = "<div style='text-align:justify;width:100%;'>The U.S. Standard Atmosphere is an atmospheric model of how the pressure, temperature, density, and viscosity of the Earth's atmosphere change over a wide range of altitudes or elevations. The model, based on an existing international standard, was first published in 1958 by the U.S. Committee on Extension to the Standard Atmosphere, and was updated to its most recent version in 1976. It is largely consistent in methodology with the International Standard Atmosphere, differing mainly in the assumed temperature distribution at higher altitudes.<br><br>This USSA calculator will determine the temperature, pressure, density, speed of sound, geopotential altitude, acceleration of gravity for any altitude between -5,000m (-16,404ft) and +86,000m (+282,152ft) given a standard sea-level temperature of 288.15 Kelvin and pressure of 101,325 pascals.<br><br><img src='USstandardAtmosphere1976.jpg' style='width:100%;'></div>"
//  return help

    return "The U.S. Standard Atmosphere is an atmospheric model of how the pressure, temperature, density, and viscosity of the Earth's atmosphere change over a wide range of altitudes or elevations. The model, based on an existing international standard, was first published in 1958 by the U.S. Committee on Extension to the Standard Atmosphere, and was updated to its most recent version in 1976. It is largely consistent in methodology with the International Standard Atmosphere, differing mainly in the assumed temperature distribution at higher altitudes.<br><br>This USSA calculator will determine the temperature, pressure, density, speed of sound, geopotential altitude, acceleration of gravity for any altitude between -5,000m (-16,404ft) and +86,000m (+282,152ft) given a standard sea-level temperature of 288.15 Kelvin and pressure of 101,325 pascals."
}
    
func (isa *isa76) showAllLayers() string {
    return "whoops!"
    /*
    if !layersInitialized {
        var postProcessed = ""
        for layer in atmLayers {
            layersDisplay = layersDisplay+"\(layer.name)".uppercaseString+" layer:\n"
            postProcessed = normalizeNumber("\(Int(layer.base_altitude))")
            
            layersDisplay = layersDisplay+"["+postProcessed
            postProcessed = normalizeNumber("\(Int(layer.top_altitude))")
            
            layersDisplay = layersDisplay+" m ➡︎ "+postProcessed+" m]\n"
            layersDisplay = layersDisplay+"Lapse Rate    = \(layer.lapse_rate) K/m\n"
            postProcessed = normalizeNumber("\(Int(layer.base_altitude))")
            layersDisplay = layersDisplay+"Base Altitude = "+postProcessed+" m\n"
            postProcessed = normalizeNumber("\(Int(layer.top_altitude))")
            layersDisplay = layersDisplay+"Top  Altitude = "+postProcessed+" m\n"
            layersDisplay = layersDisplay+"Base Temp (K) = \(round(layer.base_temperature_K*1000)/1000) K\n"
            layersDisplay = layersDisplay+"Base Temp (C) = \(round(layer.base_temperature_C*1000)/1000) C\n"
            layersDisplay = layersDisplay+"Base Temp (F) = \(round(layer.base_temperature_F*1000)/1000) F\n"
            print("before-> \(round(layer.base_pressure*1000)/1000)")
            postProcessed = normalizeNumber("\(round(layer.base_pressure*1000)/1000)")
            print("after-> "+postProcessed)
            layersDisplay = layersDisplay+"Base Pressure = "+postProcessed+" pa\n"
            layersDisplay = layersDisplay+"Base Density  = \(round(layer.base_density*1000000)/1000000) kg/m³\n"
            layersDisplay = layersDisplay+"Top  Density  = \(round(layer.top_density*1000000)/1000000) kg/m³\n\n"
            
        }
        layersInitialized = true
    }
    return layersDisplay*/
}

func mainISA() {
    var alt int
    var unit string
    var isa = isa76{}
    isa.set()
    fmt.Println(isa.showAbout()+"\n")
    isa.setSealevelTemp(STD_SEALEVEL_TEMP) //+11.5)
    for {
        fmt.Println("Enter an Altitude between -5000m (-16404ft) and +86000m (282152ft) followed by a unit (m or f)")
        unit = "m"
        fmt.Scanf("%d %s", &alt, &unit)
        fmt.Println(isa.customAltitude(strconv.Itoa(alt), unit))
        //fmt.Println(isa.showAllMetrics(""))
        //fmt.Println(isa.showCustomAltitude(""))
        fmt.Println(isa.showTemperature(""))
        fmt.Println(isa.showSoundSpeed(""))
        fmt.Println(isa.showDensity(""))
        fmt.Println(isa.showPressure(""))
        fmt.Println(isa.showEarthAcceleration())
        
        
    }
}