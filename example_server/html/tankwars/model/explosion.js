function Explosion(x, y, life, radius) {
	this.x = x;
	this.y = y;
	this.life = life;
	this.radius = radius;
}

Explosion.prototype.setClear = function(ctxt, hgt) {
	var x = this.x - this.radius-2;
	var y = hgt - (this.y + this.radius + 2);
	var w = this.radius*2 + 4;
	var h = this.radius*2 + 4;
	ctxt.clearRect(x,y,w,h);
}

Explosion.prototype.deplete = function() {
	this.life--;
}

Explosion.prototype.shouldRemove = function() {
	return this.life <= 0;
}

Explosion.prototype.render = function(ctxt, hgt) {
	var x = Math.floor(this.x);
	var y = Math.floor(this.y);
	ctxt.beginPath();
	ctxt.arc(x, hgt-y, this.radius, 0, 2*Math.PI, true);
	ctxt.closePath();
	ctxt.fill();
}
