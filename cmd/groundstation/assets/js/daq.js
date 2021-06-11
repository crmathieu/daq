const CANVAS_WIDTH = 700; //800;
const RADIUS = 2000; //1000; //600; //200;
const X_CIRCLE_CENTER = CANVAS_WIDTH/4; // /3; //400;
const Y_CIRCLE_CENTER = RADIUS + 130; //660; //400;

const CANVAS_STARTX = 0;
const CANVAS_STARTY = 0;

const CANVAS_SIDE = 420; //800;
const CANVAS_ENDX = 700; //800; //CANVAS_SIDE;
const CANVAS_ENDY = 200; //CANVAS_SIDE;

class position {
    constructor(ctx) {
        this.x0 = X_CIRCLE_CENTER+2;
        this.y0 = Y_CIRCLE_CENTER - RADIUS;
        this.curx = this.x0;
        this.cury = this.y0;
        this.ctx = ctx;
    }

    setScale() {
        this.radius = RADIUS; //Math.round((canvas.height - axis.y0 - 100) / 2);
        //alert("radius="+this.radius.toString());
        this.scale = this.radius / 6378.137;
        //alert("scale=" + axis.scale.toString());
    }

    set(x, y) {
        this.curx = x;
        this.cury = y;
    }

    plot(color, thick, range, altitude) {
        var xx, yy, x0 = this.x0, y0 = this.y0, scale = this.scale;

        if (this.curx == -1 && this.cury == -1) {
            return;
        }
        //        xx = Math.round(x0 + range/scale);
        //        yy = Math.round(y0 - altitude/scale);

        xx = Math.round(x0 + range * scale);
        yy = Math.round(y0 - altitude * scale);

        //        drawTextAtPoint("X,Y: " + xx.toFixed(2) + ", " + yy.toFixed(2), X_CIRCLE_CENTER + CANVAS_WIDTH / 3 + 15, 25);
        if (xx == this.curx && yy == this.cury) {
            return;
        }

        // reset path to this curve
        this.ctx.beginPath();
        this.ctx.lineWidth = thick;
        this.ctx.strokeStyle = color;
        this.ctx.moveTo(this.curx, this.cury);

        this.ctx.lineTo(xx, yy);
        //this.ctx.quadraticCurveTo(this.curx, this.cury, xx, yy);
        this.curx = xx;
        this.cury = yy;
        this.ctx.stroke();

    }

    plotPolar(color, thick, range, altitude) {
        // from range we can deduct the angle alpha of local referential
        var alpha = range * this.scale / RADIUS;
        console.log(alpha);

        var xx, yy, x0 = this.x0, y0 = this.y0, scale = this.scale;

        if (this.curx == -1 && this.cury == -1) {
            return;
        }
        //        xx = Math.round(x0 + range/scale);
        //        yy = Math.round(y0 - altitude/scale);
        var xe = (RADIUS + altitude * this.scale) * Math.sin(alpha);
        var ye = (RADIUS + altitude * this.scale) * Math.cos(alpha);

        var polarx = x0 + xe;
        var polary = (y0 - (ye - RADIUS));

        xx = Math.round(polarx); //range * scale);
        yy = Math.round(polary); //altitude * scale);

        //        drawTextAtPoint("X,Y: " + xx.toFixed(2) + ", " + yy.toFixed(2), X_CIRCLE_CENTER + CANVAS_WIDTH / 3 + 15, 25);
        if (xx == this.curx && yy == this.cury) {
            return;
        }

        // reset path to this curve
        this.ctx.beginPath();
        this.ctx.lineWidth = thick;
        this.ctx.strokeStyle = color;
        this.ctx.moveTo(this.curx, this.cury);

        this.ctx.lineTo(xx, yy);
        //this.ctx.quadraticCurveTo(this.curx, this.cury, xx, yy);
        this.curx = xx;
        this.cury = yy;
        this.ctx.stroke();

    }
}

class canvasElt {
    constructor(id, width, height) {
        this.width = width;
        this.height = height;
        this.canvas = document.getElementById(id);
        if (this.canvas !== null) {
            this.ctx = this.canvas.getContext('2d');
        }
    }

    drawLine(strokeColor, lineWidth, xStart, yStart, xEnd, yEnd) {
        this.ctx.strokeStyle = strokeColor;
        this.ctx.lineWidth = lineWidth;

        // position to beginning of line
        this.ctx.moveTo(xStart, yStart);
        // define end of line
        this.ctx.lineTo(xEnd, yEnd);
        // draw
        this.ctx.stroke();
    }

    drawRectangle(strokeColor, lineWidth, xStart, yStart, xEnd, yEnd) {
        this.ctx.strokeStyle = strokeColor;
        this.ctx.lineWidth = lineWidth;
        this.ctx.strokeRect(xStart, yStart, xEnd, yEnd);
    }    


    drawTextAtPoint(text, x, y) {
        this.ctx.font = "15px Arial";
        //ctx.clearRect(x, y-15, X_CIRCLE_CENTER - 16, 18);
        this.ctx.fillText(text, x, y);
    }

    clear() {
        //        this.ctx.clearRect(5, 5, this.x_circle_center - 16, this.canvas.height - 10);
        this.ctx.clearRect(5, 5, this.width - 10, TIMEOFFSET); //15, TIMEOFFSET);
    }

}



class gravityT extends canvasElt {
    constructor(id, width, height) {
        super(id, width, height);
        this.boosterPos = new position(this.ctx);
        this.stage2Pos = new position(this.ctx);

        this.boosterPos.setScale();
        this.stage2Pos.setScale();
        this.stage2Pos.set(-1, -1);
        this.x_circle_center = this.width/4;
        this.y_circle_center = RADIUS + 130;
    }
    
    drawCanvas() {
//        this.drawRectangle('#839192', 1, CANVAS_STARTX, CANVAS_STARTY, CANVAS_ENDX, CANVAS_ENDY);
//        this.drawCircle('#839192', 1, X_CIRCLE_CENTER, Y_CIRCLE_CENTER, RADIUS, 3 * Math.PI / 2, 3 * Math.PI / 2 + Math.PI / 4);
//        this.drawLine('#839192', 1, X_CIRCLE_CENTER, 0, X_CIRCLE_CENTER, CANVAS_ENDY);

        this.drawRectangle('#839192', 1, 0, 0, this.width, this.height);
        this.drawCircle('#839192', 1, this.x_circle_center, this.y_circle_center, RADIUS, 3 * Math.PI / 2, 3 * Math.PI / 2 + Math.PI / 4);
        this.drawLine('#839192', 1, this.x_circle_center, 0, this.x_circle_center, this.height);
    }

    drawCircle(strokeColor, lineWidth, xCenter, yCenter, radius, arcStart, arcEnd) {
        this.ctx.strokeStyle = strokeColor;
        this.ctx.lineWidth = lineWidth;

        this.ctx.beginPath();
        this.ctx.arc(xCenter, yCenter, radius, arcStart, arcEnd);
        this.ctx.stroke(); // to actually draw it !!
    }

    clear() {   
//        this.ctx.clearRect(5, 5, this.x_circle_center - 16, this.canvas.height - 10);
        this.ctx.clearRect(5, 5, this.x_circle_center - 16, this.height - 10);
    }

}


class flightP extends canvasElt {
    constructor(id, width, height, label) {
        super(id, width, height-25);
    
        this.RLENGTH = (Math.sqrt(Math.pow(this.width, 2) + Math.pow(this.height, 2)) / 3);
        this.color = ['red', 'yellow'];
        this.index = 0;
        this.label = label;

    }

    drawCanvas() {
        this.drawRectangle('#839192', 1, 0, 0, this.width, this.height);
        this.drawTextAtPoint(this.label, 15, this.height+18);

    }

    clear() {
//        this.ctx.clearRect(5, 5, this.canvas.width - 10, this.canvas.height - 10);
        this.ctx.clearRect(5, 5, this.width - 10, this.height - 10);
    }

    drawAngle(angle, plum) {
        // calculate upper position
        let xshift = this.RLENGTH / 2 * Math.cos(angle);
        let yshift = this.RLENGTH / 2 * Math.sin(angle);
        let xflame = this.RLENGTH / 4 * Math.cos(angle);
        let yflame = this.RLENGTH / 4 * Math.sin(angle);
        let x1 = (this.width / 2) + xshift;
        let y1 = (this.height / 3) * 2 - yshift;
        let x2 = (this.width /  2) - xflame; //xshift;
        let y2 = (this.height / 3) * 2 + yflame; //yshift;
        this.clear();
        this.ctx.beginPath();
        this.drawLine('#839192', 5, this.width / 2, (this.height / 3) * 2, x1, y1);
        this.ctx.beginPath();
        if (plum) {
            this.ctx.moveTo(this.width / 2, (this.height / 3) * 2); //x2, y2);
            this.drawLine(this.color[this.index], 2, this.width / 2, (this.height / 3) * 2, x2, y2); //x1, y1);
            this.index = (this.index + 1) % 2;
        }
    }

    drawAngleBIS(angle, plum) {
        // calculate upper position
        let xshift = this.RLENGTH / 2 * Math.cos(angle);
        let yshift = this.RLENGTH / 2 * Math.sin(angle);
        let xflame = this.RLENGTH / 4 * Math.cos(angle);
        let yflame = this.RLENGTH / 4 * Math.sin(angle);
        let x1 = (this.canvas.width / 2) + xshift;
        let y1 = (this.canvas.height / 3) * 2 - yshift;
        let x2 = (this.canvas.width / 2) - xflame; //xshift;
        let y2 = (this.canvas.height / 3) * 2 + yflame; //yshift;
        this.clear();
        this.ctx.beginPath();
        //        this.drawLine('#839192', 5, x1, y1, x2, y2);
        //        this.ctx.moveTo(x2, y2);
        //        this.drawLine('red', 2, x2, y2, x1, y1);

        //        this.drawLine('#839192', 1, 2, (this.canvas.height / 3) * 2, this.canvas.width, (this.canvas.height / 3) * 2);
        this.drawLine('#839192', 5, this.canvas.width / 2, (this.canvas.height / 3) * 2, x1, y1);
        this.ctx.beginPath();
        if (plum) {
            this.ctx.moveTo(this.canvas.width / 2, (this.canvas.height / 3) * 2); //x2, y2);
            this.drawLine(this.color[this.index], 2, this.canvas.width / 2, (this.canvas.height / 3) * 2, x2, y2); //x1, y1);
            this.index = (this.index + 1) % 2;
        }
    }
}



let canvas;
let ctx;

var gravityTurn;
var flightPathBooster;
var flightPathStage1;
var eventsRecorder;

class MousePosition {
    constructor(x, y) {
        this.x = x;
        this.y = y;
    }
}

let mousePos = new MousePosition();
document.addEventListener("DOMContentLoaded", setupCanvas);

// drawing primitives
function setupCanvas() {
//    canvas = document.getElementById("DAQcanvas");
//    ctx = canvas.getContext('2d');

//    drawCanvas();

    gravityTurn = new gravityT("DAQcanvas", CANVAS_WIDTH, CANVAS_SIDE);
    gravityTurn.drawCanvas();
    flightPathBooster = new flightP("FPBOOSTERcanvas", 100, 125, "Booster FP");
    flightPathBooster.drawCanvas();

    flightPathStage1 = new flightP("FPSTAGE1canvas", 100, 125, "Stage1 FP");
    flightPathStage1.drawCanvas();

    eventsRecorder = new canvasElt("EVENTScanvas", 150, 500);
    //eventsRecorder.drawCanvas();

    //canvas.addEventListener("mousemove", redrawCanvas); //());

}


function drawCanvas() {
    drawRectangle('#839192', 5, CANVAS_STARTX, CANVAS_STARTY, CANVAS_ENDX, CANVAS_ENDY);
//    drawCircle('#839192', 1, X_CIRCLE_CENTER, Y_CIRCLE_CENTER, RADIUS, 0, 2 * Math.PI);
    drawCircle('#839192', 1, X_CIRCLE_CENTER, Y_CIRCLE_CENTER, RADIUS, 3 * Math.PI / 2, 3 * Math.PI / 2 + Math.PI / 4);
    drawLine('#839192', 1, X_CIRCLE_CENTER, 0, X_CIRCLE_CENTER, CANVAS_ENDY); //CANVAS_SIDE);
    //drawLine('#839192', 1, 0, Y_CIRCLE_CENTER, CANVAS_ENDX, Y_CIRCLE_CENTER); //CANVAS_SIDE, Y_CIRCLE_CENTER);
}

function redrawCanvas(evt) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    drawCanvas();
    //getMousePosition(evt);
    //drawTextAtPoint("X: " + mousePos.x, 15, 25);
    //drawTextAtPoint("Y: " + mousePos.y, 15, 45);

    //let angleOfMouseDegrees = getAngleUsingXandY(mousePos.x, mousePos.y);
    //drawTextAtPoint("Degrees: " + angleOfMouseDegrees, 15, 65);
    //drawTriangle(angleOfMouseDegrees);
}

function drawRectangle(strokeColor, lineWidth, xStart, yStart, xEnd, yEnd) {
    ctx.strokeStyle = strokeColor;
    ctx.lineWidth = lineWidth;

    ctx.strokeRect(xStart, yStart, xEnd, yEnd);
}


function drawCircle(strokeColor, lineWidth, xCenter, yCenter, radius, arcStart, arcEnd) {
    ctx.strokeStyle = strokeColor;
    ctx.lineWidth = lineWidth;

    ctx.beginPath();
    ctx.arc(xCenter, yCenter, radius, arcStart, arcEnd);
    ctx.stroke(); // to actually draw it !!
}

function drawLine(strokeColor, lineWidth, xStart, yStart, xEnd, yEnd) {
    // position to beginning of line
    ctx.moveTo(xStart, yStart);
    // define end of line
    ctx.lineTo(xEnd, yEnd);
    // draw
    ctx.stroke();
}

function drawTextAtPoint(text, x, y) {
    ctx.font = "15px Arial";
    //ctx.clearRect(x, y-15, X_CIRCLE_CENTER - 16, 18);
    ctx.fillText(text, x, y);
}

function clear() {
    ctx.clearRect(5, 5, X_CIRCLE_CENTER - 16, canvas.height - 10);
}

function getMousePosition(evt) {
    // get canvas coordinate vis a vis the document
    let canvasDimensions = canvas.getBoundingClientRect();
    mousePos.x = Math.floor(evt.clientX - canvasDimensions.left);
    mousePos.y = Math.floor(evt.clientY - canvasDimensions.top);

    mousePos.x -= CANVAS_SIDE / 2;
    mousePos.y = -1 * (mousePos.y - (CANVAS_SIDE / 2));
    return mousePos;
}

function getAngleUsingXandY(x, y) {
    let adjacent = x;
    let opposite = y;
    //    let rad = Math.atan2(opposite / adjacent);
    let rad = Math.atan2(opposite, adjacent);
    return radiansToDegrees(rad);
}

function radiansToDegrees(rad) {
    if (rad < 0) {
        return (360.0 + (rad * 180 / Math.PI)).toFixed(2);
    }
    return (rad * 180 / Math.PI).toFixed(2);
}

function degreesToRadians(deg) {
    return (deg * Math.PI / 180).toFixed(2);
}

function drawTriangle(angleDegrees) {
    ctx.moveTo(X_CIRCLE_CENTER, Y_CIRCLE_CENTER);
    // cosine = adjacent / hypothenus
    let xEndPoint = X_CIRCLE_CENTER + RADIUS * Math.cos(degreesToRadians(angleDegrees));
    let yEndPoint = Y_CIRCLE_CENTER + RADIUS * -(Math.sin(degreesToRadians(angleDegrees)));

    drawTextAtPoint("Radians : " + degreesToRadians(angleDegrees), 15, 85);

    ctx.lineTo(xEndPoint, yEndPoint);
    ctx.stroke();
    ctx.moveTo(xEndPoint, yEndPoint);
    ctx.lineTo(xEndPoint, CANVAS_SIDE / 2);
    ctx.stroke();

    // write current coordinates
    drawTextAtPoint("(" + xEndPoint.toFixed(2) + ", " + yEndPoint.toFixed(2) + ")", xEndPoint + 10, yEndPoint - 10);

    let adjacentLength = getLineLength(X_CIRCLE_CENTER, Y_CIRCLE_CENTER, xEndPoint, Y_CIRCLE_CENTER);
    drawTextAtPoint("Adjacent : " + adjacentLength.toFixed(2), 15, 105);

    let oppositeLength = getLineLength(xEndPoint, yEndPoint, xEndPoint, Y_CIRCLE_CENTER);
    drawTextAtPoint("Opposite : " + oppositeLength.toFixed(2), 15, 125);

}

function getLineLength(x1, y1, x2, y2) {
    let xd = x1 - x2;
    xd = xd * xd;
    let yd = y1 - y2;
    yd = yd * yd;

    return Math.sqrt(xd + yd);
}


