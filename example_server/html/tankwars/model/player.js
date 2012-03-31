function Player(id, x, name, turretLength, initPower, minPower, maxPower, powerInc, rotateInc, health) {
	this.id = id;
	this.x = x;
	this.y = 0; // This gets set automatically by the game loop
	this.name = name;
	this.arc = 0;
	this.power = initPower;
	this.minPower = minPower;
	this.maxPower = maxPower;
	this.powerInc = powerInc;
	this.rotateInc = rotateInc;
	this.health = health;
	this.turretLength = turretLength;
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
	var x = this.x-this.turretLength;
	var y = hgt - (this.y + this.turretLength);
	var w = this.turretLength*6; // This is a cludge value to allow for clearing power % text
	var h = this.turretLength*2;
	ctxt.clearRect(x, y, w, h);
}

Player.prototype.shouldRemove = function() {
	return this.health <= 0;
}

Player.prototype.render = function(ctxt, hgt) {
	if (this.health > 0) {
		ctxt.beginPath();
		ctxt.arc(this.x, hgt-this.y, 10, 0, 2*Math.PI, true);
		ctxt.closePath();
		ctxt.fill();
		turretX = this.x+this.turretLength*Math.sin(this.arc);
		turretY = hgt-(this.y+(this.turretLength*Math.cos(this.arc)));
		ctxt.beginPath();
		ctxt.moveTo(this.x, hgt-this.y);
		ctxt.lineTo(turretX,turretY);
		ctxt.closePath();
		ctxt.stroke();
		var powerP = Math.round((this.power/this.maxPower)*100);
		ctxt.font = "20pt Calibri-bold";
		ctxt.fillText(powerP+"%",this.x+this.turretLength, hgt-this.y);
	}
}
