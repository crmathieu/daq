//function fun1(x) { return Math.sin(x); }
//function fun2(x) { return Math.cos(3 * x); }


//var canvas = document.getElementById("DAQcanvas");
var ctx;
var axis = {};

function draw() {
//    var canvas = document.getElementById("canvas");
    var canvas = document.getElementById("DAQcanvas");
    if (null == canvas || !canvas.getContext) return;

//    var axis = {};
    ctx = canvas.getContext("2d");
//    axis.x0 = .5 + .5 * canvas.width;  // x0 pixels from left to x=0
//    axis.y0 = .5 + .5 * canvas.height; // y0 pixels from top to y=0
    axis.x0 = .5;  // x0 pixels from left to x=0
    axis.y0 = canvas.height - .5; // y0 pixels from top to y=0
    axis.scale = 10; //2; //40;                 // 40 pixels from x=0 to x=1
    axis.doNegativeX = false; //true;
    axis.doNegativeY = false; //true;

    ctx.lineJoin = 'round';
    ctx.lineWidth = 4;
    
    showAxes(ctx, axis);
    //funGraph(ctx, axis, fun1, "rgb(11,153,11)", 1);
    //funGraph(ctx, axis, fun2, "rgb(66,44,255)", 2);
}

function funGraph(ctx, axis, func, color, thick) {
    var xx, yy, dx = 4, x0 = axis.x0, y0 = axis.y0, scale = axis.scale;
    var iMax = Math.round((ctx.canvas.width - x0) / dx);
    var iMin = axis.doNegativeX ? Math.round(-x0 / dx) : 0;
    ctx.beginPath();
    ctx.lineWidth = thick;
    ctx.strokeStyle = color;

    for (var i = iMin; i <= iMax; i++) {
        xx = dx * i; yy = scale * func(xx / scale);
        if (i == iMin) ctx.moveTo(x0 + xx, y0 - yy);
        else ctx.lineTo(x0 + xx, y0 - yy);
    }
    ctx.stroke();
}

function showAxes(ctx, axis) {
    var x0 = axis.x0, w = ctx.canvas.width;
    var y0 = axis.y0, h = ctx.canvas.height;
    var xmin = axis.doNegativeX ? 0 : x0;
    var ymin = axis.doNegativeY ? 0 : y0;//
    ctx.beginPath();
    ctx.strokeStyle = "#00000"; //"rgb(128,128,128)";
//    ctx.moveTo(xmin, y0); ctx.lineTo(w, y0);  // X axis
    ctx.moveTo(xmin, ymin); ctx.lineTo(w, y0);  // X axis
    ctx.moveTo(x0, 0); ctx.lineTo(x0, h);  // Y axis
    ctx.stroke();
}

function plot2(ctx, axis, color, thick, iteration, range, altitude) {
    var xx, yy, dx = 4, x0 = axis.x0, y0 = axis.y0, scale = axis.scale;
//    var iMax = Math.round((ctx.canvas.width - x0) / dx);
//    var iMin = axis.doNegativeX ? Math.round(-x0 / dx) : 0;
//    ctx.beginPath();
    ctx.lineWidth = thick;
    ctx.strokeStyle = color;

    xx = Math.round(x0 + (range/10));
//    yy = scale * altitude/1000;
    yy = Math.round(y0 - (altitude / 10));

/*    for (var i = iMin; i <= iMax; i++) {
        xx = dx * i; yy = scale * func(xx / scale);
        if (i == iMin) ctx.moveTo(x0 + xx, y0 - yy);
        else ctx.lineTo(x0 + xx, y0 - yy);
    }*/
    if (iteration < 1) {
        ctx.moveTo(xx, yy);
    } else {
        ctx.lineTo(xx, yy);  
        document.getElementById("X").innerHTML = xx;
        document.getElementById("Y").innerHTML = yy;
    }
    ctx.stroke();

}

class position {
    constructor() {
        this.curx = axis.x0;
        this.cury = axis.y0;
    }

    set(x, y) {
        this.curx = x;
        this.cury = y;
    }

    plot(ctx, axis, color, thick, range, altitude) {
        var xx, yy, x0 = axis.x0, y0 = axis.y0, scale = axis.scale;

        if (this.curx == -1 && this.cury == -1) {
            return;
        }
        xx = Math.round(x0 + range);
        yy = Math.round(y0 - altitude);

            
        // reset path to this curve
        ctx.beginPath();
        ctx.lineWidth = thick;
        ctx.strokeStyle = color;
        ctx.moveTo(this.curx, this.cury);

        ctx.lineTo(xx, yy);
        this.curx = xx;
        this.cury = yy;
        ctx.stroke();

    }
}

function plot(ctx, axis, color, thick, iteration, range, altitude) {
    var xx, yy, dx = 4, x0 = axis.x0, y0 = axis.y0, scale = axis.scale;

    //    var iMax = Math.round((ctx.canvas.width - x0) / dx);
    //    var iMin = axis.doNegativeX ? Math.round(-x0 / dx) : 0;
    //    ctx.beginPath();
    ctx.lineWidth = thick;
    ctx.strokeStyle = color;

//    xx = Math.round(x0 + (range / 1000));
//    yy = Math.round(y0 - (altitude / 1000));

    xx = Math.round(x0 + range);
    yy = Math.round(y0 - altitude);

    /*    for (var i = iMin; i <= iMax; i++) {
            xx = dx * i; yy = scale * func(xx / scale);
            if (i == iMin) ctx.moveTo(x0 + xx, y0 - yy);
            else ctx.lineTo(x0 + xx, y0 - yy);
        }*/
    if (iteration < 1) {
        ctx.moveTo(xx, yy);
    } else {
        ctx.lineTo(xx, yy);
        document.getElementById("X").innerHTML = xx;
        document.getElementById("Y").innerHTML = y0-yy;
    }
    ctx.stroke();

}