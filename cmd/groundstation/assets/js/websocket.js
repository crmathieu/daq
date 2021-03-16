var restartId = 0;
var timeout = 60000;
var restartTO = 75000;
var isConnOpen = false;
var timerId = 0;
var conn;
var bigEndian = isBigEndian()

var eventHist = 0;
var eventOffset = 20;
var nsec = 0;
var timeHist = -20;
var initStage2Position = false;

//var gravityTurn;
//var flightPath;

//var boosterPos = new position();
//var stage2Pos = new position();
//stage2Pos.set(-1, -1);

class DAQ {
    
    constructor(websocketUrl, clientToken) {
        this.websocketServerAddr = websocketUrl;
        this.clientToken = clientToken;
    }
}

function start(websocketServerAddr) {
    // first let the user know if there has been an issue on the backend
    //        if (!{{.Error }}) {
    //            errMessage.innerHTML = {{.Error }};
    //         }

    conn = new WebSocket(websocketServerAddr);
    conn.binaryType = 'arraybuffer';
    //draw();
//    boosterPos.setScale();
//    stage2Pos.setScale();

//    gravityTurn = new gravityT("DAQcanvas", CANVAS_WIDTH, CANVAS_SIDE);
//    flightPath = new flightP("FPcanvas", 100, 100);

    //-----------------------------------------------------------------------
    conn.onopen = function (evt) {
        // inside the prototype, "this" represents the websocket object, not the DAQ object
        isConnOpen = true
        //console.log('%c '+ this.clientToken + '\n IS LIVE!',
        //    'font-weight: bold; font-size: 50px;color: red; text-shadow: 3px 3px 0 rgb(217,31,38) , 6px 6px 0 rgb(226,91,14) , 9px 9px 0 rgb(245,221,8) , 12px 12px 0 rgb(5,148,68) , 15px 15px 0 rgb(2,135,206) , 18px 18px 0 rgb(4,77,145) , 21px 21px 0 rgb(42,21,113)'
        //);
        //this.conn.send(this.clientToken);
        this.send(clientToken);
        timerId = setTimeout(keepAlive, timeout, conn);
    };

    //---------------------------------------------WebSocket--------------------------
    conn.onmessage = function (evt) {
        const PACKET_SIZE = 16;
        const PACKET_GRP = 16; //2;

        let dataset = evt.data
        let littleEndian = !bigEndian
        let offset = 0;
        
        gravityTurn.clear();

        //redrawCanvas(null);
        for (; offset <= (PACKET_GRP - 1) * PACKET_SIZE; offset += PACKET_SIZE) {
            let dataView = new DataView(dataset, offset, PACKET_SIZE); //evt.data)

            let stageOffset = 0;
            let color = "purple";
            let stagePos = gravityTurn.boosterPos;

            var packetType = dataView.getUint32(0, littleEndian)
            var stage = packetType >> 16;
            var packetId = packetType & 0xffff;

            switch (packetId) {
                case IDVELOCITY: // velocity:	Id, Velocity, Acceleration, Stage
                    //let velocityView = new Float32Array(evt.data, 2, 2);
                    //var stage = dataView.getInt32(12, littleEndian).toString();
                    //document.getElementById("Stage").innerHTML = dataView.getInt32(12, littleEndian).toString(); 
                    if (stage == 1) {
                        stageOffset = STAGE1OFFSET;
                    }
                    gravityTurn.drawTextAtPoint("Velocity: " + (dataView.getFloat32(4, littleEndian) * 3600).toFixed(2), 15, 25 + stageOffset);
                    gravityTurn.drawTextAtPoint("Acceleration: " + (dataView.getFloat32(8, littleEndian)).toFixed(2), 15, 45 + stageOffset);
                    console.log(evt.data);
                    break;

                case IDPOSITION: // position: 	Id, Range, Altitude, Stage
                    //var stage = dataView.getInt32(12, littleEndian).toString();
                    if (stage == 1) {
                        stageOffset = STAGE1OFFSET;
                        color = "blue";
                        stagePos = gravityTurn.stage2Pos;
                    }
                    gravityTurn.drawTextAtPoint("Range: " + dataView.getFloat32(4, littleEndian).toFixed(2), 15, 65 + stageOffset);
                    gravityTurn.drawTextAtPoint("Altitude: " + dataView.getFloat32(8, littleEndian).toFixed(2), 15, 85 + stageOffset);
                    //document.getElementById("Range_" + stage).innerHTML = dataView.getFloat32(4, littleEndian).toFixed(2); //positionView[0];
                    //document.getElementById("Altitude_" + stage).innerHTML = dataView.getFloat32(8, littleEndian).toFixed(2); //positionView[1];

                    stagePos.plotPolar(color, 2, dataView.getFloat32(4, littleEndian), dataView.getFloat32(8, littleEndian));
                    break;

                case IDANGLES: // angles: 	Id, alpha, beta, gamma
//                    if (stage == 1) {
//                        stageOffset = STAGE1OFFSET;
//                        color = "green";
//                        stagePos = gravityTurn.stage2Pos;
//                    }
                    var plum = false;
                    if ((eventHist & E_STAGESEP) == 0) {
                        if ((eventHist & E_MEI_1) && !(eventHist & E_MECO_1)) {
                            plum = true;
                        }
                        if (stage == 0) {
                            flightPathBooster.drawAngle(dataView.getFloat32(12, littleEndian), plum);
                            flightPathBooster.drawTextAtPoint((dataView.getFloat32(12, littleEndian) * 180 / Math.PI).toFixed(2), 25, 25);
                            flightPathStage1.drawAngle(dataView.getFloat32(12, littleEndian), false);
                            flightPathStage1.drawTextAtPoint((dataView.getFloat32(12, littleEndian) * 180 / Math.PI).toFixed(2), 25, 25);
                        }
                    } else {
                        if (stage == 0) {
                            if (((eventHist & E_EBURNI) && !(eventHist & E_EBURNO)) || 
                                ((eventHist & E_BBURNI) && !(eventHist & E_BBURNO)) || 
                                ((eventHist & E_LBURNI) && !(eventHist & E_LBURNO))) {
                                plum = true;
                            }
                            flightPathBooster.drawAngle(dataView.getFloat32(12, littleEndian), plum);
                            flightPathBooster.drawTextAtPoint((dataView.getFloat32(12, littleEndian) * 180 / Math.PI).toFixed(2), 25, 25);
                        } else {
                            if ((eventHist & E_SEI_1) && !(eventHist & E_SECO_1)) {
                                plum = true;
                            }
                            flightPathStage1.drawAngle(dataView.getFloat32(12, littleEndian), plum);
                            flightPathStage1.drawTextAtPoint((dataView.getFloat32(12, littleEndian) * 180 / Math.PI).toFixed(2), 25, 25);
                        }
                    }
                    //document.getElementById("Range_" + stage).innerHTML = dataView.getFloat32(4, littleEndian).toFixed(2); //positionView[0];
                    //document.getElementById("Altitude_" + stage).innerHTML = dataView.getFloat32(8, littleEndian).toFixed(2); //positionView[1];

                    //stagePos.plot(color, 2, dataView.getFloat32(4, littleEndian), dataView.getFloat32(8, littleEndian));
                    break;

                case IDTIME:
                    var time = dataView.getFloat32(4, littleEndian);
                    nsec = Math.round(time);
                    if (nsec > timeHist) {
                        //document.getElementById("time").innerHTML = new Date(nsec * 1000).toISOString().substr(11, 8) + " (" + nsec.toString() + ")";
                        timeHist = nsec;
                    }
                    break;

                case IDEVENT:
                    var eventID = dataView.getInt32(4, littleEndian);
                    eventHist = dataView.getInt32(12, littleEndian);
                    if (lastEvent != eventID) {
                        //document.getElementById("EventsList").innerHTML += "<br>" + eventMap.get(eventID) + " (" + nsec + ")";
                        eventsRecorder.drawTextAtPoint("(" + nsec + ") " + eventMap.get(eventID), 5, eventOffset)
                        eventOffset += 20;
                        lastEvent = eventID;
                    }
                    if ((eventHist & E_STAGESEP) && !initStage2Position) {
                        gravityTurn.stage2Pos.set(gravityTurn.boosterPos.x, gravityTurn.boosterPos.y);
                        initStage2Position = true;
                    }

                default: // ignore for now
                // let thrustView = new Float32Array(evt.data, 2, 1);
                // let stageView = new Uint8Array(evt.data, 6, 1);
            }
        }
    };

    //-----------------------------------------------------------------------
    conn.onclose = function (evt) {
        isConnOpen = false;
        cancelKeepAlive()
        if (evt.code === 1006) {
            // 1006 is abnormal disconnect - Try to reconnect in 2 seconds
            //this.restartId = setTimeout(this.restartConnection, 1000);
            start(websocketServerAddr);
        }
    };

}

function keepAlive(ws) {
    if (ws.readyState == ws.OPEN) {
        ws.send('_keepalive_');
    }
    timerId = setTimeout(keepAlive, timeout, ws);
}

function cancelKeepAlive() {
    if (timerId) {
        clearTimeout(timerId);
    }
}

function restartConnection() {
    if (!isConnOpen) {
        console.log('Reconnecting WebSocket...');
        start(websocketServerAddr)
        restartId = setTimeout(restartConnection, restartTO);
    } else {
        if (restartId) {
            clearTimeout(restartId)
        }
    }
}

function isBigEndian() {
    const array = new Uint8Array(4);
    const view = new Uint32Array(array.buffer);
    return !((view[0] = 1) & array[0]);
}
