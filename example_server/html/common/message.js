function Msg(to, content) {
	return {op: "cMsg", to: to, content: content};
}

function Add(id) {
	return {op: "cAdd", id: id};
}

function Move(lat, lng) {
	return {op: "cMove", lat: lat, lng: lng};
}

function InitLoc(lat, lng) {
	return {op: "cInitLoc", lat: lat, lng: lng};
}
