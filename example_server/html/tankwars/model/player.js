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
	var x = this.x - this.turretLength;
	var y = hgt - (this.y + this.turretLength + 5);
	var w = this.turretLength * 6; // This is a cludge value to allow for clearing power % text
	var h = this.turretLength * 2;
	ctxt.clearRect(x, y, w, h);
}

Player.prototype.shouldRemove = function() {
	return this.health <= 0;
}

Player.prototype.render = function(ctxt, hgt) {
	if (this.health > 0) {
		if (this.remote) {
			var r = Math.floor(255*this.health/this.maxHealth);
			var g = Math.floor(50*this.health/this.maxHealth);
		       	var b = Math.floor(50*this.health/this.maxHealth);
			ctxt.fillStyle = "rgba("+r+","+g+","+b+",1.0)";
		} else {
			var r = Math.floor(50*this.health/this.maxHealth);
			var g = Math.floor(50*this.health/this.maxHealth);
		       	var b = Math.floor(255*this.health/this.maxHealth);
			ctxt.fillStyle = "rgba("+r+","+g+","+b+",1.0)";
		}
		var d = Math.floor(255*this.health/this.maxHealth);
		ctxt.strokeStyle = "rgba("+d+","+d+","+d+",1.0)";
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
		ctxt.font = "16pt Calibri-bold";
		if (!this.remote) ctxt.fillText(powerP+"%", this.x+this.turretLength, hgt-this.y);
		ctxt.fillText(this.name, this.x-this.turretLength, hgt-this.y+this.turretLength+10);
	}
}
