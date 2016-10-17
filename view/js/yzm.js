/**
 * Created by yang on 2016/10/17.
 */
function createYZM(id) {
    var canvas = document.getElementById(id);
    var context = canvas.getContext("2d");
    context.clearRect(0, 0, canvas.width, canvas.height);

    var random = ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R',
        'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'];
    var colors = ["red", "green", "brown", "blue", "orange", "purple", "black"];
    var code = [];
    for (var i = 0; i < 4; i++) {
        var index = Math.floor(Math.random() * 36);
        code.push(random[index]);
    }

    context.beginPath();
    // Sprinkle in some random dots
    for (var i = 0; i < 10; i++) {
        var px = Math.floor(Math.random() * canvas.width);
        var py = Math.floor(Math.random() * canvas.height);
        context.moveTo(px, py);
        context.lineTo(px + 1, py + 1);
        context.strokeStyle = colors[Math.floor(Math.random() * colors.length)];
        context.lineWidth = Math.floor(Math.random() * 2);
        context.stroke();
    }

    for (var j = 0; j < 2; j++) {
        //随机线条
        context.moveTo(0, Math.floor(Math.random() * canvas.height));//随机线的起点x坐标是画布x坐标0位置，y坐标是画布高度的随机数
        context.lineTo(canvas.width, Math.floor(Math.random() * canvas.height));//随机线的终点x坐标是画布宽度，y坐标是画布高度的随机数
        context.lineWidth = 0.3;//随机线宽
        context.strokeStyle = colors[Math.floor(Math.random() * colors.length)];
        context.stroke();//描边，即起点描到终点
    }

    var deg, cos, sin, dg;
    context.font = "20px Arial";
    var cx = (canvas.width - 30) / 3;
    for (var j = 0; j < 4; j++) {
        context.fillStyle = colors[Math.floor(Math.random() * colors.length)];
        //产生一个正负30度以内的角度值以及一个用于变形的dg值
        dg = Math.random() * 4.5 / 10;
        deg = Math.floor(Math.random() * 60);
        deg = deg > 30 ? (30 - deg) : deg;
        cos = Math.cos(deg * Math.PI / 180);
        sin = Math.sin(deg * Math.PI / 180);

        context.save();
        context.setTransform(cos, sin + dg, -sin + dg, cos, cx * (j + 1) - 12, 18);
        context.fillText(code[j], 0, 0);
        context.restore();
    }

    canvas.onclick = function () {
        context.clearRect(0, 0, canvas.width, canvas.height);
        createYZM(id);
    }
}