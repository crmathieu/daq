const RADIUS = 200;
const X_CIRCLE_CENTER = 400;
const Y_CIRCLE_CENTER = 400;

let canvas;
let ctx;

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
    canvas = document.getElementById("DAQcanvas");
    ctx = canvas.getContext('2d');

    drawCanvas();

    canvas.addEventListener("mousemove", redrawCanvas); //());

}

const CANVAS_STARTX = 0;
const CANVAS_STARTY = 0;

const CANVAS_SIDE = 800;
const CANVAS_ENDX = CANVAS_SIDE;
const CANVAS_ENDY = CANVAS_SIDE;

function drawCanvas() {
    drawRectangle('#839192', 5, CANVAS_STARTX, CANVAS_STARTY, CANVAS_ENDX, CANVAS_ENDY);
    drawCircle('#839192', 1, X_CIRCLE_CENTER, Y_CIRCLE_CENTER, RADIUS, 0, 2 * Math.PI);
    drawLine('#839192', 1, X_CIRCLE_CENTER, 0, X_CIRCLE_CENTER, CANVAS_SIDE);
    drawLine('#839192', 1, 0, Y_CIRCLE_CENTER, CANVAS_SIDE, Y_CIRCLE_CENTER);
}

function redrawCanvas(evt) {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    drawCanvas();
    getMousePosition(evt);
    drawTextAtPoint("X: " + mousePos.x, 15, 25);
    drawTextAtPoint("Y: " + mousePos.y, 15, 45);

    let angleOfMouseDegrees = getAngleUsingXandY(mousePos.x, mousePos.y);
    drawTextAtPoint("Degrees: "+ angleOfMouseDegrees, 15, 65);
    drawTriangle(angleOfMouseDegrees);
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
    ctx.fillText(text, x, y);
}


function getMousePosition(evt) {
    // get canvas coordinate vis a vis the document
    let canvasDimensions = canvas.getBoundingClientRect();
    mousePos.x = Math.floor(evt.clientX - canvasDimensions.left);    
    mousePos.y = Math.floor(evt.clientY - canvasDimensions.top);
    
    mousePos.x -= CANVAS_SIDE / 2;
    mousePos.y = -1 * (mousePos.y - (CANVAS_SIDE/2));
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
        return (360.0 + (rad * 180/Math.PI)).toFixed(2);
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
    ctx.lineTo(xEndPoint, CANVAS_SIDE/2);
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

