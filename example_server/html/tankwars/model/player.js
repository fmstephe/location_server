function Player(id, x, name, turretLength, initPower, minPower, maxPower, powerInc, rotateInc, health, remote) {
	this.id = id;
	this.x = x;
	this.y = 0; // This gets set automatically by the game loop
	this.name = name;
	this.turretLength = turretLength;
	this.arc = 0;
	this.power = initPower;
	this.minPower = minPower;
	this.maxPower = maxPower;
	this.powerInc = powerInc;
	this.rotateInc = rotateInc;
	this.health = health;
	this.maxHealth = health;
	this.remote = remote;
	this.animate = !remote;
	this.cycle = 0.0;
}

Player.prototype.incPower = function() {
	this.power += this.powerInc;
	this.power = Math.min(this.power, this.maxPower);
}

Player.prototype.decPower = function() {
	this.power -= this.powerInc;
	this.power = Math.max(this.power, this.minPower);
}

Player.prototype.rotateLeft = function() {
	this.arc -= this.rotateInc;
}

Player.prototype.rotateRight = function() {
	this.arc += this.rotateInc;
}

Player.prototype.setClear = function(ctxt, hgt) {
	var x = this.x - (this.turretLength+5);
	var y = hgt - (this.y + this.turretLength + 5);
	var w = this.turretLength * 6; // This is a cludge value to allow for clearing power % text
	var h = this.turretLength * 6;
	ctxt.clearRect(x, y, w, h);
}

Player.prototype.shouldRemove = function() {
	return this.health <= 0;
}

Player.prototype.render = function(ctxt, hgt) {
	if (this.health > 0) {
		var r, g, b, d;
		if (this.remote) {
			r = 255*this.health/this.maxHealth;
			g = 50*this.health/this.maxHealth;
		       	b = 50*this.health/this.maxHealth;
		       	d = 255;//*this.health/this.maxHealth;
		} else {
			r = 50*this.health/this.maxHealth;
			g = 50*this.health/this.maxHealth;
		       	b = 255*this.health/this.maxHealth;
		       	d = 255*this.health/this.maxHealth;
			if (this.animate) {
				var mult = (Math.cos(this.cycle)+1.5)/ 2;
				this.cycle = this.cycle+0.2;
				r = r*mult;
				g = g*mult;
				b = b*mult;
				d = d*mult;
			}
		}
		r = rgbLim(r);
		g = rgbLim(g);
		b = rgbLim(b);
		d = rgbLim(d);
		// Do tank body
		ctxt.fillStyle = "rgba("+r+","+g+","+b+",1.0)";
		ctxt.beginPath();
		ctxt.arc(this.x, hgt-this.y, 10, 0, 2*Math.PI, true);
		ctxt.closePath();
		ctxt.fill();
		// Do turret
		ctxt.strokeStyle = "rgba("+d+","+d+","+d+",1.0)";
		turretX = this.x+this.turretLength*Math.sin(this.arc);
		turretY = hgt-(this.y+(this.turretLength*Math.cos(this.arc)));
		ctxt.beginPath();
		ctxt.moveTo(this.x, hgt-this.y);
		ctxt.lineTo(turretX,turretY);
		ctxt.closePath();
		ctxt.stroke();
		var powerP = Math.round((this.power/this.maxPower)*100);
		// Do text
		ctxt.font = "16pt Calibri-bold";
		ctxt.fillStyle = "rgba(255,255,255,1.0)";
		if (!this.remote) ctxt.fillText(powerP+"%", this.x+this.turretLength, hgt-this.y);
		ctxt.fillText(formatPlayerName(this.name), this.x-this.turretLength, hgt-this.y+this.turretLength+10, 80);
	}
}

function formatPlayerName(name) {
	if (name.length > 15) {
		return name.substring(0,12) + "...";
	} else {
		return name;
	}
}

function rgbLim(v) {
	return Math.floor(Math.min(255, Math.max(0,v)));
}
