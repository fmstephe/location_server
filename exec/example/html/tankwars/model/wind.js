function Wind(width, height, particleNum) {
	this.width = width;
	this.height = height;
	this.particleNum = particleNum;
	this.wind = 0;
	this.particles = new Array();
	for (var i = 0; i < particleNum; i++) {
		var x = r(width);
		var y = r(height);
		this.particles[i] = [x,y,x,y];
	}
}

function windValue() {
	return r(30) - 15;
}

Wind.prototype.windChange = function() {
	return r(14) - 7;
}

Wind.prototype.addWind = function(diff) {
	this.wind += diff;
	if (this.wind > 15) {
		this.wind = 15;
	}
	if (this.wind < -15) {
		this.wind = -15;
	}
}

Wind.prototype.setClear = function(ctxt) {
	ctxt.fillStyle = "#000";
	ctxt.fillRect(0, 0, this.width, this.height);
}

Wind.prototype.render = function(ctxt) {
	ctxt.beginPath();
	for (var i = 0; i < this.particleNum; i++) {
		var p = this.particles[i];
		if (p[0] != 0 && p[0] != this.width) {
			ctxt.moveTo(p[0],p[1]);
			ctxt.lineTo(p[2],p[3]);
		}
	}
	ctxt.closePath();
	ctxt.stroke();
}

Wind.prototype.advance = function() {
	for (var i = 0; i < this.particleNum; i++) {
		var p = this.particles[i];
		p[3] = p[1];
		p[2] = p[0];
		p[1] = p[1]; // Vertically particles don't change
		p[0] = p[0]+this.wind;
		if (p[0] < 0) {
			p[0] = this.width;
		}
		if (p[0] > this.width) {
			p[0] = 0;
		}
	}

}
