function KeyBindings(upKey, downKey, leftKey, rightKey, firingKey) {
	this.upKey = upKey;
	this.downKey = downKey;
	this.leftKey = leftKey;
	this.rightKey = rightKey;
	this.firingKey = firingKey;
	this.up = false;
	this.down = false;
	this.left = false;
	this.right = false;
	this.firing = false;
	this.reset = resetKeyBindings;
}

function resetKeyBindings() {
	this.up = false;
	this.down = false;
	this.left = false;
	this.right = false;
	this.firing = false;
}
