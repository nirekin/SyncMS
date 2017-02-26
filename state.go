package main

import ()

var (
	stateMap map[string]string
)

func initStates() {
	stateMap = make(map[string]string)

	stateMap["Alabama"] = "AL"
	stateMap["Alaska"] = "AK"
	stateMap["Arizona"] = "AZ"
	stateMap["Arkansas"] = "AR"
	stateMap["California"] = "CA"
	stateMap["Colorado"] = "CO"
	stateMap["Connecticut"] = "CT"
	stateMap["Delaware"] = "DE"
	stateMap["District of Columbia"] = "DC"
	stateMap["Florida"] = "FL"
	stateMap["Georgia"] = "GA"
	stateMap["Hawaii"] = "HI"
	stateMap["Idaho"] = "ID"
	stateMap["Illinois"] = "IL"
	stateMap["Indiana"] = "IN"
	stateMap["Iowa"] = "IA"
	stateMap["Kansas"] = "KS"
	stateMap["Kentucky"] = "KY"
	stateMap["Louisiana"] = "LA"
	stateMap["Maine"] = "ME"
	stateMap["Maryland"] = "MD"
	stateMap["Massachusetts"] = "MA"
	stateMap["Michigan"] = "MI"
	stateMap["Minnesota"] = "MN"
	stateMap["Mississippi"] = "MS"
	stateMap["Missouri"] = "MO"
	stateMap["Montana"] = "MT"
	stateMap["Nebraska"] = "NE"
	stateMap["Nevada"] = "NV"
	stateMap["New Hampshire"] = "NH"
	stateMap["New Jersey"] = "NJ"
	stateMap["New Mexico"] = "NM"
	stateMap["New York"] = "NY"
	stateMap["North Carolina"] = "NC"
	stateMap["North Dakota"] = "ND"
	stateMap["Ohio"] = "OH"
	stateMap["Oklahoma"] = "OK"
	stateMap["Oregon"] = "OR"
	stateMap["Pennsylvania"] = "PA"
	stateMap["Rhode Island"] = "RI"
	stateMap["South Carolina"] = "SC"
	stateMap["South Dakota"] = "SD"
	stateMap["Tennessee"] = "TN"
	stateMap["Texas"] = "TX"
	stateMap["Utah"] = "UT"
	stateMap["Vermont"] = "VT"
	stateMap["Virginia"] = "VA"
	stateMap["Washington"] = "WA"
	stateMap["West Virginia"] = "WV"
	stateMap["Wisconsin"] = "WI"
	stateMap["Wyoming"] = "WY"
}

func setStateCode(incident *Incident) {

	switch incident.State {
	case "Alabama":
		incident.StateCode = "AL"
		break
	case "Alaska":
		incident.StateCode = "AK"
		break
	case "Arizona":
		incident.StateCode = "AZ"
		break
	case "Arkansas":
		incident.StateCode = "AR"
		break
	case "California":
		incident.StateCode = "CA"
		break
	case "Colorado":
		incident.StateCode = "CO"
		break
	case "Connecticut":
		incident.StateCode = "CT"
		break
	case "Delaware":
		incident.StateCode = "DE"
		break
	case "District of Columbia":
		incident.StateCode = "DC"
		break
	case "Florida":
		incident.StateCode = "FL"
		break
	case "Georgia":
		incident.StateCode = "GA"
		break
	case "Hawaii":
		incident.StateCode = "HI"
		break
	case "Idaho":
		incident.StateCode = "ID"
		break
	case "Illinois":
		incident.StateCode = "IL"
		break
	case "Indiana":
		incident.StateCode = "IN"
		break
	case "Iowa":
		incident.StateCode = "IA"
		break
	case "Kansas":
		incident.StateCode = "KS"
		break
	case "Kentucky":
		incident.StateCode = "KY"
		break
	case "Louisiana":
		incident.StateCode = "LA"
		break
	case "Maine":
		incident.StateCode = "ME"
		break
	case "Maryland":
		incident.StateCode = "MD"
		break
	case "Massachusetts":
		incident.StateCode = "MA"
		break
	case "Michigan":
		incident.StateCode = "MI"
		break
	case "Minnesota":
		incident.StateCode = "MN"
		break
	case "Mississippi":
		incident.StateCode = "MS"
		break
	case "Missouri":
		incident.StateCode = "MO"
		break
	case "Montana":
		incident.StateCode = "MT"
		break
	case "Nebraska":
		incident.StateCode = "NE"
		break
	case "Nevada":
		incident.StateCode = "NV"
		break
	case "New Hampshire":
		incident.StateCode = "NH"
		break
	case "New Jersey":
		incident.StateCode = "NJ"
		break
	case "New Mexico":
		incident.StateCode = "NM"
		break
	case "New York":
		incident.StateCode = "NY"
		break
	case "North Carolina":
		incident.StateCode = "NC"
		break
	case "North Dakota":
		incident.StateCode = "ND"
		break
	case "Ohio":
		incident.StateCode = "OH"
		break
	case "Oklahoma":
		incident.StateCode = "OK"
		break
	case "Oregon":
		incident.StateCode = "OR"
		break
	case "Pennsylvania":
		incident.StateCode = "PA"
		break
	case "Rhode Island":
		incident.StateCode = "RI"
		break
	case "South Carolina":
		incident.StateCode = "SC"
		break
	case "South Dakota":
		incident.StateCode = "SD"
		break
	case "Tennessee":
		incident.StateCode = "TN"
		break
	case "Texas":
		incident.StateCode = "TX"
		break
	case "Utah":
		incident.StateCode = "UT"
		break
	case "Vermont":
		incident.StateCode = "VT"
		break
	case "Virginia":
		incident.StateCode = "VA"
		break
	case "Washington":
		incident.StateCode = "WA"
		break
	case "West Virginia":
		incident.StateCode = "WV"
		break
	case "Wisconsin":
		incident.StateCode = "WI"
		break
	case "Wyoming":
		incident.StateCode = "WY"
		break
	}
}
