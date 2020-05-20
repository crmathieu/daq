//function fun1(x) { return Math.sin(x); }
//function fun2(x) { return Math.cos(3 * x); }

var canvas = document.getElementById("DAQcanvas");
var ctx;
var axes = {};
function draw() {
//    var canvas = document.getElementById("canvas");
    if (null == canvas || !canvas.getContext) return;

//    var axes = {};
    ctx = canvas.getContext("2d");
//    axes.x0 = .5 + .5 * canvas.width;  // x0 pixels from left to x=0
//    axes.y0 = .5 + .5 * canvas.height; // y0 pixels from top to y=0
    axes.x0 = .5;  // x0 pixels from left to x=0
    axes.y0 = canvas.height - .5; // y0 pixels from top to y=0
    axes.scale = 2; //40;                 // 40 pixels from x=0 to x=1
    axes.doNegativeX = false; //true;
    axes.doNegativeY = false; //true;

    showAxes(ctx, axes);
    //funGraph(ctx, axes, fun1, "rgb(11,153,11)", 1);
    //funGraph(ctx, axes, fun2, "rgb(66,44,255)", 2);
}

function funGraph(ctx, axes, func, color, thick) {
    var xx, yy, dx = 4, x0 = axes.x0, y0 = axes.y0, scale = axes.scale;
    var iMax = Math.round((ctx.canvas.width - x0) / dx);
    var iMin = axes.doNegativeX ? Math.round(-x0 / dx) : 0;
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

function showAxes(ctx, axes) {
    var x0 = axes.x0, w = ctx.canvas.width;
    var y0 = axes.y0, h = ctx.canvas.height;
    var xmin = axes.doNegativeX ? 0 : x0;
    var ymin = axes.doNegativeY ? 0 : y0;//
    ctx.beginPath();
    ctx.strokeStyle = "rgb(128,128,128)";
//    ctx.moveTo(xmin, y0); ctx.lineTo(w, y0);  // X axis
    ctx.moveTo(xmin, ymin); ctx.lineTo(w, y0);  // X axis
    ctx.moveTo(x0, 0); ctx.lineTo(x0, h);  // Y axis
    ctx.stroke();
}

function plot2(ctx, axes, color, thick, iteration, range, altitude) {
    var xx, yy, dx = 4, x0 = axes.x0, y0 = axes.y0, scale = axes.scale;
//    var iMax = Math.round((ctx.canvas.width - x0) / dx);
//    var iMin = axes.doNegativeX ? Math.round(-x0 / dx) : 0;
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

function plot(ctx, axes, color, thick, iteration, range, altitude) {
    var xx, yy, dx = 4, x0 = axes.x0, y0 = axes.y0, scale = axes.scale;
    //    var iMax = Math.round((ctx.canvas.width - x0) / dx);
    //    var iMin = axes.doNegativeX ? Math.round(-x0 / dx) : 0;
    //    ctx.beginPath();
    ctx.lineWidth = thick;
    ctx.strokeStyle = color;

    xx = Math.round(x0 + (range / 1000));
    //    yy = scale * altitude/1000;
    yy = Math.round(y0 - (altitude / 1000));

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