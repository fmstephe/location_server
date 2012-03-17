function Msg(to, content) {
	return {Op: "cMsg", To: to, Content: content};
}

function Add(id) {
	return {Op: "cAdd", Id: id};
}

function Move(lat, lng) {
	return {Op: "cMove", Lat: lat, Lng: lng};
}

function InitLoc(lat, lng) {
	return {Op: "cInitLoc", Lat: lat, Lng: lng};
}
