function Add(id) {
	this.Op = "cAdd";
	this.Id = id; 
}

function Msg(to, msg) {
	this.To = to
	this.Msg = msg;
}

function Move(lat, lng) {
	this.Op = "cMove";
	this.Lat = lat;
	this.Lng = lng;
}

function InitLoc(lat, lng) {
	this.Op = "cInitLoc";
	this.Lat = lat;
	this.Lng = lng;
}

function Nearby(lat, lng) {
	this.Op = "cNearby";
	this.Lat = lat;
	this.Lng = lng;
}
