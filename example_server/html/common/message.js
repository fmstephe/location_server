function ClientMessage(op, msg) {
	this.Op = op;
	this.Msg = msg;
}

function Add(id) {
	return new ClientMessage("cAdd", {Id: id});
}

function Msg(to, content) {
	return new ClientMessage("cMsg", {To: to, Content: content});
}

function Move(lat, lng) {
	return new ClientMessage("cMove", {Lat: lat, Lng: lng});
}

function InitLoc(lat, lng) {
	return new ClientMessage("cInitLoc", {Lat: lat, Lng: lng});
}

function Nearby(lat, lng) {
	return new ClientMessage("cNearby", {Lat: lat, Lng: lng});
}

